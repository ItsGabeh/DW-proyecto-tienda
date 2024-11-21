package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/ItsGabeh/DW-proyecto-tienda/internal/db"
	"github.com/ItsGabeh/DW-proyecto-tienda/internal/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetProducts(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Obtener productos de la base de datos
	productsCollection := db.Client.Database("tienda").Collection("products")
	filter := bson.M{}
	cursor, err := productsCollection.Find(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudieron obtener los productos"})
		return
	}
	defer cursor.Close(ctx)

	// Poner los productos en una lista
	var products []models.Product
	if err := cursor.All(ctx, &products); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al procesaro los productos"})
		return
	}

	// Poner los productos en un formato para el front
	var productsData []gin.H
	for _, item := range products {
		productsData = append(productsData, gin.H{
			"ID":          item.ID.Hex(),
			"Name":        item.Name,
			"Description": item.Description,
			"Price":       item.Price,
			"Stock":       item.Stock,
		})
	}

	// Respuesta formato json
	// c.JSON(http.StatusOK, gin.H{"products": products})
	// Respuesta en formato html
	c.HTML(http.StatusOK, "products.html", gin.H{"Products": productsData})
}

func GetProduct(c *gin.Context) {
	productId := c.Param("id")
	id, err := primitive.ObjectIDFromHex(productId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El id no es válido"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Obtener el producto desde la base de datos
	productCollection := db.Client.Database("tienda").Collection("products")
	var product models.Product
	err = productCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&product)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No se pudo obtener los datos del producto"})
		return
	}

	// Poner el producto en un formato para el front
	productData := gin.H{
		"ID":          product.ID.Hex(),
		"Name":        product.Name,
		"Description": product.Description,
		"Price":       product.Price,
		"Stock":       product.Stock,
	}

	// Regresar el producto en formato json
	// c.JSON(http.StatusOK, gin.H{"product": product})
	c.HTML(http.StatusOK, "product.html", gin.H{"product": productData})
}
