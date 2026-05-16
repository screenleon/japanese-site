#!/usr/bin/env python3

import importlib.util
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


def entry_with_vocab(*items):
    return {
        "annotations": {
            "furigana": {
                "vocabulary": [
                    {"kanji": kanji, "reading": reading}
                    for kanji, reading in items
                ]
            }
        }
    }


def annotate(value, *items):
    return mod.annotate_text(value, mod.entry_readings(entry_with_vocab(*items)))


def concat(tokens):
    return "".join(token.get("v", "") if token["t"] == "text" else token["k"] for token in tokens)


class ApplyAllLevelsInlineRubyTest(unittest.TestCase):
    def assertRuby(self, tokens, kanji, reading, count=1):
        matches = [token for token in tokens if token == {"t": "ruby", "k": kanji, "r": reading}]
        self.assertEqual(len(matches), count, tokens)

    def assertNoRuby(self, tokens, kanji):
        self.assertNotIn(kanji, [token.get("k") for token in tokens if token["t"] == "ruby"], tokens)

    def test_compound_kofu(self):
        tokens = annotate("古風で硬い響きです。", ("古", "ふる"), ("風", "かぜ"))
        self.assertRuby(tokens, "古風", "こふう")
        self.assertNoRuby(tokens, "風")

    def test_compound_shunkan(self):
        tokens = annotate("瞬間、次の瞬間、その瞬間。", ("瞬", "またた"), ("間", "あいだ"))
        self.assertRuby(tokens, "瞬間", "しゅんかん", count=3)
        self.assertNoRuby(tokens, "間")

    def test_compound_jinsoku(self):
        tokens = annotate("迅速な対応。", ("迅", "じん"), ("速", "はや"))
        self.assertRuby(tokens, "迅速", "じんそく")
        self.assertNoRuby(tokens, "速")

    def test_compound_shinbun(self):
        tokens = annotate("新聞を読む。", ("新", "き"), ("聞", "き"))
        self.assertRuby(tokens, "新聞", "しんぶん")
        self.assertNoRuby(tokens, "聞")

    def test_compound_hokokusho(self):
        tokens = annotate("報告書と報告書。", ("報告", "ほうこく"), ("書", "か"))
        self.assertRuby(tokens, "報告書", "ほうこくしょ", count=2)
        self.assertNoRuby(tokens, "書")

    def test_jippun_in_time_context(self):
        tokens = annotate("三時十分に始まります。", ("三", "さん"), ("時", "じ"), ("十分", "じゅうぶん"))
        self.assertEqual(
            tokens[:3],
            [
                {"t": "ruby", "k": "三", "r": "さん"},
                {"t": "ruby", "k": "時", "r": "じ"},
                {"t": "ruby", "k": "十分", "r": "じっぷん"},
            ],
        )

    def test_juppun_generic_context(self):
        tokens = annotate("準備は十分です。", ("十分", "じゅうぶん"))
        self.assertRuby(tokens, "十分", "じゅうぶん")
        self.assertNotIn({"t": "ruby", "k": "十分", "r": "じっぷん"}, tokens)

    def test_example_bullet_protection(self):
        value = "形：\n- 動詞辞書形＋こと\n\n例：\n- 新聞を読む\n- 報告書を書く"
        tokens = annotate(value, ("動詞辞書形", "どうしじしょけい"), ("新聞", "しんぶん"), ("報告書", "ほうこくしょ"))
        self.assertRuby(tokens, "動詞辞書形", "どうしじしょけい")
        self.assertNoRuby(tokens, "新聞")
        self.assertNoRuby(tokens, "報告書")
        self.assertEqual(concat(tokens), value)

    def test_rei_colon_range_protection(self):
        value = "新聞で使う。\n\n例：新聞を読む\n報告書を書く\n\n意味を確認する。"
        tokens = annotate(value, ("新聞", "しんぶん"), ("報告書", "ほうこくしょ"), ("意味", "いみ"))
        self.assertRuby(tokens, "新聞", "しんぶん")
        self.assertNoRuby(tokens, "報告書")
        self.assertRuby(tokens, "意味", "いみ")
        self.assertEqual(concat(tokens), value)

    def test_mixed_paragraph_and_list_block(self):
        entry = {
            "annotations": {"furigana": {"vocabulary": []}},
            "explanation_ja_blocks": [
                {"kind": "paragraph", "tokens": [{"t": "text", "v": "瞬間を表す。"}]},
                {
                    "kind": "list",
                    "items": [
                        {"tokens": [{"t": "text", "v": "新聞に多い。"}]},
                        {"tokens": [{"t": "text", "v": "- 報告書を書く"}]},
                    ],
                },
            ],
        }
        blocks = mod.rewrite_blocks(entry)
        self.assertRuby(blocks[0]["tokens"], "瞬間", "しゅんかん")
        self.assertRuby(blocks[1]["items"][0]["tokens"], "新聞", "しんぶん")
        self.assertNoRuby(blocks[1]["items"][1]["tokens"], "報告書")

    def test_byte_identity_inline_baseline(self):
        baseline = "三時十分に新聞を読む。\n- 報告書を書く"
        tokens = annotate(baseline, ("三", "さん"), ("時", "じ"), ("十分", "じゅうぶん"), ("新聞", "しんぶん"), ("報告書", "ほうこくしょ"))
        self.assertEqual(concat(tokens), baseline)
        self.assertEqual(
            tokens,
            [
                {"t": "ruby", "k": "三", "r": "さん"},
                {"t": "ruby", "k": "時", "r": "じ"},
                {"t": "ruby", "k": "十分", "r": "じっぷん"},
                {"t": "text", "v": "に"},
                {"t": "ruby", "k": "新聞", "r": "しんぶん"},
                {"t": "text", "v": "を読む。\n- 報告書を書く"},
            ],
        )


if __name__ == "__main__":
    unittest.main()
