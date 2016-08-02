package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/graphql-go/handler"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/cloud/trace"
)

var Tracer *trace.Client

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

		// Tracing.
		span := Tracer.SpanFromRequest(r)
		defer func() {
			err := span.FinishWait()
			if err != nil {
				log.Errorf(ctx, "Failed to trace request: %s", err)
			}
		}()
		ctx = trace.NewContext(ctx, span)

		// Authentication.
		auth := r.Header.Get("Authorization")
		if auth != "" {
			idToken, err := getIDToken(auth)
			if err != nil {
				// TODO#Errors
				log.Errorf(ctx, "%s", err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			gu, err := GetGoogleUser(ctx, idToken)
			if err != nil {
				// TODO#Errors
				log.Errorf(ctx, "%s", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			ctx = NewGoogleUserContext(ctx, gu)
		}

		// Caching.
		ctx = NewCacheContext(ctx)

		h.ContextHandler(ctx, w, r)
	})

	http.Handle("/static/", http.FileServer(http.Dir(".")))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	})
}

func getIDToken(auth string) (string, error) {
	if auth == "" {
		return "", errors.New("expected Authorization header")
	}

	if !strings.HasPrefix(auth, "Bearer ") {
		return "", errors.New("expected bearer auth")
	}

	return strings.TrimPrefix(auth, "Bearer "), nil
}
