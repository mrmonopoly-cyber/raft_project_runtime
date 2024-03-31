package AppendEntryRpc

import (
	"log"
	"raft/internal/raftstate"
	"raft/internal/rpcs"
    "raft/internal/node/nodeState"
	app_resp "raft/internal/rpcs/AppendResponse"
	"raft/pkg/rpcEncoding/out/protobuf"
	"strconv"

	"google.golang.org/protobuf/proto"
)

type AppendEntryRpc struct {
	pMex protobuf.AppendEntriesRequest
}

func GenerateHearthbeat(state raftstate.State) rpcs.Rpc {
	var entries []protobuf.LogEntry = state.GetEntries()
	prevLogIndex := len(entries)
	var prevLogTerm uint64 = 0
	if prevLogIndex > 0 {
		prevLogIndex -= 2
		prevLogTerm = entries[prevLogIndex].GetTerm()
	}

	var app = &AppendEntryRpc{
		pMex: protobuf.AppendEntriesRequest{
			Term:         state.GetTerm(),
			LeaderId:     state.GetId(),
			PrevLogIndex: uint64(prevLogIndex),
			PrevLogTerm:  prevLogTerm,
			Entries:      make([]*protobuf.LogEntry, 0),
			LeaderCommit: state.GetCommitIndex(),
		},
	}

	log.Print("hearthbit generated by ", state.GetId(), " : ", app.ToString())

	return app
}

func NewAppendEntryRPC(term uint64, leaderId string, prevLogIndex uint64,
	prevLogTerm uint64, entries []*protobuf.LogEntry,
	leaderCommit int64) rpcs.Rpc {
	return &AppendEntryRpc{
		pMex: protobuf.AppendEntriesRequest{
			Term:         term,
			LeaderId:     leaderId,
			PrevLogIndex: prevLogIndex,
			PrevLogTerm:  prevLogTerm,
			Entries:      entries,
			LeaderCommit: leaderCommit,
		},
	}
}

func checkConsistency(prevLogIndex uint64, prevLogTerm uint64, entries []protobuf.LogEntry) bool {
    if len(entries) <= 0 {
        return false
    }
	return entries[prevLogIndex].GetTerm() == prevLogTerm
}

// Manage implements rpcs.Rpc.
func (this *AppendEntryRpc) Execute(state *raftstate.State, senderState *nodeState.VolatileNodeState) *rpcs.Rpc {
	(*state).StopElectionTimeout()
	defer (*state).StartElectionTimeout()

	var role raftstate.Role = (*state).GetRole()
	var id string = (*state).GetId()
	var myTerm uint64 = (*state).GetTerm()
	var error uint64
	var success bool
	var prevLogIndex uint64 = this.pMex.GetPrevLogIndex()
	var prevLogTerm uint64 = this.pMex.GetPrevLogTerm()
	var entries []protobuf.LogEntry = (*state).GetEntries()

	if role != raftstate.FOLLOWER {
		(*state).BecomeFollower()
	}

	if this.pMex.GetTerm() < myTerm {

		success = false
		return respondeAppend(id, success, myTerm, -1)

	} else if checkConsistency(prevLogIndex, prevLogTerm, entries) {

		success = false
		error = prevLogIndex
		return respondeAppend(id, success, myTerm, int(error))

	} else {

		(*state).AppendEntries(this.pMex.GetEntries(), int(prevLogIndex))
		success = true
		var leaderCommit int64 = this.pMex.GetLeaderCommit()
		var lastNewEntryIdx int64 = int64(len(entries) - 1)
		if leaderCommit > (*state).GetCommitIndex() {
			if leaderCommit > lastNewEntryIdx {
				(*state).SetCommitIndex(lastNewEntryIdx)
			} else {
				(*state).SetCommitIndex(leaderCommit)
			}
		}
		return respondeAppend(id, success, myTerm, -1)

	}
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
	return "{term : " + strconv.Itoa(int(this.pMex.GetTerm())) + ", leaderId: " + this.pMex.GetLeaderId() + ", prevLogIndex: " + strconv.Itoa(int(this.pMex.PrevLogIndex)) + ", prevLogTerm: " + strconv.Itoa(int(this.pMex.PrevLogIndex)) + ", entries: " + entries + ", leaderCommit: " + strconv.Itoa(int(this.pMex.LeaderCommit)) + "}"
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
func (this *AppendEntryRpc) Decode(rawMex []byte) (error) {
	err := proto.Unmarshal(rawMex, &this.pMex)
    if err != nil {
        log.Panicln("error in Decoding Append Entry: ", err)
    }
	return err
}
