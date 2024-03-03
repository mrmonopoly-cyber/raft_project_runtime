package rpcs

import (
	"encoding/json"
	"log"
	l "raft/raft_log"
)

type AppendEntryRPC struct {
	Term         int       `json:"term"`
	LeaderId     int       `json:"leaderId"`
	PrevLogIndex int       `json:"prevLogIndex"`
	PrevLogTerm  int       `json:"prevLogTerm"`
	Entries      []l.Entry `json:"entries"`
	LeaderCommit int       `json:"leaderCommit"`
}

func NewAppendEntry(term int, leaderId int, prevLogIndex int, prevLogTerm int, entries []l.Entry, leaderCommit int) *AppendEntryRPC {
	var a = new(AppendEntryRPC)
	a.Term = term
	a.LeaderId = leaderId
	a.PrevLogIndex = prevLogIndex
	a.PrevLogTerm = prevLogTerm
	a.Entries = entries
	a.LeaderCommit = leaderCommit
	return a
}

func (a *AppendEntryRPC) ToJson() ([]byte, error) {
	/*data := map[string]interface{}{
		"term":         a.term,
		"leaderId":     a.leaderId,
		"prevLogIndex": a.prevLogIndex,
		"preLogTerm":   a.prevLogTerm,
		"entries":      a.entries,
		"leaderCommit": a.leaderCommit,
	}*/

	return json.Marshal(a)
}

func FromJson(data []byte) *AppendEntryRPC {

	var jsonData map[string]interface{}

	err := json.Unmarshal(data, &jsonData)

	if err != nil {
		log.Println("Error in decoding apppen")
	}

	return NewAppendEntry(jsonData["term"].(int),
		jsonData["leaderId"].(int),
		jsonData["prevLogIndex"].(int),
		jsonData["prevLogTerm"].(int),
		jsonData["entries"].([]l.Entry),
		jsonData["leaderCommit"].(int))

}
