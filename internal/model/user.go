package model

import (
	"database/sql"
	"time"
)

// User type is the main structure for user.
type User struct {
	ID        int64
	Name      string
	Email     string
	Role      string
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}
