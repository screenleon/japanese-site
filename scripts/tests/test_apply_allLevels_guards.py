#!/usr/bin/env python3

import tempfile
import unittest
from collections import Counter
from pathlib import Path

from test_apply_allLevels_inline_ruby import mod


def baseline_entry():
    return {
        "slug": "fixture",
        "annotations": {"furigana": {"vocabulary": []}},
        "explanation_ja_blocks": [
            {"kind": "paragraph", "tokens": [{"t": "text", "v": "新聞を読む。"}]},
        ],
    }


def audit_inputs():
    return (
        {level: 1 for level in mod.LEVELS},
        Counter({("新聞", "しんぶん"): 3}),
        {level: ["fixture-slug"] for level in mod.LEVELS},
        [],
    )


def grammar_root(tmpdir):
    root = Path(tmpdir) / "server" / "data" / "corpus" / "grammar"
    for level in mod.LEVELS:
        (root / level).mkdir(parents=True)
    return root


def write_level_files(root, level, count):
    for index in range(count):
        (root / level / f"entry-{index:02d}.json").write_text("{}\n", encoding="utf-8")


class ApplyAllLevelsGuardTest(unittest.TestCase):
    def setUp(self):
        self.tmpdir = tempfile.TemporaryDirectory()
        self._orig_audit_path = mod.AUDIT_PATH
        mod.AUDIT_PATH = Path(self.tmpdir.name) / "test-audit.md"

    def tearDown(self):
        mod.AUDIT_PATH = self._orig_audit_path
        self.tmpdir.cleanup()

    def test_verify_raises_on_non_explanation_drift(self):
        """verify rejects rewrites that change fields outside explanation_ja_blocks.

        Steps:
        1. Build matching baseline and rewritten entries, then mutate a non-explanation field.
        2. Call verify on the drifted pair.
        3. Assert the non-explanation drift AssertionError branch is reported.
        """
        baseline = baseline_entry()
        rewritten = baseline_entry()
        rewritten["slug"] = "changed"

        with self.assertRaises(AssertionError) as err:
            mod.verify(Path("fixture.json"), baseline, rewritten)
        self.assertIn("field slug drifted from baseline", str(err.exception))

    def test_verify_allows_post_baseline_additions(self):
        """verify accepts rewrites that add fields absent from baseline (e.g. annotations.mental_model from JS-071).

        Steps:
        1. Build matching baseline and rewritten entries; add annotations.mental_model only to rewritten.
        2. Replace text token with a ruby token so the ruby-emitted invariant holds.
        3. Call verify; no exception.
        """
        baseline = baseline_entry()
        rewritten = baseline_entry()
        rewritten["annotations"]["mental_model"] = "「あげく」は、長く迷ったり苦労したりした流れの最後に、望ましくない結果へ行き着く感じです。"
        rewritten["explanation_ja_blocks"][0]["tokens"] = [
            {"t": "ruby", "k": "新聞", "r": "しんぶん"},
            {"t": "text", "v": "を読む。"},
        ]
        baseline["explanation_ja_blocks"][0]["tokens"] = [
            {"t": "text", "v": "新聞を読む。"},
        ]

        mod.verify(Path("fixture.json"), baseline, rewritten)

    def test_verify_rejects_baseline_annotation_drop(self):
        """verify rejects rewrites that drop or change an annotation key present in baseline.

        Steps:
        1. Build baseline with annotations.furigana populated; build rewritten that drops that key.
        2. Call verify; expect the annotations.furigana drift branch to fire.
        """
        baseline = baseline_entry()
        rewritten = baseline_entry()
        rewritten["annotations"].pop("furigana")
        rewritten["explanation_ja_blocks"][0]["tokens"] = [
            {"t": "ruby", "k": "新聞", "r": "しんぶん"},
            {"t": "text", "v": "を読む。"},
        ]
        baseline["explanation_ja_blocks"][0]["tokens"] = [
            {"t": "text", "v": "新聞を読む。"},
        ]

        with self.assertRaises(AssertionError) as err:
            mod.verify(Path("fixture.json"), baseline, rewritten)
        self.assertIn("annotations.furigana drifted from baseline", str(err.exception))

    def test_verify_raises_on_explanation_text_drift(self):
        """verify rejects rewrites that change concatenated explanation text.

        Steps:
        1. Build matching baseline and rewritten entries, then mutate explanation text.
        2. Call verify on the drifted pair.
        3. Assert the explanation text drift AssertionError branch is reported.
        """
        baseline = baseline_entry()
        rewritten = baseline_entry()
        rewritten["explanation_ja_blocks"][0]["tokens"][0]["v"] = "意味を読む。"

        with self.assertRaises(AssertionError) as err:
            mod.verify(Path("fixture.json"), baseline, rewritten)
        self.assertIn("explanation text concatenation changed", str(err.exception))

    def test_verify_raises_on_no_ruby_output(self):
        """verify rejects rewrites that emit no ruby tokens.

        Steps:
        1. Build matching baseline and rewritten entries containing only text tokens.
        2. Call verify on the pair.
        3. Assert the no-ruby AssertionError branch is reported.
        """
        baseline = baseline_entry()
        rewritten = baseline_entry()

        with self.assertRaises(AssertionError) as err:
            mod.verify(Path("fixture.json"), baseline, rewritten)
        self.assertIn("no ruby tokens emitted", str(err.exception))

    def test_audit_preserves_content_outside_markers(self):
        """audit rewrites only the generated marker section of an existing audit file.

        Steps:
        1. Create a temporary audit file with content before and after the marker pair.
        2. Run audit with deterministic generated body inputs.
        3. Assert the surrounding content is preserved and the old marker body is replaced.
        """
        before = "native review notes\n"
        after = "\nfinal verification notes\n"
        mod.AUDIT_PATH.write_text(
            before + "\n" + mod.AUDIT_START + "\nold body\n" + mod.AUDIT_END + after,
            encoding="utf-8",
        )

        mod.audit(*audit_inputs())

        result = mod.AUDIT_PATH.read_text(encoding="utf-8")
        self.assertTrue(result.startswith(before + "\n" + mod.AUDIT_START + "\n"))
        self.assertTrue(result.endswith(mod.AUDIT_END + after))
        self.assertIn("| N5 | 1 | 1 |", result)
        self.assertNotIn("old body", result)

    def test_audit_inserts_markers_when_missing(self):
        """audit prepends fallback markers when an existing audit file has no marker pair.

        Steps:
        1. Create a temporary audit file containing legacy content with no markers.
        2. Run audit with deterministic generated body inputs.
        3. Assert fallback header and markers are prepended while legacy content remains after them.
        """
        legacy = "legacy audit notes\nkeep this text\n"
        mod.AUDIT_PATH.write_text(legacy, encoding="utf-8")

        mod.audit(*audit_inputs())

        result = mod.AUDIT_PATH.read_text(encoding="utf-8")
        self.assertTrue(result.startswith("# JS-106 N5+N4+N2+N1 Inline Ruby"))
        self.assertIn(mod.AUDIT_START + "\n## Per-level Counts", result)
        self.assertIn(mod.AUDIT_END + "\n\n" + legacy, result)

    def test_audit_creates_file_when_absent(self):
        """audit creates a missing audit file with fallback header and generated markers.

        Steps:
        1. Leave the temporary audit path absent.
        2. Run audit with deterministic generated body inputs.
        3. Assert the new file contains the fallback header, marker pair, and generated body.
        """
        self.assertFalse(mod.AUDIT_PATH.exists())

        mod.audit(*audit_inputs())

        result = mod.AUDIT_PATH.read_text(encoding="utf-8")
        self.assertTrue(result.startswith("# JS-106 N5+N4+N2+N1 Inline Ruby"))
        self.assertIn(mod.AUDIT_START + "\n## Per-level Counts", result)
        self.assertIn("| 新聞 | しんぶん | 3 |", result)
        self.assertIn("\n" + mod.AUDIT_END + "\n", result)

    def test_target_paths_returns_160_when_each_level_has_40(self):
        """target_paths returns all 160 grammar files when every target level has 40 files.

        Steps:
        1. Create a temporary grammar root with 40 JSON files in each target level.
        2. Call target_paths on the temporary grammar root.
        3. Assert 160 paths are returned across the four levels.
        """
        root = grammar_root(self.tmpdir.name)
        for level in mod.LEVELS:
            write_level_files(root, level, 40)

        paths = mod.target_paths(root)

        self.assertEqual(len(paths), 160)
        self.assertEqual({path.parent.name for path in paths}, set(mod.LEVELS))

    def test_target_paths_raises_when_level_short_by_one(self):
        """target_paths rejects a level with 39 grammar files.

        Steps:
        1. Create a temporary grammar root where N4 has 39 JSON files and other levels have 40.
        2. Call target_paths on the temporary grammar root.
        3. Assert the N4 count-guard AssertionError reports 39 files.
        """
        root = grammar_root(self.tmpdir.name)
        for level in mod.LEVELS:
            write_level_files(root, level, 39 if level == "N4" else 40)

        with self.assertRaisesRegex(AssertionError, r"N4: expected 40 files, found 39"):
            mod.target_paths(root)

    def test_target_paths_raises_when_level_long_by_one(self):
        """target_paths rejects a level with 41 grammar files.

        Steps:
        1. Create a temporary grammar root where N2 has 41 JSON files and other levels have 40.
        2. Call target_paths on the temporary grammar root.
        3. Assert the N2 count-guard AssertionError reports 41 files.
        """
        root = grammar_root(self.tmpdir.name)
        for level in mod.LEVELS:
            write_level_files(root, level, 41 if level == "N2" else 40)

        with self.assertRaisesRegex(AssertionError, r"N2: expected 40 files, found 41"):
            mod.target_paths(root)

    def test_rewrite_tokens_rejects_term_token(self):
        """rewrite_tokens rejects term tokens as out of scope for JS-106.

        Steps:
        1. Build a minimal token list containing one term token.
        2. Call rewrite_tokens with empty readings.
        3. Assert the term-token AssertionError branch is reported.
        """
        tokens = [{"t": "term", "v": "新聞"}]

        with self.assertRaisesRegex(AssertionError, "term tokens are out of scope for JS-106"):
            mod.rewrite_tokens(tokens, {})

    def test_rewrite_blocks_rejects_unexpected_kind(self):
        """rewrite_blocks rejects explanation blocks with an unexpected kind.

        Steps:
        1. Build a minimal entry whose explanation_ja_blocks contains a header block.
        2. Call rewrite_blocks on the entry.
        3. Assert the unexpected-kind AssertionError branch is reported.
        """
        entry = {
            "annotations": {"furigana": {"vocabulary": []}},
            "explanation_ja_blocks": [
                {"kind": "header", "tokens": [{"t": "text", "v": "新聞"}]},
            ],
        }

        with self.assertRaisesRegex(AssertionError, "unexpected block kind header"):
            mod.rewrite_blocks(entry)


if __name__ == "__main__":
    unittest.main()
