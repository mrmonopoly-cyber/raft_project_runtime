package genericmessage

import (
	"log"
	"raft/internal/rpcs"
	"raft/internal/rpcs/AppendEntryRpc"
    "raft/internal/rpcs/RequestVoteRPC"

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
    log.Println("generic message Decoded with type :", genericMex.GetPayload())
    switch genericMex.GetOpType(){
    case protobuf.MexType_APPEND_ENTRY:
        var appendEntry AppendEntryRpc.AppendEntryRpc
        appendEntry.Decode(genericMex.GetPayload())
        outRpc = &appendEntry
    }
    
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
    case *RequestVoteRPC.RequestVoteRPC:
        genericMessage.OpType = protobuf.MexType_REQUEST_VOTE
    default:
        log.Panicln("rpc not recognize: ", (*mex).ToString())
    }

    rawByteToSend,err = proto.Marshal(&genericMessage)
    if err != nil {
        log.Panicln("failed in serializing the rpc AppendEntry")
    }

    return rawByteToSend,nil

    
}
