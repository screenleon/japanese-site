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
const jsonlCache = new Map<string, Promise<unknown[]>>();
const jsonCache = new Map<string, Promise<unknown>>();

export function resetStaticApiCachesForTest() {
  jsonCache.clear();
  jsonlCache.clear();
}

function dataPath(path: string): string {
  const base =
    (import.meta as ImportMeta & { env: { BASE_URL?: string } }).env.BASE_URL || "/";
  return `${base.replace(/\/$/, "")}/data/${path.replace(/^\//, "")}`;
}

async function fetchJSON<T>(path: string): Promise<T> {
  const url = dataPath(path);
  if (!jsonCache.has(url)) {
    const promise = fetch(url)
      .then(async (response) => {
        if (!response.ok) throw staticFetchError(response, "not_found");
        return response.json();
      })
      .catch((error) => {
        jsonCache.delete(url);
        throw error;
      });
    jsonCache.set(url, promise);
  }
  return jsonCache.get(url) as Promise<T>;
}

async function fetchJSONL<T>(path: string): Promise<T[]> {
  const url = dataPath(path);
  if (!jsonlCache.has(url)) {
    const promise = fetch(url)
      .then(async (response) => {
        if (!response.ok) throw staticFetchError(response, "not_found");
        const text = await response.text();
        return text
          .split(/\r?\n/)
          .map((line) => line.trim())
          .filter(Boolean)
          .map((line) => JSON.parse(line));
      })
      .catch((error) => {
        jsonlCache.delete(url);
        throw error;
      });
    jsonlCache.set(url, promise);
  }
  return jsonlCache.get(url) as Promise<T[]>;
}

function staticFetchError(response: Response, code: string): ApiError {
  return new ApiError(response.status, response.statusText, code);
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

async function loadGrammarLevel(level: string): Promise<GrammarPoint[]> {
  return fetchJSON<GrammarPoint[]>(`grammar/${level}.json`);
}

async function loadGrammarLevels(jlpt?: string): Promise<GrammarPoint[]> {
  const groups = await Promise.all(
    normalizeLevel(jlpt).map((level) =>
      loadGrammarLevel(level).catch(() => [] as GrammarPoint[])
    )
  );
  return groups.flat();
}

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
    try {
      const examples = await fetchJSONL<GrammarExample>(
        `grammar-examples/${slug}.jsonl`
      );
      return { examples, count: examples.length };
    } catch {
      return { examples: [], count: 0 };
    }
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
