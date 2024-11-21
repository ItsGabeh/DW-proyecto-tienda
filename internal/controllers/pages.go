package controllers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func IndexPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func RegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", nil)
}

func GetNavbar(c *gin.Context) {
	// Extraer el token desde las cookies
	session := sessions.Default(c)
	email := session.Get("email")

	if email == nil {
		c.HTML(http.StatusOK, "guest_navbar.html", nil)
		return
	}
	// Aqui renderiza la navbar cuando si haya email
	c.HTML(http.StatusOK, "user_navbar.html", nil)
}
