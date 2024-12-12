package httpapi

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/giicoo/go-auth-service/internal/entity"
	"github.com/giicoo/go-auth-service/internal/services"
	mock_services "github.com/giicoo/go-auth-service/internal/services/mocks"
	"github.com/giicoo/go-auth-service/pkg/apiError"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	type mockBehavior func(r *mock_services.MockUserService, user *entity.User)

	tests := []struct {
		name                 string
		inputBody            string
		inputUser            *entity.User
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"email":"da", "password":"123"}`,
			inputUser: &entity.User{
				Email:    "da",
				Password: "123",
			},
			mockBehavior: func(r *mock_services.MockUserService, user *entity.User) {
				r.EXPECT().CreateUser(user).Return(&entity.User{
					ID:       1,
					Email:    "da",
					Password: "123",
				}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "{\"id\":1,\"email\":\"da\"}\n",
		},

		{
			name:                 "Invalid JSON",
			inputBody:            `{"email":"da" "password":"123"}`,
			inputUser:            &entity.User{},
			mockBehavior:         func(r *mock_services.MockUserService, user *entity.User) {},
			expectedStatusCode:   400,
			expectedResponseBody: "decode json: invalid character '\"' after object key:value pair\n",
		},
		{
			name:      "User already exist",
			inputBody: `{"email":"da", "password":"123"}`,
			inputUser: &entity.User{
				Email:    "da",
				Password: "123",
			},
			mockBehavior: func(r *mock_services.MockUserService, user *entity.User) {
				r.EXPECT().CreateUser(user).Return(nil, apiError.ErrUserAlreadyExists)
			},
			expectedStatusCode:   400,
			expectedResponseBody: "service create user: user already exist\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			service := mock_services.NewMockUserService(c)
			test.mockBehavior(service, test.inputUser)

			services := &services.Services{
				UserService: service,
			}
			h := NewHandler(services)

			r := mux.NewRouter()
			r.HandleFunc("/create-user", h.CreateUser).Methods("POST")

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/create-user",
				bytes.NewBufferString(test.inputBody))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, w.Body.String())
		})
	}

}
