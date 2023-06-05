package handler

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
	mock_service "todo_gorm/internal/service/mocks"
	"todo_gorm/model"
)

//func TestUpdateTask(t *testing.T) {
//	//**************************************
//	cfg := configs.DatabaseConnConfig{
//		Host:     "localhost",
//		Port:     "5432",
//		User:     "alif",
//		Password: "pass",
//		DbName: "alif_db",
//	}
//	conn, err := repository.GetDBConnection(cfg)
//	if err != nil {
//		t.Fatalf("error while opening DB. error: %s", err.Error())
//	}
//
//	repo := repository.NewRepository(conn)
//	ser := service.NewService(repo)
//	newHandler := NewHandler(ser)
//	r := newHandler.InitRoutes()
//	//*******************************************
//
//	expect := `{"msg": "successfully updated"}`
//
//	task := model.Task{
//		Name:        "chemistry",
//		Description: "study hard",
//		Done:        true,
//		IsActive:    true,
//		Deadline:    "2022-05-26 16:05:24.198879",
//	}
//	id := 2
//
//	jsonValue, _ := json.Marshal(task)
//
//	req, _ := http.NewRequest("PUT", "/main/api/task/"+strconv.Itoa(id), bytes.NewBuffer(jsonValue))
//	req.Header.Set("Content-Type", "application/json")
//	//req.Header.Set("Authorization", "")
//	w := httptest.NewRecorder()
//	r.ServeHTTP(w, req)
//
//	assert.Equal(t, http.StatusOK, w.Code)
//	assert.Equal(t, expect, w.Body.String())
//
//	//request, _ := http.NewRequest("PUT", "/task/12", bytes.NewBuffer(jsonValue))
//	//w2 := httptest.NewRecorder()
//	//r.ServeHTTP(w2, request)
//	//assert.Equal(t, http.StatusBadRequest, w2.Code)
//}

//func TestHandler_signUp(t *testing.T) {
//	//**************************************
//	cfg := configs.DatabaseConnConfig{
//		Host:     "localhost",
//		Port:     "5432",
//		User:     "alif",
//		Password: "pass",
//		DbName: "alif_db",
//	}
//	conn, err := repository.GetDBConnection(cfg)
//	if err != nil {
//		t.Fatalf("error while opening DB. error: %s", err.Error())
//	}
//
//	repo := repository.NewRepository(conn)
//	ser := service.NewService(repo)
//	newHandler := NewHandler(ser.Auth, ser.Todo)
//	r := newHandler.InitRoutes()
//	//*******************************************
//
//	user := model.User{
//		Name:     "sharif",
//		Email:    "sharifov300@gmail.com",
//		Password: "password",
//	}
//	jsonStr, err := json.Marshal(user)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	req, err := http.NewRequest("POST", "/main/auth/sign-up", bytes.NewBuffer(jsonStr))
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	req.Header.Set("Content-Type", "application/json")
//	rr := httptest.NewRecorder()
//	r.ServeHTTP(rr, req)
//	if status := rr.Code; status != http.StatusCreated {
//		t.Errorf("got %v, want %v", status, http.StatusCreated)
//	}
//}

func TestHandler_signUp(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, user model.User)

	testTable := []struct {
		name                 string
		inputBody            string
		inputUser            model.User
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"name":"Test","email":"sharifov300@gmail.com","password":"pass"}`,
			inputUser: model.User{
				Name:     "Test",
				Email:    "sharifov300@gmail.com",
				Password: "pass",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user model.User) {
				s.EXPECT().CreateUser(&user).Return(1, nil).AnyTimes()
				s.EXPECT().ValidateUser(user).Return(nil).AnyTimes()
				s.EXPECT().IsEmailUsed(user.Email).Return(false).AnyTimes()
			},
			expectedStatusCode:   201,
			expectedResponseBody: `{"id":1}`,
		},
		{
			name:      "Empty fields",
			inputBody: `{"name":"","email":"","password":""}`,
			inputUser: model.User{},
			mockBehavior: func(s *mock_service.MockAuthorization, user model.User) {
				s.EXPECT().ValidateUser(user).Return(fmt.Errorf("forbidden")).AnyTimes()
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"validate"}`,
		},
		{
			name:      "Used Email",
			inputBody: `{"name":"sharif","email":"sharif@gmail.com","password":"qwerty"}`,
			inputUser: model.User{
				Name:     "sharif",
				Email:    "sharif@gmail.com",
				Password: "qwerty",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user model.User) {
				s.EXPECT().ValidateUser(user).Return(nil).AnyTimes()
				s.EXPECT().IsEmailUsed(user.Email).Return(true).AnyTimes()
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"email is already created"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.inputUser)

			//services := &service.Service{Auth: auth}
			handler := NewHandler(auth, nil)
			//handler := Handler{services.Auth, services.Todo}

			gin.SetMode(gin.ReleaseMode)
			r := gin.New()
			r.POST("/sign-up", handler.signUp)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-up",
				bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)
		})
	}
}

func TestHandler_signIn(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, user model.User)

	testTable := []struct {
		name                 string
		inputBody            string
		inputUser            model.User
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"email":"sharifov300@gmail.com","password":"pass"}`,
			inputUser: model.User{
				Email:    "sharifov300@gmail.com",
				Password: "pass",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user model.User) {
				s.EXPECT().GenerateToken(user).Return("token", nil).AnyTimes()
				s.EXPECT().CheckUser(user).Return(user, nil).AnyTimes()
				s.EXPECT().ValidateUser(user).Return(nil).AnyTimes()
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"msg":"signed in"}`,
		},
		{
			name:      "Empty fields",
			inputBody: `{"email":"","password":""}`,
			inputUser: model.User{
				Email:    "",
				Password: "",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user model.User) {
				s.EXPECT().ValidateUser(user).Return(fmt.Errorf("forbidden")).AnyTimes()
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"validation"}`,
		},
		{
			name:      "Invalid fields",
			inputBody: `{"email":"sharifov300@gmail.com","password":"pass"}`,
			inputUser: model.User{
				Email:    "sharifov300@gmail.com",
				Password: "pass",
			},
			mockBehavior: func(s *mock_service.MockAuthorization, user model.User) {
				s.EXPECT().CheckUser(user).Return(user, errors.New("invalid email or password")).AnyTimes()
				s.EXPECT().ValidateUser(user).Return(nil).AnyTimes()
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"invalid email or password"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			testCase.mockBehavior(auth, testCase.inputUser)

			handler := NewHandler(auth, nil)

			gin.SetMode(gin.ReleaseMode)
			r := gin.New()
			r.POST("/sign-in", handler.signIn)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/sign-in", bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, testCase.expectedStatusCode)
			assert.Equal(t, w.Body.String(), testCase.expectedResponseBody)
		})
	}
}
