package store

import (
	"strings"
	"testing"
)

func TestMigrate_FreshDBHasAllTables(t *testing.T) {
	db := newMemoryDB(t)
	defer db.Close()

	if err := Migrate(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	expected := []string{
		"schema_migrations",
		"vocab",
		"kanji",
		"sentence",
		"grammar_point",
		"grammar_example",
		"question",
		"feedback_template",
		"attempt",
	}
	for _, tbl := range expected {
		var name string
		err := db.QueryRow(
			`SELECT name FROM sqlite_master WHERE type='table' AND name=?`,
			tbl,
		).Scan(&name)
		if err != nil {
			t.Errorf("table %s missing after Migrate: %v", tbl, err)
		}
	}
}

func TestMigrate_RecordsChecksums(t *testing.T) {
	db := newMemoryDB(t)
	defer db.Close()

	if err := Migrate(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	rows, err := db.Query(`SELECT version, COALESCE(checksum, '') FROM schema_migrations`)
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	defer rows.Close()
	count := 0
	for rows.Next() {
		var v, c string
		if err := rows.Scan(&v, &c); err != nil {
			t.Fatalf("scan: %v", err)
		}
		if len(c) != 64 {
			t.Errorf("%s: expected 64-char SHA-256 hex, got %q (len=%d)", v, c, len(c))
		}
		count++
	}
	if count == 0 {
		t.Fatal("no migration rows recorded")
	}
}

func TestMigrate_ChecksumMismatchFails(t *testing.T) {
	db := newMemoryDB(t)
	defer db.Close()

	if err := Migrate(db); err != nil {
		t.Fatalf("first migrate: %v", err)
	}

	// Simulate a tampered-with applied migration by replacing the stored
	// checksum with a fake value. A real edit to the embedded .sql file
	// is detected the same way: the freshly computed hash no longer
	// matches the stored one.
	tampered := strings.Repeat("ab", 32)
	if _, err := db.Exec(
		`UPDATE schema_migrations SET checksum = ? WHERE version = ?`,
		tampered, "0001_init.sql",
	); err != nil {
		t.Fatalf("tamper: %v", err)
	}

	err := Migrate(db)
	if err == nil {
		t.Fatal("expected checksum mismatch error, got nil")
	}
	if !strings.Contains(err.Error(), "checksum mismatch") {
		t.Errorf("expected error mentioning 'checksum mismatch', got: %v", err)
	}
}

func TestMigrate_BackfillsLegacyNullChecksum(t *testing.T) {
	db := newMemoryDB(t)
	defer db.Close()

	if err := Migrate(db); err != nil {
		t.Fatalf("first migrate: %v", err)
	}

	// Pre-checksum DBs stored rows with NULL in this column. Simulate.
	if _, err := db.Exec(
		`UPDATE schema_migrations SET checksum = NULL WHERE version = ?`,
		"0001_init.sql",
	); err != nil {
		t.Fatalf("null out: %v", err)
	}

	if err := Migrate(db); err != nil {
		t.Fatalf("backfill migrate: %v", err)
	}

	var sum string
	if err := db.QueryRow(
		`SELECT COALESCE(checksum, '') FROM schema_migrations WHERE version = ?`,
		"0001_init.sql",
	).Scan(&sum); err != nil {
		t.Fatalf("read: %v", err)
	}
	if len(sum) != 64 {
		t.Fatalf("checksum was not backfilled: %q (len=%d)", sum, len(sum))
	}

	// Idempotency: a third call must succeed (the just-backfilled hash is
	// canonical, so verification passes).
	if err := Migrate(db); err != nil {
		t.Fatalf("third migrate: %v", err)
	}
}

func newMemoryDB(t *testing.T) *DB {
	t.Helper()
	db, err := Open(":memory:")
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	return db
}
