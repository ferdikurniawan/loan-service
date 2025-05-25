package entity

import "time"

type Borrower struct {
	ID        int64  `json:"borrower_id"`
	Name      string `json:"name"`
	Active    bool   `json:"is_active"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
