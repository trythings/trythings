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
	"google.golang.org/cloud/trace"
)

// Task represents a particular action or piece of work to be completed.
type Task struct {
	ID         string    `json:"id"`
	CreatedAt  time.Time `json:"createdAt"`
	Title      string    `json:"title"`
	Body       string    `json:"body" datastore:",noindex"`
	IsArchived bool      `json:"isArchived"`
}

// TODO#Ranks: We might eventually want ranks here so that we can load a limited number of subtasks / searches.
// For reference, we used to do this on search creation:
// var ranks []*struct {
// 	ViewRank datastore.ByteString
// }
// _, err = datastore.NewQuery("Search").
// 	Ancestor(rootKey).
// 	Filter("ViewID =", se.ViewID).
// 	Project("ViewRank").
// 	Order("-ViewRank").
// 	Limit(1).
// 	GetAll(ctx, &ranks)
// if err != nil {
// 	return err
// }
// maxViewRank := MinRank
// if len(ranks) != 0 {
// 	maxViewRank = Rank(ranks[0].ViewRank)
// }
// rank, err := NewRank(maxViewRank, MaxRank)
// if err != nil {
// 	return err
// }
// se.ViewRank = datastore.ByteString(rank)

type TaskSubtaskRelationship struct {
	ParentTaskID string `json:"parentTaskId"`
	ChildTaskID  string `json:"childTaskId"`
}

type TaskSearchRelationship struct {
	TaskID   string `json:"taskId"`
	SearchID string `json:"searchId"`
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
		{Name: "Body", Value: t.Body},
		{Name: "IsArchived", Value: isArchived},
		// TODO#xcxc: Figure out how to populate this with appropriate ancestors (more denormalization!).
		// {Name: "SpaceID", Value: search.Atom(t.SpaceID)},
	}, nil, nil
}

type TaskService struct {
	UserService *UserService `inject:""`
}

func (s *TaskService) IsVisible(ctx context.Context, t *Task) (bool, error) {
	// TODO#AccessControl: Add back access control, which will do some form of "edge" access control.
	return true, nil
}

// TODO: For later, here is Space's old IsVisible:
// isVisible, ok := CacheFromContext(ctx).IsVisible(sp)
// 	if ok {
// 		return isVisible, nil
// 	}
// 	defer func() {
// 		if err == nil {
// 			CacheFromContext(ctx).SetIsVisible(sp, isVisible)
// 		}
// 	}()

// 	span := trace.FromContext(ctx).NewChild("trythings.space.IsVisible")
// 	defer span.Finish()

// 	su, err := IsSuperuser(ctx)
// 	if err != nil {
// 		return false, err
// 	}

// 	if su {
// 		return true, nil
// 	}

// 	u, err := s.UserService.FromContext(ctx)
// 	if err != nil {
// 		return false, err
// 	}

// 	for _, id := range sp.UserIDs {
// 		if u.ID == id {
// 			return true, nil
// 		}
// 	}

// 	return false, nil

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

// TODO#ModelRewrite: General task creation should take a parent task?

func (s *TaskService) CreateRelationship(ctx context.Context, childTask *Task, parentTask *Task) error {
	span := trace.FromContext(ctx).NewChild("trythings.task.CreateRelationship")
	defer span.Finish()

	// Do not create duplicates of a relationship.
	rootKey := datastore.NewKey(ctx, "Root", "root", 0, nil)
	var existing []*TaskSubtaskRelationship
	_, err := datastore.NewQuery("TaskSubtaskRelationship").
		Ancestor(rootKey).
		Filter("ParentTaskID =", parentTask.ID).
		Filter("ChildTaskID =", childTask.ID).
		GetAll(ctx, &existing)
	if err != nil {
		return err
	}
	if len(existing) > 0 {
		return fmt.Errorf("Task %s is already a child of %s", childTask.ID, parentTask.ID)
	}

	// Does the child exist?
	_, err = s.ByID(ctx, childTask.ID)
	if err != nil {
		return err
	}

	// Does the parent exist?
	_, err = s.ByID(ctx, parentTask.ID)
	if err != nil {
		return err
	}

	// Create the relationship
	k := datastore.NewIncompleteKey(ctx, "TaskSubtaskRelationship", rootKey)
	_, err = datastore.Put(ctx, k, &TaskSubtaskRelationship{
		ChildTaskID:  childTask.ID,
		ParentTaskID: parentTask.ID,
	})
	if err != nil {
		return err
	}

	return nil
}

// TODO#Perf: Consider caching this relationship.
func (s *TaskService) Subtasks(ctx context.Context, pt *Task) ([]*Task, error) {
	span := trace.FromContext(ctx).NewChild("trythings.task.Subtasks")
	defer span.Finish()

	rootKey := datastore.NewKey(ctx, "Root", "root", 0, nil)
	var subtaskRels []*TaskSubtaskRelationship
	_, err := datastore.NewQuery("TaskSubtaskRelationship").
		Ancestor(rootKey).
		Filter("ParentTaskID =", pt.ID).
		GetAll(ctx, &subtaskRels)
	if err != nil {
		return nil, err
	}

	var subtaskIDs []string
	for _, st := range subtaskRels {
		subtaskIDs = append(subtaskIDs, st.ChildTaskID)
	}

	return s.ByIDs(ctx, subtaskIDs)
}

// TODO#Perf: Consider caching this relationship.
// TODO#CircularDependencies: Call out to search service's ByID rather than taking it in.
func (s *TaskService) Searches(ctx context.Context, pt *Task, byID func(ctx context.Context, id string) (*Search, error)) ([]*Search, error) {
	span := trace.FromContext(ctx).NewChild("trythings.task.Searches")
	defer span.Finish()

	rootKey := datastore.NewKey(ctx, "Root", "root", 0, nil)
	var searchRels []*TaskSearchRelationship
	_, err := datastore.NewQuery("TaskSearchRelationship").
		Ancestor(rootKey).
		Filter("TaskID =", pt.ID).
		GetAll(ctx, &searchRels)
	if err != nil {
		return nil, err
	}

	var searchIDs []string
	for _, se := range searchRels {
		searchIDs = append(searchIDs, se.SearchID)
	}

	// TODO#Perf: Consider a batch ByID for searches.
	var searches []*Search
	for _, id := range searchIDs {
		se, err := byID(ctx, id)
		if err != nil {
			return nil, err
		}
		searches = append(searches, se)
	}

	return searches, nil
}

func (s *TaskService) Create(ctx context.Context, pt *Task, t *Task) error {
	span := trace.FromContext(ctx).NewChild("trythings.task.Create")
	defer span.Finish()

	if t.ID != "" {
		return fmt.Errorf("t already has id %q", t.ID)
	}

	if t.CreatedAt.IsZero() {
		t.CreatedAt = time.Now()
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

	if pt != nil {
		err = s.CreateRelationship(ctx, t, pt)
		if err != nil {
			return err
		}
	}

	err = s.Index(ctx, t)
	if err != nil {
		return err
	}

	return nil
}

// TODO#Transactional: Make sure updates are rolled back if an error is returned at any level.

// TODO#CircularDependencies: This probably actually belongs on the UserService.
// To avoid circular dependencies, I'm leaving it here. Once we have a distinction between service implementations and interfaces, move it.

func (s *TaskService) GetOrCreateRootTask(ctx context.Context, u *User) (*Task, error) {
	if u.RootTaskID != "" {
		return s.ByID(ctx, u.RootTaskID)
	}

	t := &Task{
		Title:      fmt.Sprintf("%s's Home", u.GivenName),
		IsArchived: false,
	}
	err := s.Create(ctx, nil, t)
	if err != nil {
		return nil, err
	}

	if t.ID == "" {
		return nil, errors.New("Expected newly-created task to have a non-empty ID")
	}

	u.RootTaskID = t.ID
	err = s.UserService.Update(ctx, u)
	if err != nil {
		return nil, err
	}

	return t, nil
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

	// TODO#Validation: Every task should be a root task or have a parent task.

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

func (s *TaskService) Search(ctx context.Context, pt *Task, query string) (ts []*Task, err error) {
	span := trace.FromContext(ctx).NewChild("trythings.task.Search")
	defer span.Finish()

	ts, ok := CacheFromContext(ctx).SearchResults(pt, query)
	if ok {
		return ts, nil
	}
	originalQuery := query
	defer func() {
		if err == nil {
			CacheFromContext(ctx).SetSearchResults(pt, originalQuery, ts)
		}
	}()

	ok, err = s.IsVisible(ctx, pt)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("cannot access task to search")
	}

	// Replace the fake today() expression with the actual date.
	// TODO: Have this reflect the user's time zone.
	today := time.Now().Format(" 2006-01-02 ")
	query = strings.Replace(query, " today() ", today, -1)

	// TODO#Search: Restrict the search query to pt and its subtasks.
	// if query != "" {
	// 	// Restrict the query to the space.
	// 	query = fmt.Sprintf("%s AND SpaceID: %q", query, sp.ID)
	// } else {
	// 	query = fmt.Sprintf("SpaceID: %q", sp.ID)
	// }

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
	// TODO#CircularDependencies
	SearchAPI *SearchAPI

	Type           *graphql.Object
	ConnectionType *graphql.Object
	Mutations      map[string]*graphql.Field
}

func (api *TaskAPI) AfterStart(searchAPI *SearchAPI) {
	// Our dependency-injection library doesn't support circular dependencies, so we add this manually here.
	api.SearchAPI = searchAPI

	// This doesn't strictly need to be here, but it seems to belong.
	api.Type.AddFieldConfig("subtasks", &graphql.Field{
		Type: api.ConnectionType,
		Args: relay.ConnectionArgs,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			span := trace.FromContext(p.Context).NewChild("trythings.taskAPI.subtasks")
			defer span.Finish()

			pt, ok := p.Source.(*Task)
			if !ok {
				return nil, errors.New("expected task source")
			}

			ts, err := api.TaskService.Subtasks(p.Context, pt)
			if err != nil {
				return nil, err
			}

			objs := []interface{}{}
			for _, t := range ts {
				objs = append(objs, *t)
			}

			// TODO#Performance: Run a limited query instead of filtering after the query.
			args := relay.NewConnectionArguments(p.Args)
			return relay.ConnectionFromArray(objs, args), nil
		},
	})

	api.Type.AddFieldConfig("searches", &graphql.Field{
		Type: api.SearchAPI.ConnectionType,
		Args: relay.ConnectionArgs,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			span := trace.FromContext(p.Context).NewChild("trythings.taskAPI.searches")
			defer span.Finish()

			pt, ok := p.Source.(*Task)
			if !ok {
				return nil, errors.New("expected task source")
			}

			ses, err := api.TaskService.Searches(p.Context, pt, api.SearchAPI.SearchService.ByID)
			if err != nil {
				return nil, err
			}

			objs := []interface{}{}
			for _, se := range ses {
				objs = append(objs, *se)
			}

			// TODO#Performance: Run a limited query instead of filtering after the query.
			args := relay.NewConnectionArguments(p.Args)
			return relay.ConnectionFromArray(objs, args), nil
		},
	})

	api.Type.AddFieldConfig("querySearch", &graphql.Field{
		Args: graphql.FieldConfigArgument{
			"query": &graphql.ArgumentConfig{
				Type:         graphql.String,
				DefaultValue: "",
				Description:  "query filters the result to only subtasks that contain particular terms in their title or description",
			},
		},
		Type: api.SearchAPI.Type,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			span := trace.FromContext(p.Context).NewChild("trythings.taskAPI.querySearch")
			defer span.Finish()

			t, ok := p.Source.(*Task)
			if !ok {
				return nil, errors.New("expected a task source")
			}

			q, ok := p.Args["query"].(string)
			if !ok {
				q = "" // Return all subtasks.
			}

			return &Search{
				Query:        q,
				ParentTaskID: t.ID,
			}, nil
		},
	})

	api.Type.AddFieldConfig("savedSearch", &graphql.Field{
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Type: api.SearchAPI.Type,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			span := trace.FromContext(p.Context).NewChild("trythings.taskAPI.savedSearch")
			defer span.Finish()

			id, ok := p.Args["id"].(string)
			if !ok {
				return nil, errors.New("id is required")
			}
			resolvedID := relay.FromGlobalID(id)
			if resolvedID == nil {
				return nil, fmt.Errorf("invalid id %q", id)
			}

			se, err := api.SearchAPI.SearchService.ByClientID(p.Context, resolvedID.ID)
			if err != nil {
				return nil, err
			}
			return se, nil
		},
	})
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
			"body": &graphql.Field{
				Type: graphql.String,
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

	// TODO#Features: Support creating a search.
	// TODO#Features: Support moving a task into a different task (for drag and drop).
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
				"parentTaskId": &graphql.InputObjectFieldConfig{
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

				var body string
				bodyOrNil := inputMap["body"]
				if bodyOrNil != nil {
					body, ok = bodyOrNil.(string)
					if !ok {
						return nil, errors.New("could not cast body to string")
					}
				}

				parentTaskID, ok := inputMap["parentTaskId"].(string)
				if !ok {
					return nil, errors.New("could not cast parentTaskId to string")
				}
				resolvedParentTaskID := relay.FromGlobalID(parentTaskID)
				if resolvedParentTaskID == nil {
					return nil, fmt.Errorf("invalid id %q", parentTaskID)
				}
				parentTask, err := api.TaskService.ByID(ctx, resolvedParentTaskID.ID)
				if err != nil {
					return nil, err
				}

				t := &Task{
					Title: title,
					Body:  body,
				}
				err = api.TaskService.Create(ctx, parentTask, t)
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
				"body": &graphql.InputObjectFieldConfig{
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

				body, ok := inputMap["body"].(string)
				if ok {
					t.Body = body
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
