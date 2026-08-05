/** Canonical JLPT level order for study navigation (N5 → N1). */
export const JLPT_LEVELS = ["N5", "N4", "N3", "N2", "N1"] as const;

export type JlptLevel = (typeof JLPT_LEVELS)[number];

export function isJlptLevel(value: string): value is JlptLevel {
  return (JLPT_LEVELS as readonly string[]).includes(value);
}
