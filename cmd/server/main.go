package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/stkrfx/dg-backend/internal/config"
	"github.com/stkrfx/dg-backend/internal/handler"
	"github.com/stkrfx/dg-backend/internal/repository"
	"github.com/stkrfx/dg-backend/internal/service"
)

// Response is a helper for the health check
type Response struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func main() {
	// 1. Load Config (fails if DATABASE_URL is missing)
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	// 2. Initialize Database (Connects to Supabase)
	db, err := repository.NewDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("DB connection error: %v", err)
	}
	defer db.Close()

	// 3. Dependency Injection (Wiring the layers)
	// We pass 'db' into the repository now
	repo := repository.NewClientRepository(db)
	svc := service.NewClientService(repo)
	hdl := handler.NewClientHandler(svc)

	// 4. Set up Router
	router := http.NewServeMux()

	// Register Routes
	router.HandleFunc("GET /health", healthHandler)
	router.HandleFunc("POST /clients", hdl.CreateClient)
	router.HandleFunc("GET /clients", hdl.GetClients)

	// 5. Configure the Server (Professional settings)
	// We use cfg.Port to allow dynamic port assignment (e.g., by Heroku/Render/Docker)
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	fmt.Printf("Server starting on port %s...\n", cfg.Port)

	// 6. Start Server
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %s", err)
	}
}

// healthHandler provides a quick check for monitoring tools
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(Response{
		Message: "Server is running smoothly",
		Status:  200,
	})
}
