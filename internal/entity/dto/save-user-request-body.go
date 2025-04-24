package dto

type SaveUserRequestBody struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

func (body *SaveUserRequestBody) IsValid() bool {
	return body.Email != "" && body.Name != "" && body.Role != ""
}
