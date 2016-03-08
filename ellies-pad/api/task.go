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
)

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

type TaskAPI struct {
	tasks           *TaskService
	nodeDefinitions *relay.NodeDefinitions

	typ *graphql.Object
}

func (api *TaskAPI) Start() error {
	api.typ = graphql.NewObject(graphql.ObjectConfig{
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
			api.nodeDefinitions.NodeInterface,
		},
	})

	return nil
}

func (api *TaskAPI) Type() *graphql.Object {
	return api.typ
}
