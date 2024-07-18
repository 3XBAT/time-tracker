package handlers

import (
	"errors"
	"github.com/3XBAT/time-tracker/internal/domain/models"
	"github.com/3XBAT/time-tracker/internal/storage"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// GetUsers godoc
// @Summary GetUsers
// @Tags User
// @Description Returns users according to filters and pagination
// @Accept json
// @Produce json
// @Param ID query string false "user id"
// @Param Name query string false "username"
// @Param Surname query string false "Surname"
// @Param Patronymic query string false "Patronymic"
// @Param PassportNumber query string false "PassportNumber"
// @Param Address query string false "Address"
// @Param Limit query int false "Limit"
// @Param Offset query int false "Offset"
// @Success 200 {object} map[string]int "{"id": 1}"
// @Failure 400 {object} map[string]string "{"error": "Bad Request"}"
// @Failure 404 {object} map[string]string "{"error": "Not Found"}"
// @Failure 500 {object} map[string]string "{"error": "Internal Server Error"}"
// @Router /users [get]
func (h *Handler) getUsers(c *gin.Context) {
	var queryParams models.QueryParams
	var err error

	queryParams.ID = c.Query("ID")
	queryParams.Name = c.Query("Name")
	queryParams.Surname = c.Query("Surname")
	queryParams.Patronymic = c.Query("Patronymic")
	queryParams.PassportNumber = c.Query("PassportNumber")
	queryParams.Address = c.Query("Address")

	if queryParams.Limit, err = strconv.Atoi(c.Query("Limit")); err != nil {
		if check := c.Query("Limit"); check == "" {
			queryParams.Limit = 10
		} else {
			newErrorResponse(c, http.StatusBadRequest, "invalid Limit")
			return
		}
	}
	if queryParams.Offset, err = strconv.Atoi(c.Query("Offset")); err != nil {
		if check := c.Query("Offset"); check == "" {
			queryParams.Offset = 0
		} else {
			newErrorResponse(c, http.StatusBadRequest, "invalid Offset")
			return
		}
	}

	users, err := h.service.UserProvider.Users(queryParams)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			newErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	c.JSON(http.StatusOK, users)
}

// GetUserByID godoc
// @Summary GetUserByID
// @Tags User
// @Description Returns a user by ID
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string "{"error": "Bad Request"}"
// @Failure 404 {object} map[string]string "{"error": "Not Found"}"
// @Failure 500 {object} map[string]string "{"error": "Internal Server Error"}"
// @Router /users/{id} [get]
func (h *Handler) getUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.service.UserProvider.UserById(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, user)
}

// CreateUser godoc
// @Summary CreateUser
// @Tags User
// @Description Create a new user
// @Accept json
// @Produce json
// @Param PassportNumber query string true "User info"
// @Success 200 {object} map[string]int "{"id": 1}"
// @Failure 400 {object} map[string]string "{"error": "Bad Request"}"
// @Failure 500 {object} map[string]string "{"error": "Internal Server Error"}"
// @Router /users [post]
func (h *Handler) createUser(c *gin.Context) {
	var userPassport string

	userPassport = c.Query("PassportNumber")

	id, err := h.service.UserProvider.Create(userPassport)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// UpdateUser godoc
// @Summary UpdateUser
// @Tags User
// @Description Update an existing user
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body models.UpdateUserInput true "User update info"
// @Success 200 {object} map[string]int "{"id": 1}"
// @Failure 400 {object} map[string]string "{"error": "Bad Request"}"
// @Failure 500 {object} map[string]string "{"error": "Internal Server Error"}"
// @Router /users/{id} [put]
func (h *Handler) updateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var user models.UpdateUserInput
	if err = c.BindJSON(&user); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err = h.service.UserProvider.Update(user, id); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"user_id": id,
	})
}

// DeleteUser godoc
// @Summary DeleteUser
// @Tags User
// @Description Delete a user by ID
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]int "{"id": 1}"
// @Failure 400 {object} map[string]string "{"error": "Bad Request"}"
// @Failure 500 {object} map[string]string "{"error": "Internal Server Error"}"
// @Router /users/{id} [delete]
func (h *Handler) deleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	err = h.service.UserProvider.Delete(id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})

}
