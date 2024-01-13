package packet

import (
	"time"

	"github.com/ubis/Freya/share/network"
)

// MoveBegin Packet
func MoveBegin(session *Session, reader *network.Reader) {
	if !verifyState(session, StateInGame, reader.Type) {
		return
	}

	startX := byte(reader.ReadInt16())
	startY := byte(reader.ReadInt16())
	endX := byte(reader.ReadInt16())
	endY := byte(reader.ReadInt16())
	_ = reader.ReadInt16() // pnt x
	_ = reader.ReadInt16() // pnt y
	_ = reader.ReadInt16() // world map

	pkt := network.NewWriter(NFYMoveBegin)
	pkt.WriteInt32(session.Character.Id)
	pkt.WriteUint32(time.Now().Unix()) // move begin time
	pkt.WriteInt16(startX)
	pkt.WriteInt16(startY)
	pkt.WriteInt16(endX)
	pkt.WriteInt16(endY)

	session.Character.SetMovement(startX, startY, endX, endY)

	session.Broadcast(pkt)
}

// MoveEnd Packet
func MoveEnd(session *Session, reader *network.Reader) {
	if !verifyState(session, StateInGame, reader.Type) {
		return
	}

	pntX := byte(reader.ReadInt16())
	pntY := byte(reader.ReadInt16())

	pkt := network.NewWriter(NFYMoveEnd)
	pkt.WriteInt32(session.Character.Id)
	pkt.WriteInt16(pntX)
	pkt.WriteInt16(pntY)

	session.Character.SetPosition(pntX, pntY)
	session.Character.SetMovement(pntX, pntY, pntX, pntY)

	session.Broadcast(pkt)
}

// MoveChange Packet
func MoveChange(session *Session, reader *network.Reader) {
	if !verifyState(session, StateInGame, reader.Type) {
		return
	}

	startX := byte(reader.ReadInt16())
	startY := byte(reader.ReadInt16())
	endX := byte(reader.ReadInt16())
	endY := byte(reader.ReadInt16())
	_ = reader.ReadInt16() // pnt x
	_ = reader.ReadInt16() // pnt y
	_ = reader.ReadInt16() // world map

	pkt := network.NewWriter(NFYMoveChange)
	pkt.WriteInt32(session.Character.Id)
	pkt.WriteUint32(time.Now().Unix()) // move begin time
	pkt.WriteInt16(startX)
	pkt.WriteInt16(startY)
	pkt.WriteInt16(endX)
	pkt.WriteInt16(endY)

	session.Character.SetMovement(startX, startY, endX, endY)

	session.Broadcast(pkt)
}

// MoveTile packet
func MoveTile(session *Session, reader *network.Reader) {
	if !verifyState(session, StateInGame, reader.Type) {
		return
	}

	_ = reader.ReadInt16()           // curr x
	_ = reader.ReadInt16()           // curr y
	pntX := byte(reader.ReadInt16()) // pnt x
	pntY := byte(reader.ReadInt16()) // pnt y

	session.Character.SetPosition(pntX, pntY)
	session.AdjustCell()
}

// ChangeDirection Packet
func ChangeDirection(session *Session, reader *network.Reader) {
	if !verifyState(session, StateInGame, reader.Type) {
		return
	}

	direction := reader.ReadUint32() // float

	pkt := network.NewWriter(NFYChangeDirection)
	pkt.WriteInt32(session.Character.Id)
	pkt.WriteUint32(direction) // float

	session.Broadcast(pkt)
}

// KeyMoveBegin Packet
func KeyMoveBegin(session *Session, reader *network.Reader) {
	if !verifyState(session, StateInGame, reader.Type) {
		return
	}

	startX := reader.ReadUint32() // float
	startY := reader.ReadUint32() // float
	endX := reader.ReadUint32()   // float
	endY := reader.ReadUint32()   // float
	_ = reader.ReadUint32()       // pnt x
	_ = reader.ReadUint32()       // pnt y
	dir := reader.ReadByte()

	pkt := network.NewWriter(NFYKeyMoveBegin)
	pkt.WriteInt32(session.Character.Id)
	pkt.WriteUint32(time.Now().Unix()) // move begin time
	pkt.WriteUint32(startX)
	pkt.WriteUint32(startY)
	pkt.WriteUint32(endX)
	pkt.WriteUint32(endY)
	pkt.WriteByte(dir)

	session.Broadcast(pkt)
}

// KeyMoveEnd Packet
func KeyMoveEnd(session *Session, reader *network.Reader) {
	if !verifyState(session, StateInGame, reader.Type) {
		return
	}

	pntX := reader.ReadUint32() // float
	pntY := reader.ReadUint32() // float

	pkt := network.NewWriter(NFYKeyMoveEnd)
	pkt.WriteInt32(session.Character.Id)
	pkt.WriteInt32(pntX)
	pkt.WriteInt32(pntY)

	session.Broadcast(pkt)
}

// KeyMoveChange Packet
func KeyMoveChange(session *Session, reader *network.Reader) {
	if !verifyState(session, StateInGame, reader.Type) {
		return
	}

	startX := reader.ReadUint32() // float
	startY := reader.ReadUint32() // float
	endX := reader.ReadUint32()   // float
	endY := reader.ReadUint32()   // float
	_ = reader.ReadUint32()       // pnt x
	_ = reader.ReadUint32()       // pnt y
	dir := reader.ReadByte()

	pkt := network.NewWriter(NFYKeyMoveChange)
	pkt.WriteInt32(session.Character.Id)
	pkt.WriteUint32(time.Now().Unix()) // move begin time
	pkt.WriteUint32(startX)
	pkt.WriteUint32(startY)
	pkt.WriteUint32(endX)
	pkt.WriteUint32(endY)
	pkt.WriteByte(dir)

	session.Broadcast(pkt)
}
