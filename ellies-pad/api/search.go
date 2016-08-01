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
)

type Search struct {
	ID        string               `json:"id"`
	CreatedAt time.Time            `json:"createdAt"`
	Name      string               `json:"name"`
	SpaceID   string               `json:"spaceId"`
	ViewID    string               `json:"viewId"`
	ViewRank  datastore.ByteString `json:"viewRank"`
	Query     string               `json:"query"`
}

type clientSearchID struct {
	ID      string `json:"id"`
	Query   string `json:"query"`
	SpaceID string `json:"spaceId"`
}

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
		cid.SpaceID = se.SpaceID
	}

	jsonID, err := json.Marshal(cid)
	if err != nil {
		return "", err
	}
	return string(jsonID), nil
}

type SearchService struct {
	SpaceService *SpaceService `inject:""`
	ViewService  *ViewService  `inject:""`
}

func (s *SearchService) IsVisible(ctx context.Context, se *Search) (bool, error) {
	sp, err := s.SpaceService.ByID(ctx, se.SpaceID)
	if err != nil {
		return false, err
	}
	return s.SpaceService.IsVisible(ctx, sp)
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
		SpaceID: cid.SpaceID,
		Query:   cid.Query,
	}, nil
}

func (s *SearchService) ByID(ctx context.Context, id string) (*Search, error) {
	rootKey := datastore.NewKey(ctx, "Root", "root", 0, nil)
	k := datastore.NewKey(ctx, "Search", id, 0, rootKey)

	cse, ok := CacheFromContext(ctx).Get(k).(*Search)
	if ok {
		return cse, nil
	}

	var se Search
	err := datastore.Get(ctx, k, &se)
	if err != nil {
		return nil, err
	}

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

func (s *SearchService) ByView(ctx context.Context, v *View) ([]*Search, error) {
	var ss []*Search
	_, err := datastore.NewQuery("Search").
		Ancestor(datastore.NewKey(ctx, "Root", "root", 0, nil)).
		Filter("ViewID =", v.ID).
		Order("ViewRank").
		GetAll(ctx, &ss)
	if err != nil {
		return nil, err
	}

	var ac []*Search
	for _, se := range ss {
		ok, err := s.IsVisible(ctx, se)
		if err != nil {
			// TODO use multierror
			return nil, err
		}

		if ok {
			ac = append(ac, se)
		}
	}

	return ac, nil
}

func (s *SearchService) Create(ctx context.Context, se *Search) error {
	if se.ID != "" {
		return fmt.Errorf("se already has id %q", se.ID)
	}

	if se.CreatedAt.IsZero() {
		se.CreatedAt = time.Now()
	}

	if se.Name == "" {
		return errors.New("Name is required")
	}

	if se.ViewID == "" {
		return errors.New("ViewID is required")
	}

	v, err := s.ViewService.ByID(ctx, se.ViewID)
	if err != nil {
		return err
	}

	if se.SpaceID == "" {
		se.SpaceID = v.SpaceID
	}

	if se.SpaceID != v.SpaceID {
		return errors.New("Search's SpaceID must match View's")
	}

	if len(se.ViewRank) != 0 {
		return fmt.Errorf("se already has a view rank %x", se.ViewRank)
	}

	// TODO#Performance: Add a shared or per-request cache to support these small, repeated queries.

	if se.Query == "" {
		return errors.New("Query is required")
	}

	rootKey := datastore.NewKey(ctx, "Root", "root", 0, nil)

	// Create a ViewRank for the search.
	// It should come after every other search in the view.
	var ranks []*struct {
		ViewRank datastore.ByteString
	}
	_, err = datastore.NewQuery("Search").
		Ancestor(rootKey).
		Filter("ViewID =", se.ViewID).
		Project("ViewRank").
		Order("-ViewRank").
		Limit(1).
		GetAll(ctx, &ranks)
	if err != nil {
		return err
	}

	maxViewRank := MinRank
	if len(ranks) != 0 {
		maxViewRank = Rank(ranks[0].ViewRank)
	}
	rank, err := NewRank(maxViewRank, MaxRank)
	if err != nil {
		return err
	}
	se.ViewRank = datastore.ByteString(rank)

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

	k := datastore.NewKey(ctx, "Search", se.ID, 0, rootKey)
	k, err = datastore.Put(ctx, k, se)
	if err != nil {
		return err
	}

	return nil
}

func (s *SearchService) Update(ctx context.Context, se *Search) error {
	if se.ID == "" {
		return errors.New("cannot update search with no ID")
	}

	// Make sure we have access to the search before it was modified.
	_, err := s.ByID(ctx, se.ID)
	if err != nil {
		return err
	}

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

func (s *SearchService) Space(ctx context.Context, se *Search) (*Space, error) {
	return s.SpaceService.ByID(ctx, se.SpaceID)
}

type SearchAPI struct {
	NodeInterface *graphql.Interface `inject:"node"`
	SearchService *SearchService     `inject:""`
	TaskService   *TaskService       `inject:""`
	TaskAPI       *TaskAPI           `inject:""`

	Type *graphql.Object
}

func (api *SearchAPI) Start() error {
	api.Type = graphql.NewObject(graphql.ObjectConfig{
		Name: "Search",
		Fields: graphql.Fields{
			"id": relay.GlobalIDField("Search", func(obj interface{}, info graphql.ResolveInfo, ctx context.Context) (string, error) {
				se, ok := obj.(*Search)
				if !ok {
					return "", fmt.Errorf("Search's GlobalIDField() was called with a non-Search")
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
					se, ok := p.Source.(*Search)
					if !ok {
						return nil, errors.New("expected search source")
					}

					sp, err := api.SearchService.Space(p.Context, se)
					if err != nil {
						return nil, err
					}

					// TODO#Perf: Run a count query instead of fetching all of the matches.
					ts, err := api.TaskService.Search(p.Context, sp, se.Query)
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
					se, ok := p.Source.(*Search)
					if !ok {
						return nil, errors.New("expected search source")
					}

					sp, err := api.SearchService.Space(p.Context, se)
					if err != nil {
						return nil, err
					}

					ts, err := api.TaskService.Search(p.Context, sp, se.Query)
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
	return nil
}
