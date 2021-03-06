package api

import (
	"errors"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/relay"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/delay"
	"google.golang.org/appengine/log"
	"cloud.google.com/go/trace"
)

// Migration represents a batch update to existing entities in the datastore.
// Migrations are defined in code and are only stored in the database once they have been executed.
type Migration struct {
	Version     time.Time
	Author      string
	Description string
	RunAt       time.Time
	Run         func(ctx context.Context, s *MigrationService) error `datastore:"-"`
}

func version(timeStr string) time.Time {
	loc, err := time.LoadLocation("America/Toronto")
	if err != nil {
		panic(err)
	}

	t, err := time.ParseInLocation("2006-01-02T15:04:05", timeStr, loc)
	if err != nil {
		panic(err)
	}
	return t
}

// reindexTasks adds all tasks from the datastore into the search index.
var reindexTasks = func(ctx context.Context, s *MigrationService) error {
	span := trace.FromContext(ctx).NewChild("trythings.migration.reindexTasks")
	defer span.Finish()

	var tasks []*Task
	_, err := datastore.NewQuery("Task").
		Ancestor(datastore.NewKey(ctx, "Root", "root", 0, nil)).
		GetAll(ctx, &tasks)
	if err != nil {
		return err
	}

	for _, t := range tasks {
		err = s.TaskService.Index(ctx, t)
		if err != nil {
			return err
		}
	}

	return nil
}

var migrations = []*Migration{
	{
		Version:     version("2016-02-03T18:52:00"),
		Author:      "annie",
		Description: "Add createdAt time to existing tasks, defaulting to now.",
		Run: func(ctx context.Context, s *MigrationService) error {
			// TODO#Perf: Consider using a cursor and/or a batch update.
			var tasks []*Task
			_, err := datastore.NewQuery("Task").
				Ancestor(datastore.NewKey(ctx, "Root", "root", 0, nil)).
				GetAll(ctx, &tasks)
			if err != nil {
				return err
			}

			for _, t := range tasks {
				if t.CreatedAt.IsZero() {
					t.CreatedAt = time.Now()
					err = s.TaskService.Update(ctx, t)
					if err != nil {
						return err
					}
				}
			}

			return nil
		},
	},
	{
		Version:     version("2016-02-10T16:37:00"),
		Author:      "annie, daniel",
		Description: "Add tasks to search index.",
		Run:         reindexTasks,
	},
	{
		Version:     version("2016-02-16T21:20:00"),
		Author:      "annie, daniel",
		Description: "Add task.IsArchived to search index.",
		Run:         reindexTasks,
	},
	{
		Version:     version("2016-02-27T19:20:00"),
		Author:      "annie, daniel",
		Description: "Add Annie and Daniel's space.",
		Run: func(ctx context.Context, s *MigrationService) error {
			numSpaces, err := datastore.NewQuery("Space").
				Ancestor(datastore.NewKey(ctx, "Root", "root", 0, nil)).
				Count(ctx)
			if err != nil {
				return err
			}

			if numSpaces != 0 {
				return nil
			}

			sp := &Space{
				Name: "Annie and Daniel",
			}
			err = s.SpaceService.Create(ctx, sp)
			if err != nil {
				return err
			}

			var tasks []*Task
			_, err = datastore.NewQuery("Task").
				Ancestor(datastore.NewKey(ctx, "Root", "root", 0, nil)).
				GetAll(ctx, &tasks)
			if err != nil {
				return err
			}

			for _, t := range tasks {
				if t.SpaceID == "" {
					t.SpaceID = sp.ID
					err = s.TaskService.Update(ctx, t)
					if err != nil {
						return err
					}
				}
			}

			return nil
		},
	},
	{
		Version:     version("2016-02-29T02:22:00"),
		Author:      "annie, daniel",
		Description: "Add task.SpaceID to search index.",
		Run:         reindexTasks,
	},
	{
		Version:     version("2016-02-29T22:59:00"),
		Author:      "annie, daniel",
		Description: "Add users to default space.",
		Run: func(ctx context.Context, s *MigrationService) error {
			root := datastore.NewKey(ctx, "Root", "root", 0, nil)

			var sps []*Space
			_, err := datastore.NewQuery("Space").
				Ancestor(root).
				Limit(1).
				GetAll(ctx, &sps)
			if err != nil {
				return err
			}

			if len(sps) == 0 {
				return errors.New("expected a space")
			}
			sp := sps[0]

			var us []*User
			_, err = datastore.NewQuery("User").
				Ancestor(root).
				GetAll(ctx, &us)
			if err != nil {
				return err
			}

			var ids []string
			for _, u := range us {
				ids = append(ids, u.ID)
			}

			sp.UserIDs = ids
			_, err = datastore.Put(ctx, datastore.NewKey(ctx, "Space", sp.ID, 0, root), sp)
			if err != nil {
				return err
			}

			return nil
		},
	},
	{
		Version:     version("2016-03-29T14:14:00"),
		Author:      "annie, daniel",
		Description: "Add default view to each space.",
		Run: func(ctx context.Context, s *MigrationService) error {
			root := datastore.NewKey(ctx, "Root", "root", 0, nil)

			var sps []*Space
			_, err := datastore.NewQuery("Space").
				Ancestor(root).
				GetAll(ctx, &sps)
			if err != nil {
				return err
			}

			for _, sp := range sps {
				vs, err := s.ViewService.BySpace(ctx, sp)
				if err != nil {
					return err
				}

				if len(vs) != 0 {
					continue
				}

				err = s.ViewService.Create(ctx, &View{
					Name:    "Default",
					SpaceID: sp.ID,
				})
				if err != nil {
					return err
				}
			}

			return nil
		},
	},
	{
		Version:     version("2016-03-29T15:18:00"),
		Author:      "annie, daniel",
		Description: "Add default searches to each view.",
		Run: func(ctx context.Context, s *MigrationService) error {
			root := datastore.NewKey(ctx, "Root", "root", 0, nil)

			var vs []*View
			_, err := datastore.NewQuery("View").
				Ancestor(root).
				GetAll(ctx, &vs)
			if err != nil {
				return err
			}

			for _, v := range vs {
				ss, err := s.SearchService.ByView(ctx, v)
				if err != nil {
					return err
				}

				if len(ss) != 0 {
					continue
				}

				err = s.SearchService.Create(ctx, &Search{
					Name:     "#now",
					ViewID:   v.ID,
					Query:    "#now AND IsArchived: false",
					ViewRank: datastore.ByteString("0"),
				})
				if err != nil {
					return err
				}

				err = s.SearchService.Create(ctx, &Search{
					Name:     "Incoming",
					ViewID:   v.ID,
					Query:    "NOT #now AND NOT #next AND NOT #later AND IsArchived: false",
					ViewRank: datastore.ByteString("1"),
				})
				if err != nil {
					return err
				}

				err = s.SearchService.Create(ctx, &Search{
					Name:     "#next",
					ViewID:   v.ID,
					Query:    "#next AND NOT #now AND IsArchived: false",
					ViewRank: datastore.ByteString("2"),
				})
				if err != nil {
					return err
				}

				err = s.SearchService.Create(ctx, &Search{
					Name:     "#later",
					ViewID:   v.ID,
					Query:    "#later AND NOT #next AND NOT #now AND IsArchived: false",
					ViewRank: datastore.ByteString("3"),
				})
				if err != nil {
					return err
				}

				err = s.SearchService.Create(ctx, &Search{
					Name:     "Archived",
					ViewID:   v.ID,
					Query:    "IsArchived: true",
					ViewRank: datastore.ByteString("4"),
				})
				if err != nil {
					return err
				}
			}

			return nil
		},
	},
	{
		Version:     version("2016-03-29T22:57:00"),
		Author:      "annie, daniel",
		Description: "Add ranks to searches.",
		Run: func(ctx context.Context, s *MigrationService) error {
			root := datastore.NewKey(ctx, "Root", "root", 0, nil)

			var vs []*View
			_, err := datastore.NewQuery("View").
				Ancestor(root).
				GetAll(ctx, &vs)
			if err != nil {
				return err
			}

			for _, v := range vs {
				ss, err := s.SearchService.ByView(ctx, v)
				if err != nil {
					return err
				}

				rs := NewRanks(len(ss))
				for i, se := range ss {
					se.ViewRank = datastore.ByteString(rs[i])
					err := s.SearchService.Update(ctx, se)
					if err != nil {
						return err
					}
				}
			}

			return nil
		},
	},
	{
		Version:     version("2016-07-28T15:20:00"),
		Author:      "annie, daniel",
		Description: "Add spaceID to searches",
		Run: func(ctx context.Context, s *MigrationService) error {
			root := datastore.NewKey(ctx, "Root", "root", 0, nil)

			var vs []*View
			_, err := datastore.NewQuery("View").
				Ancestor(root).
				GetAll(ctx, &vs)
			if err != nil {
				return err
			}

			for _, v := range vs {
				sp, err := s.ViewService.Space(ctx, v)
				if err != nil {
					return err
				}

				ss, err := s.SearchService.ByView(ctx, v)
				if err != nil {
					return err
				}

				for _, se := range ss {
					se.SpaceID = sp.ID
					err := s.SearchService.Update(ctx, se)
					if err != nil {
						return err
					}
				}
			}

			return nil
		},
	},
}

type MigrationService struct {
	SearchService *SearchService `inject:""`
	SpaceService  *SpaceService  `inject:""`
	TaskService   *TaskService   `inject:""`
	ViewService   *ViewService   `inject:""`
}

// latestVersion returns the largest version stored in the Migrations table.
// Since versions are expected to be strictly increasing,
// any Migration with a version > latestVersion is expected to have not yet been run.
// If no Migrations have been run against the datastore, latestVersion returns the zero time.
func (s *MigrationService) latestVersion(ctx context.Context) (time.Time, error) {
	var ms []*Migration
	_, err := datastore.NewQuery("Migration").
		Ancestor(datastore.NewKey(ctx, "Root", "root", 0, nil)).
		Project("Version").
		Order("-Version").
		Limit(1).
		GetAll(ctx, &ms)
	if err != nil {
		return time.Time{}, err
	}

	if len(ms) == 0 {
		return time.Time{}, nil
	}

	return ms[0].Version, nil
}

func (s *MigrationService) run(ctx context.Context, m *Migration) error {
	if m.RunAt.IsZero() {
		m.RunAt = time.Now()
	}

	if m.Version.IsZero() {
		return errors.New("cannot run migration without version")
	}

	// TODO: Pipe rootKey through with context.
	rootKey := datastore.NewKey(ctx, "Root", "root", 0, nil)
	k := datastore.NewIncompleteKey(ctx, "Migration", rootKey)

	err := m.Run(ctx, s)
	if err != nil {
		return err
	}

	_, err = datastore.Put(ctx, k, m)
	if err != nil {
		return err
	}

	return nil
}

func (s *MigrationService) RunAll(ctx context.Context) error {
	span := trace.FromContext(ctx).NewChild("trythings.migration.RunAll")
	defer span.Finish()

	su, err := IsSuperuser(ctx)
	if err != nil {
		return err
	}
	if !su {
		return errors.New("must run migrations as superuser")
	}

	latest, err := s.latestVersion(ctx)
	if err != nil {
		return err
	}
	log.Infof(ctx, "running all migrations. latest is %s", latest)

	for _, m := range migrations {
		if m.Version.After(latest) {
			log.Infof(ctx, "running migration version %s", m.Version)
			err = s.run(ctx, m)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

type MigrationAPI struct {
	UserService      *UserService      `inject:""`
	MigrationService *MigrationService `inject:""`
	Mutations        map[string]*graphql.Field
}

func (api *MigrationAPI) Start() error {
	runAll := delay.Func("*MigrationService.RunAll", func(ctx context.Context) error {
		return api.MigrationService.RunAll(AsSuperuser(ctx))
	})

	api.Mutations = map[string]*graphql.Field{
		"migrate": relay.MutationWithClientMutationID(relay.MutationConfig{
			Name:         "Migrate",
			InputFields:  graphql.InputObjectConfigFieldMap{},
			OutputFields: graphql.Fields{},
			MutateAndGetPayload: func(inputMap map[string]interface{}, info graphql.ResolveInfo, ctx context.Context) (map[string]interface{}, error) {
				u, err := api.UserService.FromContext(ctx)
				if err != nil {
					return nil, err
				}

				if !u.IsAdmin {
					return nil, errors.New("user must be an admin")
				}

				err = runAll.Call(ctx)
				if err != nil {
					return nil, err
				}

				return map[string]interface{}{}, nil
			},
		}),
	}

	return nil
}
