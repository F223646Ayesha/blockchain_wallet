package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"crypto-wallet/config"
	"crypto-wallet/routes"
	"crypto-wallet/services"
)

func main() {

	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ Warning: .env file not found")
	}

	// Init Firestore connection
	config.InitFirestore()

	// ----------------------------------------
	// 🔥 START MONTHLY ZAKAT SCHEDULER
	// ----------------------------------------
	services.InitScheduler()
	log.Println("⏰ Monthly Zakat Scheduler initialized")

	// Create Gin server
	r := gin.Default()

	// === CORS FIX ===
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
			"http://localhost:5174",
			"http://127.0.0.1:5173",
			"http://127.0.0.1:5174",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Register all API routes
	routes.RegisterRoutes(r)

	// Start backend server
	log.Println("🚀 Server running on http://localhost:8080")
	r.Run(":8080")
}
