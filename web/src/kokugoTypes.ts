// Kokugo (国語教室) content types — Track B per ADR-0005.
// Separate from QuizContentType / cloze loop. Corpus lives at
// server/data/corpus/kokugo/** and is validated by scripts/lint-kokugo.sh.

import type { Block } from "./apiTypes";

export const KOKUGO_SCHEMA_VERSION = 1 as const;

/**
 * Kokugo L1 provenance — stricter than shared Meta: lint-kokugo requires
 * non-empty validated_by on every unit (project-manifest Constraint 2).
 */
export interface KokugoMeta {
  source: string;
  license: string;
  validated_by: string;
  validator_score?: number;
}

/** v1 authored content allowlist (ADR-0005). Other stages are reserved only. */
export const KOKUGO_STAGES_V1 = ["e5-6"] as const;
export type KokugoStageV1 = (typeof KOKUGO_STAGES_V1)[number];

/** Reserved stages — not authorable in v1 lint. */
export const KOKUGO_STAGES_RESERVED = [
  "e1-2",
  "e3-4",
  "j1",
  "j2",
  "j3",
  "h-modern",
  "h-culture",
] as const;

export const KOKUGO_SUPPORT_PROFILES = [
  "heavy",
  "n3",
  "standard",
  "none",
] as const;
export type SupportProfile = (typeof KOKUGO_SUPPORT_PROFILES)[number];

export const KOKUGO_GENRES = [
  "story",
  "expository",
  "opinion",
  "poetry",
] as const;
export type KokugoGenre = (typeof KOKUGO_GENRES)[number];

export const KOKUGO_SKILLS = [
  "reading.predict",
  "reading.locate-evidence",
  "reading.structure",
  "reading.summary",
  "writing.claim-reason",
  "writing.revision",
] as const;
export type KokugoSkill = (typeof KOKUGO_SKILLS)[number];

/** v1 task kinds (closed). Phase 2+ kinds are intentionally excluded. */
export const KOKUGO_TASK_KINDS_V1 = [
  "predict",
  "evidence-highlight",
  "paragraph-role",
  "summary-choice",
] as const;
export type KokugoTaskKindV1 = (typeof KOKUGO_TASK_KINDS_V1)[number];

export const KOKUGO_ARTIFACT_KINDS = ["short-proposal", "summary"] as const;
export type KokugoArtifactKind = (typeof KOKUGO_ARTIFACT_KINDS)[number];

export interface KokugoChoice {
  id: string;
  text_ja: string;
}

/** Pre-reading prediction — not strictly graded. */
export interface PredictPayload {
  prompt_ja: string;
  choices: KokugoChoice[];
  /** Optional free-text box alongside choices. */
  allow_free_text?: boolean;
}

/** Highlight / select gold quote(s) from the passage. */
export interface EvidenceHighlightPayload {
  prompt_ja: string;
  /** Exact surface strings that must appear in unit text (concatenated plain text). */
  gold_quotes: string[];
}

/** Assign a role to each paragraph in `text` (paragraph blocks only, in order). */
export interface ParagraphRolePayload {
  prompt_ja: string;
  roles: string[];
  /** One gold role per paragraph block in unit.text, same order. */
  gold_by_paragraph_index: string[];
}

export interface SummaryChoicePayload {
  prompt_ja: string;
  choices: KokugoChoice[];
  correct_id: string;
}

export type KokugoTaskPayload =
  | PredictPayload
  | EvidenceHighlightPayload
  | ParagraphRolePayload
  | SummaryChoicePayload;

export type KokugoTask =
  | {
      id: string;
      skill: KokugoSkill;
      kind: "predict";
      payload: PredictPayload;
      rubric?: Record<string, unknown>;
    }
  | {
      id: string;
      skill: KokugoSkill;
      kind: "evidence-highlight";
      payload: EvidenceHighlightPayload;
      rubric?: Record<string, unknown>;
    }
  | {
      id: string;
      skill: KokugoSkill;
      kind: "paragraph-role";
      payload: ParagraphRolePayload;
      rubric?: Record<string, unknown>;
    }
  | {
      id: string;
      skill: KokugoSkill;
      kind: "summary-choice";
      payload: SummaryChoicePayload;
      rubric?: Record<string, unknown>;
    };

export interface KokugoArtifact {
  kind: KokugoArtifactKind;
  /** 0 = no minimum (progressive writing). */
  min_chars: number;
  /** 0 = no maximum. */
  max_chars: number;
  checklist: string[];
  exemplar_ja?: string;
}

/**
 * Curated “classmate” sample (JS-134). Not multi-user social — L1 content only.
 * Revealed after the learner has completed the anchored step.
 */
export type KokugoClassmateRevealAfter =
  | { kind: "task"; task_id: string }
  | { kind: "artifact" }
  | { kind: "revise" };

export interface KokugoClassmate {
  id: string;
  name_ja: string;
  reveal_after: KokugoClassmateRevealAfter;
  text_ja: string;
  /** Short pedagogical focus label (optional). */
  focus_ja?: string;
}

export interface KokugoSupport {
  default_profile: SupportProfile;
}

export interface KokugoUnit {
  id: string;
  schema_version: typeof KOKUGO_SCHEMA_VERSION;
  stage: KokugoStageV1;
  title_ja: string;
  genre: KokugoGenre;
  objectives: string[];
  estimated_minutes: number;
  text: Block[];
  support: KokugoSupport;
  tasks: KokugoTask[];
  artifact?: KokugoArtifact;
  /** Curated peer samples (JS-134); omit when none. */
  classmates?: KokugoClassmate[];
  _meta: KokugoMeta;
}

/** Grapheme-ish length for Japanese UI counters (… spread, after trim). */
export function countJaChars(text: string): number {
  return [...text.trim()].length;
}

/** Single filter for curated classmates by reveal_after (JS-134). */
export function classmatesFor(
  unit: Pick<KokugoUnit, "classmates">,
  reveal: KokugoClassmateRevealAfter
): KokugoClassmate[] {
  return (unit.classmates ?? []).filter((c) => {
    if (c.reveal_after.kind !== reveal.kind) return false;
    if (reveal.kind === "task") {
      return (
        c.reveal_after.kind === "task" && c.reveal_after.task_id === reveal.task_id
      );
    }
    return true;
  });
}

/** Classmates anchored to a specific task id. */
export function classmatesForTask(
  unit: Pick<KokugoUnit, "classmates">,
  taskId: string
): KokugoClassmate[] {
  return classmatesFor(unit, { kind: "task", task_id: taskId });
}

/** Classmates revealed after draft (artifact) phase. */
export function classmatesForArtifact(
  unit: Pick<KokugoUnit, "classmates">
): KokugoClassmate[] {
  return classmatesFor(unit, { kind: "artifact" });
}

/** Classmates revealed after revision / on done. */
export function classmatesForRevise(
  unit: Pick<KokugoUnit, "classmates">
): KokugoClassmate[] {
  return classmatesFor(unit, { kind: "revise" });
}

/** Writing classmates shown on the done summary (draft + revise samples). */
export function classmatesForWritingDone(
  unit: Pick<KokugoUnit, "classmates">
): KokugoClassmate[] {
  return [...classmatesForArtifact(unit), ...classmatesForRevise(unit)];
}

/**
 * Whether the last task's classmates should sit above the step chrome.
 * Shared by resume + live submit so reveal rules stay in one place.
 */
export function shouldShowTaskClassmates(phase: string): boolean {
  return phase === "task" || phase === "read" || phase === "artifact";
}
