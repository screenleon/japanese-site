import { useEffect, useState } from "react";
import { api, type Question, type GradeResult, type GrammarPoint } from "../api";

type Mode = "single" | "session" | "summary";

interface SessionTurn {
  question: Question;
  result: GradeResult | null;
}

export function QuizTab() {
  const [grammar, setGrammar] = useState<string>("");
  const [jlpt, setJlpt] = useState<string>("");
  const [grammars, setGrammars] = useState<GrammarPoint[]>([]);
  const [mode, setMode] = useState<Mode>("single");
  const [target, setTarget] = useState(10);

  const [turns, setTurns] = useState<SessionTurn[]>([]);
  const [current, setCurrent] = useState<Question | null>(null);
  const [answer, setAnswer] = useState("");
  const [result, setResult] = useState<GradeResult | null>(null);
  const [err, setErr] = useState("");

  useEffect(() => {
    api.listGrammar().then((r) => setGrammars(r.points || []));
  }, []);

  async function pickNext(seenIDs: string[]) {
    setErr("");
    setResult(null);
    setAnswer("");
    try {
      const q = await api.nextQuestion({
        jlpt: jlpt || undefined,
        grammar: grammar || undefined,
        exclude: seenIDs,
      });
      setCurrent(q);
    } catch (e) {
      // 404 from /api/quiz/next means the filter (or exclude list) has run
      // out of candidates. In session mode, gracefully end the session
      // early with whatever's been answered. In single mode, surface the
      // empty state so the user knows to relax the filter.
      const msg = String(e);
      const exhausted = msg.includes("404");
      setCurrent(null);
      if (exhausted && mode === "session" && turns.length > 0) {
        setMode("summary");
        return;
      }
      setErr(exhausted ? "目前的條件下沒有更多題目了，換個文法點或等級試試。" : msg);
    }
  }

  async function startSingle() {
    setMode("single");
    setTurns([]);
    pickNext([]);
  }

  async function startSession() {
    setMode("session");
    setTurns([]);
    setResult(null);
    pickNext([]);
  }

  async function submit(e: React.FormEvent) {
    e.preventDefault();
    if (!current || !answer.trim()) return;
    try {
      const r = await api.answer(current.id, answer.trim());
      setResult(r);

      if (mode === "session") {
        const newTurns = [...turns, { question: current, result: r }];
        setTurns(newTurns);
        if (newTurns.length >= target) {
          setMode("summary");
        }
      }
    } catch (e) {
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

  // refetch on filter change while not mid-session
  useEffect(() => {
    if (mode === "summary") return;
    if (mode === "session" && turns.length > 0 && !result) return;
    pickNext(mode === "session" ? turns.map((t) => t.question.id) : []);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [grammar, jlpt]);

  return (
    <section>
      <Toolbar
        mode={mode}
        target={target}
        setTarget={setTarget}
        grammar={grammar}
        setGrammar={setGrammar}
        jlpt={jlpt}
        setJlpt={setJlpt}
        grammars={grammars}
        onStartSingle={startSingle}
        onStartSession={startSession}
      />

      {mode === "session" && (
        <SessionProgress total={target} turns={turns} />
      )}

      {err && <p className="text-red-600 text-sm mb-4">{err}</p>}

      {mode === "summary" ? (
        <Summary turns={turns} onRestart={startSession} />
      ) : current ? (
        <QuestionForm
          q={current}
          answer={answer}
          setAnswer={setAnswer}
          onSubmit={submit}
          disabled={!!result}
        />
      ) : null}

      {result && mode !== "summary" && (
        <ResultBox result={result} onNext={next} sessionMode={mode === "session"} />
      )}
    </section>
  );
}

function Toolbar(props: {
  mode: Mode;
  target: number;
  setTarget: (n: number) => void;
  grammar: string;
  setGrammar: (s: string) => void;
  jlpt: string;
  setJlpt: (s: string) => void;
  grammars: GrammarPoint[];
  onStartSingle: () => void;
  onStartSession: () => void;
}) {
  return (
    <div className="flex flex-wrap items-center gap-2 mb-6">
      <div className="flex border border-slate-300 rounded-md overflow-hidden text-sm">
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
    </div>
  );
}

function SessionProgress({ total, turns }: { total: number; turns: SessionTurn[] }) {
  const correct = turns.filter((t) => t.result?.correct).length;
  return (
    <div className="mb-4">
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
      <div className="text-xs text-slate-400">
        {props.q.jlpt_level} · {props.q.grammar_point}
        {props.q.hint && <span className="ml-2">提示：{props.q.hint}</span>}
      </div>
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
  onNext,
  sessionMode,
}: {
  result: GradeResult;
  onNext: () => void;
  sessionMode: boolean;
}) {
  return (
    <div
      className={
        "mt-4 p-6 rounded-md border " +
        (result.correct
          ? "bg-green-50 border-green-200 text-green-900"
          : "bg-amber-50 border-amber-200 text-amber-900")
      }
    >
      <p className="font-medium">{result.correct ? "正解！" : "不對"}</p>
      {!result.correct && (
        <>
          <p className="text-sm mt-1">
            你的答案：<code>{result.user_answer}</code>　正確答案：
            <code>{result.expected}</code>
            {result.error_class && (
              <span className="ml-2 text-xs px-1.5 py-0.5 bg-amber-100 rounded">
                {result.error_class}
              </span>
            )}
          </p>
          <p className="text-sm mt-3 leading-relaxed">{result.explanation_zh}</p>
        </>
      )}
      <button
        onClick={onNext}
        className="mt-4 px-3 py-1.5 bg-white border border-slate-300 rounded-md text-sm"
      >
        {sessionMode ? "下一題" : "再來一題"} →
      </button>
    </div>
  );
}

function Summary({ turns, onRestart }: { turns: SessionTurn[]; onRestart: () => void }) {
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
          <h3 className="text-sm font-medium mb-3">分文法點</h3>
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
