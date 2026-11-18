package users

import "time"

type Entity struct {
	ID        string
	Name      string
	Email     string
	CreatedAt time.Time
}
