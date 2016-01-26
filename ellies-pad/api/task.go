package api

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/relay"
)

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

func (s *TaskService) Get(id string) *Task {
	for _, t := range s.tasks {
		if t.ID == id {
			return t
		}
	}
	return nil
}

func (s *TaskService) GetAll() []*Task {
	return s.tasks
}

var taskType *graphql.Object

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

	nodeDefinitions = relay.NewNodeDefinitions(relay.NodeDefinitionsConfig{
		IDFetcher: func(id string, info graphql.ResolveInfo) interface{} {
			resolvedID := relay.FromGlobalID(id)
			if resolvedID.Type == "Task" {
				return ts.Get(resolvedID.ID)
			}
			return nil
		},
		TypeResolve: func(value interface{}, info graphql.ResolveInfo) *graphql.Object {
			switch value.(type) {
			case *Task:
				return taskType
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

	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"node": nodeDefinitions.NodeField,
			"tasks": &graphql.Field{
				Description: "All tasks that need to be completed",
				Type:        graphql.NewList(taskType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return ts.GetAll(), nil
				},
			},
		},
	})

	var err error
	Schema, err = graphql.NewSchema(graphql.SchemaConfig{
		Query: queryType,
	})
	if err != nil {
		panic(err)
	}
}
