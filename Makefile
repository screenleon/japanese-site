.PHONY: help lint lint-rules lint-grammar lint-vocab backlog-render lint-backlog-render corpus-scale vet test test-dump-grammar-examples test-lint-grammar test-lint-vocab test-validate-backlog-schema clean bake-static dump-grammar-examples build-static \
        bootstrap dev start build dist dist-update \
        run web-dev web-build \
        seed-jmdict seed-kanjidic2 seed-jlpt seed-tatoeba seed-derive seed-corpus seed-all \
        db-update

# ── Quick start ─────────────────────────────────────────────────────────────
help:
	@echo "First-time setup:    make bootstrap   # install deps + seed DB"
	@echo "Daily development:   make dev         # backend + frontend hot-reload"
	@echo "One-shot start:      make start       # build everything, serve on :8080"
	@echo ""
	@echo "Other targets:"
	@echo "  build         Build production server binary + frontend bundle"
	@echo "  dist          Stage into dist/ — full rebuild + seed-all (slow, fresh deploy)"
	@echo "  dist-update   Stage into dist/ — build + seed-corpus only (fast, content update)"
	@echo "  db-update     Reload curated corpus into local dev DB (vocab + grammar questions)"
	@echo "  test          Go unit tests"
	@echo "  test-validate-backlog-schema Backlog schema validator fixtures"
	@echo "  vet           Go static analysis"
	@echo "  lint-rules    Layered-rule lint"
	@echo "  backlog-render Regenerate BACKLOG.md index + status headings"
	@echo "  lint-backlog-render Verify BACKLOG.md generated sections are current"
	@echo "  corpus-scale  Report corpus scale against current learning-content floors"
	@echo "  clean         Remove build artifacts and dev SQLite"
	@echo "  seed-all      Re-run full corpus pipeline (incl. external data download)"
	@echo "  seed-corpus   Reload only the curated L1 corpus (vocab + grammar questions)"
	@echo ""
	@echo "Deploy workflow (content update, no external data change):"
	@echo "  make dist-update   # builds binary + frontend + reloads corpus → dist/"
	@echo "  # copy dist/ to host; migrations apply automatically on first server start"

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
# Use dist-update for content-only deploys (skips slow external data pipeline).
dist: build seed-all
	rm -rf dist
	mkdir -p dist/web dist/data
	cp server/bin/api dist/api
	cp -r web/dist/* dist/web/
	cp server/data/japanese-site.sqlite dist/data/
	cp server/data/external.lock dist/data/
	@echo ""
	@echo "✓ dist/ ready. DB migrations apply automatically on first server start."
	@echo "Deploy: copy dist/ to target host, then run:"
	@echo "  STATIC_DIR=web DB_PATH=data/japanese-site.sqlite ./api"

# Fast deploy for content-only updates (new grammar/vocab, no external data change).
# Skips jmdict/kanjidic2/tatoeba download — only reloads curated corpus.
dist-update: build seed-corpus
	rm -rf dist
	mkdir -p dist/web dist/data
	cp server/bin/api dist/api
	cp -r web/dist/* dist/web/
	cp server/data/japanese-site.sqlite dist/data/
	cp server/data/external.lock dist/data/
	@echo ""
	@echo "✓ dist/ ready (corpus updated, external data unchanged)."
	@echo "Deploy: copy dist/ to target host, then run:"
	@echo "  STATIC_DIR=web DB_PATH=data/japanese-site.sqlite ./api"

# Reload curated corpus into local dev DB (run after adding new vocab/grammar content).
db-update: seed-corpus
	@echo "✓ Local dev DB updated with latest corpus."

# ── Lower-level / individual targets ────────────────────────────────────────
run:
	cd server && go run ./cmd/api

web-dev:
	cd web && npm run dev

web-build:
	cd web && npm run build

vet:
	cd server && go vet ./...

test: test-dump-grammar-examples test-lint-grammar test-lint-vocab test-validate-backlog-schema
	cd server && go test ./...

test-dump-grammar-examples:
	bash scripts/test-dump-grammar-examples.sh

test-lint-grammar:
	bash scripts/test-lint-grammar.sh

test-lint-vocab:
	bash scripts/test-lint-vocab.sh

test-validate-backlog-schema:
	bash scripts/test-validate-backlog-schema.sh

bake-static:
	@command -v jq >/dev/null || { echo "bake-static: jq is required (install via apt/brew)"; exit 1; }
	@rm -rf web/public/data
	@mkdir -p web/public/data/grammar
	@for level in N1 N2 N3 N4 N5; do \
		src="server/data/corpus/grammar/$$level"; \
		[ -d "$$src" ] || continue; \
		count=$$(ls $$src/*.json 2>/dev/null | wc -l); \
		[ "$$count" -gt 0 ] || { echo "bake-static: $$level has no .json files"; continue; }; \
		jq -s . $$src/*.json > "web/public/data/grammar/$$level.json"; \
	done
	@cp -r server/data/corpus/vocab web/public/data/
	@cp -r server/data/corpus/kanji web/public/data/
	@echo "bake-static: web/public/data populated"

dump-grammar-examples:
	bash scripts/dump-grammar-examples.sh

build-static: bake-static dump-grammar-examples
	cd web && VITE_DEPLOY_MODE=static VITE_DEPLOY_BASE=/japanese-site/ npm run build

lint: lint-rules lint-vocab

lint-rules:
	bash scripts/lint-rules.sh
	bash scripts/lint-grammar.sh

lint-grammar:
	bash scripts/lint-grammar.sh

lint-vocab:
	bash scripts/lint-vocab.sh

backlog-render:
	node scripts/generate-backlog-md.mjs

lint-backlog-render:
	@node scripts/validate-backlog-schema.mjs
	@cp BACKLOG.md /tmp/backlog.before.md
	@node scripts/generate-backlog-md.mjs
	@if ! diff -q /tmp/backlog.before.md BACKLOG.md > /dev/null; then \
		echo "BACKLOG.md is out of date with project/backlog.yml. Run: make backlog-render"; \
		mv /tmp/backlog.before.md BACKLOG.md; \
		exit 1; \
	fi
	@rm -f /tmp/backlog.before.md

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
