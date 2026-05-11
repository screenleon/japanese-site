ALTER TABLE grammar_point ADD COLUMN pattern TEXT NOT NULL DEFAULT '[]';
ALTER TABLE grammar_point ADD COLUMN audit_status TEXT NOT NULL DEFAULT '';
ALTER TABLE grammar_point ADD COLUMN explanation_ja_blocks TEXT NOT NULL DEFAULT '[]';
-- DEFAULT 1 (not 2): legacy rows stay at v1 until the loader UPSERTs them with
-- schema_version=2 via `make seed-corpus`. The API handler refuses to serve
-- rows where schema_version != 2, so a fresh migration without a corpus reload
-- never advertises false-v2 content. See JS-102 for eventual shadow-column drop.
ALTER TABLE grammar_point ADD COLUMN schema_version INTEGER NOT NULL DEFAULT 1;
