package messaging

import (
	"encoding/json"
)

type CommsRcvd struct {
	MessageId int64 `json:"message_id"`
	BaseMessageData
	Destination string `json:"destination"`
	MessageSid  string `json:"message_sid"`

	Status string `json:"status"`
}

func (m CommsRcvd) GetBody() []byte {
	b, _ := json.Marshal(m)
	return b
}
