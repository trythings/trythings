package api

import (
	"time"

	"golang.org/x/net/context"
)

type View struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Name      string    `json:"name"`
	SpaceID   string    `json:"spaceId"`
}

type ViewService struct {
	spaces *SpaceService
}

func (s *ViewService) IsVisible(ctx context.Context, v *View) (bool, error) {
	sp, err := s.spaces.ByID(ctx, v.SpaceID)
	if err != nil {
		return false, err
	}
	return s.spaces.IsVisible(ctx, sp)
}
