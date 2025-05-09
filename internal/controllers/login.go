package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/ItsGabeh/DW-proyecto-tienda/internal/db"
	"github.com/ItsGabeh/DW-proyecto-tienda/internal/models"
	"github.com/ItsGabeh/DW-proyecto-tienda/internal/utils"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("CLAVE_SECRETA")

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func LoginUser(c *gin.Context) {
	var loginData struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	// vincular el JSON con los datos para el login
	// if err := c.ShouldBindJSON(&loginData); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Datos invalidos"})
	// 	return
	// }

	// Extraer los elementos que vengan desde un formulario
	loginData.Email = c.PostForm("email")
	loginData.Password = c.PostForm("password")

	// Validar datos
	if err := validate.Struct(loginData); err != nil {
		errorMessages := utils.ValidationMessages(err)
		// c.JSON(http.StatusBadRequest, gin.H{"errors": errorMessages})
		c.HTML(http.StatusOK, "login.html", gin.H{"errors": errorMessages})
		return
	}

	// Buscar el usuario en la base de datos
	userCollection := db.Client.Database("tienda").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	filter := bson.M{"email": loginData.Email}
	err := userCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales incorrectas"})
		return
	}

	// Comparar la contraseña del login y del usuario encontrado
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Credenciales incorrectas"})
		return
	}

	// Crear el token JWT
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "Aplicacion de tienda",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al generar el token"})
		return
	}

	// configurar la cookie con el token
	c.SetCookie(
		"token",
		tokenString,
		86400,
		"/",
		"",
		false,
		true,
	)

	session := sessions.Default(c)
	session.Set("email", user.Email)
	session.Save()

	// c.JSON(http.StatusOK, gin.H{"message": "Inicio de sesión exitoso"})
	// c.HTML(http.StatusOK, "index.html", gin.H{"message": "Inicio de sesión exitoso"})
	c.Header("HX-Redirect", "/")
	c.Status(http.StatusFound)
}
