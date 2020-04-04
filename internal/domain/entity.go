package domain

import "time"

type Created struct {
	CreatedAt time.Time `json:"createdAt,omitempty" db:"created_at"`
}

type Updated struct {
	UpdatedAt time.Time `json:"updatedAt,omitempty" db:"updated_at"`
}
