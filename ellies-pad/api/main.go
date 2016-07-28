package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"golang.org/x/net/context"

	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"

	"github.com/graphql-go/handler"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

func init() {
	apis, err := NewAPIs()
	if err != nil {
		panic(err)
	}

	h := handler.New(&handler.Config{
		Schema: apis.Schema,
		Pretty: true,
	})

	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		ctx := appengine.NewContext(r)
		auth := r.Header.Get("Authorization")
		if auth != "" {
			u, err := getUser(ctx, r.Header.Get("Authorization"))
			if err != nil {
				log.Errorf(ctx, "%s", err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			ctx = NewUserContext(ctx, u)
		}
		h.ContextHandler(ctx, w, r)
	})
}

var googleKeys jose.JSONWebKeySet

func getUser(ctx context.Context, auth string) (*User, error) {
	if auth == "" {
		return nil, errors.New("expected Authorization header")
	}

	if !strings.HasPrefix(auth, "Bearer ") {
		return nil, errors.New("expected bearer auth")
	}

	auth = strings.TrimPrefix(auth, "Bearer ")
	t, err := jwt.ParseSigned(auth)
	if err != nil {
		return nil, err
	}

	if len(t.Headers) != 1 {
		return nil, errors.New("expected exactly one token header")
	}

	keys := googleKeys.Key(t.Headers[0].KeyID)
	if len(keys) == 0 {
		// Try to fetch new public keys from Google.
		client := urlfetch.Client(ctx)
		client.Timeout = 1 * time.Second
		resp, err := client.Get("https://www.googleapis.com/oauth2/v3/certs")
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(&googleKeys)
		if err != nil {
			return nil, err
		}

		keys = googleKeys.Key(t.Headers[0].KeyID)
		if len(keys) == 0 {
			return nil, errors.New("could not find key matching kid")
		}
	}

	gu := struct {
		jwt.Claims
		ID            string `json:"sub"`
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
		Name          string `json:"name"`
		GivenName     string `json:"given_name"`
		FamilyName    string `json:"family_name"`
		Picture       string `json:"picture"`
	}{}
	err = t.Claims(&gu, keys[0].Key)
	if err != nil {
		return nil, err
	}

	issuer := "accounts.google.com"
	if strings.HasPrefix(gu.Issuer, "https://") {
		issuer = "https://accounts.google.com"
	}

	err = gu.Validate(jwt.Expected{
		Issuer:   issuer,
		Audience: []string{"695504958192-8k3tf807271m7jcllcvlauddeqhbr0hg.apps.googleusercontent.com"},
		Time:     time.Now(),
	})
	if err != nil {
		return nil, err
	}

	return &User{
		GoogleID:        gu.ID,
		Email:           gu.Email,
		IsEmailVerified: gu.EmailVerified,
		Name:            gu.Name,
		GivenName:       gu.GivenName,
		FamilyName:      gu.FamilyName,
		ImageURL:        gu.Picture,
	}, nil
}
