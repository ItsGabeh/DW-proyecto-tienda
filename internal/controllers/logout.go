package controllers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
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

	session := sessions.Default(c)
	session.Clear()
	session.Save()

	// c.JSON(http.StatusOK, gin.H{"message": "Sesion cerrada exitosamente"})
	c.Header("HX-Redirect", "/")
	c.Status(http.StatusOK)
}
