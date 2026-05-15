import { render, screen } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { api } from "../api";
import { ChineseVisibilityProvider } from "../chineseVisibility";
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

const kanaOnlyVocab = {
  id: 2,
  headword: "ありがとう",
  reading: "ありがとう",
  pos: "",
  gloss_ja: "感謝を表すことば。",
  gloss_zh: "謝謝。",
  jlpt_level: "N5",
  source: "curated",
  license: "CC-BY-SA-4.0",
};

const searchVocab = vi.mocked(api.searchVocab);
const randomVocab = vi.mocked(api.randomVocab);

function renderWithChinese(visible: boolean) {
  return render(
    <ChineseVisibilityProvider initialVisible={visible}>
      <VocabTab />
    </ChineseVisibilityProvider>
  );
}

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

  it("renders kanji vocabulary readings with ruby", async () => {
    /**
     * Verifies vocab headwords with distinct readings render ruby reading markup.
     * Steps:
     * 1. Render the vocab tab with mocked kanji vocabulary API results.
     * 2. Wait for the vocab card and row content to appear.
     * 3. Assert the rendered ruby contains the headword and rt reading.
     */
    // Arrange / Act
    const { container } = render(<VocabTab />);

    await screen.findAllByText("よく話すこと。");

    // Assert
    const ruby = Array.from(container.querySelectorAll("ruby"));
    expect(ruby.length).toBeGreaterThanOrEqual(2);
    expect(ruby[0]).toHaveTextContent("お喋りおしゃべり");
    expect(ruby[0].querySelector("rt")).toHaveTextContent("おしゃべり");
  });

  it("renders kana-only vocabulary without ruby", async () => {
    /**
     * Verifies vocab headwords with matching kana readings render without ruby markup.
     * Steps:
     * 1. Render the vocab tab with mocked kana-only vocabulary API results.
     * 2. Wait for the vocab card and row content to appear.
     * 3. Assert the rendered output does not contain ruby markup.
     */
    // Arrange
    searchVocab.mockResolvedValue({ results: [kanaOnlyVocab], count: 1, total: 1 });
    randomVocab.mockResolvedValue(kanaOnlyVocab);

    // Act
    const { container } = render(<VocabTab />);

    await screen.findAllByText("ありがとう");

    // Assert
    expect(container.querySelector("ruby")).not.toBeInTheDocument();
  });

  /**
   * Verifies Chinese vocab glosses stay hidden under the default provider state.
   * Steps:
   * 1. Arrange: mock vocab API responses with Japanese and Chinese gloss text.
   * 2. Act: render VocabTab without enabling the Chinese visibility provider.
   * 3. Assert: wait for the Japanese gloss to confirm the tab has loaded.
   * 4. Assert: the Chinese gloss is absent from the DOM.
   */
  it("hides Chinese glosses by default", async () => {
    render(<VocabTab />);

    expect(await screen.findAllByText("よく話すこと。")).not.toHaveLength(0);
    expect(screen.queryByText("聊天、多話。")).not.toBeInTheDocument();
  });

  /**
   * Verifies Chinese vocab glosses render when global Chinese visibility is enabled.
   * Steps:
   * 1. Arrange: mock vocab API responses with Japanese and Chinese gloss text.
   * 2. Act: render VocabTab inside ChineseVisibilityProvider with initialVisible=true.
   * 3. Assert: wait for the Japanese gloss to confirm the tab has loaded.
   * 4. Assert: the Chinese gloss appears in both expected vocab surfaces.
   */
  it("shows Chinese glosses when globally enabled", async () => {
    renderWithChinese(true);

    expect(await screen.findAllByText("よく話すこと。")).not.toHaveLength(0);
    expect(screen.getAllByText("聊天、多話。").length).toBeGreaterThanOrEqual(2);
  });
});
