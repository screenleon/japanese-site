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
        """古風 emits a single こふう ruby instead of splitting into 古 + 風(かぜ).

        Steps:
        1. Build a vocabulary that biases the tokenizer toward single-kanji 風→かぜ.
        2. Annotate prose containing 古風.
        3. Assert one ruby{古風→こふう} token and no standalone 風 ruby.
        """
        tokens = annotate("古風で硬い響きです。", ("古", "ふる"), ("風", "かぜ"))
        self.assertRuby(tokens, "古風", "こふう")
        self.assertNoRuby(tokens, "風")

    def test_compound_shunkan(self):
        """瞬間 emits しゅんかん at every occurrence instead of splitting into 瞬 + 間(あいだ).

        Steps:
        1. Build a vocabulary that biases the tokenizer toward single-kanji 間→あいだ.
        2. Annotate prose containing 瞬間 three times.
        3. Assert three ruby{瞬間→しゅんかん} tokens and no standalone 間 ruby.
        """
        tokens = annotate("瞬間、次の瞬間、その瞬間。", ("瞬", "またた"), ("間", "あいだ"))
        self.assertRuby(tokens, "瞬間", "しゅんかん", count=3)
        self.assertNoRuby(tokens, "間")

    def test_compound_jinsoku(self):
        """迅速 emits じんそく instead of splitting into 迅 + 速(はや).

        Steps:
        1. Build a vocabulary that biases the tokenizer toward single-kanji 速→はや.
        2. Annotate prose containing 迅速.
        3. Assert one ruby{迅速→じんそく} token and no standalone 速 ruby.
        """
        tokens = annotate("迅速な対応。", ("迅", "じん"), ("速", "はや"))
        self.assertRuby(tokens, "迅速", "じんそく")
        self.assertNoRuby(tokens, "速")

    def test_compound_shinbun(self):
        """新聞 emits しんぶん instead of splitting into 新 + 聞(き).

        Steps:
        1. Build a vocabulary that biases the tokenizer toward single-kanji 聞→き.
        2. Annotate prose containing 新聞.
        3. Assert one ruby{新聞→しんぶん} token and no standalone 聞 ruby.
        """
        tokens = annotate("新聞を読む。", ("新", "き"), ("聞", "き"))
        self.assertRuby(tokens, "新聞", "しんぶん")
        self.assertNoRuby(tokens, "聞")

    def test_compound_hokokusho(self):
        """報告書 emits ほうこくしょ at every occurrence instead of splitting into 報告 + 書(か).

        Steps:
        1. Build a vocabulary that biases the tokenizer toward sub-compound 報告 and single-kanji 書→か.
        2. Annotate prose containing 報告書 twice.
        3. Assert two ruby{報告書→ほうこくしょ} tokens and no standalone 書 ruby.
        """
        tokens = annotate("報告書と報告書。", ("報告", "ほうこく"), ("書", "か"))
        self.assertRuby(tokens, "報告書", "ほうこくしょ", count=2)
        self.assertNoRuby(tokens, "書")

    def test_jippun_in_time_context(self):
        """十分 in a 時+digit context emits じっぷん (minutes), overriding generic じゅうぶん.

        Steps:
        1. Build a vocabulary whose default 十分 reading is じゅうぶん.
        2. Annotate 三時十分に始まります。 (time-of-day context).
        3. Assert the leading three tokens are 三/時/十分→じっぷん in document order.
        """
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
        """十分 outside a time-of-day context keeps the default じゅうぶん reading.

        Steps:
        1. Build a vocabulary whose default 十分 reading is じゅうぶん.
        2. Annotate generic-context prose 準備は十分です。
        3. Assert ruby{十分→じゅうぶん} is emitted and no ruby{十分→じっぷん} appears.
        """
        tokens = annotate("準備は十分です。", ("十分", "じゅうぶん"))
        self.assertRuby(tokens, "十分", "じゅうぶん")
        self.assertNotIn({"t": "ruby", "k": "十分", "r": "じっぷん"}, tokens)

    def test_example_bullet_protection(self):
        """Markdown example bullets after 例： remain a single text token and are not annotated.

        Steps:
        1. Build prose with a metalanguage section, a 例： header, and example bullets containing 新聞 / 報告書.
        2. Annotate the prose with vocabulary covering those compounds.
        3. Assert 動詞辞書形 in metalanguage IS ruby-tagged, but 新聞 / 報告書 inside example bullets are NOT, and the concat round-trips byte-for-byte.
        """
        value = "形：\n- 動詞辞書形＋こと\n\n例：\n- 新聞を読む\n- 報告書を書く"
        tokens = annotate(value, ("動詞辞書形", "どうしじしょけい"), ("新聞", "しんぶん"), ("報告書", "ほうこくしょ"))
        self.assertRuby(tokens, "動詞辞書形", "どうしじしょけい")
        self.assertNoRuby(tokens, "新聞")
        self.assertNoRuby(tokens, "報告書")
        self.assertEqual(concat(tokens), value)

    def test_rei_colon_range_protection(self):
        """The 例： example range is preserved as text; compounds before and after it are still annotated.

        Steps:
        1. Build prose with a leading sentence using 新聞, a 例： range containing 新聞 / 報告書, and a trailing sentence using 意味.
        2. Annotate with vocabulary covering 新聞 / 報告書 / 意味.
        3. Assert the leading 新聞 and trailing 意味 ARE ruby-tagged, the 報告書 inside the 例： range is NOT, and the concat round-trips byte-for-byte.
        """
        value = "新聞で使う。\n\n例：新聞を読む\n報告書を書く\n\n意味を確認する。"
        tokens = annotate(value, ("新聞", "しんぶん"), ("報告書", "ほうこくしょ"), ("意味", "いみ"))
        self.assertRuby(tokens, "新聞", "しんぶん")
        self.assertNoRuby(tokens, "報告書")
        self.assertRuby(tokens, "意味", "いみ")
        self.assertEqual(concat(tokens), value)

    def test_mixed_paragraph_and_list_block(self):
        """A mixed paragraph + list entry annotates each block kind and skips example-style list items.

        Steps:
        1. Build an entry with one paragraph block containing 瞬間 and one list block whose first item contains 新聞 and second item is an example bullet containing 報告書.
        2. Run rewrite_blocks on the entry.
        3. Assert the paragraph emits ruby{瞬間→しゅんかん}, the first list item emits ruby{新聞→しんぶん}, and the example-bullet list item does NOT emit a 報告書 ruby.
        """
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
        """Annotated tokens round-trip byte-for-byte to the baseline string AND match an exact expected token sequence.

        Steps:
        1. Build a baseline string mixing time-context 十分, the 新聞 compound, and an example-bullet 報告書.
        2. Annotate with vocabulary covering 三/時/十分/新聞/報告書.
        3. Assert concat(tokens) equals the baseline byte-for-byte AND the token list matches the exact expected sequence (three rubies, separator text, 新聞 ruby, then the example-bullet tail as one text token).
        """
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
