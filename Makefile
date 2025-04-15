.PHONY: help run migrate clean dev build logs stop rebuild

GO=go
API_PATH=cmd/api/main.go
MIGRATE_PATH=cmd/migrate/main.go

GREEN=\033[0;32m
BLUE=\033[0;34m
NC=\033[0m

help:
	@echo ""
	@echo "📦  ${BLUE}Commandes disponibles :${NC}"
	@grep -E '^[a-zA-Z_-]+:.*?##' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  ${GREEN}%-12s${NC} %s\n", $$1, $$2}'
	@echo ""

run: ## Lance l'API Go en local (hors Docker)
	$(GO) run $(API_PATH)

migrate: ## Lance le script de migration Go
	$(GO) run $(MIGRATE_PATH)

clean: ## Supprime tous les conteneurs, images, volumes et réseaux Docker
	./clean.sh

dev: ## Build & démarre les services via Docker Compose
	./dev.sh

build: ## Build les images Docker sans les lancer
	docker-compose build

logs: ## Affiche les logs live de Caddy
	docker-compose logs -f caddy

stop: ## Stoppe et supprime tous les conteneurs
	docker-compose down -v

rebuild:
	@make clean
	@make dev
