import { useState } from "react";
import { api, type VocabRow } from "../api";

export function VocabTab() {
  const [query, setQuery] = useState("");
  const [jlpt, setJlpt] = useState("");
  const [rows, setRows] = useState<VocabRow[]>([]);
  const [err, setErr] = useState("");
  const [loading, setLoading] = useState(false);

  async function search(e: React.FormEvent) {
    e.preventDefault();
    if (!query.trim()) return;
    setLoading(true);
    setErr("");
    try {
      const r = await api.searchVocab(query.trim(), jlpt || undefined);
      setRows(r.results || []);
    } catch (e) {
      setErr(String(e));
      setRows([]);
    } finally {
      setLoading(false);
    }
  }

  return (
    <section>
      <form onSubmit={search} className="flex gap-2 mb-4">
        <input
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="食べ / おはよう / arigato / etc."
          className="flex-1 px-3 py-2 border border-slate-300 rounded-md text-base"
        />
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
        <button
          type="submit"
          disabled={loading}
          className="px-4 py-2 bg-blue-600 text-white rounded-md disabled:opacity-50"
        >
          搜尋
        </button>
      </form>
      {err && <p className="text-red-600 text-sm">{err}</p>}
      <ul className="divide-y divide-slate-200 bg-white rounded-md border border-slate-200">
        {rows.map((r) => (
          <li key={r.id} className="p-4 flex items-baseline gap-4">
            <div className="flex-shrink-0 w-32">
              <div className="text-lg font-medium">{r.headword}</div>
              <div className="text-sm text-slate-500">{r.reading}</div>
            </div>
            <div className="flex-1 min-w-0">
              <div className="text-sm">{r.gloss_en}</div>
              <div className="text-xs text-slate-400 mt-1">
                {r.pos}
                {r.jlpt_level && (
                  <span className="ml-2 inline-block px-1.5 py-0.5 bg-blue-50 text-blue-700 rounded">
                    {r.jlpt_level}
                  </span>
                )}
              </div>
            </div>
          </li>
        ))}
        {!loading && rows.length === 0 && query && (
          <li className="p-4 text-sm text-slate-400">無結果</li>
        )}
      </ul>
    </section>
  );
}
