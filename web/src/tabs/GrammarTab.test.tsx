import { act, fireEvent, render, screen, waitFor } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { api } from "../api";
import { GrammarTab } from "./GrammarTab";

vi.mock("../hooks/useReadTracking", () => ({
  useReadTracking: vi.fn(),
}));

vi.mock("../api", () => ({
  api: {
    listGrammar: vi.fn(),
    randomGrammar: vi.fn(),
    getGrammarExamples: vi.fn(),
  },
}));

const points = [
  {
    slug: "sae",
    title_ja: "〜さえ",
    title_zh: "甚至；只要",
    jlpt_level: "N3",
    explanation_zh: "表示極端例或最低條件。",
  },
  {
    slug: "dake",
    title_ja: "〜だけ",
    title_zh: "只有",
    jlpt_level: "N3",
    explanation_zh: "表示限定。",
  },
];

const listGrammar = vi.mocked(api.listGrammar);
const randomGrammar = vi.mocked(api.randomGrammar);
const getGrammarExamples = vi.mocked(api.getGrammarExamples!);

describe("GrammarTab", () => {
  beforeEach(() => {
    listGrammar.mockResolvedValue({ points, count: points.length });
    randomGrammar.mockResolvedValue(points[1]);
    getGrammarExamples.mockResolvedValue({
      examples: [
        {
          id: 1,
          text_ja: "名前さえ書けばいい。",
          text_zh: "只要寫名字就好。",
        },
      ],
      count: 1,
    });
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  it("places the draw-random button inside the active article header", async () => {
    render(<GrammarTab />);

    const button = await screen.findByRole("button", { name: "抽下一個" });
    const article = button.closest("article");
    const header = button.closest("header");

    expect(article).toBeInTheDocument();
    expect(header).toBeInTheDocument();
    expect(header?.closest("article")).toBe(article);
    expect(screen.getByRole("heading", { name: "〜さえ" })).toBeVisible();
  });

  it("disables the draw-random button while a random grammar point is loading", async () => {
    let resolveRandom: (value: (typeof points)[1]) => void = () => {};
    const randomPromise = new Promise<(typeof points)[1]>((resolve) => {
      resolveRandom = resolve;
    });
    randomGrammar.mockReturnValueOnce(randomPromise);
    render(<GrammarTab />);

    fireEvent.click(await screen.findByRole("button", { name: "抽下一個" }));

    const loadingButton = await screen.findByRole("button", { name: "抽取中" });
    expect(loadingButton).toBeDisabled();
    await act(async () => {
      resolveRandom(points[1]);
      await randomPromise;
    });
  });

  it("switches the active grammar slug after drawing a random grammar point", async () => {
    render(<GrammarTab />);

    fireEvent.click(await screen.findByRole("button", { name: "抽下一個" }));

    await waitFor(() => {
      expect(screen.getByRole("heading", { name: "〜だけ" })).toBeVisible();
    });
    expect(randomGrammar).toHaveBeenCalledWith("N3");
  });

  it("loads examples with the active slug and level namespace", async () => {
    render(<GrammarTab />);

    await waitFor(() => {
      expect(getGrammarExamples).toHaveBeenCalledWith("sae", "N3");
    });
    expect(screen.getByText("名前さえ書けばいい。")).toBeVisible();
  });
});
