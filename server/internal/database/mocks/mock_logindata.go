// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/GophKeeper/server/internal/app/logindata (interfaces: ILoginDataStore)

// Package mock_logindata is a generated GoMock package.
package mocks

import (
        context "context"
        reflect "reflect"

        logindata "github.com/GophKeeper/server/internal/app/logindata"
        gomock "github.com/golang/mock/gomock"
)

// MockILoginDataStore is a mock of ILoginDataStore interface.
type MockILoginDataStore struct {
        ctrl     *gomock.Controller
        recorder *MockILoginDataStoreMockRecorder
}

// MockILoginDataStoreMockRecorder is the mock recorder for MockILoginDataStore.
type MockILoginDataStoreMockRecorder struct {
        mock *MockILoginDataStore
}

// NewMockILoginDataStore creates a new mock instance.
func NewMockILoginDataStore(ctrl *gomock.Controller) *MockILoginDataStore {
        mock := &MockILoginDataStore{ctrl: ctrl}
        mock.recorder = &MockILoginDataStoreMockRecorder{mock}
        return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockILoginDataStore) EXPECT() *MockILoginDataStoreMockRecorder {
        return m.recorder
}

// Create mocks base method.
func (m *MockILoginDataStore) Create(arg0 context.Context, arg1 logindata.LoginData) error {
        m.ctrl.T.Helper()
        ret := m.ctrl.Call(m, "Create", arg0, arg1)
        ret0, _ := ret[0].(error)
        return ret0
}

// Create indicates an expected call of Create.
func (mr *MockILoginDataStoreMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
        mr.mock.ctrl.T.Helper()
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockILoginDataStore)(nil).Create), arg0, arg1)
}

// Delete mocks base method.
func (m *MockILoginDataStore) Delete(arg0 context.Context, arg1 string) error {
        m.ctrl.T.Helper()
        ret := m.ctrl.Call(m, "Delete", arg0, arg1)
        ret0, _ := ret[0].(error)
        return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockILoginDataStoreMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
        mr.mock.ctrl.T.Helper()
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockILoginDataStore)(nil).Delete), arg0, arg1)
}

// GetList mocks base method.
func (m *MockILoginDataStore) GetList(arg0 context.Context, arg1 string) ([]logindata.LoginData, error) {
        m.ctrl.T.Helper()
        ret := m.ctrl.Call(m, "GetList", arg0, arg1)
        ret0, _ := ret[0].([]logindata.LoginData)
        ret1, _ := ret[1].(error)
        return ret0, ret1
}

// GetList indicates an expected call of GetList.
func (mr *MockILoginDataStoreMockRecorder) GetList(arg0, arg1 interface{}) *gomock.Call {
        mr.mock.ctrl.T.Helper()
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetList", reflect.TypeOf((*MockILoginDataStore)(nil).GetList), arg0, arg1)
}

// Update mocks base method.
func (m *MockILoginDataStore) Update(arg0 context.Context, arg1 string, arg2 []byte) error {
        m.ctrl.T.Helper()
        ret := m.ctrl.Call(m, "Update", arg0, arg1, arg2)
        ret0, _ := ret[0].(error)
        return ret0
}

// Update indicates an expected call of Update.
func (mr *MockILoginDataStoreMockRecorder) Update(arg0, arg1, arg2 interface{}) *gomock.Call {
        mr.mock.ctrl.T.Helper()
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockILoginDataStore)(nil).Update), arg0, arg1, arg2)
}