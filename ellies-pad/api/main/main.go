package main

import (
	"errors"
	"net/http"
	"strings"

	"github.com/graphql-go/handler"
	"github.com/trythings/trythings/ellies-pad/api"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func main() {
	apis, err := api.NewAPIs()
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
			idToken, err := getIDToken(auth)
			if err != nil {
				// TODO#Errors
				log.Errorf(ctx, "%s", err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			gu, err := api.GetGoogleUser(ctx, idToken)
			if err != nil {
				// TODO#Errors
				log.Errorf(ctx, "%s", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			ctx = api.NewGoogleUserContext(ctx, gu)
		}

		ctx = NewCacheContext(ctx)

		h.ContextHandler(ctx, w, r)
	})

	http.Handle("/static/", http.FileServer(http.Dir(".")))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	})

	appengine.Main()
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
