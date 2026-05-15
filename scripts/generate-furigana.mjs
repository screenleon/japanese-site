#!/usr/bin/env node
// Build-time furigana authoring helper for ADR-0002.
// Reads Japanese text from stdin (or `--text "<...>"`), runs Kuromoji
// (kuromoji@0.1.2 + mecab-ipadic-2.7.0-20070801, both Apache-2.0 / NAIST
// permissive) over it, and prints `Pair[]` JSON suitable for assembling
// an `annotations.furigana.vocabulary` Pair[] array by default, or a shared
// Token[] array for `annotations.furigana.title_ja` with `--emit token`.
// Per ADR-0001 each pair has non-whitespace `kanji` and `reading` strings;
// pure-kana Pair[] output is empty, while Token[] output preserves the text.

import { createRequire } from "node:module";
import { fileURLToPath } from "node:url";
import { dirname, resolve } from "node:path";

const HERE = dirname(fileURLToPath(import.meta.url));
const DEFAULT_DIC_PATH = resolve(HERE, "..", "web", "node_modules", "kuromoji", "dict");
// kuromoji is a devDependency of `web/`; anchor module resolution there so the
// script works regardless of CWD (direct `node scripts/...` from repo root, or
// `npm run generate:furigana` from `web/`).
const requireFromWeb = createRequire(new URL("../web/package.json", import.meta.url));
const kuromoji = requireFromWeb("kuromoji");

function isKanji(ch) {
  const c = ch.codePointAt(0);
  return (
    (c >= 0x3400 && c <= 0x4dbf) ||
    (c >= 0x4e00 && c <= 0x9fff) ||
    (c >= 0x20000 && c <= 0x2a6df) ||
    (c >= 0x2a700 && c <= 0x2ebef)
  );
}

function isHiragana(ch) {
  const c = ch.codePointAt(0);
  return c >= 0x3041 && c <= 0x309f;
}

function isKatakana(ch) {
  const c = ch.codePointAt(0);
  return c >= 0x30a0 && c <= 0x30ff;
}

function katakanaToHiragana(text) {
  let out = "";
  for (const ch of text) {
    const c = ch.codePointAt(0);
    if (c >= 0x30a1 && c <= 0x30f6) {
      out += String.fromCodePoint(c - 0x60);
    } else {
      out += ch;
    }
  }
  return out;
}

function hiraganaToKatakana(text) {
  let out = "";
  for (const ch of text) {
    const c = ch.codePointAt(0);
    if (c >= 0x3041 && c <= 0x3096) {
      out += String.fromCodePoint(c + 0x60);
    } else {
      out += ch;
    }
  }
  return out;
}

function surfaceSegments(surface) {
  const segments = [];
  let cur = null;
  for (const ch of surface) {
    const kind = isKanji(ch) ? "kanji" : isHiragana(ch) || isKatakana(ch) ? "kana" : "other";
    if (cur && cur.kind === kind) {
      cur.text += ch;
    } else {
      if (cur) segments.push(cur);
      cur = { kind, text: ch };
    }
  }
  if (cur) segments.push(cur);
  return segments;
}

function readingAnchors(segments, reading) {
  const anchors = [];
  let cursor = 0;
  for (let i = 0; i < segments.length; i++) {
    const seg = segments[i];
    if (seg.kind !== "kana") continue;
    const target = hiraganaToKatakana(seg.text);
    const idx = reading.indexOf(target, cursor);
    if (idx === -1) return null;
    anchors.push({ segIdx: i, readingStart: idx, readingEnd: idx + target.length });
    cursor = idx + target.length;
  }
  return anchors;
}

function kanjiReadingsFromSegments(segments, reading) {
  const anchors = readingAnchors(segments, reading);
  if (!anchors) return null;

  const readings = new Map();
  let readingPos = 0;
  let anchorIdx = 0;
  for (let i = 0; i < segments.length; i++) {
    const seg = segments[i];
    if (seg.kind === "kana") {
      const anchor = anchors[anchorIdx++];
      readingPos = anchor.readingEnd;
      continue;
    }
    if (seg.kind !== "kanji") continue;
    const nextAnchor = anchors[anchorIdx];
    const end = nextAnchor ? nextAnchor.readingStart : reading.length;
    const runReading = reading.slice(readingPos, end);
    if (runReading.length === 0) return null;
    const hira = katakanaToHiragana(runReading);
    if (hira.trim().length === 0 || seg.text.trim().length === 0) return null;
    readings.set(i, hira);
    readingPos = end;
  }
  return readings;
}

function tokenHasKanji(token) {
  return [...(token.surface_form ?? "")].some(isKanji);
}

// Walk a kuromoji token's surface_form alongside its (katakana) reading,
// splitting the surface into runs of kanji vs. kana and assigning each kanji-run
// a reading suffix-stripped from the token's full reading. Returns Pair[] with
// hiragana readings.
function pairsFromToken(token) {
  const surface = token.surface_form ?? "";
  const reading = token.reading;
  // Reading missing or "*" → unknown reading. Skip token entirely.
  if (!reading || reading === "*") return [];

  // Walk surface character-by-character, segmenting kanji vs. kana runs.
  const segments = surfaceSegments(surface).filter((segment) => segment.kind !== "other");

  if (segments.length === 0) return [];

  // No kanji runs → nothing to annotate.
  if (!segments.some((s) => s.kind === "kanji")) return [];

  // Anchor each kana-run inside the reading by exact match (kana-run converted to
  // katakana). Each kanji-run's reading is the slice of reading between the
  // previous anchor's end (or 0) and the next anchor's start (or reading length).
  const kanjiReadings = kanjiReadingsFromSegments(segments, reading);
  if (!kanjiReadings) return [];

  const pairs = [];
  for (let i = 0; i < segments.length; i++) {
    const seg = segments[i];
    if (seg.kind !== "kanji") continue;
    const readingForRun = kanjiReadings.get(i);
    if (readingForRun) pairs.push({ kanji: seg.text, reading: readingForRun });
  }
  return pairs;
}

function tokensFromToken(token) {
  const surface = token.surface_form ?? "";
  if (!surface) return [];
  const reading = token.reading;
  const segments = surfaceSegments(surface);

  if (!tokenHasKanji(token)) {
    return [{ t: "text", v: surface }];
  }
  if (!reading || reading === "*") return [];

  const kanjiReadings = kanjiReadingsFromSegments(segments, reading);
  if (!kanjiReadings) return [];

  const out = [];
  for (let i = 0; i < segments.length; i++) {
    const seg = segments[i];
    if (seg.kind === "kanji") {
      const r = kanjiReadings.get(i);
      if (!r) return [];
      out.push({ t: "ruby", k: seg.text, r });
    } else {
      out.push({ t: "text", v: seg.text });
    }
  }
  return out;
}

function mergeAdjacentTextTokens(tokens) {
  const out = [];
  for (const token of tokens) {
    if (token.t === "text" && token.v === "") continue;
    const prev = out[out.length - 1];
    if (token.t === "text" && prev?.t === "text") {
      prev.v += token.v;
    } else {
      out.push(token);
    }
  }
  return out;
}

// De-dupe consecutive identical pairs (e.g. a token boundary that splits the same
// kanji+okurigana across two tokens occasionally produces duplicate pairs). Keep
// the first occurrence.
function dedupePairs(pairs) {
  const seen = new Set();
  const out = [];
  for (const p of pairs) {
    const key = `${p.kanji}\t${p.reading}`;
    if (seen.has(key)) continue;
    seen.add(key);
    out.push(p);
  }
  return out;
}

function buildTokenizer(dicPath) {
  return new Promise((resolveBuilder, rejectBuilder) => {
    kuromoji.builder({ dicPath }).build((err, tokenizer) => {
      if (err) rejectBuilder(err);
      else resolveBuilder(tokenizer);
    });
  });
}

async function readStdin() {
  let data = "";
  for await (const chunk of process.stdin) data += chunk;
  return data;
}

function parseArgs(argv) {
  const args = { text: null, dicPath: DEFAULT_DIC_PATH, emit: "pair", debug: false };
  for (let i = 2; i < argv.length; i++) {
    const a = argv[i];
    if (a === "--text") args.text = argv[++i];
    else if (a === "--dic-path") args.dicPath = argv[++i];
    else if (a === "--emit") {
      args.emit = argv[++i];
      if (!["pair", "token"].includes(args.emit)) {
        console.error("generate-furigana: --emit must be 'pair' or 'token'");
        process.exit(2);
      }
    }
    else if (a === "--debug") args.debug = true;
    else if (a === "--help" || a === "-h") {
      console.log(
        "Usage: generate-furigana.mjs [--text \"<japanese>\"] [--dic-path <path>] [--emit pair|token] [--debug]\n" +
          "       echo \"<japanese>\" | generate-furigana.mjs\n" +
          "Reads Japanese text and prints Pair[] by default, or Token[] with --emit token.\n" +
          "Pure-kana input prints `[]` in pair mode and a text token in token mode."
      );
      process.exit(0);
    } else {
      console.error(`generate-furigana: unknown argument '${a}'`);
      process.exit(2);
    }
  }
  return args;
}

async function main() {
  const args = parseArgs(process.argv);
  const text = (args.text ?? (await readStdin())).trim();
  if (!text) {
    console.log("[]");
    return;
  }
  const tokenizer = await buildTokenizer(args.dicPath);
  const tokens = tokenizer.tokenize(text);
  if (args.emit === "token") {
    const allTokens = [];
    for (const tok of tokens) {
      if (args.debug) console.error(JSON.stringify({ surface: tok.surface_form, reading: tok.reading, pos: tok.pos }));
      for (const token of tokensFromToken(tok)) allTokens.push(token);
    }
    console.log(JSON.stringify(mergeAdjacentTextTokens(allTokens)));
    return;
  }

  const allPairs = [];
  for (const tok of tokens) {
    if (args.debug) console.error(JSON.stringify({ surface: tok.surface_form, reading: tok.reading, pos: tok.pos }));
    const pairs = pairsFromToken(tok);
    for (const p of pairs) allPairs.push(p);
  }
  console.log(JSON.stringify(dedupePairs(allPairs)));
}

main().catch((err) => {
  console.error(`generate-furigana: ${err?.message ?? err}`);
  process.exit(1);
});
