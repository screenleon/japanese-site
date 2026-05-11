import type { Annotations, GrammarPoint } from "../apiTypes";
import { BlockRenderer } from "./EntryAnnotations";

export function ClassifierContrasts({
  annotations,
  points = [],
  primaryPattern,
}: {
  annotations?: Annotations;
  points?: GrammarPoint[];
  primaryPattern?: string;
}) {
  const rules = annotations?.classifier?.rules ?? [];
  if (rules.length === 0) return null;

  return (
    <section className="rounded-md border border-slate-200 bg-white p-4">
      <h3 className="text-base font-medium text-slate-900">辨析</h3>
      <div className="mt-3 space-y-3">
        {rules.map((rule, index) => {
          const target = rule.with_slug
            ? points.find((point) => point.slug === rule.with_slug)
            : null;
          return (
            <article key={`${rule.with_pattern}-${index}`} className="rounded-md border border-slate-200 p-3">
              <div className="flex flex-wrap items-center gap-2 text-sm font-medium text-slate-900">
                <span>{primaryPattern || "本句型"}</span>
                <span aria-hidden="true">→</span>
                <span>{target?.title_ja ?? rule.with_pattern}</span>
                {target && (
                  <span className="rounded border border-slate-300 px-1.5 py-0.5 text-xs text-slate-600">
                    {target.jlpt_level}
                  </span>
                )}
              </div>
              <div className="mt-3">
                <BlockRenderer blocks={rule.rule_ja_blocks} />
              </div>
              {rule.rule_zh && (
                <p className="mt-2 text-sm leading-relaxed text-slate-600">{rule.rule_zh}</p>
              )}
              {rule.examples && rule.examples.length > 0 && (
                <table className="mt-3 w-full table-fixed border-collapse text-sm">
                  <tbody>
                    {rule.examples.map((example, exampleIndex) => (
                      <tr key={exampleIndex} className="border-t border-slate-100 align-top">
                        <td className="py-2 pr-3 text-emerald-800">{example.use_this}</td>
                        <td className="py-2 pl-3 text-rose-800">{example.use_alt}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              )}
            </article>
          );
        })}
      </div>
    </section>
  );
}
