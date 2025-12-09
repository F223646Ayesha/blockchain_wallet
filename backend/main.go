package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"crypto-wallet/config"
	"crypto-wallet/routes"
	"crypto-wallet/services"
)

func main() {
	godotenv.Load()

	config.InitFirestore()
	services.InitScheduler()
	log.Println("⏰ Monthly Zakat Scheduler initialized")

	r := gin.Default()

	// ---------------------------
	// FIX: CORS FIRST
	// ---------------------------
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
			"http://127.0.0.1:5173",
			"https://blockchain-wallet-ui.vercel.app",
			"https://blockchain-wallet-z5s8.onrender.com",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	log.Println("🚀 CORS middleware initialized")

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API Routes
	routes.RegisterRoutes(r)

	// PORT
	port := os.Getenv("PORT")
	if port == "" {
		port = "10000"
	}

	log.Println("🚀 Server running on port:", port)

	r.Run(":" + port)
}
