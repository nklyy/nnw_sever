// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	bytes "bytes"
	model "nnw_s/pkg/model"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	otp "github.com/pquerna/otp"
)

// MockAuthorization is a mock of Authorization interface.
type MockAuthorization struct {
	ctrl     *gomock.Controller
	recorder *MockAuthorizationMockRecorder
}

// MockAuthorizationMockRecorder is the mock recorder for MockAuthorization.
type MockAuthorizationMockRecorder struct {
	mock *MockAuthorization
}

// NewMockAuthorization creates a new mock instance.
func NewMockAuthorization(ctrl *gomock.Controller) *MockAuthorization {
	mock := &MockAuthorization{ctrl: ctrl}
	mock.recorder = &MockAuthorizationMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthorization) EXPECT() *MockAuthorizationMockRecorder {
	return m.recorder
}

// Check2FaCode mocks base method.
func (m *MockAuthorization) Check2FaCode(code, secret string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Check2FaCode", code, secret)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Check2FaCode indicates an expected call of Check2FaCode.
func (mr *MockAuthorizationMockRecorder) Check2FaCode(code, secret interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Check2FaCode", reflect.TypeOf((*MockAuthorization)(nil).Check2FaCode), code, secret)
}

// CheckPassword mocks base method.
func (m *MockAuthorization) CheckPassword(password, hashPassword string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckPassword", password, hashPassword)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckPassword indicates an expected call of CheckPassword.
func (mr *MockAuthorizationMockRecorder) CheckPassword(password, hashPassword interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckPassword", reflect.TypeOf((*MockAuthorization)(nil).CheckPassword), password, hashPassword)
}

// CreateJWTToken mocks base method.
func (m *MockAuthorization) CreateJWTToken(email string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateJWTToken", email)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateJWTToken indicates an expected call of CreateJWTToken.
func (mr *MockAuthorizationMockRecorder) CreateJWTToken(email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateJWTToken", reflect.TypeOf((*MockAuthorization)(nil).CreateJWTToken), email)
}

// CreateTemplateUserData mocks base method.
func (m *MockAuthorization) CreateTemplateUserData(secret string) (*string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTemplateUserData", secret)
	ret0, _ := ret[0].(*string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTemplateUserData indicates an expected call of CreateTemplateUserData.
func (mr *MockAuthorizationMockRecorder) CreateTemplateUserData(secret interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTemplateUserData", reflect.TypeOf((*MockAuthorization)(nil).CreateTemplateUserData), secret)
}

// CreateUser mocks base method.
func (m *MockAuthorization) CreateUser(email, password, OTPKey string) (*string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", email, password, OTPKey)
	ret0, _ := ret[0].(*string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockAuthorizationMockRecorder) CreateUser(email, password, OTPKey interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockAuthorization)(nil).CreateUser), email, password, OTPKey)
}

// Generate2FaImage mocks base method.
func (m *MockAuthorization) Generate2FaImage(email string) (*bytes.Buffer, *otp.Key, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Generate2FaImage", email)
	ret0, _ := ret[0].(*bytes.Buffer)
	ret1, _ := ret[1].(*otp.Key)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Generate2FaImage indicates an expected call of Generate2FaImage.
func (mr *MockAuthorizationMockRecorder) Generate2FaImage(email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Generate2FaImage", reflect.TypeOf((*MockAuthorization)(nil).Generate2FaImage), email)
}

// GetTemplateUserDataById mocks base method.
func (m *MockAuthorization) GetTemplateUserDataById(uid string) (*model.TemplateData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTemplateUserDataById", uid)
	ret0, _ := ret[0].(*model.TemplateData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTemplateUserDataById indicates an expected call of GetTemplateUserDataById.
func (mr *MockAuthorizationMockRecorder) GetTemplateUserDataById(uid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTemplateUserDataById", reflect.TypeOf((*MockAuthorization)(nil).GetTemplateUserDataById), uid)
}

// GetUserByEmail mocks base method.
func (m *MockAuthorization) GetUserByEmail(email string) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByEmail", email)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByEmail indicates an expected call of GetUserByEmail.
func (mr *MockAuthorizationMockRecorder) GetUserByEmail(email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByEmail", reflect.TypeOf((*MockAuthorization)(nil).GetUserByEmail), email)
}

// GetUserById mocks base method.
func (m *MockAuthorization) GetUserById(userId string) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserById", userId)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserById indicates an expected call of GetUserById.
func (mr *MockAuthorizationMockRecorder) GetUserById(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserById", reflect.TypeOf((*MockAuthorization)(nil).GetUserById), userId)
}

// VerifyJWTToken mocks base method.
func (m *MockAuthorization) VerifyJWTToken(id string) (*string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyJWTToken", id)
	ret0, _ := ret[0].(*string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// VerifyJWTToken indicates an expected call of VerifyJWTToken.
func (mr *MockAuthorizationMockRecorder) VerifyJWTToken(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyJWTToken", reflect.TypeOf((*MockAuthorization)(nil).VerifyJWTToken), id)
}