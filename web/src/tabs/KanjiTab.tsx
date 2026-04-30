import { useState } from "react";
import { api, type Kanji } from "../api";

export function KanjiTab() {
  const [ch, setCh] = useState("");
  const [data, setData] = useState<Kanji | null>(null);
  const [err, setErr] = useState("");

  async function lookup(e: React.FormEvent) {
    e.preventDefault();
    const q = ch.trim();
    if (!q) return;
    setErr("");
    try {
      const r = await api.getKanji(q[0]);
      setData(r);
    } catch (e) {
      setErr(String(e));
      setData(null);
    }
  }

  return (
    <section>
      <form onSubmit={lookup} className="flex gap-2 mb-4">
        <input
          value={ch}
          onChange={(e) => setCh(e.target.value)}
          placeholder="輸入單一漢字（例：山、愛、食）"
          className="flex-1 px-3 py-2 border border-slate-300 rounded-md text-base"
        />
        <button type="submit" className="px-4 py-2 bg-blue-600 text-white rounded-md">
          查詢
        </button>
      </form>
      {err && <p className="text-red-600 text-sm">{err}</p>}
      {data && (
        <div className="bg-white border border-slate-200 rounded-md p-6 grid grid-cols-[auto_1fr] gap-4 items-start">
          <div className="text-7xl font-serif">{data.character}</div>
          <dl className="space-y-2 text-sm">
            <Row label="意味" value={data.meaning_ja || "日本語の意味は準備中です。"} />
            <Row label="繁中" value={data.meaning_zh || "繁中意義待補"} />
            <Row label="音讀" value={data.onyomi} />
            <Row label="訓讀" value={data.kunyomi} />
            <Row label="JLPT" value={data.jlpt_level} />
            <Row label="筆畫" value={data.stroke_count?.toString()} />
            <Row label="教育部年級" value={data.grade?.toString()} />
            <Row label="來源" value={`${data.source} (${data.license})`} />
          </dl>
        </div>
      )}
    </section>
  );
}

function Row({ label, value }: { label: string; value?: string }) {
  if (!value) return null;
  return (
    <div className="flex">
      <dt className="w-20 text-slate-500">{label}</dt>
      <dd className="flex-1 break-all">{value}</dd>
    </div>
  );
}
