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
});
