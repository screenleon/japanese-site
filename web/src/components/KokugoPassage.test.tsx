import { fireEvent, render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import type { Block } from "../apiTypes";
import {
  buildPassageModel,
  ensureGoldSelectable,
  KokugoPassage,
  plainFromTokens,
  splitJapaneseSentences,
} from "./KokugoPassage";

describe("plainFromTokens / splitJapaneseSentences", () => {
  it("flattens text and ruby base", () => {
    expect(
      plainFromTokens([
        { t: "text", v: "図書" },
        { t: "ruby", k: "室", r: "しつ" },
        { t: "text", v: "です。" },
      ])
    ).toBe("図書室です。");
  });

  it("splits on Japanese terminators and keeps them", () => {
    expect(splitJapaneseSentences("Aです。Bます？C！")).toEqual(["Aです。", "Bます？", "C！"]);
  });
});

describe("buildPassageModel / ensureGoldSelectable", () => {
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

  it("indexes paragraphs and sentences", () => {
    const model = buildPassageModel(blocks);
    expect(model).toHaveLength(2);
    expect(model[0].sentences.map((s) => s.text)).toEqual(["問題です。", "次は原因です。"]);
    expect(model[1].sentences[0].key).toBe("1:0");
  });

  it("adds gold extras only when not covered by a sentence", () => {
    const model = buildPassageModel(blocks);
    const covered = ensureGoldSelectable(model, ["まず探しやすさを改善する必要があります。"]);
    expect(covered).toHaveLength(0);
    const missing = ensureGoldSelectable(model, ["存在しない引用"]);
    expect(missing).toHaveLength(1);
    expect(missing[0].text).toBe("存在しない引用");
    expect(missing[0].paraIndex).toBe(-1);
  });

  it("does not spoil spanning gold that already lives in plain text", () => {
    const model = buildPassageModel([
      {
        kind: "paragraph",
        tokens: [{ t: "text", v: "前文です。後文です。" }],
      },
    ]);
    // Spans two sentences — multi-select covers grading; no answer chip.
    const spanning = ensureGoldSelectable(model, ["前文です。後文です。"]);
    expect(spanning).toHaveLength(0);
    const substring = ensureGoldSelectable(model, ["後文"]);
    expect(substring).toHaveLength(0);
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

  it("renders orphan gold fallback chip and toggles it", () => {
    const onToggle = vi.fn();
    render(
      <KokugoPassage
        blocks={blocks}
        mode="sentence-select"
        selectedKeys={new Set()}
        onToggleSentence={onToggle}
        goldQuotes={["本文に無い金句"]}
      />
    );
    const chip = screen.getByRole("button", { name: /本文に無い金句/ });
    expect(chip).toHaveAttribute("data-sentence-key", "gold:orphan:0");
    fireEvent.click(chip);
    expect(onToggle).toHaveBeenCalledWith(
      expect.objectContaining({ key: "gold:orphan:0", paraIndex: -1, text: "本文に無い金句" })
    );
  });

  it("keeps list/callout blocks visible and stable paragraph indices", () => {
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

    render(
      <KokugoPassage
        blocks={mixed}
        mode="paragraph-role"
        roles={["問題", "結論"]}
        roleOptions={["問題", "原因", "提案", "結論"]}
      />
    );
    expect(screen.getByText("注意書き")).toBeVisible();
    expect(screen.getByText("箇条書き")).toBeVisible();
    expect(screen.getByText("第一段落です。")).toBeVisible();
    expect(screen.getByText("第二段落です。")).toBeVisible();
    expect(screen.getByLabelText("段落 1 の役割")).toBeVisible();
    expect(screen.getByLabelText("段落 2 の役割")).toBeVisible();
  });
});
