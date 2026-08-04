import { readFileSync } from "node:fs";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";
import { describe, expect, expectTypeOf, it } from "vitest";
import type { Block } from "./apiTypes";
import {
  classmatesFor,
  classmatesForArtifact,
  classmatesForRevise,
  classmatesForTask,
  classmatesForWritingDone,
  countJaChars,
  shouldShowTaskClassmates,
  KOKUGO_ARTIFACT_KINDS,
  KOKUGO_GENRES,
  KOKUGO_SCHEMA_VERSION,
  KOKUGO_SKILLS,
  KOKUGO_STAGES_V1,
  KOKUGO_SUPPORT_PROFILES,
  KOKUGO_TASK_KINDS_V1,
  type KokugoMeta,
  type KokugoSupport,
  type KokugoTask,
  type KokugoUnit,
} from "./kokugoTypes";

const repoRoot = join(dirname(fileURLToPath(import.meta.url)), "../..");
const fixturePath = join(
  repoRoot,
  "server/data/corpus/kokugo/e5-6/library-use.json",
);

/** Minimal Block for compile-time fixtures. */
const sampleParagraph: Block = {
  kind: "paragraph",
  tokens: [{ t: "text", v: "例" }],
};

const sampleMeta: KokugoMeta = {
  source: "curated",
  license: "CC0-1.0",
  validated_by: "fixture-v1",
};

const sampleSupport: KokugoSupport = { default_profile: "standard" };

function baseUnit(tasks: KokugoTask[]): KokugoUnit {
  return {
    id: "fixture-unit",
    schema_version: KOKUGO_SCHEMA_VERSION,
    stage: "e5-6",
    title_ja: "題",
    genre: "expository",
    objectives: ["o1"],
    estimated_minutes: 10,
    text: [sampleParagraph],
    support: sampleSupport,
    tasks,
    _meta: sampleMeta,
  };
}

const predictTask = {
  id: "t-predict",
  skill: "reading.predict" as const,
  kind: "predict" as const,
  payload: {
    prompt_ja: "次は？",
    choices: [
      { id: "a", text_ja: "A" },
      { id: "b", text_ja: "B" },
    ],
  },
};

const evidenceTask = {
  id: "t-ev",
  skill: "reading.locate-evidence" as const,
  kind: "evidence-highlight" as const,
  payload: {
    prompt_ja: "根拠は？",
    gold_quotes: ["例"],
  },
};

const roleTask = {
  id: "t-role",
  skill: "reading.structure" as const,
  kind: "paragraph-role" as const,
  payload: {
    prompt_ja: "役割は？",
    roles: ["問題", "原因"],
    gold_by_paragraph_index: ["問題"],
  },
};

const summaryTask = {
  id: "t-sum",
  skill: "reading.summary" as const,
  kind: "summary-choice" as const,
  payload: {
    prompt_ja: "要約は？",
    choices: [
      { id: "s1", text_ja: "S1" },
      { id: "s2", text_ja: "S2" },
    ],
    correct_id: "s1",
  },
};

/** Runtime shape guard — fails if required unit/task fields are missing. */
function assertKokugoUnit(value: unknown): asserts value is KokugoUnit {
  if (value === null || typeof value !== "object") {
    throw new Error("unit must be an object");
  }
  const u = value as Record<string, unknown>;
  for (const key of [
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
  ] as const) {
    if (!(key in u)) {
      throw new Error(`missing required field: ${key}`);
    }
  }
  if (u.schema_version !== KOKUGO_SCHEMA_VERSION) {
    throw new Error(`schema_version must be ${KOKUGO_SCHEMA_VERSION}`);
  }
  if (!Array.isArray(u.tasks) || u.tasks.length === 0) {
    throw new Error("tasks must be a non-empty array");
  }
  const meta = u._meta as Record<string, unknown>;
  if (
    typeof meta?.source !== "string" ||
    typeof meta?.license !== "string" ||
    typeof meta?.validated_by !== "string" ||
    !meta.validated_by
  ) {
    throw new Error("_meta.source/license/validated_by required");
  }
  for (const raw of u.tasks as unknown[]) {
    if (raw === null || typeof raw !== "object") {
      throw new Error("task must be an object");
    }
    const t = raw as Record<string, unknown>;
    if (typeof t.id !== "string" || typeof t.kind !== "string") {
      throw new Error("task id/kind required");
    }
    if (t.payload === null || typeof t.payload !== "object") {
      throw new Error(`task ${t.id} payload required`);
    }
    const payload = t.payload as Record<string, unknown>;
    switch (t.kind) {
      case "predict":
        if (typeof payload.prompt_ja !== "string" || !Array.isArray(payload.choices)) {
          throw new Error(`task ${t.id} needs prompt_ja + choices`);
        }
        break;
      case "summary-choice":
        if (
          typeof payload.prompt_ja !== "string" ||
          !Array.isArray(payload.choices) ||
          typeof payload.correct_id !== "string"
        ) {
          throw new Error(`task ${t.id} needs prompt_ja + choices + correct_id`);
        }
        break;
      case "evidence-highlight":
        if (typeof payload.prompt_ja !== "string" || !Array.isArray(payload.gold_quotes)) {
          throw new Error(`task ${t.id} needs prompt_ja + gold_quotes`);
        }
        break;
      case "paragraph-role":
        if (
          typeof payload.prompt_ja !== "string" ||
          !Array.isArray(payload.roles) ||
          !Array.isArray(payload.gold_by_paragraph_index)
        ) {
          throw new Error(`task ${t.id} needs prompt_ja + roles + gold_by_paragraph_index`);
        }
        break;
      default:
        throw new Error(`unknown task kind: ${String(t.kind)}`);
    }
  }
}

describe("kokugoTypes / corpus fixture", () => {
  // Behavior: Kokugo type exports and the library-use L1 fixture stay aligned
  // with ADR-0005 closed v1 enums and the four task kinds.
  //
  // Steps:
  // 1. Assert schema version and closed enum membership for stage/support/tasks.
  // 2. Load server/data/corpus/kokugo/e5-6/library-use.json and assert shape.
  // 3. Compile-time positive/negative fixtures for required fields and discriminants.
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

  it("library-use fixture satisfies KokugoUnit runtime contract", () => {
    const raw: unknown = JSON.parse(readFileSync(fixturePath, "utf8"));
    assertKokugoUnit(raw);
    expect(raw.schema_version).toBe(KOKUGO_SCHEMA_VERSION);
    expect(raw.id).toBe("library-use");
    expect(raw.stage).toBe("e5-6");
    expect(raw.tasks.map((t) => t.kind).sort()).toEqual(
      [...KOKUGO_TASK_KINDS_V1].sort(),
    );
    expect(raw.text.some((b) => b.kind === "paragraph")).toBe(true);
    expect(raw._meta.validated_by).toBeTruthy();
    // JS-134 classmates present and filterable.
    expect(Array.isArray(raw.classmates) && raw.classmates.length).toBeGreaterThan(0);
    expect(classmatesForTask(raw, "evidence-1").length).toBeGreaterThan(0);
    expect(classmatesForTask(raw, "summary-1").length).toBeGreaterThan(0);
    expect(classmatesForArtifact(raw).length).toBeGreaterThan(0);
    expect(classmatesForRevise(raw).length).toBeGreaterThan(0);
    expect(classmatesForTask(raw, "no-such").length).toBe(0);
    expect(classmatesFor(raw, { kind: "revise" })).toEqual(classmatesForRevise(raw));
    expect(classmatesForWritingDone(raw).length).toBe(
      classmatesForArtifact(raw).length + classmatesForRevise(raw).length,
    );
    expect(countJaChars("  あいう  ")).toBe(3);
    expect(shouldShowTaskClassmates("artifact")).toBe(true);
    expect(shouldShowTaskClassmates("revise")).toBe(false);
  });

  it("compile-time: positive unit with all four task discriminants", () => {
    const unit = baseUnit([predictTask, evidenceTask, roleTask, summaryTask]);
    expectTypeOf(unit).toMatchTypeOf<KokugoUnit>();
    expectTypeOf(unit.tasks[0]).toMatchTypeOf<KokugoTask>();
    expectTypeOf(unit._meta).toMatchTypeOf<KokugoMeta>();
    expectTypeOf(unit._meta.validated_by).toBeString();
  });

  it("compile-time: missing required unit fields is not KokugoUnit", () => {
    expectTypeOf({
      schema_version: KOKUGO_SCHEMA_VERSION,
      stage: "e5-6" as const,
      // id omitted
    }).not.toMatchTypeOf<KokugoUnit>();
    expectTypeOf({
      id: "x",
      schema_version: KOKUGO_SCHEMA_VERSION,
      stage: "e5-6" as const,
      title_ja: "t",
      genre: "story" as const,
      objectives: [] as string[],
      estimated_minutes: 1,
      text: [sampleParagraph],
      support: sampleSupport,
      tasks: [predictTask],
      _meta: { source: "s", license: "l" }, // validated_by missing
    }).not.toMatchTypeOf<KokugoUnit>();
  });

  it("compile-time: task payload discriminants reject wrong payloads", () => {
    expectTypeOf({
      id: "bad",
      skill: "reading.predict" as const,
      kind: "predict" as const,
      // gold_quotes belongs to evidence-highlight, not predict
      payload: { prompt_ja: "x", gold_quotes: [] as string[] },
    }).not.toMatchTypeOf<KokugoTask>();
    expectTypeOf(predictTask).toMatchTypeOf<KokugoTask>();
    expectTypeOf(summaryTask).toMatchTypeOf<KokugoTask>();
  });

  it("runtime: assertKokugoUnit rejects missing validated_by and bad tasks", () => {
    // Intentionally invalid shapes — cast to unknown so tsc does not reject the fixture.
    const missingMeta: unknown = {
      ...baseUnit([predictTask]),
      _meta: { source: "s", license: "l" },
    };
    expect(() => assertKokugoUnit(missingMeta)).toThrow(/validated_by/);

    const badTask: unknown = {
      ...baseUnit([predictTask]),
      tasks: [
        {
          id: "t",
          skill: "reading.predict",
          kind: "predict",
          payload: { gold_quotes: [] },
        },
      ],
    };
    expect(() => assertKokugoUnit(badTask)).toThrow(/prompt_ja|choices/);
  });
});
