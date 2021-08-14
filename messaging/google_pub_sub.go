package messaging

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/wbrush/go-common/config"
	"sync"
	"time"
)

const (
	PubSubProcessingTimeOutSecs  = 60
	PubSubRetentionDurationHours = 4
	PubsubExpirationPolicyDays   = 10
)

type GPubSub struct {
	client       *pubsub.Client
	subList      []*pubsub.Subscription
	subMut       sync.Mutex
	NumGoThreads int
}

func NewGPubSub(projectID string) (*GPubSub, error) {
	ps, err := pubsub.NewClient(context.Background(), projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %v", err)
	}

	return &GPubSub{
		client:       ps,
		NumGoThreads: -1,
	}, nil
}

//  gets/sets the Receive settings number of go threads. Use a parameter
//  with a negative value to get the current setting.
func (ps *GPubSub) SetNumGoThreads(setting int) (current int) {
	current = ps.NumGoThreads

	if setting >= 0 {
		ps.NumGoThreads = setting
	}

	return
}

func (ps *GPubSub) GetOrCreateTopic(topicName string) (topic *pubsub.Topic, err error) {
	topic = ps.client.Topic(topicName)
	if exists, err := topic.Exists(context.TODO()); err != nil || !exists {
		if err != nil {
			return nil, fmt.Errorf("can't check if topic exists: %s", err.Error())
		}

		if !exists {
			// Create the new topic.
			topic, err = ps.client.CreateTopic(context.TODO(), topicName)
			if err != nil {
				return nil, fmt.Errorf("failed to create topic: %v", err)
			}

		}
	}

	return topic, nil
}

func (ps *GPubSub) PublishToTopic(topic *pubsub.Topic, m Message) (serverID string, err error) {
	if topic == nil {
		return serverID, fmt.Errorf("topic cannot be nil")
	}
	if m == nil {
		return serverID, fmt.Errorf("message cannot be nil")
	}

	pr := topic.Publish(context.Background(), &pubsub.Message{
		Data: m.GetBody(),
	})

	serverID, err = pr.Get(context.TODO())
	if err != nil {
		return serverID, fmt.Errorf("failed to get publish result: %v", err)
	}

	return serverID, nil
}

func (ps *GPubSub) Publish(topicName string, m Message) (serverID string, err error) {
	topic, err := ps.GetOrCreateTopic(topicName)
	if err != nil {
		return serverID, err
	}
	defer topic.Stop()

	serverID, err = ps.PublishToTopic(topic, m)
	if err != nil {
		return serverID, err
	}

	return serverID, nil
}

func (ps GPubSub) MakeSubscriptionName(serviceName string, environment string, topicName string) string {
	return fmt.Sprintf(serviceName + "." + environment + "~" + topicName)
}

func (ps *GPubSub) SubscribeToTopic(topic *pubsub.Topic, subName string, callback SubscribeCallbackFunc) error {
	var err error

	sub := ps.client.Subscription(subName)

	subExists, err := sub.Exists(context.Background())
	if err != nil {
		return fmt.Errorf("cannot check subscription existance: %v", err)
	}

	if !subExists {
		sub, err = ps.client.CreateSubscription(context.Background(), subName,
			pubsub.SubscriptionConfig{Topic: topic,
				AckDeadline:       PubSubProcessingTimeOutSecs * time.Second,
				RetentionDuration: PubSubRetentionDurationHours * time.Hour,
				ExpirationPolicy:  PubsubExpirationPolicyDays * 24 * time.Hour})
		if err != nil {
			return fmt.Errorf("cannot subscribe to a topic: %v", err)
		}
	}

	if ps.NumGoThreads >= 0 {
		sub.ReceiveSettings.NumGoroutines = ps.NumGoThreads
	}

	ps.subMut.Lock()
	ps.subList = append(ps.subList, sub)
	ps.subMut.Unlock()

	err = sub.Receive(context.TODO(), func(ctx context.Context, m *pubsub.Message) {
		err = callback(m)
		if err != nil {
			m.Nack()
		}

		m.Ack()
	})
	if err != nil {
		return fmt.Errorf("cannot receive message: %v", err)
	}

	newList := make([]*pubsub.Subscription, 0)
	ps.subMut.Lock()
	for _, psSub := range ps.subList {
		if psSub != sub {
			newList = append(newList, psSub)
		}
	}
	ps.subList = newList
	ps.subMut.Unlock()

	return nil
}

func (ps *GPubSub) DeleteSubscription(topic *pubsub.Topic, subName string) error {
	var err error

	sub := ps.client.Subscription(subName)

	subExists, err := sub.Exists(context.Background())
	if err != nil {
		return fmt.Errorf("cannot check subscription existance: %v", err)
	}

	if !subExists {
		return nil
	}

	newList := make([]*pubsub.Subscription, 0)
	ps.subMut.Lock()
	for _, psSub := range ps.subList {
		if psSub != sub {
			newList = append(newList, psSub)
		}
	}
	ps.subList = newList
	ps.subMut.Unlock()

	err = sub.Delete(context.Background())
	if err != nil {
		return fmt.Errorf("cannot delete subscription: %s", err.Error())
	}

	return nil
}

func (ps *GPubSub) Subscribe(topicName string, subName string, callback SubscribeCallbackFunc) error {
	for {
		topic, err := ps.GetOrCreateTopic(topicName)
		if err != nil {
			logrus.Errorf("Subscribe Error: GetOrCreateTopic(%s) - %s", topicName, err.Error())
			time.Sleep(2 * time.Second)
			continue
			//return err
		}

		err = ps.SubscribeToTopic(topic, subName, callback)
		if err != nil {
			logrus.Errorf("Subscribe Error: SubscribeToTopic(%s) - %s", topicName, err.Error())
			//return err
		}

		topic.Stop()
		time.Sleep(2 * time.Second)
		logrus.Infof("Subscribe(): Retrying subscribing to topic (%s) after return", topicName)
	}

	return nil
}

func (ps *GPubSub) SubscribeToTopicsSubList(serviceName string, sp *config.ServiceParams, list []TopicSubscribe) {
	for _, topicSubItem := range list {
		topicName, _ := topicSubItem.Topic.GetTopicName(sp.Environment)
		subscriptName := ps.MakeSubscriptionName(serviceName, sp.Environment, topicName)

		go ps.Subscribe(topicName, subscriptName, topicSubItem.Handler)
	}
}

func (ps *GPubSub) Close() error {
	var err error

	ps.subMut.Lock()
	defer ps.subMut.Unlock()
	for _, sub := range ps.subList {
		if sub != nil {
			err = sub.Delete(context.TODO())
			if err != nil {
				return fmt.Errorf("can't delete subscription: %v", err)
			}
		}
	}

	if ps.client != nil {
		err = ps.client.Close()
		if err != nil {
			return fmt.Errorf("can't close client: %v", err)
		}
	}

	return nil
}
