package api

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/relay"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/search"
	"cloud.google.com/go/trace"
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
		{Name: "CreatedAt", Value: t.CreatedAt},
		{Name: "Title", Value: t.Title},
		{Name: "Description", Value: t.Description},
		{Name: "IsArchived", Value: isArchived},
		{Name: "SpaceID", Value: search.Atom(t.SpaceID)},
	}, nil, nil
}

type TaskService struct {
	SpaceService *SpaceService `inject:""`
}

func (s *TaskService) IsVisible(ctx context.Context, t *Task) (bool, error) {
	sp, err := s.SpaceService.ByID(ctx, t.SpaceID)
	if err != nil {
		return false, err
	}
	return s.SpaceService.IsVisible(ctx, sp)
}

func (s *TaskService) ByID(ctx context.Context, id string) (*Task, error) {
	rootKey := datastore.NewKey(ctx, "Root", "root", 0, nil)
	k := datastore.NewKey(ctx, "Task", id, 0, rootKey)

	ct, ok := CacheFromContext(ctx).Get(k).(*Task)
	if ok {
		return ct, nil
	}
	span := trace.FromContext(ctx).NewChild("trythings.task.ByID")
	defer span.Finish()

	var t Task
	err := datastore.Get(ctx, k, &t)
	if err != nil {
		return nil, err
	}

	ok, err = s.IsVisible(ctx, &t)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, errors.New("cannot access task")
	}

	CacheFromContext(ctx).Set(k, &t)
	return &t, nil
}

// ByIDs filters out Tasks that are not visible to the current User.
func (s *TaskService) ByIDs(ctx context.Context, ids []string) ([]*Task, error) {
	span := trace.FromContext(ctx).NewChild("trythings.task.ByIDs")
	defer span.Finish()

	rootKey := datastore.NewKey(ctx, "Root", "root", 0, nil)

	ks := []*datastore.Key{}
	for _, id := range ids {
		ks = append(ks, datastore.NewKey(ctx, "Task", id, 0, rootKey))
	}

	var allTasks = make([]*Task, len(ks))
	err := datastore.GetMulti(ctx, ks, allTasks)
	if err != nil {
		return nil, err
	}

	ts := []*Task{}
	for _, t := range allTasks {
		// TODO#Perf: Batch the isVisible check.
		ok, err := s.IsVisible(ctx, t)
		if err != nil {
			return nil, err
		}
		if !ok {
			continue
		}

		ts = append(ts, t)
	}

	return ts, nil
}

func (s *TaskService) Create(ctx context.Context, t *Task) error {
	span := trace.FromContext(ctx).NewChild("trythings.task.Create")
	defer span.Finish()

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
	span := trace.FromContext(ctx).NewChild("trythings.task.Update")
	defer span.Finish()

	if t.ID == "" {
		return errors.New("cannot update task with no ID")
	}

	// Make sure we have access to the task to start.
	_, err := s.ByID(ctx, t.ID)
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

	CacheFromContext(ctx).Set(k, t)
	return nil
}

func (s *TaskService) Search(ctx context.Context, sp *Space, query string) (ts []*Task, err error) {
	span := trace.FromContext(ctx).NewChild("trythings.task.Search")
	defer span.Finish()

	ts, ok := CacheFromContext(ctx).SearchResults(sp, query)
	if ok {
		return ts, nil
	}
	originalQuery := query
	defer func() {
		if err == nil {
			CacheFromContext(ctx).SetSearchResults(sp, originalQuery, ts)
		}
	}()

	ok, err = s.SpaceService.IsVisible(ctx, sp)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("cannot access space to search")
	}

	// Replace the fake today() expression with the actual date.
	// TODO: Have this reflect the user's time zone.
	today := time.Now().Format(" 2006-01-02 ")
	query = strings.Replace(query, " today() ", today, -1)

	if query != "" {
		// Restrict the query to the space.
		query = fmt.Sprintf("%s AND SpaceID: %q", query, sp.ID)
	} else {
		query = fmt.Sprintf("SpaceID: %q", sp.ID)
	}

	index, err := search.Open("Task")
	if err != nil {
		return nil, err
	}

	it := index.Search(ctx, query, &search.SearchOptions{
		IDsOnly: true,
		Sort: &search.SortOptions{
			Expressions: []search.SortExpression{
				{Expr: "CreatedAt", Reverse: true},
			},
		},
	})
	ids := []string{}
	for {
		id, err := it.Next(nil)
		if err == search.Done {
			break
		}
		if err != nil {
			// TODO: Use multierror
			return nil, err
		}
		ids = append(ids, id)
	}

	// FIXME Deleted tasks may still show up in the search index,
	// so we should just not return them.
	ts, err = s.ByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	return ts, nil
}

func (s *TaskService) Index(ctx context.Context, t *Task) error {
	span := trace.FromContext(ctx).NewChild("trythings.task.Index")
	defer span.Finish()

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
	NodeInterface *graphql.Interface `inject:"node"`
	TaskService   *TaskService       `inject:""`

	Type           *graphql.Object
	ConnectionType *graphql.Object
	Mutations      map[string]*graphql.Field
}

func (api *TaskAPI) Start() error {
	api.Type = graphql.NewObject(graphql.ObjectConfig{
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
			api.NodeInterface,
		},
	})
	api.ConnectionType = relay.ConnectionDefinitions(relay.ConnectionConfig{
		Name:     api.Type.Name(),
		NodeType: api.Type,
	}).ConnectionType

	api.Mutations = map[string]*graphql.Field{
		"addTask": relay.MutationWithClientMutationID(relay.MutationConfig{
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
					Type: api.Type,
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						payload, ok := p.Source.(map[string]interface{})
						if !ok {
							return nil, errors.New("could not cast payload to map")
						}
						id, ok := payload["taskId"].(string)
						if !ok {
							return nil, errors.New("could not cast taskId to string")
						}
						t, err := api.TaskService.ByID(p.Context, id)
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
				err := api.TaskService.Create(ctx, t)
				if err != nil {
					return nil, err
				}

				return map[string]interface{}{
					"taskId": t.ID,
				}, nil
			},
		}),
		"editTask": relay.MutationWithClientMutationID(relay.MutationConfig{
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
					Type: api.Type,
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						payload, ok := p.Source.(map[string]interface{})
						if !ok {
							return nil, errors.New("could not cast payload to map")
						}
						id, ok := payload["id"].(string)
						if !ok {
							return nil, errors.New("could not cast id to string")
						}
						t, err := api.TaskService.ByID(p.Context, id)
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

				t, err := api.TaskService.ByID(ctx, resolvedID.ID)
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

				err = api.TaskService.Update(ctx, t)
				if err != nil {
					return nil, err
				}

				return map[string]interface{}{
					"id": t.ID,
				}, nil
			},
		}),
	}

	return nil
}
