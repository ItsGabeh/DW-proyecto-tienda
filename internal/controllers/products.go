package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/ItsGabeh/DW-proyecto-tienda/internal/db"
	"github.com/ItsGabeh/DW-proyecto-tienda/internal/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func GetProducts(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	productsCollection := db.Client.Database("tienda").Collection("products")
	filter := bson.M{}
	cursor, err := productsCollection.Find(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener los productos"})
		return
	}
	defer cursor.Close(ctx)

	var products []models.Product
	if err := cursor.All(ctx, &products); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al procesaro los productos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"products": products})
}
