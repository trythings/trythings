package api

import (
	"net/http"

	"github.com/graphql-go/handler"
	"google.golang.org/appengine"
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
		h.ContextHandler(ctx, w, r)
	})
}
