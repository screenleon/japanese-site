#!/usr/bin/env python3
"""
Apply JS-106 inline-ruby tokenisation to N5, N4, N2, and N1 grammar entries.

The script is intentionally content-only:
  - rewrites only explanation_ja_blocks[*].tokens arrays
  - preserves block kind/order and all non-target fields
  - verifies token text concatenation against the baseline commit
  - verifies non-explanation fields against the baseline commit

Run from repo root:
    python3 scripts/apply-allLevels-inline-ruby.py
"""

from __future__ import annotations

import copy
import json
import random
import re
import subprocess
import sys
from collections import Counter
from pathlib import Path


BASE = "03ac4ddc9f5f6afa3bae1e65a3a888cf82c346b7"
LEVELS = ("N5", "N4", "N2", "N1")
ROOT = Path("server/data/corpus/grammar")
AUDIT_PATH = Path("audits/js-106-allLevels-codex-pass-2026-05-16.md")


def text(v: str) -> dict[str, str]:
    return {"t": "text", "v": v}


def ruby(k: str, r: str) -> dict[str, str]:
    return {"t": "ruby", "k": k, "r": r}


KANJI_RE = re.compile(r"[\u3400-\u4dbf\u4e00-\u9fff]")
KANA_RE = re.compile(r"[ぁ-ゖ]")


SKIP_WORDS = {
    # JS-106 convention: common body-prose words are left as text to avoid
    # over-cluttering explanations.
    "食べる",
    "行く",
    "来る",
    "見る",
    "聞く",
    "高い",
    "大きい",
    "雨",
    "山",
    "川",
    "人",
}


OVERRIDE_READINGS = {
    # Compound grammar metalanguage.
    "一段動詞": "いちだんどうし",
    "五段動詞": "ごだんどうし",
    "辞書形": "じしょけい",
    "普通形": "ふつうけい",
    "普通体": "ふつうたい",
    "丁寧形": "ていねいけい",
    "否定形": "ひていけい",
    "過去形": "かこけい",
    "連用形": "れんようけい",
    "仮定形": "かていけい",
    "可能形": "かのうけい",
    "受身形": "うけみけい",
    "使役形": "しえきけい",
    "意向形": "いこうけい",
    "命令形": "めいれいけい",
    "禁止形": "きんしけい",
    "ます形": "ますけい",
    "て形": "てけい",
    "た形": "たけい",
    "ば形": "ばけい",
    "ない形": "ないけい",
    "な形容詞": "なけいようし",
    "い形容詞": "いけいようし",
    "形容詞語幹": "けいようしごかん",
    "動詞辞書形": "どうしじしょけい",
    "動詞普通形": "どうしふつうけい",
    "動詞連用形": "どうしれんようけい",
    "動詞否定形": "どうしひていけい",
    "動詞意向形": "どうしいこうけい",
    "動詞可能形": "どうしかのうけい",
    "名詞修飾": "めいししゅうしょく",
    "文末表現": "ぶんまつひょうげん",
    "接続詞": "せつぞくし",
    "副詞": "ふくし",
    "助詞": "じょし",
    "疑問詞": "ぎもんし",
    "自動詞": "じどうし",
    "他動詞": "たどうし",
    "受身": "うけみ",
    "使役": "しえき",
    # Frequent content / discourse compounds.
    "意味核": "いみかく",
    "前件": "ぜんけん",
    "後件": "こうけん",
    "前提条件": "ぜんていじょうけん",
    "依存関係": "いぞんかんけい",
    "類義差別": "るいぎさべつ",
    "日常会話": "にちじょうかいわ",
    "話し言葉": "はなしことば",
    "書き言葉": "かきことば",
    "話し手": "はなして",
    "聞き手": "ききて",
    "相手": "あいて",
    "事実": "じじつ",
    "実際": "じっさい",
    "意外": "いがい",
    "不満": "ふまん",
    "譲歩": "じょうほ",
    "逆接": "ぎゃくせつ",
    "仮定": "かてい",
    "条件": "じょうけん",
    "結果": "けっか",
    "理由": "りゆう",
    "原因": "げんいん",
    "目的": "もくてき",
    "判断": "はんだん",
    "推量": "すいりょう",
    "評価": "ひょうか",
    "程度": "ていど",
    "状態": "じょうたい",
    "動作": "どうさ",
    "行為": "こうい",
    "動詞": "どうし",
    "形容詞": "けいようし",
    "名詞": "めいし",
    "否定": "ひてい",
    "肯定": "こうてい",
    "過去": "かこ",
    "現在": "げんざい",
    "丁寧": "ていねい",
    "普通": "ふつう",
    "表現": "ひょうげん",
    "接続": "せつぞく",
    "内容": "ないよう",
    "説明": "せつめい",
    "区別": "くべつ",
    "強調": "きょうちょう",
    "批判": "ひはん",
    "禁止": "きんし",
    "許可": "きょか",
    "義務": "ぎむ",
    "可能": "かのう",
    "不可能": "ふかのう",
    "必要": "ひつよう",
    "自然": "しぜん",
    "場面": "ばめん",
    "場合": "ばあい",
    "状況": "じょうきょう",
    "関係": "かんけい",
    "語幹": "ごかん",
    "習慣": "しゅうかん",
    "例文": "れいぶん",
}


ALLOW_KANA_TERMS = {
    "な形容詞",
    "い形容詞",
    "て形",
    "た形",
    "ば形",
    "ます形",
    "ない形",
    "書き言葉",
    "話し言葉",
    "話し手",
    "聞き手",
    "作り方",
    "読み方",
    "使い方",
    "気持ち",
    "一つ",
    "二つ",
    "三つ",
    "後ろ",
    "ご飯",
    "試し",
    "違い",
    "感じ",
    "生き物",
    "行く先",
    "代わり",
    "支え",
    "心構え",
}


AMBIGUOUS_TERMS = {
    "行く": "行く／行う family",
    "行う": "行く／行う family",
    "開く": "開く／開ける family",
    "開ける": "開く／開ける family",
    "表す": "表す／現す／著す family",
    "現す": "表す／現す／著す family",
    "著す": "表す／現す／著す family",
    "入る": "入る／入れる family",
    "入れる": "入る／入れる family",
    "上がる": "上がる／上げる／上る／上す family",
    "上げる": "上がる／上げる／上る／上す family",
    "上る": "上がる／上げる／上る／上す family",
    "上す": "上がる／上げる／上る／上す family",
    "立つ": "立つ／建つ／発つ family",
    "建つ": "立つ／建つ／発つ family",
    "発つ": "立つ／建つ／発つ family",
    "続く": "続く／続ける family",
    "続ける": "続く／続ける family",
    "通る": "通る／通う／通す family",
    "通う": "通る／通う／通す family",
    "通す": "通る／通う／通す family",
    "直る": "直る／直す family",
    "直す": "直る／直す family",
}


def load_json(path: Path) -> dict:
    with path.open(encoding="utf-8") as f:
        return json.load(f)


def load_baseline(path: Path) -> dict:
    proc = subprocess.run(
        ["git", "show", f"{BASE}:{path.as_posix()}"],
        check=True,
        text=True,
        stdout=subprocess.PIPE,
    )
    return json.loads(proc.stdout)


def concat_tokens(tokens: list[dict]) -> str:
    out = []
    for token in tokens:
        if token["t"] == "text":
            out.append(token["v"])
        elif token["t"] == "ruby":
            out.append(token["k"])
        else:
            raise AssertionError(f"unexpected token kind {token['t']}")
    return "".join(out)


def block_texts(blocks: list[dict]) -> list:
    result = []
    for block in blocks:
        if block["kind"] in {"paragraph", "callout"}:
            result.append((block["kind"], concat_tokens(block["tokens"])))
        elif block["kind"] == "list":
            result.append(
                (
                    "list",
                    [concat_tokens(item["tokens"]) for item in block["items"]],
                )
            )
        else:
            raise AssertionError(f"unexpected block kind {block['kind']}")
    return result


def non_explanation_projection(entry: dict) -> dict:
    projected = copy.deepcopy(entry)
    projected.pop("explanation_ja_blocks", None)
    return projected


def add_reading(readings: dict[str, str], kanji: str, reading: str) -> None:
    if not kanji or not reading:
        return
    if kanji in SKIP_WORDS:
        return
    if not KANJI_RE.search(kanji):
        return
    if kanji == "形":
        return
    if KANA_RE.search(kanji) and kanji not in ALLOW_KANA_TERMS and kanji not in OVERRIDE_READINGS:
        return
    readings[kanji] = reading


def n3_readings() -> dict[str, str]:
    readings = {}
    for path in sorted((ROOT / "N3").glob("*.json")):
        entry = load_json(path)
        for block in entry["explanation_ja_blocks"]:
            token_lists = [block["tokens"]] if block["kind"] != "list" else [item["tokens"] for item in block["items"]]
            for tokens in token_lists:
                for token in tokens:
                    if token["t"] == "ruby":
                        add_reading(readings, token["k"], token["r"])
    readings.update(OVERRIDE_READINGS)
    return readings


GLOBAL_READINGS = n3_readings()


def entry_readings(entry: dict) -> dict[str, str]:
    readings = dict(GLOBAL_READINGS)
    furigana = entry.get("annotations", {}).get("furigana", {})
    for item in furigana.get("vocabulary", []) or []:
        add_reading(readings, item.get("kanji", ""), item.get("reading", ""))
    # Overrides win after vocabulary imports so compounds are not split.
    readings.update(OVERRIDE_READINGS)
    return readings


def sorted_terms(readings: dict[str, str]) -> list[tuple[str, str]]:
    return sorted(readings.items(), key=lambda item: (-len(item[0]), item[0]))


def flush_text(tokens: list[dict], buf: list[str]) -> None:
    if not buf:
        return
    value = "".join(buf)
    if value:
        tokens.append(text(value))
    buf.clear()


def annotate_segment(segment: str, terms: list[tuple[str, str]]) -> list[dict]:
    tokens: list[dict] = []
    buf: list[str] = []
    i = 0
    while i < len(segment):
        match = None
        for kanji, reading in terms:
            if segment.startswith(kanji, i):
                match = (kanji, reading)
                break
        if match is None:
            buf.append(segment[i])
            i += 1
            continue
        flush_text(tokens, buf)
        tokens.append(ruby(match[0], match[1]))
        i += len(match[0])
    flush_text(tokens, buf)
    return tokens


def is_example_bullet(line: str) -> bool:
    stripped = line.lstrip()
    if not stripped.startswith("- "):
        return False
    # Form-definition bullets are grammar prose and should still receive ruby.
    grammar_markers = (
        "動詞",
        "名詞",
        "形容詞",
        "否定",
        "普通形",
        "辞書形",
        "接続",
        "形：",
        "V",
        "N",
        "いAdj",
        "なAdj",
    )
    return not stripped[2:].startswith(grammar_markers)


def protected_example_ranges(value: str) -> list[tuple[int, int]]:
    ranges: list[tuple[int, int]] = []
    start = 0
    while True:
        idx = value.find("例：", start)
        if idx == -1:
            break
        end = value.find("\n\n", idx)
        if end == -1:
            end = len(value)
        ranges.append((idx, end))
        start = end
    return ranges


def annotate_text(value: str, readings: dict[str, str]) -> list[dict]:
    terms = sorted_terms(readings)
    tokens: list[dict] = []
    protected = protected_example_ranges(value)
    cursor = 0

    def annotate_lines(segment: str) -> None:
        lines = segment.splitlines(keepends=True)
        for line in lines:
            if is_example_bullet(line):
                tokens.append(text(line))
            else:
                tokens.extend(annotate_segment(line, terms))
        if not lines and segment:
            tokens.extend(annotate_segment(segment, terms))

    for start, end in protected:
        if cursor < start:
            annotate_lines(value[cursor:start])
        tokens.append(text(value[start:end]))
        cursor = end
    if cursor < len(value):
        annotate_lines(value[cursor:])

    merged: list[dict] = []
    for token in tokens:
        if token["t"] == "text" and not token["v"]:
            continue
        if merged and token["t"] == "text" and merged[-1]["t"] == "text":
            merged[-1]["v"] += token["v"]
        else:
            merged.append(token)
    return merged


def rewrite_tokens(tokens: list[dict], readings: dict[str, str]) -> list[dict]:
    rewritten: list[dict] = []
    for token in tokens:
        if token["t"] == "text":
            rewritten.extend(annotate_text(token["v"], readings))
        elif token["t"] == "ruby":
            rewritten.append(token)
        else:
            raise AssertionError("term tokens are out of scope for JS-106")
    return rewritten


def rewrite_blocks(entry: dict) -> list[dict]:
    readings = entry_readings(entry)
    blocks = copy.deepcopy(entry["explanation_ja_blocks"])
    for block in blocks:
        if block["kind"] in {"paragraph", "callout"}:
            block["tokens"] = rewrite_tokens(block["tokens"], readings)
        elif block["kind"] == "list":
            for item in block["items"]:
                item["tokens"] = rewrite_tokens(item["tokens"], readings)
        else:
            raise AssertionError(f"unexpected block kind {block['kind']}")
    return blocks


def ruby_tokens(entry: dict) -> list[tuple[str, str]]:
    found = []
    for block in entry["explanation_ja_blocks"]:
        token_lists = [block["tokens"]] if block["kind"] != "list" else [item["tokens"] for item in block["items"]]
        for tokens in token_lists:
            for token in tokens:
                if token["t"] == "ruby":
                    found.append((token["k"], token["r"]))
    return found


def verify(path: Path, baseline: dict, rewritten: dict) -> None:
    if non_explanation_projection(baseline) != non_explanation_projection(rewritten):
        raise AssertionError(f"{path}: non-explanation fields changed")
    if block_texts(baseline["explanation_ja_blocks"]) != block_texts(rewritten["explanation_ja_blocks"]):
        raise AssertionError(f"{path}: explanation text concatenation changed")
    if not ruby_tokens(rewritten):
        raise AssertionError(f"{path}: no ruby tokens emitted")


def write_json(path: Path, entry: dict) -> None:
    path.write_text(
        json.dumps(entry, ensure_ascii=False, indent=2, sort_keys=False) + "\n",
        encoding="utf-8",
    )


def audit(level_counts: dict[str, int], freq: Counter, spot_checks: dict[str, list[str]], ambiguous: list[str]) -> None:
    lines = [
        "# JS-106 N5+N4+N2+N1 Inline Ruby — Codex Pass",
        "",
        "**Date**: 2026-05-16",
        "**Author**: codex",
        "**Scope**: N5, N4, N2, and N1 `explanation_ja_blocks` only. N3 was left untouched.",
        "",
        "## Per-level Counts",
        "",
        "| Level | Entries rewritten | Entries with ruby |",
        "|---|---:|---:|",
    ]
    for level in LEVELS:
        lines.append(f"| {level} | {level_counts[level]} | {level_counts[level]} |")
    lines.extend(
        [
            "",
            "## Spot-check Slugs",
            "",
        ]
    )
    for level in LEVELS:
        lines.append(f"- {level}: {', '.join(spot_checks[level])}")
    lines.extend(
        [
            "",
            "## Compound Reading Frequency",
            "",
            "| Compound | Reading | Occurrences |",
            "|---|---|---:|",
        ]
    )
    for (kanji, reading), count in sorted(freq.items(), key=lambda item: (-item[1], item[0][0], item[0][1])):
        if count >= 3:
            lines.append(f"| {kanji} | {reading} | {count} |")
    lines.extend(
        [
            "",
            "## Ambiguous Reading Flags",
            "",
        ]
    )
    if ambiguous:
        lines.extend(f"- {item}" for item in ambiguous)
    else:
        lines.append("- None from the explicit handoff watchlist were ruby-tagged outside fixed compounds.")
    lines.extend(
        [
            "",
            "## Verification",
            "",
            "- `python3 scripts/apply-allLevels-inline-ruby.py`: exit 0",
            "- `bash scripts/lint-grammar.sh`: PENDING",
            "- `make test`: PENDING",
            "- `cd web && npm test`: PENDING",
            "",
            "## Process Notes",
            "",
            "- The script verifies every rewritten file against baseline commit `03ac4ddc9f5f6afa3bae1e65a3a888cf82c346b7`.",
            "- For every block, concatenating `text.v` and `ruby.k` in document order matches the baseline text byte-for-byte.",
            "- For every entry, all non-`explanation_ja_blocks` fields match the baseline JSON projection.",
            "- Markdown example bullet lines were kept as trailing text runs unless they are form-definition bullets.",
            "- Native-perspective review should still focus on N1/N2 readings and any compound where the pedagogical context could prefer a less common reading.",
            "",
        ]
    )
    AUDIT_PATH.write_text("\n".join(lines), encoding="utf-8")


def main() -> int:
    level_counts: dict[str, int] = {}
    freq: Counter = Counter()
    ambiguous: list[str] = []
    all_paths: list[Path] = []

    for level in LEVELS:
        paths = sorted((ROOT / level).glob("*.json"))
        if len(paths) != 40:
            raise AssertionError(f"{level}: expected 40 files, found {len(paths)}")
        all_paths.extend(paths)

    for path in all_paths:
        baseline = load_baseline(path)
        rewritten = copy.deepcopy(baseline)
        rewritten["explanation_ja_blocks"] = rewrite_blocks(baseline)
        verify(path, baseline, rewritten)
        write_json(path, rewritten)
        pairs = ruby_tokens(rewritten)
        freq.update(pairs)
        level_counts[path.parent.name] = level_counts.get(path.parent.name, 0) + 1
        for kanji, reading in pairs:
            family = AMBIGUOUS_TERMS.get(kanji)
            if family:
                ambiguous.append(f"{path.parent.name}/{path.name}: {kanji} ({reading}) — {family}")

    rng = random.Random(106)
    spot_checks: dict[str, list[str]] = {}
    for level in LEVELS:
        slugs = sorted(path.stem for path in (ROOT / level).glob("*.json"))
        spot_checks[level] = sorted(rng.sample(slugs, 5))

    audit(level_counts, freq, spot_checks, sorted(set(ambiguous)))
    print(f"rewrote {sum(level_counts.values())} entries")
    return 0


if __name__ == "__main__":
    sys.exit(main())
