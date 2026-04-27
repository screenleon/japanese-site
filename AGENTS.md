# Agent Playbook — japanese-site

This repository follows the layered-rule conventions of `agent-playbook-template`.
The template's `docs/` is the source of truth for operating rules and routing.

Read these files before starting work:

0. Read `prompt-budget.yml` → `budget.profile` to determine loading depth.
   - This project ships at `standard`.
1. `../agent-playbook-template/docs/rules-quickstart.md` — Layer 1 minimal rules.
2. `../agent-playbook-template/docs/operating-rules.md` — safety, scope, validation.
   `../agent-playbook-template/docs/agent-playbook.md` — role routing.
3. `DECISIONS.md` — active architectural decisions (read for contradiction checks before planning).
4. `project/project-manifest.md` — repo-local constraints, validation commands, override registry.
5. `rules/global/*.md` and `rules/domain/*.md` — layered rules. Domain rules apply per the workspace-boundaries table in the manifest.

## Configuration layering

Constraint precedence: `rules/global/` (cross-project) → `rules/domain/` (domain-specific) → `project/project-manifest.md` (project-local).

`project/project-manifest.md` is the canonical location for repo-local constraints and validation commands. Do not duplicate them elsewhere.

## Decision log policy

`prompt-budget.yml` → `decision_log.policy: normal`. After architectural decisions, schema changes, or non-trivial tradeoffs, append an entry to `DECISIONS.md` following the existing format.

## Project at a glance

- **Goal**: A Japanese learning site covering grammar, vocabulary, and quizzes with grading + corrective feedback.
- **Tech stack**: Go (HTTP API) + SQLite/Postgres, React + Vite frontend, on-device or server-API-key LLM connector for dynamic content.
- **Content strategy**: hybrid (JMdict / KANJIDIC2 / Tatoeba for words; Tae Kim / itazuraneko for grammar; LLM-generated questions and feedback cached and validated). See `DECISIONS.md` → 2026-04-27 hybrid content sourcing.
- **Connector**: borrowed from `agent-native-pm` initially; will be extracted into a standalone `agent-connector` repo at M4. See `DECISIONS.md` → 2026-04-27 connector extraction.

## Roles and routing

See `../agent-playbook-template/docs/agent-playbook.md`. Enabled roles for this project are listed in `prompt-budget.yml` → `roles.enabled`.

Project-specific role notes:

- **content-validator** (project-local agent): every LLM-generated grammar explanation, example sentence, or quiz item must pass content-validator before being persisted to the shared cache. See `rules/domain/grading-feedback.md` and `rules/domain/jlpt-content-accuracy.md`.

## Source of truth

- `../agent-playbook-template/docs/operating-rules.md` — safety, scope, validation, review.
- `../agent-playbook-template/docs/agent-playbook.md` — role routing.
- `project/project-manifest.md` — repo-local constraints.
- `DECISIONS.md` — active decisions.
- `prompt-budget.yml` — execution mode, profile, enabled roles/skills, decision log policy.
