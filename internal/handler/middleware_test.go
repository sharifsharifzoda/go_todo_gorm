package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
	mock_service "todo_gorm/internal/service/mocks"
)

func TestHandler_tokenAuthMiddleware(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, token string)
	//
	testTable := []struct {
		name                 string
		headerName           string
		headerValue          string
		token                string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Ok",
			headerName:  "Authorization",
			headerValue: "Alif token",
			token:       "token",
			mockBehavior: func(s *mock_service.MockAuthorization, token string) {
				s.EXPECT().ParseToken(token).Return(1, nil).AnyTimes()
			},
			expectedStatusCode:   200,
			expectedResponseBody: "1",
		},
		{
			name:                 "Empty header",
			headerName:           "Authorization",
			headerValue:          "",
			token:                "",
			mockBehavior:         func(s *mock_service.MockAuthorization, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"reason":"empty auth header"}`,
		},
		{
			name:                 "Invalid Bearer",
			headerName:           "Authorization",
			headerValue:          "alif token",
			token:                "",
			mockBehavior:         func(s *mock_service.MockAuthorization, token string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"reason":"invalid auth header"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.token)

			handler := NewHandler(auth, nil)

			gin.SetMode(gin.ReleaseMode)
			r := gin.New()
			r.POST("/middleware", handler.tokenAuthMiddleware, func(c *gin.Context) {
				id, _ := c.Get("userId")
				c.String(200, fmt.Sprintf("%d", id.(int)))
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/middleware", nil)
			req.Header.Set(testCase.headerName, testCase.headerValue)

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)
		})
	}
}
