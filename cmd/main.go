package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/cors"

	api "github.com/djalben/epn-killer-mvp/internal/api"
	"github.com/djalben/epn-killer-mvp/internal/middleware"
	"github.com/djalben/epn-killer-mvp/internal/repository"
	"github.com/djalben/epn-killer-mvp/internal/telegram"
	"github.com/djalben/epn-killer-mvp/internal/usecases"
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

	api.GlobalDB = DB
	repository.GlobalDB = DB

	log.Println("Successfully connected to the database!")

	if token := os.Getenv("TELEGRAM_BOT_TOKEN"); token != "" {
		telegram.SetBotToken(token)
		log.Println("Telegram bot token set: real notifications enabled")
	}

	go usecases.StartAutoReplenishmentWorker()

	router := mux.NewRouter()

	router.HandleFunc("/health", api.HealthCheckHandler).Methods("GET")
	router.HandleFunc("/api/v1/auth/register", api.RegisterHandler).Methods("POST")
	router.HandleFunc("/api/v1/auth/login", api.LoginHandler).Methods("POST")

	protectedRouter := router.PathPrefix("/api/v1/user").Subrouter()
	protectedRouter.Use(middleware.JWTAuthMiddleware)

	protectedRouter.HandleFunc("/me", api.GetMeHandler).Methods("GET")
	protectedRouter.HandleFunc("/grade", api.GetUserGradeHandler).Methods("GET")
	protectedRouter.HandleFunc("/deposit", api.ProcessDepositHandler).Methods("POST")
	protectedRouter.HandleFunc("/cards", api.GetUserCardsHandler).Methods("GET")
	protectedRouter.HandleFunc("/cards/issue", api.MassIssueCardsHandler).Methods("POST")
	protectedRouter.HandleFunc("/cards/{id}/status", api.PatchCardStatusHandler).Methods("PATCH")
	protectedRouter.HandleFunc("/cards/{id}/auto-replenishment", api.SetCardAutoReplenishmentHandler).Methods("POST")
	protectedRouter.HandleFunc("/cards/{id}/auto-replenishment", api.UnsetCardAutoReplenishmentHandler).Methods("DELETE")
	protectedRouter.HandleFunc("/report", api.GetUserTransactionReportHandler).Methods("GET")
	protectedRouter.HandleFunc("/api-key", api.CreateAPIKeyHandler).Methods("POST")

	protectedRouter.HandleFunc("/teams", api.GetUserTeamsHandler).Methods("GET")
	protectedRouter.HandleFunc("/teams", api.CreateTeamHandler).Methods("POST")
	protectedRouter.HandleFunc("/teams/{id}", api.GetTeamHandler).Methods("GET")
	protectedRouter.HandleFunc("/teams/{id}/members", api.InviteTeamMemberHandler).Methods("POST")
	protectedRouter.HandleFunc("/teams/{id}/members/{userId}", api.RemoveTeamMemberHandler).Methods("DELETE")
	protectedRouter.HandleFunc("/teams/{id}/members/{userId}/role", api.UpdateTeamMemberRoleHandler).Methods("PATCH")

	protectedRouter.HandleFunc("/referrals", api.GetReferralStatsHandler).Methods("GET")
	protectedRouter.HandleFunc("/settings/telegram", api.UpdateTelegramChatIDHandler).Methods("POST")

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
