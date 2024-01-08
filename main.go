package golangauthsample

import (
	"fmt"
	"runtime"
	"gorm.io/gorm"
	"github.com/gin-gonic/gin"
)
func main() {

	// Set the router as the default one shipped with Gin
	router := gin.Default()

	// Initialize routes
	initializeRoutes(router)

	// Start serving the application
	router.Run(":3000")
}