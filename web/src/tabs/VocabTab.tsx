import { useEffect, useState } from "react";
import { api, type VocabRow } from "../api";
import {
  ContentFirstLayout,
  formatDirectorySummary,
} from "../components/ContentFirstLayout";
import { EntryAnnotations } from "../components/EntryAnnotations";
import { LevelPicker } from "../components/LevelPicker";
import { useChineseVisibility } from "../chineseVisibility";
import { useReadTracking } from "../hooks/useReadTracking";

export function VocabTab({ onNavigateKanji }: { onNavigateKanji?: (character: string) => void }) {
  const [query, setQuery] = useState("");
  const [selectedLevel, setSelectedLevel] = useState("N5");
  const [rows, setRows] = useState<VocabRow[]>([]);
  const [focusRow, setFocusRow] = useState<VocabRow | null>(null);
  const [err, setErr] = useState("");
  const [loadingList, setLoadingList] = useState(false);
  const [loadingRandom, setLoadingRandom] = useState(false);
  // Mobile directory open state lives in the tab (UI-001); lg+ always shows the list.
  const [directoryOpen, setDirectoryOpen] = useState(false);
  const { visible: chineseVisible } = useChineseVisibility();

  useReadTracking(
    focusRow?.headword ? { type: "vocab", headword: focusRow.headword } : null
  );

  useEffect(() => {
    void loadLevel(selectedLevel);
    void drawRandom(selectedLevel);
    setDirectoryOpen(false);
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
      setFocusRow(await api.randomVocab(level));
    } catch (e) {
      setErr(String(e));
      setFocusRow(null);
    } finally {
      setLoadingRandom(false);
    }
  }

  async function search(e: React.FormEvent) {
    e.preventDefault();
    await loadLevel(selectedLevel, query);
    setDirectoryOpen(true);
  }

  function supportGloss(row: VocabRow) {
    return row.gloss_zh?.trim() || "繁中釋義待補";
  }

  function japaneseGloss(row: VocabRow) {
    return row.gloss_ja?.trim() || "日本語の説明は準備中です。";
  }

  function hasAnnotations(row: VocabRow) {
    return Boolean(
      row.annotations &&
        Object.entries(row.annotations).some(([kind, value]) => {
          if (typeof value === "string") return value.trim().length > 0;
          if (kind === "furigana" && value) {
            const furigana = value as NonNullable<
              NonNullable<VocabRow["annotations"]>["furigana"]
            >;
            return (
              (furigana.title_ja?.length ?? 0) + (furigana.vocabulary?.length ?? 0) > 0
            );
          }
          return false;
        })
    );
  }

  function VocabHeadword({ row, size = "list" }: { row: VocabRow; size?: "card" | "list" }) {
    const headword = row.headword.trim();
    const reading = row.reading.trim();
    const textClass = size === "card" ? "text-3xl font-semibold" : "text-lg font-medium";
    const rtClass = size === "card" ? "text-base text-slate-500" : "text-sm text-slate-500";

    const headwordContent =
      size !== "card" || !onNavigateKanji
        ? row.headword
        : Array.from(row.headword).map((character, index) =>
            /\p{Script=Han}/u.test(character) ? (
              <button
                key={`${character}-${index}`}
                type="button"
                onClick={() => onNavigateKanji(character)}
                className="rounded text-blue-700 underline decoration-blue-300 underline-offset-4 hover:bg-blue-50 focus:outline-none focus:ring-2 focus:ring-blue-500"
                aria-label={`漢字「${character}」を調べる`}
              >
                {character}
              </button>
            ) : (
              <span key={`${character}-${index}`}>{character}</span>
            )
          );

    if (!reading || headword === reading) {
      return <div className={textClass}>{headwordContent}</div>;
    }

    return (
      <ruby className={textClass}>
        {headwordContent}
        <rt className={rtClass}>{row.reading}</rt>
      </ruby>
    );
  }

  const directory = (
    <div className="space-y-3">
      <form onSubmit={search} className="flex flex-wrap gap-2">
        <input
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder={`${selectedLevel} 內搜尋：食べる / おはよう / ありがとう`}
          className="min-w-0 flex-1 rounded-md border border-slate-300 px-3 py-2 text-base"
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

      {err && <p className="text-sm text-red-600">{err}</p>}

      <ul
        data-testid="vocab-directory-list"
        className="divide-y divide-slate-200 rounded-md border border-slate-200 bg-white"
      >
        {rows.map((r) => {
          const selected = focusRow?.id === r.id;
          return (
            <li key={r.id}>
              <button
                type="button"
                onClick={() => {
                  setFocusRow(r);
                  setDirectoryOpen(false);
                }}
                className={
                  "flex w-full gap-3 p-3 text-left transition-colors " +
                  (selected ? "bg-blue-50" : "hover:bg-slate-50")
                }
              >
                <div className="w-28 flex-shrink-0 sm:w-32">
                  <VocabHeadword row={r} />
                </div>
                <div className="min-w-0 flex-1">
                  <div className="truncate text-sm text-slate-700">{japaneseGloss(r)}</div>
                  {chineseVisible && (
                    <div className="mt-0.5 truncate text-xs text-slate-500">
                      {supportGloss(r)}
                    </div>
                  )}
                  <div className="mt-1 text-xs text-slate-400">
                    {r.pos}
                    {r.jlpt_level && (
                      <span className="ml-2 inline-block rounded bg-blue-50 px-1.5 py-0.5 text-blue-700">
                        {r.jlpt_level}
                      </span>
                    )}
                  </div>
                </div>
              </button>
            </li>
          );
        })}
        {loadingList && (
          <li className="p-4 text-sm text-slate-500">単語を読み込み中...</li>
        )}
        {!loadingList && rows.length === 0 && (
          <li className="p-4 text-sm text-slate-400">この条件に合う単語はありません。</li>
        )}
      </ul>
    </div>
  );

  const content = (
    <section
      data-testid="vocab-focus-card"
      className="rounded-md border border-slate-200 bg-white p-5"
    >
      <div className="mb-4 flex flex-wrap items-center justify-between gap-3">
        <div>
          <h2 className="text-lg font-medium">{selectedLevel} 單字</h2>
          <p className="text-sm text-slate-500">
            先讀日文說明，再需要時打開中文。可從下方目錄選字，或抽下一個。
          </p>
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
      {focusRow ? (
        <div className="grid gap-4 sm:grid-cols-[220px_1fr]">
          <div>
            <VocabHeadword row={focusRow} size="card" />
          </div>
          <div className="text-sm leading-relaxed">
            <div>{japaneseGloss(focusRow)}</div>
            {hasAnnotations(focusRow) && (
              <div className="mt-4">
                <EntryAnnotations annotations={focusRow.annotations} />
              </div>
            )}
            {chineseVisible && (
              <div className="mt-5 border-t border-slate-200 pt-4">
                <h3 className="text-sm font-medium text-slate-700">中文說明</h3>
                <div className="mt-4 text-slate-600">
                  {focusRow.gloss_zh?.trim() || "繁中釋義待補"}
                </div>
              </div>
            )}
            <div className="mt-2 text-xs text-slate-400">
              {focusRow.pos}
              {focusRow.jlpt_level && <span className="ml-2">{focusRow.jlpt_level}</span>}
            </div>
          </div>
        </div>
      ) : (
        <p className="text-sm text-slate-500">
          {loadingRandom ? "単語を読み込み中..." : "この等級の単語はまだありません。"}
        </p>
      )}
    </section>
  );

  return (
    <section className="space-y-6">
      <LevelPicker
        selected={selectedLevel}
        onSelect={(level) => {
          setSelectedLevel(level);
          setQuery("");
        }}
        subtitle="單字學習"
      />

      <ContentFirstLayout
        directoryOpen={directoryOpen}
        onDirectoryOpenChange={setDirectoryOpen}
        directorySummary={formatDirectorySummary(rows.length)}
        content={content}
        directory={directory}
      />
    </section>
  );
}
