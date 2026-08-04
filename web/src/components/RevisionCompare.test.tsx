import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import { RevisionCompare } from "./RevisionCompare";

describe("RevisionCompare", () => {
  it("renders nothing when both sides empty", () => {
    const { container } = render(<RevisionCompare draftBody="" revisionBody="" />);
    expect(container).toBeEmptyDOMElement();
  });

  it("shows side-by-side bodies and char delta", () => {
    render(
      <RevisionCompare draftBody="短い。" revisionBody="少し長くしました。" title="対比" />
    );
    expect(screen.getByLabelText("対比")).toBeVisible();
    expect(screen.getByText("短い。")).toBeVisible();
    expect(screen.getByText("少し長くしました。")).toBeVisible();
    // 短い。 = 3 chars; 少し長くしました。 = 9 chars → +6
    expect(screen.getByText(/下書き 3 字/)).toBeVisible();
    expect(screen.getByText(/改稿 9 字/)).toBeVisible();
    expect(screen.getByText(/\+6 字/)).toBeVisible();
  });

  it("shows placeholder when only draft exists", () => {
    render(<RevisionCompare draftBody="下書きだけ" revisionBody="" />);
    expect(screen.getByText("下書きだけ")).toBeVisible();
    expect(screen.getByText("（まだありません）")).toBeVisible();
  });
});
