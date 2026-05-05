import { useEffect, useMemo, useRef, useState } from "react";
import {
  api,
  isApiError,
  type Annotations,
  type GrammarExample,
  type GrammarPoint,
} from "../api";
import { EntryAnnotations } from "../components/EntryAnnotations";
import { useReadTracking } from "../hooks/useReadTracking";

const levels = ["N5", "N4", "N3", "N2", "N1"];

function mergeFromGrammarPoint(point: GrammarPoint): Annotations {
  return {
    ...point.annotations,
    mental_model: point.annotations?.mental_model ?? point.mental_model,
    nuance_note: point.annotations?.nuance_note ?? point.nuance_note,
  };
}

function hasAnnotations(annotations: Annotations) {
  return Object.values(annotations).some((value) => value?.trim());
}

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
  const [examplesFailed, setExamplesFailed] = useState(false);
  const [showTranslation, setShowTranslation] = useState(false);
  const [loading, setLoading] = useState(true);
  const initialSlugApplied = useRef(false);
  const skipExamplesForInitialSlug = useRef(false);
  const preserveActiveOnLevelChange = useRef(false);
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
      initialSlugApplied.current = true;
      skipExamplesForInitialSlug.current = true;
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
  const relatedPoints = useMemo(() => {
    return (active?.related_slugs || []).map((slug) => ({
      slug,
      point: points.find((p) => p.slug === slug) || null,
    }));
  }, [active?.related_slugs, points]);
  useReadTracking(active?.slug ? { type: "grammar", slug: active.slug } : null);

  useEffect(() => {
    if (initialSlugApplied.current) {
      initialSlugApplied.current = false;
      return;
    }
    if (preserveActiveOnLevelChange.current) {
      preserveActiveOnLevelChange.current = false;
      return;
    }
    setActiveSlug("");
    setShowTranslation(false);
  }, [selectedLevel]);

  useEffect(() => {
    setShowTranslation(false);
  }, [active?.slug]);

  useEffect(() => {
    let cancelled = false;
    setExamples([]);
    setExamplesFailed(false);
    if (skipExamplesForInitialSlug.current) {
      skipExamplesForInitialSlug.current = false;
      return () => {
        cancelled = true;
      };
    }
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
      .catch((err) => {
        console.warn("grammar examples fetch failed", err);
        if (!cancelled) {
          setExamples([]);
          setExamplesFailed(!isApiError(err, "not_found"));
        }
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
                  <div className="flex flex-wrap items-start justify-between gap-3">
                    <h2 className="text-2xl font-medium">{active.title_ja}</h2>
                    <button
                      type="button"
                      onClick={() => void drawRandomGrammar()}
                      disabled={loadingRandom}
                      className="rounded-md bg-blue-600 px-4 py-2 text-sm text-white disabled:opacity-50"
                    >
                      {loadingRandom ? "抽取中" : "抽下一個"}
                    </button>
                  </div>
                  <p className="text-sm text-slate-500">
                    {active.jlpt_level} · {active.title_zh}
                  </p>
                </header>
                {hasAnnotations(mergeFromGrammarPoint(active)) && (
                  <div className="mb-4">
                    <EntryAnnotations
                      annotations={mergeFromGrammarPoint(active)}
                      kinds={["nuance_note", "mental_model"]}
                    />
                  </div>
                )}
                <pre className="whitespace-pre-wrap font-sans text-sm leading-relaxed">
                  {primaryExplanation}
                </pre>
                {(examples.length > 0 || examplesFailed) && (
                  <section className="mt-6 border-t border-slate-200 pt-5">
                    <h3 className="text-base font-medium">例文</h3>
                    {examples.length === 0 ? (
                      <p className="mt-2 text-sm text-slate-500">例文を表示できません。</p>
                    ) : (
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
                    )}
                  </section>
                )}
                {relatedPoints.length > 0 && (
                  <section className="mt-6 border-t border-slate-200 pt-5">
                    <h3 className="text-base font-medium">相關用法</h3>
                    <div className="mt-3 flex flex-wrap gap-2">
                      {relatedPoints.map(({ slug, point }) =>
                        point ? (
                          <button
                            key={slug}
                            type="button"
                            onClick={() => {
                              preserveActiveOnLevelChange.current = true;
                              setSelectedLevel(point.jlpt_level);
                              setActiveSlug(point.slug);
                            }}
                            className="rounded-md border border-slate-300 px-3 py-2 text-sm text-slate-700 hover:bg-slate-50"
                          >
                            {point.title_ja} ({point.jlpt_level})
                          </button>
                        ) : (
                          <span key={slug} className="text-sm text-slate-500">
                            {slug}
                          </span>
                        )
                      )}
                    </div>
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
