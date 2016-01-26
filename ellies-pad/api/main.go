package api

import (
	"net/http"

	"github.com/graphql-go/handler"
)

func init() {
	h := handler.New(&handler.Config{
		Schema: &Schema,
		Pretty: true,
	})

	http.Handle("/graphql", h)
}
