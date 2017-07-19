package session

import "share/models/account"

type Data struct {
    AccountId   int32
    Verified    bool
    LoggedIn    bool
    SubPassword *account.SubPassword
}
