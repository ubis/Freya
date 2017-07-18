package account

type UserVerify struct {
    AuthKey   uint32
    UserIdx   uint16
    ServerId  byte
    ChannelId byte
}

type UserVerifyRecv struct {
    Verified bool
}