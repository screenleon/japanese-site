import { useEffect, useMemo, useState } from "react";
import { VocabTab } from "./tabs/VocabTab";
import { KanjiTab } from "./tabs/KanjiTab";
import { SentenceTab } from "./tabs/SentenceTab";
import { GrammarTab } from "./tabs/GrammarTab";
import { QuizTab, type QuizInitialMode } from "./tabs/QuizTab";
import { KokugoTab } from "./tabs/KokugoTab";
import { HomePage } from "./HomePage";
import { CapabilitiesProvider, useCapabilities } from "./capabilities";
import { ChineseVisibilityProvider, ChineseVisibilityToggle } from "./chineseVisibility";
import { api, type ProgressSummary } from "./api";

type Tab = "quiz" | "grammar" | "vocab" | "kanji" | "sentence" | "kokugo";
type ViewState = { view: "home" } | { view: "app"; initialMode: QuizInitialMode; initialTab?: Tab };

const tabs: { id: Tab; label: string }[] = [
  { id: "quiz", label: "練習題" },
  { id: "grammar", label: "文法" },
  { id: "vocab", label: "單字" },
  { id: "kanji", label: "漢字" },
  { id: "sentence", label: "例句" },
  { id: "kokugo", label: "国語" },
];

export function App() {
  return (
    <ChineseVisibilityProvider>
      <CapabilitiesProvider>
        <AppInner />
      </CapabilitiesProvider>
    </ChineseVisibilityProvider>
  );
}

function AppInner() {
  const { quiz, kokugo } = useCapabilities();
  const [viewState, setViewState] = useState<ViewState>({ view: "home" });

  if (viewState.view === "home") {
    return (
      <HomePage
        onStart={(initialMode, initialTab) =>
          setViewState({ view: "app", initialMode, initialTab })
        }
        quizCapable={quiz}
        kokugoCapable={kokugo}
      />
    );
  }

  return (
    <AppShell
      initialMode={viewState.initialMode}
      initialTab={viewState.initialTab}
    />
  );
}

function AppShell({
  initialMode,
  initialTab,
}: {
  initialMode: QuizInitialMode;
  initialTab?: Tab;
}) {
  const { loaded, quiz, sentence, kokugo } = useCapabilities();
  const visibleTabs = useMemo(
    () =>
      tabs.filter((tab) => {
        if (!loaded) return true;
        if (tab.id === "quiz") return quiz;
        if (tab.id === "sentence") return sentence;
        if (tab.id === "kokugo") return kokugo;
        return true;
      }),
    [loaded, quiz, sentence, kokugo]
  );
  const [active, setActive] = useState<Tab>(initialTab ?? "quiz");
  const [grammarSlug, setGrammarSlug] = useState<string | undefined>();
  const [kanjiCharacter, setKanjiCharacter] = useState<string | undefined>();

  useEffect(() => {
    if (!visibleTabs.some((tab) => tab.id === active)) {
      setActive(visibleTabs[0]?.id ?? "grammar");
    }
  }, [active, visibleTabs]);

  function navigateToGrammar(slug: string) {
    setGrammarSlug(slug);
    setActive("grammar");
  }

  function navigateToKanji(character: string) {
    setKanjiCharacter(character);
    setActive("kanji");
  }

  return (
    <div className="min-h-screen">
      <header className="bg-white border-b border-slate-200">
        <div className="max-w-4xl mx-auto px-4 py-4 flex items-center justify-between gap-4">
          <h1 className="text-xl font-semibold tracking-tight">
            {active === "kokugo" ? "国語教室" : "日本語学習"}{" "}
            <span className="text-slate-400 text-sm">— japanese-site</span>
          </h1>
          <div className="flex flex-wrap items-center justify-end gap-2">
            <ProgressBadge />
            <span className="text-xs text-slate-500">{quiz ? "M3 開發預覽" : "靜態瀏覽"}</span>
          </div>
        </div>
        <div className="max-w-4xl mx-auto px-4 flex items-end justify-between gap-4 -mb-px">
          <nav className="flex flex-wrap gap-1">
            {visibleTabs.map((t) => (
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
          <div className="pb-2">
            <ChineseVisibilityToggle />
          </div>
        </div>
      </header>
      <main className="max-w-4xl mx-auto px-4 py-8">
        {active === "quiz" && quiz && (
          <QuizTab initialMode={initialMode} onNavigateGrammar={navigateToGrammar} />
        )}
        {active === "grammar" && (
          <GrammarTab initialSlug={grammarSlug} onSlugConsumed={() => setGrammarSlug(undefined)} />
        )}
        {active === "vocab" && <VocabTab onNavigateKanji={navigateToKanji} />}
        {active === "kanji" && (
          <KanjiTab
            initialCharacter={kanjiCharacter}
            onInitialCharacterConsumed={() => setKanjiCharacter(undefined)}
          />
        )}
        {active === "sentence" && sentence && <SentenceTab />}
        {active === "kokugo" && kokugo && <KokugoTab />}
      </main>
      <footer className="text-center text-xs text-slate-400 py-6">
        資料來源 JMdict / KANJIDIC2 / Tatoeba (CC-BY-SA / CC-BY) ・ JLPT 標註 (MIT)
      </footer>
    </div>
  );
}

function ProgressBadge() {
  const { progress, progressRevisions } = useCapabilities();
  const [summary, setSummary] = useState<ProgressSummary | null>(null);

  useEffect(() => {
    if (!progress) {
      setSummary(null);
      return;
    }

    let cancelled = false;
    api
      .getProgress("grammar")
      .then((next) => {
        if (!cancelled) setSummary(next);
      })
      .catch((error) => {
        console.warn("failed to load progress summary", error);
        if (!cancelled) setSummary(null);
      });
    return () => {
      cancelled = true;
    };
  }, [progress, progressRevisions.grammar]);

  if (!progress || !summary) return null;

  return (
    <span className="rounded-md border border-emerald-200 bg-emerald-50 px-2 py-1 text-xs text-emerald-700">
      文法已讀 {summary.read}/{summary.total}
    </span>
  );
}
