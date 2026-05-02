import { renderHook, waitFor } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { api } from "../api";
import { useReadTracking } from "./useReadTracking";

const capabilitiesState = vi.hoisted(() => ({
  current: {
    progress: true,
    history: false,
    loaded: true,
    progressRevision: 0,
    bumpProgress: vi.fn(),
  },
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
    capabilitiesState.current = {
      progress: true,
      history: false,
      loaded: true,
      progressRevision: 0,
      bumpProgress: vi.fn(),
    };
    markRead.mockResolvedValue(undefined);
    warnSpy = vi.spyOn(console, "warn").mockImplementation(() => {});
  });

  afterEach(() => {
    vi.clearAllMocks();
    warnSpy.mockRestore();
  });

  it("fires markRead exactly once when slug becomes non-empty", () => {
    renderHook(
      ({ slug }) => useReadTracking(slug ? { type: "grammar", slug } : null),
      {
        initialProps: { slug: "foo" },
      }
    );

    expect(markRead).toHaveBeenCalledTimes(1);
    expect(markRead).toHaveBeenCalledWith({ type: "grammar", slug: "foo" });
  });

  it("does not re-fire on re-render with same (type, slug)", () => {
    const { rerender } = renderHook(
      ({ slug }) => useReadTracking(slug ? { type: "grammar", slug } : null),
      {
        initialProps: { slug: "foo" },
      }
    );

    capabilitiesState.current = {
      progress: true,
      history: false,
      loaded: false,
      progressRevision: 0,
      bumpProgress: vi.fn(),
    };
    rerender({ slug: "foo" });
    capabilitiesState.current = {
      progress: true,
      history: false,
      loaded: true,
      progressRevision: 0,
      bumpProgress: vi.fn(),
    };
    rerender({ slug: "foo" });

    expect(markRead).toHaveBeenCalledTimes(1);
  });

  it("fires again when slug changes", () => {
    const { rerender } = renderHook(
      ({ slug }) => useReadTracking(slug ? { type: "grammar", slug } : null),
      {
        initialProps: { slug: "foo" },
      }
    );

    rerender({ slug: "bar" });

    expect(markRead).toHaveBeenCalledTimes(2);
    expect(markRead).toHaveBeenLastCalledWith({ type: "grammar", slug: "bar" });
  });

  it("no-op when slug is empty / undefined / whitespace", () => {
    const { rerender } = renderHook(
      ({ slug }) => useReadTracking(slug ? { type: "grammar", slug } : null),
      {
        initialProps: { slug: "" as string | undefined },
      }
    );

    rerender({ slug: undefined });
    rerender({ slug: "   " });

    expect(markRead).not.toHaveBeenCalled();
  });

  it("does not fire when capabilities not yet loaded on mount", () => {
    capabilitiesState.current = {
      progress: false,
      history: false,
      loaded: false,
      progressRevision: 0,
      bumpProgress: vi.fn(),
    };

    renderHook(() => useReadTracking({ type: "grammar", slug: "foo" }));

    expect(markRead).not.toHaveBeenCalled();
  });

  it("does not fire in null mode (loaded but progress=false)", () => {
    capabilitiesState.current = {
      progress: false,
      history: false,
      loaded: true,
      progressRevision: 0,
      bumpProgress: vi.fn(),
    };

    renderHook(() => useReadTracking({ type: "grammar", slug: "foo" }));

    expect(markRead).not.toHaveBeenCalled();
  });

  it("fires after capabilities transition from loading to progress=true", () => {
    capabilitiesState.current = {
      progress: false,
      history: false,
      loaded: false,
      progressRevision: 0,
      bumpProgress: vi.fn(),
    };

    const { rerender } = renderHook(() =>
      useReadTracking({ type: "grammar", slug: "foo" })
    );
    expect(markRead).not.toHaveBeenCalled();

    capabilitiesState.current = {
      progress: true,
      history: false,
      loaded: true,
      progressRevision: 0,
      bumpProgress: vi.fn(),
    };
    rerender();

    expect(markRead).toHaveBeenCalledTimes(1);
    expect(markRead).toHaveBeenCalledWith({ type: "grammar", slug: "foo" });
  });

  it("swallows api.markRead rejection without throwing or re-firing", async () => {
    markRead.mockRejectedValueOnce(new Error("network failed"));

    const { rerender } = renderHook(
      ({ slug }) => useReadTracking(slug ? { type: "grammar", slug } : null),
      {
        initialProps: { slug: "foo" },
      }
    );

    await waitFor(() => {
      expect(warnSpy).toHaveBeenCalledWith(
        "failed to mark content as read",
        expect.any(Error)
      );
    });

    rerender({ slug: "foo" });

    expect(markRead).toHaveBeenCalledTimes(1);
  });

  it("calls bumpProgress after successful markRead", async () => {
    const bumpProgress = vi.fn();
    capabilitiesState.current = {
      progress: true,
      history: false,
      loaded: true,
      progressRevision: 0,
      bumpProgress,
    };

    renderHook(() => useReadTracking({ type: "grammar", slug: "foo" }));

    await waitFor(() => {
      expect(bumpProgress).toHaveBeenCalledTimes(1);
    });
  });
});
