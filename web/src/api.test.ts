import { afterEach, describe, expect, it, vi } from "vitest";

describe("api transport switch", () => {
  afterEach(() => {
    vi.unstubAllEnvs();
    vi.unstubAllGlobals();
    vi.resetModules();
  });

  it("uses httpApi by default", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn(() =>
        Promise.resolve(
          new Response(
            JSON.stringify({
              progress: true,
              history: false,
              quiz: true,
              sentence: true,
            }),
            { status: 200 }
          )
        )
      )
    );

    vi.resetModules();
    const { api } = await import("./api");

    await expect(api.getCapabilities()).resolves.toEqual({
      progress: true,
      history: false,
      quiz: true,
      sentence: true,
    });
    expect(fetch).toHaveBeenCalledWith("/api/capabilities");
  });

  it("uses staticApi when VITE_DEPLOY_MODE=static", async () => {
    vi.stubEnv("VITE_DEPLOY_MODE", "static");
    vi.stubGlobal("fetch", vi.fn());

    vi.resetModules();
    const { api } = await import("./api");

    await expect(api.getCapabilities()).resolves.toEqual({
      progress: false,
      history: false,
      quiz: false,
      sentence: false,
    });
    expect(fetch).not.toHaveBeenCalled();
  });
});
