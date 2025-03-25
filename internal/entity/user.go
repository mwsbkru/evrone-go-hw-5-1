package entity

import "time"

type UserRole string

const UserRoleAdmin UserRole = "admin"
const UserRoleUser UserRole = "user"
const UserRoleGuest UserRole = "guest"

type User struct {
	ID        string
	Name      string
	Email     string
	Role      UserRole
	CreatedAt time.Time
}
