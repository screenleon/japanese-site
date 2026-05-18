#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR=${ROOT_DIR:-$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)}
GRAMMAR_ROOT=${GRAMMAR_ROOT:-"$ROOT_DIR/server/data/corpus/grammar"}
ANNOTATION_KINDS_FILE=${ANNOTATION_KINDS_FILE:-"$ROOT_DIR/scripts/annotations-kinds.txt"}

if ! command -v node >/dev/null 2>&1; then
	echo "lint-grammar: node is required" >&2
	exit 1
fi

if [[ ! -d "$GRAMMAR_ROOT" ]]; then
	echo "lint-grammar: missing grammar root: $GRAMMAR_ROOT" >&2
	exit 1
fi

if [[ ! -f "$ANNOTATION_KINDS_FILE" ]]; then
	echo "lint-grammar: missing annotation kinds file: $ANNOTATION_KINDS_FILE" >&2
	exit 1
fi

GRAMMAR_ROOT="$GRAMMAR_ROOT" ANNOTATION_KINDS_FILE="$ANNOTATION_KINDS_FILE" node <<'NODE'
const fs = require("fs");
const path = require("path");

const grammarRoot = process.env.GRAMMAR_ROOT;
const kindsPath = process.env.ANNOTATION_KINDS_FILE;
const kinds = fs
  .readFileSync(kindsPath, "utf8")
  .split(/\r?\n/)
  .map((line) => line.trim())
  .filter((line) => line && !line.startsWith("#"));
const allowedKinds = new Set(kinds);
const allowedTop = new Set([
  "slug",
  "title_ja",
  "title_zh",
  "jlpt_level",
  "schema_version",
  "pattern",
  "explanation_ja_blocks",
  "explanation_zh",
  "_meta",
  "classifier_rules",
  "related_slugs",
  "annotations",
  "audit_status",
]);
const requiredTop = [
  "slug",
  "title_ja",
  "title_zh",
  "jlpt_level",
  "schema_version",
  "pattern",
  "explanation_ja_blocks",
  "explanation_zh",
  "_meta",
];
const legacyHints = {
  mental_model: "move to annotations.mental_model",
  nuance_note: "move to annotations.nuance_note",
  explanation_ja: "move to explanation_ja_blocks",
  source: "move to _meta.source",
  license: "move to _meta.license",
  validated_by: "move to _meta.validated_by",
  validator_score: "move to _meta.validator_score",
  key_terms: "move to annotations.furigana.vocabulary",
};
const pocSlugs = new Set(["youni-naru", "youni-suru"]);
const tokenKinds = new Set(["text", "ruby", "term"]);
const blockKinds = new Set(["paragraph", "list", "callout"]);
const calloutTones = new Set(["info", "warn", "tip"]);
const termKinds = new Set(["vocab", "grammar"]);

let exitCode = 0;
const slugs = new Map();
const related = [];
const files = [];

function fail(rel, msg) {
  console.error(`lint-grammar: ${rel} ${msg}`);
  exitCode = 2;
}

function walk(dir) {
  for (const ent of fs.readdirSync(dir, { withFileTypes: true })) {
    const p = path.join(dir, ent.name);
    if (ent.isDirectory()) walk(p);
    else if (ent.isFile() && ent.name.endsWith(".json") && !ent.name.endsWith(".examples.jsonl")) {
      files.push(p);
    }
  }
}

function isObj(value) {
  return value !== null && typeof value === "object" && !Array.isArray(value);
}

function nonEmptyString(value) {
  return typeof value === "string" && value.trim().length > 0;
}

function deepEqual(a, b) {
  return JSON.stringify(a) === JSON.stringify(b);
}

function validateTokens(tokens, rel, context) {
  if (!Array.isArray(tokens) || tokens.length === 0) {
    fail(rel, `${context} tokens must be a non-empty array`);
    return;
  }
  tokens.forEach((token, tokenIndex) => {
    const where = `${context} token ${tokenIndex}`;
    if (!isObj(token)) {
      fail(rel, `${where} must be an object`);
      return;
    }
    if (!tokenKinds.has(token.t)) {
      fail(rel, `${where} has invalid token kind '${token.t}'`);
      return;
    }
    if (token.t === "text" && !nonEmptyString(token.v)) {
      fail(rel, `${where} text.v must be a non-empty string`);
    }
    if (token.t === "ruby" && (!nonEmptyString(token.k) || !nonEmptyString(token.r))) {
      fail(rel, `${where} ruby.k and ruby.r must be non-empty strings`);
    }
    if (token.t === "term") {
      if (!nonEmptyString(token.slug)) fail(rel, `${where} term.slug must be a non-empty string`);
      if (!termKinds.has(token.kind)) fail(rel, `${where} term.kind must be vocab or grammar`);
      if (!nonEmptyString(token.label)) fail(rel, `${where} term.label must be a non-empty string`);
    }
  });
}

function validateTokenArray(tokens, rel, context) {
  if (!Array.isArray(tokens)) {
    fail(rel, `${context} must be a Token array`);
    return false;
  }
  tokens.forEach((token, tokenIndex) => {
    const where = `${context} token ${tokenIndex}`;
    if (!isObj(token)) {
      fail(rel, `${where} must be an object`);
      return;
    }
    if (Object.prototype.hasOwnProperty.call(token, "kanji") || Object.prototype.hasOwnProperty.call(token, "reading")) {
      fail(rel, `${context} must be a Token array; old {kanji, reading} title_ja pairs are disallowed`);
      return;
    }
    if (!tokenKinds.has(token.t)) {
      fail(rel, `${where} has invalid token kind '${token.t}'`);
      return;
    }
    if (token.t === "text" && !nonEmptyString(token.v)) {
      fail(rel, `${where} text.v must be a non-empty string`);
    }
    if (token.t === "ruby" && (!nonEmptyString(token.k) || !nonEmptyString(token.r))) {
      fail(rel, `${where} ruby.k and ruby.r must be non-empty strings`);
    }
    if (token.t === "term") {
      if (!nonEmptyString(token.slug)) fail(rel, `${where} term.slug must be a non-empty string`);
      if (!termKinds.has(token.kind)) fail(rel, `${where} term.kind must be vocab or grammar`);
      if (!nonEmptyString(token.label)) fail(rel, `${where} term.label must be a non-empty string`);
    }
  });
  return true;
}

function tokenText(token) {
  if (token.t === "text") return token.v;
  if (token.t === "ruby") return token.k;
  if (token.t === "term") return token.label;
  return "";
}

function titleSource(title) {
  return title.replace(/（[^）]*）$/, "");
}

function validateFuriganaPairArray(pairs, rel, context) {
  if (!Array.isArray(pairs)) {
    fail(rel, `${context} must be a FuriganaPair array`);
    return false;
  }
  pairs.forEach((pair, i) => {
    if (!isObj(pair) || !nonEmptyString(pair.kanji) || !nonEmptyString(pair.reading)) {
      fail(rel, `${context} pair ${i} must contain non-empty kanji and reading`);
    }
  });
  return true;
}

function validateBlocks(blocks, rel, context) {
  if (!Array.isArray(blocks) || blocks.length === 0) {
    fail(rel, `${context} must be a non-empty Block[]`);
    return;
  }
  blocks.forEach((block, blockIndex) => {
    const where = `${context} block ${blockIndex}`;
    if (!isObj(block)) {
      fail(rel, `${where} must be an object`);
      return;
    }
    if (!blockKinds.has(block.kind)) {
      fail(rel, `${where} has invalid block kind '${block.kind}'`);
      return;
    }
    if (block.kind === "paragraph") {
      validateTokens(block.tokens, rel, where);
    } else if (block.kind === "list") {
      if (!Array.isArray(block.items) || block.items.length === 0) {
        fail(rel, `${where} list.items must be a non-empty array`);
      } else {
        block.items.forEach((item, itemIndex) => {
          if (!isObj(item)) fail(rel, `${where} item ${itemIndex} must be an object`);
          else validateTokens(item.tokens, rel, `${where} item ${itemIndex}`);
        });
      }
    } else if (block.kind === "callout") {
      if (block.tone !== undefined && !calloutTones.has(block.tone)) {
        fail(rel, `${where} callout.tone must be info, warn, or tip`);
      }
      validateTokens(block.tokens, rel, where);
    }
  });
}

function validatePattern(gp, rel) {
  if (!Array.isArray(gp.pattern)) {
    fail(rel, "pattern must be a non-empty array of objects");
    return;
  }
  if (gp.pattern.length === 0) {
    fail(rel, "pattern must be a non-empty array");
    return;
  }
  gp.pattern.forEach((row, i) => {
    if (!isObj(row)) {
      fail(rel, `pattern[${i}] must be an object`);
      return;
    }
    if (!nonEmptyString(row.form)) fail(rel, `pattern[${i}].form must be a non-empty string`);
    if (!nonEmptyString(row.gloss_zh)) fail(rel, `pattern[${i}].gloss_zh must be a non-empty string`);
    if (row.notes_zh !== undefined && !nonEmptyString(row.notes_zh)) {
      fail(rel, `pattern[${i}].notes_zh must be non-empty when present`);
    }
  });
}

function validateMeta(gp, rel) {
  if (!isObj(gp._meta)) {
    fail(rel, "_meta must be an object");
    return;
  }
  if (!nonEmptyString(gp._meta.source)) fail(rel, "_meta.source must be a non-empty string");
  if (!nonEmptyString(gp._meta.license)) fail(rel, "_meta.license must be a non-empty string");
  if (gp._meta.validated_by !== undefined && !nonEmptyString(gp._meta.validated_by)) {
    fail(rel, "_meta.validated_by must be a non-empty string when present");
  }
  if (
    gp._meta.validator_score !== undefined &&
    (typeof gp._meta.validator_score !== "number" ||
      gp._meta.validator_score < 0 ||
      gp._meta.validator_score > 1)
  ) {
    fail(rel, "_meta.validator_score must be a number in [0,1]");
  }
}

function validateAnnotations(gp, rel) {
  if (gp.annotations === undefined) return;
  if (!isObj(gp.annotations)) {
    fail(rel, "annotations must be an object when present");
    return;
  }
  for (const key of Object.keys(gp.annotations)) {
    if (!allowedKinds.has(key)) fail(rel, `slug='${gp.slug || "?"}' annotations has unsupported kind '${key}'`);
    if (Object.prototype.hasOwnProperty.call(gp, key)) {
      fail(rel, `annotations key '${key}' overlaps top-level key`);
    }
  }
  const f = gp.annotations.furigana;
  if (f !== undefined) {
    if (!isObj(f)) {
      fail(rel, "annotations.furigana must be an object with title_ja/vocabulary pairs");
    } else {
      if (Object.prototype.hasOwnProperty.call(f, "key_terms")) {
        fail(rel, "annotations.furigana.key_terms is disallowed; rename to annotations.furigana.vocabulary");
      }
      let titleLen = 0;
      let vocabLen = 0;
      if (Object.prototype.hasOwnProperty.call(f, "title_ja")) {
        if (validateTokenArray(f.title_ja, rel, "annotations.furigana.title_ja")) {
          titleLen = f.title_ja.length;
          const rendered = f.title_ja.map(tokenText).join("");
          const expected = titleSource(gp.title_ja || "");
          if (titleLen > 0 && rendered !== expected) {
            fail(
              rel,
              `annotations.furigana.title_ja round-trip mismatch: tokens render '${rendered}', expected '${expected}'`
            );
          }
        }
      }
      if (Object.prototype.hasOwnProperty.call(f, "vocabulary")) {
        if (validateFuriganaPairArray(f.vocabulary, rel, "annotations.furigana.vocabulary")) {
          vocabLen = f.vocabulary.length;
        }
      }
      if (titleLen === 0 && vocabLen === 0) {
        fail(rel, "annotations.furigana must contain at least one title_ja token or vocabulary pair");
      }
    }
  }
  for (const key of ["mental_model", "mental_model_zh", "nuance_note"]) {
    if (gp.annotations[key] !== undefined && !nonEmptyString(gp.annotations[key])) {
      fail(rel, `annotations.${key} must be a non-empty string`);
    }
  }
  const classifier = gp.annotations.classifier;
  if (classifier !== undefined) {
    if (!isObj(classifier) || !Array.isArray(classifier.rules) || classifier.rules.length === 0) {
      fail(rel, "annotations.classifier must be an object with non-empty rules");
    } else {
      classifier.rules.forEach((rule, i) => {
        for (const key of Object.keys(rule)) {
          if (!["with_pattern", "with_slug", "rule_ja_blocks", "rule_zh", "examples"].includes(key)) {
            fail(rel, `annotations.classifier.rules[${i}] must not contain predicate key '${key}'`);
          }
        }
      });
    }
  }
}

function validateContrast(contrast, rel, context) {
  if (contrast === null || contrast === undefined) return;
  if (!isObj(contrast)) {
    fail(rel, `${context} contrast must be null or an object`);
    return;
  }
  if (!nonEmptyString(contrast.with_pattern)) fail(rel, `${context}.contrast.with_pattern must be a non-empty string`);
  validateBlocks(contrast.rule_ja_blocks, rel, `${context}.contrast.rule_ja_blocks`);
  if (contrast.with_slug !== undefined && !/^[a-z][a-z0-9-]*$/.test(contrast.with_slug)) {
    fail(rel, `${context}.contrast.with_slug must match ^[a-z][a-z0-9-]*$`);
  }
  if (contrast.rule_zh !== undefined && !nonEmptyString(contrast.rule_zh)) {
    fail(rel, `${context}.contrast.rule_zh must be non-empty when present`);
  }
  if (contrast.examples !== undefined) {
    if (!Array.isArray(contrast.examples)) {
      fail(rel, `${context}.contrast.examples must be an array when present`);
    } else {
      contrast.examples.forEach((ex, i) => {
        if (!isObj(ex) || !nonEmptyString(ex.use_this) || !nonEmptyString(ex.use_alt)) {
          fail(rel, `${context}.contrast.examples[${i}] use_this and use_alt must be non-empty strings`);
        }
        if (isObj(ex) && ex.gloss_zh !== undefined && !nonEmptyString(ex.gloss_zh)) {
          fail(rel, `${context}.contrast.examples[${i}].gloss_zh must be non-empty when present`);
        }
      });
    }
  }
}

function validateClassifierRules(gp, rel) {
  const rules = Array.isArray(gp.classifier_rules) ? gp.classifier_rules : [];
  const nonNullContrasts = [];
  rules.forEach((rule, i) => {
    if (!isObj(rule)) {
      fail(rel, `classifier_rules[${i}] must be an object`);
      return;
    }
    if (rule.default === true && rule.contrast !== null) {
      fail(rel, `classifier_rules[${i}] default rule must have contrast: null`);
    }
    if (rule.contrast !== null && rule.contrast !== undefined) {
      validateContrast(rule.contrast, rel, `classifier_rules[${i}]`);
      nonNullContrasts.push(rule.contrast);
    }
  });
  if (nonNullContrasts.length > 0) {
    if (!String(gp._meta?.validated_by || "").startsWith("native-reviewer-v")) {
      fail(rel, "non-null classifier contrast requires _meta.validated_by starting with native-reviewer-v");
    }
    if (gp.audit_status === "pre-redesign") {
      fail(rel, "pre-redesign entries must not carry non-null classifier contrast");
    }
  }
  const mirror = gp.annotations?.classifier?.rules;
  if (nonNullContrasts.length === 0) {
    if (mirror !== undefined) fail(rel, "annotations.classifier must be absent when classifier_rules has no non-null contrast");
    return;
  }
  if (!Array.isArray(mirror) || mirror.length !== nonNullContrasts.length) {
    fail(rel, "annotations.classifier.rules length must equal non-null classifier_rules contrast count");
    return;
  }
  mirror.forEach((rule, i) => {
    if (!deepEqual(rule, nonNullContrasts[i])) {
      fail(rel, `annotations.classifier.rules[${i}] must deep-equal classifier_rules non-null contrast payload`);
    }
  });
}

function validateFile(file) {
  const rel = path.relative(grammarRoot, file);
  let gp;
  try {
    gp = JSON.parse(fs.readFileSync(file, "utf8"));
  } catch (err) {
    fail(rel, `is not valid JSON: ${err.message}`);
    return;
  }
  if (!isObj(gp)) {
    fail(rel, "top-level JSON must be an object");
    return;
  }
  if (!nonEmptyString(gp.slug)) {
    fail(rel, "has missing or invalid slug");
  } else {
    const expectedSlug = path.basename(file, ".json");
    if (gp.slug !== expectedSlug) fail(rel, `slug '${gp.slug}' does not match file basename '${expectedSlug}'`);
    const prev = slugs.get(gp.slug);
    if (prev) fail(rel, `duplicate slug '${gp.slug}' in ${prev}, ${rel}`);
    else slugs.set(gp.slug, rel);
  }
  for (const key of Object.keys(gp)) {
    if (legacyHints[key]) fail(rel, `${key} is disallowed; ${legacyHints[key]}`);
    else if (!allowedTop.has(key)) fail(rel, `unknown top-level key '${key}'`);
  }
  for (const key of requiredTop) {
    if (!Object.prototype.hasOwnProperty.call(gp, key)) fail(rel, `missing required top-level key '${key}'`);
  }
  if (!Number.isInteger(gp.schema_version) || gp.schema_version !== 2) {
    fail(rel, "schema_version must equal integer 2");
  }
  validatePattern(gp, rel);
  validateBlocks(gp.explanation_ja_blocks, rel, "explanation_ja_blocks");
  validateMeta(gp, rel);
  validateAnnotations(gp, rel);
  validateClassifierRules(gp, rel);
  if (gp.related_slugs !== undefined) {
    if (!Array.isArray(gp.related_slugs) || gp.related_slugs.some((slug) => !nonEmptyString(slug))) {
      fail(rel, "related_slugs must be a string array when present");
    } else {
      for (const to of gp.related_slugs) related.push([gp.slug, to]);
    }
  }
  if (
    gp.audit_status !== undefined &&
    gp.audit_status !== "pre-redesign" &&
    gp.audit_status !== "post-dedup-naive"
  ) {
    fail(rel, 'audit_status must be one of "pre-redesign" or "post-dedup-naive" when present');
  }
  if (
    pocSlugs.has(gp.slug) &&
    (gp.audit_status === "pre-redesign" || gp.audit_status === "post-dedup-naive")
  ) {
    fail(rel, 'PoC entries must not carry audit_status: "pre-redesign" or "post-dedup-naive"');
  }
  const hasTBD = Array.isArray(gp.pattern) && gp.pattern.some((row) => row?.form === "_TBD" || row?.gloss_zh === "待補");
  if (hasTBD && gp.audit_status !== "pre-redesign") {
    fail(rel, 'pattern "_TBD" / "待補" is permitted only with audit_status: "pre-redesign"');
  }
}

walk(grammarRoot);
files.sort().forEach(validateFile);

if (files.length === 0) {
  console.error(`lint-grammar: no grammar JSON files found under ${grammarRoot}`);
  process.exit(1);
}

for (const [from, to] of related) {
  if (!slugs.has(to)) {
    console.error(`lint-grammar: ${from} related_slugs points to missing slug '${to}'`);
    exitCode = 2;
  }
}

if (exitCode) {
  console.error("lint-grammar: failed");
  process.exit(exitCode);
}
console.log("lint-grammar: passed");
NODE
