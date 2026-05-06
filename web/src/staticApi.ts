import { ApiError } from "./api";
import type {
  Api,
  DueCount,
  GradeResult,
  GrammarExample,
  GrammarPoint,
  Kanji,
  NextQuestionOpts,
  ProgressSummary,
  Question,
  ReadContentType,
  ReadKey,
  Sentence,
  Stats,
  VocabRow,
} from "./apiTypes";

const levels = ["N1", "N2", "N3", "N4", "N5"];
const jsonlCache = new Map<string, Promise<unknown[] | ApiError>>();
const jsonCache = new Map<string, Promise<unknown | ApiError>>();
type EmptyArrayFetchOptions = { on404: "empty-array" };
type ThrowFetchOptions = { on404?: "throw" };
type FetchOptions = EmptyArrayFetchOptions | ThrowFetchOptions;

export function resetStaticApiCachesForTest() {
  jsonCache.clear();
  jsonlCache.clear();
}

function dataPath(path: string): string {
  const base =
    (import.meta as ImportMeta & { env: { BASE_URL?: string } }).env.BASE_URL || "/";
  return `${base.replace(/\/$/, "")}/data/${path.replace(/^\//, "")}`;
}

async function fetchJSON<T>(path: string, options?: ThrowFetchOptions): Promise<T>;
async function fetchJSON<T extends readonly unknown[]>(
  path: string,
  options: EmptyArrayFetchOptions
): Promise<T>;
async function fetchJSON<T>(path: string, options: FetchOptions = {}): Promise<T> {
  const url = dataPath(path);
  const on404 = options.on404 ?? "throw";
  if (!jsonCache.has(url)) {
    jsonCache.set(url, fetchJSONRaw(url));
  }

  const result = await jsonCache.get(url)!;
  if (result instanceof ApiError) {
    if (result.status === 404 && on404 === "empty-array") {
      return [] as unknown as T;
    }
    if (result.status !== 404) jsonCache.delete(url);
    throw result;
  }
  return result as T;
}

async function fetchJSONL<T>(path: string, options?: ThrowFetchOptions): Promise<T[]>;
async function fetchJSONL<T extends readonly unknown[]>(
  path: string,
  options: EmptyArrayFetchOptions
): Promise<T>;
async function fetchJSONL<T>(
  path: string,
  options: FetchOptions = {}
): Promise<T[] | T> {
  const url = dataPath(path);
  const on404 = options.on404 ?? "throw";
  if (!jsonlCache.has(url)) {
    jsonlCache.set(url, fetchJSONLRaw(url));
  }

  const result = await jsonlCache.get(url)!;
  if (result instanceof ApiError) {
    if (result.status === 404 && on404 === "empty-array") {
      return [] as unknown as T;
    }
    if (result.status !== 404) jsonlCache.delete(url);
    throw result;
  }
  return result as T[];
}

async function fetchJSONRaw(url: string): Promise<unknown | ApiError> {
  try {
    const response = await fetch(url);
    if (!response.ok) return staticFetchError(response);
    try {
      return await response.json();
    } catch (error) {
      return parseFetchError(response.status, error);
    }
  } catch (error) {
    return networkFetchError(error);
  }
}

async function fetchJSONLRaw(url: string): Promise<unknown[] | ApiError> {
  try {
    const response = await fetch(url);
    if (!response.ok) return staticFetchError(response);
    const text = await response.text();
    try {
      return text
        .split(/\r?\n/)
        .map((line) => line.trim())
        .filter(Boolean)
        .map((line) => JSON.parse(line));
    } catch (error) {
      return parseFetchError(response.status, error);
    }
  } catch (error) {
    return networkFetchError(error);
  }
}

function staticFetchError(response: Response): ApiError {
  const code = response.status === 404 ? "not_found" : "http_error";
  return new ApiError(response.status, response.statusText, code);
}

function parseFetchError(status: number, cause: unknown): ApiError {
  const error = new ApiError(status, "Parse Error", "parse_error");
  (error as Error & { cause?: unknown }).cause = cause;
  return error;
}

function networkFetchError(cause: unknown): ApiError {
  const error = new ApiError(0, "Network Error", "network_error");
  (error as Error & { cause?: unknown }).cause = cause;
  return error;
}

function unsupported(): ApiError {
  return new ApiError(400, "Bad Request", "unsupported_in_static_mode");
}

function choose<T>(rows: T[]): T {
  if (rows.length === 0) {
    throw new ApiError(404, "Not Found", "not_found");
  }
  return rows[Math.floor(Math.random() * rows.length)];
}

function normalizeLevel(jlpt?: string): string[] {
  return jlpt ? [jlpt] : levels;
}

async function loadVocabLevels(jlpt?: string): Promise<VocabRow[]> {
  const groups = await Promise.all(
    normalizeLevel(jlpt).map(async (level) => {
      const rows = await fetchJSONL<Partial<VocabRow>>(`vocab/${level}.jsonl`);
      return rows.map((row, index) => ({
        id: row.id ?? stableID(level, index),
        headword: row.headword ?? "",
        reading: row.reading ?? "",
        pos: row.pos ?? "",
        jlpt_level: row.jlpt_level ?? level,
        source: row.source ?? "curated",
        license: row.license ?? "",
        ...row,
      })) as VocabRow[];
    })
  );
  return groups.flat();
}

async function loadKanjiLevels(): Promise<Kanji[]> {
  const groups = await Promise.all(
    levels.map(async (level) => {
      const rows = await fetchJSONL<Partial<Kanji>>(`kanji/${level}.jsonl`);
      return rows.map((row, index) => ({
        id: row.id ?? stableID(level, index),
        character: row.character ?? "",
        jlpt_level: row.jlpt_level ?? level,
        source: row.source ?? "curated",
        license: row.license ?? "",
        ...row,
      })) as Kanji[];
    })
  );
  return groups.flat();
}

function stableID(level: string, index: number): number {
  const levelOffset = Number(level.slice(1)) * 100000;
  return levelOffset + index + 1;
}

async function loadGrammarLevel(
  level: string,
  options?: ThrowFetchOptions
): Promise<GrammarPoint[]>;
async function loadGrammarLevel(
  level: string,
  options: EmptyArrayFetchOptions
): Promise<GrammarPoint[]>;
async function loadGrammarLevel(
  level: string,
  options: FetchOptions = {}
): Promise<GrammarPoint[]> {
  if (options.on404 === "empty-array") {
    return fetchJSON<GrammarPoint[]>(`grammar/${level}.json`, {
      on404: "empty-array",
    });
  }
  return fetchJSON<GrammarPoint[]>(`grammar/${level}.json`, options);
}

async function loadGrammarLevels(jlpt?: string): Promise<GrammarPoint[]> {
  if (jlpt) {
    return loadGrammarLevel(jlpt, { on404: "empty-array" });
  }

  const groups = await Promise.all(
    levels.map(async (level) => {
      try {
        return await loadGrammarLevel(level, { on404: "empty-array" });
      } catch (err) {
        console.warn("grammar level rollup fetch failed", { level, err });
        return [];
      }
    })
  );
  return groups.flat();
}

if (false) {
  // @ts-expect-error on404 empty-array is only valid for array-shaped JSON.
  void fetchJSON<GrammarPoint>("grammar/N1.json", { on404: "empty-array" });
  // @ts-expect-error on404 empty-array is only valid for JSONL array result types.
  void fetchJSONL<GrammarExample>("grammar-examples/sae.jsonl", { on404: "empty-array" });
}

export const staticApiTestHooks = { fetchJSON, fetchJSONL };

export const staticApi: Api = {
  async searchVocab(q: string, jlpt?: string) {
    const query = q.trim().toLocaleLowerCase();
    const rows = await loadVocabLevels(jlpt);
    const results = query
      ? rows.filter((row) =>
          `${row.headword} ${row.reading}`.toLocaleLowerCase().includes(query)
        )
      : rows;
    return { results, count: results.length, total: rows.length };
  },

  async randomVocab(jlpt?: string) {
    return choose(await loadVocabLevels(jlpt));
  },

  async getKanji(ch: string) {
    const rows = await loadKanjiLevels();
    const found = rows.find((row) => row.character === ch);
    if (!found) throw new ApiError(404, "Not Found", "not_found");
    return found;
  },

  async randomSentence(_jlpt?: string): Promise<Sentence> {
    throw unsupported();
  },

  async listGrammar(jlpt?: string) {
    const points = await loadGrammarLevels(jlpt);
    return { points, count: points.length };
  },

  async randomGrammar(jlpt?: string) {
    return choose(await loadGrammarLevels(jlpt));
  },

  async getGrammar(slug: string) {
    const all = await loadGrammarLevels();
    const found = all.find((p) => p.slug === slug);
    if (!found) throw new ApiError(404, "Not Found", "not_found");
    return found;
  },

  async getGrammarExamples(slug: string) {
    const examples = await fetchJSONL<GrammarExample[]>(
      `grammar-examples/${encodeURIComponent(slug)}.jsonl`,
      { on404: "empty-array" }
    );
    return { examples, count: examples.length };
  },

  async nextQuestion(_opts?: NextQuestionOpts): Promise<Question> {
    throw unsupported();
  },

  async answer(_question_id: string, _answer: string): Promise<GradeResult> {
    throw unsupported();
  },

  async stats(_days?: number): Promise<Stats> {
    throw unsupported();
  },

  async getDueCount(): Promise<DueCount> {
    throw unsupported();
  },

  async markRead(_key: ReadKey): Promise<void> {},

  async getProgress(type: ReadContentType, level?: string): Promise<ProgressSummary> {
    return {
      level: level ?? "",
      content_type: type,
      read: 0,
      total: 0,
      percent: 0,
    };
  },

  async getCapabilities() {
    return {
      progress: false,
      history: false,
      quiz: false,
      sentence: false,
    };
  },
};
