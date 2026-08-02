package store

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type DB struct {
	*sql.DB
	// Path is the filesystem path passed to Open (empty for :memory:).
	Path string
}

// liveDBPath is the last path passed to Open in this process; used by
// migration 0022 backup preflight when JAPANESE_SITE_DB_PATH is unset.
var liveDBPath string

func Open(path string) (*DB, error) {
	// busy_timeout(5000) keeps concurrent writers (e.g., HTTP handler logging
	// an attempt while the seed CLI is running) from immediately failing
	// with SQLITE_BUSY. WAL mode plus 5s timeout has been more than enough
	// for the local single-user workload.
	dsn := path + "?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}
	if path != "" && path != ":memory:" {
		liveDBPath = path
	}
	return &DB{DB: db, Path: path}, nil
}
