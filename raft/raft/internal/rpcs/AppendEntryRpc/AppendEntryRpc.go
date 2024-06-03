package AppendEntryRpc

import (
	"fmt"
	"log"
	"raft/internal/node"
	"raft/internal/raftstate"
	"raft/internal/rpcs"
	"raft/internal/raft_log"
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

func GenerateHearthbeat(state raftstate.State) rpcs.Rpc {
    var prevLogIndex = state.LastLogIndex()
    var prevLogTerm uint64 = uint64(state.LastLogTerm())

    var app = &AppendEntryRpc{
        pMex: protobuf.AppendEntriesRequest{
            Term:         state.GetTerm(),
            LeaderIdPrivate:     state.GetMyIp(raftstate.PRI),
            LeaderIdPublic:     state.GetMyIp(raftstate.PUB),
            PrevLogIndex: int64(prevLogIndex),
            PrevLogTerm:  prevLogTerm,
            Entries:      nil,
            LeaderCommit: state.GetCommitIndex(),
        },
    }

    //	log.Print("hearthbit generated by ", state.GetId(), " : ", app.ToString())

    return app
}

func NewAppendEntryRPC(state raftstate.State, prevLogIndex int64, prevLogTerm uint64, 
    entries []*protobuf.LogEntry) rpcs.Rpc {
    return &AppendEntryRpc{
        pMex: protobuf.AppendEntriesRequest{
            Term:         state.GetTerm(),
            LeaderIdPrivate:     state.GetMyIp(raftstate.PRI),
            LeaderIdPublic:     state.GetMyIp(raftstate.PUB),
            PrevLogIndex: prevLogIndex,
            PrevLogTerm:  prevLogTerm,
            Entries:      entries,
            LeaderCommit: state.GetCommitIndex(),
        },
    }
}

func checkConsistency(prevLogIndex int64, prevLogTerm uint64, entries []raft_log.LogInstance) (ERRORS, int) {
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
    entryState = entries[prevLogIndex].Entry
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
func (this *AppendEntryRpc) Execute(state raftstate.State, sender node.Node) rpcs.Rpc {

    log.Println("executing appendEntryRpc: ", this.ToString())
    var role raftstate.Role = state.GetRole()
    var id string = state.GetMyIp(raftstate.PRI)
    var myTerm uint64 = state.GetTerm()
    var nextIdx int
    var consistent ERRORS
    var prevLogIndex int64 = this.pMex.GetPrevLogIndex()
    var prevLogTerm uint64 = this.pMex.GetPrevLogTerm()
    var entries []raft_log.LogInstance = state.GetEntries()
    var newEntries []*protobuf.LogEntry = this.pMex.GetEntries()
    var newEntriesWrapper []*raft_log.LogInstance = state.NewLogInstanceBatch(newEntries,[]func(){})
    var err error

    var resp rpcs.Rpc = nil
    var leaderCommit int64

    if this.pMex.GetTerm() < myTerm { //case 1
        return respondeAppend(id, false, myTerm, -1)
    }


    if role != raftstate.FOLLOWER {
        state.SetRole(raftstate.FOLLOWER)
    }

    state.SetLeaderIp(raftstate.PRI, this.pMex.LeaderIdPrivate)
    state.SetLeaderIp(raftstate.PUB, this.pMex.LeaderIdPublic)

    if  len(newEntries) > 0{
        //log.Println("received Append Entry", newEntries)
        consistent, nextIdx = checkConsistency(prevLogIndex, prevLogTerm, entries)
        switch consistent{
            case C2:
                resp = respondeAppend(id, false, myTerm, nextIdx)
            case C3:
                state.DeleteFromEntry(uint(nextIdx))
                resp = respondeAppend(id, false, myTerm, nextIdx)
            default:
                state.AppendEntries(newEntriesWrapper)
                leaderCommit = this.pMex.GetLeaderCommit()

                if leaderCommit > state.GetCommitIndex() {
                    state.MinimumCommitIndex(uint(leaderCommit))
                }
                resp = respondeAppend(id, true , myTerm, state.LastLogIndex())
        }
    } 

    if resp == nil {
        log.Println("hearthbeat")
        resp = respondeAppend(id, true, myTerm, state.LastLogIndex())
        log.Println("hearthbit resp: ", resp.ToString())
    }

    err = state.RestartTimeout(raftstate.TIMER_ELECTION)
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
