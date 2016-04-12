package api

import "time"

type Channel struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Token     string    `json:"-" datastore:"-"`
	UserID    string    `json:"userId"`
}

// type ChannelService struct {
// }

// func (s *ChannelService) Create(ctx context.Context, c *Channel) error {

// }

// func (s *ChannelService) DeleteByID(ctx context.Context, id string) error {

// }

// func (s *ChannelService) SendToUserID(ctx context.Context, userID string) error {

// }
