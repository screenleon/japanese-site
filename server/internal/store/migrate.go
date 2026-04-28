package store

import (
	"crypto/sha256"
	"database/sql"
	"embed"
	"encoding/hex"
	"fmt"
	"io/fs"
	"log/slog"
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
