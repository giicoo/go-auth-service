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

func (h *Handler) CreateSession(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	defer body.Close()

	sessionJSON := new(models.SessionCreateRequest)
	if err := json.NewDecoder(body).Decode(sessionJSON); err != nil {
		logrus.WithFields(
			logrus.Fields{
				"url": r.URL.String(),
			},
		).Errorf("decode json: %s", err)
		httpError(w, fmt.Errorf("decode json: %w", apiError.New(err.Error(), http.StatusBadRequest)))
		return
	}
	session := &entity.Session{
		UserID:    sessionJSON.UserID,
		UserAgent: sessionJSON.UserAgent,
		UserIP:    sessionJSON.UserIP,
	}

	sessionDB, err := h.Services.SessionService.CreateSession(session)
	if err != nil {
		logrus.WithFields(
			logrus.Fields{
				"url":     r.URL.String(),
				"request": sessionJSON,
			},
		).Errorf("service create session: %s", err)
		httpError(w, fmt.Errorf("service create session: %w", err))
		return
	}

	sessionResponse := models.SessionResponse{
		ID: sessionDB.ID,
	}

	if err := httpResponse(w, sessionResponse); err != nil {
		logrus.WithFields(
			logrus.Fields{
				"url":     r.URL.String(),
				"request": sessionJSON,
			},
		).Errorf("send response: %s", err)
		httpError(w, fmt.Errorf("send response: %w", err))
		return
	}
}

func (h *Handler) GetSession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	userID, err := strconv.Atoi(vars["user_id"])
	if err != nil {
		logrus.WithFields(
			logrus.Fields{
				"url": r.URL.String(),
				"request": map[string]interface{}{
					"id":      id,
					"user_id": userID,
				},
			},
		).Errorf("service get session: %s", err)
		httpError(w, fmt.Errorf("service get session: %w", err))
		return
	}

	sessionDB, err := h.Services.SessionService.GetSession(id, userID)
	if err != nil {
		logrus.WithFields(
			logrus.Fields{
				"url": r.URL.String(),
				"request": map[string]interface{}{
					"id":      id,
					"user_id": userID,
				},
			},
		).Errorf("service get session: %s", err)
		httpError(w, fmt.Errorf("service get session: %w", err))
		return
	}

	sessionResponse := models.SessionResponseFull{
		ID:        sessionDB.ID,
		UserID:    sessionDB.UserID,
		UserAgent: sessionDB.UserAgent,
		UserIP:    sessionDB.UserIP,
	}

	if err := httpResponse(w, sessionResponse); err != nil {
		logrus.WithFields(
			logrus.Fields{
				"url": r.URL.String(),
				"request": map[string]interface{}{
					"id":      id,
					"user_id": userID,
				},
			},
		).Errorf("send response: %s", err)
		httpError(w, fmt.Errorf("send response: %w", err))
		return
	}
}

func (h *Handler) DeleteSession(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	defer body.Close()

	sessionJSON := new(models.SessionDeleteRequest)
	if err := json.NewDecoder(body).Decode(sessionJSON); err != nil {
		logrus.WithFields(
			logrus.Fields{
				"url": r.URL.String(),
			},
		).Errorf("decode json: %s", err)
		httpError(w, fmt.Errorf("decode json: %w", apiError.New(err.Error(), http.StatusBadRequest)))
		return
	}

	err := h.Services.SessionService.DeleteSession(sessionJSON.ID, sessionJSON.UserID)
	if err != nil {
		logrus.WithFields(
			logrus.Fields{
				"url":     r.URL.String(),
				"request": sessionJSON,
			},
		).Errorf("service delete session: %s", err)
		httpError(w, fmt.Errorf("service delete session: %w", err))
		return
	}

	response := models.Response{
		Message: fmt.Sprintf("session '%s:%d' deleted", sessionJSON.ID, sessionJSON.UserID),
	}

	if err := httpResponse(w, response); err != nil {
		logrus.WithFields(
			logrus.Fields{
				"url":     r.URL.String(),
				"request": sessionJSON,
			},
		).Errorf("send response: %s", err)
		httpError(w, fmt.Errorf("send response: %w", err))
		return
	}
}

func (h *Handler) GetListSession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["user_id"])
	if err != nil {
		logrus.WithFields(
			logrus.Fields{
				"url": r.URL.String(),
				"request": map[string]interface{}{
					"user_id": userID,
				},
			},
		).Errorf("service get session: %s", err)
		httpError(w, fmt.Errorf("service get session: %w", err))
		return
	}

	sessionsDB, err := h.Services.SessionService.GetListSession(userID)
	if err != nil {
		logrus.WithFields(
			logrus.Fields{
				"url": r.URL.String(),
				"request": map[string]interface{}{
					"user_id": userID,
				},
			},
		).Errorf("service get session: %s", err)
		httpError(w, fmt.Errorf("service get session: %w", err))
		return
	}

	response := []*models.SessionResponseFull{}
	for _, sessionDB := range sessionsDB {
		sessionResponse := &models.SessionResponseFull{
			ID:        sessionDB.ID,
			UserID:    sessionDB.UserID,
			UserAgent: sessionDB.UserAgent,
			UserIP:    sessionDB.UserIP,
		}
		response = append(response, sessionResponse)
	}

	if err := httpResponse(w, response); err != nil {
		logrus.WithFields(
			logrus.Fields{
				"url": r.URL.String(),
				"request": map[string]interface{}{
					"user_id": userID,
				},
			},
		).Errorf("send response: %s", err)
		httpError(w, fmt.Errorf("send response: %w", err))
		return
	}
}

func (h *Handler) DeleteListSession(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	defer body.Close()

	sessionJSON := new(models.SessionDeleteRequest)
	if err := json.NewDecoder(body).Decode(sessionJSON); err != nil {
		logrus.WithFields(
			logrus.Fields{
				"url": r.URL.String(),
			},
		).Errorf("decode json: %s", err)
		httpError(w, fmt.Errorf("decode json: %w", apiError.New(err.Error(), http.StatusBadRequest)))
		return
	}

	err := h.Services.SessionService.DeleteListSession(sessionJSON.UserID)
	if err != nil {
		logrus.WithFields(
			logrus.Fields{
				"url":     r.URL.String(),
				"request": sessionJSON,
			},
		).Errorf("service delete list session: %s", err)
		httpError(w, fmt.Errorf("service delete list session: %w", err))
		return
	}

	response := models.Response{
		Message: fmt.Sprintf("list session '%d' deleted", sessionJSON.UserID),
	}

	if err := httpResponse(w, response); err != nil {
		logrus.WithFields(
			logrus.Fields{
				"url":     r.URL.String(),
				"request": sessionJSON,
			},
		).Errorf("send response: %s", err)
		httpError(w, fmt.Errorf("send response: %w", err))
		return
	}
}
