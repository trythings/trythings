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

	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "Accept")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
		h.ServeHTTP(w, r)
	})
}
