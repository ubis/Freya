package packet

// Packet Type ID's
const (
	GETMYCHARTR         = 133
	NEWMYCHARTR         = 134
	DELMYCHARTR         = 135
	CONNECT2SVR         = 140
	VERIFYLINKS         = 141
	INITIALIZED         = 142
	UNINITIALZE         = 143
	GETSVRTIME          = 148
	MOVEBEGINED         = 190
	MOVEENDED00         = 191
	MOVECHANGED         = 192
	MOVETILEPOS         = 194
	MESSAGEEVNT         = 195
	NEWUSERLIST         = 200
	DELUSERLIST         = 201
	NFY_MOVEBEGINED     = 210
	NFY_MOVEENDED00     = 211
	NFY_MOVECHANGED     = 212
	NFY_MESSAGEEVNT     = 217
	SYSTEMMESSG         = 241
	WARPCOMMAND         = 244
	CHARGEINFO          = 324
	CHANGEDIRECTION     = 391
	NFY_CHANGEDIRECTION = 392
	KEYMOVEBEGINED      = 401
	KEYMOVEENDED00      = 402
	NFY_KEYMOVEBEGINED  = 403
	NFY_KEYMOVEENDED00  = 404
	KEYMOVECHANGED      = 405
	NFY_KEYMOVECHANGED  = 406
	SERVERENV           = 464
	CHECK_USR_PDATA     = 800
	BACK_TO_CHAR_LOBBY  = 985
	SUBPW_SET           = 1030
	SUBPW_CHECK_REQ     = 1032
	SUBPW_CHECK         = 1034
	SUBPW_FIND_REQ      = 1036
	SUBPW_FIND          = 1038
	SUBPW_DEL_REQ       = 1040
	SUBPW_DEL           = 1042
	SUBPW_CHG_QA_REQ    = 1044
	SUBPW_CHG_QA        = 1046
	SET_CHAR_SLOT_ORDER = 2001
	CHANNEL_LIST        = 2112
	CHANNEL_CHANGE      = 2141
	CHAR_DEL_CHK_SUBPW  = 2160
)
