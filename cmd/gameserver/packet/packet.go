package packet

import (
	"github.com/ubis/Freya/cmd/gameserver/def"
	"github.com/ubis/Freya/cmd/gameserver/net"
	"github.com/ubis/Freya/share/log"
)

var g_ServerConfig = def.ServerConfig
var g_ServerSettings = def.ServerSettings
var g_PacketHandler = def.PacketHandler
var g_RPCHandler = def.RPCHandler

// Registers network packets
func RegisterPackets() {
	log.Info("Registering packets...")

	var pk = g_PacketHandler
	pk.Register(net.GETMYCHARTR, "GetMyChartr", GetMyChartr)
	pk.Register(net.NEWMYCHARTR, "NewMyChartr", NewMyChartr)
	pk.Register(net.DELMYCHARTR, "DelMyChartr", DelMyChartr)
	pk.Register(net.CONNECT2SVR, "Connect2Svr", Connect2Svr)
	pk.Register(net.VERIFYLINKS, "VerifyLinks", VerifyLinks)
	pk.Register(net.INITIALIZED, "Initialized", Initialized)
	pk.Register(net.UNINITIALZE, "Uninitialze", Uninitialze)
	pk.Register(net.GETSVRTIME, "GetSvrTime", GetSvrTime)
	pk.Register(net.MOVEBEGINED, "MoveBegined", MoveBegined)
	pk.Register(net.MOVEENDED00, "MoveEnded", MoveEnded)
	pk.Register(net.MOVECHANGED, "MoveChanged", MoveChanged)
	pk.Register(net.MOVETILEPOS, "MoveTilePos", MoveTilePos)
	pk.Register(net.MESSAGEEVNT, "MessageEvnt", MessageEvnt)
	pk.Register(net.NEWUSERLIST, "NewUserList", nil)
	pk.Register(net.DELUSERLIST, "DelUserList", nil)
	pk.Register(net.NFY_NEWMOBSLIST, "NotifyNewMobsList", nil)
	pk.Register(net.NFY_DELMOBSLIST, "NotifyDelMobsList", nil)
	pk.Register(net.NFY_MOVEBEGINED, "NotifyMoveBegined", nil)
	pk.Register(net.NFY_MOVEENDED00, "NotifyMoveEnded", nil)
	pk.Register(net.NFY_MOVECHANGED, "NotifyMoveChanged", nil)
	pk.Register(net.NFY_MOBSMOVEBGN, "NotifyMobsMoveBegin", nil)
	pk.Register(net.NFY_MOBSMOVEEND, "NotifyMobsMoveEnd", nil)
	pk.Register(net.NFY_MESSAGEEVNT, "NotifyMessageEvnt", nil)
	pk.Register(net.SYSTEMMESSG, "SystemMessg", nil)
	pk.Register(net.WARPCOMMAND, "WarpCommand", WarpCommand)
	pk.Register(net.CHARGEINFO, "ChargeInfo", ChargeInfo)
	pk.Register(net.CHANGEDIRECTION, "ChangeDirection", ChangeDirection)
	pk.Register(net.NFY_CHANGEDIRECTION, "NotifyChangeDirection", nil)
	pk.Register(net.KEYMOVEBEGINED, "KeyMoveBegined", KeyMoveBegined)
	pk.Register(net.KEYMOVEENDED00, "KeyMoveEnded", KeyMoveEnded)
	pk.Register(net.NFY_KEYMOVEBEGINED, "NotifyKeyMoveBegined", nil)
	pk.Register(net.NFY_KEYMOVEENDED00, "NotifyKeyMoveEnded", nil)
	pk.Register(net.KEYMOVECHANGED, "KeyMoveChanged", KeyMoveChanged)
	pk.Register(net.NFY_KEYMOVECHANGED, "NotifyKeyMoveChanged", nil)
	pk.Register(net.SERVERENV, "ServerEnv", ServerEnv)
	pk.Register(net.CHECK_USR_PDATA, "CheckUserPrivacyData", CheckUserPrivacyData)
	pk.Register(net.BACK_TO_CHAR_LOBBY, "BackToCharLobby", BackToCharLobby)
	pk.Register(net.SUBPW_SET, "SubPasswordSet", SubPasswordSet)
	pk.Register(net.SUBPW_CHECK_REQ, "SubPasswordCheckRequest", SubPasswordCheckRequest)
	pk.Register(net.SUBPW_CHECK, "SubPasswordCheck", SubPasswordCheck)
	pk.Register(net.SUBPW_FIND_REQ, "SubPasswordFindRequest", SubPasswordFindRequest)
	pk.Register(net.SUBPW_FIND, "SubPasswordFind", SubPasswordFind)
	pk.Register(net.SUBPW_DEL_REQ, "SubPasswordDelRequest", SubPasswordDelRequest)
	pk.Register(net.SUBPW_DEL, "SubPasswordDel", SubPasswordDel)
	pk.Register(net.SUBPW_CHG_QA_REQ,
		"SubPasswordChangeQARequest", SubPasswordChangeQARequest)
	pk.Register(net.SUBPW_CHG_QA, "SubPasswordChangeQA", SubPasswordChangeQA)
	pk.Register(net.SET_CHAR_SLOT_ORDER, "SetCharacterSlotOrder", SetCharacterSlotOrder)
	pk.Register(net.CHANNEL_LIST, "ChannelList", ChannelList)
	pk.Register(net.CHANNEL_CHANGE, "ChannelChange", ChannelChange)
	pk.Register(net.CHAR_DEL_CHK_SUBPW,
		"CharacterDeleteCheckSubPassword", CharacterDeleteCheckSubPassword)
}
