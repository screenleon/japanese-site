import { useEffect, useMemo, useState } from "react";
import { api, type GrammarExample, type GrammarPoint } from "../api";
import { useReadTracking } from "../hooks/useReadTracking";

const levels = ["N5", "N4", "N3", "N2", "N1"];

export function GrammarTab({
  initialSlug,
  onSlugConsumed,
}: {
  initialSlug?: string;
  onSlugConsumed?: () => void;
} = {}) {
  const [points, setPoints] = useState<GrammarPoint[]>([]);
  const [selectedLevel, setSelectedLevel] = useState("N5");
  const [activeSlug, setActiveSlug] = useState("");
  const [examples, setExamples] = useState<GrammarExample[]>([]);
  const [showTranslation, setShowTranslation] = useState(false);
  const [loading, setLoading] = useState(true);
  const [loadingRandom, setLoadingRandom] = useState(false);
  const [err, setErr] = useState("");

  useEffect(() => {
    let cancelled = false;
    setLoading(true);
    setErr("");
    api
      .listGrammar()
      .then((r) => {
        if (cancelled) return;
        const next = r.points || [];
        setPoints(next);
        const firstWithContent = levels.find((level) => next.some((p) => p.jlpt_level === level));
        if (firstWithContent) setSelectedLevel(firstWithContent);
      })
      .catch((e) => {
        if (!cancelled) setErr(String(e));
      })
      .finally(() => {
        if (!cancelled) setLoading(false);
      });
    return () => {
      cancelled = true;
    };
  }, []);

  // When navigated here from ResultBox with a specific slug, select it.
  useEffect(() => {
    if (!initialSlug || points.length === 0) return;
    const target = points.find((p) => p.slug === initialSlug);
    if (target) {
      setSelectedLevel(target.jlpt_level);
      setActiveSlug(target.slug);
    }
    onSlugConsumed?.();
  }, [initialSlug, points]); // eslint-disable-line react-hooks/exhaustive-deps

  const grouped = useMemo(() => {
    return levels.map((level) => ({
      level,
      points: points.filter((p) => p.jlpt_level === level),
    }));
  }, [points]);

  const levelPoints = grouped.find((g) => g.level === selectedLevel)?.points || [];
  const active = levelPoints.find((p) => p.slug === activeSlug) || levelPoints[0] || null;
  const primaryExplanation = active?.explanation_ja || active?.explanation_zh || "";
  const hasJapaneseExplanation = Boolean(active?.explanation_ja?.trim());
  useReadTracking(active?.slug ? { type: "grammar", slug: active.slug } : null);

  useEffect(() => {
    setActiveSlug("");
    setShowTranslation(false);
  }, [selectedLevel]);

  useEffect(() => {
    setShowTranslation(false);
  }, [active?.slug]);

  useEffect(() => {
    let cancelled = false;
    setExamples([]);
    if (!active?.slug || !api.getGrammarExamples) {
      return () => {
        cancelled = true;
      };
    }
    api
      .getGrammarExamples(active.slug)
      .then((r) => {
        if (!cancelled) setExamples((r.examples || []).slice(0, 5));
      })
      .catch((e) => {
        if (!cancelled) setErr(String(e));
      });
    return () => {
      cancelled = true;
    };
  }, [active?.slug]);

  async function drawRandomGrammar() {
    setLoadingRandom(true);
    setErr("");
    try {
      const gp = await api.randomGrammar(selectedLevel);
      setActiveSlug(gp.slug);
    } catch (e) {
      setErr(String(e));
    } finally {
      setLoadingRandom(false);
    }
  }

  if (loading) {
    return <p className="text-sm text-slate-500">文法を読み込み中...</p>;
  }

  if (err) {
    return <p className="text-sm text-red-600">{err}</p>;
  }

  return (
    <section className="space-y-6">
      <div className="grid gap-3 sm:grid-cols-5">
        {grouped.map(({ level, points }) => (
          <button
            key={level}
            type="button"
            onClick={() => setSelectedLevel(level)}
            className={
              "rounded-md border px-4 py-3 text-left transition-colors " +
              (selectedLevel === level
                ? "border-blue-600 bg-blue-50 text-blue-700"
                : "border-slate-200 bg-white text-slate-700 hover:border-slate-300")
            }
          >
            <div className="text-lg font-semibold">{level}</div>
            <div className="text-xs text-slate-500">{points.length} 文法</div>
          </button>
        ))}
      </div>

      {levelPoints.length === 0 ? (
        <div className="rounded-md border border-slate-200 bg-white p-6 text-sm text-slate-500">
          この等級の文法はまだ準備中です。
        </div>
      ) : (
        <div className="space-y-5">
          <section className="rounded-md border border-slate-200 bg-white p-5">
            <div className="flex flex-wrap items-center justify-between gap-3">
              <div>
                <h2 className="text-lg font-medium">{selectedLevel} 隨機文法</h2>
                <p className="text-sm text-slate-500">抽一個文法點，先讀日文說明，再需要時打開中文。</p>
              </div>
              <button
                type="button"
                onClick={() => void drawRandomGrammar()}
                disabled={loadingRandom}
                className="rounded-md bg-blue-600 px-4 py-2 text-sm text-white disabled:opacity-50"
              >
                {loadingRandom ? "抽取中" : "抽下一個"}
              </button>
            </div>
            {active && (
              <div className="mt-4 border-t border-slate-100 pt-4">
                <div className="text-xl font-semibold">{active.title_ja}</div>
                <div className="mt-1 text-sm text-slate-500">{active.title_zh}</div>
              </div>
            )}
          </section>

          <div className="grid gap-6 lg:grid-cols-[240px_1fr]">
            <ul className="space-y-1">
              {levelPoints.map((p) => (
                <li key={p.slug}>
                  <button
                    type="button"
                    onClick={() => setActiveSlug(p.slug)}
                    className={
                      "w-full rounded-md px-3 py-2 text-left text-sm " +
                      (active?.slug === p.slug
                        ? "bg-blue-50 font-medium text-blue-700"
                        : "text-slate-700 hover:bg-slate-100")
                    }
                  >
                    <div>{p.title_ja}</div>
                    <div className="text-xs text-slate-400">{p.title_zh}</div>
                  </button>
                </li>
              ))}
            </ul>

            {active && (
              <article className="rounded-md border border-slate-200 bg-white p-6">
                <header className="mb-4">
                  <h2 className="text-2xl font-medium">{active.title_ja}</h2>
                  <p className="text-sm text-slate-500">
                    {active.jlpt_level} · {active.title_zh}
                  </p>
                </header>
                <pre className="whitespace-pre-wrap font-sans text-sm leading-relaxed">
                  {primaryExplanation}
                </pre>
                {examples.length > 0 && (
                  <section className="mt-6 border-t border-slate-200 pt-5">
                    <h3 className="text-base font-medium">例文</h3>
                    <ul className="mt-3 space-y-3">
                      {examples.map((example) => (
                        <li key={example.id}>
                          <p className="text-sm leading-relaxed text-slate-900">
                            {example.text_ja}
                          </p>
                          <p className="mt-1 text-sm text-slate-500">{example.text_zh}</p>
                        </li>
                      ))}
                    </ul>
                  </section>
                )}
                {hasJapaneseExplanation && (
                  <div className="mt-5 border-t border-slate-200 pt-4">
                    <button
                      type="button"
                      onClick={() => setShowTranslation((current) => !current)}
                      className="rounded-md border border-slate-300 px-3 py-2 text-sm text-slate-700 hover:bg-slate-50"
                    >
                      {showTranslation ? "隱藏中文說明" : "顯示中文說明"}
                    </button>
                    {showTranslation && (
                      <pre className="mt-4 whitespace-pre-wrap font-sans text-sm leading-relaxed text-slate-700">
                        {active.explanation_zh}
                      </pre>
                    )}
                  </div>
                )}
              </article>
            )}
          </div>
        </div>
      )}
    </section>
  );
}
