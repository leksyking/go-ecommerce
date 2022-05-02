package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/leksyking/go-ecommerce/controllers"
	"github.com/leksyking/go-ecommerce/database"
	"github.com/leksyking/go-ecommerce/middlewares"
	"github.com/leksyking/go-ecommerce/routes"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	app := controllers.NewApplication(database.ProductData(database.Client, "Products"), database.UserData(database.Client, "Users"))
	router := gin.New()
	router.Use(gin.Logger())

	//user routes
	routes.UserRoutes(router)

	router.Use(middlewares.Authentication())

	//cart routes
	router.GET("/addtocart", app.AddToCart())
	router.GET("/removeitem", app.RemoveItem())
	router.GET("/cartcheckout", app.BuyFromCart())
	router.GET("/instantbuy", app.InstantBuy())

	log.Fatal(router.Run(":" + port))
}
