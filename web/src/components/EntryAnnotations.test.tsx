import { render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
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

  it("renders mental_model_zh when explicitly requested", () => {
    render(
      <EntryAnnotations
        annotations={{
          mental_model_zh: "中文提示作為低等級的理解支架。",
        }}
        kinds={["mental_model_zh"]}
      />
    );

    expect(screen.getByRole("heading", { name: "考え方のヒント" })).toBeVisible();
    expect(screen.getByText("中文提示作為低等級的理解支架。")).toBeVisible();
  });

  it("renders nothing when annotations are absent or empty", () => {
    const { container, rerender } = render(<EntryAnnotations />);
    expect(container).toBeEmptyDOMElement();

    rerender(<EntryAnnotations annotations={{ usage: "   " }} />);
    expect(container).toBeEmptyDOMElement();
  });

  it("renders furigana annotations as title tokens", () => {
    /**
     * Verifies structured furigana annotations render title tokens with readings.
     * Steps:
     * 1. Render an entry annotation block with title tokens and vocabulary pairs.
     * 2. Query the generated ruby elements.
     * 3. Assert ruby contains the expected kanji and rt reading.
     */
    // Arrange / Act
    const { container } = render(
      <EntryAnnotations
        annotations={{
          furigana: {
            title_ja: [{ t: "ruby", k: "違", r: "ちが" }],
            vocabulary: [{ kanji: "根拠", reading: "こんきょ" }],
          },
        }}
      />
    );

    // Assert: vocabulary is intentionally NOT rendered in the detached
    // FuriganaBlock — only title_ja furigana surfaces. (vocabulary[]
    // corpus data is retained for future inline-ruby migration.)
    expect(screen.getByText("ふりがな")).toBeVisible();
    const ruby = container.querySelectorAll("ruby");
    expect(ruby).toHaveLength(1);
    expect(ruby[0]).toHaveTextContent("違ちが");
    expect(ruby[0].querySelector("rt")).toHaveTextContent("ちが");
    expect(container.textContent).not.toContain("根拠");
    expect(container.textContent).not.toContain("こんきょ");
  });

  it("renders title tokens with kana context in order", () => {
    /**
     * Verifies title furigana renders inline with surrounding kana context.
     * Steps:
     * 1. Render に違いない as text + ruby title tokens.
     * 2. Query the generated ruby element.
     * 3. Assert the text order and ruby structure match the intended pattern.
     */
    // Arrange / Act
    const { container } = render(
      <EntryAnnotations
        annotations={{
          furigana: {
            title_ja: [
              { t: "text", v: "に" },
              { t: "ruby", k: "違いない", r: "ちがいない" },
            ],
          },
        }}
      />
    );

    // Assert
    expect(screen.getByText("ふりがな")).toBeVisible();
    expect(container.textContent).toContain("に違いないちがいない");
    const ruby = container.querySelectorAll("ruby");
    expect(ruby).toHaveLength(1);
    expect(ruby[0].childNodes[0].textContent).toBe("違いない");
    expect(ruby[0].querySelector("rt")).toHaveTextContent("ちがいない");
    expect(ruby[0].previousSibling?.textContent).toBe("に");
  });

  it("filters blank title text tokens at runtime", () => {
    /**
     * Verifies title annotations do not render when all title tokens are blank text.
     * Steps:
     * 1. Render an entry annotation block with one blank title text token.
     * 2. Assert the furigana annotation block is suppressed.
     */
    // Arrange / Act
    const { container } = render(
      <EntryAnnotations
        annotations={{
          furigana: {
            title_ja: [{ t: "text", v: "   " }],
          },
        }}
      />
    );

    // Assert
    expect(container).toBeEmptyDOMElement();
  });

  it("filters old-shape title furigana at runtime", () => {
    /**
     * Verifies old Pair[] title payloads are not treated as renderable title tokens.
     * Steps:
     * 1. Render an entry annotation block with an old title pair shape.
     * 2. Assert the furigana annotation block is suppressed.
     */
    // Arrange / Act
    const { container } = render(
      <EntryAnnotations
        annotations={
          {
            furigana: {
              title_ja: [
                { kanji: "違", reading: "ちが" },
              ],
              vocabulary: [{ kanji: "根拠", reading: "こんきょ" }],
            },
          } as never
        }
      />
    );

    // Assert
    expect(container).toBeEmptyDOMElement();
  });

  it("filters malformed title tokens at runtime without throwing", () => {
    /**
     * Verifies malformed Token[] payloads are silently filtered before render.
     * Steps:
     * 1. Render an entry annotation block with malformed and valid title tokens.
     * 2. Assert React does not throw or report a render error.
     * 3. Assert only the valid ruby token remains visible.
     */
    // Arrange
    const consoleError = vi.spyOn(console, "error").mockImplementation(() => {});

    try {
      // Act / Assert
      expect(() => {
        render(
          <EntryAnnotations
            annotations={
              {
                furigana: {
                  title_ja: [
                    { t: "ruby", k: null, r: "x" },
                    { t: "text", v: "" },
                    { t: "ruby", k: "違", r: "ちが" },
                  ],
                },
              } as never
            }
          />
        );
      }).not.toThrow();
      expect(consoleError).not.toHaveBeenCalled();
    } finally {
      consoleError.mockRestore();
    }

    expect(screen.getByText("ふりがな")).toBeVisible();
    const ruby = document.querySelectorAll("ruby");
    expect(ruby).toHaveLength(1);
    expect(ruby[0]).toHaveTextContent("違ちが");
    expect(ruby[0].querySelector("rt")).toHaveTextContent("ちが");
    expect(document.body.textContent).not.toContain("x");
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

  it("does not render classifier annotations in the generic renderer", () => {
    render(
      <EntryAnnotations
        annotations={{
          classifier: {
            rules: [
              {
                with_pattern: "わけだ",
                rule_ja_blocks: [
                  { kind: "paragraph", tokens: [{ t: "text", v: "辨析本文" }] },
                ],
              },
            ],
          },
        }}
      />
    );

    expect(screen.queryByText("辨析")).not.toBeInTheDocument();
    expect(screen.queryByText("辨析本文")).not.toBeInTheDocument();
  });
});
