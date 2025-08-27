package model

import "time"

type User struct {
	ID        int64     `db:"id"`
	Email     string    `db:"email"`
	PassHash  []byte    `db:"password"`
	IsAdmin   bool      `db:"is_admin"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
