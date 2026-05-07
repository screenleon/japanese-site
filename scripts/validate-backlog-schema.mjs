#!/usr/bin/env node
import fs from "node:fs";
import path from "node:path";
import process from "node:process";

const ACTIVE_STATUSES = new Set(["todo", "doing", "blocked"]);
const VALID_MILESTONES = new Set(["M1", "M3", "M4", "DX"]);
const THEME_RE = /^[a-z0-9]+(-[a-z0-9]+)*$/;

function usage() {
  console.error("usage: node scripts/validate-backlog-schema.mjs [path/to/backlog.yml]");
  process.exit(2);
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
    throw new Error("project/backlog.yml has no items[] entries");
  }
  return items;
}

function backlogPathFromArgs(argv) {
  if (argv.length > 1) usage();
  if (argv.length === 1) return path.resolve(argv[0]);

  const rootDir = process.env.ROOT_DIR
    ? path.resolve(process.env.ROOT_DIR)
    : process.cwd();
  return path.join(rootDir, "project", "backlog.yml");
}

function main() {
  const backlogPath = backlogPathFromArgs(process.argv.slice(2));
  let yaml;
  try {
    yaml = fs.readFileSync(backlogPath, "utf8");
  } catch (error) {
    console.error(`validate-backlog-schema: cannot read ${backlogPath}: ${error.message}`);
    process.exit(2);
  }

  let items;
  try {
    items = parseBacklogItems(yaml);
  } catch (error) {
    console.error(`validate-backlog-schema: ${error.message}`);
    process.exit(2);
  }

  const errors = [];
  for (const item of items) {
    if (!ACTIVE_STATUSES.has(item.status)) continue;

    if (
      Object.prototype.hasOwnProperty.call(item, "milestone") &&
      !VALID_MILESTONES.has(item.milestone)
    ) {
      errors.push(
        `E-MILESTONE-ENUM: ${item.id} milestone='${item.milestone}' not in {M1,M3,M4,DX}`,
      );
    }

    if (
      Object.prototype.hasOwnProperty.call(item, "theme") &&
      (item.theme === null || item.theme === "" || !THEME_RE.test(item.theme))
    ) {
      errors.push(`E-THEME-SYNTAX: ${item.id} theme='${item.theme}' not lowercase-kebab-case`);
    }
  }

  if (errors.length > 0) {
    console.error(errors.join("\n"));
    process.exit(1);
  }
}

main();
