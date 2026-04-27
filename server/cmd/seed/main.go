package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/screenleon/japanese-site/server/internal/config"
	"github.com/screenleon/japanese-site/server/internal/content/corpus"
	"github.com/screenleon/japanese-site/server/internal/content/external"
	"github.com/screenleon/japanese-site/server/internal/content/jlpt"
	"github.com/screenleon/japanese-site/server/internal/content/jmdict"
	"github.com/screenleon/japanese-site/server/internal/content/kanjidic2"
	"github.com/screenleon/japanese-site/server/internal/content/tatoeba"
	"github.com/screenleon/japanese-site/server/internal/jlptderive"
	"github.com/screenleon/japanese-site/server/internal/store"
)

// Defaults assume the seed CLI is invoked from the server/ directory.
const (
	defaultLockPath = "data/external.lock"
	defaultCacheDir = "data/external"
)

// jlptOrder is N5 → N1 so the easier level wins when a word straddles two.
var jlptOrder = []string{"N5", "N4", "N3", "N2", "N1"}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "jmdict":
		run(args, runJMdict)
	case "kanjidic2":
		run(args, runKanjidic2)
	case "jlpt":
		run(args, runJLPT)
	case "tatoeba":
		run(args, runTatoeba)
	case "tag-sentences":
		run(args, runTagSentences)
	case "backfill-kanji":
		run(args, runBackfillKanji)
	case "derive":
		run(args, runDerive)
	case "corpus":
		run(args, runCorpus)
	case "all":
		run(args, runAll)
	case "-h", "--help", "help":
		usage()
	default:
		fmt.Fprintf(os.Stderr, "unknown subcommand: %s\n\n", cmd)
		usage()
		os.Exit(2)
	}
}

func run(args []string, fn func(opts) error) {
	o := parseFlags(args)
	if err := fn(o); err != nil {
		slog.Error("seed failed", "err", err)
		os.Exit(1)
	}
}

type opts struct {
	lockPath    string
	cacheDir    string
	acceptFresh bool
}

func parseFlags(args []string) opts {
	fs := flag.NewFlagSet("seed", flag.ExitOnError)
	o := opts{}
	fs.StringVar(&o.lockPath, "lock", defaultLockPath, "path to external.lock")
	fs.StringVar(&o.cacheDir, "cache-dir", defaultCacheDir, "where downloaded files are stored")
	fs.BoolVar(&o.acceptFresh, "accept-fresh-checksum", false, "compute and write back missing sha256 on first run")
	_ = fs.Parse(args)
	return o
}

func usage() {
	fmt.Print(`Usage: seed <subcommand> [flags]

Subcommands:
  jmdict        Import jmdict-simplified into the vocab table.
  kanjidic2     Import kanjidic2-simplified into the kanji table.
  jlpt          Apply JLPT-level overlay to vocab rows (N5 → N1 order).
  tatoeba       Import Tatoeba Japanese sentences into the sentence table.
  tag-sentences Tokenize sentences and assign jlpt_level from matched vocab.
  backfill-kanji Re-derive kanji jlpt_level from the easiest vocab word using it.
  derive        Run tag-sentences then backfill-kanji.
  corpus        Load curated L1 corpus (grammar points, examples, questions).
  all           Run jmdict → kanjidic2 → jlpt → tatoeba → derive → corpus.

Flags:
  --lock                  path to external.lock (default: data/external.lock)
  --cache-dir             where downloaded files are stored (default: data/external)
  --accept-fresh-checksum compute and write back missing sha256 on first run
`)
}

// ── helpers ─────────────────────────────────────────────────────────────────

func ensureDB(o opts) (*store.DB, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	db, err := store.Open(cfg.DBPath)
	if err != nil {
		return nil, err
	}
	if err := store.Migrate(db); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func ensureLocal(o opts, ds external.Dataset, manifest external.Manifest) (string, external.Manifest, error) {
	path, fresh, err := external.EnsureLocal(o.cacheDir, ds, o.acceptFresh)
	if err != nil {
		return "", manifest, err
	}
	if fresh != "" {
		slog.Info("computed fresh sha256, writing back to lock", "name", ds.Name, "sha256", fresh)
		for i := range manifest.Datasets {
			if manifest.Datasets[i].Name == ds.Name {
				manifest.Datasets[i].SHA256 = fresh
			}
		}
		if err := external.Save(o.lockPath, manifest); err != nil {
			return "", manifest, fmt.Errorf("write lock: %w", err)
		}
	}
	return path, manifest, nil
}

// ── subcommands ─────────────────────────────────────────────────────────────

func runJMdict(o opts) error {
	manifest, err := external.Load(o.lockPath)
	if err != nil {
		return err
	}
	ds, ok := manifest.Find("jmdict-simplified-eng-common")
	if !ok {
		return fmt.Errorf("dataset 'jmdict-simplified-eng-common' not in lock")
	}
	slog.Info("ensuring local copy", "name", ds.Name)
	path, manifest, err := ensureLocal(o, ds, manifest)
	if err != nil {
		return err
	}
	slog.Info("parsing jmdict")
	rows, err := jmdict.Parse(path)
	if err != nil {
		return err
	}
	slog.Info("parsed", "rows", len(rows))

	db, err := ensureDB(o)
	if err != nil {
		return err
	}
	defer db.Close()

	stats, err := jmdict.Import(context.Background(), db.DB, ds.License, rows)
	if err != nil {
		return err
	}
	slog.Info("jmdict import done", "inserted", stats.Inserted, "skipped", stats.Skipped)
	return nil
}

func runKanjidic2(o opts) error {
	manifest, err := external.Load(o.lockPath)
	if err != nil {
		return err
	}
	ds, ok := manifest.Find("kanjidic2-simplified-en")
	if !ok {
		return fmt.Errorf("dataset 'kanjidic2-simplified-en' not in lock")
	}
	slog.Info("ensuring local copy", "name", ds.Name)
	path, manifest, err := ensureLocal(o, ds, manifest)
	if err != nil {
		return err
	}
	slog.Info("parsing kanjidic2")
	rows, err := kanjidic2.Parse(path)
	if err != nil {
		return err
	}
	slog.Info("parsed", "rows", len(rows))

	db, err := ensureDB(o)
	if err != nil {
		return err
	}
	defer db.Close()

	stats, err := kanjidic2.Import(context.Background(), db.DB, ds.License, rows)
	if err != nil {
		return err
	}
	slog.Info("kanjidic2 import done", "inserted", stats.Inserted, "skipped", stats.Skipped)
	return nil
}

func runJLPT(o opts) error {
	manifest, err := external.Load(o.lockPath)
	if err != nil {
		return err
	}
	overlays := manifest.FilterBySourceID(jlpt.SourceID)
	if len(overlays) == 0 {
		return fmt.Errorf("no datasets with source_id=%s in lock", jlpt.SourceID)
	}
	byLevel := map[string]external.Dataset{}
	for _, d := range overlays {
		byLevel[d.Level] = d
	}

	db, err := ensureDB(o)
	if err != nil {
		return err
	}
	defer db.Close()

	for _, level := range jlptOrder {
		ds, ok := byLevel[level]
		if !ok {
			slog.Warn("no overlay dataset for level", "level", level)
			continue
		}
		slog.Info("applying overlay", "level", level, "name", ds.Name)
		path, m2, err := ensureLocal(o, ds, manifest)
		if err != nil {
			return err
		}
		manifest = m2
		entries, err := jlpt.ParseCSV(path)
		if err != nil {
			return fmt.Errorf("parse %s: %w", ds.Name, err)
		}
		stats, err := jlpt.Overlay(context.Background(), db.DB, level, entries)
		if err != nil {
			return err
		}
		slog.Info("overlay done", "level", level, "matched", stats.Matched, "unmatched", stats.Unmatched)
	}
	return nil
}

func runTatoeba(o opts) error {
	manifest, err := external.Load(o.lockPath)
	if err != nil {
		return err
	}
	ds, ok := manifest.Find("tatoeba-jpn-sentences")
	if !ok {
		return fmt.Errorf("dataset 'tatoeba-jpn-sentences' not in lock")
	}
	slog.Info("ensuring local copy", "name", ds.Name)
	path, manifest, err := ensureLocal(o, ds, manifest)
	if err != nil {
		return err
	}
	slog.Info("parsing tatoeba")
	rows, err := tatoeba.Parse(path)
	if err != nil {
		return err
	}
	slog.Info("parsed", "rows", len(rows))

	db, err := ensureDB(o)
	if err != nil {
		return err
	}
	defer db.Close()

	stats, err := tatoeba.Import(context.Background(), db.DB, ds.License, rows)
	if err != nil {
		return err
	}
	slog.Info("tatoeba import done", "inserted", stats.Inserted, "skipped", stats.Skipped)
	return nil
}

func runAll(o opts) error {
	for _, fn := range []func(opts) error{runJMdict, runKanjidic2, runJLPT, runTatoeba, runTagSentences, runBackfillKanji, runCorpus} {
		if err := fn(o); err != nil {
			return err
		}
	}
	return nil
}

func runCorpus(o opts) error {
	db, err := ensureDB(o)
	if err != nil {
		return err
	}
	defer db.Close()
	slog.Info("loading L1 corpus")
	stats, err := corpus.Load(context.Background(), db.DB, "data/corpus")
	if err != nil {
		return err
	}
	slog.Info("corpus load done",
		"grammar_points", stats.GrammarPoints,
		"examples", stats.GrammarExamples,
		"questions", stats.Questions)
	return nil
}

func runTagSentences(o opts) error {
	db, err := ensureDB(o)
	if err != nil {
		return err
	}
	defer db.Close()
	slog.Info("tagging sentences with JLPT")
	stats, err := jlptderive.TagSentences(context.Background(), db.DB, 0)
	if err != nil {
		return err
	}
	slog.Info("tag-sentences done", "tagged", stats.Tagged, "untaggable", stats.Untaggable)
	return nil
}

func runBackfillKanji(o opts) error {
	db, err := ensureDB(o)
	if err != nil {
		return err
	}
	defer db.Close()
	slog.Info("backfilling kanji JLPT from vocab")
	stats, err := jlptderive.BackfillKanjiFromVocab(context.Background(), db.DB)
	if err != nil {
		return err
	}
	slog.Info("backfill-kanji done", "updated", stats.Updated, "unchanged", stats.UnchangedFromOld)
	return nil
}

func runDerive(o opts) error {
	if err := runTagSentences(o); err != nil {
		return err
	}
	return runBackfillKanji(o)
}
