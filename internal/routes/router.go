package routes

import (
	"path/filepath"

	"github.com/ItsGabeh/DW-proyecto-tienda/internal/controllers"
	"github.com/ItsGabeh/DW-proyecto-tienda/internal/middlewares"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Static("/static", "../static")
	r.LoadHTMLGlob(filepath.Join("../templates", "*.html"))
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("session", store))

	//Rutas publicas
	r.POST("/register", controllers.RegisterUser)
	r.POST("/login", controllers.LoginUser)
	r.POST("/logout", controllers.LogoutUser)
	r.GET("/products", controllers.GetProducts)
	r.GET("/products/:id", controllers.GetProduct)
	r.GET("/", controllers.IndexPage)
	r.GET("/navbar", controllers.GetNavbar)
	r.GET("/login", controllers.LoginPage)
	r.GET("/register", controllers.RegisterPage)

	// Rutas protegidas
	auth := r.Group("/")
	auth.Use(middlewares.AuthMiddleware())
	auth.POST("/cart/add", controllers.AddToCart)
	auth.GET("/cart", controllers.GetCart)
	auth.POST("/cart/remove", controllers.RemoveFromCart)
	auth.POST("/cart/update", controllers.UpdateCartProduct)

	return r
}
