package jlpt

import (
	"os"
	"path/filepath"
	"testing"
)

const sampleCSV = `expression,reading,meaning,tags,guid
ああ,ああ,"Ah!, Oh!",JLPT JLPT_5 JLPT_N5,abc
会う,あう,"to meet, to see",JLPT JLPT_5 JLPT_N5,def
青い,あおい,blue,JLPT JLPT_5 JLPT_N5,ghi
,reading-only,no expression,JLPT,jkl
no-reading,,no reading,JLPT,mno
`

func TestParseCSV(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "n5.csv")
	if err := os.WriteFile(path, []byte(sampleCSV), 0o644); err != nil {
		t.Fatal(err)
	}
	entries, err := ParseCSV(path)
	if err != nil {
		t.Fatalf("ParseCSV: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("want 3 entries (empty rows skipped), got %d", len(entries))
	}
	if entries[0].Expression != "ああ" || entries[0].Reading != "ああ" {
		t.Errorf("first: %+v", entries[0])
	}
	if entries[1].Expression != "会う" || entries[1].Reading != "あう" {
		t.Errorf("second: %+v", entries[1])
	}
}
