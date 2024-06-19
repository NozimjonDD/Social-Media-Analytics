package entity

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type User struct {
	ID        uuid.UUID `json:"id" bun:"id,pk"`
	FirstName string    `json:"first_name" bun:"first_name"`
	LastName  string    `json:"last_name" bun:"last_name"`
	Username  string    `json:"username" bun:"username"`
	Password  *string   `json:"password,omitempty" bun:"password"`
	Role      string    `json:"role" bun:"-"`
}

func (u *User) MarshalJSON() ([]byte, error) {
	u.Role = "admin"
	return json.Marshal(*u)
}

func (u *User) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		u.ID = uuid.New()
	}
	return nil
}
