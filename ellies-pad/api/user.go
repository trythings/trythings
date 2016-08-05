package api

import (
	"errors"
	"fmt"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/relay"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/cloud/trace"
)

var ErrUserNotFound = errors.New("user not found")

type User struct {
	ID              string    `json:"id"`
	CreatedAt       time.Time `json:"createdAt"`
	IsAdmin         bool      `json:"isAdmin"`
	GoogleID        string    `json:"-"`
	Email           string    `json:"email"`
	IsEmailVerified bool      `json:"-"`
	Name            string    `json:"name"`
	GivenName       string    `json:"givenName"`
	FamilyName      string    `json:"familyName"`
	ImageURL        string    `json:"imageUrl"`

	RootTaskID string `json:"rootTaskId"`
}

type UserService struct {
}

func (s *UserService) IsVisible(ctx context.Context, u *User) (bool, error) {
	span := trace.FromContext(ctx).NewChild("trythings.user.IsVisible")
	defer span.Finish()

	su, err := IsSuperuser(ctx)
	if err != nil {
		return false, err
	}
	if su {
		return true, nil
	}

	me, err := s.FromContext(ctx)
	if err != nil {
		return false, err
	}
	return me.ID == u.ID, nil
}

func (s *UserService) ByID(ctx context.Context, id string) (*User, error) {
	rootKey := datastore.NewKey(ctx, "Root", "root", 0, nil)
	k := datastore.NewKey(ctx, "User", id, 0, rootKey)

	cu, ok := CacheFromContext(ctx).Get(k).(*User)
	if ok {
		return cu, nil
	}

	var u User
	err := datastore.Get(ctx, k, &u)
	if err != nil {
		return nil, err
	}

	ok, err = s.IsVisible(ctx, &u)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, errors.New("cannot access user")
	}

	CacheFromContext(ctx).Set(k, &u)
	return &u, nil
}

func (s *UserService) Create(ctx context.Context, u *User) error {
	span := trace.FromContext(ctx).NewChild("trythings.user.Create")
	defer span.Finish()

	// TODO Make sure u.GoogleID == user.Current(ctx).ID

	if u.ID != "" {
		return fmt.Errorf("u already has id %q", u.ID)
	}

	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now()
	}

	id, _, err := datastore.AllocateIDs(ctx, "User", nil, 1)
	if err != nil {
		return err
	}
	u.ID = fmt.Sprintf("%x", id)

	rootKey := datastore.NewKey(ctx, "Root", "root", 0, nil)
	k := datastore.NewKey(ctx, "User", u.ID, 0, rootKey)
	k, err = datastore.Put(ctx, k, u)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) Update(ctx context.Context, u *User) error {
	span := trace.FromContext(ctx).NewChild("trythings.user.Update")
	defer span.Finish()

	if u.ID == "" {
		return errors.New("cannot update user with no ID")
	}

	// Make sure that we have access to the user to start.
	_, err := s.ByID(ctx, u.ID)
	if err != nil {
		return err
	}

	// Make sure we continue to have access to the task after the update.
	ok, err := s.IsVisible(ctx, u)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("cannot update user to lose access")
	}

	rootKey := datastore.NewKey(ctx, "Root", "root", 0, nil)
	k := datastore.NewKey(ctx, "User", u.ID, 0, rootKey)
	k, err = datastore.Put(ctx, k, u)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) byGoogleID(ctx context.Context, googleID string) (*User, error) {
	span := trace.FromContext(ctx).NewChild("trythings.user.byGoogleID")
	defer span.Finish()

	var us []*User
	_, err := datastore.NewQuery("User").
		Ancestor(datastore.NewKey(ctx, "Root", "root", 0, nil)).
		Filter("GoogleID =", googleID).
		Limit(1).
		GetAll(ctx, &us)
	if err != nil {
		return nil, err
	}

	if len(us) == 0 {
		return nil, ErrUserNotFound
	}

	return us[0], nil
}

// FromContext should not be subject to access control,
// because it would create a circular dependency.
func (s *UserService) FromContext(ctx context.Context) (*User, error) {
	span := trace.FromContext(ctx).NewChild("trythings.user.FromContext")
	defer span.Finish()

	gu, ok := GoogleUserFromContext(ctx)
	if !ok {
		return nil, errors.New("expected google user, probably missing Authorization header")
	}
	return s.byGoogleID(ctx, gu.ID)
}

type UserAPI struct {
	NodeInterface *graphql.Interface `inject:"node"`
	SpaceService  *SpaceService      `inject:""`
	SpaceAPI      *SpaceAPI          `inject:""`
	UserService   *UserService       `inject:""`

	Type *graphql.Object
}

func (api *UserAPI) Start() error {
	api.Type = graphql.NewObject(graphql.ObjectConfig{
		Name:        "User",
		Description: "User represents a person who can interact with the app.",
		Fields: graphql.Fields{
			"id": relay.GlobalIDField("User", nil),
			"isAdmin": &graphql.Field{
				Description: "Whether or not the user is an Ellie's Pad admin.",
				Type:        graphql.Boolean,
			},
			"email": &graphql.Field{
				Description: "The user's email primary address.",
				Type:        graphql.String,
			},
			"name": &graphql.Field{
				Description: "The user's full name.",
				Type:        graphql.String,
			},
			"givenName": &graphql.Field{
				Description: "The user's given name.",
				Type:        graphql.String,
			},
			"familyName": &graphql.Field{
				Description: "The user's family name.",
				Type:        graphql.String,
			},
			"imageUrl": &graphql.Field{
				Description: "The user's profile picture URL.",
				Type:        graphql.String,
			},
			"space": &graphql.Field{
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type:         graphql.String,
						DefaultValue: "",
						Description:  "id can be omitted, which will have space resolve to the user's default space.",
					},
				},
				Description: "space is a disjoint universe of views, searches and tasks.",
				Type:        api.SpaceAPI.Type,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					span := trace.FromContext(p.Context).NewChild("trythings.userAPI.space")
					defer span.Finish()

					id, ok := p.Args["id"].(string)
					if ok {
						resolvedID := relay.FromGlobalID(id)
						if resolvedID == nil {
							return nil, fmt.Errorf("invalid id %q", id)
						}

						sp, err := api.SpaceService.ByID(p.Context, resolvedID.ID)
						if err != nil {
							return nil, err
						}
						return sp, nil
					}

					u, ok := p.Source.(*User)
					if !ok {
						return nil, errors.New("expected user source")
					}

					sps, err := api.SpaceService.ByUser(p.Context, u)
					if err != nil {
						return nil, err
					}

					if len(sps) == 0 {
						return nil, errors.New("could not find default space for user")
					}

					return sps[0], nil
				},
			},
			"spaces": &graphql.Field{
				Type: graphql.NewList(api.SpaceAPI.Type),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					span := trace.FromContext(p.Context).NewChild("trythings.userAPI.spaces")
					defer span.Finish()

					u, ok := p.Source.(*User)
					if !ok {
						return nil, errors.New("expected user source")
					}

					sps, err := api.SpaceService.ByUser(p.Context, u)
					if err != nil {
						return nil, err
					}

					return sps, nil
				},
			},
		},
		Interfaces: []*graphql.Interface{
			api.NodeInterface,
		},
	})
	return nil
}
