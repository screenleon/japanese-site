import { useEffect, useRef, useState } from "react";
import {
  api,
  isApiError,
  type Question,
  type GradeResult,
  type GrammarPoint,
  type QuizContentType,
  type Stats,
} from "../api";

type Mode = "single" | "session" | "summary";
export type QuizInitialMode = "練習" | "測試";
type StatsRange = 7 | 30 | 0;
const EXAM_DURATION_SEC = 600;

interface SessionTurn {
  question: Question;
  result: GradeResult | null;
}

export function QuizTab({
  initialMode = "練習",
  onNavigateGrammar,
}: {
  initialMode?: QuizInitialMode;
  onNavigateGrammar?: (slug: string) => void;
}) {
  const sealedExam = initialMode === "測試";
  const [contentType, setContentType] = useState<QuizContentType>("grammar");
  const [grammar, setGrammar] = useState<string>("");
  const [jlpt, setJlpt] = useState<string>("");
  const [grammars, setGrammars] = useState<GrammarPoint[]>([]);
  const [statsRange, setStatsRange] = useState<StatsRange>(7);
  const [stats, setStats] = useState<Stats | null>(null);
  const [statsLoading, setStatsLoading] = useState(true);
  const [statsErr, setStatsErr] = useState("");
  const [mode, setMode] = useState<Mode>(sealedExam ? "session" : "single");
  const [target, setTarget] = useState(10);
  const [remainingSec, setRemainingSec] = useState(EXAM_DURATION_SEC);

  const [turns, setTurns] = useState<SessionTurn[]>([]);
  const [current, setCurrent] = useState<Question | null>(null);
  const [answer, setAnswer] = useState("");
  const [result, setResult] = useState<GradeResult | null>(null);
  const [err, setErr] = useState("");
  const [notice, setNotice] = useState("");
  const examSealedRef = useRef(false);

  useEffect(() => {
    api.listGrammar().then((r) => setGrammars(r.points || []));
  }, []);

  async function loadStats(range: StatsRange = statsRange) {
    setStatsLoading(true);
    setStatsErr("");
    try {
      const s = await api.stats(range === 0 ? undefined : range);
      setStats(s);
    } catch (e) {
      setStatsErr(String(e));
    } finally {
      setStatsLoading(false);
    }
  }

  async function pickNext(seenIDs: string[], sessionTurns = turns) {
    setErr("");
    setNotice("");
    setResult(null);
    setAnswer("");
    try {
      const q = await api.nextQuestion({
        type: contentType,
        jlpt: jlpt || undefined,
        grammar: contentType === "grammar" ? grammar || undefined : undefined,
        exclude: seenIDs,
      });
      setCurrent(q);
    } catch (e) {
      // 404 from /api/quiz/next means the filter (or exclude list) has run
      // out of candidates. In session mode, gracefully end the session
      // early with whatever's been answered. In single mode, surface the
      // empty state so the user knows to relax the filter.
      const exhausted = isApiError(e) && e.status === 404;
      setCurrent(null);
      if (exhausted && mode === "session" && sessionTurns.length > 0) {
        setMode("summary");
        return;
      }
      setErr(exhausted ? "目前的條件下沒有更多題目了，換個條件試試。" : String(e));
    }
  }

  function changeContentType(next: QuizContentType) {
    setContentType(next);
    setGrammar("");
    setTurns([]);
    setResult(null);
    setAnswer("");
    setNotice("");
    setErr("");
    if (mode === "summary") {
      setMode(sealedExam ? "session" : "single");
    }
    if (sealedExam) {
      examSealedRef.current = false;
      setRemainingSec(EXAM_DURATION_SEC);
    }
  }

  async function startSingle() {
    setMode(sealedExam ? "session" : "single");
    setTurns([]);
    if (sealedExam) {
      examSealedRef.current = false;
      setRemainingSec(EXAM_DURATION_SEC);
    }
    pickNext([]);
  }

  async function startSession() {
    setMode("session");
    setTurns([]);
    setResult(null);
    if (sealedExam) {
      examSealedRef.current = false;
      setRemainingSec(EXAM_DURATION_SEC);
    }
    pickNext([]);
  }

  async function submit(e: React.FormEvent) {
    e.preventDefault();
    if (sealedExam && (remainingSec <= 0 || examSealedRef.current)) return;
    if (!current || !answer.trim()) return;
    try {
      const r = await api.answer(current.id, answer.trim());
      if (sealedExam && examSealedRef.current) return;
      setResult(sealedExam ? null : r);
      loadStats();

      if (mode === "session") {
        const newTurns = [...turns, { question: current, result: r }];
        setTurns(newTurns);
        if (newTurns.length >= target) {
          setCurrent(null);
          setAnswer("");
          setMode("summary");
          return;
        }
        if (sealedExam) {
          await pickNext(newTurns.map((t) => t.question.id), newTurns);
        }
      }
    } catch (e) {
      if (isApiError(e, "question_not_found")) {
        const seen = mode === "session"
          ? [...turns.map((t) => t.question.id), current.id]
          : [current.id];
        await pickNext(seen);
        setNotice("這題剛被更新，已替你換下一題。");
        return;
      }
      setErr(String(e));
    }
  }

  async function next() {
    if (mode === "session") {
      const seen = turns.map((t) => t.question.id);
      pickNext(seen);
    } else {
      pickNext([]);
    }
  }

  // initial mount
  useEffect(() => {
    pickNext([]);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  useEffect(() => {
    loadStats(statsRange);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [statsRange]);

  // refetch on filter change while not mid-session
  useEffect(() => {
    if (mode === "summary") return;
    if (mode === "session" && turns.length > 0 && !result) return;
    pickNext(mode === "session" ? turns.map((t) => t.question.id) : []);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [contentType, grammar, jlpt]);

  const countdownActive = sealedExam && mode === "session" && remainingSec > 0;
  const examExpired = sealedExam && mode === "session" && remainingSec <= 0;

  useEffect(() => {
    if (!countdownActive) return;
    const countdown = window.setInterval(() => {
      setRemainingSec((prev) => Math.max(0, prev - 1));
    }, 1000);
    return () => window.clearInterval(countdown);
  }, [countdownActive]);

  useEffect(() => {
    if (!examExpired || !current) return;
    examSealedRef.current = true;
    setCurrent(null);
    setAnswer("");
    setResult(null);
    setMode("summary");
  }, [current, examExpired]);

  return (
    <section>
      <StatsPanel
        contentType={contentType}
        stats={stats}
        loading={statsLoading}
        error={statsErr}
        range={statsRange}
        onRangeChange={setStatsRange}
        onRefresh={() => loadStats()}
      />

      <Toolbar
        mode={mode}
        target={target}
        setTarget={setTarget}
        contentType={contentType}
        setContentType={changeContentType}
        grammar={grammar}
        setGrammar={setGrammar}
        jlpt={jlpt}
        setJlpt={setJlpt}
        grammars={grammars}
        sealedExam={sealedExam}
        onStartSingle={startSingle}
        onStartSession={startSession}
      />

      {mode === "session" && (
        <SessionProgress
          total={target}
          turns={turns}
          remainingSec={sealedExam ? remainingSec : undefined}
        />
      )}

      {err && <p className="text-red-600 text-sm mb-4">{err}</p>}
      {notice && <p className="text-blue-700 text-sm mb-4">{notice}</p>}

      {mode === "summary" ? (
        <Summary turns={turns} contentType={contentType} onRestart={startSession} />
      ) : current ? (
        <QuestionForm
          q={current}
          answer={answer}
          setAnswer={setAnswer}
          onSubmit={submit}
          disabled={!!result || examExpired}
        />
      ) : null}

      {result && mode !== "summary" && !sealedExam && (
        <ResultBox
          result={result}
          hint={current?.hint}
          onNext={next}
          sessionMode={mode === "session"}
          onNavigateGrammar={onNavigateGrammar}
        />
      )}
    </section>
  );
}

function StatsPanel({
  contentType,
  stats,
  loading,
  error,
  range,
  onRangeChange,
  onRefresh,
}: {
  contentType: QuizContentType;
  stats: Stats | null;
  loading: boolean;
  error: string;
  range: StatsRange;
  onRangeChange: (range: StatsRange) => void;
  onRefresh: () => void;
}) {
  const weakGrammar = (stats?.by_grammar ?? [])
    .filter((g) => g.total > 0)
    .sort((a, b) => a.accuracy - b.accuracy || b.total - a.total)
    .slice(0, 3);
  const weakVocab = (stats?.by_vocab ?? [])
    .filter((v) => v.total > 0)
    .sort((a, b) => a.accuracy - b.accuracy || b.total - a.total)
    .slice(0, 3);
  const topErrors = (stats?.by_error_class ?? []).slice(0, 3);
  const recentWrong = (stats?.recent_wrong ?? []).slice(0, 3);

  return (
    <div className="bg-white border border-slate-200 rounded-md p-5 mb-6">
      <div className="flex flex-wrap items-center justify-between gap-3 mb-4">
        <div>
          <h2 className="text-base font-medium">練習狀態</h2>
          <p className="text-xs text-slate-500 mt-1">
            根據已作答紀錄找出需要複習的{contentType === "grammar" ? "文法點" : "單字"}。
          </p>
        </div>
        <div className="flex items-center gap-2">
          <select
            value={range}
            onChange={(e) => onRangeChange(Number(e.target.value) as StatsRange)}
            className="px-2 py-1.5 border border-slate-300 rounded-md text-xs bg-white"
            aria-label="成績範圍"
          >
            <option value={7}>最近 7 天</option>
            <option value={30}>最近 30 天</option>
            <option value={0}>全部</option>
          </select>
          <button
            onClick={onRefresh}
            className="px-2.5 py-1.5 border border-slate-300 rounded-md text-xs bg-white text-slate-700"
          >
            更新
          </button>
        </div>
      </div>

      {loading ? (
        <p className="text-sm text-slate-500">讀取練習紀錄中...</p>
      ) : error ? (
        <p className="text-sm text-red-600">無法讀取練習狀態：{error}</p>
      ) : !stats || stats.total_attempts === 0 ? (
        <p className="text-sm text-slate-500">尚無作答紀錄。先完成幾題後，這裡會顯示弱點與最近錯題。</p>
      ) : (
        <div className="space-y-5">
          <div className="grid grid-cols-3 gap-3 text-center">
            <Stat label="作答" value={stats.total_attempts} />
            <Stat label="答錯" value={stats.wrong} accent="text-amber-600" />
            <Stat label="正確率" value={formatPercent(stats.accuracy)} accent={stats.accuracy >= 0.8 ? "text-green-600" : ""} />
          </div>

          <div className="grid gap-4 md:grid-cols-3">
            <StatsList
              title="優先複習"
              empty={contentType === "grammar" ? "目前沒有可排序的文法點。" : "目前沒有可排序的單字。"}
            >
              {contentType === "grammar"
                ? weakGrammar.map((g) => (
                    <li key={g.grammar_point} className="space-y-1">
                      <div className="flex justify-between gap-3">
                        <span className="truncate">{g.grammar_point}</span>
                        <span className="text-slate-500">{formatPercent(g.accuracy)}</span>
                      </div>
                      <div className="h-1.5 bg-slate-100 rounded-full overflow-hidden">
                        <div className="h-full bg-amber-500" style={{ width: `${Math.round(g.accuracy * 100)}%` }} />
                      </div>
                    </li>
                  ))
                : weakVocab.map((v) => (
                    <li key={v.vocab_item} className="space-y-1">
                      <div className="flex justify-between gap-3">
                        <span className="truncate">{v.vocab_item}</span>
                        <span className="text-slate-500">{formatPercent(v.accuracy)}</span>
                      </div>
                      <div className="h-1.5 bg-slate-100 rounded-full overflow-hidden">
                        <div className="h-full bg-amber-500" style={{ width: `${Math.round(v.accuracy * 100)}%` }} />
                      </div>
                    </li>
                  ))}
            </StatsList>

            <StatsList title="常見錯誤" empty="目前沒有錯誤類型。">
              {topErrors.map((e) => (
                <li key={`${e.grammar_point}:${e.error_class}`} className="flex justify-between gap-3">
                  <span className="truncate">{e.grammar_point}</span>
                  <span className="text-slate-500 whitespace-nowrap">{e.error_class} · {e.count}</span>
                </li>
              ))}
            </StatsList>

            <StatsList title="最近錯題" empty="最近沒有答錯紀錄。">
              {recentWrong.map((w) => (
                <li key={`${w.question_id}:${w.created_at}`} className="space-y-1">
                  <p className="truncate">{w.prompt.replace("___", `[${w.expected}]`)}</p>
                  <p className="text-slate-500 truncate">
                    {w.grammar_point} · 你的答案 {w.user_answer}
                  </p>
                </li>
              ))}
            </StatsList>
          </div>
        </div>
      )}
    </div>
  );
}

function StatsList({
  title,
  empty,
  children,
}: {
  title: string;
  empty: string;
  children: React.ReactNode;
}) {
  const count = Array.isArray(children) ? children.length : 0;
  return (
    <div>
      <h3 className="text-xs font-medium text-slate-500 mb-2">{title}</h3>
      {count === 0 ? (
        <p className="text-xs text-slate-400">{empty}</p>
      ) : (
        <ul className="space-y-2 text-xs">{children}</ul>
      )}
    </div>
  );
}

function Toolbar(props: {
  mode: Mode;
  target: number;
  setTarget: (n: number) => void;
  contentType: QuizContentType;
  setContentType: (type: QuizContentType) => void;
  grammar: string;
  setGrammar: (s: string) => void;
  jlpt: string;
  setJlpt: (s: string) => void;
  grammars: GrammarPoint[];
  sealedExam: boolean;
  onStartSingle: () => void;
  onStartSession: () => void;
}) {
  return (
    <div className="flex flex-wrap items-center gap-2 mb-6">
      <div className="flex border border-slate-300 rounded-md overflow-hidden text-sm">
        {props.sealedExam ? (
          <span className="px-3 py-2 bg-blue-600 text-white">測試</span>
        ) : (
          <>
            <button
              onClick={props.onStartSingle}
              className={
                "px-3 py-2 " +
                (props.mode === "single" ? "bg-blue-600 text-white" : "bg-white text-slate-700")
              }
            >
              單題
            </button>
            <button
              onClick={props.onStartSession}
              className={
                "px-3 py-2 border-l border-slate-300 " +
                (props.mode === "session" || props.mode === "summary"
                  ? "bg-blue-600 text-white"
                  : "bg-white text-slate-700")
              }
            >
              套題 ({props.target})
            </button>
          </>
        )}
      </div>
      <div className="flex border border-slate-300 rounded-md overflow-hidden text-sm">
        <button
          onClick={() => props.setContentType("grammar")}
          className={
            "px-3 py-2 " +
            (props.contentType === "grammar" ? "bg-blue-600 text-white" : "bg-white text-slate-700")
          }
        >
          文法
        </button>
        <button
          onClick={() => props.setContentType("vocab")}
          className={
            "px-3 py-2 border-l border-slate-300 " +
            (props.contentType === "vocab" ? "bg-blue-600 text-white" : "bg-white text-slate-700")
          }
        >
          單字
        </button>
      </div>
      {props.mode !== "summary" && (
        <input
          type="number"
          value={props.target}
          onChange={(e) => props.setTarget(Math.max(1, Math.min(50, +e.target.value)))}
          min={1}
          max={50}
          className="w-16 px-2 py-2 border border-slate-300 rounded-md text-sm"
          aria-label="題目數量"
        />
      )}
      <select
        value={props.jlpt}
        onChange={(e) => props.setJlpt(e.target.value)}
        className="px-3 py-2 border border-slate-300 rounded-md text-sm bg-white"
      >
        <option value="">全等級</option>
        <option value="N5">N5</option>
        <option value="N4">N4</option>
        <option value="N3">N3</option>
        <option value="N2">N2</option>
        <option value="N1">N1</option>
      </select>
      {props.contentType === "grammar" && (
        <select
          value={props.grammar}
          onChange={(e) => props.setGrammar(e.target.value)}
          className="px-3 py-2 border border-slate-300 rounded-md text-sm bg-white"
        >
          <option value="">全文法點</option>
          {props.grammars.map((g) => (
            <option key={g.slug} value={g.slug}>
              {g.jlpt_level} · {g.title_ja}
            </option>
          ))}
        </select>
      )}
    </div>
  );
}

function SessionProgress({
  total,
  turns,
  remainingSec,
}: {
  total: number;
  turns: SessionTurn[];
  remainingSec?: number;
}) {
  const correct = turns.filter((t) => t.result?.correct).length;
  return (
    <div className="mb-4">
      {remainingSec !== undefined && (
        <div className="flex justify-end mb-2">
          <span
            className={
              "text-sm font-medium tabular-nums " +
              (remainingSec <= 60 ? "text-red-500" : "text-slate-700")
            }
          >
            {formatCountdown(remainingSec)}
          </span>
        </div>
      )}
      <div className="flex justify-between text-xs text-slate-500 mb-1">
        <span>{turns.length} / {total}</span>
        <span>對 {correct} · 錯 {turns.length - correct}</span>
      </div>
      <div className="h-1 bg-slate-200 rounded-full overflow-hidden">
        <div
          className="h-full bg-blue-600 transition-all"
          style={{ width: `${(turns.length / total) * 100}%` }}
        />
      </div>
    </div>
  );
}

function formatCountdown(seconds: number) {
  const safeSeconds = Math.max(0, seconds);
  const minutes = Math.floor(safeSeconds / 60);
  const remainder = safeSeconds % 60;
  return `${minutes.toString().padStart(2, "0")}:${remainder.toString().padStart(2, "0")}`;
}

function QuestionForm(props: {
  q: Question;
  answer: string;
  setAnswer: (s: string) => void;
  onSubmit: (e: React.FormEvent) => void;
  disabled: boolean;
}) {
  return (
    <form
      onSubmit={props.onSubmit}
      className="bg-white border border-slate-200 rounded-md p-6 space-y-4"
    >
      <p className="text-2xl leading-relaxed">
        {props.q.prompt.split("___").map((part, i, arr) => (
          <span key={i}>
            {part}
            {i < arr.length - 1 && (
              <input
                autoFocus
                value={props.answer}
                onChange={(e) => props.setAnswer(e.target.value)}
                disabled={props.disabled}
                className="inline-block mx-1 w-32 border-b-2 border-blue-500 focus:outline-none px-1 text-2xl text-center disabled:text-slate-400"
              />
            )}
          </span>
        ))}
      </p>
      <button
        type="submit"
        disabled={props.disabled}
        className="px-4 py-2 bg-blue-600 text-white rounded-md disabled:opacity-50"
      >
        送出
      </button>
    </form>
  );
}

function ResultBox({
  result,
  hint,
  onNext,
  sessionMode,
  onNavigateGrammar,
}: {
  result: GradeResult;
  hint?: string;
  onNext: () => void;
  sessionMode: boolean;
  onNavigateGrammar?: (slug: string) => void;
}) {
  const isGrammar = result.content_type === "grammar";
  const isVocab = result.content_type === "vocab";
  const showErrorClass = result.error_class && result.error_class !== "generic";
  const explanation = result.explanation_zh && !result.explanation_zh.includes("暫無針對此題")
    ? result.explanation_zh
    : null;

  return (
    <div
      className={
        "mt-4 p-6 rounded-md border space-y-3 " +
        (result.correct
          ? "bg-green-50 border-green-200 text-green-900"
          : "bg-amber-50 border-amber-200 text-amber-900")
      }
    >
      <p className="font-medium">{result.correct ? "正解！" : "不對"}</p>

      {!result.correct && (
        <>
          <p className="text-sm">
            你的答案：<code>{result.user_answer}</code>
            {"　"}正確答案：<code>{result.expected}</code>
            {showErrorClass && (
              <span className="ml-2 text-xs px-1.5 py-0.5 bg-amber-100 rounded">
                {result.error_class}
              </span>
            )}
          </p>

          {/* Specific feedback from grader (only if non-placeholder) */}
          {explanation && (
            <p className="text-sm leading-relaxed border-l-2 border-amber-400 pl-3">
              {explanation}
            </p>
          )}

          {/* Grammar explanation from corpus */}
          {isGrammar && result.item_detail_zh && (
            <div className="text-sm leading-relaxed space-y-1">
              <p className="text-xs font-medium text-amber-700">文法說明</p>
              <p className="line-clamp-4">{result.item_detail_zh}</p>
              {onNavigateGrammar && (
                <button
                  onClick={() => onNavigateGrammar(result.grammar_point)}
                  className="text-xs text-blue-700 hover:underline mt-1"
                >
                  查看完整說明 →
                </button>
              )}
            </div>
          )}

          {/* Vocab meaning */}
          {isVocab && result.item_detail_zh && (
            <div className="text-sm space-y-1">
              <p className="text-xs font-medium text-amber-700">詞義</p>
              <p>{result.item_detail_zh}</p>
            </div>
          )}

          {hint && <p className="text-xs text-slate-600">提示：{hint}</p>}
        </>
      )}

      <button
        onClick={onNext}
        className="mt-1 px-3 py-1.5 bg-white border border-slate-300 rounded-md text-sm"
      >
        {sessionMode ? "下一題" : "再來一題"} →
      </button>
    </div>
  );
}

function Summary({
  turns,
  contentType,
  onRestart,
}: {
  turns: SessionTurn[];
  contentType: QuizContentType;
  onRestart: () => void;
}) {
  const correct = turns.filter((t) => t.result?.correct).length;
  const accuracy = turns.length > 0 ? correct / turns.length : 0;
  const wrong = turns.filter((t) => t.result && !t.result.correct);

  // group by grammar_point
  const byGP = new Map<string, { total: number; correct: number }>();
  for (const t of turns) {
    if (!t.result) continue;
    const k = t.question.grammar_point;
    const e = byGP.get(k) ?? { total: 0, correct: 0 };
    e.total += 1;
    if (t.result.correct) e.correct += 1;
    byGP.set(k, e);
  }

  return (
    <div className="space-y-6">
      <div className="bg-white border border-slate-200 rounded-md p-6">
        <h2 className="text-lg font-medium mb-4">本次成績</h2>
        <div className="grid grid-cols-3 gap-4 text-center">
          <Stat label="總題數" value={turns.length} />
          <Stat label="答對" value={correct} accent="text-green-600" />
          <Stat label="正確率" value={`${Math.round(accuracy * 100)}%`} />
        </div>
      </div>

      {byGP.size > 0 && (
        <div className="bg-white border border-slate-200 rounded-md p-6">
          <h3 className="text-sm font-medium mb-3">
            {contentType === "grammar" ? "分文法點" : "分單字"}
          </h3>
          <ul className="space-y-2 text-sm">
            {Array.from(byGP.entries()).map(([k, v]) => (
              <li key={k} className="flex justify-between">
                <span>{k}</span>
                <span className="text-slate-500">
                  {v.correct}/{v.total} ({Math.round((v.correct / v.total) * 100)}%)
                </span>
              </li>
            ))}
          </ul>
        </div>
      )}

      {wrong.length > 0 && (
        <div className="bg-white border border-slate-200 rounded-md p-6">
          <h3 className="text-sm font-medium mb-3">答錯的題目</h3>
          <ul className="space-y-3 text-sm">
            {wrong.map((t, i) => (
              <li key={i} className="border-l-2 border-amber-400 pl-3">
                <p className="text-base">{t.question.prompt.replace("___", `[${t.result!.expected}]`)}</p>
                <p className="text-xs text-slate-500 mt-1">
                  你的答案 <code>{t.result!.user_answer}</code> ・ {t.question.grammar_point}
                  {t.result!.error_class && <> · {t.result!.error_class}</>}
                </p>
              </li>
            ))}
          </ul>
        </div>
      )}

      <button
        onClick={onRestart}
        className="px-4 py-2 bg-blue-600 text-white rounded-md"
      >
        再來一輪
      </button>
    </div>
  );
}

function Stat({ label, value, accent }: { label: string; value: string | number; accent?: string }) {
  return (
    <div>
      <div className={"text-3xl font-medium " + (accent ?? "")}>{value}</div>
      <div className="text-xs text-slate-500 mt-1">{label}</div>
    </div>
  );
}

function formatPercent(value: number) {
  return `${Math.round(value * 100)}%`;
}
