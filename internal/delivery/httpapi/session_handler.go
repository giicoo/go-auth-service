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

// @Summary      	Create Session
// @Description  	SessionID генерируется как 32 рандомных бит Добавляется в Redis: 1) *session <session_id>{"user_id": int, "user_ip": "string", "user_agent": "string"}* 2) *user_sessions <user_id>["session_id",...]*
// @Tags         	sessions
// @Security Bearer
// @Accept			json
// @Produce			json
// @Param			user	body	models.SessionCreateRequest	true	"Write UserID UserIP UserAgent"
// @Success		 	200		{object}	models.SessionResponse
// @Failure      	400  	{object}  	models.ErrorResponse
// @Failure      	500  	{object}  	models.ErrorResponse
// @Router       	/create-session [post]
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

// @Summary      	Get Session
// @Security Bearer
// @Description  	Получение из Redis: session <session_id>
// @Tags         	sessions
// @Accept			json
// @Produce			json
// @Param        session_id   path      string  true  "Session ID"
// @Success		 	200		{object}	models.SessionResponseFull
// @Failure      	400  	{object}  	models.ErrorResponse
// @Failure      	500  	{object}  	models.ErrorResponse
// @Router       	/get-session/{session_id} [get]
func (h *Handler) GetSession(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	session_id := vars["session_id"]

	sessionDB, err := h.Services.SessionService.GetSession(session_id)
	if err != nil {
		logrus.WithFields(
			logrus.Fields{
				"url": r.URL.String(),
				"request": map[string]interface{}{
					"session_id": session_id,
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
					"session_id": session_id,
				},
			},
		).Errorf("send response: %s", err)
		httpError(w, fmt.Errorf("send response: %w", err))
		return
	}
}

// @Summary      	Delete Session
// @Security Bearer
// @Description  	Удаление сессии session_id из Redis: session <session_id> если у сессий совпадает user_id Удаление сессии из списка сессий юзера
// @Tags         	sessions
// @Accept			json
// @Produce			json
// @Param			user	body	models.SessionRequest	true	"id"
// @Success		 	200		{object}	models.Response
// @Failure      	400  	{object}  	models.ErrorResponse
// @Failure      	500  	{object}  	models.ErrorResponse
// @Router       	/delete-session [delete]
func (h *Handler) DeleteSession(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	defer body.Close()

	sessionJSON := new(models.SessionRequest)
	if err := json.NewDecoder(body).Decode(sessionJSON); err != nil {
		logrus.WithFields(
			logrus.Fields{
				"url": r.URL.String(),
			},
		).Errorf("decode json: %s", err)
		httpError(w, fmt.Errorf("decode json: %w", apiError.New(err.Error(), http.StatusBadRequest)))
		return
	}

	err := h.Services.SessionService.DeleteSession(sessionJSON.ID)
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
		Message: fmt.Sprintf("session successfully deleted"),
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

// @Summary      	Get List Sessions
// @Security Bearer
// @Description  	Получение из Redis: user_sessions <user_id>
// @Tags         	sessions
// @Accept			json
// @Produce			json
// @Param        user_id   path      int  true  "User ID"
// @Success		 	200		{object}	[]models.SessionResponseFull
// @Failure      	400  	{object}  	models.ErrorResponse
// @Failure      	500  	{object}  	models.ErrorResponse
// @Router       	/get-sessions/{user_id} [get]
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

// @Summary      	Delete Session
// @Security Bearer
// @Description  	Удаление сессий по user_id из Redis: user_sessions <user_id>.
// @Tags         	sessions
// @Accept			json
// @Produce			json
// @Param			user	body	models.SessionListRequest	true	"id"
// @Success		 	200		{object}	models.Response
// @Failure      	400  	{object}  	models.ErrorResponse
// @Failure      	500  	{object}  	models.ErrorResponse
// @Router       	/delete-sessions [delete]
func (h *Handler) DeleteListSession(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	defer body.Close()

	sessionJSON := new(models.SessionListRequest)
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
		Message: "sessions successfully deleted",
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
