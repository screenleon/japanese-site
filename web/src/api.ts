// API client. Uses Vite dev proxy so /api → :8080 in development.

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

export interface Question {
  id: number;
  kind: string;
  jlpt_level: string;
  grammar_point: string;
  prompt: string;
  hint?: string;
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

async function getJSON<T>(path: string): Promise<T> {
  const r = await fetch(path);
  if (!r.ok) throw new Error(`${r.status} ${r.statusText}`);
  return r.json();
}

async function postJSON<T>(path: string, body: unknown): Promise<T> {
  const r = await fetch(path, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  if (!r.ok) throw new Error(`${r.status} ${r.statusText}`);
  return r.json();
}

export const api = {
  searchVocab: (q: string, jlpt?: string) =>
    getJSON<{ results: VocabRow[]; count: number }>(
      `/api/vocab/search?q=${encodeURIComponent(q)}${jlpt ? `&jlpt=${jlpt}` : ""}`
    ),
  getKanji: (ch: string) => getJSON<Kanji>(`/api/kanji/${encodeURIComponent(ch)}`),
  randomSentence: (jlpt?: string) =>
    getJSON<Sentence>(`/api/sentence/random${jlpt ? `?jlpt=${jlpt}` : ""}`),
  listGrammar: (jlpt?: string) =>
    getJSON<{ points: GrammarPoint[]; count: number }>(
      `/api/grammar${jlpt ? `?jlpt=${jlpt}` : ""}`
    ),
  getGrammar: (slug: string) => getJSON<GrammarPoint>(`/api/grammar/${slug}`),
  nextQuestion: (opts: { jlpt?: string; grammar?: string; exclude?: number[] } = {}) => {
    const q = new URLSearchParams();
    if (opts.jlpt) q.set("jlpt", opts.jlpt);
    if (opts.grammar) q.set("grammar", opts.grammar);
    if (opts.exclude && opts.exclude.length) q.set("exclude", opts.exclude.join(","));
    const qs = q.toString();
    return getJSON<Question>(`/api/quiz/next${qs ? `?${qs}` : ""}`);
  },
  answer: (question_id: number, answer: string) =>
    postJSON<GradeResult>(`/api/quiz/answer`, { question_id, answer }),
  stats: (days?: number) =>
    getJSON<Stats>(`/api/quiz/stats${days ? `?days=${days}` : ""}`),
};

export interface Stats {
  total_attempts: number;
  correct: number;
  wrong: number;
  accuracy: number;
  by_grammar: { grammar_point: string; total: number; correct: number; accuracy: number }[];
  by_error_class: { grammar_point: string; error_class: string; count: number }[];
  recent_wrong: {
    question_id: number;
    grammar_point: string;
    prompt: string;
    user_answer: string;
    expected: string;
    created_at: string;
  }[];
}
