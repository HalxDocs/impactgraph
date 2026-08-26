package models

// Repository model — see graph.go
type RepositoryProps struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	URL        string `json:"url"`
	Language   string `json:"language"`
	Visibility string `json:"visibility"`
}
