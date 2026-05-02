import { useEffect, useState } from "react";
import { api, type VocabRow } from "../api";
import { useReadTracking } from "../hooks/useReadTracking";

const levels = ["N5", "N4", "N3", "N2", "N1"];

export function VocabTab() {
  const [query, setQuery] = useState("");
  const [selectedLevel, setSelectedLevel] = useState("N5");
  const [rows, setRows] = useState<VocabRow[]>([]);
  const [randomRow, setRandomRow] = useState<VocabRow | null>(null);
  const [err, setErr] = useState("");
  const [loadingList, setLoadingList] = useState(false);
  const [loadingRandom, setLoadingRandom] = useState(false);

  useReadTracking("vocab", randomRow?.headword);

  useEffect(() => {
    void loadLevel(selectedLevel);
    void drawRandom(selectedLevel);
  }, [selectedLevel]);

  async function loadLevel(level: string, q = "") {
    setLoadingList(true);
    setErr("");
    try {
      const r = await api.searchVocab(q.trim(), level);
      setRows(r.results || []);
    } catch (e) {
      setErr(String(e));
      setRows([]);
    } finally {
      setLoadingList(false);
    }
  }

  async function drawRandom(level = selectedLevel) {
    setLoadingRandom(true);
    setErr("");
    try {
      setRandomRow(await api.randomVocab(level));
    } catch (e) {
      setErr(String(e));
      setRandomRow(null);
    } finally {
      setLoadingRandom(false);
    }
  }

  async function search(e: React.FormEvent) {
    e.preventDefault();
    await loadLevel(selectedLevel, query);
  }

  function supportGloss(row: VocabRow) {
    return row.gloss_zh?.trim() || "繁中釋義待補";
  }

  function japaneseGloss(row: VocabRow) {
    return row.gloss_ja?.trim() || "日本語の説明は準備中です。";
  }

  return (
    <section className="space-y-6">
      <div className="grid gap-3 sm:grid-cols-5">
        {levels.map((level) => (
          <button
            key={level}
            type="button"
            onClick={() => {
              setSelectedLevel(level);
              setQuery("");
            }}
            className={
              "rounded-md border px-4 py-3 text-left transition-colors " +
              (selectedLevel === level
                ? "border-blue-600 bg-blue-50 text-blue-700"
                : "border-slate-200 bg-white text-slate-700 hover:border-slate-300")
            }
          >
            <div className="text-lg font-semibold">{level}</div>
            <div className="text-xs text-slate-500">單字學習</div>
          </button>
        ))}
      </div>

      <section className="rounded-md border border-slate-200 bg-white p-5">
        <div className="mb-4 flex flex-wrap items-center justify-between gap-3">
          <div>
            <h2 className="text-lg font-medium">{selectedLevel} 隨機單字</h2>
            <p className="text-sm text-slate-500">每次抽一個單字，先看日文與讀音，再確認意思。</p>
          </div>
          <button
            type="button"
            onClick={() => void drawRandom()}
            disabled={loadingRandom}
            className="rounded-md bg-blue-600 px-4 py-2 text-sm text-white disabled:opacity-50"
          >
            {loadingRandom ? "抽取中" : "抽下一個"}
          </button>
        </div>
        {randomRow ? (
          <div className="grid gap-4 sm:grid-cols-[220px_1fr]">
            <div>
              <div className="text-3xl font-semibold">{randomRow.headword}</div>
              <div className="mt-1 text-base text-slate-500">{randomRow.reading}</div>
            </div>
            <div className="text-sm leading-relaxed">
              <div>{japaneseGloss(randomRow)}</div>
              <div className="mt-2 text-slate-600">{supportGloss(randomRow)}</div>
              <div className="mt-2 text-xs text-slate-400">
                {randomRow.pos}
                {randomRow.jlpt_level && <span className="ml-2">{randomRow.jlpt_level}</span>}
              </div>
            </div>
          </div>
        ) : (
          <p className="text-sm text-slate-500">
            {loadingRandom ? "単語を読み込み中..." : "この等級のランダム単語はまだありません。"}
          </p>
        )}
      </section>

      <section>
        <form onSubmit={search} className="mb-4 flex gap-2">
          <input
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder={`${selectedLevel} 內搜尋：食べる / おはよう / ありがとう`}
            className="flex-1 rounded-md border border-slate-300 px-3 py-2 text-base"
          />
          <button
            type="submit"
            disabled={loadingList}
            className="rounded-md bg-blue-600 px-4 py-2 text-white disabled:opacity-50"
          >
            搜尋
          </button>
          <button
            type="button"
            onClick={() => {
              setQuery("");
              void loadLevel(selectedLevel);
            }}
            className="rounded-md border border-slate-300 px-4 py-2 text-slate-700 hover:bg-slate-50"
          >
            顯示全部
          </button>
        </form>

        {err && <p className="mb-3 text-sm text-red-600">{err}</p>}
        <ul className="divide-y divide-slate-200 rounded-md border border-slate-200 bg-white">
          {rows.map((r) => (
            <li key={r.id} className="flex gap-4 p-4">
              <div className="w-32 flex-shrink-0">
                <div className="text-lg font-medium">{r.headword}</div>
                <div className="text-sm text-slate-500">{r.reading}</div>
              </div>
              <div className="min-w-0 flex-1">
                <div className="text-sm">{japaneseGloss(r)}</div>
                <div className="mt-1 text-xs text-slate-500">{supportGloss(r)}</div>
                <div className="mt-1 text-xs text-slate-400">
                  {r.pos}
                  {r.jlpt_level && (
                    <span className="ml-2 inline-block rounded bg-blue-50 px-1.5 py-0.5 text-blue-700">
                      {r.jlpt_level}
                    </span>
                  )}
                </div>
              </div>
            </li>
          ))}
          {loadingList && <li className="p-4 text-sm text-slate-500">単語を読み込み中...</li>}
          {!loadingList && rows.length === 0 && (
            <li className="p-4 text-sm text-slate-400">この条件に合う単語はありません。</li>
          )}
        </ul>
      </section>
    </section>
  );
}
