# Deploy ImpactGraph to Railway

This repo is ready for Railway (monorepo with 2 services). No new SDK — just Dockerfiles.

## 1. What Railway will run

| Service | Root dir | Dockerfile | Build arg | Health | Port |
|---------|----------|------------|-----------|--------|------|
| **backend** | `backend` | `backend/Dockerfile` | — | `GET /api/health` → 200 | `PORT` (Railway sets, app reads `os.Getenv("PORT")`) |
| **frontend** | `frontend` | `frontend/Dockerfile` | `VITE_API_URL` | `/` | `80` (nginx) |

`docker-compose.yml` is for local only. Railway ignores it and uses each Dockerfile.

## 2. Deploy via Dashboard (no CLI)

> **Fix for your screenshot error** `Railpack could not determine how to build the app.` → `Script start.sh not found`:
> You deployed the **repo root `.`** as one service. Railpack scanned `backend/` + `frontend/` at root and found no `go.mod`/`package.json` at root, so it failed. **Fix: create TWO services with Root Directory set** (see steps 3 & 5 below). After fixing, delete the failed single service. This commit adds `backend/railpack.json` + `frontend/railpack.json` + `nixpacks.toml` so both Railpack and Nixpacks now succeed when Root is set correctly.

1. Push is already done: `https://github.com/HalxDocs/impactgraph` (`main`).

2. Go to **https://railway.app/new** → **Deploy from GitHub repo** → pick `HalxDocs/impactgraph`.

3. **Add backend service**
   - Click **+ New Service → GitHub Repo** → same repo → **Settings → Service Name** = `backend` → **Root Directory** = `backend`
   - Builder = `Dockerfile` (auto-detected, `backend/Dockerfile`)
   - **Variables** (Service → Variables → RAW Editor, paste):
     ```
     COGNODB_URI=bolt+s://db-f1a4c142.bravo.databases.cognodb.cloud
     COGNODB_USERNAME=cognodb
     COGNODB_PASSWORD=2d484750fe382ed5290f3257877b9f5c
     PORT=8080
     ```
     Railway will override `PORT` at runtime — leave as is, app reads it. Use **Reference** if you prefer `${{PORT}}`.
   - **Deploy** → wait for `connected to CognoDB/Neo4j` in logs → copy **Public Domain** e.g. `https://impactgraph-backend.up.railway.app`

4. **Seed CognoDB** (once, after backend is live)
   - Locally: `cd backend && COGNODB_URI=bolt+s://... COGNODB_PASSWORD=... go run ./cmd/seed`
   - Or via Railway shell: **backend → Shell** → `go run ./cmd/seed` (needs env already set)

5. **Add frontend service** (same project, second service)
   - **+ New Service → GitHub Repo** → same `HalxDocs/impactgraph` → **Root Directory** = `frontend`
   - Builder = `Dockerfile` (`frontend/Dockerfile`)
   - **Build Arg** (important): Service → **Settings → Build → Build Args** → `VITE_API_URL=https://impactgraph-backend.up.railway.app` (paste the backend public domain, no trailing slash)
   - **Variables**: none needed (frontend is static nginx)
   - **Deploy** → copy **Public Domain** e.g. `https://impactgraph-frontend.up.railway.app` → this is your demo link for README

6. **Update README Demo link**
   ```md
   Frontend: https://impactgraph-frontend.up.railway.app
   Backend: https://impactgraph-backend.up.railway.app
   ```

7. **CORS** — backend already `AllowOrigins: ["*"]` (`backend/cmd/api/main.go` uses `cors.New(...)`). For tight prod, change to `AllowOrigins: []string{"https://impactgraph-frontend.up.railway.app"}`.

## 3. Deploy via CLI (alternative)

```bash
npm i -g @railway/cli
railway login
railway init   # pick HalxDocs/impactgraph project
# backend
railway up --service backend --root backend
railway variables --service backend set COGNODB_URI=bolt+s://... COGNODB_PASSWORD=... COGNODB_USERNAME=cognodb
# frontend (needs build arg)
railway up --service frontend --root frontend --build-arg VITE_API_URL=https://impactgraph-backend.up.railway.app
```

## 4. Local test before pushing

```bash
docker compose up --build # frontend on :5173 → nginx 80, backend on :8080 (uses backend/.env)
# or
cd backend && go run ./cmd/api        # :8080, mode=neo4j if env present else in-memory
cd frontend && VITE_API_URL=http://localhost:8080 npm run dev  # :5173 proxied to :8080
```

## 5. Common gotchas

- **If `api` is missing after clone**: was fixed in `ee17dba` (`.gitignore` had bare `api` ignoring `backend/cmd/api`). Re-clone now has it.
- **Frontend shows `in-memory` not `neo4j`**: `VITE_API_URL` not set at build time — rebuild with correct arg.
- **Backend 502**: check `COGNODB_URI` is `bolt+s://` not `bolt://` and `COGNODB_USERNAME=cognodb` (not `neo4j`) per `https://cognodb.com/docs`.
- **PORT**: don't hardcode `EXPOSE 8080` change — app reads `PORT` env, Railway injects it.

After deploy, replace screenshots in `docs/screenshots/` from `https://impactgraph-frontend.up.railway.app` at 1280×800 light+dark.
