package account

const (
    None                = 31   // no status
    Normal              = 32   // login authentication is good
    Incorrect           = 33   // wrong id or password
    Online              = 34   // account is logged in already
    OutOfService        = 35   // cannot connect at the moment; [EP16 or never - shows nothing]
    AccountExpired      = 36   // account expired
    IpBanned            = 37   // ip is blocked, though not account related
    Banned              = 38   // account id is blocked
    TestServerTrial     = 39   // "cannot use test server during free trial period, bla bla"...
    PcCafe              = 40   // "use pc cafe to login, bla bla"...
    Unverified          = 41   // account is unverified
    AccountDeleted      = 42   // inexistent or deleted account from whitelist
    AccountLocked       = 43   // too many wrong passwd attempts

    OutOfService2       = 47   // cannot connect at the moment
    AccountLockedSub    = 49   // account locked due sub pass fail
    TMS                 = 50   // [EP16+] time limit system; cannot start chatting

    EmailVerify         = 113  // please complete email verification first
    OGPTransfer         = 114  // transfer OGP information to EST accounts first
    AccountInactive     = 115  // logged out for a long time, activate first
    ChangingPasswd      = 116  // you can connect after changing password
)
