// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/KWRI/go-common/httphelper (interfaces: HttpHelperInterface)

// Package mock_httphelper is a generated GoMock package.
package test

import (
	gomock "github.com/golang/mock/gomock"
	http "net/http"
	reflect "reflect"
)

// MockHttpHelperInterface is a mock of HttpHelperInterface interface
type MockHttpHelperInterface struct {
	ctrl     *gomock.Controller
	recorder *MockHttpHelperInterfaceMockRecorder
}

// MockHttpHelperInterfaceMockRecorder is the mock recorder for MockHttpHelperInterface
type MockHttpHelperInterfaceMockRecorder struct {
	mock *MockHttpHelperInterface
}

// NewMockHttpHelperInterface creates a new mock instance
func NewMockHttpHelperInterface(ctrl *gomock.Controller) *MockHttpHelperInterface {
	mock := &MockHttpHelperInterface{ctrl: ctrl}
	mock.recorder = &MockHttpHelperInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockHttpHelperInterface) EXPECT() *MockHttpHelperInterfaceMockRecorder {
	return m.recorder
}

// Do mocks base method
func (m *MockHttpHelperInterface) Do(arg0 *http.Request) (*http.Response, error) {
	ret := m.ctrl.Call(m, "Do", arg0)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Do indicates an expected call of Do
func (mr *MockHttpHelperInterfaceMockRecorder) Do(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Do", reflect.TypeOf((*MockHttpHelperInterface)(nil).Do), arg0)
}
