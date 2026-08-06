import { fireEvent, render, screen } from "@testing-library/react";
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

const multiKanjiVocab = {
  id: 3,
  headword: "食べ物",
  reading: "たべもの",
  pos: "名詞",
  gloss_ja: "食べるもの。",
  gloss_zh: "食物。",
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

  /**
   * Verifies usage annotations appear only on the focus card, not compact directory rows.
   * Steps:
   * 1. Render VocabTab with a mocked vocab row that includes annotations.
   * 2. Wait for the 使い方 label and collect annotation text nodes.
   * 3. Assert exactly one annotation instance and that it lives under vocab-focus-card, not the directory list.
   */
  it("renders annotations on the focus vocab card only (not compact directory rows)", async () => {
    render(<VocabTab />);

    expect(await screen.findAllByText("使い方")).not.toHaveLength(0);
    const notes = screen.getAllByText(
      "「お喋り」は雑談にも、話しすぎる様子にも使う。"
    );
    // Annotations live on the focus card; directory rows stay compact.
    expect(notes).toHaveLength(1);
    expect(screen.getByTestId("vocab-focus-card")).toContainElement(notes[0]);
    expect(screen.getByTestId("vocab-directory-list")).not.toContainElement(notes[0]);
  });

  /**
   * Verifies VocabTab adopts content-first DOM order via ContentFirstLayout.
   * Steps:
   * 1. Render VocabTab with mocked vocab results.
   * 2. Wait for gloss text to confirm load.
   * 3. Assert study-directory follows study-content in document order.
   */
  it("places focus content before the directory in DOM order", async () => {
    render(<VocabTab />);
    await screen.findAllByText("よく話すこと。");

    const content = screen.getByTestId("study-content");
    const directory = screen.getByTestId("study-directory");
    const following = Node.DOCUMENT_POSITION_FOLLOWING;
    expect(content.compareDocumentPosition(directory) & following).toBe(following);
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

    // Assert — focus card + compact directory row both show headword ruby.
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
   * Verifies a focused single-kanji headword emits its selected character.
   * Steps:
   * 1. Arrange a vocab focus card with お喋り and a navigation spy.
   * 2. Act by clicking the accessible 喋 lookup button.
   * 3. Assert the callback receives 喋.
   */
  it("navigates to the selected kanji from the focus headword", async () => {
    const onNavigateKanji = vi.fn();
    render(<VocabTab onNavigateKanji={onNavigateKanji} />);

    await screen.findAllByText("よく話すこと。");
    fireEvent.click(screen.getAllByRole("button", { name: "漢字「喋」を調べる" })[0]);

    expect(onNavigateKanji).toHaveBeenCalledWith("喋");
  });

  /**
   * Verifies every kanji in a multi-kanji headword has an independent lookup action.
   * Steps:
   * 1. Arrange a focus card with 食べ物 and a navigation spy.
   * 2. Act by clicking the 食 and 物 lookup buttons independently.
   * 3. Assert the callback receives 食 then 物.
   */
  it("navigates each kanji in a multi-kanji focus headword independently", async () => {
    searchVocab.mockResolvedValue({ results: [multiKanjiVocab], count: 1, total: 1 });
    randomVocab.mockResolvedValue(multiKanjiVocab);
    const onNavigateKanji = vi.fn();
    render(<VocabTab onNavigateKanji={onNavigateKanji} />);

    await screen.findAllByText("食べるもの。");
    fireEvent.click(screen.getByRole("button", { name: "漢字「食」を調べる" }));
    fireEvent.click(screen.getByRole("button", { name: "漢字「物」を調べる" }));

    expect(onNavigateKanji).toHaveBeenNthCalledWith(1, "食");
    expect(onNavigateKanji).toHaveBeenNthCalledWith(2, "物");
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
    // Focus card (full 中文說明) + compact directory row both surface Chinese gloss.
    expect(screen.getAllByText("聊天、多話。").length).toBeGreaterThanOrEqual(2);
  });

  /**
   * Verifies selecting a directory vocab row focuses that entry and collapses the mobile directory.
   * Steps:
   * 1. Mock two list rows and render VocabTab.
   * 2. Open the mobile directory and click the 勉強 row.
   * 3. Assert the focus card shows 学ぶこと。 and the directory toggle is collapsed.
   */
  it("selecting a directory row moves it into the focus card and collapses the mobile directory", async () => {
    const second = {
      ...vocab,
      id: 99,
      headword: "勉強",
      reading: "べんきょう",
      gloss_ja: "学ぶこと。",
      gloss_zh: "學習。",
      annotations: undefined,
    };
    searchVocab.mockResolvedValue({ results: [vocab, second], count: 2, total: 2 });
    randomVocab.mockResolvedValue(vocab);

    render(<VocabTab />);
    await screen.findByText("学ぶこと。");

    // Open mobile directory, then pick the second row.
    fireEvent.click(screen.getByRole("button", { name: /本級目錄/ }));
    fireEvent.click(screen.getByRole("button", { name: /勉強/ }));

    expect(screen.getByTestId("vocab-focus-card")).toHaveTextContent("学ぶこと。");
    expect(screen.getByRole("button", { name: /本級目錄/ })).toHaveAttribute(
      "aria-expanded",
      "false"
    );
  });
});
