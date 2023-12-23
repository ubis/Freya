package packet

// Packet opcode values
const (
	CSCGetMyChar                       = 133
	CSCNewMyChar                       = 134
	CSCDelMyChar                       = 135
	CSCConnect2Svr                     = 140
	CSCVerifyLinks                     = 141
	CSCInitialize                      = 142
	CSCUnInitialize                    = 143
	CSCQuickLinkSet                    = 146
	CSCQuickLinkSwap                   = 147
	CSCGetServerTime                   = 148
	CSCItemLooting                     = 153
	CSCSkillToUser                     = 175
	CSCMoveBegin                       = 190
	CSCMoveEnd                         = 191
	CSCMoveChange                      = 192
	CSCMoveTile                        = 194
	CSCMessageEvent                    = 195
	NFYNewUserList                     = 200
	NFYDelUserList                     = 201
	NFYNewMonsterList                  = 202
	NFYDelMonsterList                  = 203
	NFYNewItemList                     = 204
	NFYDelItemList                     = 205
	NFYItemEquip                       = 206
	NFYItemUnEquip                     = 207
	NFYMoveBegin                       = 210
	NFYMoveEnd                         = 211
	NFYMoveChange                      = 212
	NFYMonsterMoveBegin                = 213
	NFYMonsterMoveEnd                  = 214
	NFYMessageEvent                    = 217
	NFYSkillToUser                     = 221
	NFYSystemMessage                   = 241
	CSCWarpCommand                     = 244
	CSCSkillToAction                   = 310
	NFYSkillToAction                   = 311
	CSCChangeStyle                     = 322
	NFYChangeStyle                     = 323
	CSCChargeInfo                      = 324
	CSCNewTargetUser                   = 350
	CSCChangeDirection                 = 391
	NFYChangeDirection                 = 392
	CSCKeyMoveBegin                    = 401
	CSCKeyMoveEnd                      = 402
	NFYKeyMoveBegin                    = 403
	NFYKeyMoveEnd                      = 404
	CSCKeyMoveChange                   = 405
	NFYKeyMoveChange                   = 406
	CSCServerEnv                       = 464
	CSCAccessoryEquip                  = 486
	CSCCheckUserPrivacyData            = 800
	CSCBackToCharacterLobby            = 985
	CSCSubPasswordSet                  = 1030
	CSCSubPasswordCheckRequest         = 1032
	CSCSubPasswordCheck                = 1034
	CSCSubPasswordFindRequest          = 1036
	CSCSubPasswordFind                 = 1038
	CSCSubPasswordDeleteRequest        = 1040
	CSCSubPasswordDelete               = 1042
	CSCSubPasswordChangeQARequest      = 1044
	CSCSubPasswordChangeQA             = 1046
	CSCSetCharacterSlotOrder           = 2001
	CSCChannelList                     = 2112
	CSCChannelChange                   = 2141
	CSCCharacterDeleteCheckSubPassword = 2160
	CSCStorageExchangeMove             = 2165
	CSCStorageItemSwap                 = 2166
	CSCStorageItemDrop                 = 2168
)

var opcodeNames = map[uint16]string{
	CSCGetMyChar:                       "GetMyChar",
	CSCNewMyChar:                       "NewMyChar",
	CSCDelMyChar:                       "DelMyChar",
	CSCConnect2Svr:                     "Connect2Svr",
	CSCVerifyLinks:                     "VerifyLinks",
	CSCInitialize:                      "Initialize",
	CSCUnInitialize:                    "UnInitialize",
	CSCQuickLinkSet:                    "QuickLinkSet",
	CSCQuickLinkSwap:                   "QuickLinkSwap",
	CSCGetServerTime:                   "GetServerTime",
	CSCItemLooting:                     "ItemLooting",
	CSCSkillToUser:                     "SkillToUser",
	CSCMoveBegin:                       "MoveBegin",
	CSCMoveEnd:                         "MoveEnd",
	CSCMoveChange:                      "MoveChange",
	CSCMoveTile:                        "MoveTile",
	CSCMessageEvent:                    "MessageEvent",
	NFYNewUserList:                     "NewUserList",
	NFYDelUserList:                     "DelUserList",
	NFYNewMonsterList:                  "NewMonsterList",
	NFYDelMonsterList:                  "DelMonsterList",
	NFYNewItemList:                     "NewItemList",
	NFYDelItemList:                     "DelItemList",
	NFYItemEquip:                       "ItemEquip",
	NFYItemUnEquip:                     "ItemUnEquip",
	NFYMoveBegin:                       "MoveBegin",
	NFYMoveEnd:                         "MoveEnd",
	NFYMoveChange:                      "MoveChange",
	NFYMonsterMoveBegin:                "MonsterMoveBegin",
	NFYMonsterMoveEnd:                  "MonsterMoveEnd",
	NFYMessageEvent:                    "MessageEvent",
	NFYSkillToUser:                     "SkillToUser",
	NFYSystemMessage:                   "SystemMessage",
	CSCWarpCommand:                     "WarpCommand",
	CSCSkillToAction:                   "SkillToAction",
	NFYSkillToAction:                   "SkillToAction",
	CSCChangeStyle:                     "ChangeStyle",
	NFYChangeStyle:                     "ChangeStyle",
	CSCChargeInfo:                      "ChargeInfo",
	CSCNewTargetUser:                   "NewTargetUser",
	CSCChangeDirection:                 "ChangeDirection",
	NFYChangeDirection:                 "ChangeDirection",
	CSCKeyMoveBegin:                    "KeyMoveBegin",
	CSCKeyMoveEnd:                      "KeyMoveEnd",
	NFYKeyMoveBegin:                    "KeyMoveBegin",
	NFYKeyMoveEnd:                      "KeyMoveEnd",
	CSCKeyMoveChange:                   "KeyMoveChange",
	NFYKeyMoveChange:                   "KeyMoveChange",
	CSCServerEnv:                       "ServerEnv",
	CSCAccessoryEquip:                  "AccessoryEquip",
	CSCCheckUserPrivacyData:            "CheckUserPrivacyData",
	CSCBackToCharacterLobby:            "BackToCharacterLobby",
	CSCSubPasswordSet:                  "SubPasswordSet",
	CSCSubPasswordCheckRequest:         "SubPasswordCheckRequest",
	CSCSubPasswordCheck:                "SubPasswordCheck",
	CSCSubPasswordFindRequest:          "SubPasswordFindRequest",
	CSCSubPasswordFind:                 "SubPasswordFind",
	CSCSubPasswordDeleteRequest:        "SubPasswordDeleteRequest",
	CSCSubPasswordDelete:               "SubPasswordDelete",
	CSCSubPasswordChangeQARequest:      "SubPasswordChangeQARequest",
	CSCSubPasswordChangeQA:             "SubPasswordChangeQA",
	CSCSetCharacterSlotOrder:           "SetCharacterSlotOrder",
	CSCChannelList:                     "ChannelList",
	CSCChannelChange:                   "ChannelChange",
	CSCCharacterDeleteCheckSubPassword: "CharacterDeleteCheckSubPassword",
	CSCStorageExchangeMove:             "StorageExchangeMove",
	CSCStorageItemSwap:                 "StorageItemSwap",
	CSCStorageItemDrop:                 "StorageItemDrop",
}

// Max in-game values
const (
	MaxCharacterSlot       = 5
	MaxCharacterNameLength = 16
)
