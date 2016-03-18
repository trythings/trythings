package api

import (
	"errors"

	"github.com/facebookgo/inject"
	"github.com/facebookgo/startstop"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/relay"
	"golang.org/x/net/context"
)

type apis struct {
	MigrationAPI *MigrationAPI `inject:""`
	SpaceAPI     *SpaceAPI     `inject:""`
	TaskAPI      *TaskAPI      `inject:""`
	UserAPI      *UserAPI      `inject:""`
	UserService  *UserService  `inject:""`
}

var Schema graphql.Schema

func init() {
	apis := &apis{}

	nodeDefinitions := relay.NewNodeDefinitions(relay.NodeDefinitionsConfig{
		IDFetcher: func(ctx context.Context, id string, info graphql.ResolveInfo) (interface{}, error) {
			return nil, errors.New("not implemented")
		},
		TypeResolve: func(value interface{}, info graphql.ResolveInfo) *graphql.Object {
			switch value.(type) {
			case *Space:
				return apis.SpaceAPI.Type()
			case *Task:
				return apis.TaskAPI.Type()
			case *User:
				return apis.UserAPI.Type()
			}
			return nil
		},
	})

	mutation := graphql.NewObject(graphql.ObjectConfig{
		Name:   "Mutation",
		Fields: graphql.Fields{},
	})

	graph := &inject.Graph{}
	err := graph.Provide(
		&inject.Object{
			Value: apis,
		},
		&inject.Object{
			Value: nodeDefinitions.NodeInterface,
			Name:  "node",
		},
		&inject.Object{
			Value: mutation,
			Name:  "mutation",
		},
	)
	if err != nil {
		panic(err)
	}

	err = graph.Populate()
	if err != nil {
		panic(err)
	}

	startstop.Start(graph.Objects(), nil)

	query := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"node": nodeDefinitions.NodeField,
			"viewer": &graphql.Field{
				Description: "viewer is the person currently interacting with the app.",
				Type:        apis.UserAPI.Type(),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					// gu := user.Current(p.Context)
					// err := apis.UserService.Create(p.Context, &User{
					// 	GoogleID: gu.ID,
					// 	Email:    gu.Email,
					// 	Name:     gu.String(),
					// })
					// if err != nil {
					// 	return nil, err
					// }
					return apis.UserService.FromContext(p.Context)
				},
			},
		},
	})

	Schema, err = graphql.NewSchema(graphql.SchemaConfig{
		Query:    query,
		Mutation: mutation,
	})
	if err != nil {
		panic(err)
	}
}
