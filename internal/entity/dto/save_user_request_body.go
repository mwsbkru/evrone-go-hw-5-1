package dto

// SaveUserRequestBody represents request body for save user action
type SaveUserRequestBody struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// IsValid checks is request body for save user action valid
func (body *SaveUserRequestBody) IsValid() bool {
	return body.Email != "" && body.Name != "" && body.Role != ""
}
