// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/wbrush/go-common/messaging (interfaces: PublisherSubscriber)

// Package mock_messaging is a generated GoMock package.
package mock_messaging

import (
	pubsub "cloud.google.com/go/pubsub"
	gomock "github.com/golang/mock/gomock"
	config "github.com/wbrush/go-common/config"
	messaging "github.com/wbrush/go-common/messaging"
	reflect "reflect"
)

// MockPublisherSubscriber is a mock of PublisherSubscriber interface
type MockPublisherSubscriber struct {
	ctrl     *gomock.Controller
	recorder *MockPublisherSubscriberMockRecorder
}

// MockPublisherSubscriberMockRecorder is the mock recorder for MockPublisherSubscriber
type MockPublisherSubscriberMockRecorder struct {
	mock *MockPublisherSubscriber
}

// NewMockPublisherSubscriber creates a new mock instance
func NewMockPublisherSubscriber(ctrl *gomock.Controller) *MockPublisherSubscriber {
	mock := &MockPublisherSubscriber{ctrl: ctrl}
	mock.recorder = &MockPublisherSubscriberMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPublisherSubscriber) EXPECT() *MockPublisherSubscriberMockRecorder {
	return m.recorder
}

// Close mocks base method
func (m *MockPublisherSubscriber) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockPublisherSubscriberMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockPublisherSubscriber)(nil).Close))
}

// GetOrCreateTopic mocks base method
func (m *MockPublisherSubscriber) GetOrCreateTopic(arg0 string) (*pubsub.Topic, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrCreateTopic", arg0)
	ret0, _ := ret[0].(*pubsub.Topic)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrCreateTopic indicates an expected call of GetOrCreateTopic
func (mr *MockPublisherSubscriberMockRecorder) GetOrCreateTopic(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrCreateTopic", reflect.TypeOf((*MockPublisherSubscriber)(nil).GetOrCreateTopic), arg0)
}

// MakeSubscriptionName mocks base method
func (m *MockPublisherSubscriber) MakeSubscriptionName(arg0, arg1, arg2 string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MakeSubscriptionName", arg0, arg1, arg2)
	ret0, _ := ret[0].(string)
	return ret0
}

// MakeSubscriptionName indicates an expected call of MakeSubscriptionName
func (mr *MockPublisherSubscriberMockRecorder) MakeSubscriptionName(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MakeSubscriptionName", reflect.TypeOf((*MockPublisherSubscriber)(nil).MakeSubscriptionName), arg0, arg1, arg2)
}

// Publish mocks base method
func (m *MockPublisherSubscriber) Publish(arg0 string, arg1 messaging.Message) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Publish", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Publish indicates an expected call of Publish
func (mr *MockPublisherSubscriberMockRecorder) Publish(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Publish", reflect.TypeOf((*MockPublisherSubscriber)(nil).Publish), arg0, arg1)
}

// PublishToTopic mocks base method
func (m *MockPublisherSubscriber) PublishToTopic(arg0 *pubsub.Topic, arg1 messaging.Message) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PublishToTopic", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PublishToTopic indicates an expected call of PublishToTopic
func (mr *MockPublisherSubscriberMockRecorder) PublishToTopic(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PublishToTopic", reflect.TypeOf((*MockPublisherSubscriber)(nil).PublishToTopic), arg0, arg1)
}

// Subscribe mocks base method
func (m *MockPublisherSubscriber) Subscribe(arg0, arg1 string, arg2 messaging.SubscribeCallbackFunc) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Subscribe", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Subscribe indicates an expected call of Subscribe
func (mr *MockPublisherSubscriberMockRecorder) Subscribe(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Subscribe", reflect.TypeOf((*MockPublisherSubscriber)(nil).Subscribe), arg0, arg1, arg2)
}

// SubscribeToTopic mocks base method
func (m *MockPublisherSubscriber) SubscribeToTopic(arg0 *pubsub.Topic, arg1 string, arg2 messaging.SubscribeCallbackFunc) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubscribeToTopic", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// SubscribeToTopic indicates an expected call of SubscribeToTopic
func (mr *MockPublisherSubscriberMockRecorder) SubscribeToTopic(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubscribeToTopic", reflect.TypeOf((*MockPublisherSubscriber)(nil).SubscribeToTopic), arg0, arg1, arg2)
}

// SubscribeToTopicsSubList mocks base method
func (m *MockPublisherSubscriber) SubscribeToTopicsSubList(arg0 string, arg1 *config.ServiceParams, arg2 []messaging.TopicSubscribe) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SubscribeToTopicsSubList", arg0, arg1, arg2)
}

// SubscribeToTopicsSubList indicates an expected call of SubscribeToTopicsSubList
func (mr *MockPublisherSubscriberMockRecorder) SubscribeToTopicsSubList(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubscribeToTopicsSubList", reflect.TypeOf((*MockPublisherSubscriber)(nil).SubscribeToTopicsSubList), arg0, arg1, arg2)
}
