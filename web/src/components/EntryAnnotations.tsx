import {
  ANNOTATION_KINDS,
  type AnnotationKind,
  type Annotations,
  type FuriganaAnnotation,
  type FuriganaPair,
} from "../apiTypes";

export const LABELS: Record<AnnotationKind, string> = {
  usage: "使い方",
  collocations: "コロケーション",
  particle_pairing: "助詞の組み合わせ",
  synonym_diff: "類義語の違い",
  mental_model: "考え方のヒント",
  nuance_note: "ニュアンス",
  furigana: "ふりがな",
};

function hasFuriganaPairs(value: FuriganaAnnotation | undefined) {
  if (!value) return false;
  return [value.title_ja, value.key_terms].some((pairs) =>
    filterRenderableFuriganaPairs(pairs).length > 0
  );
}

function isRenderableFuriganaPair(pair: unknown): pair is FuriganaPair {
  if (!pair || typeof pair !== "object") return false;
  const candidate = pair as Partial<Record<keyof FuriganaPair, unknown>>;
  return (
    typeof candidate.kanji === "string" &&
    typeof candidate.reading === "string" &&
    candidate.kanji.trim() !== "" &&
    candidate.reading.trim() !== ""
  );
}

function filterRenderableFuriganaPairs(pairs: unknown) {
  if (!Array.isArray(pairs)) return [];
  return pairs.filter(isRenderableFuriganaPair);
}

function hasAnnotationValue(kind: AnnotationKind, annotations?: Annotations) {
  const value = annotations?.[kind];
  if (kind === "furigana") {
    return hasFuriganaPairs(value as FuriganaAnnotation | undefined);
  }
  return typeof value === "string" && value.trim().length > 0;
}

function RubyList({ pairs }: { pairs: FuriganaPair[] }) {
  return (
    <>
      {pairs.map((pair, index) => (
        <span key={`${pair.kanji}-${pair.reading}-${index}`}>
          {index > 0 && "、"}
          <ruby>
            {pair.kanji}
            <rt>{pair.reading}</rt>
          </ruby>
        </span>
      ))}
    </>
  );
}

function FuriganaBlock({ value }: { value: FuriganaAnnotation }) {
  const titlePairs = filterRenderableFuriganaPairs(value.title_ja);
  const keyTermPairs = filterRenderableFuriganaPairs(value.key_terms);

  return (
    <div className="space-y-1 text-sm leading-relaxed text-sky-950">
      {titlePairs.length > 0 && (
        <div>
          <RubyList pairs={titlePairs} />
        </div>
      )}
      {keyTermPairs.length > 0 && (
        <div>
          <RubyList pairs={keyTermPairs} />
        </div>
      )}
    </div>
  );
}

export function EntryAnnotations({
  annotations,
  kinds,
}: {
  annotations?: Annotations;
  kinds?: AnnotationKind[];
}) {
  const allowedKinds = new Set(kinds ?? ANNOTATION_KINDS);
  const visibleKinds = ANNOTATION_KINDS.filter(
    (kind) => allowedKinds.has(kind) && hasAnnotationValue(kind, annotations)
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
            {kind === "furigana" ? (
              <dd className="mt-1">
                <FuriganaBlock value={annotations?.furigana as FuriganaAnnotation} />
              </dd>
            ) : (
              <dd className="mt-1 whitespace-pre-wrap text-sm leading-relaxed text-sky-950">
                {annotations?.[kind] as string}
              </dd>
            )}
          </div>
        ))}
      </dl>
    </aside>
  );
}
