import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { api } from "../api";
import { CapabilitiesProvider } from "../capabilities";
import type { KokugoUnit } from "../kokugoTypes";
import { deriveResumeState, KokugoTab } from "./KokugoTab";

vi.mock("../api", () => ({
  api: {
    listKokugoUnits: vi.fn(),
    getKokugoUnit: vi.fn(),
    getKokugoUnitState: vi.fn(),
    putKokugoProgress: vi.fn(),
    submitKokugoTask: vi.fn(),
    saveKokugoArtifact: vi.fn(),
    getCapabilities: vi.fn(),
  },
}));

const listKokugoUnits = vi.mocked(api.listKokugoUnits);
const getKokugoUnit = vi.mocked(api.getKokugoUnit);
const getKokugoUnitState = vi.mocked(api.getKokugoUnitState);
const putKokugoProgress = vi.mocked(api.putKokugoProgress);
const submitKokugoTask = vi.mocked(api.submitKokugoTask);
const saveKokugoArtifact = vi.mocked(api.saveKokugoArtifact);
const getCapabilities = vi.mocked(api.getCapabilities);

const sampleUnit: KokugoUnit = {
  id: "library-use",
  schema_version: 1,
  stage: "e5-6",
  title_ja: "学校の図書室をもっと使いやすくするには",
  genre: "expository",
  objectives: ["主張を見つける"],
  estimated_minutes: 20,
  text: [
    {
      kind: "paragraph",
      tokens: [{ t: "text", v: "図書室は大切です。" }],
    },
  ],
  support: { default_profile: "n3" },
  tasks: [
    {
      id: "predict-1",
      skill: "reading.predict",
      kind: "predict",
      payload: {
        prompt_ja: "何について？",
        choices: [
          { id: "a", text_ja: "歴史" },
          { id: "b", text_ja: "工夫" },
        ],
      },
    },
    {
      id: "summary-1",
      skill: "reading.summary",
      kind: "summary-choice",
      payload: {
        prompt_ja: "要約は？",
        choices: [
          { id: "a", text_ja: "悪い要約" },
          { id: "b", text_ja: "良い要約" },
        ],
        correct_id: "b",
      },
    },
  ],
  artifact: {
    kind: "short-proposal",
    min_chars: 0,
    max_chars: 0,
    checklist: ["提案がある", "理由がある"],
    exemplar_ja: "例です。",
  },
  _meta: { source: "original", license: "CC0-1.0", validated_by: "test" },
};

function renderTab(caps: { progress: boolean; kokugo: boolean } = { progress: true, kokugo: true }) {
  getCapabilities.mockResolvedValue({
    progress: caps.progress,
    history: false,
    quiz: true,
    sentence: true,
    kokugo: caps.kokugo,
  });
  return render(
    <CapabilitiesProvider>
      <KokugoTab />
    </CapabilitiesProvider>
  );
}

describe("deriveResumeState", () => {
  it("restores revise phase with draft body from saved state", () => {
    const r = deriveResumeState(sampleUnit, {
      progress: {
        unit_key: "e5-6/library-use",
        stage: "e5-6",
        unit_id: "library-use",
        status: "in_progress",
        step: "revise",
        started_at: "",
        updated_at: "",
      },
      attempts: [],
      artifacts: [
        {
          unit_key: "e5-6/library-use",
          revision: 0,
          body: "下書き本文",
          checklist: [true, true],
          created_at: "",
          version: 1,
          updated_at: "t1",
        },
      ],
    });
    expect(r.phase).toBe("revise");
    expect(r.draftBody).toBe("下書き本文");
    expect(r.draftVersion).toBe(1);
    expect(r.checklist).toEqual([true, true]);
  });

  it("marks completed units as done", () => {
    const r = deriveResumeState(sampleUnit, {
      progress: {
        unit_key: "e5-6/library-use",
        stage: "e5-6",
        unit_id: "library-use",
        status: "completed",
        step: "done",
        started_at: "",
        updated_at: "",
        completed_at: "now",
      },
      attempts: [],
      artifacts: [
        {
          unit_key: "e5-6/library-use",
          revision: 0,
          body: "下書き",
          checklist: [true, true],
          created_at: "",
          version: 1,
          updated_at: "t0",
        },
        {
          unit_key: "e5-6/library-use",
          revision: 1,
          body: "改稿",
          checklist: [true, true],
          created_at: "",
          version: 1,
          updated_at: "t1",
        },
      ],
    });
    expect(r.phase).toBe("done");
    expect(r.revisionBody).toBe("改稿");
  });
});

describe("KokugoTab", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    listKokugoUnits.mockResolvedValue({
      units: [
        {
          id: "library-use",
          stage: "e5-6",
          title_ja: "学校の図書室をもっと使いやすくするには",
          genre: "expository",
          estimated_minutes: 20,
          task_count: 2,
          has_artifact: true,
        },
      ],
      count: 1,
    });
    getKokugoUnit.mockResolvedValue(sampleUnit);
    getKokugoUnitState.mockResolvedValue({ attempts: [], artifacts: [] });
    putKokugoProgress.mockResolvedValue({
      unit_key: "e5-6/library-use",
      stage: "e5-6",
      unit_id: "library-use",
      status: "in_progress",
      step: "predict",
      started_at: "",
      updated_at: "",
    });
    submitKokugoTask.mockResolvedValue({
      attempt: {
        id: 1,
        unit_key: "e5-6/library-use",
        task_id: "predict-1",
        answer: {},
        created_at: "",
      },
      grade: { correct: null, explanation_ja: "予測を記録しました。" },
    });
  });

  it("lists units and opens the predict step for a fresh unit", async () => {
    renderTab();
    expect(await screen.findByText("国語教室")).toBeVisible();
    fireEvent.click(await screen.findByRole("button", { name: /学校の図書室/ }));
    await waitFor(() => {
      expect(screen.getByText("読む前の予測")).toBeVisible();
    });
    expect(getKokugoUnit).toHaveBeenCalledWith("e5-6", "library-use");
    expect(getKokugoUnitState).toHaveBeenCalledWith("e5-6", "library-use");
    expect(putKokugoProgress).toHaveBeenCalledWith("e5-6", "library-use", {
      step: "predict",
      status: "in_progress",
    });
  });

  it("restores saved draft when reopening a unit", async () => {
    getKokugoUnitState.mockResolvedValue({
      progress: {
        unit_key: "e5-6/library-use",
        stage: "e5-6",
        unit_id: "library-use",
        status: "in_progress",
        step: "revise",
        started_at: "",
        updated_at: "",
      },
      attempts: [],
      artifacts: [
        {
          unit_key: "e5-6/library-use",
          revision: 0,
          body: "保存済み下書きです。",
          checklist: [true, false],
          created_at: "",
          version: 3,
          updated_at: "2026-08-03 10:00:00",
        },
      ],
    });
    renderTab();
    fireEvent.click(await screen.findByRole("button", { name: /学校の図書室/ }));
    await waitFor(() => {
      expect(screen.getByText("改稿")).toBeVisible();
    });
    expect(screen.getByDisplayValue("保存済み下書きです。")).toBeVisible();
    // Should not reset progress to predict
    expect(putKokugoProgress).not.toHaveBeenCalledWith(
      "e5-6",
      "library-use",
      expect.objectContaining({ step: "predict" })
    );
  });

  it("advances predict → read → task and submits a task", async () => {
    renderTab();
    fireEvent.click(await screen.findByRole("button", { name: /学校の図書室/ }));
    await screen.findByText("読む前の予測");
    fireEvent.click(screen.getByLabelText("工夫"));
    fireEvent.click(screen.getByRole("button", { name: /予測を記録/ }));
    await waitFor(() => {
      expect(screen.getByText("本文を読む")).toBeVisible();
    });
    fireEvent.click(screen.getByRole("button", { name: "課題へ進む" }));
    await waitFor(() => {
      expect(screen.getByText("要約")).toBeVisible();
    });
    submitKokugoTask.mockResolvedValueOnce({
      attempt: {
        id: 2,
        unit_key: "e5-6/library-use",
        task_id: "summary-1",
        answer: {},
        created_at: "",
      },
      grade: { correct: true, explanation_ja: "正しい要約を選べました。" },
    });
    fireEvent.click(screen.getByLabelText("良い要約"));
    fireEvent.click(screen.getByRole("button", { name: "提出" }));
    await waitFor(() => {
      expect(submitKokugoTask).toHaveBeenCalledWith(
        "e5-6",
        "library-use",
        "summary-1",
        { choice_id: "b" }
      );
    });
    await waitFor(() => {
      expect(screen.getByText("作品（下書き）")).toBeVisible();
    });
  });

  it("JS-133: evidence highlight submits quotes from in-passage sentence taps", async () => {
    /**
     * Behavior: in-passage taps build quotes[]; submit disabled until selection.
     * 1. Open unit, reach evidence step; assert 提出 disabled and no submit yet.
     * 2. Tap gold sentence once, submit; assert quotes payload.
     * 3. (Separate unit reopen path covered by deselect test below.)
     */
    const evidenceUnit: KokugoUnit = {
      ...sampleUnit,
      text: [
        {
          kind: "paragraph",
          tokens: [
            {
              t: "text",
              v: "図書室は大切です。まず探しやすさを改善する必要があります。別の文です。",
            },
          ],
        },
      ],
      tasks: [
        {
          id: "predict-1",
          skill: "reading.predict",
          kind: "predict",
          payload: {
            prompt_ja: "何について？",
            choices: [
              { id: "a", text_ja: "歴史" },
              { id: "b", text_ja: "工夫" },
            ],
          },
        },
        {
          id: "evidence-1",
          skill: "reading.locate-evidence",
          kind: "evidence-highlight",
          payload: {
            prompt_ja: "根拠の文を選びなさい。",
            gold_quotes: ["まず探しやすさを改善する必要があります。"],
          },
        },
      ],
      artifact: undefined,
    };
    getKokugoUnit.mockResolvedValue(evidenceUnit);
    renderTab();
    fireEvent.click(await screen.findByRole("button", { name: /学校の図書室/ }));
    await screen.findByText("読む前の予測");
    fireEvent.click(screen.getByLabelText("工夫"));
    fireEvent.click(screen.getByRole("button", { name: /予測を記録/ }));
    await screen.findByText("本文を読む");
    fireEvent.click(screen.getByRole("button", { name: "課題へ進む" }));
    await screen.findByText("根拠を選ぶ");

    // Predict already submitted; only assert no *evidence* submit while disabled.
    submitKokugoTask.mockClear();
    const submitBtn = screen.getByRole("button", { name: "提出" });
    expect(submitBtn).toBeDisabled();
    fireEvent.click(submitBtn);
    expect(submitKokugoTask).not.toHaveBeenCalled();

    fireEvent.click(
      screen.getByRole("button", { name: "まず探しやすさを改善する必要があります。" })
    );
    expect(submitBtn).not.toBeDisabled();
    submitKokugoTask.mockResolvedValueOnce({
      attempt: {
        id: 3,
        unit_key: "e5-6/library-use",
        task_id: "evidence-1",
        answer: {},
        created_at: "",
      },
      grade: { correct: true, explanation_ja: "根拠となる文を正しく選べました。" },
    });
    fireEvent.click(submitBtn);
    await waitFor(() => {
      expect(submitKokugoTask).toHaveBeenCalledWith(
        "e5-6",
        "library-use",
        "evidence-1",
        { quotes: ["まず探しやすさを改善する必要があります。"] }
      );
    });
  });

  it("JS-133: consecutive same-kind tasks remount and reset selection state", async () => {
    /**
     * Behavior: key={task.id} resets EvidenceStep Map between consecutive evidence tasks.
     * 1. Resume on first evidence task; select a sentence (selection summary visible).
     * 2. Advance to second evidence task via successful submit mock.
     * 3. Assert selection summary is gone (state reset) and new prompt is shown.
     */
    const multiEvidence: KokugoUnit = {
      ...sampleUnit,
      text: [
        {
          kind: "paragraph",
          tokens: [{ t: "text", v: "第一の文です。第二の文です。" }],
        },
      ],
      tasks: [
        {
          id: "evidence-1",
          skill: "reading.locate-evidence",
          kind: "evidence-highlight",
          payload: {
            prompt_ja: "根拠Aを選びなさい。",
            gold_quotes: ["第一の文です。"],
          },
        },
        {
          id: "evidence-2",
          skill: "reading.locate-evidence",
          kind: "evidence-highlight",
          payload: {
            prompt_ja: "根拠Bを選びなさい。",
            gold_quotes: ["第二の文です。"],
          },
        },
      ],
      artifact: undefined,
    };
    getKokugoUnit.mockResolvedValue(multiEvidence);
    getKokugoUnitState.mockResolvedValue({
      progress: {
        unit_key: "e5-6/library-use",
        stage: "e5-6",
        unit_id: "library-use",
        status: "in_progress",
        step: "task:evidence-1",
        started_at: "",
        updated_at: "",
      },
      attempts: [],
      artifacts: [],
    });
    renderTab();
    fireEvent.click(await screen.findByRole("button", { name: /学校の図書室/ }));
    await screen.findByText("根拠Aを選びなさい。");
    fireEvent.click(screen.getByRole("button", { name: "第一の文です。" }));
    expect(screen.getByText(/選択中（1）/)).toBeVisible();
    submitKokugoTask.mockResolvedValueOnce({
      attempt: {
        id: 10,
        unit_key: "e5-6/library-use",
        task_id: "evidence-1",
        answer: {},
        created_at: "",
      },
      grade: { correct: true, explanation_ja: "ok" },
    });
    fireEvent.click(screen.getByRole("button", { name: "提出" }));
    await screen.findByText("根拠Bを選びなさい。");
    expect(screen.queryByText(/選択中/)).toBeNull();
    expect(screen.getByRole("button", { name: "提出" })).toBeDisabled();
  });

  it("JS-133: paragraph-role with zero paragraphs cannot submit", async () => {
    /**
     * Behavior: no paragraph blocks → submit disabled, no API call.
     * 1. Open unit with only callout text and a paragraph-role task.
     * 2. Assert status message and disabled 提出.
     * 3. Click 提出; submitKokugoTask not called.
     */
    const noPara: KokugoUnit = {
      ...sampleUnit,
      text: [{ kind: "callout", tokens: [{ t: "text", v: "段落なし" }] }],
      tasks: [
        {
          id: "structure-1",
          skill: "reading.structure",
          kind: "paragraph-role",
          payload: {
            prompt_ja: "役割を選びなさい。",
            roles: ["問題", "原因"],
            gold_by_paragraph_index: [],
          },
        },
      ],
      artifact: undefined,
    };
    getKokugoUnit.mockResolvedValue(noPara);
    getKokugoUnitState.mockResolvedValue({
      progress: {
        unit_key: "e5-6/library-use",
        stage: "e5-6",
        unit_id: "library-use",
        status: "in_progress",
        step: "task:structure-1",
        started_at: "",
        updated_at: "",
      },
      attempts: [],
      artifacts: [],
    });
    renderTab();
    fireEvent.click(await screen.findByRole("button", { name: /学校の図書室/ }));
    await screen.findByText(/段落がありません/);
    const submitBtn = screen.getByRole("button", { name: "提出" });
    expect(submitBtn).toBeDisabled();
    fireEvent.click(submitBtn);
    expect(submitKokugoTask).not.toHaveBeenCalled();
  });

  it("JS-133: second tap deselects evidence sentence before submit", async () => {
    /**
     * Behavior: toggle off removes quote from selection summary and payload path.
     * 1. Resume on evidence task with multi-sentence passage.
     * 2. Tap gold, assert 選択中; tap again to clear; 提出 disabled.
     * 3. Ensure submitKokugoTask never called while empty.
     */
    const evidenceUnit: KokugoUnit = {
      ...sampleUnit,
      text: [
        {
          kind: "paragraph",
          tokens: [
            {
              t: "text",
              v: "図書室は大切です。まず探しやすさを改善する必要があります。別の文です。",
            },
          ],
        },
      ],
      tasks: [
        {
          id: "evidence-1",
          skill: "reading.locate-evidence",
          kind: "evidence-highlight",
          payload: {
            prompt_ja: "根拠の文を選びなさい。",
            gold_quotes: ["まず探しやすさを改善する必要があります。"],
          },
        },
      ],
      artifact: undefined,
    };
    getKokugoUnit.mockResolvedValue(evidenceUnit);
    getKokugoUnitState.mockResolvedValue({
      progress: {
        unit_key: "e5-6/library-use",
        stage: "e5-6",
        unit_id: "library-use",
        status: "in_progress",
        step: "task:evidence-1",
        started_at: "",
        updated_at: "",
      },
      attempts: [],
      artifacts: [],
    });
    renderTab();
    fireEvent.click(await screen.findByRole("button", { name: /学校の図書室/ }));
    await screen.findByText("根拠を選ぶ");
    const gold = screen.getByRole("button", {
      name: "まず探しやすさを改善する必要があります。",
    });
    fireEvent.click(gold);
    expect(screen.getByText(/選択中（1）/)).toBeVisible();
    fireEvent.click(gold);
    expect(screen.queryByText(/選択中/)).toBeNull();
    expect(screen.getByRole("button", { name: "提出" })).toBeDisabled();
    fireEvent.click(screen.getByRole("button", { name: "提出" }));
    expect(submitKokugoTask).not.toHaveBeenCalled();
  });

  it("JS-133: paragraph-role submits roles assigned on the passage", async () => {
    /**
     * Behavior: initial roles[] length equals paragraph count; submit preserves order.
     * 1. Open unit with two paragraphs + list/callout noise + structure task.
     * 2. Assert two role selects; set 問題/原因; submit.
     * 3. Expect API body roles length 2 matching selections.
     */
    const roleUnit: KokugoUnit = {
      ...sampleUnit,
      text: [
        { kind: "callout", tokens: [{ t: "text", v: "ヒント" }] },
        { kind: "paragraph", tokens: [{ t: "text", v: "問題の段落。" }] },
        {
          kind: "list",
          items: [{ tokens: [{ t: "text", v: "メモ" }] }],
        },
        { kind: "paragraph", tokens: [{ t: "text", v: "原因の段落。" }] },
      ],
      tasks: [
        {
          id: "structure-1",
          skill: "reading.structure",
          kind: "paragraph-role",
          payload: {
            prompt_ja: "役割を選びなさい。",
            roles: ["問題", "原因", "提案", "結論"],
            gold_by_paragraph_index: ["問題", "原因"],
          },
        },
      ],
      artifact: undefined,
    };
    getKokugoUnit.mockResolvedValue(roleUnit);
    getKokugoUnitState.mockResolvedValue({
      progress: {
        unit_key: "e5-6/library-use",
        stage: "e5-6",
        unit_id: "library-use",
        status: "in_progress",
        step: "task:structure-1",
        started_at: "",
        updated_at: "",
      },
      attempts: [],
      artifacts: [],
    });
    renderTab();
    fireEvent.click(await screen.findByRole("button", { name: /学校の図書室/ }));
    await screen.findByText("段落の役割");
    expect(screen.getByText("ヒント")).toBeVisible();
    expect(screen.getAllByLabelText(/段落 \d+ の役割/)).toHaveLength(2);
    fireEvent.change(screen.getByLabelText("段落 1 の役割"), {
      target: { value: "問題" },
    });
    fireEvent.change(screen.getByLabelText("段落 2 の役割"), {
      target: { value: "原因" },
    });
    submitKokugoTask.mockResolvedValueOnce({
      attempt: {
        id: 4,
        unit_key: "e5-6/library-use",
        task_id: "structure-1",
        answer: {},
        created_at: "",
      },
      grade: { correct: true, explanation_ja: "各段落の役割を正しく整理できました。" },
    });
    fireEvent.click(screen.getByRole("button", { name: "提出" }));
    await waitFor(() => {
      expect(submitKokugoTask).toHaveBeenCalledWith(
        "e5-6",
        "library-use",
        "structure-1",
        { roles: ["問題", "原因"] }
      );
    });
  });

  it("retains phase when task submit fails", async () => {
    renderTab();
    fireEvent.click(await screen.findByRole("button", { name: /学校の図書室/ }));
    await screen.findByText("読む前の予測");
    fireEvent.click(screen.getByLabelText("工夫"));
    fireEvent.click(screen.getByRole("button", { name: /予測を記録/ }));
    await screen.findByText("本文を読む");
    fireEvent.click(screen.getByRole("button", { name: "課題へ進む" }));
    await screen.findByText("要約");
    submitKokugoTask.mockRejectedValueOnce(new Error("network down"));
    fireEvent.click(screen.getByLabelText("良い要約"));
    fireEvent.click(screen.getByRole("button", { name: "提出" }));
    await waitFor(() => {
      expect(screen.getByRole("alert")).toHaveTextContent("network down");
    });
    expect(screen.getByText("要約")).toBeVisible();
  });

  it("saves draft then revision with expected_version on second save", async () => {
    // Start already on artifact with no draft token
    getKokugoUnitState.mockResolvedValue({
      progress: {
        unit_key: "e5-6/library-use",
        stage: "e5-6",
        unit_id: "library-use",
        status: "in_progress",
        step: "artifact",
        started_at: "",
        updated_at: "",
      },
      attempts: [],
      artifacts: [],
    });
    saveKokugoArtifact
      .mockResolvedValueOnce({
        artifact: {
          unit_key: "e5-6/library-use",
          revision: 0,
          body: "下書き本文です。",
          checklist: [true, true],
          created_at: "",
          version: 1,
          updated_at: "ts-draft",
        },
        grade: {
          correct: true,
          explanation_ja: "下書きを保存しました。何度でも書き直してかまいません。",
        },
        progress: {
          unit_key: "e5-6/library-use",
          stage: "e5-6",
          unit_id: "library-use",
          status: "in_progress",
          step: "artifact",
          started_at: "",
          updated_at: "",
        },
      })
      .mockResolvedValueOnce({
        artifact: {
          unit_key: "e5-6/library-use",
          revision: 1,
          body: "改稿本文です。",
          checklist: [true, true],
          created_at: "",
          version: 1,
          updated_at: "ts-rev",
        },
        grade: { correct: true, explanation_ja: "ok" },
        progress: {
          unit_key: "e5-6/library-use",
          stage: "e5-6",
          unit_id: "library-use",
          status: "completed",
          step: "done",
          started_at: "",
          updated_at: "",
        },
      });

    renderTab();
    fireEvent.click(await screen.findByRole("button", { name: /学校の図書室/ }));
    await screen.findByText("作品（下書き）");
    fireEvent.change(screen.getByLabelText("下書き"), {
      target: { value: "下書き本文です。" },
    });
    // check both checklist boxes
    const boxes = screen.getAllByRole("checkbox");
    boxes.forEach((b) => fireEvent.click(b));
    fireEvent.click(screen.getByRole("button", { name: "下書きを保存" }));
    await waitFor(() => {
      expect(saveKokugoArtifact).toHaveBeenCalledWith(
        "e5-6",
        "library-use",
        expect.objectContaining({ revision: 0, body: "下書き本文です。" })
      );
    });
    // Progressive writing: stay on draft after save; advance only via 改稿へ.
    expect(screen.getByText("作品（下書き）")).toBeVisible();
    fireEvent.click(screen.getByRole("button", { name: "改稿へ進む" }));
    await waitFor(() => {
      expect(screen.getByText("改稿")).toBeVisible();
    });
    expect(putKokugoProgress).toHaveBeenCalledWith("e5-6", "library-use", {
      step: "revise",
    });
    fireEvent.click(screen.getByRole("button", { name: /改稿を保存/ }));
    await waitFor(() => {
      expect(saveKokugoArtifact).toHaveBeenCalledWith(
        "e5-6",
        "library-use",
        expect.objectContaining({ revision: 1 })
      );
    });
    await waitFor(() => {
      expect(screen.getByText("循環完了")).toBeVisible();
    });
  });

  it("does not show done when server progress stays in_progress after revision", async () => {
    // Missing task attempts: server accepts artifact grades but does not complete.
    getKokugoUnitState.mockResolvedValue({
      progress: {
        unit_key: "e5-6/library-use",
        stage: "e5-6",
        unit_id: "library-use",
        status: "in_progress",
        step: "artifact",
        started_at: "",
        updated_at: "",
      },
      attempts: [],
      artifacts: [],
    });
    saveKokugoArtifact
      .mockResolvedValueOnce({
        artifact: {
          unit_key: "e5-6/library-use",
          revision: 0,
          body: "下書き本文です。",
          checklist: [true, true],
          created_at: "",
          version: 1,
          updated_at: "ts-draft",
        },
        grade: { correct: true, explanation_ja: "ok" },
        progress: {
          unit_key: "e5-6/library-use",
          stage: "e5-6",
          unit_id: "library-use",
          status: "in_progress",
          step: "artifact",
          started_at: "",
          updated_at: "",
        },
      })
      .mockResolvedValueOnce({
        artifact: {
          unit_key: "e5-6/library-use",
          revision: 1,
          body: "改稿本文です。",
          checklist: [true, true],
          created_at: "",
          version: 1,
          updated_at: "ts-rev",
        },
        grade: { correct: true, explanation_ja: "ok" },
        progress: {
          unit_key: "e5-6/library-use",
          stage: "e5-6",
          unit_id: "library-use",
          status: "in_progress",
          step: "revise",
          started_at: "",
          updated_at: "",
        },
      });

    renderTab();
    fireEvent.click(await screen.findByRole("button", { name: /学校の図書室/ }));
    await screen.findByText("作品（下書き）");
    fireEvent.change(screen.getByLabelText("下書き"), {
      target: { value: "下書き本文です。" },
    });
    screen.getAllByRole("checkbox").forEach((b) => fireEvent.click(b));
    fireEvent.click(screen.getByRole("button", { name: "下書きを保存" }));
    await waitFor(() => {
      expect(saveKokugoArtifact).toHaveBeenCalledTimes(1);
    });
    expect(screen.getByText("作品（下書き）")).toBeVisible();
    fireEvent.click(screen.getByRole("button", { name: "改稿へ進む" }));
    await waitFor(() => {
      expect(screen.getByText("改稿")).toBeVisible();
    });
    fireEvent.click(screen.getByRole("button", { name: /改稿を保存/ }));
    await waitFor(() => {
      expect(saveKokugoArtifact).toHaveBeenCalledTimes(2);
    });
    expect(screen.queryByText("循環完了")).not.toBeInTheDocument();
    expect(screen.getByText("改稿")).toBeVisible();
  });

  it("draft can be re-edited after save; short text is allowed", async () => {
    /**
     * Behavior: short draft saves without char floor; stays on draft for re-edit.
     * 1. Open unit on artifact phase.
     * 2. Save short body; stay on 下書き.
     * 3. Edit again and save second time with CAS version; then 改稿へ / 下書きに戻る.
     */
    getKokugoUnitState.mockResolvedValue({
      progress: {
        unit_key: "e5-6/library-use",
        stage: "e5-6",
        unit_id: "library-use",
        status: "in_progress",
        step: "artifact",
        started_at: "",
        updated_at: "",
      },
      attempts: [],
      artifacts: [],
    });
    saveKokugoArtifact
      .mockResolvedValueOnce({
        artifact: {
          unit_key: "e5-6/library-use",
          revision: 0,
          body: "短い提案。",
          checklist: [false, false],
          created_at: "",
          version: 1,
          updated_at: "t1",
        },
        grade: {
          correct: true,
          explanation_ja: "下書きを保存しました。何度でも書き直してかまいません。",
        },
        progress: {
          unit_key: "e5-6/library-use",
          stage: "e5-6",
          unit_id: "library-use",
          status: "in_progress",
          step: "artifact",
          started_at: "",
          updated_at: "",
        },
      })
      .mockResolvedValueOnce({
        artifact: {
          unit_key: "e5-6/library-use",
          revision: 0,
          body: "短い提案。理由も足した。",
          checklist: [true, false],
          created_at: "",
          version: 2,
          updated_at: "t2",
        },
        grade: { correct: true, explanation_ja: "下書きを保存しました。" },
        progress: {
          unit_key: "e5-6/library-use",
          stage: "e5-6",
          unit_id: "library-use",
          status: "in_progress",
          step: "artifact",
          started_at: "",
          updated_at: "",
        },
      });

    renderTab();
    fireEvent.click(await screen.findByRole("button", { name: /学校の図書室/ }));
    await screen.findByText("作品（下書き）");
    fireEvent.change(screen.getByLabelText("下書き"), {
      target: { value: "短い提案。" },
    });
    fireEvent.click(screen.getByRole("button", { name: "下書きを保存" }));
    await waitFor(() => {
      expect(screen.getByText(/何度でも書き直して/)).toBeVisible();
    });
    expect(screen.getByText("作品（下書き）")).toBeVisible();
    fireEvent.change(screen.getByLabelText("下書き"), {
      target: { value: "短い提案。理由も足した。" },
    });
    fireEvent.click(screen.getByRole("button", { name: "下書きを保存" }));
    await waitFor(() => {
      expect(saveKokugoArtifact).toHaveBeenLastCalledWith(
        "e5-6",
        "library-use",
        expect.objectContaining({
          revision: 0,
          body: "短い提案。理由も足した。",
          expected_version: 1,
        })
      );
    });
    fireEvent.click(screen.getByRole("button", { name: "改稿へ進む" }));
    await screen.findByText("改稿");
    fireEvent.click(screen.getByRole("button", { name: "下書きに戻る" }));
    await screen.findByText("作品（下書き）");
  });

  it("progress-disabled fallback advances without calling submit", async () => {
    renderTab({ progress: false, kokugo: true });
    fireEvent.click(await screen.findByRole("button", { name: /学校の図書室/ }));
    await screen.findByText("読む前の予測");
    fireEvent.click(screen.getByLabelText("工夫"));
    fireEvent.click(screen.getByRole("button", { name: /予測を記録/ }));
    await waitFor(() => {
      expect(screen.getByText("本文を読む")).toBeVisible();
    });
    expect(submitKokugoTask).not.toHaveBeenCalled();
    expect(getKokugoUnitState).not.toHaveBeenCalled();
    expect(
      screen.getByText(/ローカル API が無いため/)
    ).toBeVisible();
  });
});
