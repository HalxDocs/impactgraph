package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"impactgraph/internal/config"
	"impactgraph/internal/repository"
)

func main() {
	cfg := config.Load()
	if cfg.CognoDBURI == "" || cfg.CognoDBPassword == "" {
		fmt.Println("COGNODB_URI/PASSWORD not set — seeding in-memory only")
		mem := repository.NewInMemoryStore()
		nodes, edges := mem.Export()
		fmt.Printf("in-memory nodes=%d edges=%d\n", len(nodes), len(edges))
		return
	}
	driver, err := neo4j.NewDriverWithContext(cfg.CognoDBURI, neo4j.BasicAuth(cfg.CognoDBUsername, cfg.CognoDBPassword, ""))
	if err != nil {
		log.Fatalf("driver: %v", err)
	}
	defer driver.Close(context.Background())
	ctx := context.Background()
	if err := driver.VerifyConnectivity(ctx); err != nil {
		log.Fatalf("verify: %v", err)
	}
	fmt.Println("connected to CognoDB, clearing and seeding...")

	// constraints per docs — idempotent MERGE requires constraints
	constraints := []string{
		"CREATE CONSTRAINT developer_id IF NOT EXISTS FOR (d:Developer) REQUIRE d.id IS UNIQUE",
		"CREATE CONSTRAINT repository_id IF NOT EXISTS FOR (r:Repository) REQUIRE r.id IS UNIQUE",
		"CREATE CONSTRAINT package_id IF NOT EXISTS FOR (p:Package) REQUIRE p.id IS UNIQUE",
		"CREATE CONSTRAINT service_id IF NOT EXISTS FOR (s:Service) REQUIRE s.id IS UNIQUE",
		"CREATE CONSTRAINT team_id IF NOT EXISTS FOR (t:Team) REQUIRE t.id IS UNIQUE",
		"CREATE CONSTRAINT organization_id IF NOT EXISTS FOR (o:Organization) REQUIRE o.id IS UNIQUE",
		"CREATE INDEX package_name IF NOT EXISTS FOR (p:Package) ON (p.name)",
		"CREATE INDEX service_name IF NOT EXISTS FOR (s:Service) ON (s.name)",
		"CREATE INDEX repository_name IF NOT EXISTS FOR (r:Repository) ON (r.name)",
	}
	sess := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	for _, c := range constraints {
		if _, err := sess.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			_, err := tx.Run(ctx, c, nil)
			return nil, err
		}); err != nil {
			fmt.Printf("constraint warn: %v (%s)\n", err, c)
		}
	}
	sess.Close(ctx)

	// Wipe existing (docs recommends MERGE idempotent, but clear for fresh seed)
	sess = driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	_, _ = sess.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		_, err := tx.Run(ctx, "MATCH (n) DETACH DELETE n", nil)
		return nil, err
	})
	sess.Close(ctx)

	mem := repository.NewInMemoryStore()
	nodes, edges := mem.Export()

	// Seed nodes via UNWIND + MERGE per docs
	type row struct {
		Id    string
		Label string
		Props map[string]any
	}
	var nodeRows []map[string]any
	for _, n := range nodes {
		props := map[string]any{}
		for k, v := range n.Properties {
			props[k] = v
		}
		props["id"] = n.ID
		props["name"] = n.Name
		nodeRows = append(nodeRows, map[string]any{"id": n.ID, "label": n.Type, "props": props})
	}

	sess = driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	_, err = sess.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		// Use separate MERGE per label to preserve label — do generic with apoc? Keep simple: MERGE by id then SET label via SET
		for _, r := range nodeRows {
			label := r["label"].(string)
			props := r["props"].(map[string]any)
			// MERGE with label
			q := fmt.Sprintf("MERGE (n:%s {id: $id}) SET n += $props", label)
			if _, err := tx.Run(ctx, q, map[string]any{"id": r["id"], "props": props}); err != nil {
				return nil, err
			}
		}
		return nil, nil
	})
	sess.Close(ctx)
	if err != nil {
		log.Fatalf("seed nodes: %v", err)
	}
	fmt.Printf("seeded %d nodes\n", len(nodeRows))

	// Seed edges via MATCH + MERGE
	sess = driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	_, err = sess.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		for _, e := range edges {
			q := fmt.Sprintf("MATCH (a {id: $src}), (b {id: $dst}) MERGE (a)-[r:%s]->(b) SET r += $props", e.Type)
			props := map[string]any{}
			for k, v := range e.Properties {
				props[k] = v
			}
			if _, err := tx.Run(ctx, q, map[string]any{"src": e.Source, "dst": e.Target, "props": props}); err != nil {
				return nil, err
			}
		}
		return nil, nil
	})
	sess.Close(ctx)
	if err != nil {
		log.Fatalf("seed edges: %v", err)
	}
	fmt.Printf("seeded %d edges\n", len(edges))

	// Verify
	sess = driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	counts, _ := sess.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		res, _ := tx.Run(ctx, "MATCH (n) RETURN count(n) as c", nil)
		c := 0
		if res.Next(ctx) {
			if v, ok := res.Record().Get("c"); ok {
				if iv, ok := v.(int64); ok {
					c = int(iv)
				}
			}
		}
		res2, _ := tx.Run(ctx, "MATCH ()-[r]->() RETURN count(r) as c", nil)
		rc := 0
		if res2.Next(ctx) {
			if v, ok := res2.Record().Get("c"); ok {
				if iv, ok := v.(int64); ok {
					rc = int(iv)
				}
			}
		}
		return map[string]int{"nodes": c, "edges": rc}, nil
	})
	sess.Close(ctx)
	fmt.Printf("verify: %+v\n", counts)

	if len(os.Args) > 1 && os.Args[1] == "check" {
		os.Exit(0)
	}
}
