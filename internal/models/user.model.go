package models

import "time"

type UserRole string

const (
	AdminRole  UserRole = "manager"
	MemberRole UserRole = "attendee"
)

type User struct {
	ID        string   `json:"id"`
	Email     string   `json:"email"`
	Role      UserRole `json:"role"`
	Password  string   `json:"password"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}
