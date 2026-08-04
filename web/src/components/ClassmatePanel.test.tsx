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
    const { container } = render(<ClassmatePanel classmates={[]} />);
    expect(container).toBeEmptyDOMElement();
  });

  it("renders names, focus, and bodies", () => {
    render(<ClassmatePanel classmates={samples} title="クラスメイトの答え" />);
    expect(screen.getByLabelText("クラスメイトの答え")).toBeVisible();
    expect(screen.getByText("田中さん")).toBeVisible();
    expect(screen.getByText("視点: 全体の流れ")).toBeVisible();
    expect(screen.getByText("私は b を選びました。")).toBeVisible();
    expect(screen.getByText("鈴木さん")).toBeVisible();
    expect(screen.getByText("短い提案です。")).toBeVisible();
  });
});
