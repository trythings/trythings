package api

import (
	"errors"
	"fmt"

	"golang.org/x/net/context"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/relay"
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

func (s *UserService) Get(id string) *User {
	if id != s.user.ID {
		return nil
	}
	return s.user
}

// Task represents a particular action or piece of work to be completed.
type Task struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type TaskService struct {
	tasks []*Task
}

func NewTaskService(tasks []*Task) *TaskService {
	return &TaskService{
		tasks: tasks,
	}
}

func (s *TaskService) Get(id string) (*Task, error) {
	for _, t := range s.tasks {
		if t.ID == id {
			return t, nil
		}
	}
	return nil, fmt.Errorf("could not find task with id %q", id)
}

func (s *TaskService) GetAll() []*Task {
	return s.tasks
}

func (s *TaskService) Create(t *Task) error {
	if t.ID != "" {
		return fmt.Errorf("t already has id %q", t.ID)
	}
	t.ID = "klm"
	s.tasks = append(s.tasks, t)
	return nil
}

var taskType *graphql.Object
var userType *graphql.Object

var nodeDefinitions *relay.NodeDefinitions

var Schema graphql.Schema

func init() {
	ts := NewTaskService([]*Task{
		&Task{
			ID:    "abc",
			Title: "Pick up milk",
		},
		&Task{
			ID:    "def",
			Title: "Finish working on Ellie's Pad",
		},
		&Task{
			ID:    "ghi",
			Title: "Rub the dog",
		},
	})

	us := NewUserService(&User{
		ID: "ellie",
	})

	nodeDefinitions = relay.NewNodeDefinitions(relay.NodeDefinitionsConfig{
		IDFetcher: func(id string, info graphql.ResolveInfo) interface{} {
			resolvedID := relay.FromGlobalID(id)
			switch resolvedID.Type {
			case "Task":
				t, err := ts.Get(resolvedID.ID)
				if err != nil {
					panic(err)
				}
				return t
			case "User":
				return us.Get(resolvedID.ID)
			}
			return nil
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
			"title": &graphql.Field{
				Description: "A short summary of the task",
				Type:        graphql.String,
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
					return ts.GetAll(), nil
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
					return us.Get("ellie"), nil
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
					t, err := ts.Get(id)
					if err != nil {
						return nil, err
					}
					return t, nil
				},
			},
			"viewer": &graphql.Field{
				Type: userType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return us.Get("ellie"), nil
				},
			},
		},
		MutateAndGetPayload: func(inputMap map[string]interface{}, info graphql.ResolveInfo, ctx context.Context) (map[string]interface{}, error) {
			title, ok := inputMap["title"].(string)
			if !ok {
				return nil, errors.New("could not cast title to string")
			}

			t := &Task{
				Title: title,
			}
			err := ts.Create(t)
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
			"addTask": addTaskMutation,
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
