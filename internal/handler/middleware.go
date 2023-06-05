package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func (h *Handler) tokenAuthMiddleware(c *gin.Context) {
	//cookie, err := c.Cookie("Authorization")
	//if err != nil {
	//	c.JSON(http.StatusUnauthorized, gin.H{
	//		"error": "not found any cookie",
	//	})
	//	return
	//}

	header := c.GetHeader("Authorization")

	if header == "" {
		//c.JSON(http.StatusUnauthorized, gin.H{
		//	"reason": "empty auth header",
		//})
		//c.Abort()
		//return

		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"reason": "empty auth header",
		})
		return
	}

	split := strings.Split(header, " ")
	if len(split) != 2 || split[0] != "Alif" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"reason": "invalid auth header",
		})
		return
	}

	if len(split[1]) == 0 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"reason": "token is empty",
		})
		return
	}

	userId, err := h.Auth.ParseToken(split[1])
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Set("userId", userId)
}

func getUserId(c *gin.Context) (int, error) {
	id, ok := c.Get("userId")
	if !ok {
		//c.JSON(http.StatusInternalServerError, gin.H{
		//	"error": "userId not found",
		//})
		return 0, errors.New("userId not found")
	}

	idInt, ok := id.(int)
	if !ok {
		//c.JSON(http.StatusInternalServerError, gin.H{
		//	"error": "invalid type of userId",
		//})
		return 0, errors.New("invalid type of userId")
	}

	return idInt, nil
}
