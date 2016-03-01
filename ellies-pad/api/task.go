package api

import (
	"errors"
	"fmt"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/relay"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/search"
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

// Task represents a particular action or piece of work to be completed.
type Task struct {
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	Title       string    `json:"title"`
	Description string    `json:"description" datastore:",noindex"`
	IsArchived  bool      `json:"isArchived"`
	SpaceID     string    `json:"spaceId"`
}

func (t *Task) Load(fields []search.Field, meta *search.DocumentMetadata) error {
	// You should load the fields of a Task from the datastore.
	return errors.New("task should not be loaded from search")
}

func (t *Task) Save() ([]search.Field, *search.DocumentMetadata, error) {
	isArchived := search.Atom("false")
	if t.IsArchived {
		isArchived = search.Atom("true")
	}
	return []search.Field{
		{Name: "Title", Value: t.Title},
		{Name: "Description", Value: t.Description},
		{Name: "IsArchived", Value: isArchived},
		{Name: "SpaceID", Value: search.Atom(t.SpaceID)},
	}, nil, nil
}

type TaskService struct {
	spaces *SpaceService
}

func NewTaskService(spaces *SpaceService) *TaskService {
	return &TaskService{
		spaces: spaces,
	}
}

func (s *TaskService) IsVisible(ctx context.Context, t *Task) (bool, error) {
	sp, err := s.spaces.ByID(ctx, t.SpaceID)
	if err != nil {
		return false, err
	}
	return s.spaces.IsVisible(ctx, sp)
}

type ErrAccessDenied struct{}

func (e ErrAccessDenied) Error() string {
	return "cannot access task"
}

func (s *TaskService) Get(ctx context.Context, id string) (*Task, error) {
	rootKey := datastore.NewKey(ctx, "Root", "root", 0, nil)
	k := datastore.NewKey(ctx, "Task", id, 0, rootKey)
	var t Task
	err := datastore.Get(ctx, k, &t)
	if err != nil {
		return nil, err
	}

	ok, err := s.IsVisible(ctx, &t)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, ErrAccessDenied{}
	}

	return &t, nil
}

func (s *TaskService) Create(ctx context.Context, t *Task) error {
	if t.ID != "" {
		return fmt.Errorf("t already has id %q", t.ID)
	}

	if t.CreatedAt.IsZero() {
		t.CreatedAt = time.Now()
	}

	if t.SpaceID == "" {
		return errors.New("SpaceID is required")
	}

	ok, err := s.IsVisible(ctx, t)
	if err != nil {
		return err
	}

	if !ok {
		return errors.New("cannot access space to create task")
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

func (s *TaskService) Update(ctx context.Context, t *Task) error {
	if t.ID == "" {
		return errors.New("cannot update task with no ID")
	}

	// Make sure we have access to the task to start.
	_, err := s.Get(ctx, t.ID)
	if err != nil {
		return err
	}

	// Make sure we continue to have access to the task after our update.
	ok, err := s.IsVisible(ctx, t)
	if err != nil {
		return err
	}

	if !ok {
		return errors.New("cannot update task to lose access")
	}

	rootKey := datastore.NewKey(ctx, "Root", "root", 0, nil)
	k := datastore.NewKey(ctx, "Task", t.ID, 0, rootKey)
	_, err = datastore.Put(ctx, k, t)
	if err != nil {
		return err
	}

	err = s.Index(ctx, t)
	if err != nil {
		return err
	}

	return nil
}

func (s *TaskService) Search(ctx context.Context, sp *Space, query string) ([]*Task, error) {
	ok, err := s.spaces.IsVisible(ctx, sp)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("cannot access space to search")
	}

	// Restrict the query to the space.
	query = fmt.Sprintf("%s AND SpaceID: %q", query, sp.ID)

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
		// FIXME Deleted tasks may still show up in the search index,
		// so we should just not return them.
		t, err := s.Get(ctx, id)
		if err != nil {
			if _, ok := err.(ErrAccessDenied); ok {
				continue
			}
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

	u, err := s.users.FromContext(ctx)
	if err != nil {
		return err
	}
	sp.UserIDs = []string{u.ID}

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

var taskType *graphql.Object
var spaceType *graphql.Object
var userType *graphql.Object

var nodeDefinitions *relay.NodeDefinitions

var Schema graphql.Schema

func init() {
	us := NewUserService()
	ss := NewSpaceService(us)
	ts := NewTaskService(ss)
	ms := NewMigrationService(ss, ts)

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
				Type:        graphql.NewList(taskType),
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
				Type: taskType,
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
		mutation:   mutationType,
	}
	migrationAPI.AddMutations()

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
