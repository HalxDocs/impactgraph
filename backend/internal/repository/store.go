package repository

import (
	"context"
	"impactgraph/internal/models"
)

// Store is the abstraction over Neo4j vs in-memory.
type Store interface {
	Health(ctx context.Context) error
	Search(ctx context.Context, q string) ([]models.SearchResult, error)
	GetEntity(ctx context.Context, typ, id string) (*models.EntityDetail, error)
	GetEntityGraph(ctx context.Context, typ, id string) (*models.Graph, error)
	GetImpact(ctx context.Context, packageID string) (*models.ImpactResponse, error)
	GetPath(ctx context.Context, from, to string) (*models.PathResponse, error)
	GetStats(ctx context.Context) (*models.StatsResponse, error)
	GetAllNodes(ctx context.Context) ([]models.Node, error)
}

// Entity type normalization
func normalizeType(t string) string {
	switch t {
	case "package", "Package", "PACKAGES":
		return "Package"
	case "repository", "Repository", "repo":
		return "Repository"
	case "service", "Service":
		return "Service"
	case "developer", "Developer":
		return "Developer"
	case "team", "Team":
		return "Team"
	case "organization", "Organization", "org":
		return "Organization"
	default:
		return t
	}
}
