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

var errUserNotFound = errors.New("user not found")

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
}

type UserService struct {
}

func (s *UserService) Create(ctx context.Context, u *User) error {
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

func (s *UserService) byGoogleID(ctx context.Context, googleID string) (*User, error) {
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
		return nil, errUserNotFound
	}

	return us[0], nil
}

func NewUserContext(ctx context.Context, u *User) context.Context {
	return context.WithValue(ctx, userKey, u)
}

// FromContext should not be subject to access control,
// because it would create a circular dependency.
func (s *UserService) FromContext(ctx context.Context) (*User, error) {
	ctxUser, ok := ctx.Value(userKey).(*User)
	if !ok {
		return nil, errors.New("expected user info in context, probably missing authorization header")
	}

	u, err := s.byGoogleID(ctx, ctxUser.GoogleID)
	if err == errUserNotFound {
		err := s.Create(ctx, ctxUser)
		if err != nil {
			return nil, err
		}
		return s.byGoogleID(ctx, ctxUser.GoogleID)
	}
	return u, err
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
