// JS-133 — in-passage reader for 国語 evidence highlight + paragraph roles.
// Keeps grading contracts unchanged: quotes[] / roles[] surface strings.
// UI-001: prop-driven / no persistent selection state here (tabs own Maps).
// ADR-0005: unit.text is Block[]; non-paragraph kinds stay visible via BlockRenderer.

import { useMemo } from "react";
import type { Block, Token } from "../apiTypes";
import { BlockRenderer } from "./EntryAnnotations";

/** One selectable sentence within a paragraph (plain surface text). */
export interface PassageSentence {
  /** Stable key: `${paraIndex}:${sentIndex}` or `gold:…` for fallback chips. */
  key: string;
  /** Paragraph index among paragraph blocks only, or -1 for orphan gold. */
  paraIndex: number;
  sentIndex: number;
  text: string;
}

export interface PassageParagraph {
  paraIndex: number;
  block: Extract<Block, { kind: "paragraph" }>;
  plain: string;
  sentences: PassageSentence[];
}

/** Ordered stream item preserving full Block[] for rendering. */
export type PassageStreamItem =
  | {
      kind: "paragraph";
      paraIndex: number;
      block: Extract<Block, { kind: "paragraph" }>;
      plain: string;
      sentences: PassageSentence[];
    }
  | {
      kind: "other";
      blockIndex: number;
      block: Exclude<Block, { kind: "paragraph" }>;
    };

/** Flatten Block tokens to surface Japanese (kanji base for ruby). */
export function plainFromTokens(tokens: Token[]): string {
  return tokens
    .map((t) => {
      if (t.t === "text") return t.v;
      if (t.t === "ruby") return t.k;
      return t.label;
    })
    .join("");
}

/**
 * Split Japanese prose into sentences. Keeps the terminator with the clause.
 * Empty / whitespace-only chunks are dropped.
 */
export function splitJapaneseSentences(text: string): string[] {
  if (!text) return [];
  return text
    .split(/(?<=[。．？！\n])/)
    .map((s) => s.trim())
    .filter((s) => s.length > 0);
}

/** Count paragraph blocks in unit.text order (shared by role step init). */
export function countParagraphs(blocks: Block[]): number {
  return blocks.reduce((n, b) => n + (b.kind === "paragraph" ? 1 : 0), 0);
}

/** Build ordered stream: all Block kinds in place; paragraph indices stable. */
export function buildPassageStream(blocks: Block[]): PassageStreamItem[] {
  let paraIndex = 0;
  return blocks.map((block, blockIndex) => {
    if (block.kind === "paragraph") {
      const plain = plainFromTokens(block.tokens);
      const sentences = splitJapaneseSentences(plain).map((text, sentIndex) => ({
        key: `${paraIndex}:${sentIndex}`,
        paraIndex,
        sentIndex,
        text,
      }));
      const item: PassageStreamItem = {
        kind: "paragraph",
        paraIndex,
        block,
        plain,
        sentences,
      };
      paraIndex += 1;
      return item;
    }
    return { kind: "other", blockIndex, block };
  });
}

/** Paragraph-only model (selection / roles). Derived from the full stream. */
export function buildPassageModel(blocks: Block[]): PassageParagraph[] {
  return buildPassageStream(blocks)
    .filter((item): item is Extract<PassageStreamItem, { kind: "paragraph" }> => item.kind === "paragraph")
    .map(({ paraIndex, block, plain, sentences }) => ({
      paraIndex,
      block,
      plain,
      sentences,
    }));
}

/**
 * Gold that is already a substring of some sentence needs no extra UI — the
 * learner taps that sentence. Gold that spans multiple sentences is also OK
 * without extras: multi-select + server compactSpace join still matches.
 *
 * Only when gold is **not** present in any paragraph plain text do we emit a
 * defensive fallback chip (lint should prevent this in real corpus). We never
 * surface an in-text gold string as a labeled “answer chip” (would spoil).
 */
export function ensureGoldSelectable(
  model: PassageParagraph[],
  goldQuotes: string[]
): PassageSentence[] {
  const all = model.flatMap((p) => p.sentences);
  const extras: PassageSentence[] = [];
  for (const gold of goldQuotes) {
    const g = gold.trim();
    if (!g) continue;
    if (all.some((s) => s.text.includes(g)) || extras.some((s) => s.text.includes(g))) {
      continue;
    }
    // Spanning or multi-sentence gold already in plain: no spoiler chip.
    if (model.some((p) => p.plain.includes(g))) {
      continue;
    }
    extras.push({
      key: `gold:orphan:${extras.length}`,
      paraIndex: -1,
      sentIndex: -1,
      text: g,
    });
  }
  return extras;
}

const ROLE_STYLE: Record<string, string> = {
  問題: "border-rose-300 bg-rose-50/80",
  原因: "border-amber-300 bg-amber-50/80",
  提案: "border-sky-300 bg-sky-50/80",
  結論: "border-emerald-300 bg-emerald-50/80",
};

function roleBorderClass(role: string | undefined): string {
  if (!role) return "border-slate-200 bg-white";
  return ROLE_STYLE[role] ?? "border-indigo-200 bg-indigo-50/50";
}

export type KokugoPassageMode = "readonly" | "sentence-select" | "paragraph-role";

export interface KokugoPassageProps {
  blocks: Block[];
  mode?: KokugoPassageMode;
  /** Selected sentence keys (sentence-select mode). */
  selectedKeys?: ReadonlySet<string>;
  onToggleSentence?: (sentence: PassageSentence) => void;
  /** Per-paragraph role (paragraph-role mode), same order as paragraph blocks. */
  roles?: string[];
  roleOptions?: string[];
  onRoleChange?: (paraIndex: number, role: string) => void;
  /** Optional gold quotes only for ensureGoldSelectable extras in sentence-select. */
  goldQuotes?: string[];
  className?: string;
  showParagraphIndex?: boolean;
}

function SentenceChip({
  sentence,
  selected,
  onToggle,
  variant,
}: {
  sentence: PassageSentence;
  selected: boolean;
  onToggle?: (sentence: PassageSentence) => void;
  variant: "inline" | "fallback";
}) {
  if (variant === "fallback") {
    return (
      <button
        type="button"
        data-sentence-key={sentence.key}
        aria-pressed={selected}
        onClick={() => onToggle?.(sentence)}
        className={
          "mt-1 block w-full rounded border border-dashed px-2 py-1 text-left text-sm " +
          (selected
            ? "border-amber-400 bg-amber-100"
            : "border-slate-300 bg-white hover:bg-slate-50")
        }
      >
        <span className="text-[10px] text-slate-500">補足引用 · </span>
        {sentence.text}
      </button>
    );
  }
  return (
    <button
      type="button"
      data-sentence-key={sentence.key}
      aria-pressed={selected}
      onClick={() => onToggle?.(sentence)}
      className={
        "mr-0.5 mb-1 inline rounded px-0.5 text-left transition-colors " +
        (selected
          ? "bg-amber-200/90 text-amber-950 ring-1 ring-amber-400"
          : "hover:bg-sky-100/80 focus-visible:bg-sky-100 focus-visible:outline focus-visible:outline-2 focus-visible:outline-sky-400")
      }
    >
      {sentence.text}
    </button>
  );
}

export function KokugoPassage({
  blocks,
  mode = "readonly",
  selectedKeys,
  onToggleSentence,
  roles,
  roleOptions,
  onRoleChange,
  goldQuotes,
  className = "",
  showParagraphIndex = true,
}: KokugoPassageProps) {
  const stream = useMemo(() => buildPassageStream(blocks), [blocks]);
  const model = useMemo(
    () =>
      stream
        .filter(
          (item): item is Extract<PassageStreamItem, { kind: "paragraph" }> =>
            item.kind === "paragraph"
        )
        .map(({ paraIndex, block, plain, sentences }) => ({
          paraIndex,
          block,
          plain,
          sentences,
        })),
    [stream]
  );
  const goldExtras = useMemo(() => {
    if (mode !== "sentence-select" || !goldQuotes?.length) return [] as PassageSentence[];
    return ensureGoldSelectable(model, goldQuotes);
  }, [mode, goldQuotes, model]);

  const orphanGold = useMemo(
    () => goldExtras.filter((g) => g.paraIndex < 0),
    [goldExtras]
  );

  if (blocks.length === 0) {
    return (
      <article
        className={`rounded-lg border border-slate-200 bg-white p-4 text-sm text-slate-500 ${className}`}
        aria-label="本文"
      >
        本文がありません。
      </article>
    );
  }

  return (
    <article
      className={`space-y-4 rounded-lg border border-slate-200 bg-white p-4 leading-8 ${className}`}
      aria-label="本文"
    >
      {stream.map((item) => {
        if (item.kind === "other") {
          return (
            <div
              key={`other-${item.blockIndex}`}
              className="rounded-md border border-slate-100 bg-white px-1 py-1"
              data-block-kind={item.block.kind}
              data-block-index={item.blockIndex}
            >
              <BlockRenderer blocks={[item.block]} />
            </div>
          );
        }

        const { paraIndex, block, sentences } = item;
        const role = roles?.[paraIndex];
        const border =
          mode === "paragraph-role" ? roleBorderClass(role) : "border-slate-100 bg-slate-50/40";
        const paraGold = goldExtras.filter((g) => g.paraIndex === paraIndex);

        return (
          <div
            key={`para-${paraIndex}`}
            className={`rounded-md border-l-4 pl-3 pr-2 py-2 ${border}`}
            data-para-index={paraIndex}
          >
            <div className="mb-1 flex flex-wrap items-center gap-2">
              {showParagraphIndex && (
                <span className="text-[11px] font-medium uppercase tracking-wide text-slate-500">
                  段落 {paraIndex + 1}
                </span>
              )}
              {mode === "paragraph-role" && roleOptions && roleOptions.length > 0 && (
                <label className="flex items-center gap-1.5 text-xs text-slate-700">
                  <span className="sr-only">段落 {paraIndex + 1} の役割</span>
                  <select
                    className="rounded border border-slate-300 bg-white px-2 py-0.5 text-xs font-medium"
                    value={role ?? roleOptions[0] ?? ""}
                    aria-label={`段落 ${paraIndex + 1} の役割`}
                    onChange={(e) => onRoleChange?.(paraIndex, e.target.value)}
                  >
                    {roleOptions.map((r) => (
                      <option key={r} value={r}>
                        {r}
                      </option>
                    ))}
                  </select>
                </label>
              )}
              {mode === "paragraph-role" && role && (
                <span className="rounded-full bg-white/80 px-2 py-0.5 text-[11px] font-medium text-slate-700 ring-1 ring-slate-200">
                  {role}
                </span>
              )}
            </div>

            {mode === "sentence-select" ? (
              <div className="text-sm leading-8 text-slate-900">
                <p>
                  {sentences.map((sent) => (
                    <SentenceChip
                      key={sent.key}
                      sentence={sent}
                      selected={selectedKeys?.has(sent.key) ?? false}
                      onToggle={onToggleSentence}
                      variant="inline"
                    />
                  ))}
                </p>
                {paraGold.map((sent) => (
                  <SentenceChip
                    key={sent.key}
                    sentence={sent}
                    selected={selectedKeys?.has(sent.key) ?? false}
                    onToggle={onToggleSentence}
                    variant="fallback"
                  />
                ))}
              </div>
            ) : (
              <div className="text-sm leading-8 text-slate-900">
                <BlockRenderer blocks={[block]} />
              </div>
            )}
          </div>
        );
      })}

      {mode === "sentence-select" &&
        orphanGold.map((sent) => (
          <SentenceChip
            key={sent.key}
            sentence={sent}
            selected={selectedKeys?.has(sent.key) ?? false}
            onToggle={onToggleSentence}
            variant="fallback"
          />
        ))}
    </article>
  );
}
