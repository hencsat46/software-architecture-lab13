package notifications

import "time"

type Entity struct {
	ID        string
	UserID    string
	Message   string
	CreatedAt time.Time
}
