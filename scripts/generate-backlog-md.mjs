#!/usr/bin/env node
import fs from "node:fs";
import path from "node:path";
import process from "node:process";

/*
 * Renders the generated parts of BACKLOG.md from project/backlog.yml.
 *
 * Hybrid source-of-truth model:
 * - project/backlog.yml owns structured fields:
 *   status / priority / milestone / area / completed_at / depends_on / related /
 *   source / title.
 * - BACKLOG.md owns narrative bodies and the user-facing 主題 text in each
 *   `## JS-XXX — 主題` heading.
 * - BACKLOG.md also owns 首次記錄 dates. This renderer reads them from each
 *   section's `<!-- 首次記錄: YYYY-MM-DD -->` or
 *   `<!-- 首次記錄: backfilled YYYY-MM-DD -->` comment, falling back to the
 *   existing index table's 首次記錄 cell when older sections lack the comment.
 *
 * Only two surfaces are rewritten:
 * 1. the `## Index` markdown table block, and
 * 2. each entry's heading line suffix (`✅ YYYY-MM-DD` for done items).
 *
 * The legacy project/backlog.yml `done:` bootstrap entries (JS-D001..JS-D003)
 * intentionally have no BACKLOG.md sections and are not surfaced here.
 */

const rootDir = process.env.ROOT_DIR
  ? path.resolve(process.env.ROOT_DIR)
  : process.cwd();
const backlogYamlPath = path.join(rootDir, "project", "backlog.yml");
const backlogMdPath = path.join(rootDir, "BACKLOG.md");

function fail(message) {
  console.error(`generate-backlog-md: ${message}`);
  process.exit(1);
}

function unquote(value) {
  const trimmed = String(value ?? "").trim();
  if (
    (trimmed.startsWith('"') && trimmed.endsWith('"')) ||
    (trimmed.startsWith("'") && trimmed.endsWith("'"))
  ) {
    return trimmed.slice(1, -1);
  }
  return trimmed;
}

function parseScalar(value) {
  const trimmed = String(value ?? "").trim();
  if (trimmed.startsWith("[") && trimmed.endsWith("]")) {
    const inner = trimmed.slice(1, -1).trim();
    if (!inner) return [];
    return inner.split(",").map((part) => unquote(part));
  }
  return unquote(trimmed);
}

function parseBacklogItems(yaml) {
  const items = [];
  let inItems = false;
  let current = null;

  for (const line of yaml.split(/\r?\n/)) {
    if (/^items:\s*$/.test(line)) {
      inItems = true;
      continue;
    }
    if (inItems && /^[A-Za-z_][A-Za-z0-9_]*:\s*$/.test(line)) {
      break;
    }
    if (!inItems) continue;

    const itemMatch = line.match(/^  - id:\s*(.+?)\s*$/);
    if (itemMatch) {
      current = { id: parseScalar(itemMatch[1]) };
      items.push(current);
      continue;
    }
    if (!current) continue;

    const fieldMatch = line.match(/^    ([A-Za-z_][A-Za-z0-9_]*):(?:\s*(.*))?$/);
    if (!fieldMatch) continue;

    const [, key, rawValue = ""] = fieldMatch;
    if (rawValue === "") {
      current[key] = null;
      continue;
    }
    current[key] = parseScalar(rawValue);
  }

  if (items.length === 0) {
    fail("project/backlog.yml has no items[] entries");
  }
  return items;
}

function compareBacklogIds(a, b) {
  const parseId = (id) => {
    const match = String(id).match(/^JS-(\d+)([a-z]*)$/);
    if (!match) return { numeric: Number.MAX_SAFE_INTEGER, suffix: String(id) };
    return { numeric: Number(match[1]), suffix: match[2] };
  };
  const left = parseId(a);
  const right = parseId(b);
  if (left.numeric !== right.numeric) return left.numeric - right.numeric;
  return left.suffix.localeCompare(right.suffix);
}

function statusCell(item) {
  if (item.status === "done") {
    if (!item.completed_at) {
      fail(`${item.id} has status=done but no completed_at`);
    }
    return `✅ closed ${item.completed_at}`;
  }
  if (item.status === "doing") return "🟡 in_progress";
  return "🔵 active";
}

function headingLine(item, subject) {
  if (item.status === "done") {
    if (!item.completed_at) {
      fail(`${item.id} has status=done but no completed_at`);
    }
    return `## ${item.id} — ${subject} ✅ ${item.completed_at}`;
  }
  return `## ${item.id} — ${subject}`;
}

function extractSubject(heading, id) {
  const match = heading.match(new RegExp(`^## ${id} — (.*?)(?: ✅ \\d{4}-\\d{2}-\\d{2})?$`));
  if (!match) {
    fail(`${id} heading is malformed: ${heading}`);
  }
  return match[1];
}

function parseExistingIndexFirstRecorded(indexBlock) {
  const dates = new Map();
  for (const line of indexBlock.split(/\r?\n/)) {
    if (!line.startsWith("| JS-")) continue;
    const cells = line
      .split("|")
      .slice(1, -1)
      .map((cell) => cell.trim());
    if (cells.length < 5) continue;
    dates.set(cells[0], cells[4]);
  }
  return dates;
}

function extractFirstRecorded(section, id, fallbackDates) {
  const comment = section.match(
    /<!--\s*首次記錄:\s*(?:backfilled\s+)?(\d{4}-\d{2}-\d{2})\s*-->/,
  );
  if (comment) return comment[1];

  const fallback = fallbackDates.get(id);
  if (fallback) return fallback;
  fail(`${id} is missing 首次記錄 comment and has no existing index fallback`);
}

function parseSections(sectionsBlock) {
  const matches = [...sectionsBlock.matchAll(/^## (JS-[0-9]+[a-z]?) — .*$/gm)];
  if (matches.length === 0) fail("BACKLOG.md has no JS entry sections");

  const sections = new Map();
  for (let i = 0; i < matches.length; i += 1) {
    const start = matches[i].index;
    const end = i + 1 < matches.length ? matches[i + 1].index : sectionsBlock.length;
    const id = matches[i][1];
    sections.set(id, sectionsBlock.slice(start, end));
  }
  return sections;
}

function render() {
  const yaml = fs.readFileSync(backlogYamlPath, "utf8");
  const markdown = fs.readFileSync(backlogMdPath, "utf8");
  const items = parseBacklogItems(yaml).sort((a, b) => compareBacklogIds(a.id, b.id));

  const indexHeading = markdown.match(/^## Index$/m);
  if (!indexHeading) fail("BACKLOG.md is missing ## Index");
  const indexHeadingEnd = indexHeading.index + indexHeading[0].length;

  const separatorMatch = markdown.slice(indexHeadingEnd).match(/\n---\n/);
  if (!separatorMatch) fail("BACKLOG.md is missing --- separator after ## Index");
  const separatorStart = indexHeadingEnd + separatorMatch.index;
  const separatorEnd = separatorStart + separatorMatch[0].length;

  const preamble = markdown.slice(0, indexHeadingEnd);
  const existingIndexBlock = markdown.slice(indexHeadingEnd, separatorStart);
  const afterSeparator = markdown.slice(separatorEnd);
  const fallbackDates = parseExistingIndexFirstRecorded(existingIndexBlock);
  const sections = parseSections(afterSeparator);

  const renderedSections = [];
  const indexRows = [];

  for (const item of items) {
    const section = sections.get(item.id);
    if (!section) fail(`${item.id} exists in project/backlog.yml items[] but has no BACKLOG.md section`);

    const newlineIndex = section.indexOf("\n");
    const oldHeading = newlineIndex === -1 ? section : section.slice(0, newlineIndex);
    const body = newlineIndex === -1 ? "" : section.slice(newlineIndex);
    const subject = extractSubject(oldHeading, item.id);
    const firstRecorded = extractFirstRecorded(section, item.id, fallbackDates);

    indexRows.push(
      `| ${item.id} | ${statusCell(item)} | ${subject} | ${item.area ?? ""} | ${firstRecorded} | ${item.source ?? ""} |`,
    );
    renderedSections.push(`${headingLine(item, subject)}${body}`);
  }

  const ymlIds = new Set(items.map((item) => item.id));
  for (const id of sections.keys()) {
    if (!ymlIds.has(id)) {
      fail(`${id} has a BACKLOG.md section but is missing from project/backlog.yml items[]`);
    }
  }

  const renderedIndex = [
    "",
    "",
    "<!-- JS-D001..JS-D003 bootstrap entries live in project/backlog.yml `done:` only; not surfaced here. -->",
    "",
    "| # | Status | 主題 | 影響面 | 首次記錄 | Refs |",
    "|---|--------|------|--------|----------|------|",
    ...indexRows,
    "",
  ].join("\n");

  const output = `${preamble}${renderedIndex}\n---\n\n${renderedSections.join("")}`;
  fs.writeFileSync(backlogMdPath, output);
}

render();
