package config

import (
	"os"
	"testing"
)

func TestLoadKokugoDirDefaultAndOverride(t *testing.T) {
	t.Setenv("JS_PROGRESS_STORE", "sqlite")
	t.Setenv("LISTEN_ADDR", "")
	t.Setenv("DB_PATH", "")
	t.Setenv("STATIC_DIR", "")

	// Unset KOKUGO_DIR → default corpus path (enabled).
	// t.Setenv cannot express "unset"; clear after a temporary set.
	t.Setenv("KOKUGO_DIR", "tmp")
	if err := os.Unsetenv("KOKUGO_DIR"); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.KokugoDir != "data/corpus/kokugo" {
		t.Fatalf("default KokugoDir=%q", cfg.KokugoDir)
	}

	t.Setenv("KOKUGO_DIR", "/tmp/custom-kokugo")
	cfg2, err := Load()
	if err != nil {
		t.Fatalf("Load override: %v", err)
	}
	if cfg2.KokugoDir != "/tmp/custom-kokugo" {
		t.Fatalf("override KokugoDir=%q", cfg2.KokugoDir)
	}

	// Explicit empty disables Kokugo (distinct from unset).
	t.Setenv("KOKUGO_DIR", "")
	cfg3, err := Load()
	if err != nil {
		t.Fatalf("Load empty: %v", err)
	}
	if cfg3.KokugoDir != "" {
		t.Fatalf("explicit empty KokugoDir=%q want \"\"", cfg3.KokugoDir)
	}
}

func TestLoadProgressModeFallback(t *testing.T) {
	t.Setenv("JS_PROGRESS_STORE", "bogus")
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.ProgressMode != "null" {
		t.Fatalf("ProgressMode=%q", cfg.ProgressMode)
	}
}
