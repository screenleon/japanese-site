package kanjidic2

import (
	"os"
	"path/filepath"
	"testing"
)

const sampleDoc = `{
  "version": "test",
  "characters": [
    {
      "literal": "山",
      "misc": {"grade": 1, "strokeCounts": [3], "jlptLevel": 4},
      "readingMeaning": {
        "groups": [
          {
            "readings": [
              {"type": "ja_on", "value": "サン"},
              {"type": "ja_kun", "value": "やま"}
            ],
            "meanings": [
              {"lang": "en", "value": "mountain"}
            ]
          }
        ]
      }
    },
    {
      "literal": "火",
      "misc": {"grade": 1, "strokeCounts": [4], "jlptLevel": 4},
      "readingMeaning": {
        "groups": [
          {
            "readings": [{"type": "ja_on", "value": "カ"}, {"type": "ja_kun", "value": "ひ"}],
            "meanings": [{"lang": "en", "value": "fire"}]
          }
        ]
      }
    }
  ]
}`

func TestParse(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sample.json")
	if err := os.WriteFile(path, []byte(sampleDoc), 0o644); err != nil {
		t.Fatal(err)
	}
	rows, err := Parse(path)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(rows) != 2 {
		t.Fatalf("want 2 rows, got %d", len(rows))
	}
	mountain := rows[0]
	if mountain.Character != "山" {
		t.Errorf("character: %q", mountain.Character)
	}
	if mountain.Onyomi != "サン" || mountain.Kunyomi != "やま" {
		t.Errorf("readings: on=%q kun=%q", mountain.Onyomi, mountain.Kunyomi)
	}
	if mountain.MeaningEN != "mountain" {
		t.Errorf("meaning: %q", mountain.MeaningEN)
	}
	if mountain.JLPTLevelOld != 4 || mountain.Grade != 1 || mountain.StrokeCount != 3 {
		t.Errorf("misc: %+v", mountain)
	}
}

func TestMapOldJLPT(t *testing.T) {
	cases := map[int]string{
		4: "N5",
		3: "N4",
		2: "N2",
		1: "N1",
		0: "",
	}
	for in, want := range cases {
		if got := MapOldJLPT(in); got != want {
			t.Errorf("MapOldJLPT(%d) = %q, want %q", in, got, want)
		}
	}
}
