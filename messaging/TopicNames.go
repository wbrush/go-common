package messaging

import (
	"errors"
	"github.com/wbrush/go-common/config"
	"strings"
)

type PubSubTopic string

const (
	PubSubCmdTopic        PubSubTopic = "pubsub-cmd-topic"
	DepartmentTopic       PubSubTopic = "department-topic"
	PropertyTopic         PubSubTopic = "property-topic"
	EmployeeTopic         PubSubTopic = "employee-topic"
	UserInfoTopic         PubSubTopic = "user-info-topic"
	RolesTopic            PubSubTopic = "roles-topic"
	LocationsTopic        PubSubTopic = "locations-topic"
	EmploymentTypesTopic  PubSubTopic = "employment-types-topic"
	InvitationTopic       PubSubTopic = "invitation-topic"
	ScheduleTopic         PubSubTopic = "schedule-topic"
	ScheduleEntryTopic    PubSubTopic = "schedule-entry-topic"
	AdhocJobsTopic        PubSubTopic = "adhoc-jobs-topic"
	LongTermStorageTopic  PubSubTopic = "longterm-storage-topic"
	V2XLateTopic          PubSubTopic = "v2-xlate-topic"
	NotificationTopic     PubSubTopic = "notification-send-topic"
	JobItemsTopic         PubSubTopic = "job-items-topic"
	JobSettingsTopic      PubSubTopic = "job-settings-topic"
	IdaasProxyTopic       PubSubTopic = "idaas-proxy-topic"
	EscalationNotifyTopic PubSubTopic = "escalation-notification-topic"
	TemplateServiceTopic  PubSubTopic = "template-svc-topic"
	TimekeeperTopic       PubSubTopic = "timekeeper-topic"

	CommsSendTopic PubSubTopic = "comms-send-topic"
	CommsRcvdTopic PubSubTopic = "comms-rcvd-topic"
)

const (
	topicSeparator string = "." // the first letter is lowercase because I want this local only!
)

//  allow for validating both topic and topic name
func (str PubSubTopic) Validate() (rv bool) {
	rv = true
	switch PubSubTopic(str) {
	case PubSubCmdTopic, DepartmentTopic:
		break
	default:
		rv = false
	}

	return
}

func (str PubSubTopic) GetTopicName(env string) (string, error) {
	if env == "" {
		return "", errors.New("environment Suffix can not be null or empty")
	}
	if !config.IsEnviromentValid(env) {
		return "", errors.New("invalid Environment value")
	}
	return string(str) + topicSeparator + string(env), nil
}

func (str *PubSubTopic) GetTopicFromName(name string) error {
	if name == "" {
		return errors.New("topic name can not be null or empty")
	}
	list := strings.Split(name, topicSeparator)
	if len(list) != 2 {
		return errors.New("invalid format for Topic Name")
	}
	if !PubSubTopic(list[0]).Validate() {
		return errors.New("invalid topic")
	}

	*str = PubSubTopic(list[0])
	return nil
}
