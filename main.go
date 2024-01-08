package main

import (
	"github.com/gin-gonic/gin"
)
func main() {
	router := gin.Default()

	initializeRoutes(router)

	if err := router.Run(":5500"); err != nil {
		panic(err)
	}
}