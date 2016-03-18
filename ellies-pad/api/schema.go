package api

import (
	"errors"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/relay"
	"golang.org/x/net/context"
)

var Schema graphql.Schema

func init() {
	us := NewUserService()
	ss := NewSpaceService(us)
	ts := NewTaskService(ss)
	ms := NewMigrationService(ss, ts)

	var spaceAPI *SpaceAPI
	var taskAPI *TaskAPI
	var userAPI *UserAPI

	nodeDefinitions := relay.NewNodeDefinitions(relay.NodeDefinitionsConfig{
		IDFetcher: func(ctx context.Context, id string, info graphql.ResolveInfo) (interface{}, error) {
			return nil, errors.New("not implemented")
		},
		TypeResolve: func(value interface{}, info graphql.ResolveInfo) *graphql.Object {
			switch value.(type) {
			case *Space:
				return spaceAPI.Type()
			case *Task:
				return taskAPI.Type()
			case *User:
				return userAPI.Type()
			}
			return nil
		},
	})

	mutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
	})

	migrationAPI := &MigrationAPI{
		migrations: ms,
		mutation:   mutation,
	}
	migrationAPI.Start()

	taskAPI = &TaskAPI{
		tasks:           ts,
		nodeDefinitions: nodeDefinitions,
		mutation:        mutation,
	}
	taskAPI.Start()

	spaceAPI = &SpaceAPI{
		spaces:          ss,
		tasks:           ts,
		taskAPI:         taskAPI,
		nodeDefinitions: nodeDefinitions,
	}
	spaceAPI.Start()

	userAPI = &UserAPI{
		users:           us,
		spaces:          ss,
		spaceAPI:        spaceAPI,
		nodeDefinitions: nodeDefinitions,
	}
	userAPI.Start()

	query := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"node": nodeDefinitions.NodeField,
			"viewer": &graphql.Field{
				Description: "viewer is the person currently interacting with the app.",
				Type:        userAPI.Type(),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					// gu := user.Current(p.Context)
					// err := us.Create(p.Context, &User{
					// 	GoogleID: gu.ID,
					// 	Email:    gu.Email,
					// 	Name:     gu.String(),
					// })
					// if err != nil {
					// 	return nil, err
					// }
					return us.FromContext(p.Context)
				},
			},
		},
	})

	var err error
	Schema, err = graphql.NewSchema(graphql.SchemaConfig{
		Query:    query,
		Mutation: mutation,
	})
	if err != nil {
		panic(err)
	}
}
