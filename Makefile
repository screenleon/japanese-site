.PHONY: help lint-rules corpus-scale vet test clean bake-static build-static \
        bootstrap dev start build dist \
        run web-dev web-build \
        seed-jmdict seed-kanjidic2 seed-jlpt seed-tatoeba seed-derive seed-corpus seed-all

# ── Quick start ─────────────────────────────────────────────────────────────
help:
	@echo "First-time setup:    make bootstrap   # install deps + seed DB"
	@echo "Daily development:   make dev         # backend + frontend hot-reload"
	@echo "One-shot start:      make start       # build everything, serve on :8080"
	@echo ""
	@echo "Other targets:"
	@echo "  build         Build production server binary + frontend bundle"
	@echo "  dist          Stage everything into dist/ (binary + web/dist + data)"
	@echo "  test          Go unit tests"
	@echo "  vet           Go static analysis"
	@echo "  lint-rules    Layered-rule lint"
	@echo "  corpus-scale  Report corpus scale against current learning-content floors"
	@echo "  clean         Remove build artifacts and dev SQLite"
	@echo "  seed-all      Re-run full corpus pipeline"
	@echo "  seed-corpus   Reload only the curated L1 corpus"

# ── Bootstrap & dev ─────────────────────────────────────────────────────────
bootstrap: bootstrap-server bootstrap-web seed-all
	@echo ""
	@echo "✓ bootstrap done. Run 'make dev' or 'make start'."

bootstrap-server:
	cd server && go mod download

bootstrap-web:
	cd web && npm install --no-fund --no-audit

# Run backend on :8080 and Vite dev server on :5173 simultaneously.
# Ctrl-C cleans up both.
dev:
	@echo "API on :8080  ・  Web on :5173 (proxies /api → backend)"
	@trap 'kill 0' INT TERM; \
	(cd server && go run ./cmd/api) & \
	(cd web && npm run dev -- --host 127.0.0.1) & \
	wait

# ── One-shot production-style start ─────────────────────────────────────────
# Build frontend bundle and run the Go server with STATIC_DIR pointing at it.
# Single port (:8080) serves both the SPA and /api/*.
start: build
	@echo "Serving on http://localhost:8080  ・  Ctrl-C to stop"
	cd server && STATIC_DIR=../web/dist ./bin/api

# ── Build (production-style) ────────────────────────────────────────────────
build: build-server build-web

build-server:
	cd server && go build -o bin/api ./cmd/api

build-web:
	cd web && npm run build

# Stage everything into dist/ — copy this directory to deploy.
dist: build seed-all
	rm -rf dist
	mkdir -p dist/web dist/data
	cp server/bin/api dist/api
	cp -r web/dist/* dist/web/
	cp server/data/japanese-site.sqlite dist/data/
	cp server/data/external.lock dist/data/
	@echo ""
	@echo "✓ dist/ ready."
	@echo "Deploy: copy dist/ to target host, then run there:"
	@echo "  STATIC_DIR=web DB_PATH=data/japanese-site.sqlite ./api"

# ── Lower-level / individual targets ────────────────────────────────────────
run:
	cd server && go run ./cmd/api

web-dev:
	cd web && npm run dev

web-build:
	cd web && npm run build

vet:
	cd server && go vet ./...

test:
	cd server && go test ./...

bake-static:
	@rm -rf web/public/data
	@mkdir -p web/public/data/grammar
	@for level in N1 N2 N3 N4 N5; do \
		src="server/data/corpus/grammar/$$level"; \
		[ -d "$$src" ] || continue; \
		dst="web/public/data/grammar/$$level"; \
		mkdir -p "$$dst"; \
		cp $$src/*.json $$src/*.examples.jsonl "$$dst/" 2>/dev/null || true; \
		ls $$src/*.json 2>/dev/null | xargs -n1 basename | sed 's/\.json$$//' | jq -R . | jq -s . > "$$dst/_index.json"; \
	done
	@cp -r server/data/corpus/vocab web/public/data/
	@cp -r server/data/corpus/kanji web/public/data/
	@echo "bake-static: web/public/data populated"

build-static: bake-static
	cd web && VITE_DEPLOY_MODE=static VITE_DEPLOY_BASE=/japanese-site/ npm run build

lint-rules:
	bash scripts/lint-rules.sh

corpus-scale:
	bash scripts/check-corpus-scale.sh

clean:
	rm -rf server/bin web/dist dist
	rm -f server/data/*.sqlite server/data/*.sqlite-*

# ── Seed pipeline ──────────────────────────────────────────────────────────
seed-jmdict:
	cd server && go run ./cmd/seed jmdict $(SEED_ARGS)

seed-kanjidic2:
	cd server && go run ./cmd/seed kanjidic2 $(SEED_ARGS)

seed-jlpt:
	cd server && go run ./cmd/seed jlpt $(SEED_ARGS)

seed-tatoeba:
	cd server && go run ./cmd/seed tatoeba $(SEED_ARGS)

seed-derive:
	cd server && go run ./cmd/seed derive $(SEED_ARGS)

seed-corpus:
	cd server && go run ./cmd/seed corpus $(SEED_ARGS)

seed-all:
	cd server && go run ./cmd/seed all $(SEED_ARGS)
