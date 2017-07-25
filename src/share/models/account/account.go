package account

type AuthRequest struct {
    UserId   string
    Password string
}

type AuthResponse struct {
    Id       int32
    Status   byte
    AuthKey  string `db:"auth_key"`
    SubPassChar byte
}

type VerifyReq struct {
    AuthKey   uint32
    UserIdx   uint16
    ServerId  byte
    ChannelId byte
    DBIdx     int32
}

type VerifyRes struct {
    Verified bool
}

type AuthCheckReq struct {
    Id       int32
    Password string
}

type AuthCheckRes struct {
    Result bool
}