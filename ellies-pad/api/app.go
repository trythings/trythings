package api

import (
	"errors"
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/relay"
	"golang.org/x/net/context"
)

var nodeDefinitions *relay.NodeDefinitions

var Schema graphql.Schema

func init() {
	us := NewUserService()
	ss := NewSpaceService(us)
	ts := NewTaskService(ss)
	ms := NewMigrationService(ss, ts)

	var spaceAPI *SpaceAPI
	var taskAPI *TaskAPI
	var userAPI *UserAPI

	nodeDefinitions = relay.NewNodeDefinitions(relay.NodeDefinitionsConfig{
		IDFetcher: func(ctx context.Context, id string, info graphql.ResolveInfo) (interface{}, error) {
			resolvedID := relay.FromGlobalID(id)
			// relay.FromGlobalID does not return an error if it encounters one.
			// Instead, it just returns a nil pointer.
			if resolvedID == nil {
				return nil, fmt.Errorf("invalid id %q", id)
			}

			switch resolvedID.Type {
			case "Space":
				return nil, errors.New("not implemented")
			case "Task":
				return ts.Get(ctx, resolvedID.ID)
			case "User":
				return nil, errors.New("not implemented")
			}
			return nil, fmt.Errorf("unknown type %q", resolvedID.Type)
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

	taskAPI = &TaskAPI{
		tasks:           ts,
		nodeDefinitions: nodeDefinitions,
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

	queryType := graphql.NewObject(graphql.ObjectConfig{
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

	addTaskMutation := relay.MutationWithClientMutationID(relay.MutationConfig{
		Name: "AddTask",
		InputFields: graphql.InputObjectConfigFieldMap{
			"title": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"description": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"spaceId": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
		OutputFields: graphql.Fields{
			"task": &graphql.Field{
				Type: taskAPI.Type(),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					payload, ok := p.Source.(map[string]interface{})
					if !ok {
						return nil, errors.New("could not cast payload to map")
					}
					id, ok := payload["taskId"].(string)
					if !ok {
						return nil, errors.New("could not cast taskId to string")
					}
					t, err := ts.Get(p.Context, id)
					if err != nil {
						return nil, err
					}
					return t, nil
				},
			},
		},
		MutateAndGetPayload: func(inputMap map[string]interface{}, info graphql.ResolveInfo, ctx context.Context) (map[string]interface{}, error) {
			title, ok := inputMap["title"].(string)
			if !ok {
				return nil, errors.New("could not cast title to string")
			}

			var desc string
			descOrNil := inputMap["description"]
			if descOrNil != nil {
				desc, ok = descOrNil.(string)
				if !ok {
					return nil, errors.New("could not cast description to string")
				}
			}

			spaceID, ok := inputMap["spaceId"].(string)
			if !ok {
				return nil, errors.New("could not cast spaceId to string")
			}
			resolvedSpaceID := relay.FromGlobalID(spaceID)
			if resolvedSpaceID == nil {
				return nil, fmt.Errorf("invalid id %q", spaceID)
			}

			t := &Task{
				Title:       title,
				Description: desc,
				SpaceID:     resolvedSpaceID.ID,
			}
			err := ts.Create(ctx, t)
			if err != nil {
				return nil, err
			}

			return map[string]interface{}{
				"taskId": t.ID,
			}, nil
		},
	})

	editTaskMutation := relay.MutationWithClientMutationID(relay.MutationConfig{
		Name: "EditTask",
		InputFields: graphql.InputObjectConfigFieldMap{
			"id": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"title": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"description": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"isArchived": &graphql.InputObjectFieldConfig{
				Type: graphql.Boolean,
			},
		},
		OutputFields: graphql.Fields{
			"task": &graphql.Field{
				Type: taskAPI.Type(),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					payload, ok := p.Source.(map[string]interface{})
					if !ok {
						return nil, errors.New("could not cast payload to map")
					}
					id, ok := payload["id"].(string)
					if !ok {
						return nil, errors.New("could not cast id to string")
					}
					t, err := ts.Get(p.Context, id)
					if err != nil {
						return nil, err
					}
					return t, nil
				},
			},
		},
		MutateAndGetPayload: func(inputMap map[string]interface{}, info graphql.ResolveInfo, ctx context.Context) (map[string]interface{}, error) {
			id, ok := inputMap["id"].(string)
			if !ok {
				return nil, errors.New("could not cast id to string")
			}

			resolvedID := relay.FromGlobalID(id)
			if resolvedID == nil {
				return nil, fmt.Errorf("invalid id %q", id)
			}

			t, err := ts.Get(ctx, resolvedID.ID)
			if err != nil {
				return nil, err
			}

			title, ok := inputMap["title"].(string)
			if ok {
				t.Title = title
			}

			description, ok := inputMap["description"].(string)
			if ok {
				t.Description = description
			}

			isArchived, ok := inputMap["isArchived"].(bool)
			if ok {
				t.IsArchived = isArchived
			}

			err = ts.Update(ctx, t)
			if err != nil {
				return nil, err
			}

			return map[string]interface{}{
				"id": t.ID,
			}, nil
		},
	})

	mutationType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"addTask":  addTaskMutation,
			"editTask": editTaskMutation,
		},
	})

	migrationAPI := &MigrationAPI{
		migrations: ms,
	}
	migrationAPI.AddMutations(mutationType)

	var err error
	Schema, err = graphql.NewSchema(graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	})
	if err != nil {
		panic(err)
	}
}

type key int

const superuserKey key = 0

func isSuperuser(ctx context.Context) (bool, error) {
	v := ctx.Value(superuserKey)
	if v == nil {
		return false, nil
	}

	su, ok := v.(bool)
	if !ok {
		return false, errors.New("unexpected superuser type")
	}

	return su, nil
}
