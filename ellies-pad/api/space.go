package api

import (
	"errors"
	"fmt"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/relay"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
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

func (s *SpaceService) IsVisible(ctx context.Context, sp *Space) (bool, error) {
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
	SpaceService  *SpaceService      `inject:""`
	TaskService   *TaskService       `inject:""`
	TaskAPI       *TaskAPI           `inject:""`

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
			"tasks": &graphql.Field{
				Args: graphql.FieldConfigArgument{
					"query": &graphql.ArgumentConfig{
						Type:         graphql.String,
						DefaultValue: "",
						Description:  "query filters the result to only tasks that contain particular terms in their title or description",
					},
				},
				Description: "tasks are all pieces of work that need to be completed for the user.",
				Type:        graphql.NewList(api.TaskAPI.Type),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					sp, ok := p.Source.(*Space)
					if !ok {
						return nil, errors.New("expected a space source")
					}

					q, ok := p.Args["query"].(string)
					if !ok {
						q = "" // Return all tasks.
					}

					return api.TaskService.Search(p.Context, sp, q)
				},
			},
		},
		Interfaces: []*graphql.Interface{
			api.NodeInterface,
		},
	})

	return nil
}
