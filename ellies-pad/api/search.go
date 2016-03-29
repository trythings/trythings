package api

import (
	"errors"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/relay"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type Search struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Name      string    `json:"name"`
	ViewID    string    `json:"viewId"`
	Query     string    `json:"query"`
}

type SearchService struct {
	ViewService *ViewService `inject:""`
}

func (s *SearchService) IsVisible(ctx context.Context, se *Search) (bool, error) {
	v, err := s.ViewService.ByID(ctx, se.ViewID)
	if err != nil {
		return false, err
	}
	return s.ViewService.IsVisible(ctx, v)
}

func (s *SearchService) ByView(ctx context.Context, v *View) ([]*Search, error) {
	var ss []*Search
	_, err := datastore.NewQuery("Search").
		Ancestor(datastore.NewKey(ctx, "Root", "root", 0, nil)).
		Filter("ViewID =", v.ID).
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

func (s *SearchService) Space(ctx context.Context, se *Search) (*Space, error) {
	v, err := s.ViewService.ByID(ctx, se.ViewID)
	if err != nil {
		return nil, err
	}
	return s.ViewService.Space(ctx, v)
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
			"id": relay.GlobalIDField("Search", nil),
			"createdAt": &graphql.Field{
				Description: "When the search was first saved.",
				Type:        graphql.String,
			},
			"name": &graphql.Field{
				Description: "The name to display for the search.",
				Type:        graphql.String,
			},
			"query": &graphql.Field{
				Description: "The query used to search for tasks.",
				Type:        graphql.String,
			},
			"tasks": &graphql.Field{
				Description: "The tasks that match the query.",
				Type:        graphql.NewList(api.TaskAPI.Type),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					se, ok := p.Source.(*Search)
					if !ok {
						return nil, errors.New("expected view source")
					}

					sp, err := api.SearchService.Space(p.Context, se)
					if err != nil {
						return nil, err
					}

					ts, err := api.TaskService.Search(p.Context, sp, se.Query)
					if err != nil {
						return nil, err
					}

					return ts, nil
				},
			},
		},
		Interfaces: []*graphql.Interface{
			api.NodeInterface,
		},
	})
	return nil
}
