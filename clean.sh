#!/bin/bash

echo "🧹 Clean complet de Docker (containers, images, volumes, networks)"

docker stop $(docker ps -aq) 2>/dev/null
docker rm -f $(docker ps -aq) 2>/dev/null

docker rmi -f $(docker images -aq) 2>/dev/null

docker volume prune -f

docker network prune -f

echo "✅ Docker clean terminé"
