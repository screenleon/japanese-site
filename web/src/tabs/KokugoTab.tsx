import { useCallback, useEffect, useMemo, useState } from "react";
import { api } from "../api";
import type {
  KokugoGradeResult,
  KokugoUnitState,
  KokugoUnitSummary,
} from "../apiTypes";
import { ClassmatePanel } from "../components/ClassmatePanel";
import {
  countParagraphs,
  KokugoPassage,
  type PassageSentence,
} from "../components/KokugoPassage";
import { RevisionCompare } from "../components/RevisionCompare";
import type {
  EvidenceHighlightPayload,
  KokugoTask,
  KokugoUnit,
  ParagraphRolePayload,
  PredictPayload,
  SummaryChoicePayload,
} from "../kokugoTypes";
import {
  classmatesForArtifact,
  classmatesForRevise,
  classmatesForTask,
} from "../kokugoTypes";
import { useCapabilities } from "../capabilities";

export type Phase =
  | "list"
  | "predict"
  | "read"
  | "task"
  | "artifact"
  | "revise"
  | "done";

/** Map server step/status + artifacts into UI phase and task index. */
export function deriveResumeState(
  unit: KokugoUnit,
  state: KokugoUnitState | null | undefined
): {
  phase: Phase;
  taskIndex: number;
  draftBody: string;
  revisionBody: string;
  checklist: boolean[];
  draftVersion: number;
  revisionVersion: number;
} {
  const checklistLen = unit.artifact?.checklist.length ?? 0;
  const emptyChecklist = Array.from({ length: checklistLen }, () => false);
  const draft = state?.artifacts?.find((a) => a.revision === 0);
  const rev = state?.artifacts?.find((a) => a.revision === 1);
  let checklist = emptyChecklist;
  const checklistSrc = rev?.checklist ?? draft?.checklist;
  if (Array.isArray(checklistSrc)) {
    checklist = emptyChecklist.map((_, i) => Boolean(checklistSrc[i]));
  }

  const draftBody = draft?.body ?? "";
  const revisionBody = rev?.body ?? "";
  const draftVersion = draft?.version ?? 0;
  const revisionVersion = rev?.version ?? 0;
  const base = {
    draftBody,
    revisionBody,
    checklist,
    draftVersion,
    revisionVersion,
  };

  const progress = state?.progress;
  if (!progress) {
    return { phase: "predict", taskIndex: 0, ...base };
  }
  if (progress.status === "completed" || progress.step === "done") {
    return {
      phase: "done",
      taskIndex: Math.max(0, unit.tasks.length - 1),
      ...base,
    };
  }

  const step = progress.step || "predict";
  switch (step) {
    case "predict":
      return { phase: "predict", taskIndex: 0, ...base };
    case "read":
      return { phase: "read", taskIndex: 0, ...base };
    case "revise":
      return {
        phase: "revise",
        taskIndex: firstNonPredictIndex(unit),
        ...base,
        revisionBody: revisionBody || draftBody,
      };
    case "artifact":
      // Stay on draft phase (may already have a body if grade failed and learner re-edits).
      return { phase: "artifact", taskIndex: firstNonPredictIndex(unit), ...base };
    default:
      if (step === "task" || step.startsWith("task:")) {
        return {
          phase: "task",
          taskIndex: taskIndexFromState(unit, state, step),
          ...base,
        };
      }
      // Unknown step: infer from artifacts/attempts.
      if (rev) {
        return { phase: "done", taskIndex: firstNonPredictIndex(unit), ...base };
      }
      if (draft) {
        return {
          phase: "revise",
          taskIndex: firstNonPredictIndex(unit),
          ...base,
          revisionBody: revisionBody || draftBody,
        };
      }
      return {
        phase: "task",
        taskIndex: taskIndexFromState(unit, state, step),
        ...base,
      };
  }
}

function firstNonPredictIndex(unit: KokugoUnit): number {
  const i = unit.tasks.findIndex((t) => t.kind !== "predict");
  return i >= 0 ? i : 0;
}

function taskIndexFromState(
  unit: KokugoUnit,
  state: KokugoUnitState | null | undefined,
  step: string
): number {
  if (step.startsWith("task:")) {
    const id = step.slice("task:".length);
    const idx = unit.tasks.findIndex((t) => t.id === id);
    if (idx >= 0) return idx;
  }
  const attempts = state?.attempts ?? [];
  for (let i = 0; i < unit.tasks.length; i++) {
    const t = unit.tasks[i];
    if (t.kind === "predict") continue;
    if (!attempts.some((a) => a.task_id === t.id)) return i;
  }
  return firstNonPredictIndex(unit);
}

export function KokugoTab() {
  const { progress, kokugo, loaded: capsLoaded } = useCapabilities();
  const [units, setUnits] = useState<KokugoUnitSummary[]>([]);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(true);
  const [unit, setUnit] = useState<KokugoUnit | null>(null);
  const [phase, setPhase] = useState<Phase>("list");
  const [taskIndex, setTaskIndex] = useState(0);
  const [feedback, setFeedback] = useState<KokugoGradeResult | null>(null);
  const [draftBody, setDraftBody] = useState("");
  const [revisionBody, setRevisionBody] = useState("");
  const [checklist, setChecklist] = useState<boolean[]>([]);
  const [draftVersion, setDraftVersion] = useState(0);
  const [revisionVersion, setRevisionVersion] = useState(0);
  /** Last task the learner submitted — drives classmate reveal (JS-134). */
  const [lastSubmittedTaskId, setLastSubmittedTaskId] = useState<string | null>(null);
  /** Draft was saved at least once this open (or restored with version > 0). */
  const [draftSavedOnce, setDraftSavedOnce] = useState(false);

  const loadUnits = useCallback(async () => {
    setLoading(true);
    setError("");
    try {
      const res = await api.listKokugoUnits();
      setUnits(res.units);
    } catch (e) {
      setError(String(e));
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    void loadUnits();
  }, [loadUnits]);

  async function openUnit(summary: KokugoUnitSummary) {
    setError("");
    setFeedback(null);
    try {
      const u = await api.getKokugoUnit(summary.stage, summary.id);
      setUnit(u);

      let state: KokugoUnitState | null = null;
      if (progress) {
        try {
          state = await api.getKokugoUnitState(u.stage, u.id);
        } catch {
          state = null;
        }
      }

      const resumed = deriveResumeState(u, state);
      setTaskIndex(resumed.taskIndex);
      setDraftBody(resumed.draftBody);
      setRevisionBody(resumed.revisionBody);
      setChecklist(resumed.checklist);
      setDraftVersion(resumed.draftVersion);
      setRevisionVersion(resumed.revisionVersion);
      setPhase(resumed.phase);
      setDraftSavedOnce(resumed.draftVersion > 0);
      // Resume: show classmates for the most recent attempted task (if any).
      const attempts = state?.attempts ?? [];
      setLastSubmittedTaskId(attempts.length > 0 ? attempts[attempts.length - 1].task_id : null);

      // Only seed progress when starting fresh (no prior row).
      if (progress && !state?.progress) {
        await api.putKokugoProgress(u.stage, u.id, { step: "predict", status: "in_progress" });
      }
    } catch (e) {
      setError(String(e));
    }
  }

  function backToList() {
    setUnit(null);
    setPhase("list");
    setFeedback(null);
    setDraftVersion(0);
    setRevisionVersion(0);
    setLastSubmittedTaskId(null);
    setDraftSavedOnce(false);
    void loadUnits();
  }

  const tasks = unit?.tasks ?? [];
  const currentTask = tasks[taskIndex];

  async function submitTask(answer: unknown, after: "read" | "next-task") {
    if (!unit || !currentTask) return;
    setError("");
    setFeedback(null);
    if (!progress) {
      setFeedback({
        correct: null,
        explanation_ja: "ローカル API が無いため採点を保存できません。内容の確認だけ進めます。",
      });
      setLastSubmittedTaskId(currentTask.id);
      advanceAfterTask(after);
      return;
    }
    try {
      const res = await api.submitKokugoTask(unit.stage, unit.id, currentTask.id, answer);
      setFeedback(res.grade);
      setLastSubmittedTaskId(currentTask.id);
      advanceAfterTask(after);
    } catch (e) {
      setError(String(e));
      // Retain phase on save/submit error (qa-tester-F002).
    }
  }

  function advanceAfterTask(after: "read" | "next-task") {
    if (!unit) return;
    if (after === "read") {
      setPhase("read");
      if (progress) {
        void api.putKokugoProgress(unit.stage, unit.id, { step: "read" });
      }
      return;
    }
    let next = taskIndex + 1;
    while (next < tasks.length && tasks[next]?.kind === "predict") next++;
    if (next < tasks.length) {
      setTaskIndex(next);
      setPhase("task");
      if (progress) {
        void api.putKokugoProgress(unit.stage, unit.id, {
          step: `task:${tasks[next].id}`,
        });
      }
      return;
    }
    if (unit.artifact) {
      setPhase("artifact");
      if (progress) {
        void api.putKokugoProgress(unit.stage, unit.id, { step: "artifact" });
      }
    } else {
      setPhase("done");
      if (progress) {
        void api.putKokugoProgress(unit.stage, unit.id, { step: "done", status: "completed" });
      }
    }
  }

  async function saveArtifact(revision: 0 | 1) {
    if (!unit) return;
    const body = revision === 0 ? draftBody : revisionBody;
    const checks = checklist;
    setError("");
    if (!progress) {
      setFeedback({
        correct: null,
        explanation_ja: "静的モードでは作品を保存できません。",
      });
      return;
    }
    const expected = revision === 0 ? draftVersion : revisionVersion;
    try {
      const res = await api.saveKokugoArtifact(unit.stage, unit.id, {
        revision,
        body,
        checklist_checked: checks,
        ...(expected > 0 ? { expected_version: expected } : {}),
      });
      setFeedback(res.grade);
      if (revision === 0) {
        setDraftVersion(res.artifact.version);
        setDraftBody(res.artifact.body);
        setDraftSavedOnce(true);
        // Progressive writing: stay on draft after save — learner re-edits freely.
        // Advance to revise only via explicit 「改稿へ」.
        // Until rev1 is persisted, keep the revision editor seed in sync with latest draft.
        if (res.grade.correct && revisionVersion === 0) {
          setRevisionBody(res.artifact.body);
        }
      } else {
        setRevisionVersion(res.artifact.version);
        setRevisionBody(res.artifact.body);
        // Done requires server-confirmed completed progress (critic-F001).
        if (res.grade.correct) {
          const serverDone =
            res.progress?.status === "completed" || res.progress?.step === "done";
          if (serverDone) {
            setPhase("done");
          }
        }
      }
    } catch (e) {
      setError(String(e));
    }
  }

  function goToRevise() {
    if (!unit) return;
    setFeedback(null);
    // Prefer latest draft when rev1 has never been saved; keep learner rev1 text otherwise.
    if (revisionVersion === 0) {
      setRevisionBody(draftBody);
    }
    setPhase("revise");
    if (progress) {
      void api.putKokugoProgress(unit.stage, unit.id, { step: "revise" });
    }
  }

  function goToDraft() {
    if (!unit) return;
    setFeedback(null);
    setPhase("artifact");
    if (progress) {
      void api.putKokugoProgress(unit.stage, unit.id, { step: "artifact" });
    }
  }

  const draftCharCount = [...draftBody.trim()].length;
  const revisionCharCount = [...revisionBody.trim()].length;
  const activeCharCount = phase === "revise" ? revisionCharCount : draftCharCount;
  // Server-persisted draft only — local typing alone must not unlock 改稿へ
  // (rev1 save requires draft_required on the server).
  const draftSaved = draftVersion > 0;

  if (capsLoaded && !kokugo) {
    return (
      <section className="p-4 space-y-2">
        <h2 className="text-lg font-semibold">国語教室</h2>
        <p className="text-sm text-slate-600">
          国語ユニットはこの環境では利用できません（静的ビルド、またはコーパス未設定）。
        </p>
      </section>
    );
  }

  if (phase === "list" || !unit) {
    return (
      <section className="p-4 space-y-4" aria-label="国語ユニット一覧">
        <header>
          <h2 className="text-lg font-semibold">国語教室</h2>
          <p className="text-sm text-slate-600 mt-1">
            予測 → 本文 → 課題 → 作品 → 改稿の最小循環（e5-6）
          </p>
        </header>
        {error && <p className="text-sm text-red-600">{error}</p>}
        {!capsLoaded || loading ? (
          <p className="text-sm text-slate-500">読み込み中…</p>
        ) : units.length === 0 ? (
          <p className="text-sm text-slate-500">ユニットがありません。</p>
        ) : (
          <ul className="space-y-2">
            {units.map((u) => (
              <li key={`${u.stage}/${u.id}`}>
                <button
                  type="button"
                  onClick={() => void openUnit(u)}
                  className="w-full text-left p-4 rounded-lg border border-slate-200 bg-white hover:border-blue-400 hover:bg-blue-50/40 transition-colors"
                >
                  <div className="flex items-start justify-between gap-3">
                    <div>
                      <p className="font-medium text-slate-900">{u.title_ja}</p>
                      <p className="text-xs text-slate-500 mt-1">
                        {u.stage} · {u.genre} · 約 {u.estimated_minutes} 分 · 課題 {u.task_count}
                      </p>
                    </div>
                    <span className="text-xs text-blue-700 shrink-0">開く</span>
                  </div>
                </button>
              </li>
            ))}
          </ul>
        )}
      </section>
    );
  }

  return (
    <section className="p-4 space-y-4" aria-label="国語ユニット">
      <div className="flex items-center justify-between gap-3">
        <button type="button" onClick={backToList} className="text-sm text-blue-700 hover:underline">
          ← 一覧
        </button>
        <PhaseBadge phase={phase} taskIndex={taskIndex} taskCount={tasks.length} />
      </div>
      <header>
        <h2 className="text-lg font-semibold">{unit.title_ja}</h2>
        <p className="text-xs text-slate-500 mt-1">
          {unit.stage} · {unit.genre}
        </p>
      </header>
      {error && (
        <p className="text-sm text-red-600" role="alert">
          {error}
        </p>
      )}
      {feedback && (
        <div
          className={
            "text-sm rounded-md border px-3 py-2 " +
            (feedback.correct === false
              ? "border-amber-200 bg-amber-50 text-amber-900"
              : "border-slate-200 bg-slate-50 text-slate-800")
          }
        >
          {feedback.explanation_ja}
        </div>
      )}

      {lastSubmittedTaskId &&
        (phase === "task" || phase === "read" || phase === "artifact") && (
          <ClassmatePanel classmates={classmatesForTask(unit, lastSubmittedTaskId)} />
        )}

      {phase === "predict" && currentTask?.kind === "predict" && (
        <PredictStep
          payload={currentTask.payload}
          onSubmit={(choiceId) => {
            void submitTask({ choice_id: choiceId }, "read");
          }}
        />
      )}

      {phase === "read" && (
        <div className="space-y-4">
          <h3 className="text-sm font-medium text-slate-800">本文を読む</h3>
          <p className="text-xs text-slate-600">
            段落ごとに読み、主張・理由・提案がどこにあるか意識しましょう。
          </p>
          <KokugoPassage blocks={unit.text} mode="readonly" showParagraphIndex />
          <button
            type="button"
            className="px-4 py-2 rounded-md bg-blue-600 text-white text-sm font-medium hover:bg-blue-700"
            onClick={() => {
              setFeedback(null);
              const first = tasks.findIndex((t) => t.kind !== "predict");
              setTaskIndex(first >= 0 ? first : 0);
              setPhase("task");
              if (progress) {
                const tid = tasks[first >= 0 ? first : 0]?.id;
                void api.putKokugoProgress(unit.stage, unit.id, {
                  step: tid ? `task:${tid}` : "task",
                });
              }
            }}
          >
            課題へ進む
          </button>
        </div>
      )}

      {phase === "task" && currentTask && (
        <TaskStep
          key={currentTask.id}
          task={currentTask}
          unit={unit}
          onSubmit={(answer) => void submitTask(answer, "next-task")}
        />
      )}

      {(phase === "artifact" || phase === "revise") && unit.artifact && (
        <div className="space-y-3">
          <h3 className="text-sm font-medium text-slate-800">
            {phase === "artifact" ? "作品（下書き）" : "改稿"}
          </h3>
          <p className="text-xs text-slate-600">
            {unit.artifact.kind}
            {" · "}
            現在 {activeCharCount} 字
            {unit.artifact.min_chars > 0 || unit.artifact.max_chars > 0
              ? ` · 目安 ${unit.artifact.min_chars > 0 ? unit.artifact.min_chars : "—"}〜${
                  unit.artifact.max_chars > 0 ? unit.artifact.max_chars : "—"
                } 字`
              : " · 字数制限なし（短く書いてから少しずつ増やせます）"}
          </p>
          {phase === "artifact" && (
            <p className="text-xs text-slate-500">
              下書きは何度でも保存・書き直せます。準備ができたら「改稿へ」へ進んでください。
            </p>
          )}
          {phase === "revise" && (
            <p className="text-xs text-slate-500">
              改稿ではチェックリストを確認してから保存します。必要なら下書きに戻って直せます。
            </p>
          )}
          <ul className="space-y-1">
            {unit.artifact.checklist.map((item, i) => (
              <li key={item} className="flex items-start gap-2 text-sm">
                <input
                  type="checkbox"
                  checked={checklist[i] ?? false}
                  onChange={(e) => {
                    setChecklist((prev) => {
                      const next = [...prev];
                      next[i] = e.target.checked;
                      return next;
                    });
                  }}
                />
                <span>
                  {item}
                  {phase === "artifact" ? (
                    <span className="text-xs text-slate-400">（下書きでは任意）</span>
                  ) : null}
                </span>
              </li>
            ))}
          </ul>
          <textarea
            className="w-full min-h-[8rem] rounded-md border border-slate-300 p-3 text-sm"
            value={phase === "artifact" ? draftBody : revisionBody}
            onChange={(e) =>
              phase === "artifact" ? setDraftBody(e.target.value) : setRevisionBody(e.target.value)
            }
            placeholder="ここに提案文を書いてください（短くて大丈夫です）"
            aria-label={phase === "artifact" ? "下書き" : "改稿本文"}
          />
          {unit.artifact.exemplar_ja && (
            <details className="text-sm text-slate-600">
              <summary className="cursor-pointer">例（参考）</summary>
              <p className="mt-2 p-2 bg-slate-50 rounded border border-slate-100">
                {unit.artifact.exemplar_ja}
              </p>
            </details>
          )}
          {phase === "artifact" && draftSavedOnce && (
            <ClassmatePanel
              classmates={classmatesForArtifact(unit)}
              title="クラスメイトの下書き"
            />
          )}
          {phase === "revise" && (
            <>
              <RevisionCompare draftBody={draftBody} revisionBody={revisionBody} />
              <ClassmatePanel
                classmates={classmatesForRevise(unit)}
                title="クラスメイトの改稿例"
              />
            </>
          )}
          <div className="flex flex-wrap gap-2">
            {phase === "artifact" ? (
              <>
                <button
                  type="button"
                  className="px-4 py-2 rounded-md bg-blue-600 text-white text-sm font-medium hover:bg-blue-700 disabled:opacity-40"
                  disabled={!draftBody.trim()}
                  onClick={() => void saveArtifact(0)}
                >
                  下書きを保存
                </button>
                <button
                  type="button"
                  className="px-4 py-2 rounded-md border border-slate-300 bg-white text-sm font-medium text-slate-800 hover:bg-slate-50 disabled:opacity-40"
                  disabled={!draftSaved || !draftBody.trim()}
                  onClick={goToRevise}
                >
                  改稿へ進む
                </button>
              </>
            ) : (
              <>
                <button
                  type="button"
                  className="px-4 py-2 rounded-md bg-blue-600 text-white text-sm font-medium hover:bg-blue-700 disabled:opacity-40"
                  disabled={!revisionBody.trim()}
                  onClick={() => void saveArtifact(1)}
                >
                  改稿を保存して完了
                </button>
                <button
                  type="button"
                  className="px-4 py-2 rounded-md border border-slate-300 bg-white text-sm font-medium text-slate-800 hover:bg-slate-50"
                  onClick={goToDraft}
                >
                  下書きに戻る
                </button>
              </>
            )}
          </div>
        </div>
      )}

      {phase === "done" && (
        <div className="space-y-3 rounded-lg border border-green-200 bg-green-50 p-4">
          <h3 className="font-medium text-green-900">循環完了</h3>
          <p className="text-sm text-green-800">
            予測・読解・課題・作品・改稿の一巡が終わりました。
          </p>
          <RevisionCompare
            draftBody={draftBody}
            revisionBody={revisionBody}
            title="完成した作品の対比"
          />
          <ClassmatePanel
            classmates={[
              ...classmatesForArtifact(unit),
              ...classmatesForRevise(unit),
            ]}
            title="クラスメイトの作品例"
          />
          <button
            type="button"
            className="px-4 py-2 rounded-md bg-white border border-green-300 text-sm"
            onClick={backToList}
          >
            一覧に戻る
          </button>
        </div>
      )}
    </section>
  );
}

function PhaseBadge({
  phase,
  taskIndex,
  taskCount,
}: {
  phase: Phase;
  taskIndex: number;
  taskCount: number;
}) {
  const label = useMemo(() => {
    switch (phase) {
      case "predict":
        return "1 予測";
      case "read":
        return "2 本文";
      case "task":
        return `3 課題 ${taskIndex + 1}/${taskCount}`;
      case "artifact":
        return "4 作品";
      case "revise":
        return "5 改稿";
      case "done":
        return "完了";
      default:
        return "";
    }
  }, [phase, taskIndex, taskCount]);
  return (
    <span className="text-xs font-medium px-2 py-1 rounded-full bg-slate-100 text-slate-700">
      {label}
    </span>
  );
}

function PredictStep({
  payload,
  onSubmit,
}: {
  payload: PredictPayload;
  onSubmit: (choiceId: string) => void;
}) {
  const [choice, setChoice] = useState<string>("");
  return (
    <div className="space-y-3">
      <h3 className="text-sm font-medium">読む前の予測</h3>
      <p className="text-sm text-slate-800">{payload.prompt_ja}</p>
      <ul className="space-y-2">
        {payload.choices.map((c) => (
          <li key={c.id}>
            <label className="flex items-start gap-2 text-sm cursor-pointer">
              <input
                type="radio"
                name="predict"
                checked={choice === c.id}
                onChange={() => setChoice(c.id)}
              />
              <span>{c.text_ja}</span>
            </label>
          </li>
        ))}
      </ul>
      <button
        type="button"
        disabled={!choice}
        className="px-4 py-2 rounded-md bg-blue-600 text-white text-sm font-medium disabled:opacity-40"
        onClick={() => onSubmit(choice)}
      >
        予測を記録して本文へ
      </button>
    </div>
  );
}

function TaskStep({
  task,
  unit,
  onSubmit,
}: {
  task: KokugoTask;
  unit: KokugoUnit;
  onSubmit: (answer: unknown) => void;
}) {
  if (task.kind === "predict") {
    return (
      <PredictStep
        key={task.id}
        payload={task.payload}
        onSubmit={(choiceId) => onSubmit({ choice_id: choiceId })}
      />
    );
  }
  if (task.kind === "summary-choice") {
    return <SummaryChoiceStep key={task.id} payload={task.payload} onSubmit={onSubmit} />;
  }
  if (task.kind === "paragraph-role") {
    return (
      <ParagraphRoleStep key={task.id} payload={task.payload} unit={unit} onSubmit={onSubmit} />
    );
  }
  if (task.kind === "evidence-highlight") {
    return <EvidenceStep key={task.id} payload={task.payload} unit={unit} onSubmit={onSubmit} />;
  }
  return <p className="text-sm text-red-600">未対応の課題種別です。</p>;
}

function SummaryChoiceStep({
  payload,
  onSubmit,
}: {
  payload: SummaryChoicePayload;
  onSubmit: (answer: unknown) => void;
}) {
  const [choice, setChoice] = useState("");
  return (
    <div className="space-y-3">
      <h3 className="text-sm font-medium">要約</h3>
      <p className="text-sm">{payload.prompt_ja}</p>
      <ul className="space-y-2">
        {payload.choices.map((c) => (
          <li key={c.id}>
            <label className="flex items-start gap-2 text-sm cursor-pointer">
              <input
                type="radio"
                name="summary"
                checked={choice === c.id}
                onChange={() => setChoice(c.id)}
              />
              <span>{c.text_ja}</span>
            </label>
          </li>
        ))}
      </ul>
      <button
        type="button"
        disabled={!choice}
        className="px-4 py-2 rounded-md bg-blue-600 text-white text-sm disabled:opacity-40"
        onClick={() => onSubmit({ choice_id: choice })}
      >
        提出
      </button>
    </div>
  );
}

function ParagraphRoleStep({
  payload,
  unit,
  onSubmit,
}: {
  payload: ParagraphRolePayload;
  unit: KokugoUnit;
  onSubmit: (answer: unknown) => void;
}) {
  const paraCount = countParagraphs(unit.text);
  const [roles, setRoles] = useState<string[]>(() =>
    Array.from({ length: paraCount }, () => payload.roles[0] ?? "")
  );

  return (
    <div className="space-y-3">
      <h3 className="text-sm font-medium">段落の役割</h3>
      <p className="text-sm text-slate-800">{payload.prompt_ja}</p>
      {paraCount === 0 ? (
        <p className="text-sm text-amber-800 rounded-md border border-amber-200 bg-amber-50 px-3 py-2" role="status">
          このユニットに段落がありません。役割課題を提出できません。
        </p>
      ) : (
        <p className="text-xs text-slate-600">
          各段落の見出しで役割を選びます。色が変わると割り当てが反映されます。
        </p>
      )}
      <KokugoPassage
        blocks={unit.text}
        mode="paragraph-role"
        roles={roles}
        roleOptions={payload.roles}
        onRoleChange={(i, role) => {
          setRoles((prev) => {
            const next = [...prev];
            next[i] = role;
            return next;
          });
        }}
        showParagraphIndex
      />
      <button
        type="button"
        disabled={paraCount === 0 || roles.length !== paraCount}
        className="px-4 py-2 rounded-md bg-blue-600 text-white text-sm disabled:opacity-40"
        onClick={() => {
          if (paraCount === 0 || roles.length !== paraCount) return;
          onSubmit({ roles });
        }}
      >
        提出
      </button>
    </div>
  );
}

function EvidenceStep({
  payload,
  unit,
  onSubmit,
}: {
  payload: EvidenceHighlightPayload;
  unit: KokugoUnit;
  onSubmit: (answer: unknown) => void;
}) {
  // key → surface quote text (for API quotes[])
  const [selected, setSelected] = useState<Map<string, string>>(() => new Map());

  function toggle(sent: PassageSentence) {
    setSelected((prev) => {
      const next = new Map(prev);
      if (next.has(sent.key)) next.delete(sent.key);
      else next.set(sent.key, sent.text);
      return next;
    });
  }

  const selectedKeys = useMemo(() => new Set(selected.keys()), [selected]);
  const selectedEntries = useMemo(() => Array.from(selected.entries()), [selected]);
  const quotes = useMemo(() => selectedEntries.map(([, text]) => text), [selectedEntries]);

  return (
    <div className="space-y-3">
      <h3 className="text-sm font-medium">根拠を選ぶ</h3>
      <p className="text-sm text-slate-800">{payload.prompt_ja}</p>
      <p className="text-xs text-slate-600">
        本文の文をタップして根拠をマークします（複数可）。もう一度タップすると解除します。
      </p>
      <KokugoPassage
        blocks={unit.text}
        mode="sentence-select"
        selectedKeys={selectedKeys}
        onToggleSentence={toggle}
        goldQuotes={payload.gold_quotes}
        showParagraphIndex
      />
      {selectedEntries.length > 0 && (
        <div className="rounded-md border border-amber-200 bg-amber-50/60 px-3 py-2 text-xs text-amber-950">
          <p className="font-medium mb-1">選択中（{selectedEntries.length}）</p>
          <ul className="list-disc pl-4 space-y-0.5">
            {selectedEntries.map(([key, text]) => (
              <li key={key}>{text}</li>
            ))}
          </ul>
        </div>
      )}
      <button
        type="button"
        disabled={quotes.length === 0}
        className="px-4 py-2 rounded-md bg-blue-600 text-white text-sm disabled:opacity-40"
        onClick={() => onSubmit({ quotes })}
      >
        提出
      </button>
    </div>
  );
}
