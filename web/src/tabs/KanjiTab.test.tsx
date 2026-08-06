import { fireEvent, render, screen } from "@testing-library/react";
import { StrictMode } from "react";
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

  /**
   * Verifies an initial character hydrates the form, loads its record, and is consumed once.
   * Steps:
   * 1. Arrange KanjiTab with 山 and a consumption callback.
   * 2. Act by rendering the tab and allowing its initial lookup to resolve.
   * 3. Assert the input, API call, and callback each reflect 山 exactly once.
   */
  it("loads and consumes an initial character once", async () => {
    const onInitialCharacterConsumed = vi.fn();
    render(
      <KanjiTab
        initialCharacter="山"
        onInitialCharacterConsumed={onInitialCharacterConsumed}
      />
    );

    expect(await screen.findByText("やま。")).toBeVisible();
    expect(screen.getByPlaceholderText("輸入單一漢字（例：山、愛、食）")).toHaveValue("山");
    expect(getKanji).toHaveBeenCalledWith("山");
    expect(onInitialCharacterConsumed).toHaveBeenCalledTimes(1);
  });

  /**
   * Verifies a failed initial deep-link lookup preserves the requested character and reports an error.
   * Steps:
   * 1. Arrange getKanji to reject while KanjiTab receives 山 as its initial character.
   * 2. Act by rendering the tab and allowing the automatic lookup to fail.
   * 3. Assert the input still shows 山, the failure is visible, and no kanji card is rendered.
   */
  it("shows an error when an initial character lookup fails", async () => {
    getKanji.mockRejectedValueOnce(new Error("lookup failed"));
    render(<KanjiTab initialCharacter="山" />);

    expect(await screen.findByText("Error: lookup failed")).toBeVisible();
    expect(screen.getByPlaceholderText("輸入單一漢字（例：山、愛、食）")).toHaveValue("山");
    expect(screen.queryByText("やま。")).not.toBeInTheDocument();
  });

  /**
   * Verifies an unchanged initial character is not consumed again when its callback identity changes.
   * Steps:
   * 1. Arrange KanjiTab with 山 and an initial callback.
   * 2. Act by rerendering with the same character and a new callback.
   * 3. Assert the lookup and first callback run once while the new callback is not invoked.
   */
  it("does not repeat an unchanged initial character lookup after rerender", async () => {
    const firstCallback = vi.fn();
    const { rerender } = render(
      <KanjiTab initialCharacter="山" onInitialCharacterConsumed={firstCallback} />
    );

    await screen.findByText("やま。");
    const replacementCallback = vi.fn();
    rerender(
      <KanjiTab initialCharacter="山" onInitialCharacterConsumed={replacementCallback} />
    );

    expect(getKanji).toHaveBeenCalledTimes(1);
    expect(firstCallback).toHaveBeenCalledTimes(1);
    expect(replacementCallback).not.toHaveBeenCalled();
  });

  /**
   * Verifies StrictMode effect replay does not duplicate a one-time deep-link lookup.
   * Steps:
   * 1. Arrange KanjiTab under StrictMode with 山 and a consumption callback.
   * 2. Act by rendering and allowing the initial lookup to resolve.
   * 3. Assert the API lookup and consumption callback each occur once.
   */
  it("consumes an initial character once under StrictMode", async () => {
    const onInitialCharacterConsumed = vi.fn();
    render(
      <StrictMode>
        <KanjiTab
          initialCharacter="山"
          onInitialCharacterConsumed={onInitialCharacterConsumed}
        />
      </StrictMode>
    );

    expect(await screen.findByText("やま。")).toBeVisible();
    expect(getKanji).toHaveBeenCalledTimes(1);
    expect(onInitialCharacterConsumed).toHaveBeenCalledTimes(1);
  });
});
