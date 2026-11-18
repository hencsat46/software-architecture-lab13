package orders

import "time"

type Entity struct {
	ID        string
	UserID    string
	Product   string
	Quantity  int
	CreatedAt time.Time
}
