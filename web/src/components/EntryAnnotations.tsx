import {
  ANNOTATION_KINDS,
  type AnnotationKind,
  type Annotations,
  type Block,
  type FuriganaAnnotation,
  type FuriganaPair,
  type Token,
} from "../apiTypes";

export const LABELS: Record<AnnotationKind, string> = {
  usage: "使い方",
  collocations: "コロケーション",
  particle_pairing: "助詞の組み合わせ",
  synonym_diff: "類義語の違い",
  mental_model: "考え方のヒント",
  nuance_note: "ニュアンス",
  furigana: "ふりがな",
  classifier: "辨析",
};

function hasFuriganaPairs(value: FuriganaAnnotation | undefined) {
  if (!value) return false;
  return [value.title_ja, value.vocabulary].some((pairs) =>
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
  const vocabularyPairs = filterRenderableFuriganaPairs(value.vocabulary);

  return (
    <div className="space-y-1 text-sm leading-relaxed text-sky-950">
      {titlePairs.length > 0 && (
        <div>
          <RubyList pairs={titlePairs} />
        </div>
      )}
      {vocabularyPairs.length > 0 && (
        <div>
          <RubyList pairs={vocabularyPairs} />
        </div>
      )}
    </div>
  );
}

function TokenRenderer({ token }: { token: Token }) {
  if (token.t === "ruby") {
    return (
      <ruby>
        {token.k}
        <rt>{token.r}</rt>
      </ruby>
    );
  }
  if (token.t === "term") {
    const href = token.kind === "grammar" ? `#grammar/${token.slug}` : `#vocab/${token.slug}`;
    return (
      <a className="font-medium text-blue-700 underline-offset-2 hover:underline" href={href}>
        {token.label}
      </a>
    );
  }
  return <>{token.v}</>;
}

function Tokens({ tokens }: { tokens: Token[] }) {
  return (
    <>
      {tokens.map((token, index) => (
        <TokenRenderer key={`${token.t}-${index}`} token={token} />
      ))}
    </>
  );
}

export function BlockRenderer({ blocks }: { blocks: Block[] }) {
  if (!Array.isArray(blocks) || blocks.length === 0) return null;

  return (
    <div className="space-y-3 text-sm leading-relaxed text-slate-900">
      {blocks.map((block, index) => {
        if (block.kind === "list") {
          return (
            <ul key={index} className="list-disc space-y-1 pl-5">
              {block.items.map((item, itemIndex) => (
                <li key={itemIndex}>
                  <Tokens tokens={item.tokens} />
                </li>
              ))}
            </ul>
          );
        }
        if (block.kind === "callout") {
          const toneClass =
            block.tone === "warn"
              ? "border-amber-300 bg-amber-50 text-amber-950"
              : block.tone === "tip"
                ? "border-emerald-300 bg-emerald-50 text-emerald-950"
                : "border-sky-300 bg-sky-50 text-sky-950";
          return (
            <aside key={index} className={`rounded-md border px-3 py-2 ${toneClass}`}>
              <Tokens tokens={block.tokens} />
            </aside>
          );
        }
        return (
          <p key={index} className="whitespace-pre-wrap">
            <Tokens tokens={block.tokens} />
          </p>
        );
      })}
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
    (kind) => kind !== "classifier" && allowedKinds.has(kind) && hasAnnotationValue(kind, annotations)
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
