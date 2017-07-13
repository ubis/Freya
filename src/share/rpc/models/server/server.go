package server

const (
    LOGIN_SERVER_TYPE = 1
)

type RegRequest struct {
    Type byte
}

type RegResponse struct {
    Registered bool
}