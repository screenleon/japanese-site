package store

import (
	"context"
	"path/filepath"
	"testing"
	"time"
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
				time.Sleep(1100 * time.Millisecond)
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
				if second.LastReadAt <= first.LastReadAt {
					t.Fatalf("LastReadAt did not advance: before=%q after=%q", first.LastReadAt, second.LastReadAt)
				}

				s, err := ps.Progress(ctx, "N5", "grammar")
				if err != nil {
					t.Fatalf("Progress: %v", err)
				}
				if s.Read != 1 || s.Total != 2 || s.Percent != 0.5 || s.Level != "N5" || s.ContentType != "grammar" {
					t.Fatalf("Progress = %+v, want read=1 total=2 percent=0.5 level=N5 content_type=grammar", s)
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
