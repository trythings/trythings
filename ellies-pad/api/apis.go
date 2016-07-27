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

	SearchAPI *SearchAPI `inject:""`
	SpaceAPI  *SpaceAPI  `inject:""`
	TaskAPI   *TaskAPI   `inject:""`
	UserAPI   *UserAPI   `inject:""`

	UserService *UserService `inject:""`

	Schema          *graphql.Schema
	nodeDefinitions *relay.NodeDefinitions
}

func NewAPIs() (*apis, error) {
	apis := &apis{}
	apis.nodeDefinitions = relay.NewNodeDefinitions(relay.NodeDefinitionsConfig{
		IDFetcher: func(ctx context.Context, id string, info graphql.ResolveInfo) (interface{}, error) {
			resolvedID := relay.FromGlobalID(id)
			switch resolvedID.Type {
			case "Search":
				return apis.SearchAPI.SearchService.ByID(ctx, resolvedID.ID)
			default:
				return nil, errors.New("Unknown node type")
			}
		},
		TypeResolve: func(p graphql.ResolveTypeParams) *graphql.Object {
			switch p.Value.(type) {
			case *Search:
				return apis.SearchAPI.Type
			case *Space:
				return apis.SpaceAPI.Type
			case *Task:
				return apis.TaskAPI.Type
			case *User:
				return apis.UserAPI.Type
			}
			return nil
		},
	})

	graph := &inject.Graph{}
	err := graph.Provide(
		&inject.Object{
			Value: apis,
		},
		&inject.Object{
			Value: apis.nodeDefinitions.NodeInterface,
			Name:  "node",
		},
	)
	if err != nil {
		return nil, err
	}

	err = graph.Populate()
	if err != nil {
		return nil, err
	}

	err = startstop.Start(graph.Objects(), nil)
	if err != nil {
		return nil, err
	}

	return apis, nil
}

func (apis *apis) Start() error {
	query := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"node": apis.nodeDefinitions.NodeField,
			"viewer": &graphql.Field{
				Description: "viewer is the person currently interacting with the app.",
				Type:        apis.UserAPI.Type,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					u, err := apis.UserService.FromContext(p.Context)
					if err == ErrUserNotFound {
						gu, ok := GoogleUserFromContext(p.Context)
						if !ok {
							return nil, errors.New("expected google user, probably missing Authorization header")
						}
						// TODO Some of the google user's fields could change after user creation.
						// Consider updating the user to reflect those changes (e.g. IsEmailVerified).
						err := apis.UserService.Create(p.Context, &User{
							GoogleID:        gu.ID,
							Email:           gu.Email,
							IsEmailVerified: gu.EmailVerified,
							Name:            gu.Name,
							GivenName:       gu.GivenName,
							FamilyName:      gu.FamilyName,
							ImageURL:        gu.Picture,
						})
						if err != nil {
							return nil, err
						}
						// TODO create default space and view
					}
					return u, err
				},
			},
		},
	})

	mutation := graphql.NewObject(graphql.ObjectConfig{
		Name:   "Mutation",
		Fields: graphql.Fields{},
	})

	for n, f := range apis.MigrationAPI.Mutations {
		mutation.AddFieldConfig(n, f)
	}
	for n, f := range apis.TaskAPI.Mutations {
		mutation.AddFieldConfig(n, f)
	}

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    query,
		Mutation: mutation,
	})

	if err != nil {
		return err
	}

	apis.Schema = &schema

	return nil
}
