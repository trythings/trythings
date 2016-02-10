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
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/search"
)

// Migration represents a batch update to existing entities in the datastore.
// Migrations are defined in code and are only stored in the database once they have been executed.
type Migration struct {
	Version     time.Time
	Author      string
	Description string
	RunAt       time.Time
	Run         func(ctx context.Context, ts *TaskService) error `datastore:"-"`
}

func version(timeStr string) time.Time {
	loc, err := time.LoadLocation("America/Toronto")
	if err != nil {
		panic(err)
	}

	t, err := time.ParseInLocation("2006-01-02T15:04:05", timeStr, loc)
	if err != nil {
		panic(err)
	}
	return t
}

var migrations = []*Migration{
	{
		Version:     version("2016-02-03T18:52:00"),
		Author:      "annie",
		Description: "Add createdAt time to existing tasks, defaulting to now.",
		Run: func(ctx context.Context, ts *TaskService) error {
			// TODO#Perf: Consider using a cursor and/or a batch update.
			var tasks []*Task
			_, err := datastore.NewQuery("Task").
				Ancestor(datastore.NewKey(ctx, "Root", "root", 0, nil)).
				GetAll(ctx, &tasks)
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
		},
	},
	{
		Version:     version("2016-02-10T16:37:00"),
		Author:      "annie, daniel",
		Description: "Add tasks to search index.",
		Run: func(ctx context.Context, ts *TaskService) error {
			var tasks []*Task
			_, err := datastore.NewQuery("Task").
				Ancestor(datastore.NewKey(ctx, "Root", "root", 0, nil)).
				GetAll(ctx, &tasks)
			if err != nil {
				return err
			}

			for _, t := range tasks {
				err = ts.Index(ctx, t)
				if err != nil {
					return err
				}
			}

			return nil
		},
	},
}

type MigrationService struct {
	ts *TaskService
}

func NewMigrationService(ts *TaskService) *MigrationService {
	return &MigrationService{
		ts: ts,
	}
}

// latestVersion returns the largest version stored in the Migrations table.
// Since versions are expected to be strictly increasing, any Migration with a version > latestVersion is expected to have not yet been run.
// If no Migrations have been run against the datastore, latestVersion returns the zero time.
func (s *MigrationService) latestVersion(ctx context.Context) (time.Time, error) {
	var ms []*Migration
	_, err := datastore.NewQuery("Migration").
		Ancestor(datastore.NewKey(ctx, "Root", "root", 0, nil)).
		Project("Version").
		Order("-Version").
		Limit(1).
		GetAll(ctx, &ms)
	if err != nil {
		return time.Time{}, err
	}

	if len(ms) == 0 {
		return time.Time{}, nil
	}

	return ms[0].Version, nil
}

func (s *MigrationService) run(ctx context.Context, m *Migration) error {
	if m.RunAt.IsZero() {
		m.RunAt = time.Now()
	}

	if m.Version.IsZero() {
		return errors.New("cannot run migration without version")
	}

	// TODO: Pipe rootKey through with context.
	rootKey := datastore.NewKey(ctx, "Root", "root", 0, nil)
	k := datastore.NewIncompleteKey(ctx, "Migration", rootKey)

	err := m.Run(ctx, s.ts)
	if err != nil {
		return err
	}

	_, err = datastore.Put(ctx, k, m)
	if err != nil {
		return err
	}

	return nil
}

func (s *MigrationService) RunAll(ctx context.Context) error {
	latest, err := s.latestVersion(ctx)
	if err != nil {
		return err
	}
	log.Infof(ctx, "running all migrations. latest is %s", latest)

	for _, m := range migrations {
		if m.Version.After(latest) {
			log.Infof(ctx, "running migration version %s", m.Version)
			err = s.run(ctx, m)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

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

func (t *Task) Load(fields []search.Field, meta *search.DocumentMetadata) error {
	// You should load the fields of a Task from the datastore.
	return errors.New("task should not be loaded from search")
}

func (t *Task) Save() ([]search.Field, *search.DocumentMetadata, error) {
	return []search.Field{
		{Name: "Title", Value: t.Title},
		{Name: "Description", Value: t.Description},
	}, nil, nil
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
		Order("-CreatedAt").
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

	err = s.Index(ctx, t)
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

	err = s.Index(ctx, t)
	if err != nil {
		return err
	}

	return nil
}

var numSearch = expvar.NewInt("api.*TaskService.Search")

func (s *TaskService) Search(ctx context.Context, query string) ([]*Task, error) {
	numSearch.Add(1)

	index, err := search.Open("Task")
	if err != nil {
		return nil, err
	}

	var ts []*Task
	for it := index.Search(ctx, query, &search.SearchOptions{
		IDsOnly: true,
	}); ; {
		id, err := it.Next(nil)
		if err == search.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		// TODO Use GetMulti.
		t, err := s.Get(ctx, id)
		if err != nil {
			return nil, err
		}

		ts = append(ts, t)
	}

	return ts, nil
}

func (s *TaskService) Index(ctx context.Context, t *Task) error {
	index, err := search.Open("Task")
	if err != nil {
		return err
	}
	_, err = index.Put(ctx, t.ID, t)
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

	ms := NewMigrationService(ts)

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
			"addTask":     addTaskMutation,
			"archiveTask": archiveTaskMutation,

			"migrate": relay.MutationWithClientMutationID(migrateMutationConfig(ms)),
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

// migrateMutation should only be called once, from init().
func migrateMutationConfig(ms *MigrationService) relay.MutationConfig {
	// TODO Do we really want this to be separate from init()?
	runAll := delay.Func("*MigrationService.RunAll", func(ctx context.Context) error {
		return ms.RunAll(ctx)
	})
	return relay.MutationConfig{
		Name:         "Migrate",
		InputFields:  graphql.InputObjectConfigFieldMap{},
		OutputFields: graphql.Fields{},
		MutateAndGetPayload: func(inputMap map[string]interface{}, info graphql.ResolveInfo, ctx context.Context) (map[string]interface{}, error) {
			err := runAll.Call(ctx)
			if err != nil {
				return nil, err
			}

			return map[string]interface{}{}, nil
		},
	}
}
