// Code generated by MockGen. DO NOT EDIT.
// Source: ./Service/sms/types.go
//
// Generated by this command:
//
//      mockgen --source=./Service/sms/types.go --package= smsmock --destination= ./Service/mocks/sms/types.mock.go
//

// Package mock_sms is a generated GoMock package.
package mock_sms

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// Send mocks base method.
func (m *MockService) Send(ctx context.Context, tplID string, args []string, number ...string) error {
	m.ctrl.T.Helper()
	varargs := []any{ctx, tplID, args}
	for _, a := range number {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Send", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send.
func (mr *MockServiceMockRecorder) Send(ctx, tplID, args any, number ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, tplID, args}, number...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockService)(nil).Send), varargs...)
}