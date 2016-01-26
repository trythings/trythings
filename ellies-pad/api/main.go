package api

import (
	"fmt"
	"net/http"

	"github.com/graphql-go/handler"
)

func init() {
	h := handler.New(&handler.Config{
		Schema: &Schema,
		Pretty: true,
	})

	http.Handle("/graphql", h)
	http.HandleFunc("/hello", handlerFunc)
}

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world!")
}
