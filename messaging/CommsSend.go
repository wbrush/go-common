package messaging

import "encoding/json"

const (
	CommsProviderTwilio   = "twilio"
	CommsProviderSendGrid = "sendgrid"
)

const (
	CommsActionSend   = "send"
	CommsActionStatus = "status"
	CommsActionRcvd   = "rcvd"
)

type CommsSend struct {
	MessageId int64 `json:"message_id"`
	BaseMessageData
	Source          string `json:"source"`
	SourceName      string `json:"sourceName"`
	Destination     string `json:"destination"`
	DestinationName string `json:"destinationName"`
	Subject         string `json:"subject"`
	MessageBody     string `json:"body"`
}

func (m CommsSend) GetBody() []byte {
	b, _ := json.Marshal(m)
	return b
}
