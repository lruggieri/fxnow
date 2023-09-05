package logic

type CreateAPIKeyRequest struct {
}

type CreateAPIKeyResponse struct {
	APIKeyID string
}

type DeleteAPIKeyRequest struct {
	APIKeyID string
}

type DeleteAPIKeyResponse struct{}

type CreateUserRequest struct {
	FirstName string
	LastName  string
	Email     string
}

type CreateUserResponse struct {
	UserID string
}
