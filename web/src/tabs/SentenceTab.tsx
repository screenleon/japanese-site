import { useEffect, useState } from "react";
import { api, type Sentence } from "../api";

export function SentenceTab() {
  const [jlpt, setJlpt] = useState("");
  const [s, setS] = useState<Sentence | null>(null);
  const [err, setErr] = useState("");

  async function refresh() {
    setErr("");
    try {
      const r = await api.randomSentence(jlpt || undefined);
      setS(r);
    } catch (e) {
      setErr(String(e));
      setS(null);
    }
  }

  useEffect(() => {
    refresh();
  }, []);

  return (
    <section>
      <div className="flex gap-2 mb-4">
        <select
          value={jlpt}
          onChange={(e) => setJlpt(e.target.value)}
          className="px-3 py-2 border border-slate-300 rounded-md text-sm bg-white"
        >
          <option value="">全等級</option>
          <option value="N5">N5</option>
          <option value="N4">N4</option>
          <option value="N3">N3</option>
          <option value="N2">N2</option>
          <option value="N1">N1</option>
        </select>
        <button onClick={refresh} className="px-4 py-2 bg-blue-600 text-white rounded-md">
          下一句
        </button>
      </div>
      {err && <p className="text-red-600 text-sm">{err}</p>}
      {s && (
        <div className="bg-white border border-slate-200 rounded-md p-6">
          <p className="text-2xl leading-relaxed">{s.text_ja}</p>
          <div className="text-xs text-slate-400 mt-4 space-x-2">
            <span>id={s.id}</span>
            {s.jlpt_level && (
              <span className="inline-block px-1.5 py-0.5 bg-blue-50 text-blue-700 rounded">
                {s.jlpt_level}
              </span>
            )}
            <span>{s.source}</span>
          </div>
        </div>
      )}
    </section>
  );
}
