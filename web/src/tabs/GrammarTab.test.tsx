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
    slug: "monono",
    title_ja: "ものの",
    title_zh: "ものの（雖然／但是）",
    jlpt_level: "N3",
    annotations: {
      nuance_note: "口語・くだけた逆接。前文の予想と異なる結果を続ける。",
      mental_model:
        "前件を事実として置いたうえで、後件で予想から外れる結果を示す。逆接を一つの流れとして読むことで、単なる接続詞ではなく判断の向きを意識できる。",
    },
    related_slugs: ["monono-formal"],
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
    slug: "monono-formal",
    title_ja: "〜ものの",
    title_zh: "〜ものの（雖然…但是…）",
    jlpt_level: "N2",
    annotations: {
      nuance_note: "文語的な逆接で、書き言葉や改まった場面で使う。",
    },
    related_slugs: ["monono"],
    explanation_zh: "書面逆接。",
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
  expect(screen.getByText("相關用法")).toBeVisible();
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

  it("hides Chinese grammar surfaces by default", async () => {
    render(<GrammarTab />);

    expect(await screen.findByRole("heading", { name: "〜さえ" })).toBeVisible();
    expect(screen.queryByText("甚至；只要")).not.toBeInTheDocument();
    expect(screen.queryByText("測試句型")).not.toBeInTheDocument();
    expect(screen.queryByText("只要寫名字就好。")).not.toBeInTheDocument();
    expect(screen.getByText("N3 · 〜さえ")).toBeVisible();
  });

  it("shows Chinese grammar surfaces when globally enabled", async () => {
    renderWithChinese(true, <GrammarTab />);

    expect(await screen.findByRole("heading", { name: "〜さえ" })).toBeVisible();
    expect(screen.getAllByText("甚至；只要").length).toBeGreaterThanOrEqual(1);
    expect(screen.getByText("測試句型")).toBeVisible();
    expect(screen.getByText("只要寫名字就好。")).toBeVisible();
    expect(screen.getByText("N3 · 甚至；只要")).toBeVisible();
  });

  it("examples 404 — page renders main entry without crashing", async () => {
    getGrammarExamples.mockRejectedValue(
      new ApiError(404, "Not Found", "not_found")
    );

    render(<GrammarTab initialSlug="monono" />);

    await expectMainEntryWithoutExamples(/not_found/);
    expect(getGrammarExamples).toHaveBeenCalledWith("monono");
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

    render(<GrammarTab initialSlug="monono" />);

    expect(await screen.findByRole("heading", { name: "ものの" })).toBeVisible();
    expect(await screen.findByText("例文を表示できません。")).toBeVisible();
    expect(screen.queryByText(/http_error/)).not.toBeInTheDocument();
    expect(screen.queryByRole("alert")).toBeNull();
    expect(getGrammarExamples).toHaveBeenCalledWith("monono");
    expect(getGrammarExamples).toHaveBeenCalledTimes(1);
    expect(warn).toHaveBeenCalledWith(
      "grammar examples fetch failed",
      expect.objectContaining({ status: 500, code: "http_error" })
    );
  });

  it("examples network reject — page renders main entry without crashing", async () => {
    getGrammarExamples.mockRejectedValue(new Error("network down"));

    render(<GrammarTab initialSlug="monono" />);

    expect(await screen.findByRole("heading", { name: "ものの" })).toBeVisible();
    expect(await screen.findByText("例文を表示できません。")).toBeVisible();
    expect(screen.queryByText(/network down/)).not.toBeInTheDocument();
    expect(screen.queryByRole("alert")).toBeNull();
    expect(getGrammarExamples).toHaveBeenCalledWith("monono");
    expect(getGrammarExamples).toHaveBeenCalledTimes(1);
    expect(warn).toHaveBeenCalledWith(
      "grammar examples fetch failed",
      expect.objectContaining({ message: "network down" })
    );
  });

  it("renders nuance_note for the active grammar point", async () => {
    render(<GrammarTab initialSlug="monono" />);

    expect(
      await screen.findByText("口語・くだけた逆接。前文の予想と異なる結果を続ける。")
    ).toBeVisible();
  });

  it("renders mental_model for the active grammar point", async () => {
    render(<GrammarTab initialSlug="monono" />);

    expect(await screen.findByRole("heading", { name: "考え方のヒント" })).toBeVisible();
    expect(
      screen.getByText(
        "前件を事実として置いたうえで、後件で予想から外れる結果を示す。逆接を一つの流れとして読むことで、単なる接続詞ではなく判断の向きを意識できる。"
      )
    ).toBeVisible();
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

  it("renders related grammar buttons and navigates to the variant", async () => {
    render(<GrammarTab initialSlug="monono" />);

    expect(await screen.findByRole("heading", { name: "ものの" })).toBeVisible();
    expect(screen.getByText("相關用法")).toBeVisible();

    fireEvent.click(screen.getByRole("button", { name: "〜ものの (N2)" }));

    await waitFor(() => {
      expect(screen.getByRole("heading", { name: "〜ものの" })).toBeVisible();
    });
    expect(screen.queryByRole("heading", { name: "デコイ" })).not.toBeInTheDocument();
    expect(getGrammarExamples).toHaveBeenLastCalledWith("monono-formal");
  });
});
