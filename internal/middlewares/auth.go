package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("CLAVE_SECRETA")

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extraer el token del header de autorizacion
		// authHeader := c.GetHeader("Authorization")
		// if authHeader == "" {
		// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Falta autenticacion"})
		// 	c.Abort()
		// 	return
		// }

		// Eliminar el prefijo Bearer
		// tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		// if tokenString == authHeader {
		// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "token de autenticacion no válido"})
		// 	c.Abort()
		// 	return
		// }

		// Extraer el token desde las cookies
		tokenString, err := c.Cookie("token")
		if err != nil {
			// c.JSON(http.StatusUnauthorized, gin.H{"error": "Falta token de autenticación"})
			c.HTML(http.StatusOK, "register.html", gin.H{"errors": []string{"Falta autenticacion"}})
			c.Abort()
			return
		}

		// Validar el token
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token de autenticacion no válido"})
			c.Abort()
			return
		}

		// Pasar los datos del usuario al contexto
		c.Set("email", claims.Email)
		c.Next()

	}
}
