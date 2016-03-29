package api

import (
	"errors"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/relay"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type View struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Name      string    `json:"name"`
	SpaceID   string    `json:"spaceId"`
}

type ViewService struct {
	SpaceService *SpaceService `inject:""`
}

func (s *ViewService) IsVisible(ctx context.Context, v *View) (bool, error) {
	sp, err := s.SpaceService.ByID(ctx, v.SpaceID)
	if err != nil {
		return false, err
	}
	return s.SpaceService.IsVisible(ctx, sp)
}

func (s *ViewService) ByID(ctx context.Context, id string) (*View, error) {
	rootKey := datastore.NewKey(ctx, "Root", "root", 0, nil)
	k := datastore.NewKey(ctx, "View", id, 0, rootKey)
	var v View
	err := datastore.Get(ctx, k, &v)
	if err != nil {
		return nil, err
	}

	ok, err := s.IsVisible(ctx, &v)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("cannot access view")
	}

	return &v, nil
}

func (s *ViewService) BySpace(ctx context.Context, sp *Space) ([]*View, error) {
	var vs []*View
	_, err := datastore.NewQuery("View").
		Ancestor(datastore.NewKey(ctx, "Root", "root", 0, nil)).
		Filter("SpaceID =", sp.ID).
		GetAll(ctx, &vs)
	if err != nil {
		return nil, err
	}

	var ac []*View
	for _, v := range vs {
		ok, err := s.IsVisible(ctx, v)
		if err != nil {
			// TODO use multierror
			return nil, err
		}

		if ok {
			ac = append(ac, v)
		}
	}

	return ac, nil
}

type ViewAPI struct {
	NodeInterface *graphql.Interface `inject:"node"`

	Type *graphql.Object
}

func (api *ViewAPI) Start() error {
	api.Type = graphql.NewObject(graphql.ObjectConfig{
		Name: "View",
		Fields: graphql.Fields{
			"id": relay.GlobalIDField("View", nil),
			"createdAt": &graphql.Field{
				Description: "When the view was first created.",
				Type:        graphql.String,
			},
			"name": &graphql.Field{
				Description: "The name to display for the view.",
				Type:        graphql.String,
			},
		},
		Interfaces: []*graphql.Interface{
			api.NodeInterface,
		},
	})
	return nil
}
