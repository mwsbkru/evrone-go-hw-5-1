package entity

import "time"

// UserRole represents user role
type UserRole string

// UserRoleAdmin is admin role
const UserRoleAdmin UserRole = "admin"

// UserRoleUser is member user role
const UserRoleUser UserRole = "user"

// UserRoleGuest is not registered user role
const UserRoleGuest UserRole = "guest"

// User represents user entity
type User struct {
	ID        string
	Name      string
	Email     string
	Role      UserRole
	CreatedAt time.Time
}
