import { fireEvent, render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import {
  ContentFirstLayout,
  formatDirectorySummary,
} from "./ContentFirstLayout";

describe("formatDirectorySummary", () => {
  /**
   * Verifies directory toggle labels include optional item counts.
   * Steps:
   * 1. Call formatDirectorySummary with positive, zero, and omitted counts.
   * 2. Assert the exact Chinese summary strings for each input.
   */
  it("formats count and bare label", () => {
    expect(formatDirectorySummary(3)).toBe("本級目錄 (3)");
    expect(formatDirectorySummary(0)).toBe("本級目錄 (0)");
    expect(formatDirectorySummary()).toBe("本級目錄");
  });
});

describe("ContentFirstLayout", () => {
  /**
   * Verifies content precedes directory in document order (mobile content-first).
   * Steps:
   * 1. Render layout with content and directory children.
   * 2. Compare document positions of the two panes.
   * 3. Assert directory follows content (DOCUMENT_POSITION_FOLLOWING bit equals the following mask).
   */
  it("places content before directory in DOM order for mobile content-first reading", () => {
    render(
      <ContentFirstLayout
        content={<article>active entry</article>}
        directory={
          <ul>
            <li>item a</li>
          </ul>
        }
        directorySummary={formatDirectorySummary(1)}
        directoryOpen={false}
        onDirectoryOpenChange={() => {}}
      />
    );

    const content = screen.getByTestId("study-content");
    const directory = screen.getByTestId("study-directory");
    const following = Node.DOCUMENT_POSITION_FOLLOWING;
    expect(content.compareDocumentPosition(directory) & following).toBe(following);
  });

  /**
   * Verifies the controlled mobile directory panel starts hidden and opens via callback.
   * Steps:
   * 1. Render with directoryOpen=false and assert panel classes and aria-expanded=false.
   * 2. Click the summary toggle and assert onDirectoryOpenChange(true).
   * 3. Rerender with directoryOpen=true and assert expanded attributes and block class.
   */
  it("hides the directory panel when closed (mobile), shows when open", () => {
    const onChange = vi.fn();
    const { rerender } = render(
      <ContentFirstLayout
        content={<div>body</div>}
        directory={<div>directory body</div>}
        directorySummary={formatDirectorySummary(3)}
        directoryOpen={false}
        onDirectoryOpenChange={onChange}
      />
    );

    const panel = document.getElementById("study-directory-panel");
    expect(panel).toHaveClass("hidden");
    expect(panel).toHaveClass("lg:block");

    const toggle = screen.getByRole("button", { name: /本級目錄 \(3\)/ });
    expect(toggle).toHaveAttribute("aria-expanded", "false");
    fireEvent.click(toggle);
    expect(onChange).toHaveBeenCalledWith(true);

    rerender(
      <ContentFirstLayout
        content={<div>body</div>}
        directory={<div>directory body</div>}
        directorySummary={formatDirectorySummary(3)}
        directoryOpen={true}
        onDirectoryOpenChange={onChange}
      />
    );
    expect(screen.getByRole("button", { name: /本級目錄 \(3\)/ })).toHaveAttribute(
      "aria-expanded",
      "true"
    );
    expect(document.getElementById("study-directory-panel")).toHaveClass("block");
  });
});
