package model

import (
	"database/sql"
	"time"
)

// User type is the main structure for user.
type User struct {
	ID        int64        `db:"id"`
	Name      string       `db:"name"`
	Email     string       `db:"email"`
	Role      string       `db:"role"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}

// AuthInfo type is the structure for user authentication data from storage.
type AuthInfo struct {
	Username string `db:"name"`
	Password string `db:"password"`
	Role     string `db:"role"`
}
