import { fireEvent, render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import type { Block } from "../apiTypes";
import {
  assessGoldCoverage,
  buildPassageModel,
  countParagraphs,
  ensureGoldSelectable,
  KokugoPassage,
  plainFromTokens,
  splitJapaneseSentences,
} from "./KokugoPassage";

describe("plainFromTokens / splitJapaneseSentences", () => {
  it("flattens text and ruby base", () => {
    /**
     * Behavior: ruby tokens contribute kanji base only.
     * 1. Build mixed text/ruby tokens.
     * 2. Flatten via plainFromTokens.
     * 3. Expect concatenated surface without furigana readings.
     */
    expect(
      plainFromTokens([
        { t: "text", v: "図書" },
        { t: "ruby", k: "室", r: "しつ" },
        { t: "text", v: "です。" },
      ])
    ).toBe("図書室です。");
  });

  it("splits on Japanese terminators and keeps them", () => {
    /**
     * Behavior: sentence split retains terminators.
     * 1. Provide multi-clause Japanese string.
     * 2. Split with splitJapaneseSentences.
     * 3. Expect three clauses each keeping 。？！
     */
    expect(splitJapaneseSentences("Aです。Bます？C！")).toEqual(["Aです。", "Bます？", "C！"]);
  });

  it("drops empty and whitespace-only sentence chunks", () => {
    /**
     * Behavior: empty/whitespace-only chunks are never selectable.
     * 1. Call split on empty string, whitespace-only, and newline-only separators.
     * 2. Assert empty arrays or only real sentences remain.
     * 3. Confirm no zero-length entries.
     */
    expect(splitJapaneseSentences("")).toEqual([]);
    expect(splitJapaneseSentences("   \n\t  ")).toEqual([]);
    expect(splitJapaneseSentences("文です。\n\n次です。")).toEqual(["文です。", "次です。"]);
    expect(splitJapaneseSentences("文です。   \n   ")).toEqual(["文です。"]);
  });
});

describe("countParagraphs / buildPassageModel", () => {
  const blocks: Block[] = [
    {
      kind: "paragraph",
      tokens: [{ t: "text", v: "問題です。次は原因です。" }],
    },
    {
      kind: "paragraph",
      tokens: [{ t: "text", v: "まず探しやすさを改善する必要があります。結論です。" }],
    },
  ];

  it("countParagraphs ignores non-paragraph blocks and handles empty", () => {
    /**
     * Behavior: roles array length must match paragraph count only.
     * 1. Call countParagraphs on [] → 0.
     * 2. Call on mixed callout/list/paragraph blocks.
     * 3. Expect only paragraphs counted (2).
     */
    expect(countParagraphs([])).toBe(0);
    expect(
      countParagraphs([
        { kind: "callout", tokens: [{ t: "text", v: "x" }] },
        { kind: "paragraph", tokens: [{ t: "text", v: "a。" }] },
        {
          kind: "list",
          items: [{ tokens: [{ t: "text", v: "i" }] }],
        },
        { kind: "paragraph", tokens: [{ t: "text", v: "b。" }] },
      ])
    ).toBe(2);
    expect(countParagraphs(blocks)).toBe(2);
  });

  it("indexes paragraphs and sentences", () => {
    /**
     * Behavior: paragraph indices are stable among paragraph blocks only.
     * 1. Build model from two paragraphs.
     * 2. Read sentence texts and keys.
     * 3. Expect ordered sentences and key 1:0 for second paragraph first sentence.
     */
    const model = buildPassageModel(blocks);
    expect(model).toHaveLength(2);
    expect(model[0].sentences.map((s) => s.text)).toEqual(["問題です。", "次は原因です。"]);
    expect(model[1].sentences[0].key).toBe("1:0");
  });

  it("assessGoldCoverage fails closed when gold is missing from plain", () => {
    /**
     * Behavior: missing gold must not yield selectable answer chips.
     * 1. Build model with known sentences.
     * 2. Assess gold that is not in plain; assess gold present in sentence.
     * 3. Expect missing_from_plain vs ok; ensureGoldSelectable always [].
     */
    const model = buildPassageModel(blocks);
    expect(assessGoldCoverage(model, ["まず探しやすさを改善する必要があります。"])).toEqual({
      status: "ok",
    });
    expect(assessGoldCoverage(model, ["存在しない引用"])).toEqual({
      status: "missing_from_plain",
      count: 1,
    });
    expect(ensureGoldSelectable(model, ["存在しない引用"])).toEqual([]);
  });

  it("does not spoil spanning gold that already lives in plain text", () => {
    /**
     * Behavior: multi-sentence gold in plain is ok without chips; multi-select represents it.
     * 1. Two short sentences in one paragraph; assess spanning gold → ok.
     * 2. Model exposes two independently selectable sentences covering the span.
     * 3. Joined quote surfaces compact to the spanning gold (grader compactSpace).
     */
    const model = buildPassageModel([
      {
        kind: "paragraph",
        tokens: [{ t: "text", v: "前文です。後文です。" }],
      },
    ]);
    expect(assessGoldCoverage(model, ["前文です。後文です。"])).toEqual({ status: "ok" });
    expect(assessGoldCoverage(model, ["後文"])).toEqual({ status: "ok" });
    const sents = model[0].sentences.map((s) => s.text);
    expect(sents).toEqual(["前文です。", "後文です。"]);
    const joined = sents.join("");
    expect(joined.includes("前文です。後文です。")).toBe(true);
  });
});

describe("KokugoPassage UI", () => {
  const blocks: Block[] = [
    {
      kind: "paragraph",
      tokens: [{ t: "text", v: "根拠の文です。別の文です。" }],
    },
  ];

  it("toggles sentence selection in sentence-select mode", () => {
    /**
     * Behavior: tapping a sentence toggles selection callback.
     * 1. Render sentence-select with first sentence selected.
     * 2. Click second sentence.
     * 3. Expect onToggle with key 0:1.
     */
    const onToggle = vi.fn();
    render(
      <KokugoPassage
        blocks={blocks}
        mode="sentence-select"
        selectedKeys={new Set(["0:0"])}
        onToggleSentence={onToggle}
      />
    );
    const first = screen.getByRole("button", { name: "根拠の文です。" });
    expect(first).toHaveAttribute("aria-pressed", "true");
    fireEvent.click(screen.getByRole("button", { name: "別の文です。" }));
    expect(onToggle).toHaveBeenCalledWith(
      expect.objectContaining({ key: "0:1", text: "別の文です。" })
    );
  });

  it("changes paragraph role via select", () => {
    /**
     * Behavior: role select reports paragraph index and role.
     * 1. Render paragraph-role with options.
     * 2. Change select to 提案.
     * 3. Expect onRoleChange(0, 提案).
     */
    const onRole = vi.fn();
    render(
      <KokugoPassage
        blocks={blocks}
        mode="paragraph-role"
        roles={["問題"]}
        roleOptions={["問題", "原因", "提案", "結論"]}
        onRoleChange={onRole}
      />
    );
    fireEvent.change(screen.getByLabelText("段落 1 の役割"), {
      target: { value: "提案" },
    });
    expect(onRole).toHaveBeenCalledWith(0, "提案");
  });

  it("does not expose missing gold text; shows corpus inconsistency notice", () => {
    /**
     * Behavior: invalid corpus gold must fail closed in the learner UI.
     * 1. Render sentence-select with goldQuotes not in plain.
     * 2. Assert gold string is not in the document.
     * 3. Assert generic inconsistency status is visible.
     */
    render(
      <KokugoPassage
        blocks={blocks}
        mode="sentence-select"
        selectedKeys={new Set()}
        goldQuotes={["本文に無い金句"]}
      />
    );
    expect(screen.queryByText(/本文に無い金句/)).toBeNull();
    expect(screen.getByRole("status")).toHaveAttribute("data-corpus-inconsistency", "gold_missing");
  });

  it("keeps list/callout blocks visible and stable paragraph indices", () => {
    /**
     * Behavior: non-paragraph Block kinds remain visible; para keys stable.
     * 1. Build mixed callout/paragraph/list/paragraph blocks.
     * 2. Assert model para indices 0,1 and keys 0:0 / 1:0.
     * 3. Render paragraph-role and assert all surfaces + two role selects.
     */
    const mixed: Block[] = [
      { kind: "callout", tone: "info", tokens: [{ t: "text", v: "注意書き" }] },
      { kind: "paragraph", tokens: [{ t: "text", v: "第一段落です。" }] },
      {
        kind: "list",
        items: [{ tokens: [{ t: "text", v: "箇条書き" }] }],
      },
      { kind: "paragraph", tokens: [{ t: "text", v: "第二段落です。" }] },
    ];
    const model = buildPassageModel(mixed);
    expect(model.map((p) => p.paraIndex)).toEqual([0, 1]);
    expect(model[0].sentences[0].key).toBe("0:0");
    expect(model[1].sentences[0].key).toBe("1:0");
    expect(countParagraphs(mixed)).toBe(2);

    render(
      <KokugoPassage
        blocks={mixed}
        mode="paragraph-role"
        roles={Array.from({ length: countParagraphs(mixed) }, () => "問題")}
        roleOptions={["問題", "原因", "提案", "結論"]}
      />
    );
    expect(screen.getByText("注意書き")).toBeVisible();
    expect(screen.getByText("箇条書き")).toBeVisible();
    expect(screen.getByText("第一段落です。")).toBeVisible();
    expect(screen.getByText("第二段落です。")).toBeVisible();
    expect(screen.getAllByLabelText(/段落 \d+ の役割/)).toHaveLength(2);
  });
});
