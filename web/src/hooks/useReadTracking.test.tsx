import { renderHook, waitFor } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { api } from "../api";
import { useReadTracking } from "./useReadTracking";

const capabilitiesState = vi.hoisted(() => ({
  current: { progress: true, history: false, loaded: true },
}));

// These tests exercise the hook in isolation, so useCapabilities is mocked
// directly instead of wrapping every renderHook call in CapabilitiesProvider.
vi.mock("../capabilities", () => ({
  useCapabilities: () => capabilitiesState.current,
}));

vi.mock("../api", () => ({
  api: {
    markRead: vi.fn(),
  },
}));

const markRead = vi.mocked(api.markRead);

describe("useReadTracking", () => {
  let warnSpy: ReturnType<typeof vi.spyOn>;

  beforeEach(() => {
    capabilitiesState.current = { progress: true, history: false, loaded: true };
    markRead.mockResolvedValue(undefined);
    warnSpy = vi.spyOn(console, "warn").mockImplementation(() => {});
  });

  afterEach(() => {
    vi.clearAllMocks();
    warnSpy.mockRestore();
  });

  it("fires markRead exactly once when slug becomes non-empty", () => {
    renderHook(({ slug }) => useReadTracking("grammar", slug), {
      initialProps: { slug: "foo" },
    });

    expect(markRead).toHaveBeenCalledTimes(1);
    expect(markRead).toHaveBeenCalledWith("grammar", "foo");
  });

  it("does not re-fire on re-render with same (type, slug)", () => {
    const { rerender } = renderHook(({ slug }) => useReadTracking("grammar", slug), {
      initialProps: { slug: "foo" },
    });

    capabilitiesState.current = { progress: true, history: false, loaded: false };
    rerender({ slug: "foo" });
    capabilitiesState.current = { progress: true, history: false, loaded: true };
    rerender({ slug: "foo" });

    expect(markRead).toHaveBeenCalledTimes(1);
  });

  it("fires again when slug changes", () => {
    const { rerender } = renderHook(({ slug }) => useReadTracking("grammar", slug), {
      initialProps: { slug: "foo" },
    });

    rerender({ slug: "bar" });

    expect(markRead).toHaveBeenCalledTimes(2);
    expect(markRead).toHaveBeenLastCalledWith("grammar", "bar");
  });

  it("no-op when slug is empty / undefined / whitespace", () => {
    const { rerender } = renderHook(({ slug }) => useReadTracking("grammar", slug), {
      initialProps: { slug: "" as string | undefined },
    });

    rerender({ slug: undefined });
    rerender({ slug: "   " });

    expect(markRead).not.toHaveBeenCalled();
  });

  it("swallows api.markRead rejection without throwing or re-firing", async () => {
    markRead.mockRejectedValueOnce(new Error("network failed"));

    const { rerender } = renderHook(({ slug }) => useReadTracking("grammar", slug), {
      initialProps: { slug: "foo" },
    });

    await waitFor(() => {
      expect(warnSpy).toHaveBeenCalledWith(
        "failed to mark content as read",
        expect.any(Error)
      );
    });

    rerender({ slug: "foo" });

    expect(markRead).toHaveBeenCalledTimes(1);
  });
});
