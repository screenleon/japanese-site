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
  source: string;
  license: string;
  validated_by?: string;
}

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
  nuance_note?: string;
  related_slugs?: string[];
  explanation_ja?: string;
  explanation_zh: string;
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
}
