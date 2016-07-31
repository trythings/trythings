package api

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/urlfetch"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

type GoogleUser struct {
	ID            string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
}

var googleKeys jose.JSONWebKeySet

func updateGoogleKeys(ctx context.Context) error {
	// Try to fetch new public keys from Google.
	client := urlfetch.Client(ctx)
	client.Timeout = 1 * time.Second
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/certs")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&googleKeys)
	if err != nil {
		return err
	}

	return nil
}

func GetGoogleUser(ctx context.Context, idToken string) (*GoogleUser, error) {
	tok, err := jwt.ParseSigned(idToken)
	if err != nil {
		return nil, err
	}

	if len(tok.Headers) != 1 {
		// We must have a header to specify a kid.
		// We don't know how to handle multiple headers,
		// since it's unclear which kid to use.
		return nil, errors.New("expected exactly one token header")
	}

	keys := googleKeys.Key(tok.Headers[0].KeyID)
	if len(keys) == 0 {
		err := updateGoogleKeys(ctx)
		if err != nil {
			return nil, err
		}
		keys = googleKeys.Key(tok.Headers[0].KeyID)
	}

	if len(keys) != 1 {
		// We must have a key to check the signature.
		// We don't know how to deal with multiple keys matching the same kid.
		return nil, errors.New("expected exactly one key matching kid")
	}
	key := keys[0]

	var payload struct {
		jwt.Claims
		GoogleUser
	}
	err = tok.Claims(&payload, key.Key)
	if err != nil {
		return nil, err
	}

	expectedIssuer := "accounts.google.com"
	if strings.HasPrefix(payload.Issuer, "https://") {
		expectedIssuer = "https://accounts.google.com"
	}

	err = payload.Validate(jwt.Expected{
		Issuer:   expectedIssuer,
		Audience: []string{"695504958192-8k3tf807271m7jcllcvlauddeqhbr0hg.apps.googleusercontent.com"},
		Time:     time.Now(),
	})
	if err != nil {
		return nil, err
	}

	return &payload.GoogleUser, nil
}

func NewGoogleUserContext(ctx context.Context, gu *GoogleUser) context.Context {
	return context.WithValue(ctx, googleUserKey, gu)
}

func GoogleUserFromContext(ctx context.Context) (*GoogleUser, bool) {
	gu, ok := ctx.Value(googleUserKey).(*GoogleUser)
	return gu, ok
}
