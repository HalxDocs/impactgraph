package repository

// queries.go — central Cypher query definitions per spec §9-11
// Actual execution lives in neo4j.go; this file documents the precise queries
// used for impact analysis and demonstrates why a graph is needed.

const (
	// Find repository dependencies
	QueryRepoDeps = `MATCH (r:Repository {id: $repository_id})-[:DEPENDS_ON]->(p:Package) RETURN p ORDER BY p.name`

	// Multi-hop impact (reverse DEPENDS_ON then forward DEPLOYED_AS/CALLS)
	// See GetImpact in neo4j.go for UNION variant handling Package/Product depth.
	QueryImpact = `MATCH path = (p:Package {id: $package_id})<-[:DEPENDS_ON*0..3]-(:Repository)-[:DEPLOYED_AS*0..1]->(:Service)-[:CALLS*0..4]->(target:Service) RETURN path`

	// Relationally-awkward: services potentially affected
	QueryServicesAffected = `MATCH path = (p:Package {id: $package_id})<-[:DEPENDS_ON*1..3]-(:Repository)-[:DEPLOYED_AS]->(s:Service) RETURN DISTINCT s, length(path) AS distance ORDER BY distance`

	// Dependency path
	QueryShortestPath = `MATCH (a {id: $from}), (b {id: $to}), path = shortestPath((a)-[:DEPENDS_ON|DEPLOYED_AS|CALLS|OWNS|BELONGS_TO|CONTRIBUTED_TO|MAINTAINS*..8]-(b)) RETURN path`
)
