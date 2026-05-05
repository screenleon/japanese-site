import { render, screen } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { api } from "../api";
import { VocabTab } from "./VocabTab";

vi.mock("../hooks/useReadTracking", () => ({
  useReadTracking: vi.fn(),
}));

vi.mock("../api", () => ({
  api: {
    searchVocab: vi.fn(),
    randomVocab: vi.fn(),
  },
}));

const vocab = {
  id: 1,
  headword: "お喋り",
  reading: "おしゃべり",
  pos: "",
  gloss_ja: "よく話すこと。",
  gloss_zh: "聊天、多話。",
  jlpt_level: "N3",
  annotations: {
    usage: "「お喋り」は雑談にも、話しすぎる様子にも使う。",
  },
  source: "curated",
  license: "CC-BY-SA-4.0",
};

const searchVocab = vi.mocked(api.searchVocab);
const randomVocab = vi.mocked(api.randomVocab);

describe("VocabTab", () => {
  beforeEach(() => {
    searchVocab.mockResolvedValue({ results: [vocab], count: 1, total: 1 });
    randomVocab.mockResolvedValue(vocab);
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  it("renders annotations on the random vocab card", async () => {
    render(<VocabTab />);

    expect(await screen.findAllByText("使い方")).not.toHaveLength(0);
    expect(
      screen.getAllByText("「お喋り」は雑談にも、話しすぎる様子にも使う。")[0]
    ).toBeVisible();
  });

  it("renders annotations on search result rows", async () => {
    render(<VocabTab />);

    const notes = await screen.findAllByText(
      "「お喋り」は雑談にも、話しすぎる様子にも使う。"
    );
    expect(notes.length).toBeGreaterThanOrEqual(2);
  });
});
