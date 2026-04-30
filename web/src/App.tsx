import { useState } from "react";
import { VocabTab } from "./tabs/VocabTab";
import { KanjiTab } from "./tabs/KanjiTab";
import { SentenceTab } from "./tabs/SentenceTab";
import { GrammarTab } from "./tabs/GrammarTab";
import { QuizTab } from "./tabs/QuizTab";

type Tab = "quiz" | "grammar" | "vocab" | "kanji" | "sentence";

const tabs: { id: Tab; label: string }[] = [
  { id: "quiz", label: "練習題" },
  { id: "grammar", label: "文法" },
  { id: "vocab", label: "單字" },
  { id: "kanji", label: "漢字" },
  { id: "sentence", label: "例句" },
];

export function App() {
  const [active, setActive] = useState<Tab>("quiz");

  return (
    <div className="min-h-screen">
      <header className="bg-white border-b border-slate-200">
        <div className="max-w-4xl mx-auto px-4 py-4 flex items-center justify-between">
          <h1 className="text-xl font-semibold tracking-tight">
            日本語学習 <span className="text-slate-400 text-sm">— japanese-site</span>
          </h1>
          <span className="text-xs text-slate-500">M3 開發預覽</span>
        </div>
        <nav className="max-w-4xl mx-auto px-4 flex gap-1 -mb-px">
          {tabs.map((t) => (
            <button
              key={t.id}
              onClick={() => setActive(t.id)}
              className={
                "px-4 py-2 text-sm font-medium border-b-2 transition-colors " +
                (active === t.id
                  ? "border-blue-600 text-blue-600"
                  : "border-transparent text-slate-600 hover:text-slate-900")
              }
            >
              {t.label}
            </button>
          ))}
        </nav>
      </header>
      <main className="max-w-4xl mx-auto px-4 py-8">
        {active === "quiz" && <QuizTab />}
        {active === "grammar" && <GrammarTab />}
        {active === "vocab" && <VocabTab />}
        {active === "kanji" && <KanjiTab />}
        {active === "sentence" && <SentenceTab />}
      </main>
      <footer className="text-center text-xs text-slate-400 py-6">
        資料來源 JMdict / KANJIDIC2 / Tatoeba (CC-BY-SA / CC-BY) ・ JLPT 標註 (MIT)
      </footer>
    </div>
  );
}
