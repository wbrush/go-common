package messaging

import (
	"encoding/json"
)

type UserStateChangedMessage struct {
	BaseMessageData
	UserId      int64  `json:"userId"`
	IsOnline    bool   `json:"isOnline"`
	Fingerprint string `json:"fingerprint"`
}

func (cm *UserStateChangedMessage) GetBody() []byte {
	b, _ := json.Marshal(cm)
	return b
}
