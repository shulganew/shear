// Code generated by MockGen. DO NOT EDIT.
// Source: internal/service/shortener.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/shulganew/shear.git/internal/model"
)

// MockStorageURL is a mock of StorageURL interface.
type MockStorageURL struct {
	ctrl     *gomock.Controller
	recorder *MockStorageURLMockRecorder
}

// MockStorageURLMockRecorder is the mock recorder for MockStorageURL.
type MockStorageURLMockRecorder struct {
	mock *MockStorageURL
}

// NewMockStorageURL creates a new mock instance.
func NewMockStorageURL(ctrl *gomock.Controller) *MockStorageURL {
	mock := &MockStorageURL{ctrl: ctrl}
	mock.recorder = &MockStorageURLMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorageURL) EXPECT() *MockStorageURLMockRecorder {
	return m.recorder
}

// DelelteBatch mocks base method.
func (m *MockStorageURL) DelelteBatch(ctx context.Context, userID string, briefs []string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "DelelteBatch", ctx, userID, briefs)
}

// DelelteBatch indicates an expected call of DelelteBatch.
func (mr *MockStorageURLMockRecorder) DelelteBatch(ctx, userID, briefs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DelelteBatch", reflect.TypeOf((*MockStorageURL)(nil).DelelteBatch), ctx, userID, briefs)
}

// GetAll mocks base method.
func (m *MockStorageURL) GetAll(ctx context.Context) []model.Short {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx)
	ret0, _ := ret[0].([]model.Short)
	return ret0
}

// GetAll indicates an expected call of GetAll.
func (mr *MockStorageURLMockRecorder) GetAll(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockStorageURL)(nil).GetAll), ctx)
}

// GetBrief mocks base method.
func (m *MockStorageURL) GetBrief(ctx context.Context, origin string) (string, bool, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBrief", ctx, origin)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(bool)
	return ret0, ret1, ret2
}

// GetBrief indicates an expected call of GetBrief.
func (mr *MockStorageURLMockRecorder) GetBrief(ctx, origin interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBrief", reflect.TypeOf((*MockStorageURL)(nil).GetBrief), ctx, origin)
}

// GetOrigin mocks base method.
func (m *MockStorageURL) GetOrigin(ctx context.Context, brief string) (string, bool, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrigin", ctx, brief)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(bool)
	return ret0, ret1, ret2
}

// GetOrigin indicates an expected call of GetOrigin.
func (mr *MockStorageURLMockRecorder) GetOrigin(ctx, brief interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrigin", reflect.TypeOf((*MockStorageURL)(nil).GetOrigin), ctx, brief)
}

// GetUserAll mocks base method.
func (m *MockStorageURL) GetUserAll(ctx context.Context, userID string) []model.Short {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserAll", ctx, userID)
	ret0, _ := ret[0].([]model.Short)
	return ret0
}

// GetUserAll indicates an expected call of GetUserAll.
func (mr *MockStorageURLMockRecorder) GetUserAll(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserAll", reflect.TypeOf((*MockStorageURL)(nil).GetUserAll), ctx, userID)
}

// Set mocks base method.
func (m *MockStorageURL) Set(ctx context.Context, userID, brief, origin string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", ctx, userID, brief, origin)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set.
func (mr *MockStorageURLMockRecorder) Set(ctx, userID, brief, origin interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockStorageURL)(nil).Set), ctx, userID, brief, origin)
}

// SetAll mocks base method.
func (m *MockStorageURL) SetAll(ctx context.Context, short []model.Short) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetAll", ctx, short)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetAll indicates an expected call of SetAll.
func (mr *MockStorageURLMockRecorder) SetAll(ctx, short interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetAll", reflect.TypeOf((*MockStorageURL)(nil).SetAll), ctx, short)
}