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
)

func AddToCart(c *gin.Context) {
	email := c.GetString("email") // Obtener el email del contexto

	// Obtener los datos del producto a añadir
	var cartProduct models.CartProduct
	if err := c.ShouldBindJSON(&cartProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos invalidos"})
		return
	}

	// Validar los datos
	if err := validate.Struct(cartProduct); err != nil {
		var errorMessages []string
		validationErrors := err.(validator.ValidationErrors)
		for _, e := range validationErrors {
			errorMessages = append(errorMessages, e.Error())
		}
		c.JSON(http.StatusBadRequest, gin.H{"errors": errorMessages})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Buscar el usuario en la base de datos
	userCollection := db.Client.Database("tienda").Collection("users")
	var user models.User
	filter := bson.M{"email": email}
	err := userCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no encontrado"})
		return
	}

	// Obtener el carrito del usuario o crearlo si no existe
	cartCollection := db.Client.Database("tienda").Collection("carts")
	var cart models.Cart
	if err := cartCollection.FindOne(ctx, bson.M{"userId": user.ID}).Decode(&cart); err != nil {
		// Crear un nuevo carrito
		cart := models.Cart{
			ID:       primitive.NewObjectID(),
			UserID:   user.ID,
			Products: []models.CartProduct{cartProduct},
		}
		_, err := cartCollection.InsertOne(ctx, cart)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear carrito"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Producto añadido al carrito"})
		return
	}

	// Actualizar carrito existente
	updated := false
	for i, item := range cart.Products {
		if item.ProductID == cartProduct.ProductID {
			cart.Products[i].Quantity += cartProduct.Quantity
			updated = true
			break
		}
	}
	if !updated {
		cart.Products = append(cart.Products, cartProduct)
	}

	_, err = cartCollection.UpdateOne(ctx, bson.M{"_id": cart.ID}, bson.M{"$set": bson.M{"products": cart.Products}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el carrito"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Producto añadido al carrito"})
}

func GetCart(c *gin.Context) {
	email := c.GetString("email")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Buscar el usuario en la base de datos
	userCollection := db.Client.Database("tienda").Collection("users")
	var user models.User
	filter := bson.M{"email": email}
	err := userCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no encontrado"})
		return
	}

	//Obtener carrito
	cartCollection := db.Client.Database("tienda").Collection("carts")
	var cart models.Cart
	if err := cartCollection.FindOne(ctx, bson.M{"userId": user.ID}).Decode(&cart); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Carrito no encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"cart": cart})
}

func RemoveFromCart(c *gin.Context) {
	email := c.GetString("email") // Obe¿tener el email del usuario

	var productData struct {
		ID string `json:"productId" validate:"required"`
	}
	if err := c.ShouldBindJSON(&productData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos invalidos"})
		return
	}

	// Validar los datos
	if err := validate.Struct(productData); err != nil {
		var errorMessages []string
		validationErrors := err.(validator.ValidationErrors)
		for _, e := range validationErrors {
			errorMessages = append(errorMessages, e.Error())
		}
		c.JSON(http.StatusBadRequest, gin.H{"errors": errorMessages})
		return
	}

	// convertir el productID a ObjectID
	productID, err := primitive.ObjectIDFromHex(productData.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": "ID del producto no válido"})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Buscar el usuario en la base de datos
	userCollection := db.Client.Database("tienda").Collection("users")
	var user models.User
	filter := bson.M{"email": email}
	err = userCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no encontrado"})
		return
	}

	// Actualizar el carrito del usuario
	cartCollection := db.Client.Database("tienda").Collection("carts")
	filter = bson.M{"userId": user.ID}
	update := bson.M{"$pull": bson.M{"products": bson.M{"productId": productID}}}
	_, err = cartCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar producto del carrito"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Producto eliminado del carrito"})
}

func UpdateCartProduct(c *gin.Context) {
	email := c.GetString("email") // Obetener el email del usuario

	var productData struct {
		ID       string `json:"productId" validate:"required"`
		Quantity int    `json:"quantity" validate:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&productData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos invalidos"})
		return
	}

	// Validar los datos
	if err := validate.Struct(productData); err != nil {
		var errorMessages []string
		validationErrors := err.(validator.ValidationErrors)
		for _, e := range validationErrors {
			errorMessages = append(errorMessages, e.Error())
		}
		c.JSON(http.StatusBadRequest, gin.H{"errors": errorMessages})
		return
	}

	// convertir el productID a ObjectID
	productID, err := primitive.ObjectIDFromHex(productData.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": "ID del producto no válido"})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Buscar el usuario en la base de datos
	userCollection := db.Client.Database("tienda").Collection("users")
	var user models.User
	filter := bson.M{"email": email}
	err = userCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no encontrado"})
		return
	}

	// Actualizar el carrito del usuario
	cartCollection := db.Client.Database("tienda").Collection("carts")
	filter = bson.M{"userId": user.ID, "products.productId": productID}
	update := bson.M{"$set": bson.M{"products.$.quantity": productData.Quantity}}
	_, err = cartCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar producto del carrito"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Producto actualizado correctamente"})
}
