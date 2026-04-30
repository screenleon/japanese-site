# QA Test Matrix

This matrix records the test types expected for japanese-site development.
It is a developer artifact, not learner-facing product content.

## Automated Coverage

| QA type | Coverage | Command |
|---|---|---|
| Smoke tests | Public HTTP API happy paths: health, version, vocab, kanji, sentence, grammar, quiz next/answer/stats | `cd server && go test ./internal/handlers` |
| Positive functional tests | Corpus load, deterministic ids, grading, classifier rules, question picker, stats | `cd server && go test ./...` |
| Negative functional tests | Missing query params, unknown resources, invalid JSON, missing answer fields, stale question id, no matching questions | `cd server && go test ./internal/handlers` |
| Boundary tests | Answer body size cap, exclude-list cap, deterministic id whitespace policy, malformed classifier rules | `cd server && go test ./...` |
| Security tests | Stable JSON error codes, no raw handler errors, security headers, static path traversal rejection | `cd server && go test ./internal/handlers` |
| Regression tests | Migration checksums, orphan sweep safety, LLM-generated question preservation, corpus classifier behavior | `cd server && go test ./...` |
| Data integrity tests | Migration fresh DB, checksum mismatch/backfill, next-due scheduling, stats date windows | `cd server && go test ./internal/store` |
| Build/type checks | Go vet, Go build, TypeScript build, Vite production build | `make vet && make build` |
| Rule/document lint | Layered rules structure | `make lint-rules` |
| Corpus-load acceptance | Real L1 corpus loads into SQLite after migrations | `make seed-corpus` |

## Manual / Not Yet Automated

| QA type | Current status | Recommended next automation |
|---|---|---|
| Browser UI interaction | Covered by TypeScript/Vite build only; no DOM/component runner is installed | Add Vitest + Testing Library for `QuizTab` state transitions |
| End-to-end browser flow | Manual only | Add Playwright smoke for quiz answer, stats refresh, stale-question notice |
| Visual responsive review | Manual only | Add Playwright screenshots for mobile and desktop quiz tab |
| Accessibility review | Manual only | Add axe checks once browser test runner exists |
| Performance/load | Not relevant for current single-user M3 scope | Add API benchmark/load tests if multi-user hosting starts |

## Required Before Marking a Backlog Item Done

- New behavior must have at least one positive test.
- Error-handling behavior must have at least one negative test.
- Schema or migration changes must have a migration/fresh-DB or seed-path test.
- UI changes must at minimum pass `npm run build`; add browser/component tests once the frontend test runner exists.
- Run the full validation set before PR:
  - `make lint-rules`
  - `env GOCACHE=/tmp/go-build-cache make vet`
  - `env GOCACHE=/tmp/go-build-cache make test`
  - `env GOCACHE=/tmp/go-build-cache make build`
  - `env GOCACHE=/tmp/go-build-cache make seed-corpus`
