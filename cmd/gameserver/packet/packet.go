package packet

import (
	"github.com/ubis/Freya/share/log"
	"github.com/ubis/Freya/share/network"
	"github.com/ubis/Freya/share/script"
)

type PacketFunc func(*Session, *network.Reader)

func register(p *network.PacketHandler, opcode uint16, method PacketFunc) {
	name, ok := opcodeNames[opcode]
	if !ok {
		name = "Unknown"
	}

	if method == nil {
		p.Register(opcode, name, nil)
		return
	}

	p.Register(opcode, name, func(s *network.Session, r *network.Reader) {
		session, ok := s.Retrieve().(*Session)
		if !ok {
			log.Error("Unable to parse client session!")
			return
		}

		method(session, r)
	})
}

func verifyState(session *Session, state SessionState, opcode uint16) bool {
	currentState := session.GetState()

	if currentState != state {
		session.LogFatalf("Invalid client state for %d (%s) packet "+
			"[need: %d ; have: %d]",
			opcode, opcodeNames[opcode], state, currentState)
		return false
	}

	return true
}

// Registers network packets
func RegisterPackets(h *network.PacketHandler) {
	log.Info("Registering packets...")

	register(h, CSCGetMyChar, GetMyChar)
	register(h, CSCNewMyChar, NewMyChar)
	register(h, CSCDelMyChar, DelMyChar)
	register(h, CSCConnect2Svr, Connect2Svr)
	register(h, CSCVerifyLinks, VerifyLinks)
	register(h, CSCInitialize, Initialize)
	register(h, CSCUnInitialize, UnInitialize)
	register(h, CSCQuickLinkSet, QuickLinkSet)
	register(h, CSCQuickLinkSwap, QuickLinkSwap)
	register(h, CSCGetServerTime, GetServerTime)
	register(h, CSCItemLooting, ItemLooting)
	register(h, CSCSkillToUser, SkillToUser)
	register(h, CSCMoveBegin, MoveBegin)
	register(h, CSCMoveEnd, MoveEnd)
	register(h, CSCMoveChange, MoveChange)
	register(h, CSCMoveTile, MoveTile)
	register(h, CSCMessageEvent, MessageEvent)
	register(h, NFYNewUserList, nil)
	register(h, NFYDelUserList, nil)
	register(h, NFYItemEquip, nil)
	register(h, NFYItemUnEquip, nil)
	register(h, NFYNewMonsterList, nil)
	register(h, NFYDelMonsterList, nil)
	register(h, NFYNewItemList, nil)
	register(h, NFYDelItemList, nil)
	register(h, NFYMoveBegin, nil)
	register(h, NFYMoveEnd, nil)
	register(h, NFYMoveChange, nil)
	register(h, NFYMonsterMoveBegin, nil)
	register(h, NFYMonsterMoveEnd, nil)
	register(h, NFYMessageEvent, nil)
	register(h, NFYSkillToUser, nil)
	register(h, NFYSystemMessage, nil)
	register(h, CSCWarpCommand, WarpCommand)
	register(h, CSCSkillToAction, SkillToAction)
	register(h, NFYSkillToAction, nil)
	register(h, CSCChangeStyle, ChangeStyle)
	register(h, NFYChangeStyle, nil)
	register(h, CSCChargeInfo, ChargeInfo)
	register(h, CSCNewTargetUser, NewTargetUser)
	register(h, CSCChangeDirection, ChangeDirection)
	register(h, NFYChangeDirection, nil)
	register(h, CSCKeyMoveBegin, KeyMoveBegin)
	register(h, CSCKeyMoveEnd, KeyMoveEnd)
	register(h, NFYKeyMoveBegin, nil)
	register(h, NFYKeyMoveEnd, nil)
	register(h, CSCKeyMoveChange, KeyMoveChange)
	register(h, NFYKeyMoveChange, nil)
	register(h, CSCServerEnv, ServerEnv)
	register(h, CSCAccessoryEquip, AccessoryEquip)
	register(h, CSCCheckUserPrivacyData, CheckUserPrivacyData)
	register(h, CSCBackToCharacterLobby, BackToCharacterLobby)
	register(h, CSCSubPasswordSet, SubPasswordSet)
	register(h, CSCSubPasswordCheckRequest, SubPasswordCheckRequest)
	register(h, CSCSubPasswordCheck, SubPasswordCheck)
	register(h, CSCSubPasswordFindRequest, SubPasswordFindRequest)
	register(h, CSCSubPasswordFind, SubPasswordFind)
	register(h, CSCSubPasswordDeleteRequest, SubPasswordDeleteRequest)
	register(h, CSCSubPasswordDelete, SubPasswordDelete)
	register(h, CSCSubPasswordChangeQARequest, SubPasswordChangeQARequest)
	register(h, CSCSubPasswordChangeQA, SubPasswordChangeQA)
	register(h, CSCSetCharacterSlotOrder, SetCharacterSlotOrder)
	register(h, CSCChannelList, ChannelList)
	register(h, CSCChannelChange, ChannelChange)
	register(h, CSCCharacterDeleteCheckSubPassword, CharacterDeleteCheckSubPassword)
	register(h, CSCStorageExchangeMove, StorageExchangeMove)
	register(h, CSCStorageItemSwap, StorageItemSwap)
	register(h, CSCStorageItemDrop, StorageItemDrop)
}

func RegisterFunc() {
	script.RegisterFunc("sendClientPacket", sessionPacketFunc{})
	script.RegisterFunc("sendClientMessage", clientMessageFunc{Fn: SendMessage})

	script.RegisterFunc("getPlayerLevel", playerGetLevelFunc{Fn: GetPlayerLevel})
	script.RegisterFunc("setPlayerLevel", playerSetLevelFunc{Fn: SetPlayerLevel})
	script.RegisterFunc("getPlayerPosition", playerPositionFunc{})
	script.RegisterFunc("dropItem", playerDropItemFunc{})
}
