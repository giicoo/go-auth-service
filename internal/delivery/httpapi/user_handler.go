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

// @Summary      	Create User
// @Description  	Проверка отсутствия аккаунта с этим email. Пароль хешируется с помощью `bcrypt`. *User* записывается записывается в БД
// @Tags         	users
// @Security Bearer
// @Accept			json
// @Produce			json
// @Param			user	body	models.UserCreateRequest	true	"Write Email and Password"
// @Success		 	200		{object}	models.UserResponse
// @Failure      	400  	{object}  	models.ErrorResponse
// @Failure      	500  	{object}  	models.ErrorResponse
// @Router       	/create-user [post]
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
		UserID: userDB.ID,
		Email:  user.Email,
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

// @Summary      	Delete User
// @Security Bearer
// @Description  	Удаление из БД по ID из Session
// @Tags         	users
// @Accept			json
// @Produce			json
// @Param			user	body	models.UserIdRequest	true	"id"
// @Success		 	200		{object}	models.Response
// @Failure      	400  	{object}  	models.ErrorResponse
// @Failure      	500  	{object}  	models.ErrorResponse
// @Router       	/delete-user [delete]
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	defer body.Close()

	requestJSON := new(models.UserIdRequest)
	if err := json.NewDecoder(body).Decode(requestJSON); err != nil {
		logrus.WithFields(
			logrus.Fields{
				"url": r.URL.String(),
			},
		).Errorf("decode json: %s", err)
		httpError(w, fmt.Errorf("decode json: %w", apiError.New(err.Error(), http.StatusBadRequest)))
		return
	}

	if err := h.Services.UserService.DeleteUser(requestJSON.UserID); err != nil {
		logrus.WithFields(
			logrus.Fields{
				"url":     r.URL.String(),
				"request": requestJSON,
			},
		).Errorf("service delete user: %s", err)
		httpError(w, fmt.Errorf("delete user: %w", err))
		return
	}

	response := models.Response{
		Message: "user successfully deleted",
	}
	if err := httpResponse(w, response); err != nil {
		logrus.WithFields(
			logrus.Fields{
				"url":     r.URL.String(),
				"request": requestJSON,
			},
		).Errorf("send response: %s", err)
		httpError(w, fmt.Errorf("send response: %w", err))
		return
	}
}

// @Summary      	Update Email
// @Security Bearer
// @Description  	Проверка не занят ли новый email. Обновление полей в БД по ID из Session
// @Tags         	users
// @Accept			json
// @Produce			json
// @Param			user	body	models.UserUpdateEmailRequest	true	"id and new email"
// @Success		 	200		{object}	models.UserResponse
// @Failure      	400  	{object}  	models.ErrorResponse
// @Failure      	500  	{object}  	models.ErrorResponse
// @Router       	/update-email [put]
func (h *Handler) UpdateEmailUser(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	defer body.Close()

	userJSON := new(models.UserUpdateEmailRequest)
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
		ID:    userJSON.UserID,
		Email: userJSON.Email,
	}

	userDB, err := h.Services.UserService.UpdateEmailUser(user)
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
		UserID: userDB.ID,
		Email:  userDB.Email,
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

// @Summary      	Update Password
// @Security Bearer
// @Description  	Хеширование нового пароля с помощью bcrypt. Обновление полей в БД по ID из Session
// @Tags         	users
// @Accept			json
// @Produce			json
// @Param			user	body	models.UserUpdatePasswordRequest	true	"id and new password"
// @Success		 	200		{object}	models.UserResponse
// @Failure      	400  	{object}  	models.ErrorResponse
// @Failure      	500  	{object}  	models.ErrorResponse
// @Router       	/update-password [put]
func (h *Handler) UpdatePasswordUser(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	defer body.Close()

	userJSON := new(models.UserUpdatePasswordRequest)
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
		ID:       userJSON.UserID,
		Password: userJSON.Password,
	}

	userDB, err := h.Services.UserService.UpdatePasswordUser(user)
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
		UserID: userDB.ID,
		Email:  userDB.Email,
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

// @Summary      	Get User By ID
// @Security Bearer
// @Description  	Получение user из БД
// @Tags         	users
// @Accept			json
// @Produce			json
// @Param        id   path      int  true  "User ID"
// @Success		 	200		{object}	models.UserResponse
// @Failure      	400  	{object}  	models.ErrorResponse
// @Failure      	500  	{object}  	models.ErrorResponse
// @Router       	/get-user-by-id/{id} [get]
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

	userDB, err := h.Services.UserService.GetUserByID(id)
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
		UserID: userDB.ID,
		Email:  userDB.Email,
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

// @Summary      	Create User
// @Security Bearer
// @Description  	Проверка пароля
// @Tags         	users
// @Accept			json
// @Produce			json
// @Param			user	body	models.UserCheckRequest	true	"Write Email and Password"
// @Success		 	200		{object}	models.UserResponse
// @Failure      	400  	{object}  	models.ErrorResponse
// @Failure      	500  	{object}  	models.ErrorResponse
// @Router       	/check-user [post]
func (h *Handler) CheckUser(w http.ResponseWriter, r *http.Request) {
	body := r.Body
	defer body.Close()

	userJSON := new(models.UserCheckRequest)
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

	userDB, err := h.Services.UserService.CheckUser(user)
	if err != nil {
		logrus.WithFields(
			logrus.Fields{
				"url":     r.URL.String(),
				"request": userJSON,
			},
		).Errorf("service check user: %s", err)
		httpError(w, fmt.Errorf("check user: %w", err))
		return
	}

	userResponse := models.UserResponse{
		UserID: userDB.ID,
		Email:  userDB.Email,
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
