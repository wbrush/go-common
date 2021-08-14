package messaging

import (
	"encoding/json"
	"strings"
)

type ScheduleProcessedMessageAction string

const (
	SchedulePublishedAction        ScheduleProcessedMessageAction = "published"
	ScheduleEntriesProcessedAction ScheduleProcessedMessageAction = "entries_processed"
	ScheduleRepublishedAction      ScheduleProcessedMessageAction = "republished"
	ScheduleCopiedAction           ScheduleProcessedMessageAction = "copied"
	ScheduleEmployeeRemovedAction  ScheduleProcessedMessageAction = "employee_removed"
	ScheduleProcessingError        ScheduleProcessedMessageAction = "error"
)

func (status ScheduleProcessedMessageAction) IsValid() bool {
	s := strings.ToLower(string(status))
	switch ScheduleProcessedMessageAction(s) {
	case SchedulePublishedAction,
		ScheduleEntriesProcessedAction,
		ScheduleRepublishedAction,
		ScheduleCopiedAction,
		ScheduleProcessingError,
		ScheduleEmployeeRemovedAction:
		return true
	default:
		return false
	}
}

type ScheduleProcessedMessage struct {
	BaseMessageData
	Action        ScheduleProcessedMessageAction `json:"action"` //  what actually happened to the record
	TargetId      int64                          `json:"targetId"`
	TargetSelf    string                         `json:"targetSelf"`
	SecondaryId   int64                          `json:"secondaryId"`
	SecondarySelf string                         `json:"secondarySelf"`
}

func (cm *ScheduleProcessedMessage) GetBody() []byte {
	b, _ := json.Marshal(cm)
	return b
}
