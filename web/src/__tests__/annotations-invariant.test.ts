// @ts-nocheck
import { readFileSync } from "node:fs";
import { resolve } from "node:path";
import { describe, expect, it } from "vitest";
import { ANNOTATION_KINDS } from "../apiTypes";
import { LABELS } from "../components/EntryAnnotations";

function readKindsFile() {
  return readFileSync(resolve(process.cwd(), "../scripts/annotations-kinds.txt"), "utf8")
    .split(/\r?\n/)
    .map((line) => line.trim())
    .filter(Boolean);
}

describe("annotation kind invariants", () => {
  it("keeps API types, labels, lint allowlist, and ADR in sync", () => {
    const kinds = [...ANNOTATION_KINDS];
    expect(kinds).toEqual(readKindsFile());
    expect(Object.keys(LABELS)).toEqual(kinds);

    const adr = readFileSync(
      resolve(process.cwd(), "../docs/adr/0001-vocab-annotations-schema.md"),
      "utf8"
    );
    for (const kind of kinds) {
      expect(adr).toMatch(new RegExp(`\\b${kind}\\b`));
    }
  });
});
