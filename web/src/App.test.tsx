import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { api } from "./api";
import { App } from "./App";
import type { Capabilities, GrammarPoint } from "./apiTypes";

vi.mock("./tabs/QuizTab", () => ({
  QuizTab: () => <section aria-label="quiz panel">quiz panel</section>,
}));

vi.mock("./tabs/VocabTab", () => ({
  VocabTab: ({ onNavigateKanji }: { onNavigateKanji?: (character: string) => void }) => (
    <section aria-label="vocab panel">
      <button type="button" onClick={() => onNavigateKanji?.("食")}>
        漢字へ
      </button>
    </section>
  ),
}));

vi.mock("./tabs/KanjiTab", () => ({
  KanjiTab: ({ initialCharacter }: { initialCharacter?: string }) => (
    <section aria-label="kanji panel" data-character={initialCharacter}>
      kanji panel
    </section>
  ),
}));

vi.mock("./tabs/SentenceTab", () => ({
  SentenceTab: () => <section aria-label="sentence panel">sentence panel</section>,
}));

vi.mock("./tabs/KokugoTab", () => ({
  KokugoTab: () => <section aria-label="kokugo panel">kokugo panel</section>,
}));

vi.mock("./api", () => ({
  api: {
    searchVocab: vi.fn().mockResolvedValue({ results: [], count: 0 }),
    getDueCount: vi.fn().mockResolvedValue({ grammar: 0, vocab: 0 }),
    randomVocab: vi.fn().mockResolvedValue(null),
    getKanji: vi.fn().mockResolvedValue(null),
    randomSentence: vi.fn().mockResolvedValue(null),
    listGrammar: vi.fn().mockResolvedValue({ points: [], count: 0 }),
    randomGrammar: vi.fn().mockResolvedValue(null),
    getGrammar: vi.fn().mockResolvedValue(null),
    getGrammarExamples: vi.fn().mockResolvedValue({ examples: [], count: 0 }),
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
    listKokugoUnits: vi.fn().mockResolvedValue({ units: [], count: 0 }),
    getKokugoSkills: vi.fn().mockResolvedValue({ skills: [], review_queue: [] }),
  },
}));

const getCapabilities = vi.mocked(api.getCapabilities);
const listGrammar = vi.mocked(api.listGrammar);
const getGrammarExamples = vi.mocked(api.getGrammarExamples!);

const grammarPoint: GrammarPoint = {
  slug: "sae",
  title_ja: "〜さえ",
  title_zh: "甚至；只要",
  jlpt_level: "N3",
  schema_version: 2,
  pattern: [{ form: "普通形＋さえ", gloss_zh: "甚至連這種極端例也包含" }],
  explanation_ja_blocks: [{ kind: "paragraph", tokens: [{ t: "text", v: "日本語の説明" }] }],
  explanation_zh: "表示極端例或最低條件。",
  _meta: { source: "curated", license: "CC0-1.0" },
};

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
    localStorage.clear();
    vi.clearAllMocks();
    listGrammar.mockResolvedValue({ points: [grammarPoint], count: 1 });
    getGrammarExamples.mockResolvedValue({ examples: [], count: 0 });
  });

  it("renders all 5 tabs when capabilities load with quiz and sentence enabled", async () => {
    mockCapabilities({
      progress: false,
      history: false,
      quiz: true,
      sentence: true,
      kokugo: false,
    });

    render(<App />);
    fireEvent.click(await screen.findByRole("button", { name: "開始練習" }));

    await waitFor(() => {
      expectVisibleTabs(["練習題", "文法", "單字", "漢字", "例句"]);
      expect(screen.getByLabelText("quiz panel")).toBeVisible();
    });
    expect(screen.queryByRole("button", { name: "国語" })).not.toBeInTheDocument();
  });

  it("shows 国語 tab and panel when kokugo capability is enabled", async () => {
    mockCapabilities({
      progress: true,
      history: false,
      quiz: true,
      sentence: true,
      kokugo: true,
    });

    render(<App />);
    fireEvent.click(await screen.findByRole("button", { name: /^国語教室/ }));

    await waitFor(() => {
      expectVisibleTabs(["練習題", "文法", "單字", "漢字", "例句", "国語"]);
      expect(screen.getByLabelText("kokugo panel")).toBeVisible();
    });
  });

  it("hides quiz and sentence tabs when capabilities report them disabled", async () => {
    mockCapabilities({
      progress: false,
      history: false,
      quiz: false,
      sentence: false,
      kokugo: false,
    });

    render(<App />);
    fireEvent.click(await screen.findByRole("button", { name: /^文法/ }));

    // The reference nav cards remain available even when quiz capability is disabled;
    // after capabilities resolve, quiz and sentence tabs are filtered out.
    await waitFor(() => {
      expect(screen.queryByRole("button", { name: "練習題" })).not.toBeInTheDocument();
      expect(screen.queryByRole("button", { name: "例句" })).not.toBeInTheDocument();
      expect(screen.getByRole("heading", { name: "〜さえ" })).toBeVisible();
    });
    expectVisibleTabs(["文法", "單字", "漢字"]);
    expect(screen.queryByLabelText("quiz panel")).not.toBeInTheDocument();
  });

  it("enters grammar through the nav card while capabilities are loading", async () => {
    let resolveCapabilities: (capabilities: Capabilities) => void = () => {};
    getCapabilities.mockReturnValueOnce(
      new Promise((resolve) => {
        resolveCapabilities = resolve;
      })
    );

    render(<App />);
    fireEvent.click(await screen.findByRole("button", { name: /^文法/ }));

    expect(screen.getByRole("button", { name: "練習題" })).toBeVisible();
    expect(await screen.findByRole("heading", { name: "〜さえ" })).toBeVisible();

    resolveCapabilities({
      progress: false,
      history: false,
      quiz: false,
      sentence: false,
      kokugo: false,
    });

    await waitFor(() => {
      expect(screen.queryByRole("button", { name: "練習題" })).not.toBeInTheDocument();
      expect(screen.getByRole("heading", { name: "〜さえ" })).toBeVisible();
    });
  });

  /**
   * Verifies a vocab-selected kanji opens KanjiTab with that exact character.
   * Steps:
   * 1. Arrange an App with reference tabs available.
   * 2. Act by opening vocab and invoking its kanji navigation action.
   * 3. Assert KanjiTab is active and receives the selected character.
   */
  it("passes a vocab-selected kanji to KanjiTab and opens that tab", async () => {
    mockCapabilities({
      progress: false,
      history: false,
      quiz: false,
      sentence: false,
      kokugo: false,
    });

    render(<App />);
    fireEvent.click(await screen.findByRole("button", { name: /^單字/ }));
    fireEvent.click(screen.getByRole("button", { name: "漢字へ" }));

    expect(await screen.findByLabelText("kanji panel")).toHaveAttribute(
      "data-character",
      "食"
    );
  });
});

describe("Chinese visibility toggle", () => {
  beforeEach(() => {
    localStorage.clear();
    vi.clearAllMocks();
    listGrammar.mockResolvedValue({ points: [grammarPoint], count: 1 });
    getGrammarExamples.mockResolvedValue({ examples: [], count: 0 });
  });

  /**
   * Integration: verifies the App-level header `<ChineseVisibilityToggle />` is wired
   * through the real `<ChineseVisibilityProvider>` and controls Chinese content in
   * the active tab.
   * Steps:
   * 1. Arrange: clear localStorage; mock api to return a fixture grammar entry with
   *    non-empty title_zh / explanation_zh / pattern.gloss_zh.
   * 2. Act: render full <App />, navigate from HomePage to the Grammar tab.
   * 3. Assert (default-hidden): some Chinese surface (e.g. title_zh value or
   *    explanation_zh value or pattern.gloss_zh value) is NOT in the DOM.
   * 4. Act: click the header `中文` button.
   * 5. Assert (toggled-visible): the same Chinese surface IS now in the DOM.
   */
  it("header 中文 toggle reveals Chinese surfaces across tabs", async () => {
    mockCapabilities({
      progress: false,
      history: false,
      quiz: false,
      sentence: false,
      kokugo: false,
    });

    render(<App />);
    fireEvent.click(await screen.findByRole("button", { name: /^文法/ }));

    expect(await screen.findByRole("heading", { name: "〜さえ" })).toBeVisible();
    expect(screen.queryByText("表示極端例或最低條件。")).not.toBeInTheDocument();

    fireEvent.click(screen.getByRole("button", { name: "中文" }));

    expect(await screen.findByText("表示極端例或最低條件。")).toBeVisible();
  });
});
