// Pure type declarations for the API layer. Components and tabs MUST import
// types from here (or from `api.ts` which re-exports them) and depend on the
// `Api` interface, never the concrete `httpApi`. See
// rules/domain/frontend-components.md → UI-003.
//
// Scope: M3 read/write surface only. M4 (LLM connector) will add free-form
// grading methods (e.g. `gradeFreeForm`); add to `Api` then.

export interface VocabRow {
  id: number;
  headword: string;
  reading: string;
  pos: string;
  gloss_en?: string;
  gloss_ja?: string;
  gloss_zh?: string;
  jlpt_level?: string;
  frequency_rank?: number;
  annotations?: Annotations;
  source: string;
  license: string;
  validated_by?: string;
}

export const ANNOTATION_KINDS = [
  "usage",
  "collocations",
  "particle_pairing",
  "synonym_diff",
  "mental_model",
  "mental_model_zh",
  "nuance_note",
  "furigana",
  "classifier",
] as const;

export type AnnotationKind = typeof ANNOTATION_KINDS[number];
export interface FuriganaPair {
  kanji: string;
  reading: string;
}

export interface FuriganaAnnotation {
  title_ja?: Token[];
  vocabulary?: FuriganaPair[];
}

export type BlockKind = "paragraph" | "list" | "callout";
export type Block =
  | { kind: "paragraph"; tokens: Token[] }
  | { kind: "list"; items: { tokens: Token[] }[] }
  | { kind: "callout"; tone?: "info" | "warn" | "tip"; tokens: Token[] };

export type TokenKind = "text" | "ruby" | "term";
export type Token =
  | { t: "text"; v: string }
  | { t: "ruby"; k: string; r: string }
  | { t: "term"; slug: string; kind: "vocab" | "grammar"; label: string };

export interface PatternRow {
  form: string;
  gloss_zh: string;
  notes_zh?: string;
}

export interface Meta {
  source: string;
  license: string;
  validated_by?: string;
  validator_score?: number;
}

export interface ClassifierContrastExample {
  use_this: string;
  use_alt: string;
  gloss_zh?: string;
}

export interface ClassifierContrast {
  with_pattern: string;
  with_slug?: string;
  rule_ja_blocks: Block[];
  rule_zh?: string;
  examples?: ClassifierContrastExample[];
}

export interface ClassifierRule {
  error_class?: string;
  if_answer_equals_any?: string[];
  if_answer_suffix_any?: string[];
  if_answer_contains_any?: string[];
  if_answer_not_contains_any?: string[];
  if_answer_dictionary_form?: boolean;
  default?: boolean;
  contrast?: ClassifierContrast | null;
}

export interface ClassifierAnnotation {
  rules: ClassifierContrast[];
}

export type AnnotationValue<K extends AnnotationKind = AnnotationKind> =
  K extends "furigana" ? FuriganaAnnotation
  : K extends "classifier" ? ClassifierAnnotation
  : string;

export type Annotations = {
  [K in AnnotationKind]?: AnnotationValue<K>;
};

export interface Kanji {
  id: number;
  character: string;
  onyomi?: string;
  kunyomi?: string;
  meaning_en?: string;
  meaning_ja?: string;
  meaning_zh?: string;
  jlpt_level?: string;
  grade?: number;
  stroke_count?: number;
  source: string;
  license: string;
}

export interface Sentence {
  id: number;
  text_ja: string;
  text_en?: string;
  text_zh?: string;
  jlpt_level?: string;
  source: string;
  license: string;
}

export interface GrammarPoint {
  slug: string;
  title_ja: string;
  title_zh: string;
  jlpt_level: string;
  schema_version: 2;
  pattern: PatternRow[];
  explanation_ja_blocks: Block[];
  explanation_zh: string;
  _meta: Meta;
  classifier_rules?: ClassifierRule[];
  related_slugs?: string[];
  annotations?: Annotations;
  audit_status?: "pre-redesign" | "post-dedup-naive";
}

export interface GrammarExample {
  id: number;
  text_ja: string;
  text_zh: string;
}

// Question.id is the deterministic id derived server-side from
// (slug | prompt | expected) — a 16-char hex string. See
// `server/internal/content/corpus/load.go`'s `QuestionID`.
//
// `payload` is reserved for non-cloze kinds (multiple-choice distractors,
// ordering chunks, listening audio refs). Cloze questions omit it. The
// shape is kind-specific and intentionally untyped at this layer; tabs
// that consume non-cloze kinds will narrow it locally.
export interface Question {
  id: string;
  content_type: QuizContentType;
  kind: string;
  jlpt_level: string;
  grammar_point: string;
  prompt: string;
  hint?: string;
  payload?: unknown;
}

export interface GradeResult {
  correct: boolean;
  user_answer: string;
  expected: string;
  explanation_zh: string;
  grammar_point: string;
  error_class?: string;
  suggested_next: string[];
  content_type?: string;
  item_detail_zh?: string;
}

export interface GrammarStat {
  grammar_point: string;
  total: number;
  correct: number;
  accuracy: number;
}

export interface VocabStat {
  vocab_item: string;
  total: number;
  correct: number;
  accuracy: number;
}

export interface ErrorClassStat {
  grammar_point: string;
  error_class: string;
  count: number;
}

export interface RecentWrongAttempt {
  question_id: string;
  grammar_point: string;
  prompt: string;
  user_answer: string;
  expected: string;
  created_at: string;
}

export interface Stats {
  total_attempts: number;
  correct: number;
  wrong: number;
  accuracy: number;
  by_grammar: GrammarStat[];
  by_vocab: VocabStat[];
  by_error_class: ErrorClassStat[];
  recent_wrong: RecentWrongAttempt[];
}

export interface DueCount {
  grammar: number;
  vocab: number;
}

export type QuizContentType = "grammar" | "vocab";

export interface NextQuestionOpts {
  type?: QuizContentType;
  jlpt?: string;
  grammar?: string;
  exclude?: string[];
}

export type ReadKey =
  | { type: "grammar"; slug: string }
  | { type: "vocab"; headword: string }
  | { type: "kanji"; character: string };

// Derived from ReadKey so a 4th content type only needs to extend the union.
export type ReadContentType = ReadKey["type"];

export function readKeyIdentifier(key: ReadKey): string {
  switch (key.type) {
    case "grammar":
      return key.slug;
    case "vocab":
      return key.headword;
    case "kanji":
      return key.character;
  }
}

export interface Capabilities {
  progress: boolean;
  history: boolean;
  quiz: boolean;
  sentence: boolean;
  /** Track B 国語教室 unit listing available (corpus present). */
  kokugo: boolean;
}

/** List-row for a Kokugo unit (JS-131). */
export interface KokugoUnitSummary {
  id: string;
  stage: string;
  title_ja: string;
  genre: string;
  estimated_minutes: number;
  task_count: number;
  has_artifact: boolean;
}

export interface KokugoUnitProgress {
  unit_key: string;
  stage: string;
  unit_id: string;
  status: "in_progress" | "completed";
  step: string;
  started_at: string;
  updated_at: string;
  completed_at?: string;
}

export interface KokugoTaskAttempt {
  id: number;
  unit_key: string;
  task_id: string;
  answer: unknown;
  correct?: boolean | null;
  feedback?: unknown;
  created_at: string;
}

export interface KokugoArtifactRow {
  unit_key: string;
  revision: number;
  body: string;
  checklist: unknown;
  /** Monotonic optimistic-concurrency token (starts at 1). */
  version: number;
  created_at: string;
  updated_at: string;
}

export interface KokugoUnitState {
  progress?: KokugoUnitProgress;
  attempts: KokugoTaskAttempt[];
  artifacts: KokugoArtifactRow[];
}

export interface KokugoGradeResult {
  correct?: boolean | null;
  explanation_ja: string;
  detail?: Record<string, unknown>;
}

export interface ProgressSummary {
  level?: string;
  content_type: ReadContentType;
  read: number;
  total: number;
  percent: number;
}

/** Contract for any API client implementation. The default `httpApi` makes
 *  fetch() calls; tests can substitute a mock that satisfies this shape. */
export interface Api {
  searchVocab(q: string, jlpt?: string): Promise<{ results: VocabRow[]; count: number; total: number }>;
  randomVocab(jlpt?: string): Promise<VocabRow>;
  getKanji(ch: string): Promise<Kanji>;
  randomSentence(jlpt?: string): Promise<Sentence>;
  listGrammar(jlpt?: string): Promise<{ points: GrammarPoint[]; count: number }>;
  randomGrammar(jlpt?: string): Promise<GrammarPoint>;
  getGrammar(slug: string): Promise<GrammarPoint>;
  getGrammarExamples?(slug: string): Promise<{ examples: GrammarExample[]; count: number }>;
  nextQuestion(opts?: NextQuestionOpts): Promise<Question>;
  answer(question_id: string, answer: string): Promise<GradeResult>;
  stats(days?: number): Promise<Stats>;
  getDueCount(): Promise<DueCount>;
  markRead(key: ReadKey): Promise<void>;
  getProgress(type: ReadContentType, level?: string): Promise<ProgressSummary>;
  getCapabilities(): Promise<Capabilities>;
  /** 国語教室 (Track B) — optional on staticApi (throws / empty). */
  listKokugoUnits?(): Promise<{ units: KokugoUnitSummary[]; count: number }>;
  getKokugoUnit?(stage: string, id: string): Promise<import("./kokugoTypes").KokugoUnit>;
  getKokugoUnitState?(stage: string, id: string): Promise<KokugoUnitState>;
  putKokugoProgress?(
    stage: string,
    id: string,
    body: { step: string; status?: "in_progress" | "completed" }
  ): Promise<KokugoUnitProgress>;
  submitKokugoTask?(
    stage: string,
    id: string,
    taskId: string,
    answer: unknown
  ): Promise<{ attempt: KokugoTaskAttempt; grade: KokugoGradeResult }>;
  saveKokugoArtifact?(
    stage: string,
    id: string,
    body: {
      revision: 0 | 1;
      body: string;
      checklist_checked: boolean[];
      /** Required when updating an existing revision row (CAS on version). */
      expected_version?: number;
    }
  ): Promise<{
    artifact: KokugoArtifactRow;
    grade: KokugoGradeResult;
    /** Server progress after the transactional artifact write (JS-132). */
    progress?: KokugoUnitProgress;
  }>;
}
