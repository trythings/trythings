package api

import (
	"errors"
	"fmt"

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

	SearchService *SearchService `inject:""`
	SpaceService  *SpaceService  `inject:""`
	TaskService   *TaskService   `inject:""`
	UserService   *UserService   `inject:""`
	ViewService   *ViewService   `inject:""`

	Schema          *graphql.Schema
	nodeDefinitions *relay.NodeDefinitions
}

func NewAPIs() (*apis, error) {
	apis := &apis{}
	apis.nodeDefinitions = relay.NewNodeDefinitions(relay.NodeDefinitionsConfig{
		IDFetcher: func(id string, info graphql.ResolveInfo, ctx context.Context) (interface{}, error) {
			resolvedID := relay.FromGlobalID(id)
			switch resolvedID.Type {
			case "Search":
				return apis.SearchService.ByClientID(ctx, resolvedID.ID)
			case "Space":
				return apis.SpaceService.ByID(ctx, resolvedID.ID)
			case "User":
				return apis.UserService.ByID(ctx, resolvedID.ID)
			default:
				return nil, fmt.Errorf("Unknown node type %s", resolvedID.Type)
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
						u := &User{
							GoogleID:        gu.ID,
							Email:           gu.Email,
							IsEmailVerified: gu.EmailVerified,
							Name:            gu.Name,
							GivenName:       gu.GivenName,
							FamilyName:      gu.FamilyName,
							ImageURL:        gu.Picture,
						}
						err := apis.UserService.Create(p.Context, u)
						if err != nil {
							return nil, err
						}

						rt, err := apis.TaskService.GetOrCreateRootTask(ctx, u)
						if err != nil {
							return nil, err
						}

						se := &Search{
							Name:         "#welcome Search",
							ParentTaskID: rt.ID,
							Query:        "#welcome",
						}
						err = apis.SearchService.Create(p.Context, se)
						if err != nil {
							return nil, err
						}

						se = &Search{
							Name:         "Recent Search",
							ParentTaskID: rt.ID,
							Query:        "CreatedAt >= today() ",
						}
						err = apis.SearchService.Create(p.Context, se)
						if err != nil {
							return nil, err
						}

						se = &Search{
							Name:         "Everything Else Search",
							ParentTaskID: rt.ID,
							Query:        "CreatedAt < today() ",
						}
						err = apis.SearchService.Create(p.Context, se)
						if err != nil {
							return nil, err
						}

						var t *Task

						t = &Task{
							Title:   "Recently changed or added tasks will show up in 'Recent'",
							SpaceID: sp.ID,
						}
						err = apis.TaskService.Create(p.Context, t)
						if err != nil {
							return nil, err
						}

						t = &Task{
							Title:       "Searches help you find and organize tasks",
							Description: "#welcome",
							SpaceID:     sp.ID,
						}
						err = apis.TaskService.Create(p.Context, t)
						if err != nil {
							return nil, err
						}

						t = &Task{
							Title:       "The same task can show up in multiple searches",
							Description: "#welcome",
							SpaceID:     sp.ID,
						}
						err = apis.TaskService.Create(p.Context, t)
						if err != nil {
							return nil, err
						}

						t = &Task{
							Title:       fmt.Sprintf("Tasks in '%s' are only visible to you", sp.Name),
							Description: "You can also add people to your space to share it. #welcome",
							SpaceID:     sp.ID,
						}
						err = apis.TaskService.Create(p.Context, t)
						if err != nil {
							return nil, err
						}

						t = &Task{
							Title:       "Add #tags to help you find tasks later",
							Description: "#welcome",
							SpaceID:     sp.ID,
						}
						err = apis.TaskService.Create(p.Context, t)
						if err != nil {
							return nil, err
						}

						v = &View{
							SpaceID: sp.ID,
							Name:    "Priority View",
						}
						err = apis.ViewService.Create(p.Context, v)
						if err != nil {
							return nil, err
						}

						se = &Search{
							Name:   "Unprioritized Search",
							ViewID: v.ID,
							Query:  "NOT (#now OR #next OR #later OR IsArchived: true)",
						}
						err = apis.SearchService.Create(p.Context, se)
						if err != nil {
							return nil, err
						}

						se = &Search{
							Name:   "#now Search",
							ViewID: v.ID,
							Query:  "#now",
						}
						err = apis.SearchService.Create(p.Context, se)
						if err != nil {
							return nil, err
						}

						se = &Search{
							Name:   "#next Search",
							ViewID: v.ID,
							Query:  "#next AND NOT #now",
						}
						err = apis.SearchService.Create(p.Context, se)
						if err != nil {
							return nil, err
						}

						se = &Search{
							Name:   "#later Search",
							ViewID: v.ID,
							Query:  "#later AND NOT (#now OR #next)",
						}
						err = apis.SearchService.Create(p.Context, se)
						if err != nil {
							return nil, err
						}

						t = &Task{
							Title:       "Use #now, #next, and #later to prioritize tasks",
							Description: "https://medium.com/@noah_weiss/now-next-later-roadmaps-without-the-drudgery-1cfe65656645#.lcwurwozj",
							SpaceID:     sp.ID,
						}
						err = apis.TaskService.Create(p.Context, t)
						if err != nil {
							return nil, err
						}

						t = &Task{
							Title:       "#now is the next 2–4 weeks",
							Description: "For many people that use bi-weekly sprints, this fits perfectly into their planning cadence.",
							SpaceID:     sp.ID,
						}
						err = apis.TaskService.Create(p.Context, t)
						if err != nil {
							return nil, err
						}

						t = &Task{
							Title:       "#next is 1-3 months out",
							Description: "Effectively, it’s the rest of the quarter after now.",
							SpaceID:     sp.ID,
						}
						err = apis.TaskService.Create(p.Context, t)
						if err != nil {
							return nil, err
						}

						t = &Task{
							Title:       "#later is 3+ months out",
							Description: "It’s a useful place to park ideas you're passionate about.",
							SpaceID:     sp.ID,
						}
						err = apis.TaskService.Create(p.Context, t)
						if err != nil {
							return nil, err
						}

						return u, nil
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
