package account

type AuthRequest struct {
    UserId   string
    Password string
}

type AuthResponse struct {
    Id       int
    Status   byte
    AuthKey  string `db:"auth_key"`
}