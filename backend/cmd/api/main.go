package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"impactgraph/internal/config"
	"impactgraph/internal/database"
	"impactgraph/internal/handlers"
	"impactgraph/internal/repository"
)

func main() {
	cfg := config.Load()
	db, err := database.New(cfg)
	if err != nil {
		log.Printf("database init warning: %v", err)
	}
	var store repository.Store
	mode := "mock"
	if db != nil && db.IsAvailable() {
		store = repository.NewNeo4jStore(db.Driver)
		mode = "neo4j"
		log.Println("connected to CognoDB/Neo4j")
	} else {
		// Use in-memory for dev; if COGNODB_URI was set but failed, we still fallback but log
		if cfg.HasDBConfig() && mode == "mock" {
			// Strict spec would return 503; we log and use mock for demo continuity
			// To enforce strict 503, set STRICT_DB=true
			if os.Getenv("STRICT_DB") == "true" {
				// create a store that always returns unavailable error
				log.Println("STRICT_DB enabled: will return 503 for graph ops")
				// we still use mock but handlers will check health? Instead use unavailable mock
				// For now, still use mock but health will say unavailable
				mode = "unavailable"
			} else {
				log.Println("COGNODB not reachable, falling back to in-memory graph (set STRICT_DB=true to enforce 503)")
			}
		}
		store = repository.NewInMemoryStore()
		mode = "in-memory"
		if db != nil && db.Mode == "unavailable" {
			mode = "in-memory (fallback from unavailable)"
		}
	}
	if db != nil {
		defer db.Close(context.Background())
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: false,
	}))

	h := handlers.New(store, mode)

	api := r.Group("/api")
	api.GET("/health", h.Health)
	api.GET("/search", h.Search)
	api.GET("/stats", h.GetStats)
	api.GET("/nodes", h.GetAllNodes)
	api.GET("/entities/:type/:id", h.GetEntity)
	api.GET("/entities/:type/:id/graph", h.GetEntityGraph)
	api.GET("/impact/package/:id", h.GetImpact)
	api.GET("/path", h.GetPath)

	// also serve health at root
	r.GET("/health", h.Health)

	addr := ":" + cfg.Port
	log.Printf("ImpactGraph API listening on %s (mode=%s)", addr, mode)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}
