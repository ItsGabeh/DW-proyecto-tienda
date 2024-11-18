package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/ItsGabeh/DW-proyecto-tienda/internal/db"
	"github.com/ItsGabeh/DW-proyecto-tienda/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

func RegisterUser(c *gin.Context) {
	var user models.User

	// vincular el JSON recibido con la estructura del usuario
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos invalidos"})
		return
	}

	// Validar los datos del usuario
	if err := validate.Struct(user); err != nil {
		var errorMessages []string
		validationErrors := err.(validator.ValidationErrors)
		for _, e := range validationErrors {
			errorMessages = append(errorMessages, e.Error())
		}
		c.JSON(http.StatusBadRequest, gin.H{"errors": errorMessages})
		return
	}

	// Verificar si el email ya existe
	userCollection := db.Client.Database("tienda").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"email": user.Email}
	count, err := userCollection.CountDocuments(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al verificar el email"})
		return
	}
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El email ya está registrado"})
		return
	}

	// Encriptar la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al encriptar la contraseña"})
		return
	}
	user.Password = string(hashedPassword)

	// Insertar el usuario en la base de datos
	user.ID = primitive.NewObjectID()
	_, err = userCollection.InsertOne(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al registrar el usuario"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Usuario correctamente registrado"})

}
