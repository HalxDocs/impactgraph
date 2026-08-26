package database

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"impactgraph/internal/config"
)

type DB struct {
	Driver neo4j.DriverWithContext
	Available bool
	Mode string // "neo4j" or "mock"
}

func New(cfg *config.Config) (*DB, error) {
	if !cfg.HasDBConfig() {
		return &DB{Available: false, Mode: "mock"}, nil
	}
	driver, err := neo4j.NewDriverWithContext(cfg.CognoDBURI, neo4j.BasicAuth(cfg.CognoDBUsername, cfg.CognoDBPassword, ""))
	if err != nil {
		return &DB{Available: false, Mode: "mock"}, fmt.Errorf("neo4j driver: %w", err)
	}
	ctx := context.Background()
	if err := driver.VerifyConnectivity(ctx); err != nil {
		// Do not hard-fail; caller will treat as unavailable and use mock or return 503 if strict.
		_ = driver.Close(ctx)
		return &DB{Available: false, Mode: "unavailable"}, fmt.Errorf("verify connectivity: %w", err)
	}
	return &DB{Driver: driver, Available: true, Mode: "neo4j"}, nil
}

func (d *DB) Close(ctx context.Context) {
	if d.Driver != nil {
		_ = d.Driver.Close(ctx)
	}
}

func (d *DB) IsAvailable() bool { return d.Available && d.Driver != nil }
