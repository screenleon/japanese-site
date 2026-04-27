package tatoeba

import (
	"os"
	"path/filepath"
	"testing"
)

const sample = `1	jpn	おはよう。
2	eng	Hello.
3	jpn	今日は良い天気です。
4	jpn
5	jpn	寿司を食べます。
`

func TestParse(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.tsv")
	if err := os.WriteFile(path, []byte(sample), 0o644); err != nil {
		t.Fatal(err)
	}
	rows, err := Parse(path)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(rows) != 3 {
		t.Fatalf("want 3 jpn rows (empty + non-jpn skipped), got %d", len(rows))
	}
	if rows[0].TatoebaID != 1 || rows[0].TextJA != "おはよう。" {
		t.Errorf("first: %+v", rows[0])
	}
	if rows[2].TatoebaID != 5 || rows[2].TextJA != "寿司を食べます。" {
		t.Errorf("third: %+v", rows[2])
	}
}
