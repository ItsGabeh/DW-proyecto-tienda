package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func LogoutUser(c *gin.Context) {
	// Eliminar el token de la cookie
	c.SetCookie(
		"token",
		"",
		-1,
		"/",
		"",
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{"message": "Sesion cerrada exitosamente"})
}
