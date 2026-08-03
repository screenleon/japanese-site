#!/usr/bin/env bash
# Validate School Kokugo L1 units under server/data/corpus/kokugo/**
# Contract: docs/adr/0005-kokugo-track.md + web/src/kokugoTypes.ts
set -euo pipefail

ROOT_DIR=${ROOT_DIR:-$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)}
KOKUGO_ROOT=${KOKUGO_ROOT:-"$ROOT_DIR/server/data/corpus/kokugo"}

if ! command -v node >/dev/null 2>&1; then
	echo "lint-kokugo: node is required" >&2
	exit 1
fi

if [[ ! -d "$KOKUGO_ROOT" ]]; then
	echo "lint-kokugo: missing kokugo root: $KOKUGO_ROOT" >&2
	exit 1
fi

KOKUGO_ROOT="$KOKUGO_ROOT" ROOT_DIR="$ROOT_DIR" node <<'NODE'
const fs = require("fs");
const path = require("path");

const kokugoRoot = process.env.KOKUGO_ROOT;
const rootDir = process.env.ROOT_DIR;

const STAGES_V1 = new Set(["e5-6"]);
const SUPPORT = new Set(["heavy", "n3", "standard", "none"]);
const GENRES = new Set(["story", "expository", "opinion", "poetry"]);
const SKILLS = new Set([
  "reading.predict",
  "reading.locate-evidence",
  "reading.structure",
  "reading.summary",
  "writing.claim-reason",
  "writing.revision",
]);
const TASK_KINDS = new Set([
  "predict",
  "evidence-highlight",
  "paragraph-role",
  "summary-choice",
]);
const ARTIFACT_KINDS = new Set(["short-proposal", "summary"]);
const TOKEN_KINDS = new Set(["text", "ruby", "term"]);
const BLOCK_KINDS = new Set(["paragraph", "list", "callout"]);
const CALLOUT_TONES = new Set(["info", "warn", "tip"]);
const TERM_KINDS = new Set(["vocab", "grammar"]);

const REQUIRED_TOP = [
  "id",
  "schema_version",
  "stage",
  "title_ja",
  "genre",
  "objectives",
  "estimated_minutes",
  "text",
  "support",
  "tasks",
  "_meta",
];
const ALLOWED_TOP = new Set([
  ...REQUIRED_TOP,
  "artifact",
  "classmates",
]);

let exitCode = 0;
const ids = new Map();
const files = [];

function fail(rel, msg) {
  console.error(`lint-kokugo: ${rel}: ${msg}`);
  exitCode = 2;
}

function isObj(value) {
  return value !== null && typeof value === "object" && !Array.isArray(value);
}

function nonEmptyString(value) {
  return typeof value === "string" && value.trim().length > 0;
}

function walk(dir) {
  for (const ent of fs.readdirSync(dir, { withFileTypes: true })) {
    const p = path.join(dir, ent.name);
    if (ent.isDirectory()) walk(p);
    else if (ent.isFile() && ent.name.endsWith(".json")) files.push(p);
  }
}

function validateTokens(tokens, rel, context) {
  if (!Array.isArray(tokens) || tokens.length === 0) {
    fail(rel, `${context} tokens must be a non-empty array`);
    return;
  }
  tokens.forEach((token, i) => {
    const where = `${context} token ${i}`;
    if (!isObj(token)) {
      fail(rel, `${where} must be an object`);
      return;
    }
    if (!TOKEN_KINDS.has(token.t)) {
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
      if (!TERM_KINDS.has(token.kind)) fail(rel, `${where} term.kind must be vocab or grammar`);
      if (!nonEmptyString(token.label)) fail(rel, `${where} term.label must be a non-empty string`);
    }
  });
}

function validateBlocks(blocks, rel, context) {
  if (!Array.isArray(blocks) || blocks.length === 0) {
    fail(rel, `${context} must be a non-empty Block[]`);
    return 0;
  }
  let paragraphCount = 0;
  blocks.forEach((block, i) => {
    const where = `${context} block ${i}`;
    if (!isObj(block)) {
      fail(rel, `${where} must be an object`);
      return;
    }
    if (!BLOCK_KINDS.has(block.kind)) {
      fail(rel, `${where} has invalid kind '${block.kind}'`);
      return;
    }
    if (block.kind === "paragraph") {
      paragraphCount += 1;
      validateTokens(block.tokens, rel, where);
    } else if (block.kind === "list") {
      if (!Array.isArray(block.items) || block.items.length === 0) {
        fail(rel, `${where} list.items must be a non-empty array`);
        return;
      }
      block.items.forEach((item, j) => {
        if (!isObj(item)) {
          fail(rel, `${where} item ${j} must be an object`);
          return;
        }
        validateTokens(item.tokens, rel, `${where} item ${j}`);
      });
    } else if (block.kind === "callout") {
      if (block.tone !== undefined && !CALLOUT_TONES.has(block.tone)) {
        fail(rel, `${where} callout.tone must be info|warn|tip when set`);
      }
      validateTokens(block.tokens, rel, where);
    }
  });
  return paragraphCount;
}

function plainTextFromBlocks(blocks) {
  const parts = [];
  for (const block of blocks) {
    if (!isObj(block)) continue;
    if (block.kind === "paragraph" || block.kind === "callout") {
      for (const t of block.tokens || []) {
        if (t.t === "text") parts.push(t.v);
        else if (t.t === "ruby") parts.push(t.k);
        else if (t.t === "term") parts.push(t.label);
      }
    } else if (block.kind === "list") {
      for (const item of block.items || []) {
        for (const t of item.tokens || []) {
          if (t.t === "text") parts.push(t.v);
          else if (t.t === "ruby") parts.push(t.k);
          else if (t.t === "term") parts.push(t.label);
        }
      }
    }
  }
  return parts.join("");
}

function validateChoices(choices, rel, context) {
  if (!Array.isArray(choices) || choices.length < 2) {
    fail(rel, `${context} choices must be an array of length >= 2`);
    return;
  }
  const seen = new Set();
  choices.forEach((c, i) => {
    if (!isObj(c) || !nonEmptyString(c.id) || !nonEmptyString(c.text_ja)) {
      fail(rel, `${context} choice ${i} needs non-empty id and text_ja`);
      return;
    }
    if (seen.has(c.id)) fail(rel, `${context} duplicate choice id '${c.id}'`);
    seen.add(c.id);
  });
}

function validatePredict(payload, rel, where) {
  if (!nonEmptyString(payload.prompt_ja)) fail(rel, `${where} prompt_ja required`);
  validateChoices(payload.choices, rel, where);
  if (payload.allow_free_text !== undefined && typeof payload.allow_free_text !== "boolean") {
    fail(rel, `${where} allow_free_text must be boolean when set`);
  }
}

function validateEvidence(payload, rel, where, plainText) {
  if (!nonEmptyString(payload.prompt_ja)) fail(rel, `${where} prompt_ja required`);
  if (!Array.isArray(payload.gold_quotes) || payload.gold_quotes.length === 0) {
    fail(rel, `${where} gold_quotes must be a non-empty string array`);
    return;
  }
  payload.gold_quotes.forEach((q, i) => {
    if (!nonEmptyString(q)) {
      fail(rel, `${where} gold_quotes[${i}] must be a non-empty string`);
      return;
    }
    if (!plainText.includes(q)) {
      fail(rel, `${where} gold_quotes[${i}] not found in unit text`);
    }
  });
}

function validateParagraphRole(payload, rel, where, paragraphCount) {
  if (!nonEmptyString(payload.prompt_ja)) fail(rel, `${where} prompt_ja required`);
  if (!Array.isArray(payload.roles) || payload.roles.length < 2) {
    fail(rel, `${where} roles must be an array of length >= 2`);
    return;
  }
  payload.roles.forEach((r, i) => {
    if (!nonEmptyString(r)) fail(rel, `${where} roles[${i}] must be non-empty`);
  });
  const roleSet = new Set(payload.roles);
  if (!Array.isArray(payload.gold_by_paragraph_index)) {
    fail(rel, `${where} gold_by_paragraph_index must be an array`);
    return;
  }
  if (payload.gold_by_paragraph_index.length !== paragraphCount) {
    fail(
      rel,
      `${where} gold_by_paragraph_index length ${payload.gold_by_paragraph_index.length} != paragraph count ${paragraphCount}`,
    );
  }
  payload.gold_by_paragraph_index.forEach((g, i) => {
    if (!nonEmptyString(g) || !roleSet.has(g)) {
      fail(rel, `${where} gold_by_paragraph_index[${i}] must be one of roles`);
    }
  });
}

function validateSummaryChoice(payload, rel, where) {
  if (!nonEmptyString(payload.prompt_ja)) fail(rel, `${where} prompt_ja required`);
  validateChoices(payload.choices, rel, where);
  if (!nonEmptyString(payload.correct_id)) {
    fail(rel, `${where} correct_id required`);
    return;
  }
  const ids = new Set((payload.choices || []).map((c) => c && c.id));
  if (!ids.has(payload.correct_id)) {
    fail(rel, `${where} correct_id '${payload.correct_id}' not in choices`);
  }
}

function validateTask(task, rel, index, plainText, paragraphCount, taskIds) {
  const where = `tasks[${index}]`;
  if (!isObj(task)) {
    fail(rel, `${where} must be an object`);
    return;
  }
  if (!nonEmptyString(task.id)) fail(rel, `${where} id required`);
  else if (taskIds.has(task.id)) fail(rel, `${where} duplicate task id '${task.id}'`);
  else taskIds.add(task.id);

  if (!SKILLS.has(task.skill)) fail(rel, `${where} invalid skill '${task.skill}'`);
  if (!TASK_KINDS.has(task.kind)) fail(rel, `${where} invalid kind '${task.kind}' (v1 closed set)`);
  if (!isObj(task.payload)) {
    fail(rel, `${where} payload must be an object`);
    return;
  }
  if (task.rubric !== undefined && !isObj(task.rubric)) {
    fail(rel, `${where} rubric must be an object when set`);
  }

  switch (task.kind) {
    case "predict":
      validatePredict(task.payload, rel, where);
      break;
    case "evidence-highlight":
      validateEvidence(task.payload, rel, where, plainText);
      break;
    case "paragraph-role":
      validateParagraphRole(task.payload, rel, where, paragraphCount);
      break;
    case "summary-choice":
      validateSummaryChoice(task.payload, rel, where);
      break;
    default:
      fail(rel, `${where} unhandled kind`);
  }
}

function validateArtifact(artifact, rel) {
  if (!isObj(artifact)) {
    fail(rel, "artifact must be an object when set");
    return;
  }
  if (!ARTIFACT_KINDS.has(artifact.kind)) {
    fail(rel, `artifact.kind must be short-proposal|summary`);
  }
  // min_chars / max_chars: 0 means "no bound" (progressive writing). When both
  // are > 0, min must be <= max. Non-integers or negatives fail.
  if (typeof artifact.min_chars !== "number" || !Number.isInteger(artifact.min_chars) || artifact.min_chars < 0) {
    fail(rel, "artifact.min_chars must be an integer >= 0 (0 = no minimum)");
  }
  if (typeof artifact.max_chars !== "number" || !Number.isInteger(artifact.max_chars) || artifact.max_chars < 0) {
    fail(rel, "artifact.max_chars must be an integer >= 0 (0 = no maximum)");
  }
  if (
    artifact.min_chars > 0 &&
    artifact.max_chars > 0 &&
    artifact.min_chars > artifact.max_chars
  ) {
    fail(rel, "artifact.min_chars must be <= max_chars when both are positive");
  }
  if (!Array.isArray(artifact.checklist) || artifact.checklist.length === 0) {
    fail(rel, "artifact.checklist must be a non-empty string array");
  } else {
    artifact.checklist.forEach((c, i) => {
      if (!nonEmptyString(c)) fail(rel, `artifact.checklist[${i}] must be non-empty`);
    });
  }
  if (artifact.exemplar_ja !== undefined && !nonEmptyString(artifact.exemplar_ja)) {
    fail(rel, "artifact.exemplar_ja must be non-empty when set");
  }
}

function validateMeta(meta, rel) {
  if (!isObj(meta)) {
    fail(rel, "_meta must be an object");
    return;
  }
  if (!nonEmptyString(meta.source)) fail(rel, "_meta.source required");
  if (!nonEmptyString(meta.license)) fail(rel, "_meta.license required");
  // Content contract (project-manifest Constraint 2 / corpus-storage): every
  // learner-facing L1 row carries validated_by so reviewed content is accountable.
  if (!nonEmptyString(meta.validated_by)) {
    fail(rel, "_meta.validated_by required (non-empty string)");
  }
  if (meta.validator_score !== undefined && typeof meta.validator_score !== "number") {
    fail(rel, "_meta.validator_score must be a number when set");
  }
}

function validateUnit(data, absPath) {
  const rel = path.relative(rootDir, absPath);
  const base = path.basename(absPath, ".json");
  const parent = path.basename(path.dirname(absPath));

  if (!isObj(data)) {
    fail(rel, "root must be a JSON object");
    return;
  }

  for (const key of Object.keys(data)) {
    if (!ALLOWED_TOP.has(key)) fail(rel, `unknown top-level key '${key}'`);
  }
  for (const key of REQUIRED_TOP) {
    if (!Object.prototype.hasOwnProperty.call(data, key)) {
      fail(rel, `missing required key '${key}'`);
    }
  }

  if (data.schema_version !== 1) {
    fail(rel, `schema_version must be 1 (got ${JSON.stringify(data.schema_version)})`);
  }
  if (!nonEmptyString(data.id)) fail(rel, "id must be a non-empty string");
  else {
    if (data.id !== base) {
      fail(rel, `id '${data.id}' must match filename stem '${base}'`);
    }
    if (ids.has(data.id)) {
      fail(rel, `duplicate id '${data.id}' (also ${ids.get(data.id)})`);
    } else {
      ids.set(data.id, rel);
    }
  }

  if (!STAGES_V1.has(data.stage)) {
    fail(rel, `stage '${data.stage}' not in v1 allowlist {e5-6}`);
  } else if (parent !== "kokugo" && parent !== data.stage) {
    // Allow flat layout under kokugo/ OR stage subdir matching stage field.
    fail(rel, `parent directory '${parent}' must be 'kokugo' or match stage '${data.stage}'`);
  }

  if (!nonEmptyString(data.title_ja)) fail(rel, "title_ja required");
  if (!GENRES.has(data.genre)) fail(rel, `invalid genre '${data.genre}'`);

  if (!Array.isArray(data.objectives) || data.objectives.length === 0) {
    fail(rel, "objectives must be a non-empty string array");
  } else {
    data.objectives.forEach((o, i) => {
      if (!nonEmptyString(o)) fail(rel, `objectives[${i}] must be non-empty`);
    });
  }

  if (
    typeof data.estimated_minutes !== "number" ||
    !Number.isInteger(data.estimated_minutes) ||
    data.estimated_minutes < 1
  ) {
    fail(rel, "estimated_minutes must be a positive integer");
  }

  const paragraphCount = validateBlocks(data.text, rel, "text");
  const plainText = Array.isArray(data.text) ? plainTextFromBlocks(data.text) : "";

  if (!isObj(data.support) || !SUPPORT.has(data.support.default_profile)) {
    fail(rel, "support.default_profile must be heavy|n3|standard|none");
  }

  if (!Array.isArray(data.tasks) || data.tasks.length === 0) {
    fail(rel, "tasks must be a non-empty array");
  } else {
    const taskIds = new Set();
    const kindsSeen = new Set();
    data.tasks.forEach((task, i) => {
      validateTask(task, rel, i, plainText, paragraphCount, taskIds);
      if (isObj(task) && task.kind) kindsSeen.add(task.kind);
    });
    // Soft guidance: PoC-quality units should cover all v1 kinds; warn via fail only if zero kinds valid.
    // Require at least one task; full four-kind coverage is recommended for PoC (JS-130) but not hard-failed here.
  }

  if (data.artifact !== undefined) validateArtifact(data.artifact, rel);
  if (data.classmates !== undefined && !Array.isArray(data.classmates)) {
    fail(rel, "classmates must be an array when set");
  }

  validateMeta(data._meta, rel);
}

walk(kokugoRoot);

if (files.length === 0) {
  console.error(`lint-kokugo: no .json units under ${path.relative(rootDir, kokugoRoot) || kokugoRoot}`);
  process.exit(2);
}

for (const file of files.sort()) {
  let data;
  try {
    data = JSON.parse(fs.readFileSync(file, "utf8"));
  } catch (err) {
    fail(path.relative(rootDir, file), `invalid JSON: ${err.message}`);
    // Fail-fast on unparseable units so later units are not scanned under a
    // partially trusted inventory (and so CI does not hide the first parse error).
    process.exit(exitCode);
  }
  validateUnit(data, file);
}

if (exitCode === 0) {
  console.log(`lint-kokugo: ok (${files.length} unit${files.length === 1 ? "" : "s"})`);
}
process.exit(exitCode);
NODE
