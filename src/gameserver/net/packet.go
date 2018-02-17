package net

import (
	"share/log"
	"share/network/packet"
	"share/rpc"
)

// Packet structure
type Packet struct {
	packet.List

	RPC *rpc.Client

	ServerID   int
	GroupID    int
	ServerType int
}

// Register network packets
func (p *Packet) Register() {
	// register packets
	log.Info("Registering packets...")

	p.Add(GetMyChartr, "GetMyChartr", p.GetMyChartr)
	p.Add(NewMyChartr, "NewMyChartr", p.NewMyChartr)
	p.Add(DelMyChartr, "DelMyChartr", p.DelMyChartr)
	p.Add(Connect2Svr, "Connect2Svr", p.Connect2Svr)
	p.Add(VerifyLinks, "VerifyLinks", p.VerifyLinks)
	p.Add(GetSvrTime, "GetSvrTime", p.GetSvrTime)
	p.Add(SystemMessg, "SystemMessg", nil)
	p.Add(ChargeInfo, "ChargeInfo", p.ChargeInfo)
	p.Add(ServerEnv, "ServerEnv", p.ServerEnv)
	p.Add(CheckUserPrivacyData,
		"CheckUserPrivacyData", p.CheckUserPrivacyData)
	p.Add(SubPasswordSet, "SubPasswordSet", p.SubPasswordSet)
	p.Add(SubPasswordCheckRequest,
		"SubPasswordCheckRequest", p.SubPasswordCheckRequest)
	p.Add(SubPasswordCheck, "SubPasswordCheck", p.SubPasswordCheck)
	p.Add(SubPasswordFindRequest,
		"SubPasswordFindRequest", p.SubPasswordFindRequest)
	p.Add(SubPasswordFind, "SubPasswordFind", p.SubPasswordFind)
	p.Add(SubPasswordDelRequest,
		"SubPasswordDelRequest", p.SubPasswordDelRequest)
	p.Add(SubPasswordDel, "SubPasswordDel", p.SubPasswordDel)
	p.Add(SubPasswordChangeQARequest,
		"SubPasswordChangeQARequest", p.SubPasswordChangeQARequest)
	p.Add(SubPasswordChangeQA,
		"SubPasswordChangeQA", p.SubPasswordChangeQA)
	p.Add(SetCharacterSlotOrder,
		"SetCharacterSlotOrder", p.SetCharacterSlotOrder)
	p.Add(CharacterDeleteCheckSubPassword,
		"CharacterDeleteCheckSubPassword", p.CharacterDeleteCheckSubPassword)
}
