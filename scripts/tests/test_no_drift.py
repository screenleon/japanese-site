#!/usr/bin/env python3

import difflib
import filecmp
import importlib.util
import tempfile
import unittest
from pathlib import Path


ROOT = Path(__file__).resolve().parents[2]
SCRIPT = ROOT / "scripts" / "apply-allLevels-inline-ruby.py"


def load_script():
    spec = importlib.util.spec_from_file_location("apply_allLevels_inline_ruby", SCRIPT)
    module = importlib.util.module_from_spec(spec)
    assert spec.loader is not None
    spec.loader.exec_module(module)
    return module


mod = load_script()


def diff_snippet(expected_path: Path, actual_path: Path) -> str:
    expected = expected_path.read_text(encoding="utf-8").splitlines(keepends=True)
    actual = actual_path.read_text(encoding="utf-8").splitlines(keepends=True)
    return "".join(
        difflib.unified_diff(
            expected,
            actual,
            fromfile=str(expected_path),
            tofile=str(actual_path),
            n=3,
        )
    )


class NoDriftTest(unittest.TestCase):
    def test_generated_all_levels_match_real_corpus(self):
        """Regenerating the corpus from baseline produces files byte-identical to the on-disk corpus.

        Steps:
        1. Run rewrite_all into a fresh temp directory (no audit write).
        2. For every target path the script claims to produce, compare the temp output to the on-disk corpus file with shallow=False byte comparison.
        3. Assert each generated file exists and matches; on mismatch, fail the test with a unified-diff snippet for the first divergence.
        """
        with tempfile.TemporaryDirectory(prefix="js106-inline-ruby-") as tmp:
            out_root = Path(tmp)
            mod.rewrite_all(output_root=out_root, write_audit=False)
            for real_path in mod.target_paths(mod.ROOT):
                generated_path = out_root / real_path.relative_to(mod.ROOT)
                self.assertTrue(generated_path.exists(), f"missing generated file: {generated_path}")
                if not filecmp.cmp(generated_path, real_path, shallow=False):
                    self.fail(
                        "generated corpus drifted for "
                        f"{real_path.relative_to(mod.ROOT)}\n"
                        + diff_snippet(real_path, generated_path)
                    )


if __name__ == "__main__":
    unittest.main()
