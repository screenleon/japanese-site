package jmdict

import (
	"os"
	"path/filepath"
	"testing"
)

const sampleDoc = `{
  "version": "test",
  "words": [
    {
      "id": "100001",
      "kanji": [{"text": "食べる", "common": true}],
      "kana":  [{"text": "たべる", "common": true}],
      "sense": [
        {"partOfSpeech": ["v1", "vt"],
         "gloss": [
           {"lang": "eng", "text": "to eat"},
           {"lang": "eng", "text": "to live on"}
         ]}
      ]
    },
    {
      "id": "100002",
      "kanji": [],
      "kana":  [{"text": "おはよう", "common": true}],
      "sense": [
        {"partOfSpeech": ["int"],
         "gloss": [{"lang": "eng", "text": "good morning"}]}
      ]
    },
    {
      "id": "100003",
      "kanji": [{"text": "稀", "common": false}, {"text": "希", "common": true}],
      "kana":  [{"text": "まれ", "common": true}],
      "sense": [
        {"partOfSpeech": ["adj-na", "n"],
         "gloss": [{"lang": "eng", "text": "rare"}]}
      ]
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
	if len(rows) != 3 {
		t.Fatalf("want 3 rows, got %d", len(rows))
	}

	taberu := rows[0]
	if taberu.Headword != "食べる" || taberu.Reading != "たべる" {
		t.Errorf("taberu wrong: %+v", taberu)
	}
	if taberu.POS != "v1,vt" {
		t.Errorf("taberu POS: %q", taberu.POS)
	}
	if taberu.GlossEN != "to eat; to live on" {
		t.Errorf("taberu gloss: %q", taberu.GlossEN)
	}
	if !taberu.Common {
		t.Errorf("taberu should be common")
	}

	ohayou := rows[1]
	if ohayou.Headword != "おはよう" || ohayou.Reading != "おはよう" {
		t.Errorf("kana-only entry: %+v", ohayou)
	}

	mare := rows[2]
	if mare.Headword != "希" {
		t.Errorf("mare should pick first common kanji form, got %q", mare.Headword)
	}
}
