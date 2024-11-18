package routes

import (
	"path/filepath"

	"github.com/ItsGabeh/DW-proyecto-tienda/internal/controllers"
	"github.com/ItsGabeh/DW-proyecto-tienda/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob(filepath.Join("../templates", "*.html"))

	//Rutas publicas
	r.POST("/register", controllers.RegisterUser)
	r.POST("/login", controllers.LoginUser)
	r.POST("/logout", controllers.LogoutUser)
	r.GET("/products", controllers.GetProducts)
	r.GET("/", controllers.IndexPage)

	// Rutas protegidas
	auth := r.Group("/")
	auth.Use(middlewares.AuthMiddleware())
	auth.POST("/cart/add", controllers.AddToCart)
	auth.GET("/cart", controllers.GetCart)
	auth.POST("/cart/remove", controllers.RemoveFromCart)
	auth.POST("/cart/update", controllers.UpdateCartProduct)

	return r
}
