package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/giicoo/go-auth-service/internal/entity"
	"github.com/giicoo/go-auth-service/internal/entity/models"
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
		httpError(w, err, http.StatusBadRequest)
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
		httpError(w, err, http.StatusInternalServerError)
		return
	}

	if err := httpResponse(w, userDB); err != nil {
		logrus.WithFields(
			logrus.Fields{
				"url":     r.URL.String(),
				"request": userJSON,
			},
		).Errorf("send response: %s", err)
		httpError(w, err, http.StatusInternalServerError)
		return
	}
}
