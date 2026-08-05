import { describe, expect, it } from "vitest";
import { isJlptLevel, JLPT_LEVELS } from "./jlptLevels";

describe("JLPT_LEVELS", () => {
  /**
   * Verifies the canonical study navigation order is N5 → N1.
   * Steps:
   * 1. Read JLPT_LEVELS.
   * 2. Assert exact length and sequence.
   */
  it("lists N5 through N1 in ascending difficulty order", () => {
    expect([...JLPT_LEVELS]).toEqual(["N5", "N4", "N3", "N2", "N1"]);
  });
});

describe("isJlptLevel", () => {
  /**
   * Verifies every canonical level is accepted by the type guard.
   * Steps:
   * 1. Call isJlptLevel for each JLPT_LEVELS entry.
   * 2. Assert true for each.
   */
  it("accepts every canonical N5–N1 level", () => {
    for (const level of JLPT_LEVELS) {
      expect(isJlptLevel(level)).toBe(true);
    }
  });

  /**
   * Verifies non-canonical strings are rejected.
   * Steps:
   * 1. Call isJlptLevel with empty, wrong case, out-of-range, and unrelated inputs.
   * 2. Assert false for each.
   */
  it("rejects invalid or non-canonical inputs", () => {
    expect(isJlptLevel("")).toBe(false);
    expect(isJlptLevel("n5")).toBe(false);
    expect(isJlptLevel("N6")).toBe(false);
    expect(isJlptLevel("N0")).toBe(false);
    expect(isJlptLevel("N5 ")).toBe(false);
    expect(isJlptLevel("JLPT-N5")).toBe(false);
  });
});
