import { fireEvent, render, screen } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { api, type Kanji } from "../api";
import { ChineseVisibilityProvider } from "../chineseVisibility";
import { KanjiTab } from "./KanjiTab";

vi.mock("../hooks/useReadTracking", () => ({
  useReadTracking: vi.fn(),
}));

vi.mock("../api", () => ({
  api: {
    getKanji: vi.fn(),
  },
}));

const kanji: Kanji = {
  id: 1,
  character: "山",
  onyomi: "サン",
  kunyomi: "やま",
  meaning_ja: "やま。",
  meaning_zh: "山、山脈。",
  jlpt_level: "N5",
  grade: 1,
  stroke_count: 3,
  source: "kanjidic2",
  license: "EDRDG",
};

const getKanji = vi.mocked(api.getKanji);

function renderWithChinese(visible: boolean) {
  return render(
    <ChineseVisibilityProvider initialVisible={visible}>
      <KanjiTab />
    </ChineseVisibilityProvider>
  );
}

async function lookupKanji() {
  fireEvent.change(screen.getByPlaceholderText("輸入單一漢字（例：山、愛、食）"), {
    target: { value: "山" },
  });
  fireEvent.click(screen.getByRole("button", { name: "查詢" }));

  await screen.findByText("山");
}

describe("KanjiTab", () => {
  beforeEach(() => {
    getKanji.mockResolvedValue(kanji);
  });

  afterEach(() => {
    vi.clearAllMocks();
    window.localStorage.clear();
  });

  /**
   * Verifies Traditional Chinese kanji meaning stays hidden without a visibility provider.
   * Steps:
   * 1. Arrange the kanji lookup API to return a row with meaning_zh.
   * 2. Act by rendering KanjiTab without ChineseVisibilityProvider and submitting a lookup.
   * 3. Assert Japanese/reading fields render while the 繁中 label and meaning_zh do not.
   */
  it("hides Chinese kanji meaning by default", async () => {
    // Arrange / Act
    render(<KanjiTab />);
    await lookupKanji();

    // Assert
    expect(screen.getByText("やま。")).toBeVisible();
    expect(screen.getByText("サン")).toBeVisible();
    expect(screen.getByText("やま")).toBeVisible();
    expect(screen.queryByText("繁中")).not.toBeInTheDocument();
    expect(screen.queryByText("山、山脈。")).not.toBeInTheDocument();
  });

  /**
   * Verifies Traditional Chinese kanji meaning renders when visibility is enabled globally.
   * Steps:
   * 1. Arrange the kanji lookup API to return a row with meaning_zh.
   * 2. Act by rendering KanjiTab inside ChineseVisibilityProvider with initialVisible=true.
   * 3. Assert the 繁中 label and meaning_zh are visible with the kanji entry.
   */
  it("shows Chinese kanji meaning when globally enabled", async () => {
    // Arrange / Act
    renderWithChinese(true);
    await lookupKanji();

    // Assert
    expect(screen.getByText("やま。")).toBeVisible();
    expect(screen.getByText("繁中")).toBeVisible();
    expect(screen.getByText("山、山脈。")).toBeVisible();
  });
});
