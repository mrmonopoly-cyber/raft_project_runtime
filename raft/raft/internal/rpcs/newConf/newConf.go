package newConf

import (
	"log"
	"raft/internal/node"
	"raft/internal/raft_log"
	"raft/internal/raftstate"
	"raft/internal/rpcs"
	ClientReturnValue "raft/internal/rpcs/clientReturnValue"
	"raft/pkg/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"

	"google.golang.org/protobuf/proto"
)

type NewConf struct {
    pMex protobuf.ChangeConfReq
}

func NewnewConfRPC(op protobuf.AdminOp, conf []string) rpcs.Rpc {
    return &NewConf{
        pMex: protobuf.ChangeConfReq{
            Op: op,
            Conf: &protobuf.ClusterConf{
                Conf: conf,
            },
        },
    }
}

// Manage implements rpcs.Rpc.
func (this *NewConf) Execute(state raftstate.State, sender node.Node) rpcs.Rpc {
    var exitSucess rpcs.Rpc = ClientReturnValue.NewclientReturnValueRPC(protobuf.STATUS_SUCCESS,"")
    var myPrivateIp = state.GetMyIp(raftstate.PRI)

    switch this.pMex.Op{
    case protobuf.AdminOp_CHANGE_CONF_NEW:
        if state.GetConfig() == nil && myPrivateIp == *this.pMex.Conf.Leader{
            var newConfEntry = protobuf.LogEntry{
                    Term: state.GetTerm(),
                    OpType: protobuf.Operation_JOIN_CONF_ADD,
                    Description: "New configuration of the new cluster",
                }
            var newConfAppEntry = state.NewLogInstance(&newConfEntry)

            for _, v := range this.pMex.Conf.GetConf(){
                var ele string = v + " "
                newConfAppEntry.Entry.Payload = append(newConfAppEntry.Entry.Payload,ele...)
            }

            state.AppendEntries([]*raft_log.LogInstance{newConfAppEntry})
            state.SetRole(raftstate.LEADER)
            return exitSucess
        }
        var failureDescr = `cluster already created and settend, New conf can only be applied 
        when the cluster is still yet to be configured in the first place,
        in all other cases this call will do noting.
        If you want to change the conf use ADD or REM`
        log.Println(failureDescr)
        return ClientReturnValue.NewclientReturnValueRPC(protobuf.STATUS_FAILURE, failureDescr)
    case protobuf.AdminOp_CHANGE_CONF_ADD:
        panic("not implemented")
    case protobuf.AdminOp_CHANGE_CONF_REM:
        panic("not implemented")
    default:
        var failDescr = "invalid new conf request type: " + this.pMex.Op.String()
        log.Println(failDescr)
        return ClientReturnValue.NewclientReturnValueRPC(protobuf.STATUS_FAILURE, failDescr)
    }
}

// ToString implements rpcs.Rpc.
func (this *NewConf) ToString() string {
    return this.pMex.String()
}

func (this *NewConf) Encode() ([]byte, error) {
    var mess []byte
    var err error

    mess, err = proto.Marshal(&(*this).pMex)
    if err != nil {
        log.Panicln("error in Encoding Request Vote: ", err)
    }

	return mess, err
}
func (this *NewConf) Decode(b []byte) error {
	err := proto.Unmarshal(b,&this.pMex)
    if err != nil {
        log.Panicln("error in Encoding Request Vote: ", err)
    }
	return err
}
