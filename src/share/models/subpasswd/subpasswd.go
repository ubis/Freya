package subpasswd

import "time"

type FetchReq struct {
    Account int32
}

type FetchRes struct {
    Details
}

type SetReq struct {
    Account     int32
    Details
}

type SetRes struct {
    Success bool
}

type Details struct {
    Password  string
    Answer    string
    Question  byte
    Expires   time.Time
    Verified  bool
    FailTimes byte
}