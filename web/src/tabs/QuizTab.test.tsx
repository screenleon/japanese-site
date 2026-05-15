import { fireEvent, render, screen } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import {
  api,
  type GradeResult,
  type GrammarPoint,
  type Question,
  type Stats,
} from "../api";
import { ChineseVisibilityProvider } from "../chineseVisibility";
import { QuizTab } from "./QuizTab";

vi.mock("../api", async () => {
  const actual = await vi.importActual<typeof import("../api")>("../api");
  return {
    ...actual,
    api: {
      listGrammar: vi.fn(),
      nextQuestion: vi.fn(),
      answer: vi.fn(),
      stats: vi.fn(),
    },
  };
});

const grammarPoint: GrammarPoint = {
  slug: "wake-da",
  title_ja: "わけだ",
  title_zh: "怪不得",
  jlpt_level: "N3",
  schema_version: 2,
  pattern: [{ form: "普通形＋わけだ", gloss_zh: "自然推論" }],
  explanation_ja_blocks: [{ kind: "paragraph", tokens: [{ t: "text", v: "日本語の説明" }] }],
  explanation_zh: "事情自然如此的文法說明。",
  _meta: { source: "curated", license: "CC0-1.0" },
};

const question: Question = {
  id: "question-1",
  content_type: "grammar",
  kind: "cloze",
  jlpt_level: "N3",
  grammar_point: "wake-da",
  prompt: "雨が降った___。",
  hint: "自然な結論",
};

const emptyStats: Stats = {
  total_attempts: 0,
  correct: 0,
  wrong: 0,
  accuracy: 0,
  by_grammar: [],
  by_vocab: [],
  by_error_class: [],
  recent_wrong: [],
};

const grammarResult: GradeResult = {
  correct: false,
  user_answer: "はず",
  expected: "わけだ",
  explanation_zh: "這裡是在說自然得出的結論。",
  grammar_point: "wake-da",
  error_class: "grammar_confusion",
  suggested_next: ["wake-da"],
  content_type: "grammar",
  item_detail_zh: "文法項目的繁體中文詳細說明。",
};

const vocabResult: GradeResult = {
  correct: false,
  user_answer: "錯誤",
  expected: "正解",
  explanation_zh: "這裡要選符合語境的單字。",
  grammar_point: "vocab:山",
  error_class: "vocab_confusion",
  suggested_next: ["山"],
  content_type: "vocab",
  item_detail_zh: "單字的繁體中文詞義。",
};

const listGrammar = vi.mocked(api.listGrammar);
const nextQuestion = vi.mocked(api.nextQuestion);
const answer = vi.mocked(api.answer);
const stats = vi.mocked(api.stats);

function renderWithChinese(visible: boolean) {
  return render(
    <ChineseVisibilityProvider initialVisible={visible}>
      <QuizTab />
    </ChineseVisibilityProvider>
  );
}

async function submitAnswer(result: GradeResult) {
  answer.mockResolvedValue(result);

  await screen.findByText("雨が降った");
  fireEvent.change(screen.getByRole("textbox"), {
    target: { value: "はず" },
  });
  fireEvent.click(screen.getByRole("button", { name: "送出" }));

  await screen.findByText("不對");
}

describe("QuizTab", () => {
  beforeEach(() => {
    listGrammar.mockResolvedValue({ points: [grammarPoint], count: 1 });
    nextQuestion.mockResolvedValue(question);
    stats.mockResolvedValue(emptyStats);
  });

  afterEach(() => {
    vi.clearAllMocks();
    window.localStorage.clear();
  });

  /**
   * Verifies ResultBox hides grader explanation_zh without a visibility provider.
   * Steps:
   * 1. Arrange quiz API mocks so answering returns a wrong result with explanation_zh.
   * 2. Act by rendering QuizTab without ChineseVisibilityProvider and submitting an answer.
   * 3. Assert the result appears while the Traditional Chinese explanation does not.
   */
  it("hides Chinese grading explanation by default", async () => {
    // Arrange / Act
    render(<QuizTab />);
    await submitAnswer(grammarResult);

    // Assert
    expect(screen.getByText("不對")).toBeVisible();
    expect(screen.queryByText("這裡是在說自然得出的結論。")).not.toBeInTheDocument();
  });

  /**
   * Verifies ResultBox shows grader explanation_zh when visibility is enabled globally.
   * Steps:
   * 1. Arrange quiz API mocks so answering returns a wrong result with explanation_zh.
   * 2. Act by rendering QuizTab inside ChineseVisibilityProvider with initialVisible=true and submitting an answer.
   * 3. Assert the Traditional Chinese explanation is visible in the result.
   */
  it("shows Chinese grading explanation when globally enabled", async () => {
    // Arrange / Act
    renderWithChinese(true);
    await submitAnswer(grammarResult);

    // Assert
    expect(screen.getByText("這裡是在說自然得出的結論。")).toBeVisible();
  });

  /**
   * Verifies grammar item_detail_zh stays hidden without a visibility provider.
   * Steps:
   * 1. Arrange quiz API mocks so answering returns a grammar result with item_detail_zh.
   * 2. Act by rendering QuizTab without ChineseVisibilityProvider and submitting an answer.
   * 3. Assert the grammar detail label and Traditional Chinese detail are absent.
   */
  it("hides Chinese grammar item detail by default", async () => {
    // Arrange / Act
    render(<QuizTab />);
    await submitAnswer(grammarResult);

    // Assert
    expect(screen.queryByText("文法說明")).not.toBeInTheDocument();
    expect(screen.queryByText("文法項目的繁體中文詳細說明。")).not.toBeInTheDocument();
  });

  /**
   * Verifies grammar item_detail_zh renders when visibility is enabled globally.
   * Steps:
   * 1. Arrange quiz API mocks so answering returns a grammar result with item_detail_zh.
   * 2. Act by rendering QuizTab inside ChineseVisibilityProvider with initialVisible=true and submitting an answer.
   * 3. Assert the grammar detail label and Traditional Chinese detail are visible.
   */
  it("shows Chinese grammar item detail when globally enabled", async () => {
    // Arrange / Act
    renderWithChinese(true);
    await submitAnswer(grammarResult);

    // Assert
    expect(screen.getByText("文法說明")).toBeVisible();
    expect(screen.getByText("文法項目的繁體中文詳細說明。")).toBeVisible();
  });

  /**
   * Verifies vocab item_detail_zh stays hidden without a visibility provider.
   * Steps:
   * 1. Arrange quiz API mocks so answering returns a vocab result with item_detail_zh.
   * 2. Act by rendering QuizTab without ChineseVisibilityProvider and submitting an answer.
   * 3. Assert the vocab detail label and Traditional Chinese detail are absent.
   */
  it("hides Chinese vocab item detail by default", async () => {
    // Arrange / Act
    render(<QuizTab />);
    await submitAnswer(vocabResult);

    // Assert
    expect(screen.queryByText("詞義")).not.toBeInTheDocument();
    expect(screen.queryByText("單字的繁體中文詞義。")).not.toBeInTheDocument();
  });

  /**
   * Verifies vocab item_detail_zh renders when visibility is enabled globally.
   * Steps:
   * 1. Arrange quiz API mocks so answering returns a vocab result with item_detail_zh.
   * 2. Act by rendering QuizTab inside ChineseVisibilityProvider with initialVisible=true and submitting an answer.
   * 3. Assert the vocab detail label and Traditional Chinese detail are visible.
   */
  it("shows Chinese vocab item detail when globally enabled", async () => {
    // Arrange / Act
    renderWithChinese(true);
    await submitAnswer(vocabResult);

    // Assert
    expect(screen.getByText("詞義")).toBeVisible();
    expect(screen.getByText("單字的繁體中文詞義。")).toBeVisible();
  });
});
