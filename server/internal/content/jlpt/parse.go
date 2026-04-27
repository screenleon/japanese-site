// Package jlpt parses the jamsinclair/open-anki-jlpt-decks per-level CSV
// files. The CSV columns are: expression, reading, meaning, tags, guid.
// Some entries use only kana for both columns (no kanji form), which is fine
// — we match on (expression, reading) when overlaying onto the vocab table.
package jlpt

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

type Entry struct {
	Expression string
	Reading    string
	Meaning    string
}

func ParseCSV(path string) ([]Entry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.FieldsPerRecord = -1 // tolerate trailing fields

	header, err := r.Read()
	if err != nil {
		return nil, fmt.Errorf("read header: %w", err)
	}
	col := map[string]int{}
	for i, h := range header {
		col[h] = i
	}
	expIdx, ok := col["expression"]
	if !ok {
		return nil, fmt.Errorf("missing 'expression' column")
	}
	readIdx, ok := col["reading"]
	if !ok {
		return nil, fmt.Errorf("missing 'reading' column")
	}
	meanIdx := col["meaning"]

	var out []Entry
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read row: %w", err)
		}
		if len(row) <= readIdx {
			continue
		}
		e := Entry{
			Expression: row[expIdx],
			Reading:    row[readIdx],
		}
		if meanIdx > 0 && meanIdx < len(row) {
			e.Meaning = row[meanIdx]
		}
		if e.Expression == "" || e.Reading == "" {
			continue
		}
		out = append(out, e)
	}
	return out, nil
}
