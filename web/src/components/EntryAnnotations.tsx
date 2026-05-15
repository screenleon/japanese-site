import {
  ANNOTATION_KINDS,
  type AnnotationKind,
  type Annotations,
  type Block,
  type FuriganaAnnotation,
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

function hasFuriganaTitleTokens(value: FuriganaAnnotation | undefined) {
  if (!value) return false;
  return value.title_ja?.some(isRenderableToken) ?? false;
}

function hasAnnotationValue(kind: AnnotationKind, annotations?: Annotations) {
  const value = annotations?.[kind];
  if (kind === "furigana") {
    return hasFuriganaTitleTokens(value as FuriganaAnnotation | undefined);
  }
  return typeof value === "string" && value.trim().length > 0;
}

function FuriganaBlock({ value }: { value: FuriganaAnnotation }) {
  if (!hasFuriganaTitleTokens(value)) return null;

  return (
    <div className="space-y-1 text-sm leading-relaxed text-sky-950">
      <div>
        <Tokens tokens={value.title_ja ?? []} />
      </div>
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

function isNonBlankString(value: unknown): value is string {
  return typeof value === "string" && value.trim() !== "";
}

function isRenderableToken(value: unknown): value is Token {
  if (!value || typeof value !== "object" || !("t" in value)) return false;
  const token = value as Record<string, unknown>;
  if (token.t === "text") return isNonBlankString(token.v);
  if (token.t === "ruby") return isNonBlankString(token.k) && isNonBlankString(token.r);
  if (token.t === "term") {
    return (
      isNonBlankString(token.slug) &&
      isNonBlankString(token.label) &&
      (token.kind === "grammar" || token.kind === "vocab")
    );
  }
  return false;
}

function Tokens({ tokens }: { tokens: Token[] }) {
  const renderableTokens = tokens.filter(isRenderableToken);
  return (
    <>
      {renderableTokens.map((token, index) => (
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
