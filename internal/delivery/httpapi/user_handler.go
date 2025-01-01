package httpapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/giicoo/go-auth-service/internal/entity"
	"github.com/giicoo/go-auth-service/internal/entity/models"
	"github.com/giicoo/go-auth-service/pkg/apiError"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	defer body.Close()

	userJSON := new(models.UserCreateRequest)
	if err := json.NewDecoder(body).Decode(userJSON); err != nil {
		logrus.WithFields(
			logrus.Fields{
				"url": r.URL.String(),
			},
		).Errorf("decode json: %s", err)
		httpError(w, fmt.Errorf("decode json: %w", apiError.New(err.Error(), http.StatusBadRequest)))
		return
	}

	user := &entity.User{
		Email:    userJSON.Email,
		Password: userJSON.Password,
	}

	userDB, err := h.Services.UserService.CreateUser(user)
	if err != nil {
		logrus.WithFields(
			logrus.Fields{
				"url":     r.URL.String(),
				"request": userJSON,
			},
		).Errorf("service create user: %s", err)
		httpError(w, fmt.Errorf("service create user: %w", err))
		return
	}

	userResponse := models.UserResponse{
		ID:    userDB.ID,
		Email: userDB.Email,
	}

	if err := httpResponse(w, userResponse); err != nil {
		logrus.WithFields(
			logrus.Fields{
				"url":     r.URL.String(),
				"request": userJSON,
			},
		).Errorf("send response: %s", err)
		httpError(w, fmt.Errorf("send response: %w", err))
		return
	}
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	defer body.Close()

	userJSON := new(models.UserDeleteRequest)
	if err := json.NewDecoder(body).Decode(userJSON); err != nil {
		logrus.WithFields(
			logrus.Fields{
				"url":     r.URL.String(),
				"request": userJSON,
			},
		).Errorf("decode json: %s", err)
		httpError(w, fmt.Errorf("decode json: %w", apiError.New(err.Error(), http.StatusBadRequest)))
		return
	}

	user := &entity.User{
		Email:    userJSON.Email,
		Password: userJSON.Password,
	}

	if err := h.Services.UserService.DeleteUser(user); err != nil {
		logrus.WithFields(
			logrus.Fields{
				"url":     r.URL.String(),
				"request": userJSON,
			},
		).Errorf("service delete user: %s", err)
		httpError(w, fmt.Errorf("delete user: %w", err))
		return
	}

	response := models.Response{
		Message: fmt.Sprintf("user with %s email success delete", user.Email),
	}
	if err := httpResponse(w, response); err != nil {
		logrus.WithFields(
			logrus.Fields{
				"url":     r.URL.String(),
				"request": userJSON,
			},
		).Errorf("send response: %s", err)
		httpError(w, fmt.Errorf("send response: %w", err))
		return
	}
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	defer body.Close()

	userJSON := new(models.UserUpdateRequest)
	if err := json.NewDecoder(body).Decode(userJSON); err != nil {
		logrus.WithFields(
			logrus.Fields{
				"url":     r.URL.String(),
				"request": userJSON,
			},
		).Errorf("decode json: %s", err)
		httpError(w, fmt.Errorf("decode json: %w", apiError.New(err.Error(), http.StatusBadRequest)))
		return
	}

	user := &entity.User{
		ID:       userJSON.ID,
		Email:    userJSON.Email,
		Password: userJSON.Password,
	}

	userDB, err := h.Services.UserService.UpdateUser(user)
	if err != nil {
		logrus.WithFields(
			logrus.Fields{
				"url":     r.URL.String(),
				"request": userJSON,
			},
		).Errorf("service update user: %s", err)
		httpError(w, fmt.Errorf("update user: %w", err))
		return
	}

	userResponse := models.UserResponse{
		ID:    userDB.ID,
		Email: userDB.Email,
	}

	if err := httpResponse(w, userResponse); err != nil {
		logrus.WithFields(
			logrus.Fields{
				"url":     r.URL.String(),
				"request": userJSON,
			},
		).Errorf("send response: %s", err)
		httpError(w, fmt.Errorf("send response: %w", err))
		return
	}
}

func (h *Handler) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	email := vars["email"]

	user := &entity.User{
		Email: email,
	}

	userDB, err := h.Services.UserService.GetUserByEmail(user)
	if err != nil {
		logrus.WithFields(
			logrus.Fields{
				"url":     r.URL.String(),
				"request": email,
			},
		).Errorf("service get user: %s", err)
		httpError(w, fmt.Errorf("get user: %w", err))
		return
	}

	userResponse := models.UserResponse{
		ID:    userDB.ID,
		Email: userDB.Email,
	}

	if err := httpResponse(w, userResponse); err != nil {
		logrus.WithFields(
			logrus.Fields{
				"url":     r.URL.String(),
				"request": email,
			},
		).Errorf("send response: %s", err)
		httpError(w, fmt.Errorf("send response: %w", err))
		return
	}
}

func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		logrus.WithFields(
			logrus.Fields{
				"url":     r.URL.String(),
				"request": id,
			},
		).Errorf("service get user: %s", err)
		httpError(w, fmt.Errorf("get user: %w", err))
		return
	}

	user := &entity.User{
		ID: id,
	}

	userDB, err := h.Services.UserService.GetUserByID(user)
	if err != nil {
		logrus.WithFields(
			logrus.Fields{
				"url":     r.URL.String(),
				"request": id,
			},
		).Errorf("service get user: %s", err)
		httpError(w, fmt.Errorf("get user: %w", err))
		return
	}

	userResponse := models.UserResponse{
		ID:    userDB.ID,
		Email: userDB.Email,
	}

	if err := httpResponse(w, userResponse); err != nil {
		logrus.WithFields(
			logrus.Fields{
				"url":     r.URL.String(),
				"request": id,
			},
		).Errorf("send response: %s", err)
		httpError(w, fmt.Errorf("send response: %w", err))
		return
	}
}
