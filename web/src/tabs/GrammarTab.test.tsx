import { act, fireEvent, render, screen, waitFor } from "@testing-library/react";
import type { ReactElement } from "react";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { api, ApiError, type GrammarPoint } from "../api";
import { ChineseVisibilityProvider } from "../chineseVisibility";
import { GrammarTab } from "./GrammarTab";

vi.mock("../hooks/useReadTracking", () => ({
  useReadTracking: vi.fn(),
}));

vi.mock("../api", async () => {
  const actual = await vi.importActual<typeof import("../api")>("../api");
  return {
    ...actual,
    api: {
      listGrammar: vi.fn(),
      randomGrammar: vi.fn(),
      getGrammarExamples: vi.fn(),
    },
  };
});

const basePoint: Pick<GrammarPoint, "schema_version" | "pattern" | "explanation_ja_blocks" | "_meta"> = {
  schema_version: 2,
  pattern: [{ form: "普通形＋テスト", gloss_zh: "測試句型" }],
  explanation_ja_blocks: [{ kind: "paragraph", tokens: [{ t: "text", v: "日本語の説明" }] }],
  _meta: { source: "curated", license: "CC0-1.0" },
};

const points = [
  {
    ...basePoint,
    slug: "sae",
    title_ja: "〜さえ",
    title_zh: "甚至；只要",
    jlpt_level: "N3",
    explanation_zh: "表示極端例或最低條件。",
  },
  {
    ...basePoint,
    slug: "dake",
    title_ja: "〜だけ",
    title_zh: "只有",
    jlpt_level: "N3",
    explanation_zh: "表示限定。",
  },
  {
    ...basePoint,
    slug: "mono-no",
    title_ja: "ものの",
    title_zh: "ものの（雖然／但是）",
    jlpt_level: "N3",
    annotations: {
      nuance_note: "口語・くだけた逆接。前文の予想と異なる結果を続ける。",
      mental_model:
        "前件を事実として置いたうえで、後件で予想から外れる結果を示す。逆接を一つの流れとして読むことで、単なる接続詞ではなく判断の向きを意識できる。",
    },
    explanation_zh: "雖然但是。",
  },
  {
    ...basePoint,
    slug: "zzz-decoy-n2-foo",
    title_ja: "デコイ",
    title_zh: "decoy N2 entry",
    jlpt_level: "N2",
    explanation_zh: "如果 related-slug navigation 沒有設定 active slug，會落到這一筆。",
  },
  {
    ...basePoint,
    slug: "te-iru",
    title_ja: "ている",
    title_zh: "ている",
    jlpt_level: "N3",
    annotations: {
      mental_model: "nested mental model wins",
      nuance_note: "nested nuance renders",
    },
    explanation_zh: "正在／狀態。",
  },
  {
    ...basePoint,
    slug: "empty-nested",
    title_ja: "空文字",
    title_zh: "empty nested",
    jlpt_level: "N3",
    annotations: {
      mental_model: "",
    },
    explanation_zh: "空文字の遷移確認。",
  },
];

const lowLevelPoint: GrammarPoint = {
  ...basePoint,
  slug: "ka-question",
  title_ja: "か",
  title_zh: "疑問助詞",
  jlpt_level: "N5",
  annotations: {
    mental_model: "「か」は、文を開いた問いに変える合図です。",
    mental_model_zh: "句尾的「か」把陳述打開成問題，留下等待回答的空位。",
  },
  explanation_zh: "疑問。",
};

const listGrammar = vi.mocked(api.listGrammar);
const randomGrammar = vi.mocked(api.randomGrammar);
const getGrammarExamples = vi.mocked(api.getGrammarExamples!);

function renderWithChinese(visible: boolean, ui: ReactElement) {
  return render(
    <ChineseVisibilityProvider initialVisible={visible}>
      {ui}
    </ChineseVisibilityProvider>
  );
}

async function expectMainEntryWithoutExamples(errorText: RegExp) {
  expect(await screen.findByRole("heading", { name: "ものの" })).toBeVisible();
  expect(screen.getByText("日本語の説明")).toBeVisible();
  expect(
    screen.getByText("口語・くだけた逆接。前文の予想と異なる結果を続ける。")
  ).toBeVisible();
  expect(screen.getByRole("heading", { name: "考え方のヒント" })).toBeVisible();
  // Fixture has no related_slugs; section must not appear after cross-level dedup.
  expect(screen.queryByText("相關用法")).not.toBeInTheDocument();
  expect(screen.queryByText("例文")).not.toBeInTheDocument();
  expect(screen.queryByText("例文を表示できません。")).not.toBeInTheDocument();
  expect(screen.queryByText(errorText)).not.toBeInTheDocument();
  expect(screen.queryByRole("alert")).toBeNull();
}

describe("GrammarTab", () => {
  let warn: ReturnType<typeof vi.spyOn>;

  beforeEach(() => {
    warn = vi.spyOn(console, "warn").mockImplementation(() => {});
    listGrammar.mockResolvedValue({ points, count: points.length });
    randomGrammar.mockResolvedValue(points[1]);
    getGrammarExamples.mockResolvedValue({
      examples: [
        {
          id: 1,
          text_ja: "名前さえ書けばいい。",
          text_zh: "只要寫名字就好。",
        },
      ],
      count: 1,
    });
  });

  afterEach(() => {
    vi.clearAllMocks();
    warn.mockRestore();
  });

  it("places the draw-random button inside the active article header", async () => {
    render(<GrammarTab />);

    const button = await screen.findByRole("button", { name: "抽下一個" });
    const article = button.closest("article");
    const header = button.closest("header");

    expect(article).toBeInTheDocument();
    expect(header).toBeInTheDocument();
    expect(header?.closest("article")).toBe(article);
    expect(screen.getByRole("heading", { name: "〜さえ" })).toBeVisible();
  });

  /**
   * Verifies GrammarTab adopts content-first DOM order via ContentFirstLayout.
   * Steps:
   * 1. Render GrammarTab with mocked grammar points.
   * 2. Wait for the active entry heading.
   * 3. Assert study-directory follows study-content in document order.
   */
  it("places study content before the directory in DOM order (content-first)", async () => {
    render(<GrammarTab />);
    await screen.findByRole("heading", { name: "〜さえ" });

    const content = screen.getByTestId("study-content");
    const directory = screen.getByTestId("study-directory");
    const following = Node.DOCUMENT_POSITION_FOLLOWING;
    expect(content.compareDocumentPosition(directory) & following).toBe(following);
  });

  /**
   * Verifies picking a grammar directory item focuses that entry and closes the mobile directory.
   * Steps:
   * 1. Render GrammarTab and wait for the default active heading.
   * 2. Open the mobile directory and click the ものの entry.
   * 3. Assert ものの is the active heading and the directory toggle is collapsed.
   */
  it("selecting a directory entry activates it and collapses the mobile directory", async () => {
    render(<GrammarTab />);
    await screen.findByRole("heading", { name: "〜さえ" });

    fireEvent.click(screen.getByRole("button", { name: /本級目錄/ }));
    fireEvent.click(screen.getByRole("button", { name: "ものの" }));

    expect(await screen.findByRole("heading", { name: "ものの" })).toBeVisible();
    expect(screen.getByRole("button", { name: /本級目錄/ })).toHaveAttribute(
      "aria-expanded",
      "false"
    );
  });

  it("disables the draw-random button while a random grammar point is loading", async () => {
    let resolveRandom: (value: (typeof points)[1]) => void = () => {};
    const randomPromise = new Promise<(typeof points)[1]>((resolve) => {
      resolveRandom = resolve;
    });
    randomGrammar.mockReturnValueOnce(randomPromise);
    render(<GrammarTab />);

    fireEvent.click(await screen.findByRole("button", { name: "抽下一個" }));

    const loadingButton = await screen.findByRole("button", { name: "抽取中" });
    expect(loadingButton).toBeDisabled();
    await act(async () => {
      resolveRandom(points[1]);
      await randomPromise;
    });
  });

  it("switches the active grammar slug after drawing a random grammar point", async () => {
    render(<GrammarTab />);

    fireEvent.click(await screen.findByRole("button", { name: "抽下一個" }));

    await waitFor(() => {
      expect(screen.getByRole("heading", { name: "〜だけ" })).toBeVisible();
    });
    expect(randomGrammar).toHaveBeenCalledWith("N3");
  });

  it("loads examples with only the active slug", async () => {
    render(<GrammarTab />);

    await waitFor(() => {
      expect(getGrammarExamples).toHaveBeenCalledWith("sae");
    });
    expect(screen.getByText("名前さえ書けばいい。")).toBeVisible();
  });

  /**
   * Verifies Chinese grammar surfaces stay hidden under the default provider state.
   * Steps:
   * 1. Arrange: mock grammar and examples API responses with Chinese title, pattern, example, and explanation text.
   * 2. Act: render GrammarTab without enabling the Chinese visibility provider.
   * 3. Assert: wait for the active Japanese heading so the tab has loaded.
   * 4. Assert: Chinese title, pattern gloss, example translation, and explanation block are absent.
   */
  it("hides Chinese grammar surfaces by default", async () => {
    render(<GrammarTab />);

    expect(await screen.findByRole("heading", { name: "〜さえ" })).toBeVisible();
    expect(screen.queryByText("甚至；只要")).not.toBeInTheDocument();
    expect(screen.queryByText("測試句型")).not.toBeInTheDocument();
    expect(screen.queryByText("只要寫名字就好。")).not.toBeInTheDocument();
    expect(screen.queryByText("中文說明")).not.toBeInTheDocument();
    expect(screen.queryByText("表示極端例或最低條件。")).not.toBeInTheDocument();
    expect(screen.getByText("N3 · 〜さえ")).toBeVisible();
  });

  /**
   * Verifies Chinese grammar surfaces render when global Chinese visibility is enabled.
   * Steps:
   * 1. Arrange: mock grammar and examples API responses with Chinese title, pattern, example, and explanation text.
   * 2. Act: render GrammarTab inside ChineseVisibilityProvider with initialVisible=true.
   * 3. Assert: wait for the active Japanese heading and async example translation.
   * 4. Assert: Chinese title, pattern gloss, example translation, and explanation block are visible.
   */
  it("shows Chinese grammar surfaces when globally enabled", async () => {
    renderWithChinese(true, <GrammarTab />);

    expect(await screen.findByRole("heading", { name: "〜さえ" })).toBeVisible();
    // example.text_zh comes from getGrammarExamples (separate async call) — await it specifically.
    expect(await screen.findByText("只要寫名字就好。")).toBeVisible();
    expect(screen.getAllByText("甚至；只要").length).toBeGreaterThanOrEqual(1);
    expect(screen.getByText("測試句型")).toBeVisible();
    expect(screen.getByText("中文說明")).toBeVisible();
    expect(screen.getByText("表示極端例或最低條件。")).toBeVisible();
    expect(screen.getByText("N3 · 甚至；只要")).toBeVisible();
  });

  it("examples 404 — page renders main entry without crashing", async () => {
    getGrammarExamples.mockRejectedValue(
      new ApiError(404, "Not Found", "not_found")
    );

    render(<GrammarTab initialSlug="mono-no" />);

    await expectMainEntryWithoutExamples(/not_found/);
    expect(getGrammarExamples).toHaveBeenCalledWith("mono-no");
    expect(getGrammarExamples).toHaveBeenCalledTimes(1);
    expect(warn).toHaveBeenCalledWith(
      "grammar examples fetch failed",
      expect.objectContaining({ status: 404, code: "not_found" })
    );
  });

  it("examples 500 — page renders main entry without crashing", async () => {
    getGrammarExamples.mockRejectedValue(
      new ApiError(500, "Server Error", "http_error")
    );

    render(<GrammarTab initialSlug="mono-no" />);

    expect(await screen.findByRole("heading", { name: "ものの" })).toBeVisible();
    expect(await screen.findByText("例文を表示できません。")).toBeVisible();
    expect(screen.queryByText(/http_error/)).not.toBeInTheDocument();
    expect(screen.queryByRole("alert")).toBeNull();
    expect(getGrammarExamples).toHaveBeenCalledWith("mono-no");
    expect(getGrammarExamples).toHaveBeenCalledTimes(1);
    expect(warn).toHaveBeenCalledWith(
      "grammar examples fetch failed",
      expect.objectContaining({ status: 500, code: "http_error" })
    );
  });

  it("examples network reject — page renders main entry without crashing", async () => {
    getGrammarExamples.mockRejectedValue(new Error("network down"));

    render(<GrammarTab initialSlug="mono-no" />);

    expect(await screen.findByRole("heading", { name: "ものの" })).toBeVisible();
    expect(await screen.findByText("例文を表示できません。")).toBeVisible();
    expect(screen.queryByText(/network down/)).not.toBeInTheDocument();
    expect(screen.queryByRole("alert")).toBeNull();
    expect(getGrammarExamples).toHaveBeenCalledWith("mono-no");
    expect(getGrammarExamples).toHaveBeenCalledTimes(1);
    expect(warn).toHaveBeenCalledWith(
      "grammar examples fetch failed",
      expect.objectContaining({ message: "network down" })
    );
  });

  it("renders nuance_note for the active grammar point", async () => {
    render(<GrammarTab initialSlug="mono-no" />);

    expect(
      await screen.findByText("口語・くだけた逆接。前文の予想と異なる結果を続ける。")
    ).toBeVisible();
  });

  it("renders mental_model for the active grammar point", async () => {
    render(<GrammarTab initialSlug="mono-no" />);

    expect(await screen.findByRole("heading", { name: "考え方のヒント" })).toBeVisible();
    expect(
      screen.getByText(
        "前件を事実として置いたうえで、後件で予想から外れる結果を示す。逆接を一つの流れとして読むことで、単なる接続詞ではなく判断の向きを意識できる。"
      )
    ).toBeVisible();
  });

  it("renders low-level mental_model_zh as primary before Japanese when Chinese is enabled", async () => {
    listGrammar.mockResolvedValueOnce({ points: [lowLevelPoint, ...points], count: points.length + 1 });

    const { container } = renderWithChinese(true, <GrammarTab initialSlug="ka-question" />);

    expect(await screen.findByRole("heading", { name: "か" })).toBeVisible();
    const chinese = screen.getByText("句尾的「か」把陳述打開成問題，留下等待回答的空位。");
    const japanese = screen.getByText("「か」は、文を開いた問いに変える合図です。");
    expect(chinese).toBeVisible();
    expect(japanese).toBeVisible();
    expect(chinese.compareDocumentPosition(japanese) & Node.DOCUMENT_POSITION_FOLLOWING).toBeTruthy();
    expect(chinese).toHaveClass("text-base");
    expect(japanese).toHaveClass("text-sm");
    expect(container.querySelectorAll("aside h3")).toHaveLength(1);
  });

  it("keeps N3 mental_model Japanese-only when Chinese is enabled and mental_model_zh is absent", async () => {
    renderWithChinese(true, <GrammarTab initialSlug="mono-no" />);

    expect(await screen.findByRole("heading", { name: "ものの" })).toBeVisible();
    expect(
      screen.getByText(
        "前件を事実として置いたうえで、後件で予想から外れる結果を示す。逆接を一つの流れとして読むことで、単なる接続詞ではなく判断の向きを意識できる。"
      )
    ).toBeVisible();
    expect(screen.queryByText("句尾的「か」把陳述打開成問題，留下等待回答的空位。")).not.toBeInTheDocument();
  });

  it("keeps N5 mental_model Japanese-only when Chinese is disabled", async () => {
    listGrammar.mockResolvedValueOnce({ points: [lowLevelPoint, ...points], count: points.length + 1 });

    renderWithChinese(false, <GrammarTab initialSlug="ka-question" />);

    expect(await screen.findByRole("heading", { name: "か" })).toBeVisible();
    expect(screen.getByText("「か」は、文を開いた問いに変える合図です。")).toBeVisible();
    expect(screen.queryByText("句尾的「か」把陳述打開成問題，留下等待回答的空位。")).not.toBeInTheDocument();
  });

  it("omits the mental_model heading when the active grammar point has none", async () => {
    render(<GrammarTab initialSlug="sae" />);

    expect(await screen.findByRole("heading", { name: "〜さえ" })).toBeVisible();
    expect(
      screen.queryByText(
        "前件を事実として置いたうえで、後件で予想から外れる結果を示す。逆接を一つの流れとして読むことで、単なる接続詞ではなく判断の向きを意識できる。"
      )
    ).not.toBeInTheDocument();
  });

  it("prefers nested annotations over flat transition fields", async () => {
    render(<GrammarTab initialSlug="te-iru" />);

    expect(await screen.findByText("nested mental model wins")).toBeVisible();
    expect(screen.getByText("nested nuance renders")).toBeVisible();
    expect(screen.queryByText("flat mental model should not render when nested exists")).not.toBeInTheDocument();
  });

  it("does not fall back to flat fields when nested annotations contain an empty string", async () => {
    render(<GrammarTab initialSlug="empty-nested" />);

    expect(await screen.findByRole("heading", { name: "空文字" })).toBeVisible();
    expect(
      screen.queryByText("flat mental model must not replace empty nested")
    ).not.toBeInTheDocument();
  });

  it("omits related grammar buttons when no related slugs are present", async () => {
    render(<GrammarTab initialSlug="mono-no" />);

    expect(await screen.findByRole("heading", { name: "ものの" })).toBeVisible();

    expect(screen.queryByText("相關用法")).not.toBeInTheDocument();
    expect(screen.queryByRole("button", { name: "〜ものの (N2)" })).not.toBeInTheDocument();
    expect(getGrammarExamples).toHaveBeenLastCalledWith("mono-no");
  });
});
