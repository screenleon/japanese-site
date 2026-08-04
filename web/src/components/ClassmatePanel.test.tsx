import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import type { KokugoClassmate } from "../kokugoTypes";
import { ClassmatePanel } from "./ClassmatePanel";

const samples: KokugoClassmate[] = [
  {
    id: "c1",
    name_ja: "田中さん",
    reveal_after: { kind: "task", task_id: "summary-1" },
    text_ja: "私は b を選びました。",
    focus_ja: "全体の流れ",
  },
  {
    id: "c2",
    name_ja: "鈴木さん",
    reveal_after: { kind: "artifact" },
    text_ja: "短い提案です。",
  },
];

describe("ClassmatePanel", () => {
  it("renders nothing when classmates empty", () => {
    /**
     * Behavior: empty classmate list produces no DOM (parent can mount unconditionally).
     * Steps:
     * 1. Arrange ClassmatePanel with classmates=[].
     * 2. Act render.
     * 3. Assert container is empty (no title, no list).
     */
    const { container } = render(<ClassmatePanel classmates={[]} />);
    expect(container).toBeEmptyDOMElement();
  });

  it("renders names, focus, and bodies", () => {
    /**
     * Behavior: curated samples show name, optional focus label, and body text.
     * Steps:
     * 1. Arrange two classmates (with and without focus_ja) and a custom title.
     * 2. Act render ClassmatePanel.
     * 3. Assert section aria-label, both names, focus line, and both bodies.
     */
    render(<ClassmatePanel classmates={samples} title="クラスメイトの答え" />);
    expect(screen.getByLabelText("クラスメイトの答え")).toBeVisible();
    expect(screen.getByText("田中さん")).toBeVisible();
    expect(screen.getByText("視点: 全体の流れ")).toBeVisible();
    expect(screen.getByText("私は b を選びました。")).toBeVisible();
    expect(screen.getByText("鈴木さん")).toBeVisible();
    expect(screen.getByText("短い提案です。")).toBeVisible();
  });
});
