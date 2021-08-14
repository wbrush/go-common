package messaging

import (
	"errors"
	"strings"

	"cloud.google.com/go/pubsub"
	"github.com/wbrush/go-common/config"
)

type SubscribeCallbackFunc func(m *pubsub.Message) error

type MessageType string

const (
	MessageTypeBaseMessage       MessageType = "BaseMessage"
	MessageTypeChangeMessage     MessageType = "ChangeMessage"
	MessageTypeCommsSend         MessageType = "CommsSend"
	MessageTypeCommsRcvd         MessageType = "CommsRcvd"
	MessageTypeScheduleProcessed MessageType = "ScheduleProcessedMessage"
	MessageTypeStateChange       MessageType = "StateChange"
)

type (
	Message interface {
		GetBody() []byte
	}

	TopicSubscribe struct {
		Topic   PubSubTopic
		Handler SubscribeCallbackFunc
	}

	PublisherSubscriber interface {
		GetOrCreateTopic(topicName string) (topic *pubsub.Topic, err error)
		PublishToTopic(topic *pubsub.Topic, m Message) (serverID string, err error)
		Publish(topicName string, m Message) (serverID string, err error)
		MakeSubscriptionName(serviceName string, environment string, topicName string) string
		SubscribeToTopic(topic *pubsub.Topic, subName string, callback SubscribeCallbackFunc) error
		Subscribe(topicName string, subName string, callback SubscribeCallbackFunc) error
		SubscribeToTopicsSubList(serviceName string, sp *config.ServiceParams, list []TopicSubscribe)
		Close() error
	}

	BaseMessageData struct {
		Filter FilterType `json:"filter"` /* {provider}.{media}.{action}.{property}.{id} - should be included in message
		but preferably this would be added to the message so it can be filtered by the receiver */

		Type    MessageType `json:"type"`
		SelfUrl string      `json:"self"` // url to the record that changed so I can get the new data if I want to
	}
)

type FilterType string

const OldMessageNumFilters = int(5)
const MessageNumFilters = int(6)

func (filter *FilterType) New(provider, media, action, property, id string) (err error) {
	filter.NewV2(provider, media, action, property, id, "tbi")
	return
}

func (filter *FilterType) NewV2(provider, media, action, property, id, sender string) (err error) {
	*filter = FilterType(provider + "." + media + "." + action + "." + property + "." + id + "." + sender)
	return
}

func (filter FilterType) GetProvider() (provider string, err error) {
	list := strings.Split(string(filter), ".")
	if len(list) != OldMessageNumFilters && len(list) != MessageNumFilters {
		err = errors.New("invalid Filter Format")
	} else {
		provider = list[0]
	}
	return
}

func (filter FilterType) GetMedia() (media string, err error) {
	list := strings.Split(string(filter), ".")
	if len(list) != OldMessageNumFilters && len(list) != MessageNumFilters {
		err = errors.New("invalid Filter Format")
	} else {
		media = list[1]
	}
	return
}

func (filter FilterType) GetAction() (action string, err error) {
	list := strings.Split(string(filter), ".")
	if len(list) != OldMessageNumFilters && len(list) != MessageNumFilters {
		err = errors.New("invalid Filter Format")
	} else {
		action = list[2]
	}
	return
}

func (filter FilterType) GetProperty() (id string, err error) {
	list := strings.Split(string(filter), ".")
	if len(list) != OldMessageNumFilters && len(list) != MessageNumFilters {
		err = errors.New("invalid Filter Format")
	} else {
		id = list[3]
	}
	return
}

func (filter FilterType) GetId() (id string, err error) {
	list := strings.Split(string(filter), ".")
	if len(list) != OldMessageNumFilters && len(list) != MessageNumFilters {
		err = errors.New("invalid Filter Format")
	} else {
		id = list[4]
	}
	return
}

func (filter FilterType) GetSender() (sender string, err error) {
	list := strings.Split(string(filter), ".")
	if len(list) != OldMessageNumFilters && len(list) != MessageNumFilters {
		err = errors.New("invalid Filter Format")
	} else {
		sender = list[5]
	}
	return
}

func (filter FilterType) String() string {
	return string(filter)
}
