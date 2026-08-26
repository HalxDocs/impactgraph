package repository

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"impactgraph/internal/models"
)

// InMemoryStore implements Store without a real DB. Uses same seed as neo4j seed.
type InMemoryStore struct {
	nodes map[string]models.Node // id -> Node
	edges []models.Edge
	// adjacency for traversals
	out map[string][]models.Edge
	in  map[string][]models.Edge
	recent []models.RecentRel
}

func NewInMemoryStore() *InMemoryStore {
	s := &InMemoryStore{
		nodes: make(map[string]models.Node),
		out:   make(map[string][]models.Edge),
		in:    make(map[string][]models.Edge),
	}
	seed(s)
	return s
}

func (s *InMemoryStore) Health(ctx context.Context) error { return nil }

func (s *InMemoryStore) addNode(n models.Node) {
	s.nodes[n.ID] = n
}
func (s *InMemoryStore) addEdge(e models.Edge) {
	s.edges = append(s.edges, e)
	s.out[e.Source] = append(s.out[e.Source], e)
	s.in[e.Target] = append(s.in[e.Target], e)
}

// Search across all types
func (s *InMemoryStore) Search(ctx context.Context, q string) ([]models.SearchResult, error) {
	if strings.TrimSpace(q) == "" {
		return []models.SearchResult{}, nil
	}
	lq := strings.ToLower(q)
	var res []models.SearchResult
	for _, n := range s.nodes {
		if strings.Contains(strings.ToLower(n.Name), lq) || strings.Contains(strings.ToLower(n.ID), lq) {
			props := n.Properties
			if props == nil {
				props = map[string]interface{}{}
			}
			res = append(res, models.SearchResult{ID: n.ID, Type: n.Type, Name: n.Name, Label: n.Label, Properties: props})
		}
		// also match properties like language, ecosystem
		if n.Properties != nil {
			for _, v := range n.Properties {
				if vs, ok := v.(string); ok && strings.Contains(strings.ToLower(vs), lq) {
					// avoid duplicate
					found := false
					for _, r := range res {
						if r.ID == n.ID {
							found = true
							break
						}
					}
					if !found {
						res = append(res, models.SearchResult{ID: n.ID, Type: n.Type, Name: n.Name, Label: n.Label, Properties: n.Properties})
					}
				}
			}
		}
	}
	sort.Slice(res, func(i, j int) bool { return res[i].Name < res[j].Name })
	if len(res) > 20 {
		res = res[:20]
	}
	return res, nil
}

func (s *InMemoryStore) GetEntity(ctx context.Context, typ, id string) (*models.EntityDetail, error) {
	n, ok := s.nodes[id]
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	// normalize type check leniently
	_ = normalizeType(typ)
	stats := s.computeStats(id, n.Type)
	return &models.EntityDetail{
		ID: n.ID, Type: n.Type, Name: n.Name, Properties: n.Properties, Stats: stats,
	}, nil
}

func (s *InMemoryStore) computeStats(id, typ string) *models.EntityStats {
	// For Package: count distinct repos/services/teams reachable via reverse impact and direct dependents
	// Simple: count in-neighbors by type
	st := &models.EntityStats{}
	// collect all nodes reachable forward for package usage counts
	// For demo, compute via BFS outward
	impact, _ := s.bfsImpact(id)
	for _, nid := range impact {
		node := s.nodes[nid]
		switch node.Type {
		case "Repository":
			st.UsedByRepos++
		case "Service":
			st.UsedByServices++
		case "Team":
			st.UsedByTeams++
		case "Package":
			if nid != id {
				st.UsedByPackages++
			}
		}
	}
	return st
}

func (s *InMemoryStore) GetEntityGraph(ctx context.Context, typ, id string) (*models.Graph, error) {
	n, ok := s.nodes[id]
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	// 1-hop neighborhood
	nodeSet := map[string]models.Node{n.ID: n}
	var edges []models.Edge
	for _, e := range s.edges {
		if e.Source == id || e.Target == id {
			edges = append(edges, e)
			if _, ok := nodeSet[e.Source]; !ok {
				if nn, ok := s.nodes[e.Source]; ok {
					nodeSet[e.Source] = nn
				}
			}
			if _, ok := nodeSet[e.Target]; !ok {
				if nn, ok := s.nodes[e.Target]; ok {
					nodeSet[e.Target] = nn
				}
			}
		}
	}
	nodes := make([]models.Node, 0, len(nodeSet))
	for _, v := range nodeSet {
		nodes = append(nodes, v)
	}
	sort.Slice(nodes, func(i, j int) bool { return nodes[i].ID < nodes[j].ID })
	return &models.Graph{Nodes: nodes, Edges: edges}, nil
}

func (s *InMemoryStore) GetAllNodes(ctx context.Context) ([]models.Node, error) {
	nodes := make([]models.Node, 0, len(s.nodes))
	for _, n := range s.nodes {
		nodes = append(nodes, n)
	}
	sort.Slice(nodes, func(i, j int) bool { return nodes[i].ID < nodes[j].ID })
	return nodes, nil
}

// bfsImpact does blast-radius traversal: Package -> (reverse DEPENDS_ON) -> Repository/Package -> DEPLOYED_AS -> Service -> CALLS -> Service
func (s *InMemoryStore) bfsImpact(start string) ([]string, map[string]int) {
	visited := map[string]bool{start: true}
	depth := map[string]int{start: 0}
	queue := []string{start}
	var order []string
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		curDepth := depth[cur]
		if curDepth >= 6 {
			continue
		}
		// Reverse DEPENDS_ON: find dependents (nodes that DEPENDS_ON cur)
		for _, e := range s.in[cur] {
			if e.Type == "DEPENDS_ON" {
				if !visited[e.Source] {
					visited[e.Source] = true
					depth[e.Source] = curDepth + 1
					queue = append(queue, e.Source)
					order = append(order, e.Source)
				}
			}
		}
		// Forward DEPLOYED_AS / CALLS from current (repo->service, service->service)
		for _, e := range s.out[cur] {
			if e.Type == "DEPLOYED_AS" || e.Type == "CALLS" {
				if !visited[e.Target] {
					visited[e.Target] = true
					depth[e.Target] = curDepth + 1
					queue = append(queue, e.Target)
					order = append(order, e.Target)
				}
			}
			// Also handle forward DEPENDS_ON for Package->Package transitive when start is downstream package?
			// For impact of http-client, stripe-sdk DEPENDS_ON http-client via reverse already covered.
			// For completeness, if cur is a Repository that DEPENDS_ON Package, we already traversed reverse, not forward.
		}
	}
	// also collect teams that own affected repos/services
	var teamIDs []string
	affectedSet := map[string]bool{}
	for _, id := range order {
		affectedSet[id] = true
	}
	for _, e := range s.edges {
		if e.Type == "OWNS" {
			if affectedSet[e.Target] {
				if !visited[e.Source] {
					visited[e.Source] = true
					depth[e.Source] = depth[e.Target] + 1
					teamIDs = append(teamIDs, e.Source)
				}
			}
		}
	}
	order = append(order, teamIDs...)
	return order, depth
}

func (s *InMemoryStore) GetImpact(ctx context.Context, packageID string) (*models.ImpactResponse, error) {
	root, ok := s.nodes[packageID]
	if !ok {
		return nil, fmt.Errorf("package not found")
	}
	// allow any node type as root but spec is package
	order, depth := s.bfsImpact(packageID)
	nodeSet := map[string]models.Node{root.ID: root}
	edgeSet := []models.Edge{}
	// collect edges along BFS tree
	for _, id := range order {
		if n, ok := s.nodes[id]; ok {
			nodeSet[id] = n
		}
	}
	// edges between visited nodes
	visited := map[string]bool{packageID: true}
	for _, id := range order {
		visited[id] = true
	}
	for _, e := range s.edges {
		if visited[e.Source] && visited[e.Target] {
			// only include impact-relevant + OWNS for team
			if e.Type == "DEPENDS_ON" || e.Type == "DEPLOYED_AS" || e.Type == "CALLS" || e.Type == "OWNS" {
				edgeSet = append(edgeSet, e)
			}
		}
	}
	nodes := make([]models.Node, 0, len(nodeSet))
	for _, n := range nodeSet {
		nodes = append(nodes, n)
	}
	// counts
	counts := models.ImpactCounts{}
	maxDepth := 0
	for _, id := range order {
		d := depth[id]
		if d > maxDepth {
			maxDepth = d
		}
		n := s.nodes[id]
		switch n.Type {
		case "Repository":
			counts.Repositories++
		case "Service":
			counts.Services++
		case "Team":
			counts.Teams++
		case "Package":
			counts.Packages++
		}
	}
	// sort for deterministic
	sort.Slice(nodes, func(i, j int) bool { return nodes[i].ID < nodes[j].ID })
	sort.Slice(edgeSet, func(i, j int) bool { return edgeSet[i].Source < edgeSet[j].Source })

	// build paths for UI: shortest path to each affected service/repo
	return &models.ImpactResponse{
		Root:     models.ImpactRoot{ID: root.ID, Name: root.Name, Type: root.Type},
		Affected: counts,
		MaxDepth: maxDepth,
		Nodes:    nodes,
		Edges:    edgeSet,
	}, nil
}

func (s *InMemoryStore) GetPath(ctx context.Context, from, to string) (*models.PathResponse, error) {
	if _, ok := s.nodes[from]; !ok {
		return &models.PathResponse{Exists: false, Reason: "source not found"}, nil
	}
	if _, ok := s.nodes[to]; !ok {
		return &models.PathResponse{Exists: false, Reason: "target not found"}, nil
	}
	// BFS shortest path undirected-ish but directed for dependency? Allow both directions for discovery but prefer directed.
	// Build undirected adjacency for shortestPath
	undirected := map[string][]string{}
	edgeMap := map[string]models.Edge{} // key source-target-type
	for _, e := range s.edges {
		undirected[e.Source] = append(undirected[e.Source], e.Target)
		undirected[e.Target] = append(undirected[e.Target], e.Source)
		edgeMap[e.Source+"|"+e.Target] = e
		edgeMap[e.Target+"|"+e.Source] = e // for undirected lookup
	}
	// BFS
	prev := map[string]string{from: ""}
	queue := []string{from}
	visited := map[string]bool{from: true}
	found := false
	for len(queue) > 0 && !found {
		cur := queue[0]
		queue = queue[1:]
		for _, nb := range undirected[cur] {
			if !visited[nb] {
				visited[nb] = true
				prev[nb] = cur
				if nb == to {
					found = true
					break
				}
				queue = append(queue, nb)
			}
		}
	}
	if !found && from != to {
		return &models.PathResponse{Exists: false, Hops: 0, Path: []models.Node{}, Edges: []models.Edge{}}, nil
	}
	// reconstruct
	var pathIDs []string
	cur := to
	for cur != "" {
		pathIDs = append([]string{cur}, pathIDs...)
		cur = prev[cur]
		if cur == "" && pathIDs[0] != from {
			break
		}
	}
	// ensure starts with from
	if len(pathIDs) == 0 || pathIDs[0] != from {
		pathIDs = append([]string{from}, pathIDs...)
	}
	var pathNodes []models.Node
	var pathEdges []models.Edge
	for _, id := range pathIDs {
		pathNodes = append(pathNodes, s.nodes[id])
	}
	for i := 0; i < len(pathIDs)-1; i++ {
		a, b := pathIDs[i], pathIDs[i+1]
		if e, ok := edgeMap[a+"|"+b]; ok {
			// ensure direction matches stored edge
			// if edge is OWNS (Team->Repo) but path goes Repo->Team, keep as stored but include
			pathEdges = append(pathEdges, e)
		} else {
			pathEdges = append(pathEdges, models.Edge{ID: fmt.Sprintf("e-%s-%s", a, b), Source: a, Target: b, Type: "RELATED"})
		}
	}
	return &models.PathResponse{Exists: true, Hops: len(pathNodes) - 1, Path: pathNodes, Edges: pathEdges}, nil
}

func (s *InMemoryStore) GetStats(ctx context.Context) (*models.StatsResponse, error) {
	counts := map[string]int{}
	for _, n := range s.nodes {
		counts[n.Type]++
	}
	// most connected: degree = in+out edges
	deg := map[string]int{}
	for _, e := range s.edges {
		deg[e.Source]++
		deg[e.Target]++
	}
	type kv struct {
		id, name, typ string
		c             int
	}
	var kvs []kv
	for id, c := range deg {
		n := s.nodes[id]
		// only packages for "most connected packages" per spec, but include all
		if n.Type == "Package" {
			kvs = append(kvs, kv{id, n.Name, n.Type, c})
		}
	}
	sort.Slice(kvs, func(i, j int) bool { return kvs[i].c > kvs[j].c })
	var most []models.MostConnected
	for i := 0; i < len(kvs) && i < 5; i++ {
		most = append(most, models.MostConnected{ID: kvs[i].id, Name: kvs[i].name, Type: kvs[i].typ, Count: kvs[i].c})
	}
	recent := s.recent
	if len(recent) > 10 {
		recent = recent[:10]
	}
	return &models.StatsResponse{
		Repositories:  counts["Repository"],
		Packages:      counts["Package"],
		Services:      counts["Service"],
		Developers:    counts["Developer"],
		Teams:         counts["Team"],
		Organizations: counts["Organization"],
		MostConnected: most,
		Recent:        recent,
	}, nil
}

func (s *InMemoryStore) Export() (map[string]models.Node, []models.Edge) {
	return s.nodes, s.edges
}
