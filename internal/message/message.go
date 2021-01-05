package message

import (
	"remote/api"
	flatbuffers "github.com/google/flatbuffers/go"
	"encoding/binary"
)


const (
	ActionNone     = 0
	ActionRequest  = 1
	ActionResponse = 2
)



type Message struct {
	Action int
	Subject string
	Content string
}

func FromBytes(in []byte) Message {
	m := api.GetRootAsMessageBuffer(in, 0)
	mS := Message{int(m.Action()), string(m.Subject()), string(m.Content())}
	return mS
}

func ToBytes(in Message) []byte {
	b := flatbuffers.NewBuilder(1024)

	bSubj := b.CreateString(in.Subject)
	bCont := b.CreateString(in.Content)

	api.MessageBufferStart(b)
	api.MessageBufferAddSubject(b, bSubj)
	api.MessageBufferAddContent(b, bCont)
	api.MessageBufferAddAction(b, api.Action(in.Action))
	msg := api.MessageBufferEnd(b)

	b.Finish(msg)
	buf := b.FinishedBytes()
	return buf
}

func IntToBinary(in int) []byte {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(in))
	return buf
}

func BinaryToInt(in []byte) int {
	return int(binary.LittleEndian.Uint32(in))
}