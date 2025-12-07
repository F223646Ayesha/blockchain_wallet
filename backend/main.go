package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"crypto-wallet/config"
	"crypto-wallet/routes"
	"crypto-wallet/services"
)

func main() {

	// Load .env only for local development
	_ = godotenv.Load()

	// Initialize Firestore connection
	config.InitFirestore()

	// Start monthly Zakat scheduler
	services.InitScheduler()
	log.Println("⏰ Monthly Zakat Scheduler initialized")

	// Create Gin router
	r := gin.Default()

	// ======================================================
	// ✅ CORS CONFIG — ALLOW LOCAL + RENDER FRONTEND
	// ======================================================
	frontendURL := os.Getenv("FRONTEND_URL") // Set this in Render
	if frontendURL == "" {
		frontendURL = "http://localhost:5173" // fallback for local dev
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			frontendURL,
			"http://localhost:5173",
			"http://localhost:5174",
			"http://127.0.0.1:5173",
			"http://127.0.0.1:5174",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Register all API routes
	routes.RegisterRoutes(r)

	// ======================================================
	// ✅ PORT FIX — Required for Render Deployment
	// ======================================================
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // local fallback
	}

	log.Println("🚀 Server running on port:", port)

	// Start server
	if err := r.Run(":" + port); err != nil {
		log.Fatal("❌ Failed to start server: ", err)
	}
}
