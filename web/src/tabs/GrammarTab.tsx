import { useEffect, useState } from "react";
import { api, type GrammarPoint } from "../api";

export function GrammarTab() {
  const [points, setPoints] = useState<GrammarPoint[]>([]);
  const [active, setActive] = useState<GrammarPoint | null>(null);

  useEffect(() => {
    api.listGrammar().then((r) => {
      setPoints(r.points || []);
      if (r.points && r.points[0]) setActive(r.points[0]);
    });
  }, []);

  return (
    <section className="grid grid-cols-[200px_1fr] gap-6">
      <ul className="space-y-1">
        {points.map((p) => (
          <li key={p.slug}>
            <button
              onClick={() => setActive(p)}
              className={
                "w-full text-left px-3 py-2 rounded-md text-sm " +
                (active?.slug === p.slug
                  ? "bg-blue-50 text-blue-700 font-medium"
                  : "text-slate-700 hover:bg-slate-100")
              }
            >
              <div>{p.title_ja}</div>
              <div className="text-xs text-slate-400">{p.jlpt_level} · {p.title_zh}</div>
            </button>
          </li>
        ))}
      </ul>
      {active && (
        <article className="bg-white border border-slate-200 rounded-md p-6">
          <header className="mb-4">
            <h2 className="text-2xl font-medium">{active.title_ja}</h2>
            <p className="text-sm text-slate-500">{active.jlpt_level} · {active.title_zh}</p>
          </header>
          <pre className="text-sm whitespace-pre-wrap font-sans leading-relaxed">{active.explanation_zh}</pre>
        </article>
      )}
    </section>
  );
}
