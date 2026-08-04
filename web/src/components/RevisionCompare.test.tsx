import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import { RevisionCompare } from "./RevisionCompare";

describe("RevisionCompare", () => {
  it("renders nothing when both sides empty", () => {
    /**
     * Behavior: blank draft and revision yield no compare UI.
     * Steps:
     * 1. Arrange empty draftBody and revisionBody.
     * 2. Act render RevisionCompare.
     * 3. Assert empty DOM (no labels, no counters).
     */
    const { container } = render(<RevisionCompare draftBody="" revisionBody="" />);
    expect(container).toBeEmptyDOMElement();
  });

  it("shows side-by-side bodies and char delta", () => {
    /**
     * Behavior: non-empty draft and revision render side-by-side with char counts and delta.
     * Steps:
     * 1. Arrange short draft (3 字) and longer revision (9 字) with title 対比.
     * 2. Act render RevisionCompare.
     * 3. Assert both bodies, 下書き/改稿 counters, and +6 字 delta.
     */
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
    /**
     * Behavior: draft-only compare still shows revision column placeholder.
     * Steps:
     * 1. Arrange non-empty draftBody and empty revisionBody.
     * 2. Act render RevisionCompare.
     * 3. Assert draft text and 「まだありません」 on the revision side.
     */
    render(<RevisionCompare draftBody="下書きだけ" revisionBody="" />);
    expect(screen.getByText("下書きだけ")).toBeVisible();
    expect(screen.getByText("（まだありません）")).toBeVisible();
  });
});
