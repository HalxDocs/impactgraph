package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"impactgraph/internal/models"
)

// Neo4jStore implements Store via real CognoDB.
type Neo4jStore struct {
	driver neo4j.DriverWithContext
}

func NewNeo4jStore(driver neo4j.DriverWithContext) *Neo4jStore {
	return &Neo4jStore{driver: driver}
}

func (s *Neo4jStore) Health(ctx context.Context) error {
	return s.driver.VerifyConnectivity(ctx)
}

// helpers to map neo4j records to models — keep simple; fall back to generic queries
func (s *Neo4jStore) Search(ctx context.Context, q string) ([]models.SearchResult, error) {
	if strings.TrimSpace(q) == "" {
		return []models.SearchResult{}, nil
	}
	session := s.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)
	// Cross-label search
	cypher := `
		MATCH (n)
		WHERE (n:Package OR n:Repository OR n:Service OR n:Developer OR n:Team OR n:Organization)
		  AND (toLower(coalesce(n.name,'')) CONTAINS toLower($q) OR toLower(coalesce(n.id,'')) CONTAINS toLower($q))
		RETURN labels(n)[0] as type, n.id as id, n.name as name, n LIMIT 20`
	res, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, cypher, map[string]any{"q": q})
		if err != nil {
			return nil, err
		}
		var out []models.SearchResult
		for result.Next(ctx) {
			rec := result.Record()
			typ, _ := rec.Get("type")
			id, _ := rec.Get("id")
			name, _ := rec.Get("name")
			nNode, _ := rec.Get("n")
			props := map[string]any{}
			if node, ok := nNode.(neo4j.Node); ok {
				props = node.Props
			}
			out = append(out, models.SearchResult{
				ID: fmt.Sprint(id), Type: fmt.Sprint(typ), Name: fmt.Sprint(name), Label: fmt.Sprint(typ), Properties: props,
			})
		}
		return out, result.Err()
	})
	if err != nil {
		return nil, err
	}
	return res.([]models.SearchResult), nil
}

func (s *Neo4jStore) GetEntity(ctx context.Context, typ, id string) (*models.EntityDetail, error) {
	typ = normalizeType(typ)
	session := s.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)
	cypher := fmt.Sprintf("MATCH (n:%s {id: $id}) RETURN n", typ)
	// fallback: if type unknown, match any
	if typ == "" {
		cypher = "MATCH (n {id: $id}) RETURN n"
	}
	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		res, err := tx.Run(ctx, cypher, map[string]any{"id": id})
		if err != nil {
			return nil, err
		}
		if res.Next(ctx) {
			rec := res.Record()
			val, _ := rec.Get("n")
			if node, ok := val.(neo4j.Node); ok {
				labels := node.Labels
				t := typ
				if len(labels) > 0 {
					t = labels[0]
				}
				name := fmt.Sprint(node.Props["name"])
				return &models.EntityDetail{ID: fmt.Sprint(node.Props["id"]), Type: t, Name: name, Properties: node.Props}, nil
			}
		}
		return nil, fmt.Errorf("not found")
	})
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, fmt.Errorf("not found")
	}
	return result.(*models.EntityDetail), nil
}

func (s *Neo4jStore) GetEntityGraph(ctx context.Context, typ, id string) (*models.Graph, error) {
	session := s.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)
	cypher := `MATCH (n {id: $id})-[r]-(m) RETURN n, r, m LIMIT 50`
	res, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, cypher, map[string]any{"id": id})
		if err != nil {
			return nil, err
		}
		nodes := map[string]models.Node{}
		internalToProps := map[int64]string{}
		var rels []neo4j.Relationship
		for result.Next(ctx) {
			rec := result.Record()
			if v, ok := rec.Get("n"); ok {
				if node, ok := v.(neo4j.Node); ok {
					idStr := fmt.Sprint(node.Props["id"])
					nodes[idStr] = neoNodeToModel(node)
					internalToProps[node.Id] = idStr
				}
			}
			if v, ok := rec.Get("m"); ok {
				if node, ok := v.(neo4j.Node); ok {
					idStr := fmt.Sprint(node.Props["id"])
					nodes[idStr] = neoNodeToModel(node)
					internalToProps[node.Id] = idStr
				}
			}
			if v, ok := rec.Get("r"); ok {
				if rel, ok := v.(neo4j.Relationship); ok {
					rels = append(rels, rel)
				}
			}
		}
		var edges []models.Edge
		for _, rel := range rels {
			src := internalToProps[rel.StartId]
			dst := internalToProps[rel.EndId]
			if src == "" {
				src = fmt.Sprint(rel.StartId)
			}
			if dst == "" {
				dst = fmt.Sprint(rel.EndId)
			}
			edges = append(edges, models.Edge{ID: fmt.Sprint(rel.Id), Source: src, Target: dst, Type: rel.Type, Properties: rel.Props})
		}
		var ns []models.Node
		for _, n := range nodes {
			ns = append(ns, n)
		}
		return &models.Graph{Nodes: ns, Edges: edges}, result.Err()
	})
	if err != nil {
		return nil, err
	}
	return res.(*models.Graph), nil
}

func (s *Neo4jStore) GetImpact(ctx context.Context, packageID string) (*models.ImpactResponse, error) {
	session := s.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)
	// Blast radius: reverse DEPENDS_ON to find dependents, then forward DEPLOYED_AS/CALLS
	// Package <-DEPENDS_ON- Repository -DEPLOYED_AS-> Service -CALLS-> Service
	cypher := `
		MATCH (p:Package {id: $id})
		MATCH path = (p)<-[:DEPENDS_ON*0..3]-(:Repository)-[:DEPLOYED_AS*0..1]->(:Service)-[:CALLS*0..4]->(target:Service)
		RETURN path, length(path) as depth
		UNION
		MATCH (p:Package {id: $id})
		MATCH path = (p)<-[:DEPENDS_ON*1..3]-(r:Repository)
		RETURN path, length(path) as depth
		UNION
		MATCH (p:Package {id: $id})
		MATCH path = (p)<-[:DEPENDS_ON*1..3]-(pkg:Package)
		RETURN path, length(path) as depth`
	res, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, cypher, map[string]any{"id": packageID})
		if err != nil {
			return nil, err
		}
		nodes := map[string]models.Node{}
		edges := map[string]models.Edge{}
		maxDepth := 0
		counts := models.ImpactCounts{}
		// also fetch root
		rootRes, _ := tx.Run(ctx, "MATCH (p:Package {id: $id}) RETURN p", map[string]any{"id": packageID})
		var root models.ImpactRoot
		if rootRes.Next(ctx) {
			if v, ok := rootRes.Record().Get("p"); ok {
				if node, ok := v.(neo4j.Node); ok {
					root = models.ImpactRoot{ID: fmt.Sprint(node.Props["id"]), Name: fmt.Sprint(node.Props["name"]), Type: "Package"}
					nodes[root.ID] = neoNodeToModel(node)
				}
			}
		}
		for result.Next(ctx) {
			rec := result.Record()
			depth, _ := rec.Get("depth")
			if d, ok := depth.(int64); ok && int(d) > maxDepth {
				maxDepth = int(d)
			}
			if v, ok := rec.Get("path"); ok {
				if p, ok := v.(neo4j.Path); ok {
					internalToProps := map[int64]string{}
					for _, n := range p.Nodes {
						idStr := fmt.Sprint(n.Props["id"])
						nodes[idStr] = neoNodeToModel(n)
						internalToProps[n.Id] = idStr
					}
					for _, r := range p.Relationships {
						src := internalToProps[r.StartId]
						dst := internalToProps[r.EndId]
						if src == "" {
							src = fmt.Sprint(r.StartId)
						}
						if dst == "" {
							dst = fmt.Sprint(r.EndId)
						}
						edges[fmt.Sprint(r.Id)] = models.Edge{ID: fmt.Sprint(r.Id), Source: src, Target: dst, Type: r.Type, Properties: r.Props}
					}
				}
			}
		}
		if err := result.Err(); err != nil {
			return nil, err
		}
		// counts
		for _, n := range nodes {
			switch n.Type {
			case "Repository":
				counts.Repositories++
			case "Service":
				counts.Services++
			case "Team":
				counts.Teams++
			case "Package":
				if n.ID != packageID {
					counts.Packages++
				}
			}
		}
		// teams owning affected services/repos
		teamCypher := `MATCH (t:Team)-[:OWNS]->(x) WHERE x.id IN $ids RETURN DISTINCT t`
		var ids []string
		for id := range nodes {
			ids = append(ids, id)
		}
		if len(ids) > 0 {
			tr, _ := tx.Run(ctx, teamCypher, map[string]any{"ids": ids})
			for tr.Next(ctx) {
				if v, ok := tr.Record().Get("t"); ok {
					if node, ok := v.(neo4j.Node); ok {
						idStr := fmt.Sprint(node.Props["id"])
						if _, exists := nodes[idStr]; !exists {
							nodes[idStr] = neoNodeToModel(node)
							counts.Teams++
							if maxDepth < 1 {
								maxDepth = 1
							}
						}
					}
				}
			}
		}
		var ns []models.Node
		for _, n := range nodes {
			ns = append(ns, n)
		}
		var es []models.Edge
		for _, e := range edges {
			es = append(es, e)
		}
		return &models.ImpactResponse{Root: root, Affected: counts, MaxDepth: maxDepth, Nodes: ns, Edges: es}, nil
	})
	if err != nil {
		return nil, err
	}
	return res.(*models.ImpactResponse), nil
}

func (s *Neo4jStore) GetPath(ctx context.Context, from, to string) (*models.PathResponse, error) {
	session := s.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)
	cypher := `MATCH (a {id: $from}), (b {id: $to}), path = shortestPath((a)-[:DEPENDS_ON|DEPLOYED_AS|CALLS|OWNS|BELONGS_TO|CONTRIBUTED_TO|MAINTAINS*..8]-(b)) RETURN path, length(path) as hops`
	res, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, cypher, map[string]any{"from": from, "to": to})
		if err != nil {
			return nil, err
		}
		if result.Next(ctx) {
			rec := result.Record()
			hopsVal, _ := rec.Get("hops")
			hops := 0
			if h, ok := hopsVal.(int64); ok {
				hops = int(h)
			}
			var nodes []models.Node
			var edges []models.Edge
			if v, ok := rec.Get("path"); ok {
				if p, ok := v.(neo4j.Path); ok {
					internalToProps := map[int64]string{}
					for _, n := range p.Nodes {
						nodes = append(nodes, neoNodeToModel(n))
						internalToProps[n.Id] = fmt.Sprint(n.Props["id"])
					}
					for _, r := range p.Relationships {
						src := internalToProps[r.StartId]
						dst := internalToProps[r.EndId]
						if src == "" {
							src = fmt.Sprint(r.StartId)
						}
						if dst == "" {
							dst = fmt.Sprint(r.EndId)
						}
						edges = append(edges, models.Edge{ID: fmt.Sprint(r.Id), Source: src, Target: dst, Type: r.Type, Properties: r.Props})
					}
				}
			}
			return &models.PathResponse{Exists: true, Hops: hops, Path: nodes, Edges: edges}, nil
		}
		if err := result.Err(); err != nil {
			return nil, err
		}
		return &models.PathResponse{Exists: false}, nil
	})
	if err != nil {
		return nil, err
	}
	return res.(*models.PathResponse), nil
}

func (s *Neo4jStore) GetStats(ctx context.Context) (*models.StatsResponse, error) {
	session := s.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)
	res, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		// counts per label
		labels := []string{"Repository", "Package", "Service", "Developer", "Team", "Organization"}
		counts := map[string]int{}
		for _, l := range labels {
			r, err := tx.Run(ctx, fmt.Sprintf("MATCH (n:%s) RETURN count(n) as c", l), nil)
			if err != nil {
				return nil, err
			}
			if r.Next(ctx) {
				if v, ok := r.Record().Get("c"); ok {
					if c, ok := v.(int64); ok {
						counts[l] = int(c)
					}
				}
			}
		}
		// most connected packages
		r, err := tx.Run(ctx, `MATCH (p:Package)-[d]-() RETURN p.id as id, p.name as name, count(d) as degree ORDER BY degree DESC LIMIT 5`, nil)
		if err != nil {
			return nil, err
		}
		var most []models.MostConnected
		for r.Next(ctx) {
			rec := r.Record()
			id, _ := rec.Get("id")
			name, _ := rec.Get("name")
			deg, _ := rec.Get("degree")
			c := 0
			if d, ok := deg.(int64); ok {
				c = int(d)
			}
			most = append(most, models.MostConnected{ID: fmt.Sprint(id), Name: fmt.Sprint(name), Type: "Package", Count: c})
		}
		if err := r.Err(); err != nil {
			return nil, err
		}
		return &models.StatsResponse{
			Repositories: counts["Repository"], Packages: counts["Package"], Services: counts["Service"],
			Developers: counts["Developer"], Teams: counts["Team"], Organizations: counts["Organization"],
			MostConnected: most, Recent: []models.RecentRel{},
		}, nil
	})
	if err != nil {
		return nil, err
	}
	return res.(*models.StatsResponse), nil
}

func (s *Neo4jStore) GetAllNodes(ctx context.Context) ([]models.Node, error) {
	session := s.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)
	res, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, "MATCH (n) RETURN n LIMIT 200", nil)
		if err != nil {
			return nil, err
		}
		var nodes []models.Node
		for result.Next(ctx) {
			if v, ok := result.Record().Get("n"); ok {
				if node, ok := v.(neo4j.Node); ok {
					nodes = append(nodes, neoNodeToModel(node))
				}
			}
		}
		return nodes, result.Err()
	})
	if err != nil {
		return nil, err
	}
	return res.([]models.Node), nil
}

func neoNodeToModel(n neo4j.Node) models.Node {
	id := fmt.Sprint(n.Props["id"])
	name := fmt.Sprint(n.Props["name"])
	if name == "<nil>" || name == "" {
		name = id
	}
	typ := ""
	if len(n.Labels) > 0 {
		typ = n.Labels[0]
	}
	return models.Node{ID: id, Type: typ, Label: typ, Name: name, Properties: n.Props}
}
func neoRelToModel(r neo4j.Relationship) models.Edge {
	return models.Edge{ID: fmt.Sprint(r.Id), Source: fmt.Sprint(r.StartId), Target: fmt.Sprint(r.EndId), Type: r.Type, Properties: r.Props}
}
