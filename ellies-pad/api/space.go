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

type Space struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Name      string    `json:"name"`
	UserIDs   []string  `json:"userIds"`
}

type SpaceService struct {
	UserService *UserService `inject:""`
}

func (s *SpaceService) ByID(ctx context.Context, id string) (*Space, error) {
	rootKey := datastore.NewKey(ctx, "Root", "root", 0, nil)
	k := datastore.NewKey(ctx, "Space", id, 0, rootKey)

	csp, ok := CacheFromContext(ctx).Get(k).(*Space)
	if ok {
		return csp, nil
	}
	span := trace.FromContext(ctx).NewChild("trythings.space.ByID")
	defer span.Finish()

	var sp Space
	err := datastore.Get(ctx, k, &sp)
	if err != nil {
		return nil, err
	}

	ok, err = s.IsVisible(ctx, &sp)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("cannot access space")
	}

	CacheFromContext(ctx).Set(k, &sp)
	return &sp, nil
}

func (s *SpaceService) ByUser(ctx context.Context, u *User) ([]*Space, error) {
	span := trace.FromContext(ctx).NewChild("trythings.space.ByUser")
	defer span.Finish()

	var sps []*Space
	_, err := datastore.NewQuery("Space").
		Ancestor(datastore.NewKey(ctx, "Root", "root", 0, nil)).
		Filter("UserIDs =", u.ID).
		GetAll(ctx, &sps)
	if err != nil {
		return nil, err
	}

	var ac []*Space
	for _, sp := range sps {
		ok, err := s.IsVisible(ctx, sp)
		if err != nil {
			// TODO use multierror
			return nil, err
		}

		if ok {
			ac = append(ac, sp)
		}
	}

	return ac, nil
}

func (s *SpaceService) Create(ctx context.Context, sp *Space) error {
	span := trace.FromContext(ctx).NewChild("trythings.space.Create")
	defer span.Finish()

	if sp.ID != "" {
		return fmt.Errorf("sp already has id %q", sp.ID)
	}

	if sp.CreatedAt.IsZero() {
		sp.CreatedAt = time.Now()
	}

	if len(sp.UserIDs) > 0 {
		return errors.New("UserIDs must be empty")
	}

	su, err := IsSuperuser(ctx)
	if err != nil {
		return err
	}
	if !su {
		u, err := s.UserService.FromContext(ctx)
		if err != nil {
			return err
		}
		sp.UserIDs = []string{u.ID}
	}

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

func (s *SpaceService) IsVisible(ctx context.Context, sp *Space) (isVisible bool, err error) {
	isVisible, ok := CacheFromContext(ctx).IsVisible(sp)
	if ok {
		return isVisible, nil
	}
	defer func() {
		if err == nil {
			CacheFromContext(ctx).SetIsVisible(sp, isVisible)
		}
	}()

	span := trace.FromContext(ctx).NewChild("trythings.space.IsVisible")
	defer span.Finish()

	su, err := IsSuperuser(ctx)
	if err != nil {
		return false, err
	}

	if su {
		return true, nil
	}

	u, err := s.UserService.FromContext(ctx)
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

type SpaceAPI struct {
	NodeInterface *graphql.Interface `inject:"node"`
	SearchService *SearchService     `inject:""`
	SearchAPI     *SearchAPI         `inject:""`
	SpaceService  *SpaceService      `inject:""`
	TaskService   *TaskService       `inject:""`
	TaskAPI       *TaskAPI           `inject:""`
	ViewService   *ViewService       `inject:""`
	ViewAPI       *ViewAPI           `inject:""`

	Type *graphql.Object
}

func (api *SpaceAPI) Start() error {
	api.Type = graphql.NewObject(graphql.ObjectConfig{
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
			"savedSearch": &graphql.Field{
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Type: api.SearchAPI.Type,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					span := trace.FromContext(p.Context).NewChild("trythings.spaceAPI.savedSearch")
					defer span.Finish()

					id, ok := p.Args["id"].(string)
					if !ok {
						return nil, errors.New("id is required")
					}
					resolvedID := relay.FromGlobalID(id)
					if resolvedID == nil {
						return nil, fmt.Errorf("invalid id %q", id)
					}

					se, err := api.SearchService.ByClientID(p.Context, resolvedID.ID)
					if err != nil {
						return nil, err
					}
					return se, nil
				},
			},
			"querySearch": &graphql.Field{
				Args: graphql.FieldConfigArgument{
					"query": &graphql.ArgumentConfig{
						Type:         graphql.String,
						DefaultValue: "",
						Description:  "query filters the result to only tasks that contain particular terms in their title or description",
					},
				},
				Type: api.SearchAPI.Type,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					span := trace.FromContext(p.Context).NewChild("trythings.spaceAPI.querySearch")
					defer span.Finish()

					sp, ok := p.Source.(*Space)
					if !ok {
						return nil, errors.New("expected a space source")
					}

					q, ok := p.Args["query"].(string)
					if !ok {
						q = "" // Return all tasks.
					}

					return &Search{
						Query:   q,
						SpaceID: sp.ID,
					}, nil
				},
			},
			"view": &graphql.Field{
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type:         graphql.String,
						DefaultValue: "",
						Description:  "id can be omitted, which will have view resolve to the space's default view.",
					},
				},
				Type: api.ViewAPI.Type,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					span := trace.FromContext(p.Context).NewChild("trythings.spaceAPI.view")
					defer span.Finish()

					id, ok := p.Args["id"].(string)
					if ok {
						resolvedID := relay.FromGlobalID(id)
						if resolvedID == nil {
							return nil, fmt.Errorf("invalid id %q", id)
						}

						v, err := api.ViewService.ByID(p.Context, resolvedID.ID)
						if err != nil {
							return nil, err
						}
						return v, nil
					}

					sp, ok := p.Source.(*Space)
					if !ok {
						return nil, errors.New("expected space source")
					}

					vs, err := api.ViewService.BySpace(p.Context, sp)
					if err != nil {
						return nil, err
					}

					if len(vs) == 0 {
						return nil, errors.New("could not find default view for space")
					}

					return vs[0], nil
				},
			},
			"views": &graphql.Field{
				Type: graphql.NewList(api.ViewAPI.Type),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					span := trace.FromContext(p.Context).NewChild("trythings.spaceAPI.views")
					defer span.Finish()

					sp, ok := p.Source.(*Space)
					if !ok {
						return nil, errors.New("expected space source")
					}

					vs, err := api.ViewService.BySpace(p.Context, sp)
					if err != nil {
						return nil, err
					}

					return vs, nil
				},
			},
		},
		Interfaces: []*graphql.Interface{
			api.NodeInterface,
		},
	})

	return nil
}
