package messaging

import (
	"encoding/json"
	"strings"
)

type ChangeMessageAction string

const (
	RecordUpdateNew      ChangeMessageAction = "new"
	RecordUpdateModified ChangeMessageAction = "modified"
	RecordUpdateSync     ChangeMessageAction = "sync" //  used to initialize new caches and transmit the current state
	RecordUpdateArchived ChangeMessageAction = "archived"
	RecordUpdateDeleted  ChangeMessageAction = "deleted"
)

func (status ChangeMessageAction) IsValid() bool {
	s := strings.ToLower(string(status))
	switch ChangeMessageAction(s) {
	case RecordUpdateNew,
		RecordUpdateSync,
		RecordUpdateModified,
		RecordUpdateArchived,
		RecordUpdateDeleted:
		return true
	default:
		return false
	}
}

// as this is for MQs, I would not think it would need a swagger notation
type ChangeMessage struct {
	BaseMessageData
	Action ChangeMessageAction `json:"action"` //  what actually happened to the record
}

func (cm *ChangeMessage) GetBody() []byte {
	b, _ := json.Marshal(cm)
	return b
}
