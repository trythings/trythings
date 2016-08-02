package main

import (
	"github.com/trythings/trythings/ellies-pad/api"
	"google.golang.org/appengine"
	"google.golang.org/cloud/trace"
)

func main() {
	var err error
	api.Tracer, err = trace.NewClient(appengine.BackgroundContext(), "ellies-pad")
	if err != nil {
		panic(err)
	}

	appengine.Main()
}
