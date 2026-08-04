import type { KokugoClassmate } from "../kokugoTypes";

/**
 * Curated classmate samples (JS-134). Parent decides when to mount
 * (after learner response). Pure presentational — no fetch.
 */
export function ClassmatePanel({
  classmates,
  title = "クラスメイトの答え",
}: {
  classmates: KokugoClassmate[];
  title?: string;
}) {
  if (classmates.length === 0) return null;

  return (
    <section
      className="rounded-lg border border-indigo-100 bg-indigo-50/50 p-3 space-y-2"
      aria-label={title}
    >
      <h4 className="text-sm font-medium text-indigo-900">{title}</h4>
      <p className="text-xs text-indigo-800/80">
        作答後に見る参考例です。そのまま写す必要はありません。
      </p>
      <ul className="space-y-2">
        {classmates.map((c) => (
          <li
            key={c.id}
            className="rounded-md border border-indigo-100 bg-white px-3 py-2 text-sm text-slate-800"
            data-classmate-id={c.id}
          >
            <div className="flex flex-wrap items-baseline gap-x-2 gap-y-0.5">
              <span className="font-medium text-indigo-900">{c.name_ja}</span>
              {c.focus_ja ? (
                <span className="text-xs text-slate-500">視点: {c.focus_ja}</span>
              ) : null}
            </div>
            <p className="mt-1 whitespace-pre-wrap text-slate-700">{c.text_ja}</p>
          </li>
        ))}
      </ul>
    </section>
  );
}
