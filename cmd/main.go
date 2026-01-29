package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	_ "github.com/lib/pq"

	"github.com/djalben/epn-killer-mvp/internal/handlers"
	"github.com/djalben/epn-killer-mvp/internal/middleware"
	"github.com/djalben/epn-killer-mvp/internal/repository"
	"github.com/djalben/epn-killer-mvp/internal/core"
	"github.com/djalben/epn-killer-mvp/internal/telegram"
)

var DB *sql.DB

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	var err error
	DB, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}
	defer DB.Close()

	if err = DB.PingContext(context.Background()); err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	handlers.GlobalDB = DB
	repository.GlobalDB = DB

	log.Println("Successfully connected to the database!")

	if token := os.Getenv("TELEGRAM_BOT_TOKEN"); token != "" {
		telegram.SetBotToken(token)
		log.Println("Telegram bot token set: real notifications enabled")
	}

	go core.StartAutoReplenishmentWorker()

	router := mux.NewRouter()

	router.HandleFunc("/health", handlers.HealthCheckHandler).Methods("GET")
	router.HandleFunc("/api/v1/auth/register", handlers.RegisterHandler).Methods("POST")
	router.HandleFunc("/api/v1/auth/login", handlers.LoginHandler).Methods("POST")

	protectedRouter := router.PathPrefix("/api/v1/user").Subrouter()
	protectedRouter.Use(middleware.JWTAuthMiddleware)

	protectedRouter.HandleFunc("/me", handlers.GetMeHandler).Methods("GET")
	protectedRouter.HandleFunc("/grade", handlers.GetUserGradeHandler).Methods("GET")
	protectedRouter.HandleFunc("/deposit", handlers.ProcessDepositHandler).Methods("POST")
	protectedRouter.HandleFunc("/cards", handlers.GetUserCardsHandler).Methods("GET")
	protectedRouter.HandleFunc("/cards/issue", handlers.MassIssueCardsHandler).Methods("POST")
	protectedRouter.HandleFunc("/cards/{id}/status", handlers.PatchCardStatusHandler).Methods("PATCH")
	protectedRouter.HandleFunc("/cards/{id}/auto-replenishment", handlers.SetCardAutoReplenishmentHandler).Methods("POST")
	protectedRouter.HandleFunc("/cards/{id}/auto-replenishment", handlers.UnsetCardAutoReplenishmentHandler).Methods("DELETE")
	protectedRouter.HandleFunc("/report", handlers.GetUserTransactionReportHandler).Methods("GET")
	protectedRouter.HandleFunc("/api-key", handlers.CreateAPIKeyHandler).Methods("POST")

	protectedRouter.HandleFunc("/teams", handlers.GetUserTeamsHandler).Methods("GET")
	protectedRouter.HandleFunc("/teams", handlers.CreateTeamHandler).Methods("POST")
	protectedRouter.HandleFunc("/teams/{id}", handlers.GetTeamHandler).Methods("GET")
	protectedRouter.HandleFunc("/teams/{id}/members", handlers.InviteTeamMemberHandler).Methods("POST")
	protectedRouter.HandleFunc("/teams/{id}/members/{userId}", handlers.RemoveTeamMemberHandler).Methods("DELETE")
	protectedRouter.HandleFunc("/teams/{id}/members/{userId}/role", handlers.UpdateTeamMemberRoleHandler).Methods("PATCH")

	protectedRouter.HandleFunc("/referrals", handlers.GetReferralStatsHandler).Methods("GET")
	protectedRouter.HandleFunc("/settings/telegram", handlers.UpdateTelegramChatIDHandler).Methods("POST")

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	}).Handler(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server listening on port %s...", port)
	if err := http.ListenAndServe(":"+port, corsHandler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
