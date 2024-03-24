package genericmessage

import (
	"log"
	"raft/internal/rpcs"
	"raft/internal/rpcs/AppendEntryRpc"
	"raft/pkg/rpcEncoding/out/protobuf"

	"google.golang.org/protobuf/proto"
)

func Decode(raw_mex []byte) (*rpcs.Rpc){
    var genericMex protobuf.Entry
    var err error
    var outRpc rpcs.Rpc

    err = proto.Unmarshal(raw_mex,&genericMex)
    if err != nil {
        log.Panicln("error decoding generic message")
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
