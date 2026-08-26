# Root Dockerfile — DO NOT USE for Railway monorepo
# This repo is a monorepo: backend/ and frontend/ each have their own Dockerfile.
# If you deploy Root "." you will get Railpack error "Script start.sh not found".
# Fix: In Railway dashboard, create TWO services:
#   Service 1: Root Directory = backend, Dockerfile = backend/Dockerfile
#   Service 2: Root Directory = frontend, Dockerfile = frontend/Dockerfile, Build Arg VITE_API_URL=<backend-url>
# See RAILWAY.md for full steps. This file just makes the error obvious if you deploy root.
FROM alpine:3.19
RUN echo "ERROR: Deploying root '.' is wrong for this monorepo. Set Root Directory to backend or frontend." && exit 1
