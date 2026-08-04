import { afterEach, describe, expect, it, vi } from "vitest";
import { ApiError } from "./api";

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
              kokugo: true,
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
      kokugo: true,
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
      kokugo: false,
    });
    expect(fetch).not.toHaveBeenCalled();
  });
});

describe("httpApi Kokugo transport", () => {
  afterEach(() => {
    vi.unstubAllEnvs();
    vi.unstubAllGlobals();
    vi.resetModules();
  });

  async function loadHttpApi() {
    vi.stubEnv("VITE_DEPLOY_MODE", "api");
    vi.resetModules();
    const mod = await import("./api");
    return mod.httpApi;
  }

  it("listKokugoUnits GETs /api/kokugo/units", async () => {
    const fetchMock = vi.fn(() =>
      Promise.resolve(new Response(JSON.stringify({ units: [], count: 0 }), { status: 200 }))
    );
    vi.stubGlobal("fetch", fetchMock);
    const httpApi = await loadHttpApi();

    await expect(httpApi.listKokugoUnits()).resolves.toEqual({ units: [], count: 0 });
    expect(fetchMock).toHaveBeenCalledWith("/api/kokugo/units");
  });

  it("getKokugoSkills GETs /api/kokugo/skills and decodes skill map payload", async () => {
    const payload = {
      skills: [
        {
          skill: "reading.summary",
          label_ja: "要約する",
          status: "weak",
          graded: 2,
          correct: 0,
          practiced: 2,
          units_touching: ["e5-6/library-use"],
        },
      ],
      review_queue: [
        {
          stage: "e5-6",
          unit_id: "library-use",
          title_ja: "学校の図書室をもっと使いやすくするには",
          genre: "expository",
          target_skills: ["reading.summary"],
          unit_completed: false,
        },
      ],
    };
    const fetchMock = vi.fn(() =>
      Promise.resolve(new Response(JSON.stringify(payload), { status: 200 }))
    );
    vi.stubGlobal("fetch", fetchMock);
    const httpApi = await loadHttpApi();

    await expect(httpApi.getKokugoSkills()).resolves.toEqual(payload);
    expect(fetchMock).toHaveBeenCalledTimes(1);
    expect(fetchMock).toHaveBeenCalledWith("/api/kokugo/skills");
  });

  it("getKokugoUnit encodes stage and id path segments", async () => {
    const fetchMock = vi.fn(() =>
      Promise.resolve(new Response(JSON.stringify({ id: "library-use" }), { status: 200 }))
    );
    vi.stubGlobal("fetch", fetchMock);
    const httpApi = await loadHttpApi();

    await httpApi.getKokugoUnit("e5-6", "library-use");
    expect(fetchMock).toHaveBeenCalledWith("/api/kokugo/units/e5-6/library-use");
  });

  it("getKokugoUnitState GETs progress path", async () => {
    const fetchMock = vi.fn(() =>
      Promise.resolve(
        new Response(JSON.stringify({ attempts: [], artifacts: [] }), { status: 200 })
      )
    );
    vi.stubGlobal("fetch", fetchMock);
    const httpApi = await loadHttpApi();

    await httpApi.getKokugoUnitState("e5-6", "library-use");
    expect(fetchMock).toHaveBeenCalledWith("/api/kokugo/progress/e5-6/library-use");
  });

  it("putKokugoProgress PUTs JSON body", async () => {
    const fetchMock = vi.fn(() =>
      Promise.resolve(
        new Response(
          JSON.stringify({
            unit_key: "e5-6/library-use",
            stage: "e5-6",
            unit_id: "library-use",
            status: "in_progress",
            step: "read",
            started_at: "",
            updated_at: "",
          }),
          { status: 200 }
        )
      )
    );
    vi.stubGlobal("fetch", fetchMock);
    const httpApi = await loadHttpApi();

    await httpApi.putKokugoProgress("e5-6", "library-use", { step: "read" });
    expect(fetchMock).toHaveBeenCalledWith("/api/kokugo/progress/e5-6/library-use", {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ step: "read" }),
    });
  });

  it("submitKokugoTask POSTs answer payload", async () => {
    const fetchMock = vi.fn(() =>
      Promise.resolve(
        new Response(
          JSON.stringify({
            attempt: { id: 1, unit_key: "k", task_id: "t", answer: {}, created_at: "" },
            grade: { explanation_ja: "ok" },
          }),
          { status: 200 }
        )
      )
    );
    vi.stubGlobal("fetch", fetchMock);
    const httpApi = await loadHttpApi();

    await httpApi.submitKokugoTask("e5-6", "library-use", "summary-1", { choice_id: "b" });
    expect(fetchMock).toHaveBeenCalledWith(
      "/api/kokugo/progress/e5-6/library-use/tasks/summary-1",
      {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ answer: { choice_id: "b" } }),
      }
    );
  });

  it("saveKokugoArtifact PUTs expected_version and maps HTTP errors", async () => {
    const fetchMock = vi.fn(() =>
      Promise.resolve(
        new Response(JSON.stringify({ error: "stale_write" }), {
          status: 409,
          statusText: "Conflict",
        })
      )
    );
    vi.stubGlobal("fetch", fetchMock);
    const httpApi = await loadHttpApi();

    await expect(
      httpApi.saveKokugoArtifact("e5-6", "library-use", {
        revision: 0,
        body: "x",
        checklist_checked: [true],
        expected_version: 2,
      })
    ).rejects.toMatchObject({ name: "ApiError", status: 409, code: "stale_write" } satisfies Partial<ApiError>);

    expect(fetchMock).toHaveBeenCalledWith(
      "/api/kokugo/progress/e5-6/library-use/artifact",
      {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          revision: 0,
          body: "x",
          checklist_checked: [true],
          expected_version: 2,
        }),
      }
    );
  });
});
