// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go

// Package mock_jwt is a generated GoMock package.
package mock_jwt

import (
	context "context"
	jwt "nnw_s/internal/auth/jwt"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// GetJWT mocks base method.
func (m *MockRepository) GetJWT(ctx context.Context, id string) (*jwt.JWT, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetJWT", ctx, id)
	ret0, _ := ret[0].(*jwt.JWT)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetJWT indicates an expected call of GetJWT.
func (mr *MockRepositoryMockRecorder) GetJWT(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetJWT", reflect.TypeOf((*MockRepository)(nil).GetJWT), ctx, id)
}

// SaveJWT mocks base method.
func (m *MockRepository) SaveJWT(ctx context.Context, jwt *jwt.JWT) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveJWT", ctx, jwt)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SaveJWT indicates an expected call of SaveJWT.
func (mr *MockRepositoryMockRecorder) SaveJWT(ctx, jwt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveJWT", reflect.TypeOf((*MockRepository)(nil).SaveJWT), ctx, jwt)
}
