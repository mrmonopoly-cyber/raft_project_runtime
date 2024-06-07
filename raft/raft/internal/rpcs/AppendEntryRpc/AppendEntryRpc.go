package AppendEntryRpc

import (
	"fmt"
	"log"
	"raft/internal/raft_log"
	clustermetadata "raft/internal/raftstate/clusterMetadata"
	nodestate "raft/internal/raftstate/confPool/NodeIndexPool/nodeState"
	"raft/internal/rpcs"
	app_resp "raft/internal/rpcs/AppendResponse"
	"raft/pkg/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"
	"strconv"

	"google.golang.org/protobuf/proto"
)

type ERRORS uint8
const(
    C0 ERRORS = iota
    C2 ERRORS = iota
    C3 ERRORS = iota
)

type AppendEntryRpc struct {
    pMex protobuf.AppendEntriesRequest
}

func GenerateHearthbeat(intLog raft_log.LogEntryRead, state clustermetadata.ClusterMetadata) rpcs.Rpc {
    var prevLogIndex = intLog.LastLogIndex()
    var prevLogTerm uint64 = uint64(intLog.LastLogTerm())

    var app = &AppendEntryRpc{
        pMex: protobuf.AppendEntriesRequest{
            Term:         state.GetTerm(),
            LeaderIdPrivate:     state.GetMyIp(clustermetadata.PRI),
            LeaderIdPublic:     state.GetMyIp(clustermetadata.PUB),
            PrevLogIndex: int64(prevLogIndex),
            PrevLogTerm:  prevLogTerm,
            Entries:      nil,
            LeaderCommit: intLog.GetCommitIndex(),
        },
    }

    //	log.Print("hearthbit generated by ", state.GetId(), " : ", app.ToString())

    return app
}

func NewAppendEntryRPC( state clustermetadata.ClusterMetadata, 
                        intLog raft_log.LogEntryRead,
                        prevLogIndex int64, 
                        prevLogTerm uint64, 
    entries []*protobuf.LogEntry) rpcs.Rpc {
    return &AppendEntryRpc{
        pMex: protobuf.AppendEntriesRequest{
            Term:         state.GetTerm(),
            LeaderIdPrivate:     state.GetMyIp(clustermetadata.PRI),
            LeaderIdPublic:     state.GetMyIp(clustermetadata.PUB),
            PrevLogIndex: prevLogIndex,
            PrevLogTerm:  prevLogTerm,
            Entries:      entries,
            LeaderCommit: intLog.GetCommitIndex(),
        },
    }
}

func checkConsistency(prevLogIndex int64, prevLogTerm uint64, entries []*protobuf.LogEntry) (ERRORS, int) {
    var logSize = len(entries)
    var entryState *protobuf.LogEntry =  nil

    if prevLogIndex < 0 {
        return C0, 0
    }

    if logSize == 0  && prevLogIndex > 0{
        log.Println("case 2: logSize = 0")
        return C2, 0
    }

    if logSize-1 < int(prevLogIndex) {
        log.Println("case 2")
        log.Printf("logSize - 1: %d, and prevLogIndex: %d", (logSize-1), int(prevLogIndex))
        return C2, (logSize - 1)
    }
    entryState = entries[prevLogIndex]
    fmt.Println("case 3")
    log.Println(entries)
    log.Printf("prevLogTerm: %d,, getTerm: %d, getDescr: %s,, getType: %o", prevLogTerm, entryState.GetTerm(), entryState.GetDescription(), entryState.GetOpType())
    
    consistent := entryState.GetTerm() == prevLogTerm

    if consistent {
        return C0, (int(prevLogIndex) + 1)
    } else {
        return C3, int(prevLogIndex)
    }
}


//Manage implements rpcs.Rpc.
func (this *AppendEntryRpc) Execute( 
            intLog raft_log.LogEntry,
            metadata clustermetadata.ClusterMetadata,
            senderState nodestate.NodeState)rpcs.Rpc {

    log.Println("executing appendEntryRpc: ", this.ToString())
    var role = metadata.GetRole()
    var id string = metadata.GetMyIp(clustermetadata.PRI)
    var myTerm uint64 = metadata.GetTerm()
    var nextIdx int
    var consistent ERRORS
    var prevLogIndex int64 = this.pMex.GetPrevLogIndex()
    var prevLogTerm uint64 = this.pMex.GetPrevLogTerm()
    var entries []*protobuf.LogEntry = intLog.GetEntries()
    var newEntries []*protobuf.LogEntry = this.pMex.GetEntries()
    var newEntriesWrapper []*raft_log.LogInstance = intLog.NewLogInstanceBatch(newEntries,nil)
    var err error

    var resp rpcs.Rpc = nil
    var leaderCommit int64

    if this.pMex.GetTerm() < myTerm { //case 1
        return respondeAppend(id, false, myTerm, -1)
    }


    if role != clustermetadata.FOLLOWER {
        metadata.SetRole(clustermetadata.FOLLOWER)
    }

    metadata.SetLeaderIp(clustermetadata.PRI, this.pMex.LeaderIdPrivate)
    metadata.SetLeaderIp(clustermetadata.PUB, this.pMex.LeaderIdPublic)

    if  len(newEntries) > 0{
        //log.Println("received Append Entry", newEntries)
        consistent, nextIdx = checkConsistency(prevLogIndex, prevLogTerm, entries)
        switch consistent{
            case C2:
                resp = respondeAppend(id, false, myTerm, nextIdx)
            case C3:
                intLog.DeleteFromEntry(uint(nextIdx))
                resp = respondeAppend(id, false, myTerm, nextIdx)
            default:
                intLog.AppendEntry(newEntriesWrapper, int(prevLogIndex))
                //FIX: check if the entry are already in the log
                leaderCommit = this.pMex.GetLeaderCommit()

                if leaderCommit > intLog.GetCommitIndex() {
                    intLog.MinimumCommitIndex(uint(leaderCommit))
                }
                resp = respondeAppend(id, true , myTerm, intLog.LastLogIndex())
        }
    } 

    if resp == nil {
        log.Println("hearthbeat")
        resp = respondeAppend(id, true, myTerm, intLog.LastLogIndex())
        log.Println("hearthbit resp: ", resp.ToString())
    }

    err = metadata.RestartTimeout(clustermetadata.TIMER_ELECTION)
    if err != nil {
        log.Panicln("failed restarting election timer")
    }
    return resp
}

func respondeAppend(id string, success bool, term uint64, error int) rpcs.Rpc {
    var appendEntryResp rpcs.Rpc = app_resp.NewAppendResponseRPC(
        id,
        success,
        term,
        error)
    return appendEntryResp
}

// ToString implements rpcs.Rpc.
func (this *AppendEntryRpc) ToString() string {
    var entries string
    for _, el := range this.pMex.Entries {
        entries += el.String()
    }
    return "{term : " + strconv.Itoa(int(this.pMex.GetTerm())) + 
    ", leaderIdPrivate: " + this.pMex.GetLeaderIdPrivate() +
    ", leaderIdPublic: " + this.pMex.GetLeaderIdPublic() +
    ", prevLogIndex: " + strconv.Itoa(int(this.pMex.PrevLogIndex)) +
    ", prevLogTerm: " + strconv.Itoa(int(this.pMex.PrevLogTerm)) + 
    ", entries: " + entries +
    ", leaderCommit: " + strconv.Itoa(int(this.pMex.LeaderCommit)) + "}"
}

func (this *AppendEntryRpc) Encode() ([]byte, error) {

    var mess []byte
    var err error
    mess, err = proto.Marshal(&(*this).pMex)
    if err != nil {
        log.Panicln("error in Encoding Append Entry: ", err)
    }
    return mess, err
}

func (this *AppendEntryRpc) Decode(rawMex []byte) error {
    err := proto.Unmarshal(rawMex, &this.pMex)
    if err != nil {
        log.Panicln("error in Decoding Append Entry: ", err)
    }
    return err
}
