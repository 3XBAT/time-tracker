package handlers

import (
	"errors"
	"fmt"
	"time"

	"github.com/3XBAT/time-tracker/internal/domain/models"
	"github.com/3XBAT/time-tracker/internal/storage"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// CreateTask godoc
// @Summary CreateTask
// @Tags Task
// @Description Create task. Example start_time: 2024-07-15T13:35:35.481207+03:00
// @Accept json
// @Produce json
// @Param input body models.InputTaskCreate true "task info"
// @Success 200 {object} map[string]int "{"id": 1}"
// @Failure 400 {object} map[string]string "{"error": "Bad Request"}"
// @Failure 500 {object} map[string]string "{"error": "Internal Server Error"}"
// @Router /tasks/ [post]
func (h *Handler) createTask(c *gin.Context) {
	var input models.InputTaskCreate

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	taskId, err := h.service.TaskProvider.Create(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": taskId,
	})
}

// @Summary UpdateTask
// @Tags Task
// @Description Update an existing task
// @Accept json
// @Produce json
// @Param input body models.InputTaskUpdate true "task update info"
// @Success 200 {object} map[string]int "{"task id": 1}"
// @Failure 400 {object} map[string]string "{"error": "Bad Request"}"
// @Failure 404 {object} map[string]string "{"error": "Not Found"}"
// @Failure 208 {object} map[string]string "{"error": "Task Already Ended"}"
// @Failure 500 {object} map[string]string "{"error": "Internal Server Error"}"
// @Router /tasks/{id} [put]
func (h *Handler) updateTask(c *gin.Context) {
	var input models.InputTaskUpdate

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if _, err := h.service.UserProvider.UserById(input.UserID); err != nil {
		newErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	if err := h.service.TaskProvider.Update(input); err != nil {
		if errors.Is(err, storage.ErrTaskEnded) {
			newErrorResponse(c, http.StatusAlreadyReported, err.Error())
			return
		}
		if errors.Is(err, storage.ErrTaskNotFound) {
			newErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"task id": input.Id,
	})

}

// @Summary DeleteTask
// @Tags Task
// @Description Delete an existing task
// @Accept json
// @Produce json
// @Param input body models.InputTaskDelete true "task delete info"
// @Success 200 {object} map[string]int "{"id": 1}"
// @Failure 400 {object} map[string]string "{"error": "Bad Request"}"
// @Failure 404 {object} map[string]string "{"error": "Not Found"}"
// @Failure 500 {object} map[string]string "{"error": "Internal Server Error"}"
// @Router /tasks/{id} [delete]
func (h *Handler) deleteTask(c *gin.Context) {

	var input models.InputTaskDelete
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "error with body"+err.Error())
		return
	}

	if _, err := h.service.UserProvider.UserById(input.UserID); err != nil {
		newErrorResponse(c, http.StatusNotFound, err.Error())
	}

	err := h.service.TaskProvider.Delete(input)
	if err != nil {
		if errors.Is(err, storage.ErrTaskNotFound) {
			newErrorResponse(c, http.StatusNotFound, err.Error())
			return
		}

		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": input.TaskID,
	})
}

// @Summary GetTasks
// @Tags Task
// @Description Get tasks for a user within a specified period. EXAMPLE: 2024-07-15T13:35:35.481207+03:00
// @Accept json
// @Produce json
// @Param user_id query int true "User ID"
// @Param start_time query string false "Start Time" format(date-time)
// @Param end_time query string false "End Time" format(date-time)
// @Success 200 {object} map[string][]models.OutputTask "{"tasks": [...]}"
// @Failure 400 {object} map[string]string "{"error": "Bad Request"}"
// @Failure 500 {object} map[string]string "{"error": "Internal Server Error"}"
// @Router /tasks [get]
func (h *Handler) getTasks(c *gin.Context) {

	var input models.InputTask

	userID, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid user ID"+err.Error())
		return
	}
	input.UserID = userID

	checkStartTime := c.Query("start_time")

	if checkStartTime != "" {
		startTime, err := time.Parse(time.RFC3339, checkStartTime)
		if err != nil {
			newErrorResponse(c, http.StatusBadRequest, "invalid start time"+err.Error())
			return
		}
		input.StartPeriod = &startTime
	} else {
		input.StartPeriod = nil
	}

	checkEndTime := c.Query("end_time")
	if checkEndTime != "" {
		endTime, err := time.Parse(time.RFC3339, c.Query("end_time"))
		if err != nil {
			newErrorResponse(c, http.StatusBadRequest, "invalid end time"+err.Error())
			return
		}
		input.EndPeriod = &endTime
	} else {
		input.EndPeriod = nil
	}

	fmt.Println(userID, checkEndTime, checkStartTime)
	
	tasks, err := h.service.TaskProvider.Tasks(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, tasks)
}
