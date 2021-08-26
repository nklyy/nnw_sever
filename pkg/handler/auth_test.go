package handler

import (
	"bytes"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"nnw_s/pkg/model"
	"nnw_s/pkg/service"
	mock_service "nnw_s/pkg/service/mocks"
	"testing"
)

func TestHandler_checkUserName(t *testing.T) {
	type mockBehavior func(r *mock_service.MockAuthorization, tUser *model.User)

	tests := []struct {
		name           string
		inputBody      string
		inputUser      *model.User
		urlPath        string
		mockBehavior   mockBehavior
		expectedStatus int
	}{
		{
			name:      "OK",
			inputBody: `{"login": "login"}`,
			inputUser: &model.User{
				Login: "login",
			},
			urlPath: "/v1/checkLogin",
			mockBehavior: func(r *mock_service.MockAuthorization, User *model.User) {
				r.EXPECT().GetUserByLogin(User.Login).Return(nil, nil)
			},
			expectedStatus: 200,
		},
		{
			name:      "FAIL",
			inputBody: `{"login": "login"}`,
			inputUser: &model.User{
				Login: "login",
			},
			urlPath: "/v1/checkLogin",
			mockBehavior: func(r *mock_service.MockAuthorization, User *model.User) {
				r.EXPECT().GetUserByLogin(User.Login).Return(User, nil)
			},
			expectedStatus: 400,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			repo := mock_service.NewMockAuthorization(controller)
			test.mockBehavior(repo, test.inputUser)

			services := &service.Service{Authorization: repo}
			handler := Handler{services: services}

			app := echo.New()
			handler.InitialRoute(app)

			req := httptest.NewRequest(http.MethodPost, test.urlPath, bytes.NewBufferString(test.inputBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := app.NewContext(req, rec)

			if assert.NoError(t, handler.checkLogin(c)) {
				assert.Equal(t, test.expectedStatus, rec.Code)
			}
		})
	}
}
