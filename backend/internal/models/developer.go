package models

// Developer model — see graph.go for shared Node/Edge types
type DeveloperProps struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Github string `json:"github"`
}
