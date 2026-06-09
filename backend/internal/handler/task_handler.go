package handler

import (
	"net/http"
	"strconv"

	"taskmaster/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TaskHandler struct {
	service *service.TaskService
}

func NewTaskHandler(service *service.TaskService) *TaskHandler {
	return &TaskHandler{service}
}

func (h *TaskHandler) Create(c *gin.Context) {
	var input service.CreateTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("user_id").(uuid.UUID)
	teamIDStr := c.MustGet("team_id").(string)
	var teamID uuid.UUID
	if teamIDStr != "" {
		teamID = uuid.MustParse(teamIDStr)
	}

	task, err := h.service.Create(c.Request.Context(), input, userID, teamID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"task": task})
}

func (h *TaskHandler) List(c *gin.Context) {
	teamIDStr := c.MustGet("team_id").(string)
	teamID := uuid.MustParse(teamIDStr)

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	tasks, err := h.service.List(c.Request.Context(), teamID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

func (h *TaskHandler) Get(c *gin.Context) {
	id := uuid.MustParse(c.Param("id"))
	task, err := h.service.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"task": task})
}

func (h *TaskHandler) Update(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "update not implemented in starter"})
}

func (h *TaskHandler) Delete(c *gin.Context) {
	id := uuid.MustParse(c.Param("id"))
	teamIDStr := c.MustGet("team_id").(string)
	teamID := uuid.MustParse(teamIDStr)

	if err := h.service.Delete(c.Request.Context(), id, teamID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *TaskHandler) Assign(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "assign not implemented in starter"})
}
