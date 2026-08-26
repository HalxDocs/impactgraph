package models

// Node represents a generic graph node returned to the frontend.
type Node struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Label      string                 `json:"label"`
	Name       string                 `json:"name"`
	Properties map[string]interface{} `json:"properties"`
}

// Edge represents a relationship.
type Edge struct {
	ID         string                 `json:"id"`
	Source     string                 `json:"source"`
	Target     string                 `json:"target"`
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}

// Graph is the payload for visualization.
type Graph struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}

// EntityDetail is the response for GET /api/entities/:type/:id
type EntityDetail struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Name       string                 `json:"name"`
	Properties map[string]interface{} `json:"properties"`
	Stats      *EntityStats           `json:"stats,omitempty"`
}

type EntityStats struct {
	UsedByRepos    int `json:"usedByRepos"`
	UsedByServices int `json:"usedByServices"`
	UsedByTeams    int `json:"usedByTeams"`
	UsedByPackages int `json:"usedByPackages"`
}

// SearchResult is one hit from /api/search
type SearchResult struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Name       string                 `json:"name"`
	Label      string                 `json:"label"`
	Properties map[string]interface{} `json:"properties"`
}

// ImpactResponse for GET /api/impact/package/:id
type ImpactResponse struct {
	Root     ImpactRoot    `json:"root"`
	Affected ImpactCounts  `json:"affected"`
	MaxDepth int           `json:"maxDepth"`
	Nodes    []Node        `json:"nodes"`
	Edges    []Edge        `json:"edges"`
	Paths    [][]string    `json:"paths,omitempty"`
}

type ImpactRoot struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type ImpactCounts struct {
	Repositories int `json:"repositories"`
	Services     int `json:"services"`
	Teams        int `json:"teams"`
	Packages     int `json:"packages"`
}

// PathResponse for GET /api/path?from=&to=
type PathResponse struct {
	Exists bool   `json:"exists"`
	Hops   int    `json:"hops"`
	Path   []Node `json:"path"`
	Edges  []Edge `json:"edges"`
	Reason string `json:"reason,omitempty"`
}

// StatsResponse for dashboard
type StatsResponse struct {
	Repositories int            `json:"repositories"`
	Packages     int            `json:"packages"`
	Services     int            `json:"services"`
	Developers   int            `json:"developers"`
	Teams        int            `json:"teams"`
	Organizations int           `json:"organizations"`
	MostConnected []MostConnected `json:"mostConnected"`
	Recent       []RecentRel     `json:"recent"`
}

type MostConnected struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	Count int    `json:"count"`
}

type RecentRel struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Type   string `json:"type"`
	Detail string `json:"detail"`
}
