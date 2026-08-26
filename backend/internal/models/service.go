package models

// Service model — see graph.go
type ServiceProps struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Environment string `json:"environment"`
	Description string `json:"description"`
}
