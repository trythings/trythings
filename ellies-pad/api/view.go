package api

import (
	"errors"
	"fmt"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/relay"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"cloud.google.com/go/trace"
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

	cv, ok := CacheFromContext(ctx).Get(k).(*View)
	if ok {
		return cv, nil
	}
	span := trace.FromContext(ctx).NewChild("trythings.view.ByID")
	defer span.Finish()

	var v View
	err := datastore.Get(ctx, k, &v)
	if err != nil {
		return nil, err
	}

	ok, err = s.IsVisible(ctx, &v)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("cannot access view")
	}

	CacheFromContext(ctx).Set(k, &v)
	return &v, nil
}

func (s *ViewService) BySpace(ctx context.Context, sp *Space) ([]*View, error) {
	span := trace.FromContext(ctx).NewChild("trythings.view.BySpace")
	defer span.Finish()

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

func (s *ViewService) Create(ctx context.Context, v *View) error {
	span := trace.FromContext(ctx).NewChild("trythings.view.Create")
	defer span.Finish()

	if v.ID != "" {
		return fmt.Errorf("v already has id %q", v.ID)
	}

	if v.CreatedAt.IsZero() {
		v.CreatedAt = time.Now()
	}

	if v.Name == "" {
		return errors.New("Name is required")
	}

	if v.SpaceID == "" {
		return errors.New("SpaceID is required")
	}

	ok, err := s.IsVisible(ctx, v)
	if err != nil {
		return err
	}

	if !ok {
		return errors.New("cannot access space to create view")
	}

	id, _, err := datastore.AllocateIDs(ctx, "View", nil, 1)
	if err != nil {
		return err
	}
	v.ID = fmt.Sprintf("%x", id)

	rootKey := datastore.NewKey(ctx, "Root", "root", 0, nil)
	k := datastore.NewKey(ctx, "View", v.ID, 0, rootKey)
	k, err = datastore.Put(ctx, k, v)
	if err != nil {
		return err
	}

	return nil
}

func (s *ViewService) Space(ctx context.Context, v *View) (*Space, error) {
	return s.SpaceService.ByID(ctx, v.SpaceID)
}

type ViewAPI struct {
	NodeInterface *graphql.Interface `inject:"node"`
	SearchService *SearchService     `inject:""`
	SearchAPI     *SearchAPI         `inject:""`

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
			"searches": &graphql.Field{
				Type: graphql.NewList(api.SearchAPI.Type),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					span := trace.FromContext(p.Context).NewChild("trythings.viewAPI.searches")
					defer span.Finish()

					v, ok := p.Source.(*View)
					if !ok {
						return nil, errors.New("expected view source")
					}

					ss, err := api.SearchService.ByView(p.Context, v)
					if err != nil {
						return nil, err
					}

					return ss, nil
				},
			},
		},
		Interfaces: []*graphql.Interface{
			api.NodeInterface,
		},
	})
	return nil
}
