import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { api } from "./api";
import { HomePage } from "./HomePage";

vi.mock("./api", () => ({
  api: {
    listGrammar: vi.fn(),
    searchVocab: vi.fn(),
    getDueCount: vi.fn(),
  },
}));

const listGrammar = vi.mocked(api.listGrammar);
const searchVocab = vi.mocked(api.searchVocab);
const getDueCount = vi.mocked(api.getDueCount);

describe("HomePage", () => {
  beforeEach(() => {
    vi.stubEnv("VITE_DEPLOY_MODE", "api");
    listGrammar.mockResolvedValue({ points: [], count: 12 });
    searchVocab.mockResolvedValue({ results: [], count: 0, total: 34 });
    getDueCount.mockResolvedValue({ grammar: 0, vocab: 0 });
  });

  afterEach(() => {
    vi.unstubAllEnvs();
    vi.clearAllMocks();
  });

  it("renders the practice subtitle and quiz CTAs in dev mode", async () => {
    render(<HomePage onStart={vi.fn()} quizCapable={true} />);

    expect(
      screen.getByText("用文法、單字與測驗建立穩定的日文練習節奏。")
    ).toBeVisible();
    expect(screen.getByRole("button", { name: "開始練習" })).toBeVisible();
    expect(screen.getByRole("button", { name: "開始測試" })).toBeVisible();
    expect(screen.getByRole("button", { name: /^文法/ })).toBeVisible();
    expect(screen.getByRole("button", { name: /^單字/ })).toBeVisible();
    expect(screen.getByRole("button", { name: /^漢字/ })).toBeVisible();

    await waitFor(() => {
      expect(screen.getByText("12")).toBeVisible();
      expect(screen.getByText("34")).toBeVisible();
    });
    expect(screen.getByText("JLPT N5–N1 文法說明")).toBeVisible();
    expect(screen.getByText("JLPT N5–N1 詞彙與例句")).toBeVisible();
    expect(screen.getByText("讀音・筆畫・部首")).toBeVisible();
    expect(listGrammar).toHaveBeenCalledTimes(1);
    expect(searchVocab).toHaveBeenCalledWith("", undefined);
  });

  it("renders the reference subtitle and hides quiz CTAs in static mode", async () => {
    vi.stubEnv("VITE_DEPLOY_MODE", "static");

    render(<HomePage onStart={vi.fn()} quizCapable={true} />);

    expect(
      screen.getByText("查閱文法說明、單字與漢字，隨時作為學習參考。")
    ).toBeVisible();
    expect(screen.queryByRole("button", { name: "開始練習" })).not.toBeInTheDocument();
    expect(screen.queryByRole("button", { name: "開始測試" })).not.toBeInTheDocument();
    expect(screen.queryByText("需要複習")).not.toBeInTheDocument();
    expect(screen.getByRole("button", { name: /^文法/ })).toBeVisible();
    expect(screen.getByRole("button", { name: /^單字/ })).toBeVisible();
    expect(screen.getByRole("button", { name: /^漢字/ })).toBeVisible();

    await waitFor(() => {
      expect(screen.getByText("12")).toBeVisible();
      expect(screen.getByText("34")).toBeVisible();
    });
    expect(getDueCount).not.toHaveBeenCalled();
  });

  it("hides quiz CTAs when capabilities report quiz disabled in dev mode", async () => {
    render(<HomePage onStart={vi.fn()} quizCapable={false} />);

    await waitFor(() => {
      expect(screen.getByText("12")).toBeVisible();
    });

    expect(screen.queryByRole("button", { name: "開始練習" })).not.toBeInTheDocument();
    expect(screen.queryByRole("button", { name: "開始測試" })).not.toBeInTheDocument();
    expect(screen.queryByText("需要複習")).not.toBeInTheDocument();
    expect(getDueCount).not.toHaveBeenCalled();

    expect(screen.getByRole("button", { name: /^文法/ })).toBeVisible();
    expect(screen.getByRole("button", { name: /^單字/ })).toBeVisible();
    expect(screen.getByRole("button", { name: /^漢字/ })).toBeVisible();
  });

  it.each([
    ["文法", "grammar"],
    ["單字", "vocab"],
    ["漢字", "kanji"],
  ] as const)("calls onStart with %s tab when its card is clicked", async (label, tab) => {
    const onStart = vi.fn();
    render(<HomePage onStart={onStart} quizCapable={true} />);

    fireEvent.click(await screen.findByRole("button", { name: new RegExp(`^${label}`) }));

    expect(onStart).toHaveBeenCalledWith("練習", tab);
  });
});
