// Kokugo (国語教室) content types — Track B per ADR-0005.
// Separate from QuizContentType / cloze loop. Corpus lives at
// server/data/corpus/kokugo/** and is validated by scripts/lint-kokugo.sh.

import type { Block, Meta } from "./apiTypes";

export const KOKUGO_SCHEMA_VERSION = 1 as const;

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

export interface KokugoTask {
  id: string;
  skill: KokugoSkill;
  kind: KokugoTaskKindV1;
  payload: KokugoTaskPayload;
  /** Optional deterministic rubric notes for future graders; free-form for now. */
  rubric?: Record<string, unknown>;
}

export interface KokugoArtifact {
  kind: KokugoArtifactKind;
  min_chars: number;
  max_chars: number;
  checklist: string[];
  exemplar_ja?: string;
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
  /** Phase 2 — allowed as array but not deeply validated in v1 lint. */
  classmates?: unknown[];
  _meta: Meta;
}
