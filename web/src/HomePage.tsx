/// <reference types="vite/client" />

import { useEffect, useState } from "react";
import { api } from "./api";
import type { DueCount } from "./apiTypes";
import type { QuizInitialMode } from "./tabs/QuizTab";

type ContentTab = "grammar" | "vocab" | "kanji" | "kokugo";

interface HomePageProps {
  onStart: (initialMode: QuizInitialMode, initialTab?: ContentTab) => void;
  quizCapable: boolean;
  kokugoCapable?: boolean;
}

export function HomePage({ onStart, quizCapable, kokugoCapable = false }: HomePageProps) {
  const [grammarCount, setGrammarCount] = useState<number | null>(null);
  const [vocabCount, setVocabCount] = useState<number | null>(null);
  const [dueCount, setDueCount] = useState<DueCount | null>(null);
  const [error, setError] = useState("");
  // Read process.env first so test stubs (vi.stubEnv → process.env in node) win
  // over the import.meta.env path that Vite injects at build time.
  const deployMode =
    (globalThis as { process?: { env?: { VITE_DEPLOY_MODE?: string } } }).process
      ?.env?.VITE_DEPLOY_MODE ??
    (import.meta as ImportMeta & { env: { VITE_DEPLOY_MODE?: string } }).env
      .VITE_DEPLOY_MODE;
  const isStaticBuild = deployMode === "static";
  const showQuizControls = !isStaticBuild && quizCapable;

  useEffect(() => {
    let cancelled = false;
    Promise.all([api.listGrammar(), api.searchVocab("", undefined)])
      .then(([grammar, vocab]) => {
        if (cancelled) return;
        setGrammarCount(grammar.count);
        setVocabCount(vocab.total);
      })
      .catch((e) => {
        if (!cancelled) setError(String(e));
      });
    return () => {
      cancelled = true;
    };
  }, []);

  useEffect(() => {
    if (isStaticBuild || !quizCapable) return;
    let cancelled = false;
    api
      .getDueCount()
      .then((d) => {
        if (!cancelled) setDueCount(d);
      })
      .catch(() => {
        // Transient API failures simply hide the badge.
      });
    return () => {
      cancelled = true;
    };
  }, [isStaticBuild, quizCapable]);

  return (
    <main className="min-h-screen bg-slate-50">
      <section className="max-w-4xl mx-auto px-4 py-16 md:py-24">
        <div className="mb-10">
          <p className="text-sm font-medium text-blue-700 mb-3">japanese-site</p>
          <h1 className="text-4xl md:text-5xl font-semibold tracking-tight text-slate-950">
            日本語学習
          </h1>
          <p className="mt-4 max-w-2xl text-base leading-7 text-slate-600">
            {isStaticBuild
              ? "查閱文法說明、單字與漢字，隨時作為學習參考。"
              : "用文法、單字與測驗建立穩定的日文練習節奏。"}
          </p>
        </div>

        <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-4 mb-8">
          <NavCard
            label="文法"
            count={grammarCount}
            description="JLPT N5–N1 文法說明"
            onClick={() => onStart("練習", "grammar")}
          />
          <NavCard
            label="單字"
            count={vocabCount}
            description="JLPT N5–N1 詞彙與例句"
            onClick={() => onStart("練習", "vocab")}
          />
          <NavCard
            label="漢字"
            count={undefined}
            description="讀音・筆畫・部首"
            onClick={() => onStart("練習", "kanji")}
          />
          {kokugoCapable ? (
            <NavCard
              label="国語教室"
              count={undefined}
              description="讀→據→寫→改 の最小循環"
              onClick={() => onStart("練習", "kokugo")}
            />
          ) : null}
        </div>

        {showQuizControls ? (
          <>
            {dueCount !== null && (dueCount.grammar > 0 || dueCount.vocab > 0) && (
              <div className="mb-6 p-4 bg-amber-50 border border-amber-200 rounded-md flex items-center justify-between gap-4">
                <div>
                  <p className="text-sm font-medium text-amber-800">需要複習</p>
                  <p className="text-xs text-amber-700 mt-0.5">
                    文法 {dueCount.grammar} 題・單字 {dueCount.vocab} 題
                  </p>
                </div>
                <button
                  onClick={() => onStart("複習")}
                  className="shrink-0 px-4 py-2 rounded-md bg-amber-600 text-white text-sm font-medium hover:bg-amber-700"
                >
                  開始複習 ({dueCount.grammar + dueCount.vocab})
                </button>
              </div>
            )}

            {error && (
              <p className="mb-6 text-sm text-red-600">無法讀取首頁統計：{error}</p>
            )}

            <div className="flex flex-wrap gap-3">
              <button
                onClick={() => onStart("練習")}
                className="px-5 py-3 rounded-md bg-blue-600 text-white text-sm font-medium hover:bg-blue-700"
              >
                開始練習
              </button>
              <button
                onClick={() => onStart("測試")}
                className="px-5 py-3 rounded-md border border-slate-300 bg-white text-slate-800 text-sm font-medium hover:bg-slate-100"
              >
                開始測試
              </button>
            </div>
          </>
        ) : null}
      </section>
    </main>
  );
}

function NavCard({
  label,
  count,
  description,
  onClick,
}: {
  label: string;
  count?: number | null;
  description: string;
  onClick: () => void;
}) {
  return (
    <button
      type="button"
      onClick={onClick}
      className="w-full text-left bg-white border border-slate-200 rounded-md p-5 hover:border-slate-300 hover:shadow-sm transition"
    >
      <div className="flex items-baseline justify-between gap-3">
        <div className="text-2xl font-semibold text-slate-950">{label}</div>
        {count !== undefined && (
          <div data-testid="nav-count" className="text-3xl font-semibold text-slate-950">
            {count === null ? "..." : count}
          </div>
        )}
      </div>
      <div className="mt-2 text-sm text-slate-500">{description}</div>
    </button>
  );
}
