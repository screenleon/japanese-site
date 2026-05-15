import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { afterEach, describe, expect, it, vi } from "vitest";
import {
  ChineseVisibilityProvider,
  ChineseVisibilityToggle,
  IfChinese,
  useChineseVisibility,
} from "./chineseVisibility";

const STORAGE_KEY = "japanese-site:chineseVisible";

function Consumer() {
  const { visible, setVisible, toggle } = useChineseVisibility();
  return (
    <div>
      <div data-testid="visible">{String(visible)}</div>
      <button type="button" onClick={() => setVisible(true)}>
        set true
      </button>
      <button type="button" onClick={() => setVisible(false)}>
        set false
      </button>
      <button type="button" onClick={toggle}>
        toggle
      </button>
    </div>
  );
}

describe("ChineseVisibilityProvider", () => {
  afterEach(() => {
    vi.restoreAllMocks();
    window.localStorage.clear();
  });

  it("defaults visible to false when localStorage has no value", () => {
    render(
      <ChineseVisibilityProvider>
        <Consumer />
      </ChineseVisibilityProvider>
    );

    expect(screen.getByTestId("visible")).toHaveTextContent("false");
  });

  it("reads true from localStorage on mount", () => {
    window.localStorage.setItem(STORAGE_KEY, "true");

    render(
      <ChineseVisibilityProvider>
        <Consumer />
      </ChineseVisibilityProvider>
    );

    expect(screen.getByTestId("visible")).toHaveTextContent("true");
  });

  it("toggle flips the state and a second toggle returns to the original value", () => {
    render(
      <ChineseVisibilityProvider>
        <Consumer />
      </ChineseVisibilityProvider>
    );

    fireEvent.click(screen.getByRole("button", { name: "toggle" }));
    expect(screen.getByTestId("visible")).toHaveTextContent("true");

    fireEvent.click(screen.getByRole("button", { name: "toggle" }));
    expect(screen.getByTestId("visible")).toHaveTextContent("false");
  });

  it("setVisible writes true and false to localStorage", async () => {
    render(
      <ChineseVisibilityProvider>
        <Consumer />
      </ChineseVisibilityProvider>
    );

    fireEvent.click(screen.getByRole("button", { name: "set true" }));
    await waitFor(() => {
      expect(window.localStorage.getItem(STORAGE_KEY)).toBe("true");
    });

    fireEvent.click(screen.getByRole("button", { name: "set false" }));
    await waitFor(() => {
      expect(window.localStorage.getItem(STORAGE_KEY)).toBe("false");
    });
  });

  it("ChineseVisibilityToggle button click flips state", () => {
    render(
      <ChineseVisibilityProvider>
        <ChineseVisibilityToggle />
      </ChineseVisibilityProvider>
    );

    const button = screen.getByRole("button", { name: "中文" });
    expect(button).toHaveAttribute("aria-pressed", "false");

    fireEvent.click(button);
    expect(button).toHaveAttribute("aria-pressed", "true");
  });

  it("IfChinese renders children only when visible", () => {
    const { rerender } = render(
      <ChineseVisibilityProvider key="hidden" initialVisible={false}>
        <IfChinese>中文內容</IfChinese>
      </ChineseVisibilityProvider>
    );

    expect(screen.queryByText("中文內容")).not.toBeInTheDocument();

    rerender(
      <ChineseVisibilityProvider key="visible" initialVisible={true}>
        <IfChinese>中文內容</IfChinese>
      </ChineseVisibilityProvider>
    );

    expect(screen.getByText("中文內容")).toBeVisible();
  });

  it("tolerates localStorage write failure", () => {
    vi.spyOn(Storage.prototype, "setItem").mockImplementation(() => {
      throw new Error("quota exceeded");
    });

    render(
      <ChineseVisibilityProvider>
        <Consumer />
      </ChineseVisibilityProvider>
    );

    fireEvent.click(screen.getByRole("button", { name: "toggle" }));
    expect(screen.getByTestId("visible")).toHaveTextContent("true");
  });
});
