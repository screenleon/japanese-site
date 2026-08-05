import { JLPT_LEVELS, type JlptLevel } from "../jlptLevels";

/**
 * Shared N5–N1 level chip row for reference tabs (grammar, vocab, future kanji browse).
 * Selection state is owned by the parent tab (UI-001).
 */
export function LevelPicker({
  levels = JLPT_LEVELS,
  selected,
  onSelect,
  subtitle,
}: {
  levels?: readonly JlptLevel[];
  selected: string;
  onSelect: (level: JlptLevel) => void;
  /** Static label under each chip, or per-level (e.g. count). */
  subtitle: string | ((level: JlptLevel) => string);
}) {
  return (
    <div data-testid="level-picker" className="grid gap-3 sm:grid-cols-5">
      {levels.map((level) => {
        const active = selected === level;
        const sub = typeof subtitle === "function" ? subtitle(level) : subtitle;
        return (
          <button
            key={level}
            type="button"
            onClick={() => onSelect(level)}
            className={
              "rounded-md border px-4 py-3 text-left transition-colors " +
              (active
                ? "border-blue-600 bg-blue-50 text-blue-700"
                : "border-slate-200 bg-white text-slate-700 hover:border-slate-300")
            }
          >
            <div className="text-lg font-semibold">{level}</div>
            <div className="text-xs text-slate-500">{sub}</div>
          </button>
        );
      })}
    </div>
  );
}

export type { JlptLevel };
