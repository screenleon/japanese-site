import { fireEvent, render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import { LevelPicker } from "./LevelPicker";

describe("LevelPicker", () => {
  /**
   * Verifies the default N5–N1 chip row reports selection to the parent.
   * Steps:
   * 1. Render LevelPicker with selected N5 and a static subtitle.
   * 2. Assert five level buttons are present.
   * 3. Click N3 and assert onSelect receives "N3".
   */
  it("renders N5–N1 chips and reports selection", () => {
    const onSelect = vi.fn();
    render(
      <LevelPicker selected="N5" onSelect={onSelect} subtitle="單字學習" />
    );

    expect(screen.getByTestId("level-picker").querySelectorAll("button")).toHaveLength(5);
    fireEvent.click(screen.getByRole("button", { name: /N3/ }));
    expect(onSelect).toHaveBeenCalledWith("N3");
  });

  /**
   * Verifies per-level subtitle callbacks render under the matching chip.
   * Steps:
   * 1. Render LevelPicker with a subtitle function returning count text for N4.
   * 2. Locate the N4 button.
   * 3. Assert it contains the per-level subtitle string.
   */
  it("supports per-level subtitles", () => {
    render(
      <LevelPicker
        selected="N4"
        onSelect={() => {}}
        subtitle={(level) => (level === "N4" ? "12 文法" : "0 文法")}
      />
    );
    expect(screen.getByRole("button", { name: /N4/ })).toHaveTextContent("12 文法");
  });
});
