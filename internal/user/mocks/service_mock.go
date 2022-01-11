// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock_user is a generated GoMock package.
package mock_user

import (
	context "context"
	user "nnw_s/internal/user"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
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

// CreateUser mocks base method.
func (m *MockService) CreateUser(ctx context.Context, dto *user.CreateUserDTO) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", ctx, dto)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockServiceMockRecorder) CreateUser(ctx, dto interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockService)(nil).CreateUser), ctx, dto)
}

// DeleteUserByEmail mocks base method.
func (m *MockService) DeleteUserByEmail(ctx context.Context, email string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUserByEmail", ctx, email)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUserByEmail indicates an expected call of DeleteUserByEmail.
func (mr *MockServiceMockRecorder) DeleteUserByEmail(ctx, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUserByEmail", reflect.TypeOf((*MockService)(nil).DeleteUserByEmail), ctx, email)
}

// GetUserByEmail mocks base method.
func (m *MockService) GetUserByEmail(ctx context.Context, email string) (*user.DTO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByEmail", ctx, email)
	ret0, _ := ret[0].(*user.DTO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByEmail indicates an expected call of GetUserByEmail.
func (mr *MockServiceMockRecorder) GetUserByEmail(ctx, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByEmail", reflect.TypeOf((*MockService)(nil).GetUserByEmail), ctx, email)
}

// GetUserByID mocks base method.
func (m *MockService) GetUserByID(ctx context.Context, userID string) (*user.DTO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByID", ctx, userID)
	ret0, _ := ret[0].(*user.DTO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByID indicates an expected call of GetUserByID.
func (mr *MockServiceMockRecorder) GetUserByID(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByID", reflect.TypeOf((*MockService)(nil).GetUserByID), ctx, userID)
}

// GetUserByWalletID mocks base method.
func (m *MockService) GetUserByWalletID(ctx context.Context, email, walletId string) (*user.DTO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByWalletID", ctx, email, walletId)
	ret0, _ := ret[0].(*user.DTO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByWalletID indicates an expected call of GetUserByWalletID.
func (mr *MockServiceMockRecorder) GetUserByWalletID(ctx, email, walletId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByWalletID", reflect.TypeOf((*MockService)(nil).GetUserByWalletID), ctx, email, walletId)
}

// UpdateUser mocks base method.
func (m *MockService) UpdateUser(ctx context.Context, dto *user.DTO) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", ctx, dto)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockServiceMockRecorder) UpdateUser(ctx, dto interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockService)(nil).UpdateUser), ctx, dto)
}
