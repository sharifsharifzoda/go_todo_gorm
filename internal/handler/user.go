package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"todo_gorm/model"
)

func (h *Handler) signUp(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid JSON provided",
		})
		return
	}

	err := h.Auth.ValidateUser(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "validate",
		})
		return
	}

	isUsed := h.Auth.IsEmailUsed(user.Email)
	if isUsed {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "email is already created",
		})
		return
	}

	id, err := h.Auth.CreateUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create user",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id": id,
	})
}

func (h *Handler) signIn(c *gin.Context) {
	var user model.User
	if c.ShouldBindJSON(&user) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid JSON provided",
		})
		return
	}

	if err := h.Auth.ValidateUser(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "validation",
		})
		return
	}

	checkedUser, err := h.Auth.CheckUser(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	token, err := h.Auth.GenerateToken(checkedUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	//c.SetSameSite(http.SameSiteLaxMode)
	//c.SetCookie("Authorization", token, 3600*24*30, "", "", false, true)

	c.Writer.Header().Set("Authorization", token)

	c.JSON(http.StatusOK, gin.H{
		"msg": "signed in",
	})
}
