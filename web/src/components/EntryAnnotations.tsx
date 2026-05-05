import { ANNOTATION_KINDS, type AnnotationKind, type Annotations } from "../apiTypes";

export const LABELS: Record<AnnotationKind, string> = {
  usage: "使い方",
  collocations: "コロケーション",
  particle_pairing: "助詞の組み合わせ",
  synonym_diff: "類義語の違い",
  mental_model: "考え方のヒント",
  nuance_note: "ニュアンス",
};

export function EntryAnnotations({
  annotations,
  kinds,
}: {
  annotations?: Annotations;
  kinds?: AnnotationKind[];
}) {
  const allowedKinds = new Set(kinds ?? ANNOTATION_KINDS);
  const visibleKinds = ANNOTATION_KINDS.filter(
    (kind) => allowedKinds.has(kind) && annotations?.[kind]?.trim()
  );

  if (visibleKinds.length === 0) return null;

  return (
    <aside className="rounded-md border border-sky-200 bg-sky-50 px-4 py-3">
      <dl className="grid gap-3 sm:grid-cols-2">
        {visibleKinds.map((kind) => (
          <div key={kind}>
            <dt
              className="text-sm font-semibold text-sky-900"
              role="heading"
              aria-level={3}
            >
              {LABELS[kind]}
            </dt>
            <dd className="mt-1 whitespace-pre-wrap text-sm leading-relaxed text-sky-950">
              {annotations?.[kind]}
            </dd>
          </div>
        ))}
      </dl>
    </aside>
  );
}
