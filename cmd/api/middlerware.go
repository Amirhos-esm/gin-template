package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware handles CORS settings for incoming requests
func (app * application) CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set CORS headers
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins
                  // Allow credentials

		// Handle preflight OPTIONS request
		if c.Request.Method == http.MethodOptions {
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS") // Allowed HTTP methods
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")    // Allowed headers
			c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")                // Exposed headers
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")   
			c.AbortWithStatus(http.StatusNoContent) // Return 204 for preflight 
			return
		}

		c.Next() // Proceed to the next middleware/handler
	}
}

func (app *application) AuthRequiredMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, _, err := app.auth.GetTokenFromHeaderAndVerify(c.Writer,c.Request)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// Continue to the next middleware or handler
		c.Next()
	}
}