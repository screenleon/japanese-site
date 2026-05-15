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

  /**
   * Verifies the provider uses the restrictive default when storage has no saved preference.
   * Steps:
   * 1. Arrange by leaving localStorage empty.
   * 2. Act by rendering a consumer inside ChineseVisibilityProvider.
   * 3. Assert the visible state is false.
   */
  it("defaults visible to false when localStorage has no value", () => {
    // Arrange / Act
    render(
      <ChineseVisibilityProvider>
        <Consumer />
      </ChineseVisibilityProvider>
    );

    // Assert
    expect(screen.getByTestId("visible")).toHaveTextContent("false");
  });

  /**
   * Verifies the provider restores a saved visible=true preference on mount.
   * Steps:
   * 1. Arrange by writing true to the provider's localStorage key.
   * 2. Act by rendering a consumer inside ChineseVisibilityProvider.
   * 3. Assert the visible state is true.
   */
  it("reads true from localStorage on mount", () => {
    // Arrange
    window.localStorage.setItem(STORAGE_KEY, "true");

    // Act
    render(
      <ChineseVisibilityProvider>
        <Consumer />
      </ChineseVisibilityProvider>
    );

    // Assert
    expect(screen.getByTestId("visible")).toHaveTextContent("true");
  });

  /**
   * Verifies toggle mutates visibility in both directions.
   * Steps:
   * 1. Arrange a consumer inside the provider with the default hidden state.
   * 2. Act by clicking toggle twice.
   * 3. Assert the first click shows Chinese content and the second hides it again.
   */
  it("toggle flips the state and a second toggle returns to the original value", () => {
    // Arrange
    render(
      <ChineseVisibilityProvider>
        <Consumer />
      </ChineseVisibilityProvider>
    );

    // Act / Assert
    fireEvent.click(screen.getByRole("button", { name: "toggle" }));
    expect(screen.getByTestId("visible")).toHaveTextContent("true");

    fireEvent.click(screen.getByRole("button", { name: "toggle" }));
    expect(screen.getByTestId("visible")).toHaveTextContent("false");
  });

  /**
   * Verifies explicit setVisible calls are persisted to localStorage.
   * Steps:
   * 1. Arrange a consumer inside the provider.
   * 2. Act by clicking the true and false setter buttons.
   * 3. Assert localStorage records each requested boolean value.
   */
  it("setVisible writes true and false to localStorage", async () => {
    // Arrange
    render(
      <ChineseVisibilityProvider>
        <Consumer />
      </ChineseVisibilityProvider>
    );

    // Act / Assert
    fireEvent.click(screen.getByRole("button", { name: "set true" }));
    await waitFor(() => {
      expect(window.localStorage.getItem(STORAGE_KEY)).toBe("true");
    });

    fireEvent.click(screen.getByRole("button", { name: "set false" }));
    await waitFor(() => {
      expect(window.localStorage.getItem(STORAGE_KEY)).toBe("false");
    });
  });

  /**
   * Verifies the shared toggle button reflects provider state through aria-pressed.
   * Steps:
   * 1. Arrange ChineseVisibilityToggle inside the provider.
   * 2. Act by clicking the toggle button.
   * 3. Assert aria-pressed changes from false to true.
   */
  it("ChineseVisibilityToggle button click flips state", () => {
    // Arrange
    render(
      <ChineseVisibilityProvider>
        <ChineseVisibilityToggle />
      </ChineseVisibilityProvider>
    );

    // Act / Assert
    const button = screen.getByRole("button", { name: "中文" });
    expect(button).toHaveAttribute("aria-pressed", "false");

    fireEvent.click(button);
    expect(button).toHaveAttribute("aria-pressed", "true");
  });

  /**
   * Verifies IfChinese gates children based on the provider's visible state.
   * Steps:
   * 1. Arrange IfChinese under a hidden provider.
   * 2. Act by rerendering the same child under a visible provider.
   * 3. Assert the child is absent when hidden and visible when enabled.
   */
  it("IfChinese renders children only when visible", () => {
    // Arrange / Act
    const { rerender } = render(
      <ChineseVisibilityProvider key="hidden" initialVisible={false}>
        <IfChinese>中文內容</IfChinese>
      </ChineseVisibilityProvider>
    );

    // Assert
    expect(screen.queryByText("中文內容")).not.toBeInTheDocument();

    // Act
    rerender(
      <ChineseVisibilityProvider key="visible" initialVisible={true}>
        <IfChinese>中文內容</IfChinese>
      </ChineseVisibilityProvider>
    );

    // Assert
    expect(screen.getByText("中文內容")).toBeVisible();
  });

  /**
   * Verifies provider state still updates when localStorage persistence throws.
   * Steps:
   * 1. Arrange Storage.setItem to throw a quota-style error.
   * 2. Act by rendering a consumer and clicking toggle.
   * 3. Assert visible state changes to true despite the storage failure.
   */
  it("tolerates localStorage write failure", () => {
    // Arrange
    vi.spyOn(Storage.prototype, "setItem").mockImplementation(() => {
      throw new Error("quota exceeded");
    });

    // Act
    render(
      <ChineseVisibilityProvider>
        <Consumer />
      </ChineseVisibilityProvider>
    );

    fireEvent.click(screen.getByRole("button", { name: "toggle" }));

    // Assert
    expect(screen.getByTestId("visible")).toHaveTextContent("true");
  });
});
