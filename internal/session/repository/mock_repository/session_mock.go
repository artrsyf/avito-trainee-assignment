// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	reflect "reflect"

	entity "github.com/artrsyf/avito-trainee-assignment/internal/session/domain/entity"
	model "github.com/artrsyf/avito-trainee-assignment/internal/session/domain/model"
	gomock "github.com/golang/mock/gomock"
)

// MockSessionRepositoryI is a mock of SessionRepositoryI interface.
type MockSessionRepositoryI struct {
	ctrl     *gomock.Controller
	recorder *MockSessionRepositoryIMockRecorder
}

// MockSessionRepositoryIMockRecorder is the mock recorder for MockSessionRepositoryI.
type MockSessionRepositoryIMockRecorder struct {
	mock *MockSessionRepositoryI
}

// NewMockSessionRepositoryI creates a new mock instance.
func NewMockSessionRepositoryI(ctrl *gomock.Controller) *MockSessionRepositoryI {
	mock := &MockSessionRepositoryI{ctrl: ctrl}
	mock.recorder = &MockSessionRepositoryIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSessionRepositoryI) EXPECT() *MockSessionRepositoryIMockRecorder {
	return m.recorder
}

// Check mocks base method.
func (m *MockSessionRepositoryI) Check(ctx context.Context, userID uint) (*model.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Check", ctx, userID)
	ret0, _ := ret[0].(*model.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Check indicates an expected call of Check.
func (mr *MockSessionRepositoryIMockRecorder) Check(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Check", reflect.TypeOf((*MockSessionRepositoryI)(nil).Check), ctx, userID)
}

// Create mocks base method.
func (m *MockSessionRepositoryI) Create(ctx context.Context, sessionEntity *entity.Session) (*model.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, sessionEntity)
	ret0, _ := ret[0].(*model.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockSessionRepositoryIMockRecorder) Create(ctx, sessionEntity interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockSessionRepositoryI)(nil).Create), ctx, sessionEntity)
}
