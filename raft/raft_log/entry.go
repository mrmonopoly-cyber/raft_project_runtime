package raft_log

import (
	"encoding/json"
	"log"
)

type Operation int

const (
	READ Operation = iota
	WRITE
	DELETE
	JOINT_CONF
	NODE_UPDATED
	NODE_DELETED
)

type Entry struct {
	Term        int       `json:"term"`
	Operation   Operation `json:"operation"`
	Description string    `json:"description"`
}

func NewEntry(term int, operation Operation, description string) *Entry {
	var e = new(Entry)
	e.Description = description
	e.Term = term
	e.Operation = operation
	return e
}

func (self Entry) GetTerm() int {
	return self.Term
}

func (self Entry) GetDescription() string {
	return self.Description
}

func (self Entry) GetOpType() Operation {
	return self.Operation
}

func (e *Entry) ToJson() ([]byte, error) {
	return json.Marshal(e)
}

func FromJson(data []byte) *Entry {
	var jsonData map[string]interface{}

	err := json.Unmarshal(data, &jsonData)

	if err != nil {
		log.Println("Error in decoding")
	}

	return NewEntry(jsonData["term"].(int),
		jsonData["operation"].(Operation),
		jsonData["description"].(string))
}
