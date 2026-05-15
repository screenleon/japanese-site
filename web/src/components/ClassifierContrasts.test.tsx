import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import { ChineseVisibilityProvider } from "../chineseVisibility";
import { ClassifierContrasts } from "./ClassifierContrasts";

describe("ClassifierContrasts", () => {
  function renderContrasts(visible = false) {
    return render(
      <ChineseVisibilityProvider initialVisible={visible}>
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
                  rule_zh: "中文辨析。",
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
      </ChineseVisibilityProvider>
    );
  }

  /**
   * Verifies Chinese classifier contrast copy stays hidden under the default provider state.
   * Steps:
   * 1. Arrange: render contrast cards with Japanese rule blocks and Chinese contrast copy.
   * 2. Act: use the default ChineseVisibilityProvider state.
   * 3. Assert: Japanese contrast content, examples, and cross-level badge are visible.
   * 4. Assert: the Chinese contrast copy is absent from the DOM.
   */
  it("renders contrast cards with cross-level badge", () => {
    renderContrasts();

    expect(screen.getByRole("heading", { name: "辨析" })).toBeVisible();
    expect(screen.getByText("否定推量との違い。")).toBeVisible();
    expect(screen.getByText("推量と納得の違い。")).toBeVisible();
    expect(screen.getByText("N2")).toBeVisible();
    expect(screen.getByText("休みのはずです。")).toBeVisible();
    expect(screen.queryByText("中文辨析。")).not.toBeInTheDocument();
  });

  /**
   * Verifies Chinese classifier contrast copy renders when global Chinese visibility is enabled.
   * Steps:
   * 1. Arrange: render contrast cards with Japanese rule blocks and Chinese contrast copy.
   * 2. Act: enable ChineseVisibilityProvider with initialVisible=true.
   * 3. Assert: the Chinese contrast copy is visible.
   */
  it("shows Chinese contrast copy when globally enabled", () => {
    renderContrasts(true);

    expect(screen.getByText("中文辨析。")).toBeVisible();
  });
});
