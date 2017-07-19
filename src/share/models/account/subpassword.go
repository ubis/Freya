package account

import "time"

type SubPasswordReq struct {
    Account int32
}

type SubPassword struct {
    Password  string
    Answer    string
    Question  byte
    Expires    time.Time
    Verified  bool
    FailTimes byte
}

type SetSubPass struct {
    Account int32
    SubPassword
}

type SubPassResp struct {
    Success bool
}