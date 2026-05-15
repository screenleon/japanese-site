import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import { ClassifierContrasts } from "./ClassifierContrasts";

describe("ClassifierContrasts", () => {
  it("renders contrast cards with cross-level badge", () => {
    render(
      <ClassifierContrasts
        primaryPattern="普通形＋はずだ"
        points={[
          {
            slug: "wakeda",
            title_ja: "わけだ",
            title_zh: "理所當然",
            jlpt_level: "N2",
            schema_version: 2,
            pattern: [{ form: "普通形＋わけだ", gloss_zh: "邏輯上可理解" }],
            explanation_ja_blocks: [{ kind: "paragraph", tokens: [{ t: "text", v: "x" }] }],
            explanation_zh: "x",
            _meta: { source: "curated", license: "CC0-1.0" },
          },
        ]}
        annotations={{
          classifier: {
            rules: [
              {
                with_pattern: "はずがない",
                rule_ja_blocks: [
                  { kind: "paragraph", tokens: [{ t: "text", v: "否定推量との違い。" }] },
                ],
              },
              {
                with_pattern: "わけだ",
                with_slug: "wakeda",
                rule_ja_blocks: [
                  { kind: "paragraph", tokens: [{ t: "text", v: "推量と納得の違い。" }] },
                ],
                examples: [{ use_this: "休みのはずです。", use_alt: "休みなわけです。" }],
              },
            ],
          },
        }}
      />
    );

    expect(screen.getByRole("heading", { name: "辨析" })).toBeVisible();
    expect(screen.getByText("否定推量との違い。")).toBeVisible();
    expect(screen.getByText("推量と納得の違い。")).toBeVisible();
    expect(screen.getByText("N2")).toBeVisible();
    expect(screen.getByText("休みのはずです。")).toBeVisible();
  });
});
