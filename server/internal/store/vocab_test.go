package store

import (
	"context"
	"testing"
)

func TestCountVocab(t *testing.T) {
	db := newMemoryDB(t)
	defer db.Close()
	if err := Migrate(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	if _, err := db.Exec(`INSERT INTO vocab (headword, reading, pos, jlpt_level, source, license)
	                      VALUES ('テスト', 'てすと', 'n', 'N5', 'test', 'CC0')`); err != nil {
		t.Fatalf("seed vocab: %v", err)
	}

	got, err := CountVocab(context.Background(), db, "N5")
	if err != nil {
		t.Fatalf("CountVocab N5: %v", err)
	}
	if got != 1 {
		t.Fatalf("CountVocab N5 = %d, want 1", got)
	}

	got, err = CountVocab(context.Background(), db, "")
	if err != nil {
		t.Fatalf("CountVocab total: %v", err)
	}
	if got != 1 {
		t.Fatalf("CountVocab total = %d, want 1", got)
	}

	got, err = CountVocab(context.Background(), db, "N4")
	if err != nil {
		t.Fatalf("CountVocab N4: %v", err)
	}
	if got != 0 {
		t.Fatalf("CountVocab N4 = %d, want 0", got)
	}
}
