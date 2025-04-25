package dto

type UserResponseBody struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type UsersResponseBody struct {
	Data []UserResponseBody `json:"data"`
}
