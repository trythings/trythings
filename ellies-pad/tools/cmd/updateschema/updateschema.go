package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/testutil"
	"github.com/trythings/trythings/ellies-pad/api"
)

func main() {
	apis, err := api.NewAPIs()
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
