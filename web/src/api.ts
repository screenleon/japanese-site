// Concrete HTTP client implementing the `Api` interface from `apiTypes.ts`.
// Uses Vite dev proxy so /api → :8080 in development.
//
// UI-003 (rules/domain/frontend-components.md): tabs and shared components
// MUST depend on the `Api` interface. Do NOT import `httpApi` outside this
// file or future top-level wiring; the `api` value alias is preserved only
// so existing tab imports keep working until they're switched to inject the
// Api interface via context.

import type {
  Api,
  GrammarPoint,
  Kanji,
  NextQuestionOpts,
  Question,
  Sentence,
  Stats,
  GradeResult,
  VocabRow,
} from "./apiTypes";

export type * from "./apiTypes";

export class ApiError extends Error {
  readonly status: number;
  readonly code: string;

  constructor(status: number, statusText: string, code: string) {
    super(`${status} ${statusText}${code ? `: ${code}` : ""}`);
    this.name = "ApiError";
    this.status = status;
    this.code = code;
  }
}

export function isApiError(error: unknown, code?: string): error is ApiError {
  return error instanceof ApiError && (!code || error.code === code);
}

async function getJSON<T>(path: string): Promise<T> {
  const r = await fetch(path);
  if (!r.ok) throw await apiError(r);
  return r.json();
}

async function postJSON<T>(path: string, body: unknown): Promise<T> {
  const r = await fetch(path, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  if (!r.ok) throw await apiError(r);
  return r.json();
}

async function apiError(response: Response) {
  let code = "";
  try {
    const body = await response.json();
    if (body && typeof body.error === "string") {
      code = body.error;
    }
  } catch {
    // Non-JSON errors still carry HTTP status and status text.
  }
  return new ApiError(response.status, response.statusText, code);
}

export const httpApi: Api = {
  searchVocab: (q, jlpt) =>
    getJSON<{ results: VocabRow[]; count: number }>(
      `/api/vocab/search?q=${encodeURIComponent(q)}${jlpt ? `&jlpt=${jlpt}` : ""}`
    ),
  randomVocab: (jlpt) =>
    getJSON<VocabRow>(`/api/vocab/random${jlpt ? `?jlpt=${jlpt}` : ""}`),
  getKanji: (ch) => getJSON<Kanji>(`/api/kanji/${encodeURIComponent(ch)}`),
  randomSentence: (jlpt) =>
    getJSON<Sentence>(`/api/sentence/random${jlpt ? `?jlpt=${jlpt}` : ""}`),
  listGrammar: (jlpt) =>
    getJSON<{ points: GrammarPoint[]; count: number }>(
      `/api/grammar${jlpt ? `?jlpt=${jlpt}` : ""}`
    ),
  randomGrammar: (jlpt) =>
    getJSON<GrammarPoint>(`/api/grammar/random${jlpt ? `?jlpt=${jlpt}` : ""}`),
  getGrammar: (slug) => getJSON<GrammarPoint>(`/api/grammar/${slug}`),
  nextQuestion: (opts: NextQuestionOpts = {}) => {
    const q = new URLSearchParams();
    if (opts.jlpt) q.set("jlpt", opts.jlpt);
    if (opts.grammar) q.set("grammar", opts.grammar);
    if (opts.exclude && opts.exclude.length) q.set("exclude", opts.exclude.join(","));
    const qs = q.toString();
    return getJSON<Question>(`/api/quiz/next${qs ? `?${qs}` : ""}`);
  },
  answer: (question_id, answer) =>
    postJSON<GradeResult>(`/api/quiz/answer`, { question_id, answer }),
  stats: (days) =>
    getJSON<Stats>(`/api/quiz/stats${days ? `?days=${days}` : ""}`),
};

/** @deprecated import the `Api` interface and inject `httpApi` instead. */
export const api: Api = httpApi;
