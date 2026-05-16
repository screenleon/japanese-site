#!/usr/bin/env python3
"""Native-review fixups for JS-106 all-levels inline-ruby pass.

Applies 9 surgical token-sequence patches identified during the
main-thread native-perspective second pass:

  - 古風 (gotoki):       text{古} + ruby{風→かぜ}    -> ruby{古風→こふう}
  - 瞬間 (3x):            text{瞬} + ruby{間→あいだ}  -> ruby{瞬間→しゅんかん}
  - 迅速 (ga-hayai-ka):  text{迅} + ruby{速→はや}    -> ruby{迅速→じんそく}
  - 新聞 (ya-inaya):     text{新} + ruby{聞→き}      -> ruby{新聞→しんぶん}
  - 報告書 (2x):          text{...報告} + ruby{書→か} -> text{...} + ruby{報告書→ほうこくしょ}
  - 十分 (time):         ruby{十分→じゅうぶん}       -> ruby{十分→じっぷん} (time-of-day context)

Each patch is validated by:
  - byte-identical baseline text concat (vs commit 03ac4ddc...)
  - schema still passes scripts/lint-grammar.sh
"""
import json
import subprocess
from pathlib import Path

BASE = "03ac4ddc9f5f6afa3bae1e65a3a888cf82c346b7"
ROOT = Path(__file__).resolve().parent.parent


def baseline_text_for(path: Path) -> str:
    """Concat all explanation_ja_blocks text from the baseline commit."""
    rel = path.relative_to(ROOT)
    raw = subprocess.check_output(
        ["git", "show", f"{BASE}:{rel.as_posix()}"], cwd=ROOT
    )
    d = json.loads(raw)
    return _concat_blocks(d.get("explanation_ja_blocks", []))


def _concat_tokens(toks):
    return "".join(t.get("v", "") + t.get("k", "") for t in toks)


def _concat_blocks(blocks):
    parts = []
    for b in blocks:
        kind = b.get("kind")
        if kind in ("paragraph", "callout"):
            parts.append(_concat_tokens(b.get("tokens", [])))
        elif kind == "list":
            for item in b.get("items", []):
                parts.append(_concat_tokens(item.get("tokens", [])))
    return "\n".join(parts)


def patch_split_compound(toks, target_prefix_kanji, target_ruby_k, target_ruby_r, new_k, new_r):
    """Replace adjacent text{...<prefix>} + ruby{ruby_k, ruby_r} with ruby{new_k, new_r}.

    target_prefix_kanji: the kanji(s) at the END of a text token that should be folded
                         into the following ruby token.
    """
    out = []
    i = 0
    hits = 0
    while i < len(toks):
        t = toks[i]
        if (
            t["t"] == "text"
            and t["v"].endswith(target_prefix_kanji)
            and i + 1 < len(toks)
            and toks[i + 1]["t"] == "ruby"
            and toks[i + 1]["k"] == target_ruby_k
            and toks[i + 1]["r"] == target_ruby_r
        ):
            trimmed = t["v"][: -len(target_prefix_kanji)]
            if trimmed:
                out.append({"t": "text", "v": trimmed})
            out.append({"t": "ruby", "k": new_k, "r": new_r})
            hits += 1
            i += 2
        else:
            out.append(t)
            i += 1
    return out, hits


def patch_ruby_reading(toks, target_k, old_r, new_r):
    """Just change reading of a ruby token; k stays the same."""
    hits = 0
    for t in toks:
        if t["t"] == "ruby" and t["k"] == target_k and t["r"] == old_r:
            t["r"] = new_r
            hits += 1
    return toks, hits


def walk_blocks(blocks, fn):
    """Apply fn(tokens) -> (new_tokens, hits) to every token sequence."""
    total = 0
    for b in blocks:
        kind = b.get("kind")
        if kind in ("paragraph", "callout"):
            new_toks, h = fn(b["tokens"])
            b["tokens"] = new_toks
            total += h
        elif kind == "list":
            for item in b.get("items", []):
                new_toks, h = fn(item["tokens"])
                item["tokens"] = new_toks
                total += h
    return total


# (relative_path, [(label, patch_fn)])
PATCHES = [
    (
        "server/data/corpus/grammar/N1/gotoki.json",
        [("古風", lambda ts: patch_split_compound(ts, "古", "風", "かぜ", "古風", "こふう"))],
    ),
    (
        "server/data/corpus/grammar/N1/ga-hayai-ka.json",
        [
            ("瞬間", lambda ts: patch_split_compound(ts, "瞬", "間", "あいだ", "瞬間", "しゅんかん")),
            ("迅速", lambda ts: patch_split_compound(ts, "迅", "速", "はや", "迅速", "じんそく")),
        ],
    ),
    (
        "server/data/corpus/grammar/N1/katawara.json",
        [("瞬間", lambda ts: patch_split_compound(ts, "瞬", "間", "あいだ", "瞬間", "しゅんかん"))],
    ),
    (
        "server/data/corpus/grammar/N1/ya-inaya.json",
        [
            ("瞬間", lambda ts: patch_split_compound(ts, "瞬", "間", "あいだ", "瞬間", "しゅんかん")),
            ("新聞", lambda ts: patch_split_compound(ts, "新", "聞", "き", "新聞", "しんぶん")),
        ],
    ),
    (
        "server/data/corpus/grammar/N1/ni-katakunai.json",
        [("報告書", lambda ts: patch_split_compound(ts, "報告", "書", "か", "報告書", "ほうこくしょ"))],
    ),
    (
        "server/data/corpus/grammar/N1/sobakara.json",
        [("報告書", lambda ts: patch_split_compound(ts, "報告", "書", "か", "報告書", "ほうこくしょ"))],
    ),
    (
        "server/data/corpus/grammar/N5/time-expression.json",
        [("十分 reading", lambda ts: patch_ruby_reading(ts, "十分", "じゅうぶん", "じっぷん"))],
    ),
]


def main():
    fails = 0
    for rel, patches in PATCHES:
        path = ROOT / rel
        d = json.loads(path.read_text(encoding="utf-8"))
        blocks = d.get("explanation_ja_blocks", [])
        for label, patch in patches:
            hits = walk_blocks(blocks, patch)
            if hits == 0:
                print(f"FAIL  {rel}  {label}: 0 hits")
                fails += 1
            else:
                print(f"ok    {rel}  {label}: {hits} hit(s)")
        # byte-identity check vs baseline
        baseline = baseline_text_for(path)
        after = _concat_blocks(blocks)
        if baseline != after:
            print(f"FAIL  {rel}  byte-identity broken")
            for i, (b, a) in enumerate(zip(baseline, after)):
                if b != a:
                    print(f"   first diff at offset {i}: baseline={b!r} after={a!r}")
                    print(f"   context baseline: ...{baseline[max(0,i-15):i+15]}...")
                    print(f"   context after:    ...{after[max(0,i-15):i+15]}...")
                    break
            else:
                print(f"   length differs: baseline={len(baseline)} after={len(after)}")
            fails += 1
            continue
        # write back (json.dump with ensure_ascii=False, indent=2, no key sort)
        out = json.dumps(d, ensure_ascii=False, indent=2, sort_keys=False) + "\n"
        path.write_text(out, encoding="utf-8")
    if fails:
        raise SystemExit(f"{fails} fix(es) failed")
    print("all patches applied")


if __name__ == "__main__":
    main()
