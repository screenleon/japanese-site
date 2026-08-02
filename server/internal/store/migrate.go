package store

import (
	"crypto/sha256"
	"database/sql"
	"embed"
	"encoding/hex"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"sort"
	"strings"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// Migrate applies pending migrations and verifies that previously-applied
// migrations have not been edited on disk.
//
// Each migration's body is hashed (SHA-256) and the hash is stored in
// `schema_migrations.checksum`. On a DB that pre-dates the checksum column,
// the column is added on first run and any row with NULL/empty checksum is
// backfilled with the live body's hash — that body is then treated as
// canonical. A subsequent edit to an applied migration produces a hash
// mismatch and aborts startup.
//
// Concurrency: the per-migration apply path uses INSERT OR IGNORE on
// schema_migrations to claim the slot before running the body. With SQLite
// + WAL + busy_timeout, two processes calling Migrate() at the same time
// serialise on the schema_migrations PK; the loser's claim is a no-op and
// it falls through to the verify-existing path. CRLF/LF line-ending drift
// is prevented by `.gitattributes` (`*.sql text eol=lf`).
//
// Backfill trust model: see DECISIONS.md "first-boot disk body is canonical
// for migration checksum backfill" (2026-04-28). Backfill triggers only on
// rows with NULL/empty checksum; new rows always insert with a checksum.
func Migrate(db *DB) error {
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (
		version TEXT PRIMARY KEY,
		applied_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
		checksum TEXT
	)`); err != nil {
		return fmt.Errorf("create schema_migrations: %w", err)
	}

	hasChecksum, err := columnExists(db, "schema_migrations", "checksum")
	if err != nil {
		return fmt.Errorf("inspect schema_migrations: %w", err)
	}
	if !hasChecksum {
		if _, err := db.Exec(`ALTER TABLE schema_migrations ADD COLUMN checksum TEXT`); err != nil {
			// Race-safe: if another process added the column between our
			// check and our ALTER, re-check before erroring.
			has2, recheckErr := columnExists(db, "schema_migrations", "checksum")
			if recheckErr != nil {
				return fmt.Errorf("add checksum column: %w (recheck failed: %v)", err, recheckErr)
			}
			if !has2 {
				return fmt.Errorf("add checksum column: %w", err)
			}
		}
	}

	entries, err := fs.ReadDir(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("read migrations: %w", err)
	}
	files := make([]string, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			files = append(files, e.Name())
		}
	}
	sort.Strings(files)

	for _, name := range files {
		body, err := migrationsFS.ReadFile("migrations/" + name)
		if err != nil {
			return fmt.Errorf("read %s: %w", name, err)
		}
		want := checksumOf(body)

		if err := applyOrVerify(db, name, body, want); err != nil {
			return err
		}
	}
	return nil
}

// applyOrVerify performs the per-migration race-safe step: claim the row,
// and either apply (claim won) or verify the stored checksum (claim lost).
func applyOrVerify(db *DB, name string, body []byte, want string) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin %s: %w", name, err)
	}
	defer tx.Rollback() // safe: a successful Commit makes Rollback a no-op.

	res, err := tx.Exec(
		`INSERT OR IGNORE INTO schema_migrations(version, checksum) VALUES (?, ?)`,
		name, want,
	)
	if err != nil {
		return fmt.Errorf("claim %s: %w", name, err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected %s: %w", name, err)
	}

	if rowsAffected == 0 {
		// Already applied (by us in a previous run, or by a racing process).
		var stored sql.NullString
		if err := tx.QueryRow(
			`SELECT checksum FROM schema_migrations WHERE version = ?`,
			name,
		).Scan(&stored); err != nil {
			return fmt.Errorf("read stored checksum %s: %w", name, err)
		}

		if !stored.Valid || stored.String == "" {
			// Pre-checksum DB. Backfill the live body's hash as canonical.
			if _, err := tx.Exec(
				`UPDATE schema_migrations SET checksum = ? WHERE version = ?`,
				want, name,
			); err != nil {
				return fmt.Errorf("backfill %s: %w", name, err)
			}
			if err := tx.Commit(); err != nil {
				return fmt.Errorf("commit backfill %s: %w", name, err)
			}
			slog.Warn("backfilled migration checksum from disk body",
				"version", name, "checksum", want)
			return nil
		}

		if stored.String != want {
			return fmt.Errorf(
				"migration %s checksum mismatch: stored=%s embedded=%s -- applied migration files must not be edited",
				name, stored.String, want,
			)
		}
		// Match: nothing to do; deferred Rollback releases the no-op tx.
		return nil
	}

	// Won the claim. Apply the body inside the same tx so a failure rolls
	// back both the claim row and any partial DDL.
	if _, err := tx.Exec(string(body)); err != nil {
		return fmt.Errorf("apply %s: %w", name, err)
	}
	// JS-114a: require a verified pre-upgrade backup when learner history exists,
	// then rewrite deterministic question ids so a subsequent curated corpus load
	// does not orphan-sweep legacy ids and CASCADE-delete learner attempts.
	if name == "0022_grammar_slug_dedup.sql" {
		if err := requireSlugMigrationBackup(tx); err != nil {
			return fmt.Errorf("apply %s backup preflight: %w", name, err)
		}
		if n, err := rekeyGrammarQuestionIDs(tx); err != nil {
			return fmt.Errorf("apply %s question id rekey: %w", name, err)
		} else if n > 0 {
			slog.Info("rekeyed grammar question ids after slug dedup", "count", n)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit %s: %w", name, err)
	}
	// Audit point: every migration application is logged so destructive
	// migrations (e.g. 0007_question_text_id.sql, which DROPs question +
	// attempt) are visible in the operator's startup log without parsing
	// .sql comments.
	slog.Info("applied migration", "version", name, "checksum", want)
	return nil
}

// requireSlugMigrationBackup blocks 0022 on non-empty learner history unless
// an operator-supplied backup path is verified, or an explicit allow flag is set
// (tests / empty personal DBs). See DECISIONS.md JS-114a rollback plan.
func requireSlugMigrationBackup(tx *sql.Tx) error {
	var attempts int
	if err := tx.QueryRow(`SELECT COUNT(*) FROM attempt`).Scan(&attempts); err != nil {
		// attempt table may not exist on ancient DBs; treat as empty.
		if strings.Contains(err.Error(), "no such table") {
			return nil
		}
		return err
	}
	if attempts == 0 {
		return nil
	}
	if os.Getenv("JAPANESE_SITE_ALLOW_SLUG_MIGRATION") == "1" {
		slog.Warn("applying 0022 with JAPANESE_SITE_ALLOW_SLUG_MIGRATION=1",
			"attempts", attempts)
		return nil
	}
	path := os.Getenv("JAPANESE_SITE_DB_BACKUP_PATH")
	if path == "" {
		return fmt.Errorf(
			"learner attempts present (%d): set JAPANESE_SITE_DB_BACKUP_PATH to a pre-upgrade SQLite copy, or JAPANESE_SITE_ALLOW_SLUG_MIGRATION=1 to opt in",
			attempts,
		)
	}
	st, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("JAPANESE_SITE_DB_BACKUP_PATH %q: %w", path, err)
	}
	if st.Size() == 0 {
		return fmt.Errorf("JAPANESE_SITE_DB_BACKUP_PATH %q is empty", path)
	}
	slog.Info("slug migration backup verified",
		"path", path, "bytes", st.Size(), "attempts", attempts)
	return nil
}

// rekeyGrammarQuestionIDs updates question.id to corpus.QuestionID after
// grammar_point slug renames so attempt history survives seed-corpus reloads.
// Returns the number of source rows that required an id change.
func rekeyGrammarQuestionIDs(tx *sql.Tx) (int, error) {
	rows, err := tx.Query(`
		SELECT id, grammar_point, prompt, expected
		FROM question
		WHERE content_type = 'grammar'`)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	type qrow struct {
		id, gp, prompt, expected string
	}
	var list []qrow
	for rows.Next() {
		var r qrow
		if err := rows.Scan(&r.id, &r.gp, &r.prompt, &r.expected); err != nil {
			return 0, err
		}
		list = append(list, r)
	}
	if err := rows.Err(); err != nil {
		return 0, err
	}

	changed := 0
	for _, r := range list {
		want := questionID(r.gp, r.prompt, r.expected)
		if want == r.id {
			continue
		}
		var destExists int
		if err := tx.QueryRow(`SELECT COUNT(*) FROM question WHERE id = ?`, want).Scan(&destExists); err != nil {
			return changed, err
		}
		if destExists == 0 {
			// Clone row under the canonical id, then remount attempts, then drop legacy.
			if _, err := tx.Exec(`
				INSERT INTO question (
					id, kind, jlpt_level, grammar_point, prompt, expected, hint,
					source, license, validated_by, validator_score, validated_at,
					created_at, payload, content_type
				)
				SELECT
					?, kind, jlpt_level, grammar_point, prompt, expected, hint,
					source, license, validated_by, validator_score, validated_at,
					created_at, payload, content_type
				FROM question WHERE id = ?`, want, r.id); err != nil {
				return changed, fmt.Errorf("clone question %s -> %s: %w", r.id, want, err)
			}
		}
		if _, err := tx.Exec(`UPDATE attempt SET question_id = ? WHERE question_id = ?`, want, r.id); err != nil {
			return changed, fmt.Errorf("move attempts %s -> %s: %w", r.id, want, err)
		}
		if _, err := tx.Exec(`DELETE FROM question WHERE id = ?`, r.id); err != nil {
			return changed, fmt.Errorf("delete obsolete question %s: %w", r.id, err)
		}
		changed++
	}
	return changed, nil
}

// questionID mirrors corpus.QuestionID (sha256(slug|trim(prompt)|trim(expected))[:8]
// as 16 hex chars). Duplicated here to avoid an import cycle: corpus load_test
// imports store, so store must not import corpus.
func questionID(slug, prompt, expected string) string {
	h := sha256.New()
	h.Write([]byte(slug))
	h.Write([]byte("|"))
	h.Write([]byte(strings.TrimSpace(prompt)))
	h.Write([]byte("|"))
	h.Write([]byte(strings.TrimSpace(expected)))
	sum := h.Sum(nil)
	return hex.EncodeToString(sum[:8])
}

func checksumOf(body []byte) string {
	sum := sha256.Sum256(body)
	return hex.EncodeToString(sum[:])
}

func columnExists(db *DB, table, column string) (bool, error) {
	rows, err := db.Query(`SELECT 1 FROM pragma_table_info(?) WHERE name = ?`, table, column)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	return rows.Next(), nil
}
