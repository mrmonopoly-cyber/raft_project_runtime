package AppendEntryRPC

import (
	"raft/internal/messages"
	appendEntryResponse "raft/internal/messages/AppendEntryResponse"
	//"raft/internal/node"
	"raft/internal/raftstate"
	p "raft/pkg/protobuf"
	"strconv"
	//"sync"

	"google.golang.org/protobuf/proto"
)

type AppendEntryRPC struct {
	term         uint64
	leaderId     string
	prevLogIndex uint64
	prevLogTerm  uint64
	entries      []*p.Entry
	leaderCommit uint64
}

func checkConsistency(prevLogIndex uint64, prevLogTerm uint64, state raftstate.State) bool {
  return state.GetEntries()[prevLogIndex].GetTerm() == prevLogTerm
}

// Manage implements messages.Rpc.
func (this *AppendEntryRPC) Execute(state *raftstate.State) *messages.Rpc{

    if ((*state).GetRole() != raftstate.FOLLOWER) {
        (*state).SetRole(raftstate.FOLLOWER)
    }

    var appendEntryResp messages.Rpc

    if (this.term < (*state).GetTerm()) || !checkConsistency(this.prevLogIndex, this.prevLogTerm, *state) { 
        appendEntryResp = appendEntryResponse.NewAppendEntryResponse("",false, 
        (*state).GetTerm(), 
        uint64(len((*state).GetEntries()) - 1))
    } else {

    }
    return &appendEntryResp
}

// ToString implements messages.Rpc.
func (this *AppendEntryRPC) ToString() string {
	var entries string
	for _, el := range this.entries {
		entries += el.String()
	}
	return "{term : " + strconv.Itoa(int(this.term)) + ", \nleaderId: " + this.leaderId + ",\nprevLogIndex: " + strconv.Itoa(int(this.prevLogIndex)) + ", \nprevLogTerm: " + strconv.Itoa(int(this.prevLogTerm)) + ", \nentries: " + entries + ", \nleaderCommit: " + strconv.Itoa(int(this.leaderCommit)) + "}"
}

func NewAppendEntryRPC(term uint64, leaderId string, prevLogIndex uint64,
	prevLogTerm uint64, entries []*p.Entry,
	leaderCommit uint64) messages.Rpc {
	return &AppendEntryRPC{
		term,
		leaderId,
		prevLogIndex,
		prevLogTerm,
		entries,
		leaderCommit,
	}
}

func GenerateHearthbeat(state raftstate.State) messages.Rpc {
	entries := state.GetEntries()
	prevLogIndex := len(entries) - 2
	prevLogTerm := entries[prevLogIndex].GetTerm()
  return &AppendEntryRPC{
		term:         state.GetTerm(),
		leaderId:     state.GetId(),
		prevLogIndex: uint64(prevLogIndex),
		prevLogTerm:  prevLogTerm,
		entries:      make([]*p.Entry, 0),
		leaderCommit: state.GetCommitIndex(),
	}
}


func (this AppendEntryRPC) GetId() string {
  return this.leaderId
}

func (this AppendEntryRPC) GetTerm() uint64 {
	return this.term
}
func (this AppendEntryRPC) Encode() ([]byte, error) {
	appendEntry := &p.AppendEntriesRequest{
		Term:         proto.Uint64(this.term),
		PrevLogIndex: proto.Uint64(this.prevLogIndex),
		PrevLogTerm:  proto.Uint64(this.prevLogTerm),
		CommitIndex:  proto.Uint64(this.leaderCommit),
		LeaderId:     proto.String(this.leaderId),
		Entries:      this.entries,
	}

	mess, err := proto.Marshal(appendEntry)
	return mess, err
}
func (this AppendEntryRPC) Decode(b []byte) error {
	pb := new(p.AppendEntriesRequest)
	err := proto.Unmarshal(b, pb)

	if err != nil {
		this.term = pb.GetTerm()
		this.leaderId = pb.GetLeaderId()
		this.leaderCommit = pb.GetCommitIndex()
		this.entries = pb.GetEntries()
		this.prevLogTerm = pb.GetPrevLogTerm()
		this.prevLogIndex = pb.GetPrevLogIndex()
	}

	return err
}
func (this AppendEntryRPC) GetEntries() []*p.Entry {
	return this.entries
}
func (this AppendEntryRPC) GetLeaderId() string {
	return this.leaderId
}
func (this AppendEntryRPC) GetPrevLogTerm() uint64 {
	return this.prevLogTerm
}
func (this AppendEntryRPC) GetPrevLogIndex() uint64 {
	return this.prevLogIndex
}
func (this AppendEntryRPC) GetLeaderCommit() uint64 {
	return this.leaderCommit
}
