package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"todo_gorm/model"
)

func (h *Handler) createTask(c *gin.Context) {
	var task *model.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid JSON provided",
		})
		return
	}

	if task.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid task format provided",
		})
		return
	}

	userId, err := getUserId(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	id, err := h.Todo.CreateTask(userId, task)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to create task",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

func (h *Handler) getTasks(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	tasks, err := h.Todo.GetAll(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to get the list of tasks",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tasks": tasks,
	})
}

func (h *Handler) getTaskById(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id provided",
		})
		return
	}

	task, err := h.Todo.GetTaskById(userId, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"task": task,
	})
}

func (h *Handler) updateTask(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id provided",
		})
		return
	}

	var task *model.Task
	if err = c.BindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid JSON provided",
		})
		return
	}

	if err := h.Todo.ValidateTask(task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := h.Todo.UpdateTask(userId, id, task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "successfully updated",
	})
}

func (h *Handler) deleteTask(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id provided",
		})
		return
	}

	if err := h.Todo.DeleteTask(userId, id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "successfully deleted",
	})
}
