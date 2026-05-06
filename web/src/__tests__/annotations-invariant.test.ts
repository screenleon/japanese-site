import { readFileSync } from "node:fs";
import { dirname, resolve } from "node:path";
import { fileURLToPath } from "node:url";
import { describe, expect, it } from "vitest";
import { ANNOTATION_KINDS } from "../apiTypes";
import { LABELS } from "../components/EntryAnnotations";

const testDir = dirname(fileURLToPath(import.meta.url));
const repoRoot = resolve(testDir, "../../..");

function readKindsFile() {
  return readFileSync(resolve(repoRoot, "scripts/annotations-kinds.txt"), "utf8")
    .split(/\r?\n/)
    .map((line) => line.trim())
    .filter((line) => line && !line.startsWith("#"));
}

describe("annotation kind invariants", () => {
  it("keeps API types, labels, lint allowlist, and ADR in sync", () => {
    const kinds = [...ANNOTATION_KINDS];
    expect(kinds).toEqual(readKindsFile());
    expect(Object.keys(LABELS)).toEqual(kinds);

    const adr = readFileSync(
      resolve(repoRoot, "docs/adr/0001-vocab-annotations-schema.md"),
      "utf8"
    );
    for (const kind of kinds) {
      expect(adr).toMatch(new RegExp(`\\b${kind}\\b`));
    }
  });
});
