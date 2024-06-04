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
        if state.GetConf() == nil && myPrivateIp == *this.pMex.Conf.Leader{
            var newEntryBaseEntry = protobuf.LogEntry{
                Term: state.GetTerm(),
                OpType: protobuf.Operation_JOIN_CONF_ADD,
            }

            for _, v := range this.pMex.Conf.Conf {
                newEntryBaseEntry.Payload = append(newEntryBaseEntry.Payload, v...)
                newEntryBaseEntry.Payload = append(newEntryBaseEntry.Payload, raft_log.SEPARATOR...)
            }

            var newConfLog = state.NewLogInstance(&newEntryBaseEntry,func() {
                var commit = protobuf.LogEntry{
                    Term: state.GetTerm(),
                    OpType: protobuf.Operation_COMMIT_CONFIG_ADD,
                }

                state.AppendEntry(state.NewLogInstance(&commit,nil))
            }) 
            log.Println("appending log entry: ",newConfLog)

            state.SetRole(raftstate.LEADER)
            state.AppendEntry(newConfLog)

            return exitSucess
        }
        var failureDescr = `cluster already created and settend, New conf can only be applied 
        when the cluster is still yet to be configured in the first place,
        in all other cases this call will do noting.
        If you want to change the conf use ADD or REM`
        log.Println(failureDescr)
        return ClientReturnValue.NewclientReturnValueRPC(protobuf.STATUS_FAILURE, failureDescr)
    case protobuf.AdminOp_CHANGE_CONF_ADD:
        if state.GetRole() == raftstate.LEADER{
            panic("not implemented")
            // return exitSucess
        }
        var failureMex = "i'm not leader, i cannot change conf"
        log.Println(failureMex)
        return ClientReturnValue.NewclientReturnValueRPC(protobuf.STATUS_FAILURE, failureMex)
    case protobuf.AdminOp_CHANGE_CONF_REM:
        if state.GetRole() == raftstate.LEADER{
            panic("not implemented")
            // return exitSucess
        }
        var failureMex = "i'm not leader, i cannot change conf"
        log.Println(failureMex)
        return ClientReturnValue.NewclientReturnValueRPC(protobuf.STATUS_FAILURE, failureMex)
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
