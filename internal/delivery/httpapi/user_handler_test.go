package httpapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/giicoo/go-auth-service/internal/entity"
	"github.com/giicoo/go-auth-service/internal/entity/models"
	"github.com/giicoo/go-auth-service/internal/services"
	mock_services "github.com/giicoo/go-auth-service/internal/services/mocks"
	"github.com/giicoo/go-auth-service/pkg/apiError"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func marshal(t *testing.T, o interface{}) string {
	b, err := json.Marshal(o)
	if err != nil {
		t.Logf("marshal: %s", err.Error())
	}
	return string(b) + "\n"
}

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
			expectedStatusCode: 200,
			expectedResponseBody: marshal(t, models.UserResponse{
				ID:    1,
				Email: "da",
			}),
		},

		{
			name:               "Invalid JSON",
			inputBody:          `{"email":"da" "password":"123"}`,
			inputUser:          &entity.User{},
			mockBehavior:       func(r *mock_services.MockUserService, user *entity.User) {},
			expectedStatusCode: 400,
			expectedResponseBody: marshal(t, models.ErrorResponse{
				Error: "decode json: invalid character '\"' after object key:value pair",
			}),
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
			expectedStatusCode: 400,
			expectedResponseBody: marshal(t, models.ErrorResponse{
				Error: "service create user: user already exist",
			}),
		},
		{
			name:      "Internal error",
			inputBody: `{"email":"da", "password":"123"}`,
			inputUser: &entity.User{
				Email:    "da",
				Password: "123",
			},
			mockBehavior: func(r *mock_services.MockUserService, user *entity.User) {
				r.EXPECT().CreateUser(user).Return(nil, errors.New("any"))
			},
			expectedStatusCode: 500,
			expectedResponseBody: marshal(t, models.ErrorResponse{
				Error: "internal error",
			}),
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

func TestDeleteUser(t *testing.T) {
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
				r.EXPECT().DeleteUser(user).Return(nil)
			},
			expectedStatusCode: 200,
			expectedResponseBody: marshal(t, models.Response{
				Message: "user with da email success delete",
			}),
		},

		{
			name:      "Invalid JSON",
			inputBody: `{"email":"da" "password":"123"}`,
			inputUser: &entity.User{
				Email:    "da",
				Password: "123",
			},
			mockBehavior:       func(r *mock_services.MockUserService, user *entity.User) {},
			expectedStatusCode: 400,
			expectedResponseBody: marshal(t, models.ErrorResponse{
				Error: "decode json: invalid character '\"' after object key:value pair",
			}),
		},
		{
			name:      "User not exist",
			inputBody: `{"email":"da", "password":"123"}`,
			inputUser: &entity.User{
				Email:    "da",
				Password: "123",
			},
			mockBehavior: func(r *mock_services.MockUserService, user *entity.User) {
				r.EXPECT().DeleteUser(user).Return(apiError.ErrUserNotExists)
			},
			expectedStatusCode: 400,
			expectedResponseBody: marshal(t, models.ErrorResponse{
				Error: "delete user: user not exist",
			}),
		},
		{
			name:      "Internal error",
			inputBody: `{"email":"da", "password":"123"}`,
			inputUser: &entity.User{
				Email:    "da",
				Password: "123",
			},
			mockBehavior: func(r *mock_services.MockUserService, user *entity.User) {
				r.EXPECT().DeleteUser(user).Return(errors.New("any"))
			},
			expectedStatusCode: 500,
			expectedResponseBody: marshal(t, models.ErrorResponse{
				Error: "internal error",
			}),
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
			r.HandleFunc("/delete-user", h.DeleteUser).Methods("DELETE")

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", "/delete-user",
				bytes.NewBufferString(test.inputBody))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, w.Body.String())
		})
	}

}

func TestUpdateUser(t *testing.T) {
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
			inputBody: `{"id":1, "email":"dad", "password":"123"}`,
			inputUser: &entity.User{
				ID:       1,
				Email:    "dad",
				Password: "123",
			},
			mockBehavior: func(r *mock_services.MockUserService, user *entity.User) {
				r.EXPECT().UpdateUser(user).Return(&entity.User{
					ID:       1,
					Email:    "dad",
					Password: "123",
				}, nil)
			},
			expectedStatusCode: 200,
			expectedResponseBody: marshal(t, models.UserResponse{
				ID:    1,
				Email: "dad",
			}),
		},
		{
			name:               "Invalid Json",
			inputBody:          `{"id":1, "email":"dad" "password":"123"}`,
			inputUser:          &entity.User{},
			mockBehavior:       func(r *mock_services.MockUserService, user *entity.User) {},
			expectedStatusCode: 400,
			expectedResponseBody: marshal(t, models.ErrorResponse{
				Error: "decode json: invalid character '\"' after object key:value pair",
			}),
		},
		{
			name:      "User not exist",
			inputBody: `{"id":1, "email":"dad", "password":"123"}`,
			inputUser: &entity.User{
				ID:       1,
				Email:    "dad",
				Password: "123",
			},
			mockBehavior: func(r *mock_services.MockUserService, user *entity.User) {
				r.EXPECT().UpdateUser(user).Return(nil, apiError.ErrUserNotExists)
			},
			expectedStatusCode: 400,
			expectedResponseBody: marshal(t, models.ErrorResponse{
				Error: "update user: user not exist",
			}),
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
			r.HandleFunc("/update-user", h.UpdateUser).Methods("PUT")

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/update-user",
				bytes.NewBufferString(test.inputBody))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, w.Body.String())
		})
	}
}
