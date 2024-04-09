package genericmessage

import (
	"log"
	"raft/internal/rpcs"
	"raft/internal/rpcs/AppendEntryRpc"
	"raft/internal/rpcs/AppendResponse"
	"raft/internal/rpcs/RequestVoteRPC"
	"raft/internal/rpcs/RequestVoteResponse"

	"raft/pkg/rpcEncoding/out/protobuf"

	"google.golang.org/protobuf/proto"
)

func Decode(raw_mex []byte) (*rpcs.Rpc){
    var genericMex protobuf.Entry
    var err error
    var outRpc rpcs.Rpc

    err = proto.Unmarshal(raw_mex,&genericMex)
    if err != nil {
        log.Panicln("error decoding generic message: ", string(raw_mex))
        return nil
    }
    //log.Println("generic message Decoded with type :", genericMex.GetOpType())

    var payload []byte = genericMex.GetPayload()

    switch genericMex.GetOpType(){
    case protobuf.MexType_APPEND_ENTRY:
        var appendEntry AppendEntryRpc.AppendEntryRpc
        outRpc = &appendEntry
    case protobuf.MexType_APPEND_ENTRY_RESPONSE:
        var appendResponse AppendResponse.AppendResponse
        outRpc = &appendResponse
    case protobuf.MexType_REQUEST_VOTE:
        var requestVote RequestVoteRPC.RequestVoteRPC
        outRpc = &requestVote
    case protobuf.MexType_REQUEST_VOTE_RESPONSE:
        var requestVoteResponse RequestVoteResponse.RequestVoteResponse
        outRpc = &requestVoteResponse
    default:
        log.Panicln("rpc type not recognize in decoing generic message")
    }
    
    outRpc.Decode(payload)
    return &outRpc

}

func Encode(mex *rpcs.Rpc) ([]byte,error){
    var err error
    var rawByte []byte
    var rawByteToSend []byte
    var genericMessage protobuf.Entry

    rawByte, err= (*mex).Encode()
    genericMessage.Payload = rawByte

    if err != nil {
        log.Panicln("error encoding this message :", (*mex).ToString())
    }

    switch (*mex).(type){
    case *AppendEntryRpc.AppendEntryRpc:
        genericMessage.OpType = protobuf.MexType_APPEND_ENTRY
    case *AppendResponse.AppendResponse:
        genericMessage.OpType = protobuf.MexType_APPEND_ENTRY_RESPONSE
    case *RequestVoteRPC.RequestVoteRPC:
        genericMessage.OpType = protobuf.MexType_REQUEST_VOTE
    case *RequestVoteResponse.RequestVoteResponse:
        genericMessage.OpType = protobuf.MexType_REQUEST_VOTE_RESPONSE
    default:
        log.Panicln("rpc not recognize: ", (*mex).ToString())
    }

    rawByteToSend,err = proto.Marshal(&genericMessage)
    if err != nil {
        log.Panicln("failed in serializing the rpc AppendEntry")
    }

    return rawByteToSend,nil

    
}
