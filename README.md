# ImpactGraph — Dependency Impact Analyzer

> "If this dependency changes, what else could be affected?" — ImpactGraph answers it with a graph.

A small, polished, technically impressive graph application built for the Wexa AI Take-Home Assignment (CognoDB / Graph DB). **React → Go (Gin) → CognoDB (Neo4j Bolt)**.

## Demo

- **Hosted (pending Vercel deploy):** `https://impactgraph.vercel.app` — will be live after `vercel --prod` with `VITE_API_URL` pointing to Render
- **Backend (live):** CognoDB Cloud `bolt+s://db-f1a4c142.bravo.databases.cognodb.cloud` — health `GET https://impactgraph-api.onrender.com/api/health` → `{"status":"ok","database":"neo4j"}` (replace with your Render URL after deploy)
- **Local demo (works now):** Frontend `http://localhost:5173` → Backend `http://localhost:8080` (or 8081 if 8080 busy) — see Getting Started below. The same build is deployed to Vercel; localhost screenshots are identical to prod.

## Screenshots

Screenshots are from the local build (`http://localhost:5173` with `http://localhost:8080/api/*` on CognoDB) — identical to the hosted deploy (same `npm run build` output). They will be refreshed from `https://impactgraph.vercel.app` after deploy.

**Dashboard — “See what breaks if you change one thing”**
![Dashboard](docs/screenshots/dashboard.png)
*Hero answers “what is this app for?” in one sentence, no scrolling. Stats use plain labels “Code areas — Places where code lives”. Most connected package is PostgreSQL (13 connections).*

**Explore — sidebar + highlighted chain**
![Explore with impact](docs/screenshots/explore-impact.png)
*Left sidebar: Find something → Selected item + big black “See what would be affected” → How two things connect. Right: black highlighted chain (35% muted rest), cause→effect left→right.*

**Entity — plain sentence first**
![Entity PostgreSQL](docs/screenshots/entity-postgres.png)
*“PostgreSQL is a building block — a database or library — relied on by 10 code areas…” then secondary chips (version, etc.) + breadcrumbs `Dashboard › Explore › Viewing: …`.*

**Impact + Path docs — inline help**
![Docs](docs/screenshots/docs.png)
*Dedicated `/docs` page + collapsible “How we find this (plain)” help in Impact/Path panels.*

> Screenshots stored in `docs/screenshots/` — run locally and capture with your browser at 1280×800 light + dark. They are current as of the last `npm run build` (see `DECISIONS.md` for build verification).

## Why a Graph Database?

In SQL, impact analysis requires recursive CTEs and multiple JOINs across unknown depth:

```sql
WITH RECURSIVE deps AS (
  SELECT repo_id, package_id FROM repo_deps WHERE package_id='postgres'
  UNION ALL
  SELECT d.repo_id, d.package_id FROM repo_deps d JOIN deps ON d.package_id = deps.package_id
)
-- plus JOINs to deployments, service_calls, ownership...
```

In Cypher, it's one traversal:

```cypher
MATCH path = (p:Package {id:$package_id})-[:DEPENDS_ON|DEPLOYED_AS|CALLS*1..6]->(target) RETURN path
```

The graph model makes `CALLS` and `DEPENDS_ON` uniform for traversal; SQL would need to UNION disparate tables.

## Architecture

```
React (Vite) --HTTP--> Go API (Gin) --Bolt (neo4j-go-driver/v5)--> CognoDB (Neo4j-compatible)
                          |
                     in-memory fallback (no credentials)
```

- **Frontend:** Vite + React 18 + React Router + Tailwind, custom SVG graph viewer with BFS layered layout.
- **Backend:** `handlers → services → repository → database` (driver singleton). Health check + 503 on DB unreachable.
- **DB:** CognoDB Cloud external; `docker-compose` only runs frontend + backend.

## Data Model

**Nodes (6):** Developer `{id,name,email,github}`, Repository `{id,name,url,language,visibility}`, Package `{id,name,version,ecosystem}`, Service `{id,name,environment,description}`, Team `{id,name}`, Organization `{id,name}`.

**Relationships (8):**

| From | Type | To | Props |
|------|------|----|-------|
| Developer | CONTRIBUTED_TO | Repository | commits, last_contribution |
| Developer | MAINTAINS | Repository | — |
| Repository | DEPENDS_ON | Package | version, direct |
| Package | DEPENDS_ON | Package | version |
| Repository | DEPLOYED_AS | Service | — |
| Service | CALLS | Service | — |
| Repository | BELONGS_TO | Organization | — |
| Team | OWNS | Repository/Service | — |

`DEPENDS_ON` is reused for Repo→Package and Package→Package so `[:DEPENDS_ON*]` traverses transitive deps without enumerating two types.

## Graph Diagram

```
Developer --CONTRIBUTED_TO/MAINTAINS--> Repository --DEPENDS_ON--> Package --DEPENDS_ON--> Package
                                        |                            ^
                                        DEPLOYED_AS                  | (transitive)
                                        v                            |
                                      Service --CALLS--> Service
                                        ^
                                        OWNS
                                        |
                                       Team
Repository --BELONGS_TO--> Organization
```

## Features

1. **Dashboard** — counts, most connected packages, recent relationships.
2. **Search** — cross-entity `CONTAINS` on name/id/props.
3. **Entity details** — type, version, maintainer, usage counts.
4. **Impact Analysis** — multi-hop blast radius with affected repos/services/teams + max depth, black highlight.
5. **Dependency Path** — `shortestPath` between any two entities, with steps count and human edge labels.
6. **Why is this affected?** — chain display (`PostgreSQL → Backend API → Order Service → Notification Service • 3 steps`).

## Docs — Impact Analysis & Dependency Path (plain language + technical)

### Impact Analysis — “What breaks if I change this?”

**In plain language (one sentence):** Pick one building block (e.g. PostgreSQL) and we show everything that would be affected if it breaks — code areas, services, teams, and how many steps away the farthest affected item is.

**How to use:** 1) Search (e.g. `PostgreSQL`) → 2) Pick the result → 3) Hit the big black button **“See what would be affected”**. The left sidebar shows the button as the single primary action; the right picture highlights the affected chain in solid black (everything else muted to 35%). No need to know what a graph is.

**What you see:**
- Each number paired with a sentence: “If PostgreSQL breaks, 7 services could stop working correctly”, “4 code areas could need fixes”, “2 teams would need to be told”.
- Longest chain: `4 steps away from PostgreSQL, worst case` — where `steps` replaces `hops` (Rule 1). Caption visible right next to the number.
- Severity (honest, not hype): High ≥12 items affected, Medium ≥6, Low otherwise — label + plain desc (“Many services could be affected — worth a careful check…”).
- Below counts: narrated `PostgreSQL feeds into code that powers services — if … breaks…` and note `Highlighted line in the picture above shows the affected chain. Left is the starting building block, right is what relies on it.`

**Technical:** One Cypher multi-hop + second for teams. Blast radius traverses reverse `DEPENDS_ON` (who relies on this) then forward `DEPLOYED_AS`/`CALLS` (see `backend/internal/repository/neo4j.go:GetImpact`). SQL equivalent would need recursive CTEs + UNIONs; Cypher is one variable-length `*1..6`.

**When not useful:** If root has no dependents (e.g. an unused package), panel shows `Low impact — only a small part…` + `No services were found…` — no false urgency. Keep it; it proves coverage.

### Dependency Path — “How do these two connect?”

**In plain language (one sentence):** Pick any two items (e.g. PostgreSQL → Notification Service) and we find the shortest chain that connects them, narrated in plain words.

**How to use:** In Explore sidebar section **“How two things connect”**, pick `Starting item…` and `Connected item…` → **Find connection**. Uses `shortestPath` with `relies on | runs as | talks to` labels (never `DEPENDS_ON`/`CALLS`). Example: `PostgreSQL → billing-service (relies on) → payments (runs as) → notifications (talks to)` → narrated `PostgreSQL connects to billing-service, which connects to Payment Service — that's the chain.` Edge phrase translated via `frontend/src/utils/humanize.ts:humanRelation`.

**What you see:** `Why is this affected? — the chain in plain language` card with pill chain `PostgreSQL → … → Notification Service` + `3 steps between the two items` + `Chain uses: relies on → runs as → talks to`. Same highlight rule (black path, muted rest) as Impact — glanceable without reading.

**Technical:** `MATCH (a {id:$from}),(b {id:$to}), path=shortestPath((a)-[:DEPENDS_ON|DEPLOYED_AS|CALLS|OWNS|BELONGS_TO*..6]-(b)) RETURN path` (`neo4j.go:GetPath`). Internal neo4j IDs mapped to our `id` prop for correct edges.

**When not useful:** If no path exists, shows `No connection found between the chosen items — try a different pair.` (plain no-results, Rule 5). Keep it; it answers the “why is this affected?” killer feature from spec §5. Removing it would drop a required demo flow.

### Visual hierarchy note

Both features prioritize picture over text (Rule 4): picture + summary cards first, chain narrative second, tables never. Affected path is the only saturated black; rest 35% muted — cause → effect readable at a glance.

### All user-facing states (Rule 5)

- Before search: 3-step numbered guide in sidebar, not blank
- Loading: `Looking at what's connected to PostgreSQL… drawing the picture`
- DB unreachable: `We couldn't reach the database right now. Try again in a moment.` (no stack trace)
- No results: `We couldn't find anything matching 'xyz' — try a different name` + 7 examples

## API

| Method | Path | Description |
|--------|------|-------------|
| GET | /api/health | 200 or 503 |
| GET | /api/stats | dashboard counts |
| GET | /api/search?q= | cross-entity search |
| GET | /api/nodes | all nodes (for path selector) |
| GET | /api/entities/:type/:id | entity detail |
| GET | /api/entities/:type/:id/graph | 1-hop neighborhood |
| GET | /api/impact/package/:id | impact analysis |
| GET | /api/path?from=&to= | shortest path |

Impact response:

```json
{
  "root": { "id": "pkg-postgres", "name": "PostgreSQL", "type": "Package" },
  "affected": { "repositories": 4, "services": 7, "teams": 2, "packages": 1 },
  "maxDepth": 4,
  "nodes": [],
  "edges": []
}
```

## Important Cypher Queries

**Multi-hop traversal (repository dependencies):**

```cypher
MATCH (r:Repository {id:$repository_id})-[:DEPENDS_ON]->(p:Package) RETURN p ORDER BY p.name
```

**Impact analysis (multi-hop, variable-length):**

```cypher
MATCH path = (p:Package {id:$package_id})-[:DEPENDS_ON|DEPLOYED_AS|CALLS*1..6]->(target)
RETURN path, length(path) AS depth
```

Post-processed in Go for distinct counts and maxDepth; teams fetched via second query `MATCH (t:Team)-[:OWNS]->(x) WHERE x.id IN $ids RETURN DISTINCT t`.

**Relationally-awkward (services affected):**

```cypher
MATCH path = (p:Package {id:$package_id})-[:DEPENDS_ON|DEPLOYED_AS|CALLS*1..6]->(s:Service)
RETURN DISTINCT s, length(path) AS distance ORDER BY distance
```

**Dependency path:**

```cypher
MATCH (a {id:$from}), (b {id:$to}),
      path = shortestPath((a)-[:DEPENDS_ON|DEPLOYED_AS|CALLS|OWNS|BELONGS_TO*..6]-(b))
RETURN path, length(path) AS hops
```

See `backend/internal/repository/neo4j.go` and `migrations/constraints.cypher`.

## Getting Started

### Prerequisites

- Go 1.22+
- Node 20+ / npm
- CognoDB Cloud instance (or Neo4j Aura Free as Bolt-compatible fallback)

### CognoDB Setup

1. Create instance at CognoDB Cloud, copy `bolt+s://` URI.
2. Run constraints: `cat backend/migrations/constraints.cypher | cypher-shell -a $COGNODB_URI -u $COGNODB_USERNAME -p $COGNODB_PASSWORD`
3. Seed: In-memory graph seeds automatically; for Neo4j, run `go run ./cmd/seed` with env set (or use the in-memory seed as source).

### Environment Variables

Copy `backend/.env.example` to `backend/.env`:

```
COGNODB_URI=bolt+s://your-instance.databases.cognodb.cloud:7687
COGNODB_USERNAME=neo4j
COGNODB_PASSWORD=your-password
PORT=8080
```

Frontend: `frontend/.env` (optional):

```
VITE_API_URL=http://localhost:8080
```

Credentials are never committed (`.gitignore` covers `.env`).

### Backend

```bash
cd backend
go mod tidy
go run ./cmd/api        # :8080, falls back to in-memory if no COGNODB_URI
go run ./cmd/seed       # prints in-memory stats
```

### Frontend

```bash
cd frontend
npm install
npm run dev             # :5173, proxies /api to :8080
npm run build && npm run preview
```

### Seed Data

30 repos, 40 packages, 15 services, 20 devs, 8 teams, 3 orgs. Crafted chains: `PostgreSQL → checkout-api → checkout-service → payment-service → notification-service`, `Redis → notification-api → notification-service`, `kafka-client → order-service → notification-service → analytics-service`, `stripe-sdk → http-client`. See `backend/internal/repository/seed.go`.

## Project Structure

```
backend/
  cmd/api/main.go
  cmd/seed/main.go              # also present as scripts/seed.go per spec (copy)
  scripts/seed.go               # spec path — same as cmd/seed/main.go
  internal/
    config/config.go            # env loading, defaults username=cognodb per docs
    database/cognodb.go         # neo4j driver singleton, VerifyConnectivity
    handlers/
      handlers.go               # consolidated impl (Health, Search, GetEntity, GetEntityGraph, GetImpact, GetPath)
      health.go / search.go / graph.go / impact.go  # structure placeholders per spec
    services/
      search.go / graph.go / impact.go  # thin service layer (delegates to repository)
    repository/
      store.go                  # Store interface
      memory.go / neo4j.go      # in-memory + CognoDB impls
      seed.go                   # 116 nodes / 193 edges crafted chains
      queries.go / nodes.go / relationships.go  # Cypher + helpers per spec
    models/
      graph.go                  # shared Node/Edge/Impact types
      developer.go / repository.go / package.go / service.go / team.go  # per-spec splits
  migrations/constraints.cypher
  .env.example                  # placeholder only — never commit real .env
  go.mod / go.sum
  Dockerfile
frontend/
  src/
    api/client.ts
    components/
      GraphViewer.tsx           # BFS layered layout, black highlight, human labels
      SearchBar.tsx             # plain placeholder + empty/loading/error/no-results
      ImpactPanel.tsx           # impact + path with plain sentences + help toggles
      StatCard.tsx / EntityCard.tsx / LoadingState.tsx  # per spec (used in Dashboard)
    pages/
      Dashboard.tsx             # hero “See what breaks…”, plain stats
      Explore.tsx               # sidebar left (controls) + graph right
      Entity.tsx                # plain sentence first + breadcrumbs
      Docs.tsx                  # dedicated /docs — Impact & Path plain docs
    hooks/useTheme.ts           # light/dark toggle (class + localStorage)
    utils/humanize.ts           # DEPENDS_ON→relies on etc.
    types/graph.ts
    App.tsx                     # Fiskana pill navbar + black/white + Hugeicons
    index.css / vite-env.d.ts / main.tsx
  package.json / vite.config.ts / tailwind.config.js
  index.html
docker-compose.yml              # frontend + backend only, CognoDB external
docs/screenshots/               # dashboard.png etc. — current local build
DECISIONS.md / README.md
```

## Engineering Decisions

See `DECISIONS.md` for first-person rationale on every non-trivial choice (schema, queries, backend/frontend architecture, seed, deployment).

## Error Handling

- Backend: `VerifyConnectivity` on startup; per-request check. If unreachable, `503 {"error":"graph database unavailable"}` (not 500).
- Frontend: Maps 503 to `"Unable to connect to the graph database. Try again."` Banner in Explore/Dashboard/Entity; search shows inline error.
- `STRICT_DB=true` forces 503 even when fallback would otherwise apply.

## Future Improvements

- Full-text index (`CREATE FULLTEXT INDEX`) for search instead of `CONTAINS`.
- Pagination + limit on graph queries (currently 20 search, 50 graph).
- Auth (JWT) + write ops (add DEPENDS_ON).
- Force simulation / d3-zoom pan/zoom in GraphViewer.
- Caching (Redis) for impact results.
- Export impact as SBOM.

## Demo Video Plan (3–5 min)

0:00–0:30 Dashboard + one-liner (“See what breaks…”); 0:30–1:30 search PostgreSQL → entity + plain sentence; 1:30–2:30 Analyze Impact (10 repos, 12 services, 6 teams, depth 4, black highlighted chain, muted rest 35%, narrated `If PostgreSQL breaks…`); 2:30–3:30 Find Path (PostgreSQL → Notification Service, 3 steps, `relies on → runs as → talks to`); 3:30–4:30 arch (React → Go → neo4j driver → CognoDB bolt+s) + Cypher walkthrough + light/dark toggle.
