import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import { BlockRenderer } from "../components/EntryAnnotations";

describe("BlockRenderer", () => {
  it("renders every spike block and token kind", () => {
    const { container } = render(
      <BlockRenderer
        blocks={[
          {
            kind: "paragraph",
            tokens: [
              { t: "text", v: "plain" },
              { t: "ruby", k: "根拠", r: "こんきょ" },
              { t: "term", kind: "grammar", slug: "hazuda", label: "はずだ" },
            ],
          },
          { kind: "list", items: [{ tokens: [{ t: "text", v: "item" }] }] },
          { kind: "callout", tone: "tip", tokens: [{ t: "text", v: "tip text" }] },
        ]}
      />
    );

    expect(screen.getByText("plain")).toBeVisible();
    expect(container.querySelector("ruby")).toHaveTextContent("根拠こんきょ");
    expect(screen.getByRole("link", { name: "はずだ" })).toBeVisible();
    expect(screen.getByText("item")).toBeVisible();
    expect(screen.getByText("tip text")).toBeVisible();
  });
});
