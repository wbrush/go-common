package messaging

import (
	"encoding/json"
	"github.com/wbrush/go-common/datamodels"
)

type CommsRcvd struct {
	MessageId int64 `json:"message_id"`
	BaseMessageData
	Destination string `json:"destination"`
	MessageSid  string `json:"message_sid"`

	Status datamodels.MessageStatus `json:"status"`
}

func (m CommsRcvd) GetBody() []byte {
	b, _ := json.Marshal(m)
	return b
}
