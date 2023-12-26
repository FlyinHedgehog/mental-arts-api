package main

import (
	"fmt"
	"mentalartsapi_hw/database"
	"mentalartsapi_hw/handlers"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := database.InitDB()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	handler := handlers.NewHandler(db)

	router := gin.Default()

	handler.InitRoutes(router)

	router.Run(":8080")
}
