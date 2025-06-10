package dto

// UserResponseBody represents response body for get single user action
type UserResponseBody struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// UsersResponseBody represents response body for get list of users action
type UsersResponseBody struct {
	Data []*UserResponseBody `json:"data"`
}
