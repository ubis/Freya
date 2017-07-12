package server

const (
    LOGIN_SERVER_TYPE = 1
)

type RegisterRequest struct {
    Type byte
}

type RegisterResponse struct {
    Registered bool
}
