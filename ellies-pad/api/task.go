package api

import (
	"errors"
	"expvar"
	"fmt"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/relay"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/delay"
)

type User struct {
	ID string `json:"id"`
}

type UserService struct {
	user *User
}

func NewUserService(user *User) *UserService {
	return &UserService{
		user: user,
	}
}

func (s *UserService) Get(ctx context.Context, id string) (*User, error) {
	if id != s.user.ID {
		return nil, fmt.Errorf("could not find user with id %q", id)
	}
	return s.user, nil
}

// Task represents a particular action or piece of work to be completed.
type Task struct {
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	Title       string    `json:"title"`
	Description string    `json:"description" datastore:",noindex"`
	IsArchived  bool      `json:"isArchived"`
}

type TaskService struct {
}

func NewTaskService() *TaskService {
	return &TaskService{}
}

var numGet = expvar.NewInt("api.*TaskService.Get")

func (s *TaskService) Get(ctx context.Context, id string) (*Task, error) {
	numGet.Add(1)

	rootKey := datastore.NewKey(ctx, "Root", "root", 0, nil)
	k := datastore.NewKey(ctx, "Task", id, 0, rootKey)
	var t Task
	err := datastore.Get(ctx, k, &t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

var numGetAll = expvar.NewInt("api.*TaskService.GetAll")

func (s *TaskService) GetAll(ctx context.Context) ([]*Task, error) {
	numGetAll.Add(1)

	var ts []*Task
	_, err := datastore.NewQuery("Task").
		Ancestor(datastore.NewKey(ctx, "Root", "root", 0, nil)).
		// Order("-CreatedAt").
		GetAll(ctx, &ts)
	if err != nil {
		return nil, err
	}
	return ts, nil
}

var numCreate = expvar.NewInt("api.*TaskService.Create")

func (s *TaskService) Create(ctx context.Context, t *Task) error {
	numCreate.Add(1)

	if t.ID != "" {
		return fmt.Errorf("t already has id %q", t.ID)
	}

	if t.CreatedAt.IsZero() {
		t.CreatedAt = time.Now()
	}

	id, _, err := datastore.AllocateIDs(ctx, "Task", nil, 1)
	if err != nil {
		return err
	}
	t.ID = fmt.Sprintf("%x", id)

	rootKey := datastore.NewKey(ctx, "Root", "root", 0, nil)
	k := datastore.NewKey(ctx, "Task", t.ID, 0, rootKey)
	k, err = datastore.Put(ctx, k, t)
	if err != nil {
		return err
	}
	return nil
}

var numUpdate = expvar.NewInt("api.*TaskService.Update")

func (s *TaskService) Update(ctx context.Context, t *Task) error {
	numUpdate.Add(1)

	if t.ID == "" {
		return errors.New("cannot update task with no ID")
	}

	rootKey := datastore.NewKey(ctx, "Root", "root", 0, nil)
	k := datastore.NewKey(ctx, "Task", t.ID, 0, rootKey)
	_, err := datastore.Put(ctx, k, t)
	if err != nil {
		return err
	}

	return nil
}

var taskType *graphql.Object
var userType *graphql.Object

var nodeDefinitions *relay.NodeDefinitions

var Schema graphql.Schema

func init() {
	ts := NewTaskService()

	us := NewUserService(&User{
		ID: "ellie",
	})

	nodeDefinitions = relay.NewNodeDefinitions(relay.NodeDefinitionsConfig{
		IDFetcher: func(ctx context.Context, id string, info graphql.ResolveInfo) (interface{}, error) {
			resolvedID := relay.FromGlobalID(id)
			// relay.FromGlobalID does not return an error if it encounters one.
			// Instead, it just returns a nil pointer.
			if resolvedID == nil {
				return nil, fmt.Errorf("invalid id %q", id)
			}

			switch resolvedID.Type {
			case "Task":
				return ts.Get(ctx, resolvedID.ID)
			case "User":
				return us.Get(ctx, resolvedID.ID)
			}
			return nil, fmt.Errorf("unknown type %q", resolvedID.Type)
		},
		TypeResolve: func(value interface{}, info graphql.ResolveInfo) *graphql.Object {
			switch value.(type) {
			case *Task:
				return taskType
			case *User:
				return userType
			}
			return nil
		},
	})

	taskType = graphql.NewObject(graphql.ObjectConfig{
		Name:        "Task",
		Description: "Task represents a particular action or piece of work to be completed.",
		Fields: graphql.Fields{
			"id": relay.GlobalIDField("Task", nil),
			"createdAt": &graphql.Field{
				Description: "When the task was first added",
				Type:        graphql.String,
			},
			"title": &graphql.Field{
				Description: "A short summary of the task",
				Type:        graphql.String,
			},
			"description": &graphql.Field{
				Description: "A more detailed explanation of the task",
				Type:        graphql.String,
			},
			"isArchived": &graphql.Field{
				Description: "Whether this task requires attention",
				Type:        graphql.Boolean,
			},
		},
		Interfaces: []*graphql.Interface{
			nodeDefinitions.NodeInterface,
		},
	})

	userType = graphql.NewObject(graphql.ObjectConfig{
		Name:        "User",
		Description: "User represents a person who can interact with the app.",
		Fields: graphql.Fields{
			"id": relay.GlobalIDField("User", nil),
			"tasks": &graphql.Field{
				Description: "tasks are all pieces of work that need to be completed for the user.",
				Type:        graphql.NewList(taskType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return ts.GetAll(p.Context)
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
					return us.Get(p.Context, "ellie")
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
		},
		OutputFields: graphql.Fields{
			"task": &graphql.Field{
				Type: taskType,
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
			"viewer": &graphql.Field{
				Type: userType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return us.Get(p.Context, "ellie")
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

			t := &Task{
				Title:       title,
				Description: desc,
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

	archiveTaskMutation := relay.MutationWithClientMutationID(relay.MutationConfig{
		Name: "ArchiveTask",
		InputFields: graphql.InputObjectConfigFieldMap{
			"taskId": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"newIsArchived": &graphql.InputObjectFieldConfig{
				Type: graphql.Boolean,
			},
		},
		OutputFields: graphql.Fields{
			"task": &graphql.Field{
				Type: taskType,
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
			id, ok := inputMap["taskId"].(string)
			if !ok {
				return nil, errors.New("could not cast taskId to string")
			}

			resolvedID := relay.FromGlobalID(id)
			if resolvedID == nil {
				return nil, fmt.Errorf("invalid id %q", id)
			}

			newIsArchived, ok := inputMap["newIsArchived"].(bool)
			if !ok {
				newIsArchived = true
			}

			t, err := ts.Get(ctx, resolvedID.ID)
			if err != nil {
				return nil, err
			}

			t.IsArchived = newIsArchived

			// TODO Move this into the service.
			rootKey := datastore.NewKey(ctx, "Root", "root", 0, nil)
			k := datastore.NewKey(ctx, "Task", t.ID, 0, rootKey)
			k, err = datastore.Put(ctx, k, t)
			if err != nil {
				return nil, err
			}

			return map[string]interface{}{
				"taskId": t.ID,
			}, nil
		},
	})

	mutationType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"addTask":                      addTaskMutation,
			"archiveTask":                  archiveTaskMutation,
			"addCreatedAtToTasksMigration": relay.MutationWithClientMutationID(addCreatedAtToTasksMigrationMutation(ts)),
		},
	})

	var err error
	Schema, err = graphql.NewSchema(graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	})
	if err != nil {
		panic(err)
	}
}

func addCreatedAtToTasksMigrationMutation(ts *TaskService) relay.MutationConfig {
	return relay.MutationConfig{
		Name:         "AddCreatedAtToTasksMigration",
		InputFields:  graphql.InputObjectConfigFieldMap{},
		OutputFields: graphql.Fields{},
		MutateAndGetPayload: func(inputMap map[string]interface{}, info graphql.ResolveInfo, ctx context.Context) (map[string]interface{}, error) {
			err := addCreatedAtToTasksMigration.Call(ctx)
			if err != nil {
				return nil, err
			}

			return map[string]interface{}{}, nil
		},
	}
}

var addCreatedAtToTasksMigration = delay.Func("AddCreatedAtToTasks", func(ctx context.Context, ts *TaskService) error {
	// TODO#Perf: Consider using a cursor and/or a batch update.

	tasks, err := ts.GetAll(ctx)
	if err != nil {
		return err
	}

	for _, t := range tasks {
		if t.CreatedAt.IsZero() {
			t.CreatedAt = time.Now()
			err = ts.Update(ctx, t)
			if err != nil {
				return err
			}
		}
	}

	return nil
})
