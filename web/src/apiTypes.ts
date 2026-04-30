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
  explanation_zh: string;
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
}

export interface GrammarStat {
  grammar_point: string;
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
  by_error_class: ErrorClassStat[];
  recent_wrong: RecentWrongAttempt[];
}

export interface NextQuestionOpts {
  jlpt?: string;
  grammar?: string;
  exclude?: string[];
}

/** Contract for any API client implementation. The default `httpApi` makes
 *  fetch() calls; tests can substitute a mock that satisfies this shape. */
export interface Api {
  searchVocab(q: string, jlpt?: string): Promise<{ results: VocabRow[]; count: number }>;
  getKanji(ch: string): Promise<Kanji>;
  randomSentence(jlpt?: string): Promise<Sentence>;
  listGrammar(jlpt?: string): Promise<{ points: GrammarPoint[]; count: number }>;
  getGrammar(slug: string): Promise<GrammarPoint>;
  nextQuestion(opts?: NextQuestionOpts): Promise<Question>;
  answer(question_id: string, answer: string): Promise<GradeResult>;
  stats(days?: number): Promise<Stats>;
}
