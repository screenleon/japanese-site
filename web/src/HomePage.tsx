import { useEffect, useState } from "react";
import { api } from "./api";

type InitialMode = "練習" | "測試";

interface HomePageProps {
  onStart: (initialMode: InitialMode) => void;
}

export function HomePage({ onStart }: HomePageProps) {
  const [grammarCount, setGrammarCount] = useState<number | null>(null);
  const [vocabCount, setVocabCount] = useState<number | null>(null);
  const [error, setError] = useState("");

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

  return (
    <main className="min-h-screen bg-slate-50">
      <section className="max-w-4xl mx-auto px-4 py-16 md:py-24">
        <div className="mb-10">
          <p className="text-sm font-medium text-blue-700 mb-3">japanese-site</p>
          <h1 className="text-4xl md:text-5xl font-semibold tracking-tight text-slate-950">
            日本語学習
          </h1>
          <p className="mt-4 max-w-2xl text-base leading-7 text-slate-600">
            用文法、單字與測驗建立穩定的日文練習節奏。
          </p>
        </div>

        <div className="grid gap-3 sm:grid-cols-2 mb-8">
          <CountPanel label="文法點" value={grammarCount} />
          <CountPanel label="單字" value={vocabCount} />
        </div>

        {error && <p className="mb-6 text-sm text-red-600">無法讀取首頁統計：{error}</p>}

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
      </section>
    </main>
  );
}

function CountPanel({ label, value }: { label: string; value: number | null }) {
  return (
    <div className="bg-white border border-slate-200 rounded-md p-5">
      <div className="text-3xl font-semibold text-slate-950">
        {value === null ? "..." : value}
      </div>
      <div className="mt-1 text-sm text-slate-500">{label}</div>
    </div>
  );
}
