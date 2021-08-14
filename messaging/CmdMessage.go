package messaging

import (
	"encoding/json"
	"strings"
)

type CommandMessageAction string

const (
	CommandMessageActionSync CommandMessageAction = "sync"
)

func (status CommandMessageAction) IsValid() bool {
	s := strings.ToLower(string(status))
	switch CommandMessageAction(s) {
	case CommandMessageActionSync:
		return true
	default:
		return false
	}
}

type CommandMessage struct {
	BaseMessageData
	Action CommandMessageAction `json:"action"` //  what actually happened to the record
}

func (cm *CommandMessage) GetBody() []byte {
	b, _ := json.Marshal(cm)
	return b
}
