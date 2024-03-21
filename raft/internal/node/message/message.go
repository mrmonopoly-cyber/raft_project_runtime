package message

import (
	"encoding/json"
	"log"
	"raft/internal/messages"
	"raft/internal/messages/AppendEntryRPC"
	"raft/internal/messages/AppendEntryResponse"
	"raft/internal/messages/CopyStateRPC"
	"raft/internal/messages/RequestVoteRPC"
	"raft/internal/messages/RequestVoteResponse"
	p "raft/pkg/protobuf"

	"google.golang.org/protobuf/proto"
)

type EnumType int

const (
	APPEND_ENTRY EnumType = iota
	REQUEST_VOTE
	APPEND_RESPONSE
	VOTE_RESPONSE
	COPY_STATE
)

type Message struct {
	Kind      p.Ty
	Payload []byte
}

func (m *Message) ToMessage(b []byte) {
	json.Unmarshal(b, m)
}

func (this *Message) ToByte() []byte {
	var t p.Ty

	switch this.Kind {
	case p.Ty_APPEND_ENTRY:
		t = p.Ty_APPEND_ENTRY
	case p.Ty_APPEND_RESPONSE:
		t = p.Ty_APPEND_RESPONSE
	case p.Ty_REQUEST_VOTE:
		t = p.Ty_REQUEST_VOTE
	case p.Ty_VOTE_RESPONSE:
		t = p.Ty_VOTE_RESPONSE
	case p.Ty_COPY_STATE:
		t = p.Ty_COPY_STATE
	}

	mess := &p.Message{
		Kind:    t.Enum(),
		Payload: this.Payload,
	}
	mex, _ := proto.Marshal(mess)
	return mex
}

func NewMessage(data []byte) *Message {
	var mess p.Message
    var mex_type p.Ty 

    var err error =proto.Unmarshal(data, &mess)
    if err != nil {
        log.Panic("error Unmarshal", err)
    }
    mex_type = *mess.Kind

    log.Println("mex type :", mex_type)
    log.Println("mex payload:", mess.Payload)

	return &Message{
		Kind:       mex_type,
		Payload: mess.Payload,
	}
}

func FromRpc(rpc messages.Rpc) Message {
	var ty p.Ty 
	switch rpc.(type) {
	case *AppendEntryRPC.AppendEntryRPC:
		ty = p.Ty_APPEND_ENTRY
	case *RequestVoteRPC.RequestVoteRPC:
		ty = p.Ty_REQUEST_VOTE
	case *AppendEntryResponse.AppendEntryResponse:
		ty = p.Ty_APPEND_RESPONSE
	case *RequestVoteResponse.RequestVoteResponse:
		ty = p.Ty_VOTE_RESPONSE
	case *CopyStateRPC.CopyStateRPC:
		ty = p.Ty_COPY_STATE
	}

	var enc []byte
	var err error
	enc, err = rpc.Encode()

	if err != nil {
		log.Panicf("Error in encoding rpc: %s", err)
	}

	return Message{
		Kind:    ty,
		Payload: enc,
	}
}

func (this *Message) ToRpc() *messages.Rpc {

	var recRpc messages.Rpc
	switch this.Kind {
	case p.Ty_APPEND_ENTRY:
		var rpc AppendEntryRPC.AppendEntryRPC
		recRpc = &rpc
	case p.Ty_APPEND_RESPONSE:
		var rpc AppendEntryResponse.AppendEntryResponse
		recRpc = &rpc
	case p.Ty_REQUEST_VOTE:
		var rpc RequestVoteRPC.RequestVoteRPC
		recRpc = &rpc
	case p.Ty_VOTE_RESPONSE:
		var rpc RequestVoteResponse.RequestVoteResponse
		recRpc = &rpc
	case p.Ty_COPY_STATE:
		var rpc CopyStateRPC.CopyStateRPC
		recRpc = &rpc
	default:
		panic("impossible case invalid type RPC " + string(rune(this.Kind)))
	}
	recRpc.Decode(this.Payload)

	return &recRpc
}
