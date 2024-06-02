package genericmessage

import (
	"errors"
	"raft/internal/rpcs"
	"raft/internal/rpcs/AppendEntryRpc"
	"raft/internal/rpcs/AppendResponse"
	"raft/internal/rpcs/RequestVoteRPC"
	"raft/internal/rpcs/RequestVoteResponse"
	"raft/internal/rpcs/UpdateNode"
	"raft/internal/rpcs/UpdateNodeResp"
	"raft/internal/rpcs/redirection"

	"raft/pkg/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"

	"google.golang.org/protobuf/proto"
)

func Decode(raw_mex []byte) (rpcs.Rpc,error){
    var genericMex protobuf.Entry
    var err error
    var outRpc rpcs.Rpc

    err = proto.Unmarshal(raw_mex,&genericMex)
    if err != nil {
        return nil, errors.New("error decoding generic message: "+ string(raw_mex))
    }
    //log.Println("generic message Decoded with type :", genericMex.GetOpType())

    var payload []byte = genericMex.GetPayload()

    switch genericMex.GetOpType(){
    case protobuf.MexType_APPEND_ENTRY:
        outRpc = &AppendEntryRpc.AppendEntryRpc{}
    case protobuf.MexType_APPEND_ENTRY_RESPONSE:
        outRpc = &AppendResponse.AppendResponse{}
    case protobuf.MexType_REQUEST_VOTE:
        outRpc = &RequestVoteRPC.RequestVoteRPC{}
    case protobuf.MexType_REQUEST_VOTE_RESPONSE:
        outRpc = &RequestVoteResponse.RequestVoteResponse{}
    case protobuf.MexType_UPDATE_NODE:
        outRpc= &UpdateNode.UpdateNode{}
    case protobuf.MexType_UPDATE_NODE_RESP:
        outRpc = &UpdateNodeResp.UpdateNodeResp{}
    case protobuf.MexType_REDIRECTION:
        outRpc = &Redirection.Redirection{}
    default:
        return nil, errors.New("rpc not recognized")
    }
    
    outRpc.Decode(payload)
    return outRpc,nil

}

func Encode(mex rpcs.Rpc) ([]byte,error){
    var err error
    var rawByte []byte
    var rawByteToSend []byte
    var genericMessage protobuf.Entry

    rawByte, err= mex.Encode()
    genericMessage.Payload = rawByte

    if err != nil {
        return nil,errors.New("error encoding this message :"+ mex.ToString())
    }

    switch mex.(type){
    case *AppendEntryRpc.AppendEntryRpc:
        genericMessage.OpType = protobuf.MexType_APPEND_ENTRY
    case *AppendResponse.AppendResponse:
        genericMessage.OpType = protobuf.MexType_APPEND_ENTRY_RESPONSE
    case *RequestVoteRPC.RequestVoteRPC:
        genericMessage.OpType = protobuf.MexType_REQUEST_VOTE
    case *RequestVoteResponse.RequestVoteResponse:
        genericMessage.OpType = protobuf.MexType_REQUEST_VOTE_RESPONSE
    case *UpdateNode.UpdateNode:
        genericMessage.OpType = protobuf.MexType_UPDATE_NODE
    case *UpdateNodeResp.UpdateNodeResp:
        genericMessage.OpType = protobuf.MexType_UPDATE_NODE_RESP
    case *Redirection.Redirection:
        genericMessage.OpType = protobuf.MexType_REDIRECTION
    default:
        return nil, errors.New("rpc not recognized")
    }

    rawByteToSend,err = proto.Marshal(&genericMessage)
    if err != nil {
        return nil,errors.New("failed in serializing the rpc AppendEntry")
    }

    return rawByteToSend,nil

    
}
