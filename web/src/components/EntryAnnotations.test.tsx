import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import { EntryAnnotations } from "./EntryAnnotations";

describe("EntryAnnotations", () => {
  it("renders present annotation kinds with Japanese labels", () => {
    render(
      <EntryAnnotations
        annotations={{
          usage: "会話でよく使う。",
          mental_model: "形より働きを見る。",
        }}
      />
    );

    expect(screen.getByText("使い方")).toBeVisible();
    expect(screen.getByText("会話でよく使う。")).toBeVisible();
    expect(screen.getByText("考え方のヒント")).toBeVisible();
    expect(screen.getByText("形より働きを見る。")).toBeVisible();
  });

  it("filters rendered kinds when kinds is provided", () => {
    render(
      <EntryAnnotations
        annotations={{
          usage: "会話でよく使う。",
          nuance_note: "少しくだけた響き。",
        }}
        kinds={["nuance_note"]}
      />
    );

    expect(screen.getByText("ニュアンス")).toBeVisible();
    expect(screen.queryByText("使い方")).not.toBeInTheDocument();
  });

  it("renders nothing when annotations are absent or empty", () => {
    const { container, rerender } = render(<EntryAnnotations />);
    expect(container).toBeEmptyDOMElement();

    rerender(<EntryAnnotations annotations={{ usage: "   " }} />);
    expect(container).toBeEmptyDOMElement();
  });

  it("renders furigana annotations as ruby pairs", () => {
    /**
     * Verifies structured furigana annotations render as ruby pairs with readings.
     * Steps:
     * 1. Render an entry annotation block with title and key-term furigana pairs.
     * 2. Query the generated ruby elements.
     * 3. Assert each ruby contains the expected kanji and rt reading.
     */
    // Arrange / Act
    const { container } = render(
      <EntryAnnotations
        annotations={{
          furigana: {
            title_ja: [{ kanji: "違", reading: "ちが" }],
            key_terms: [{ kanji: "根拠", reading: "こんきょ" }],
          },
        }}
      />
    );

    // Assert
    expect(screen.getByText("ふりがな")).toBeVisible();
    const ruby = container.querySelectorAll("ruby");
    expect(ruby).toHaveLength(2);
    expect(ruby[0]).toHaveTextContent("違ちが");
    expect(ruby[0].querySelector("rt")).toHaveTextContent("ちが");
    expect(ruby[1]).toHaveTextContent("根拠こんきょ");
    expect(ruby[1].querySelector("rt")).toHaveTextContent("こんきょ");
  });

  it("filters blank furigana pairs at runtime", () => {
    /**
     * Verifies mixed furigana annotations render only non-blank ruby pairs.
     * Steps:
     * 1. Render an entry annotation block with one valid and one blank furigana pair.
     * 2. Query the generated ruby elements.
     * 3. Assert only the valid ruby pair is rendered.
     */
    // Arrange / Act
    const { container } = render(
      <EntryAnnotations
        annotations={{
          furigana: {
            title_ja: [
              { kanji: "違", reading: "ちが" },
              { kanji: "   ", reading: "ちが" },
            ],
          },
        }}
      />
    );

    // Assert
    expect(screen.getByText("ふりがな")).toBeVisible();
    const ruby = container.querySelectorAll("ruby");
    expect(ruby).toHaveLength(1);
    expect(ruby[0]).toHaveTextContent("違ちが");
    expect(ruby[0].querySelector("rt")).toHaveTextContent("ちが");
  });

  it("filters non-string furigana pairs at runtime", () => {
    /**
     * Verifies mixed furigana annotations render only pairs with string kanji and reading values.
     * Steps:
     * 1. Render an entry annotation block with one valid and two non-string furigana pairs.
     * 2. Query the generated ruby elements.
     * 3. Assert only the valid ruby pair is rendered.
     */
    // Arrange / Act
    const { container } = render(
      <EntryAnnotations
        annotations={
          {
            furigana: {
              title_ja: [
                { kanji: "違", reading: "ちが" },
                { kanji: null, reading: "ちが" },
              ],
              key_terms: [{ kanji: 123, reading: "こんきょ" }],
            },
          } as never
        }
      />
    );

    // Assert
    expect(screen.getByText("ふりがな")).toBeVisible();
    const ruby = container.querySelectorAll("ruby");
    expect(ruby).toHaveLength(1);
    expect(ruby[0]).toHaveTextContent("違ちが");
    expect(ruby[0].querySelector("rt")).toHaveTextContent("ちが");
  });

  it("ignores unknown runtime annotation keys", () => {
    render(
      <EntryAnnotations
        annotations={
          {
            usage: "会話でよく使う。",
            foo: "bar",
          } as never
        }
      />
    );

    expect(screen.getByText("使い方")).toBeVisible();
    expect(screen.getByText("会話でよく使う。")).toBeVisible();
    expect(screen.queryByText("bar")).not.toBeInTheDocument();
    expect(screen.queryByText("undefined")).not.toBeInTheDocument();
  });
});
