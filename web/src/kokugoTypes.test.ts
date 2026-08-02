import { readFileSync } from "node:fs";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";
import { describe, expect, it } from "vitest";
import {
  KOKUGO_ARTIFACT_KINDS,
  KOKUGO_GENRES,
  KOKUGO_SCHEMA_VERSION,
  KOKUGO_SKILLS,
  KOKUGO_STAGES_V1,
  KOKUGO_SUPPORT_PROFILES,
  KOKUGO_TASK_KINDS_V1,
  type KokugoUnit,
} from "./kokugoTypes";

const repoRoot = join(dirname(fileURLToPath(import.meta.url)), "../..");
const fixturePath = join(
  repoRoot,
  "server/data/corpus/kokugo/e5-6/library-use.json",
);

describe("kokugoTypes / corpus fixture", () => {
  // Behavior: Kokugo type exports and the library-use L1 fixture stay aligned
  // with ADR-0005 closed v1 enums and the four task kinds.
  //
  // Steps:
  // 1. Assert schema version and closed enum membership for stage/support/tasks.
  // 2. Load server/data/corpus/kokugo/e5-6/library-use.json as KokugoUnit.
  // 3. Assert stage, task kinds, Block text, and _meta provenance fields.
  it("exports v1 closed enums from ADR-0005", () => {
    expect(KOKUGO_SCHEMA_VERSION).toBe(1);
    expect([...KOKUGO_STAGES_V1]).toEqual(["e5-6"]);
    expect([...KOKUGO_SUPPORT_PROFILES]).toEqual([
      "heavy",
      "n3",
      "standard",
      "none",
    ]);
    expect([...KOKUGO_TASK_KINDS_V1]).toEqual([
      "predict",
      "evidence-highlight",
      "paragraph-role",
      "summary-choice",
    ]);
    expect(KOKUGO_GENRES.length).toBe(4);
    expect(KOKUGO_SKILLS.length).toBe(6);
    expect([...KOKUGO_ARTIFACT_KINDS]).toEqual(["short-proposal", "summary"]);
  });

  it("library-use fixture is assignable to KokugoUnit shape", () => {
    const unit = JSON.parse(readFileSync(fixturePath, "utf8")) as KokugoUnit;
    expect(unit.schema_version).toBe(KOKUGO_SCHEMA_VERSION);
    expect(unit.id).toBe("library-use");
    expect(unit.stage).toBe("e5-6");
    expect(unit.tasks.map((t) => t.kind).sort()).toEqual(
      [...KOKUGO_TASK_KINDS_V1].sort(),
    );
    expect(unit.text.some((b) => b.kind === "paragraph")).toBe(true);
    expect(unit._meta.source).toBeTruthy();
    expect(unit._meta.license).toBeTruthy();
  });
});
