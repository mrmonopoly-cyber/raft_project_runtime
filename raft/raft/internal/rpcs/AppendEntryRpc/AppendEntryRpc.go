package AppendEntryRpc

import (
    "fmt"
    "log"
    "raft/internal/node/nodeState"
    "raft/internal/raftstate"
    "raft/internal/rpcs"
    app_resp "raft/internal/rpcs/AppendResponse"
	"raft/pkg/raft-rpcProtobuf-messages/rpcEncoding/out/protobuf"
    "strconv"

    "google.golang.org/protobuf/proto"
)

type AppendEntryRpc struct {
    pMex protobuf.AppendEntriesRequest
}

func GenerateHearthbeat(state raftstate.State) rpcs.Rpc {
    var entries []*protobuf.LogEntry = state.GetEntries()
    prevLogIndex := len(entries) - 1
    var prevLogTerm uint64 = 0
    if len(entries) > 0 {
        prevLogTerm = entries[prevLogIndex].GetTerm()
    }

    var app = &AppendEntryRpc{
        pMex: protobuf.AppendEntriesRequest{
            Term:         state.GetTerm(),
            LeaderIdPrivate:     state.GetIdPrivate(),
            LeaderIdPublic:     state.GetIdPublic(),
            PrevLogIndex: int64(prevLogIndex),
            PrevLogTerm:  prevLogTerm,
            Entries:      make([]*protobuf.LogEntry, 0),
            LeaderCommit: state.GetCommitIndex(),
        },
    }

    //	log.Print("hearthbit generated by ", state.GetId(), " : ", app.ToString())

    return app
}

func NewAppendEntryRPC(term uint64, leaderIdPrivate string, leaderIdPublic string, prevLogIndex int64,
prevLogTerm uint64, entries []*protobuf.LogEntry,
leaderCommit int64) rpcs.Rpc {
    return &AppendEntryRpc{
        pMex: protobuf.AppendEntriesRequest{
            Term:         term,
            LeaderIdPrivate: leaderIdPrivate,
            LeaderIdPublic: leaderIdPublic,
            PrevLogIndex: prevLogIndex,
            PrevLogTerm:  prevLogTerm,
            Entries:      entries,
            LeaderCommit: leaderCommit,
        },
    }
}

func checkConsistency(prevLogIndex int64, prevLogTerm uint64, entries []*protobuf.LogEntry) (bool, int) {
    var logSize = len(entries)
    if prevLogIndex < 0 {
        return true, 0
    }

    if logSize == 0  && prevLogIndex > 0{
        log.Println("case 2: logSize = 0")
        return false, 0
    }

    if logSize-1 < int(prevLogIndex) {
        log.Println("case 2")
        log.Printf("logSize - 1: %d, and prevLogIndex: %d", (logSize-1), int(prevLogIndex))
        return false, (logSize - 1)
    }
    fmt.Print("case 3")
  log.Println(entries)
  log.Printf("prevLogTerm: %d,, getTerm: %d, getDescr: %s,, getType: %o", prevLogTerm, entries[prevLogIndex].GetTerm(), entries[prevLogIndex].GetDescription(), entries[prevLogIndex].GetOpType())
    
    consistent := entries[prevLogIndex].GetTerm() == prevLogTerm

    if consistent {
        return true, (int(prevLogIndex) + 1)
    } else {
        return false, int(prevLogIndex)
    }
}

// Manage implements rpcs.Rpc.
func (this *AppendEntryRpc) Execute(state *raftstate.State, senderState *nodeState.VolatileNodeState) *rpcs.Rpc {

    var role raftstate.Role = (*state).GetRole()
    var id string = (*state).GetIdPrivate()
    var myTerm uint64 = (*state).GetTerm()
    var nextIdx int
    var consistent bool
    var prevLogIndex int64 = this.pMex.GetPrevLogIndex()
    var prevLogTerm uint64 = this.pMex.GetPrevLogTerm()
    var entries []*protobuf.LogEntry = (*state).GetEntries()
    var newEntries []*protobuf.LogEntry = this.pMex.GetEntries()

    var resp *rpcs.Rpc = nil
    var leaderCommit int64
    var lastNewEntryIdx int64

    if this.pMex.GetTerm() < myTerm {
        return respondeAppend(id, false, myTerm, -1)
    }

    (*state).StopElectionTimeout()

    if role != raftstate.FOLLOWER {
        (*state).SetRole(raftstate.FOLLOWER)
    }

    (*state).SetLeaderIpPrivate(this.pMex.LeaderIdPrivate)
    (*state).SetLeaderIpPublic(this.pMex.LeaderIdPublic)

    if len(newEntries) > 0 {

        //log.Println("received Append Entry", newEntries)
        consistent, nextIdx = checkConsistency(prevLogIndex, prevLogTerm, entries)
        fmt.Println(!consistent)

        if !consistent {
            fmt.Println("Not consistent")
            resp = respondeAppend(id, false, myTerm, nextIdx)
        } else {
            (*state).AppendEntries(newEntries)
            leaderCommit = this.pMex.GetLeaderCommit()
            lastNewEntryIdx = int64(len(entries) - 1)

            if leaderCommit > (*state).GetCommitIndex() {
                if leaderCommit > lastNewEntryIdx {
                    (*state).SetCommitIndex(lastNewEntryIdx)
                } else {
                    (*state).SetCommitIndex(leaderCommit)
                }
            }
            resp = respondeAppend(id, true , myTerm, (*state).LastLogIndex())
        }
    } 

    if resp == nil {
        log.Println("hearthbeat")
        resp = respondeAppend(id, true, myTerm, (*state).LastLogIndex())
    }

    (*state).StartElectionTimeout()
    return resp
}

func respondeAppend(id string, success bool, term uint64, error int) *rpcs.Rpc {
    var appendEntryResp rpcs.Rpc = app_resp.NewAppendResponseRPC(
        id,
        success,
        term,
        error)
        return &appendEntryResp
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
        ", prevLogTerm: " + strconv.Itoa(int(this.pMex.PrevLogIndex)) + 
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
