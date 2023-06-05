package handler

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http/httptest"
	"strings"
	"testing"
	mock_service "todo_gorm/internal/service/mocks"
	"todo_gorm/model"
)

func TestHandler_createTask(t *testing.T) {
	type mockBehavior func(todo *mock_service.MockTodoTask, a *mock_service.MockAuthorization,
		task model.Task, token string)

	testTable := []struct {
		name                 string
		inputTask            model.Task
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
		headerName           string
		headerValue          string
	}{
		{
			name: "OK",
			inputTask: model.Task{
				Name:        "go",
				Description: "study harder",
				Deadline:    "2023-06-22",
			},
			mockBehavior: func(todo *mock_service.MockTodoTask, a *mock_service.MockAuthorization,
				task model.Task, token string) {
				a.EXPECT().ParseToken(token).Return(1, nil).AnyTimes()
				todo.EXPECT().CreateTask(1, &task).Return(1, nil).AnyTimes()
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1}`,
			headerName:           "Authorization",
			headerValue:          "Alif token",
		},
		{
			name: "Empty fields",
			inputTask: model.Task{
				Name:        "",
				Description: "",
				Deadline:    "",
			},
			mockBehavior: func(todo *mock_service.MockTodoTask, a *mock_service.MockAuthorization,
				task model.Task, token string) {
				a.EXPECT().ParseToken(token).Return(1, nil).AnyTimes()
				todo.EXPECT().CreateTask(1, &task).Return(1, nil).AnyTimes()
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid task format provided"}`,
			headerName:           "Authorization",
			headerValue:          "Alif token",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			todo := mock_service.NewMockTodoTask(c)
			auth := mock_service.NewMockAuthorization(c)

			split := strings.Split(testCase.headerValue, " ")
			testCase.mockBehavior(todo, auth, testCase.inputTask, split[1])

			handler := NewHandler(auth, todo)

			gin.SetMode(gin.ReleaseMode)
			r := gin.New()
			r.POST("/task", handler.tokenAuthMiddleware, handler.createTask)

			marshal, err := json.Marshal(testCase.inputTask)
			if err != nil {
				log.Fatalf("error while marshaling. error is %v", err.Error())
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/task", bytes.NewBuffer(marshal))
			req.Header.Set(testCase.headerName, testCase.headerValue)

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)
		})
	}
}

func TestHandler_getTasks(t *testing.T) {
	type mockBehavior func(a *mock_service.MockAuthorization, td *mock_service.MockTodoTask, token string)

	testTable := []struct {
		name                 string
		mockBehavior         mockBehavior
		expectedStatusCodes  int
		expectedResponseBody string
		headerName           string
		headerValue          string
	}{
		{
			name: "OK",
			mockBehavior: func(a *mock_service.MockAuthorization, td *mock_service.MockTodoTask, token string) {
				a.EXPECT().ParseToken(token).Return(1, nil).AnyTimes()
				td.EXPECT().GetAll(1).Return(model.Tasks{
					{Id: 1, Name: "mathematics", Description: "study hard", Done: false, IsActive: true,
						Deadline: "2022-05-23T00:00:00Z", Username: "Sharif"},
					{Id: 2, Name: "go", Description: "study smart", Done: false, IsActive: true,
						Deadline: "2022-06-22T00:00:00Z", Username: "Sharif"},
				}, nil).AnyTimes()
			},
			expectedStatusCodes: 200,
			expectedResponseBody: `{"tasks":[{"id":1,"name":"mathematics","description":"study hard","done":false,` +
				`"is_active":true,"deadline":"2022-05-23T00:00:00Z","username":"Sharif"},` +
				`{"id":2,"name":"go","description":"study smart","done":false,"is_active":true,` +
				`"deadline":"2022-06-22T00:00:00Z","username":"Sharif"}]}`,
			headerName:  "Authorization",
			headerValue: "Alif token",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			todo := mock_service.NewMockTodoTask(c)

			split := strings.Split(testCase.headerValue, " ")

			testCase.mockBehavior(auth, todo, split[1])

			handler := NewHandler(auth, todo)

			gin.SetMode(gin.ReleaseMode)
			r := gin.New()
			r.GET("/getTasks", handler.tokenAuthMiddleware, handler.getTasks)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/getTasks", nil)
			req.Header.Set(testCase.headerName, testCase.headerValue)

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCodes)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)
		})
	}
}
