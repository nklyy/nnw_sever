package handler

import (
	"bytes"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"nnw_s/pkg/auth/service"
	"nnw_s/pkg/auth/service/mocks"
	"nnw_s/pkg/user/model"
	"testing"
)

func TestHandler_verifyLoginCode(t *testing.T) {
	type mockBehavior func(r *mock_service.MockAuthorization, tUser *model.User)

	tests := []struct {
		name                string
		inputBody           string
		inputUser           *model.User
		urlPath             string
		mockBehavior        mockBehavior
		expectedStatus      int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"email": "email", "password": "passwordA1", "code": "code"}`,
			inputUser: &model.User{
				Email:        "email",
				Password:     "password",
				SecretOTPKey: "secret",
			},
			urlPath: "/v1/verifyLogin2fa",
			mockBehavior: func(r *mock_service.MockAuthorization, User *model.User) {
				r.EXPECT().GetUserByEmail(User.Email).Return(User, nil)
				r.EXPECT().Check2FaCode("code", User.SecretOTPKey).Return(true)
				r.EXPECT().CreateJWTToken(User.Email).Return("token", nil)
			},
			expectedStatus: 200,
			expectedRequestBody: `{"token":"token"}
`,
		},
		{
			name:           "Invalid json",
			inputBody:      `{"email"}`,
			urlPath:        "/v1/verifyLogin2fa",
			mockBehavior:   func(r *mock_service.MockAuthorization, User *model.User) {},
			expectedStatus: 400,
			expectedRequestBody: `{"error":" Invalid json!"}
`,
		},
		{
			name:           "Required Login field!",
			inputBody:      `{"emaillll": "email", "passwordd": "passwordA1", "codee": "code"}`,
			urlPath:        "/v1/verifyLogin2fa",
			mockBehavior:   func(r *mock_service.MockAuthorization, User *model.User) {},
			expectedStatus: 400,
			expectedRequestBody: `{"error":["Email is a required field","Password is a required field","Code is a required field"]}
`,
		},
		{
			name:      "Not found user!",
			inputBody: `{"email": "email", "password": "passwordA1", "code": "code"}`,
			inputUser: &model.User{
				Email:        "email",
				Password:     "password",
				SecretOTPKey: "secret",
			},
			urlPath: "/v1/verifyLogin2fa",
			mockBehavior: func(r *mock_service.MockAuthorization, User *model.User) {
				r.EXPECT().GetUserByEmail(User.Email).Return(nil, errors.New(" User not found!"))
			},
			expectedStatus: 400,
			expectedRequestBody: `{"error":" User not found!"}
`,
		},
		{
			name:      "Invalid code!",
			inputBody: `{"email": "email", "password": "passwordA1", "code": "code"}`,
			inputUser: &model.User{
				Email:        "email",
				Password:     "password",
				SecretOTPKey: "secret",
			},
			urlPath: "/v1/verifyLogin2fa",
			mockBehavior: func(r *mock_service.MockAuthorization, User *model.User) {
				r.EXPECT().GetUserByEmail(User.Email).Return(User, nil)
				r.EXPECT().Check2FaCode("code", User.SecretOTPKey).Return(false)
			},
			expectedStatus: 400,
			expectedRequestBody: `{"error":" Invalid code!"}
`,
		},
		{
			name:      "Fail create JWT!",
			inputBody: `{"email": "email", "password": "passwordA1", "code": "code"}`,
			inputUser: &model.User{
				Email:        "email",
				Password:     "password",
				SecretOTPKey: "secret",
			},
			urlPath: "/v1/verifyLogin2fa",
			mockBehavior: func(r *mock_service.MockAuthorization, User *model.User) {
				r.EXPECT().GetUserByEmail(User.Email).Return(User, nil)
				r.EXPECT().Check2FaCode("code", User.SecretOTPKey).Return(true)
				r.EXPECT().CreateJWTToken(User.Email).Return("", errors.New(" Something wrong!"))
			},
			expectedStatus: 500,
			expectedRequestBody: `{"error":" Something wrong!"}
`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			repo := mock_service.NewMockAuthorization(controller)
			test.mockBehavior(repo, test.inputUser)

			services := &service.Service{Authorization: repo}

			v := validator.New()
			handler := Handler{authService: services, validate: v}

			app := echo.New()
			handler.InitialRoute(app)

			req := httptest.NewRequest(http.MethodPost, test.urlPath, bytes.NewBufferString(test.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := app.NewContext(req, rec)

			if assert.NoError(t, handler.verifyLogin2fa(c)) {
				assert.Equal(t, test.expectedStatus, rec.Code)
				assert.Equal(t, test.expectedRequestBody, rec.Body.String())
			}
		})
	}
}

// Check Login
func TestHandler_checkEmail(t *testing.T) {
	type mockBehavior func(r *mock_service.MockAuthorization, tUser *model.User)

	tests := []struct {
		name                string
		inputBody           string
		inputUser           *model.User
		urlPath             string
		mockBehavior        mockBehavior
		expectedStatus      int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"email": "email"}`,
			inputUser: &model.User{
				Email: "email",
			},
			urlPath: "/v1/checkEmail",
			mockBehavior: func(r *mock_service.MockAuthorization, User *model.User) {
				r.EXPECT().GetUserByEmail(User.Email).Return(nil, nil)
			},
			expectedStatus:      200,
			expectedRequestBody: "",
		},
		{
			name:      "FAIL",
			inputBody: `{"email": "email"}`,
			inputUser: &model.User{
				Email: "email",
			},
			urlPath: "/v1/checkEmail",
			mockBehavior: func(r *mock_service.MockAuthorization, User *model.User) {
				r.EXPECT().GetUserByEmail(User.Email).Return(User, nil)
			},
			expectedStatus:      400,
			expectedRequestBody: "",
		},
		{
			name:           "Required Email field!",
			inputBody:      `{"emailll": "email"}`,
			urlPath:        "/v1/checkEmail",
			mockBehavior:   func(r *mock_service.MockAuthorization, User *model.User) {},
			expectedStatus: 400,
			expectedRequestBody: `{"error":["Email is a required field"]}
`,
		},
		{
			name:           "Invalid json",
			inputBody:      `{"email"}`,
			urlPath:        "/v1/checkEmail",
			mockBehavior:   func(r *mock_service.MockAuthorization, User *model.User) {},
			expectedStatus: 400,
			expectedRequestBody: `{"error":" Invalid json!"}
`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			repo := mock_service.NewMockAuthorization(controller)
			test.mockBehavior(repo, test.inputUser)

			services := &service.Service{Authorization: repo}

			v := validator.New()
			handler := Handler{authService: services, validate: v}

			app := echo.New()
			handler.InitialRoute(app)

			req := httptest.NewRequest(http.MethodPost, test.urlPath, bytes.NewBufferString(test.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := app.NewContext(req, rec)

			if assert.NoError(t, handler.checkEmail(c)) {
				assert.Equal(t, test.expectedStatus, rec.Code)
				assert.Equal(t, test.expectedRequestBody, rec.Body.String())
			}
		})
	}
}

// Check Jwt Token
func TestHandler_checkJwt(t *testing.T) {
	type mockBehavior func(r *mock_service.MockAuthorization)

	tests := []struct {
		name                string
		inputBody           string
		urlPath             string
		mockBehavior        mockBehavior
		expectedStatus      int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"token": "token"}`,
			urlPath:   "/v1/checkJwt",
			mockBehavior: func(r *mock_service.MockAuthorization) {
				r.EXPECT().VerifyJWTToken("token").Return(nil, nil)
			},
			expectedStatus:      200,
			expectedRequestBody: "",
		},
		{
			name:      "FAIL",
			inputBody: `{"token": "token"}`,
			urlPath:   "/v1/checkJwt",
			mockBehavior: func(r *mock_service.MockAuthorization) {
				r.EXPECT().VerifyJWTToken("token").Return(nil, errors.New(" Wrong token!"))
			},
			expectedStatus: 400,
			expectedRequestBody: `{"error":" Wrong token!"}
`,
		},
		{
			name:           "Required Login field!",
			inputBody:      `{"tokennn": "token"}`,
			urlPath:        "/v1/checkJwt",
			mockBehavior:   func(r *mock_service.MockAuthorization) {},
			expectedStatus: 400,
			expectedRequestBody: `{"error":["Token is a required field"]}
`,
		},
		{
			name:           "Invalid json",
			inputBody:      `{"token"}`,
			urlPath:        "/v1/checkJwt",
			mockBehavior:   func(r *mock_service.MockAuthorization) {},
			expectedStatus: 400,
			expectedRequestBody: `{"error":" Invalid json!"}
`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			repo := mock_service.NewMockAuthorization(controller)
			test.mockBehavior(repo)

			services := &service.Service{Authorization: repo}

			v := validator.New()
			handler := Handler{authService: services, validate: v}

			app := echo.New()
			handler.InitialRoute(app)

			req := httptest.NewRequest(http.MethodPost, test.urlPath, bytes.NewBufferString(test.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := app.NewContext(req, rec)

			if assert.NoError(t, handler.checkJwt(c)) {
				assert.Equal(t, test.expectedStatus, rec.Code)
				assert.Equal(t, test.expectedRequestBody, rec.Body.String())
			}
		})
	}
}
