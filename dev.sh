#!/bin/bash

echo "🔧 Build & Run ton stack Docker"

docker-compose build

docker-compose up -d

echo "📜 Logs Caddy :"
docker-compose logs -f caddy
