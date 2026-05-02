import { render, screen, waitFor } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { api } from "./api";
import { App } from "./App";
import type { Capabilities } from "./apiTypes";

vi.mock("./tabs/QuizTab", () => ({
  QuizTab: () => <section aria-label="quiz panel">quiz panel</section>,
}));

vi.mock("./tabs/GrammarTab", () => ({
  GrammarTab: () => <section aria-label="grammar panel">grammar panel</section>,
}));

vi.mock("./tabs/VocabTab", () => ({
  VocabTab: () => <section aria-label="vocab panel">vocab panel</section>,
}));

vi.mock("./tabs/KanjiTab", () => ({
  KanjiTab: () => <section aria-label="kanji panel">kanji panel</section>,
}));

vi.mock("./tabs/SentenceTab", () => ({
  SentenceTab: () => <section aria-label="sentence panel">sentence panel</section>,
}));

vi.mock("./api", () => ({
  api: {
    searchVocab: vi.fn().mockResolvedValue({ results: [], count: 0 }),
    randomVocab: vi.fn().mockResolvedValue(null),
    getKanji: vi.fn().mockResolvedValue(null),
    randomSentence: vi.fn().mockResolvedValue(null),
    listGrammar: vi.fn().mockResolvedValue({ points: [], count: 0 }),
    randomGrammar: vi.fn().mockResolvedValue(null),
    getGrammar: vi.fn().mockResolvedValue(null),
    nextQuestion: vi.fn().mockResolvedValue(null),
    answer: vi.fn().mockResolvedValue(null),
    stats: vi.fn().mockResolvedValue({
      total_attempts: 0,
      correct: 0,
      wrong: 0,
      accuracy: 0,
      by_grammar: [],
      by_error_class: [],
      recent_wrong: [],
    }),
    markRead: vi.fn().mockResolvedValue(undefined),
    getProgress: vi.fn().mockResolvedValue({
      content_type: "grammar",
      read: 0,
      total: 0,
      percent: 0,
    }),
    getCapabilities: vi.fn(),
  },
}));

const getCapabilities = vi.mocked(api.getCapabilities);

function mockCapabilities(capabilities: Capabilities) {
  getCapabilities.mockResolvedValueOnce(capabilities);
}

function expectVisibleTabs(labels: string[]) {
  for (const label of labels) {
    expect(screen.getByRole("button", { name: label })).toBeVisible();
  }
}

describe("App tab filtering", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("renders all 5 tabs when capabilities load with quiz and sentence enabled", async () => {
    mockCapabilities({
      progress: false,
      history: false,
      quiz: true,
      sentence: true,
    });

    render(<App />);

    await waitFor(() => {
      expectVisibleTabs(["練習題", "文法", "單字", "漢字", "例句"]);
      expect(screen.getByLabelText("quiz panel")).toBeVisible();
    });
  });

  it("hides quiz and sentence tabs when capabilities report them disabled", async () => {
    mockCapabilities({
      progress: false,
      history: false,
      quiz: false,
      sentence: false,
    });

    render(<App />);

    // Wait for the auto-reset effect to land grammar panel after capabilities resolve.
    // Tabs disappear in render N; setActive("grammar") fires in the post-commit
    // effect so grammar panel renders in render N+1. Both bullets must be in
    // waitFor or the second assertion races the effect.
    await waitFor(() => {
      expect(screen.queryByRole("button", { name: "練習題" })).not.toBeInTheDocument();
      expect(screen.queryByRole("button", { name: "例句" })).not.toBeInTheDocument();
      expect(screen.getByLabelText("grammar panel")).toBeVisible();
    });
    expectVisibleTabs(["文法", "單字", "漢字"]);
    expect(screen.queryByLabelText("quiz panel")).not.toBeInTheDocument();
  });

  it("falls back to grammar when initial active tab disappears", async () => {
    let resolveCapabilities: (capabilities: Capabilities) => void = () => {};
    getCapabilities.mockReturnValueOnce(
      new Promise((resolve) => {
        resolveCapabilities = resolve;
      })
    );

    render(<App />);

    expect(screen.getByRole("button", { name: "練習題" })).toBeVisible();

    resolveCapabilities({
      progress: false,
      history: false,
      quiz: false,
      sentence: false,
    });

    await waitFor(() => {
      expect(screen.queryByRole("button", { name: "練習題" })).not.toBeInTheDocument();
      expect(screen.getByLabelText("grammar panel")).toBeVisible();
    });
  });
});
