package api

import (
	"net/http"

	"github.com/facebookgo/inject"
	"github.com/facebookgo/startstop"
	"github.com/graphql-go/handler"
	"google.golang.org/appengine"
)

func init() {
	apis := NewAPIs()

	graph := &inject.Graph{}
	err := graph.Provide(
		&inject.Object{
			Value: apis,
		},
		&inject.Object{
			Value: apis.NodeDefinitions.NodeInterface,
			Name:  "node",
		},
	)
	if err != nil {
		panic(err)
	}

	err = graph.Populate()
	if err != nil {
		panic(err)
	}

	err = startstop.Start(graph.Objects(), nil)
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
