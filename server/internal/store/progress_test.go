package store

import (
	"context"
	"path/filepath"
	"testing"
)

func TestMigrate_CreatesReadLog(t *testing.T) {
	db := newProgressTestDB(t)

	var name string
	if err := db.QueryRow(
		`SELECT name FROM sqlite_master WHERE type='table' AND name='read_log'`,
	).Scan(&name); err != nil {
		t.Fatalf("read_log table missing after Migrate: %v", err)
	}
}

func TestProgressStoreImplementations(t *testing.T) {
	tests := []struct {
		name    string
		make    func(*DB) ProgressStore
		enabled bool
	}{
		{
			name:    "sqlite",
			make:    func(db *DB) ProgressStore { return NewSQLiteProgressStore(db) },
			enabled: true,
		},
		{
			name:    "null",
			make:    func(db *DB) ProgressStore { return NullProgressStore{} },
			enabled: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			db := newProgressTestDB(t)
			seedProgressTestDB(t, db)
			ps := tt.make(db)

			if ps.Enabled() != tt.enabled {
				t.Fatalf("Enabled() = %v, want %v", ps.Enabled(), tt.enabled)
			}

			if err := ps.MarkRead(ctx, "grammar", "test-gp"); err != nil {
				t.Fatalf("first MarkRead: %v", err)
			}
			if tt.enabled {
				first := onlyReadEntry(t, ps, ctx)
				if err := ps.MarkRead(ctx, "grammar", "test-gp"); err != nil {
					t.Fatalf("second MarkRead: %v", err)
				}
				second := onlyReadEntry(t, ps, ctx)
				if second.ReadCount != 2 {
					t.Fatalf("ReadCount = %d, want 2", second.ReadCount)
				}
				if second.FirstReadAt != first.FirstReadAt {
					t.Fatalf("FirstReadAt changed: before=%q after=%q", first.FirstReadAt, second.FirstReadAt)
				}
				// LastReadAt advancement intentionally not asserted: SQLite's
				// CURRENT_TIMESTAMP has 1-second resolution, and forcing a
				// sleep to make the test pass would add wall-clock latency
				// without exercising SUT logic. ReadCount == 2 is the
				// behavioral contract of the upsert.

				s, err := ps.Progress(ctx, "N5", "grammar")
				if err != nil {
					t.Fatalf("Progress: %v", err)
				}
				if s.Read != 1 || s.Total != 2 || s.Percent != 0.5 || s.Level != "N5" || s.ContentType != "grammar" {
					t.Fatalf("Progress = %+v, want read=1 total=2 percent=0.5 level=N5 content_type=grammar", s)
				}

				// Cross-level isolation: marking an N1 grammar point must
				// not leak into N5 progress (kills a plausible mutation
				// where the level filter is dropped from the read-count
				// JOIN).
				if err := ps.MarkRead(ctx, "grammar", "n1-gp"); err != nil {
					t.Fatalf("N1 MarkRead: %v", err)
				}
				n5, err := ps.Progress(ctx, "N5", "grammar")
				if err != nil {
					t.Fatalf("Progress N5 after N1 mark: %v", err)
				}
				if n5.Read != 1 {
					t.Fatalf("N5 Read = %d after marking N1, want 1 (level isolation broken)", n5.Read)
				}
				n1, err := ps.Progress(ctx, "N1", "grammar")
				if err != nil {
					t.Fatalf("Progress N1: %v", err)
				}
				if n1.Read != 1 || n1.Total != 1 || n1.Percent != 1 {
					t.Fatalf("N1 progress = %+v, want read=1 total=1 percent=1", n1)
				}

				// total=0 percent guard: querying a level with no seeded
				// rows must return zero counts without panicking on
				// division-by-zero.
				zero, err := ps.Progress(ctx, "N3", "grammar")
				if err != nil {
					t.Fatalf("Progress N3 (no rows): %v", err)
				}
				if zero.Read != 0 || zero.Total != 0 || zero.Percent != 0 {
					t.Fatalf("N3 progress = %+v, want all zero (total=0 percent guard)", zero)
				}
				return
			}

			rows, err := ps.ListRead(ctx, "grammar")
			if err != nil {
				t.Fatalf("ListRead: %v", err)
			}
			if len(rows) != 0 {
				t.Fatalf("ListRead len = %d, want 0", len(rows))
			}
			s, err := ps.Progress(ctx, "N5", "grammar")
			if err != nil {
				t.Fatalf("Progress: %v", err)
			}
			if s.Read != 0 || s.Total != 0 || s.Percent != 0 || s.Level != "N5" || s.ContentType != "grammar" {
				t.Fatalf("Null Progress = %+v, want zero counts with level/content_type echoed", s)
			}
		})
	}
}

func onlyReadEntry(t *testing.T, ps ProgressStore, ctx context.Context) ReadEntry {
	t.Helper()
	rows, err := ps.ListRead(ctx, "grammar")
	if err != nil {
		t.Fatalf("ListRead: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("ListRead len = %d, want 1", len(rows))
	}
	return rows[0]
}

func newProgressTestDB(t *testing.T) *DB {
	t.Helper()
	db, err := Open(filepath.Join(t.TempDir(), "progress.sqlite"))
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	if err := Migrate(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

func seedProgressTestDB(t *testing.T, db *DB) {
	t.Helper()
	statements := []string{
		`INSERT INTO grammar_point (slug, title_ja, title_zh, jlpt_level, explanation_zh, source, license, validated_by)
		 VALUES ('test-gp', '〜ば', '條件形', 'N5', '測試文法說明', 'test', 'CC0', 'test-validator')`,
		`INSERT INTO grammar_point (slug, title_ja, title_zh, jlpt_level, explanation_zh, source, license, validated_by)
		 VALUES ('other-gp', '〜なら', '假定', 'N5', '測試文法說明', 'test', 'CC0', 'test-validator')`,
		`INSERT INTO grammar_point (slug, title_ja, title_zh, jlpt_level, explanation_zh, source, license, validated_by)
		 VALUES ('n1-gp', '〜いかん', '取決於', 'N1', '測試文法說明', 'test', 'CC0', 'test-validator')`,
		`INSERT INTO vocab (headword, reading, pos, gloss_en, jlpt_level, source, license, validated_by)
		 VALUES ('食べる', 'たべる', 'v1', 'eat', 'N5', 'test', 'CC0', 'test-validator')`,
		`INSERT INTO kanji (character, meaning_en, jlpt_level, source, license, validated_by)
		 VALUES ('食', 'eat', 'N5', 'test', 'CC0', 'test-validator')`,
	}
	for _, stmt := range statements {
		if _, err := db.Exec(stmt); err != nil {
			t.Fatalf("seed statement failed: %v\n%s", err, stmt)
		}
	}
}
