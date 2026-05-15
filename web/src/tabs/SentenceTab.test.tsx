import { render, screen } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { api, type Sentence } from "../api";
import { ChineseVisibilityProvider } from "../chineseVisibility";
import { SentenceTab } from "./SentenceTab";

vi.mock("../api", () => ({
  api: {
    randomSentence: vi.fn(),
  },
}));

const sentence: Sentence = {
  id: 42,
  text_ja: "山が見えます。",
  text_zh: "看得到山。",
  jlpt_level: "N5",
  source: "tatoeba",
  license: "CC-BY-2.0",
};

const randomSentence = vi.mocked(api.randomSentence);

function renderWithChinese(visible: boolean) {
  return render(
    <ChineseVisibilityProvider initialVisible={visible}>
      <SentenceTab />
    </ChineseVisibilityProvider>
  );
}

describe("SentenceTab", () => {
  beforeEach(() => {
    randomSentence.mockResolvedValue(sentence);
  });

  afterEach(() => {
    vi.clearAllMocks();
    window.localStorage.clear();
  });

  /**
   * Verifies Traditional Chinese sentence text stays hidden without a visibility provider.
   * Steps:
   * 1. Arrange the random sentence API to return a sentence with text_zh.
   * 2. Act by rendering SentenceTab without ChineseVisibilityProvider.
   * 3. Assert text_ja renders while text_zh does not.
   */
  it("hides Chinese sentence text by default", async () => {
    // Arrange / Act
    render(<SentenceTab />);

    // Assert
    expect(await screen.findByText("山が見えます。")).toBeVisible();
    expect(screen.queryByText("看得到山。")).not.toBeInTheDocument();
  });

  /**
   * Verifies Traditional Chinese sentence text renders when visibility is enabled globally.
   * Steps:
   * 1. Arrange the random sentence API to return a sentence with text_zh.
   * 2. Act by rendering SentenceTab inside ChineseVisibilityProvider with initialVisible=true.
   * 3. Assert both text_ja and text_zh render.
   */
  it("shows Chinese sentence text when globally enabled", async () => {
    // Arrange / Act
    renderWithChinese(true);

    // Assert
    expect(await screen.findByText("山が見えます。")).toBeVisible();
    expect(screen.getByText("看得到山。")).toBeVisible();
  });
});
