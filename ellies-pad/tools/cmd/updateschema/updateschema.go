package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/testutil"
	"github.com/trythings/trythings/ellies-pad/api"
	"github.com/trythings/trythings/vendor/github.com/facebookgo/inject"
	"github.com/trythings/trythings/vendor/github.com/facebookgo/startstop"
)

func main() {
	apis := api.NewAPIs()

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

	res := graphql.Do(graphql.Params{
		Schema:        apis.Schema,
		RequestString: testutil.IntrospectionQuery,
	})
	if res.HasErrors() {
		panic(res.Errors)
	}

	b, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile("./ellies-pad/schema.json", b, 0644) // File mode -rw-r--r--
	if err != nil {
		panic(err)
	}
}
