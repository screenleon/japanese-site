package store

import (
	"context"
	"database/sql"
	"errors"
)

type Kanji struct {
	ID           int64  `json:"id"`
	Character    string `json:"character"`
	Onyomi       string `json:"onyomi,omitempty"`
	Kunyomi      string `json:"kunyomi,omitempty"`
	MeaningEN    string `json:"meaning_en,omitempty"`
	JLPTLevel    string `json:"jlpt_level,omitempty"`
	Grade        *int64 `json:"grade,omitempty"`
	StrokeCount  *int64 `json:"stroke_count,omitempty"`
	Source       string `json:"source"`
	License      string `json:"license"`
	ValidatedBy  string `json:"validated_by,omitempty"`
}

var ErrKanjiNotFound = errors.New("kanji not found")

func GetKanji(ctx context.Context, db *DB, character string) (Kanji, error) {
	row := db.QueryRowContext(ctx, `
		SELECT id, character,
		       COALESCE(onyomi, ''), COALESCE(kunyomi, ''), COALESCE(meaning_en, ''),
		       COALESCE(jlpt_level, ''), grade, stroke_count,
		       source, license, COALESCE(validated_by, '')
		FROM kanji WHERE character = ?`, character)
	var k Kanji
	var grade, stroke sql.NullInt64
	if err := row.Scan(&k.ID, &k.Character, &k.Onyomi, &k.Kunyomi, &k.MeaningEN,
		&k.JLPTLevel, &grade, &stroke, &k.Source, &k.License, &k.ValidatedBy); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Kanji{}, ErrKanjiNotFound
		}
		return Kanji{}, err
	}
	if grade.Valid {
		v := grade.Int64
		k.Grade = &v
	}
	if stroke.Valid {
		v := stroke.Int64
		k.StrokeCount = &v
	}
	return k, nil
}
