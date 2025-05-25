package entity

import "github.com/google/uuid"

type Investor struct {
	ID   uuid.UUID `json:"investor_id"`
	Name string    `json:"investor_name"`
}
