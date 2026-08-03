-- 0023_kokugo_progress.sql
-- Local-API progress for Track B 国語教室 (JS-132 / ADR-0005).
-- Content units remain L1 JSON under data/corpus/kokugo/**;
-- only learner progress / attempts / artifacts live in SQLite.
-- Static deploy (JS-018) does not use these tables.

CREATE TABLE kokugo_unit_progress (
    unit_key     TEXT PRIMARY KEY, -- "<stage>/<unit_id>" e.g. e5-6/library-use
    stage        TEXT NOT NULL,
    unit_id      TEXT NOT NULL,
    status       TEXT NOT NULL DEFAULT 'in_progress'
                 CHECK (status IN ('in_progress', 'completed')),
    step         TEXT NOT NULL DEFAULT 'predict',
    started_at   TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at TEXT
);

CREATE INDEX idx_kokugo_progress_stage ON kokugo_unit_progress(stage);

CREATE TABLE kokugo_task_attempt (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    unit_key      TEXT NOT NULL,
    task_id       TEXT NOT NULL,
    answer_json   TEXT NOT NULL,
    correct       INTEGER, -- 0/1; NULL when ungraded (e.g. predict)
    feedback_json TEXT,
    created_at    TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_kokugo_attempt_unit_task
    ON kokugo_task_attempt(unit_key, task_id, created_at DESC);

CREATE TABLE kokugo_artifact (
    unit_key       TEXT NOT NULL,
    revision       INTEGER NOT NULL CHECK (revision IN (0, 1)), -- 0 draft, 1 revision
    body           TEXT NOT NULL,
    checklist_json TEXT NOT NULL DEFAULT '[]',
    -- Monotonic optimistic-concurrency token (not wall-clock). Clients send
    -- expected_version; conditional UPDATE fails closed on mismatch (409).
    version        INTEGER NOT NULL DEFAULT 1,
    created_at     TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (unit_key, revision)
);
