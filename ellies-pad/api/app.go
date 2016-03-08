package api

import (
	"errors"
	"fmt"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/relay"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/user"
)

type User struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	GoogleID  string    `json:"-"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
}

type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) Create(ctx context.Context, u *User) error {
	// TODO Make sure u.GoogleID == user.Current(ctx).ID

	if u.ID != "" {
		return fmt.Errorf("u already has id %q", u.ID)
	}

	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now()
	}

	id, _, err := datastore.AllocateIDs(ctx, "User", nil, 1)
	if err != nil {
		return err
	}
	u.ID = fmt.Sprintf("%x", id)

	rootKey := datastore.NewKey(ctx, "Root", "root", 0, nil)
	k := datastore.NewKey(ctx, "User", u.ID, 0, rootKey)
	k, err = datastore.Put(ctx, k, u)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) byGoogleID(ctx context.Context, googleID string) (*User, error) {
	var us []*User
	_, err := datastore.NewQuery("User").
		Ancestor(datastore.NewKey(ctx, "Root", "root", 0, nil)).
		Filter("GoogleID =", googleID).
		Limit(1).
		GetAll(ctx, &us)
	if err != nil {
		return nil, err
	}

	if len(us) == 0 {
		return nil, errors.New("user not found")
	}

	return us[0], nil
}

// FromContext should not be subject to access control,
// because it would create a circular dependency.
func (s *UserService) FromContext(ctx context.Context) (*User, error) {
	gu := user.Current(ctx)
	return s.byGoogleID(ctx, gu.ID)
}

type Space struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Name      string    `json:"name"`
	UserIDs   []string  `json:"userIds"`
}

type SpaceService struct {
	users *UserService
}

func NewSpaceService(users *UserService) *SpaceService {
	return &SpaceService{
		users: users,
	}
}

func (s *SpaceService) ByID(ctx context.Context, id string) (*Space, error) {
	rootKey := datastore.NewKey(ctx, "Root", "root", 0, nil)
	k := datastore.NewKey(ctx, "Space", id, 0, rootKey)
	var sp Space
	err := datastore.Get(ctx, k, &sp)
	if err != nil {
		return nil, err
	}

	ok, err := s.IsVisible(ctx, &sp)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("cannot access space")
	}

	return &sp, nil
}

func (s *SpaceService) Create(ctx context.Context, sp *Space) error {
	if sp.ID != "" {
		return fmt.Errorf("sp already has id %q", sp.ID)
	}

	if sp.CreatedAt.IsZero() {
		sp.CreatedAt = time.Now()
	}

	if len(sp.UserIDs) > 0 {
		return errors.New("UserIDs must be empty")
	}

	su, err := isSuperuser(ctx)
	if err != nil {
		return err
	}
	if !su {
		u, err := s.users.FromContext(ctx)
		if err != nil {
			return err
		}
		sp.UserIDs = []string{u.ID}
	}

	id, _, err := datastore.AllocateIDs(ctx, "Space", nil, 1)
	if err != nil {
		return err
	}
	sp.ID = fmt.Sprintf("%x", id)

	rootKey := datastore.NewKey(ctx, "Root", "root", 0, nil)
	k := datastore.NewKey(ctx, "Space", sp.ID, 0, rootKey)
	k, err = datastore.Put(ctx, k, sp)
	if err != nil {
		return err
	}

	return nil
}

func (s *SpaceService) IsVisible(ctx context.Context, sp *Space) (bool, error) {
	su, err := isSuperuser(ctx)
	if err != nil {
		return false, err
	}

	if su {
		return true, nil
	}

	u, err := s.users.FromContext(ctx)
	if err != nil {
		return false, err
	}

	for _, id := range sp.UserIDs {
		if u.ID == id {
			return true, nil
		}
	}

	return false, nil
}

var spaceType *graphql.Object
var userType *graphql.Object

var nodeDefinitions *relay.NodeDefinitions

var Schema graphql.Schema

func init() {
	us := NewUserService()
	ss := NewSpaceService(us)
	ts := NewTaskService(ss)
	ms := NewMigrationService(ss, ts)

	var taskAPI *TaskAPI

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
				return spaceType
			case *Task:
				return taskAPI.Type()
			case *User:
				return userType
			}
			return nil
		},
	})

	taskAPI = &TaskAPI{
		tasks:           ts,
		nodeDefinitions: nodeDefinitions,
	}
	taskAPI.Start()

	spaceType = graphql.NewObject(graphql.ObjectConfig{
		Name:        "Space",
		Description: "Space represents an access-controlled universe of tasks.",
		Fields: graphql.Fields{
			"id": relay.GlobalIDField("Space", nil),
			"createdAt": &graphql.Field{
				Description: "When the space was first created.",
				Type:        graphql.String,
			},
			"name": &graphql.Field{
				Description: "The name to display for the space.",
				Type:        graphql.String,
			},
			"tasks": &graphql.Field{
				Args: graphql.FieldConfigArgument{
					"query": &graphql.ArgumentConfig{
						Type:         graphql.String,
						DefaultValue: "",
						Description:  "query filters the result to only tasks that contain particular terms in their title or description",
					},
				},
				Description: "tasks are all pieces of work that need to be completed for the user.",
				Type:        graphql.NewList(taskAPI.Type()),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					sp, ok := p.Source.(*Space)
					if !ok {
						return nil, errors.New("expected a space source")
					}

					q, ok := p.Args["query"].(string)
					if !ok {
						q = "" // Return all tasks.
					}

					return ts.Search(p.Context, sp, q)
				},
			},
		},
	})

	userType = graphql.NewObject(graphql.ObjectConfig{
		Name:        "User",
		Description: "User represents a person who can interact with the app.",
		Fields: graphql.Fields{
			"id": relay.GlobalIDField("User", nil),
			"email": &graphql.Field{
				Description: "The user's email primary address",
				Type:        graphql.String,
			},
			"space": &graphql.Field{
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type:         graphql.String,
						DefaultValue: "",
						Description:  "id can be omitted, which will have space resolve to the user's default space.",
					},
				},
				Description: "space is a disjoint universe of views, searches and tasks.",
				Type:        spaceType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, ok := p.Args["id"].(string)
					if ok {
						resolvedID := relay.FromGlobalID(id)
						if resolvedID == nil {
							return nil, fmt.Errorf("invalid id %q", id)
						}

						sp, err := ss.ByID(p.Context, resolvedID.ID)
						if err != nil {
							return nil, err
						}
						return sp, nil
					}

					u, ok := p.Source.(*User)
					if !ok {
						return nil, errors.New("expected user source")
					}

					var sps []*Space
					_, err := datastore.NewQuery("Space").
						Ancestor(datastore.NewKey(p.Context, "Root", "root", 0, nil)).
						Filter("UserIDs =", u.ID).
						Limit(1).
						GetAll(p.Context, &sps)
					if err != nil {
						return nil, err
					}

					if len(sps) == 0 {
						return nil, errors.New("could not find default space for user")
					}

					return sps[0], nil
				},
			},
			"spaces": &graphql.Field{
				Type: graphql.NewList(spaceType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					u, ok := p.Source.(*User)
					if !ok {
						return nil, errors.New("expected user source")
					}

					var sps []*Space
					_, err := datastore.NewQuery("Space").
						Ancestor(datastore.NewKey(p.Context, "Root", "root", 0, nil)).
						Filter("UserIDs =", u.ID).
						GetAll(p.Context, &sps)
					if err != nil {
						return nil, err
					}
					return sps, nil
				},
			},
		},
		Interfaces: []*graphql.Interface{
			nodeDefinitions.NodeInterface,
		},
	})

	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"node": nodeDefinitions.NodeField,
			"viewer": &graphql.Field{
				Description: "viewer is the person currently interacting with the app.",
				Type:        userType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
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
