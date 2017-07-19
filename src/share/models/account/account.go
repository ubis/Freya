package account

type AuthRequest struct {
    UserId   string
    Password string
}

type AuthResponse struct {
    Id          int32
    Status      byte
    AuthKey     string `db:"auth_key"`
    SubPassChar byte
}