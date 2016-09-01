package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/relay"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/cloud/trace"
)

// TODO#DatamodelCleanup: Consider composing the database properties into some union GraphqlSearch struct.
type Search struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Name      string    `json:"name"`
	Query     string    `json:"query"`

	// TODO#DatamodelCleanup: Consider using the PropertyLoadSaver to populate these.
	ParentTaskID string `datastore:"-" json:"-"`
}

type clientSearchID struct {
	ID    string `json:"id"`
	Query string `json:"query"`
	// TODO#NestedSearch: Make this a search path instead.
	ParentTaskID string `json:"parentTaskId"`
}

// TODO#Cleanup: This should probably live on the service.
func (se *Search) ClientID() (string, error) {
	var cid clientSearchID

	// For now, these IDs are unstable.
	// That is, converting a temporary search into a saved search will return a search with a different client id.
	if se.ID != "" {
		// Saved search
		cid.ID = se.ID
	} else {
		// Temporary search
		cid.Query = se.Query

		if se.ParentTaskID == "" {
			return "", fmt.Errorf("Expected search to have a non-empty ParentTaskID")
		}
		cid.ParentTaskID = se.ParentTaskID
	}

	jsonID, err := json.Marshal(cid)
	if err != nil {
		return "", err
	}
	return string(jsonID), nil
}

type SearchService struct {
	TaskService *TaskService `inject:""`
}

func (s *SearchService) IsVisible(ctx context.Context, se *Search) (bool, error) {
	if se.ParentTaskID == "" {
		// TODO: Consider trying to denormalize this a second time before error-ing.
		return false, fmt.Errorf("Expected search to have a non-empty ParentTaskID")
	}

	t, err := s.TaskService.ByID(ctx, se.ParentTaskID)
	if err != nil {
		return false, err
	}

	return s.TaskService.IsVisible(ctx, t)
}

func (s *SearchService) ByClientID(ctx context.Context, clientID string) (*Search, error) {
	var cid clientSearchID
	err := json.Unmarshal([]byte(clientID), &cid)
	if err != nil {
		return nil, err
	}

	// Saved search
	if cid.ID != "" {
		return s.ByID(ctx, cid.ID)
	}

	// Temporary search
	return &Search{
		Query:        cid.Query,
		ParentTaskID: cid.ParentTaskID,
	}, nil
}

func (s *SearchService) ByID(ctx context.Context, id string) (*Search, error) {
	rootKey := datastore.NewKey(ctx, "Root", "root", 0, nil)
	k := datastore.NewKey(ctx, "Search", id, 0, rootKey)

	cse, ok := CacheFromContext(ctx).Get(k).(*Search)
	if ok {
		return cse, nil
	}
	span := trace.FromContext(ctx).NewChild("trythings.search.ByID")
	defer span.Finish()

	var se Search
	err := datastore.Get(ctx, k, &se)
	if err != nil {
		return nil, err
	}

	// Denormalize the parent task id onto the search.
	var tsrels []*TaskSearchRelationship
	_, err = datastore.NewQuery("TaskSearchRelationship").
		Ancestor(rootKey).
		Filter("SearchID =", id).
		Project("TaskID").
		GetAll(ctx, &tsrels)
	if err != nil {
		return nil, err
	}
	if len(tsrels) == 0 {
		// TODO#errors: This error should (maybe) be swallowed instead.
		return nil, errors.New("a search must have a parent task")
	}
	if len(tsrels) > 1 {
		// TODO#errors: This error should be swallowed instead.
		return nil, errors.New("a search cannot have more than one parent task")
	}
	se.ParentTaskID = tsrels[0].TaskID

	ok, err = s.IsVisible(ctx, &se)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("cannot access search")
	}

	CacheFromContext(ctx).Set(k, &se)
	return &se, nil
}

func (s *SearchService) Create(ctx context.Context, se *Search) error {
	span := trace.FromContext(ctx).NewChild("trythings.search.Create")
	defer span.Finish()

	if se.ID != "" {
		return fmt.Errorf("se already has id %q", se.ID)
	}

	if se.CreatedAt.IsZero() {
		se.CreatedAt = time.Now()
	}

	if se.Name == "" {
		return errors.New("Name is required")
	}

	if se.ParentTaskID == "" {
		return errors.New("ParentTaskID is required")
	}
	_, err := s.TaskService.ByID(ctx, se.ParentTaskID)
	if err != nil {
		return err
	}

	if se.Query == "" {
		return errors.New("Query is required")
	}

	ok, err := s.IsVisible(ctx, se)
	if err != nil {
		return err
	}

	if !ok {
		return errors.New("cannot access view to create search")
	}

	id, _, err := datastore.AllocateIDs(ctx, "Search", nil, 1)
	if err != nil {
		return err
	}
	se.ID = fmt.Sprintf("%x", id)

	rootKey := datastore.NewKey(ctx, "Root", "root", 0, nil)
	k := datastore.NewKey(ctx, "Search", se.ID, 0, rootKey)
	k, err = datastore.Put(ctx, k, se)
	if err != nil {
		return err
	}
	// TODO: This should maybe live on the task service.
	krel := datastore.NewIncompleteKey(ctx, "TaskSearchRelationship", rootKey)
	tsrel := &TaskSearchRelationship{
		TaskID:   se.ParentTaskID,
		SearchID: se.ID,
	}
	krel, err = datastore.Put(ctx, krel, tsrel)
	if err != nil {
		return err
	}

	return nil
}

func (s *SearchService) Update(ctx context.Context, se *Search) error {
	span := trace.FromContext(ctx).NewChild("trythings.search.Update")
	defer span.Finish()

	if se.ID == "" {
		return errors.New("cannot update search with no ID")
	}

	// Make sure we have access to the search before it was modified.
	_, err := s.ByID(ctx, se.ID)
	if err != nil {
		return err
	}

	// TODO#Validation: You cannot change the parent task ID.

	// Make sure we continue to have access to the task after our update.
	ok, err := s.IsVisible(ctx, se)
	if err != nil {
		return err
	}

	if !ok {
		return errors.New("cannot update search to lose access")
	}

	rootKey := datastore.NewKey(ctx, "Root", "root", 0, nil)
	k := datastore.NewKey(ctx, "Search", se.ID, 0, rootKey)
	_, err = datastore.Put(ctx, k, se)
	if err != nil {
		return err
	}

	CacheFromContext(ctx).Set(k, se)
	return nil
}

type SearchAPI struct {
	NodeInterface *graphql.Interface `inject:"node"`
	SearchService *SearchService     `inject:""`
	TaskService   *TaskService       `inject:""`
	TaskAPI       *TaskAPI           `inject:""`

	Type           *graphql.Object
	ConnectionType *graphql.Object
}

func (api *SearchAPI) Start() error {
	api.Type = graphql.NewObject(graphql.ObjectConfig{
		Name: "Search",
		Fields: graphql.Fields{
			"id": relay.GlobalIDField("Search", func(obj interface{}, info graphql.ResolveInfo, ctx context.Context) (string, error) {
				se, ok := obj.(*Search)
				if !ok {
					// TODO#Bug: We *should* only be seeing *Search objs. Figure out who is calling this with a Search.
					se2, ok := obj.(Search)
					if !ok {
						return "", fmt.Errorf("Search's GlobalIDField() was called with a non-Search")
					}
					se = &se2
				}

				cid, err := se.ClientID()
				if err != nil {
					return "", fmt.Errorf("Failed to create a ClientID for %v", se)
				}
				return cid, nil
			}),
			"createdAt": &graphql.Field{
				Description: "When the search was first saved.",
				Type:        graphql.String,
			},
			"name": &graphql.Field{
				Description: "The name to display for the search.",
				Type:        graphql.String,
			},
			// TODO#Perf: Consider storing the search results on the context or ResolveInfo to avoid computing them twice (numResults and results).
			"numResults": &graphql.Field{
				Description: "The total number of results that match the query",
				Type:        graphql.Int,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					span := trace.FromContext(p.Context).NewChild("trythings.searchAPI.numResults")
					defer span.Finish()

					se, ok := p.Source.(*Search)
					if !ok {
						return nil, errors.New("expected search source")
					}

					// TODO: What about if the parent task id is not set? Would be nice if these methods all took the search instead.

					pt, err := api.TaskService.ByID(p.Context, se.ParentTaskID)
					if err != nil {
						return nil, err
					}

					// TODO#Perf: Run a count query instead of fetching all of the matches.
					ts, err := api.TaskService.Search(p.Context, pt, se.Query)
					if err != nil {
						return nil, err
					}

					return len(ts), nil
				},
			},
			"query": &graphql.Field{
				Description: "The query used to search for tasks.",
				Type:        graphql.String,
			},
			"results": &graphql.Field{
				Description: "The tasks that match the query",
				Type:        api.TaskAPI.ConnectionType,
				Args:        relay.ConnectionArgs,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					span := trace.FromContext(p.Context).NewChild("trythings.searchAPI.results")
					defer span.Finish()

					se, ok := p.Source.(*Search)
					if !ok {
						return nil, errors.New("expected search source")
					}

					pt, err := api.TaskService.ByID(p.Context, se.ParentTaskID)
					if err != nil {
						return nil, err
					}

					ts, err := api.TaskService.Search(p.Context, pt, se.Query)
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
	return nil
}
