/**
 * Side-by-side draft vs revision compare (JS-134).
 * Pure presentational; empty sides still render labels for orientation.
 */
export function RevisionCompare({
  draftBody,
  revisionBody,
  title = "下書きと改稿の対比",
}: {
  draftBody: string;
  revisionBody: string;
  title?: string;
}) {
  const draft = draftBody.trim();
  const revision = revisionBody.trim();
  if (!draft && !revision) return null;

  const draftChars = [...draft].length;
  const revChars = [...revision].length;
  const delta = revChars - draftChars;

  return (
    <section
      className="rounded-lg border border-slate-200 bg-white p-3 space-y-2"
      aria-label={title}
    >
      <div className="flex flex-wrap items-baseline justify-between gap-2">
        <h4 className="text-sm font-medium text-slate-800">{title}</h4>
        <p className="text-xs text-slate-500">
          下書き {draftChars} 字 → 改稿 {revChars} 字
          {draft && revision ? (
            <span className="ml-1">
              （{delta === 0 ? "字数同じ" : delta > 0 ? `+${delta} 字` : `${delta} 字`}）
            </span>
          ) : null}
        </p>
      </div>
      <div className="grid gap-2 sm:grid-cols-2">
        <article className="rounded-md border border-slate-100 bg-slate-50 p-2">
          <h5 className="text-xs font-medium text-slate-600 mb-1">下書き</h5>
          <pre className="whitespace-pre-wrap text-xs text-slate-800 font-sans m-0">
            {draft || "（まだありません）"}
          </pre>
        </article>
        <article className="rounded-md border border-emerald-100 bg-emerald-50/40 p-2">
          <h5 className="text-xs font-medium text-emerald-800 mb-1">改稿後</h5>
          <pre className="whitespace-pre-wrap text-xs text-slate-800 font-sans m-0">
            {revision || "（まだありません）"}
          </pre>
        </article>
      </div>
    </section>
  );
}
