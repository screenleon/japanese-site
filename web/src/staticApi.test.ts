import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { ApiError, isApiError } from "./api";
import { resetStaticApiCachesForTest, staticApi, staticApiTestHooks } from "./staticApi";

const responses = new Map<string, string>();

function jsonResponse(body: unknown) {
  return new Response(JSON.stringify(body), {
    status: 200,
    headers: { "Content-Type": "application/json" },
  });
}

function textResponse(body: string) {
  return new Response(body, {
    status: 200,
    headers: { "Content-Type": "application/jsonl" },
  });
}

function grammarPoint(slug: string, level: string) {
  return {
    slug,
    title_ja: slug,
    title_zh: slug,
    jlpt_level: level,
    explanation_zh: slug,
  };
}

describe("staticApi", () => {
  beforeEach(() => {
    resetStaticApiCachesForTest();
    responses.clear();
    vi.stubGlobal(
      "fetch",
      vi.fn((url: string) => {
        const body = responses.get(url);
        if (body === undefined) {
          return Promise.resolve(new Response("", { status: 404, statusText: "Not Found" }));
        }
        if (url.endsWith(".jsonl")) return Promise.resolve(textResponse(body));
        return Promise.resolve(jsonResponse(JSON.parse(body)));
      })
    );
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it("returns static capabilities with quiz and sentence disabled", async () => {
    await expect(staticApi.getCapabilities()).resolves.toEqual({
      progress: false,
      history: false,
      quiz: false,
      sentence: false,
    });
  });

  it("markRead is a no-op", async () => {
    await expect(
      staticApi.markRead({ type: "grammar", slug: "aru-majiki" })
    ).resolves.toBeUndefined();
    expect(fetch).not.toHaveBeenCalled();
  });

  it("getProgress returns zero progress", async () => {
    await expect(staticApi.getProgress("grammar", "N1")).resolves.toEqual({
      level: "N1",
      content_type: "grammar",
      read: 0,
      total: 0,
      percent: 0,
    });
  });

  it("randomGrammar loads the level rollup and picks a record", async () => {
    responses.set(
      "/data/grammar/N1.json",
      JSON.stringify([
        {
          slug: "aru-majiki",
          title_ja: "〜まじき",
          title_zh: "〜まじき",
          jlpt_level: "N1",
          explanation_zh: "不該有的",
        },
      ])
    );

    await expect(staticApi.randomGrammar("N1")).resolves.toMatchObject({
      slug: "aru-majiki",
      title_ja: "〜まじき",
      jlpt_level: "N1",
    });
    expect(fetch).toHaveBeenCalledTimes(1);
    expect(fetch).toHaveBeenCalledWith("/data/grammar/N1.json");
  });

  it("retries fetch after a transient failure", async () => {
    vi.stubGlobal(
      "fetch",
      vi
        .fn()
        .mockResolvedValueOnce(
          new Response("", { status: 503, statusText: "Service Unavailable" })
        )
        .mockResolvedValueOnce(
          jsonResponse([
            {
              slug: "aru-majiki",
              title_ja: "〜まじき",
              title_zh: "〜まじき",
              jlpt_level: "N1",
              explanation_zh: "不該有的",
            },
          ])
        )
    );

    await expect(staticApi.randomGrammar("N1")).rejects.toSatisfy((error) => {
      return error instanceof ApiError && error.status === 503;
    });
    await expect(staticApi.randomGrammar("N1")).resolves.toMatchObject({
      slug: "aru-majiki",
    });
    expect(fetch).toHaveBeenCalledTimes(2);
  });

  it("listGrammar returns full GrammarPoint records (not slug stubs)", async () => {
    responses.set(
      "/data/grammar/N1.json",
      JSON.stringify([
        {
          slug: "aru-majiki",
          title_ja: "〜まじき",
          title_zh: "〜まじき",
          jlpt_level: "N1",
          explanation_zh: "不該有的",
        },
      ])
    );

    const result = await staticApi.listGrammar("N1");
    expect(result.count).toBe(1);
    expect(result.points[0]).toMatchObject({
      slug: "aru-majiki",
      title_ja: "〜まじき",
      title_zh: "〜まじき",
      explanation_zh: "不該有的",
    });
  });

  it("getKanji finds a character across levels", async () => {
    responses.set("/data/kanji/N1.jsonl", "");
    responses.set("/data/kanji/N2.jsonl", "");
    responses.set("/data/kanji/N3.jsonl", "");
    responses.set(
      "/data/kanji/N4.jsonl",
      '{"character":"会","meaning_zh":"會","source":"curated","license":"CC-BY-SA-4.0"}'
    );
    responses.set("/data/kanji/N5.jsonl", "");

    await expect(staticApi.getKanji("会")).resolves.toMatchObject({
      character: "会",
      jlpt_level: "N4",
    });
  });

  it("getKanji throws not_found when character is missing in every level", async () => {
    responses.set("/data/kanji/N1.jsonl", "");
    responses.set("/data/kanji/N2.jsonl", "");
    responses.set(
      "/data/kanji/N3.jsonl",
      '{"character":"会","meaning_zh":"會","source":"curated","license":"CC-BY-SA-4.0"}'
    );
    responses.set("/data/kanji/N4.jsonl", "");
    responses.set("/data/kanji/N5.jsonl", "");

    await expect(staticApi.getKanji("食")).rejects.toSatisfy((error) =>
      isApiError(error, "not_found")
    );
  });

  it("getGrammar throws not_found when slug is missing in every level rollup", async () => {
    responses.set("/data/grammar/N1.json", "[]");
    responses.set("/data/grammar/N2.json", "[]");
    responses.set("/data/grammar/N3.json", "[]");
    responses.set("/data/grammar/N4.json", "[]");
    responses.set("/data/grammar/N5.json", "[]");

    await expect(staticApi.getGrammar("nonexistent-slug")).rejects.toSatisfy(
      (error) => isApiError(error, "not_found")
    );
  });

  it("choose throws when level rollup is empty", async () => {
    responses.set("/data/grammar/N1.json", "[]");

    await expect(staticApi.randomGrammar("N1")).rejects.toSatisfy((error) =>
      isApiError(error, "not_found")
    );
  });

  it("fetchJSON propagates HTTP error as ApiError", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn(() =>
        Promise.resolve(new Response("", { status: 404, statusText: "Not Found" }))
      )
    );

    await expect(staticApi.randomGrammar("N1")).rejects.toSatisfy((error) => {
      return (
        error instanceof ApiError &&
        error.status === 404 &&
        error.code === "not_found"
      );
    });
  });

  it("fetchJSON propagates 500 HTTP error with http_error code", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn(() =>
        Promise.resolve(new Response("", { status: 500, statusText: "Server Error" }))
      )
    );

    await expect(staticApi.randomGrammar("N1")).rejects.toSatisfy((error) => {
      return (
        error instanceof ApiError &&
        error.status === 500 &&
        error.code === "http_error"
      );
    });
  });

  it("fetchJSON wraps malformed JSON parse failures as ApiError", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn(() =>
        Promise.resolve(
          new Response("{not-json}", {
            status: 200,
            headers: { "Content-Type": "application/json" },
          })
        )
      )
    );

    await expect(staticApi.randomGrammar("N1")).rejects.toSatisfy((error) => {
      return (
        error instanceof ApiError &&
        error.status === 200 &&
        error.code === "parse_error"
      );
    });
  });

  it("rejects unsupported quiz and sentence methods in static mode", async () => {
    await expect(staticApi.randomSentence()).rejects.toSatisfy((error) =>
      isApiError(error, "unsupported_in_static_mode")
    );
    await expect(staticApi.nextQuestion()).rejects.toSatisfy((error) =>
      isApiError(error, "unsupported_in_static_mode")
    );
  });

  it("getGrammarExamples reads flat jsonl and returns compact examples", async () => {
    responses.set(
      "/data/grammar-examples/sae.jsonl",
      [
        '{"id":1,"text_ja":"名前さえ書けばいい。","text_zh":"只要寫名字就好。"}',
        '{"id":2,"text_ja":"水さえあれば足りる。","text_zh":"只要有水就夠了。"}',
      ].join("\n")
    );

    const result = await staticApi.getGrammarExamples?.("sae");

    expect(result).toEqual({
      count: 2,
      examples: [
        { id: 1, text_ja: "名前さえ書けばいい。", text_zh: "只要寫名字就好。" },
        { id: 2, text_ja: "水さえあれば足りる。", text_zh: "只要有水就夠了。" },
      ],
    });
  });

  it("getGrammarExamples encodes the slug for the URL path", async () => {
    // Defense-in-depth: lint-grammar enforces /^[a-z0-9-]+$/ slug shape, so
    // production slugs need no encoding. This test pins the invariant that
    // staticApi still encodes per JS-039 (parity with httpApi).
    const requested: string[] = [];
    const origFetch = globalThis.fetch;
    globalThis.fetch = ((url: string) => {
      requested.push(url);
      return Promise.resolve(new Response("", { status: 404 }));
    }) as typeof fetch;
    try {
      await staticApi.getGrammarExamples?.("foo/bar");
      expect(requested).toEqual(["/data/grammar-examples/foo%2Fbar.jsonl"]);
    } finally {
      globalThis.fetch = origFetch;
    }
  });

  it("getGrammarExamples returns an empty result on 404", async () => {
    await expect(staticApi.getGrammarExamples?.("missing")).resolves.toEqual({
      examples: [],
      count: 0,
    });
  });

  it("getGrammarExamples propagates non-404 HTTP errors as ApiError", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn(() =>
        Promise.resolve(new Response("", { status: 500, statusText: "Server Error" }))
      )
    );

    await expect(staticApi.getGrammarExamples?.("missing")).rejects.toSatisfy(
      (error) => {
        return error instanceof ApiError && error.status === 500;
      }
    );
  });

  it("getGrammarExamples wraps malformed JSONL parse failures as ApiError", async () => {
    responses.set("/data/grammar-examples/bad.jsonl", "{not-json}");

    await expect(staticApi.getGrammarExamples?.("bad")).rejects.toSatisfy((error) => {
      return (
        error instanceof ApiError &&
        error.status === 200 &&
        error.code === "parse_error"
      );
    });
  });

  it("listGrammar propagates non-404 level rollup errors as ApiError", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn(() =>
        Promise.resolve(new Response("", { status: 500, statusText: "Server Error" }))
      )
    );

    await expect(staticApi.listGrammar("N1")).rejects.toSatisfy((error) => {
      return error instanceof ApiError && error.status === 500;
    });
  });

  it("listGrammar rollup swallows a 404 level and returns available levels", async () => {
    responses.set("/data/grammar/N1.json", JSON.stringify([grammarPoint("n1-a", "N1")]));
    responses.set("/data/grammar/N3.json", JSON.stringify([grammarPoint("n3-a", "N3")]));
    responses.set("/data/grammar/N4.json", JSON.stringify([grammarPoint("n4-a", "N4")]));
    responses.set("/data/grammar/N5.json", JSON.stringify([grammarPoint("n5-a", "N5")]));

    const result = await staticApi.listGrammar();

    expect(result.points.map((point) => point.slug)).toEqual([
      "n1-a",
      "n3-a",
      "n4-a",
      "n5-a",
    ]);
  });

  it("listGrammar rollup warns and degrades when one level returns 500", async () => {
    const warn = vi.spyOn(console, "warn").mockImplementation(() => {});
    vi.stubGlobal(
      "fetch",
      vi.fn((url: string) => {
        if (url === "/data/grammar/N2.json") {
          return Promise.resolve(
            new Response("", { status: 500, statusText: "Server Error" })
          );
        }
        const level = url.match(/\/data\/grammar\/(N[1-5])\.json$/)?.[1] ?? "N1";
        return Promise.resolve(jsonResponse([grammarPoint(`${level.toLowerCase()}-a`, level)]));
      })
    );

    const result = await staticApi.listGrammar();

    expect(result.points.map((point) => point.slug)).toEqual([
      "n1-a",
      "n3-a",
      "n4-a",
      "n5-a",
    ]);
    expect(warn).toHaveBeenCalledTimes(1);
    expect(warn).toHaveBeenCalledWith("grammar level rollup fetch failed", {
      level: "N2",
      err: expect.objectContaining({ status: 500, code: "http_error" }),
    });
  });

  it("uses a URL-only JSON cache key across different on404 policies", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn(() =>
        Promise.resolve(new Response("", { status: 404, statusText: "Not Found" }))
      )
    );

    await expect(
      staticApiTestHooks.fetchJSON<unknown[]>("/same-url.json", {
        on404: "empty-array",
      })
    ).resolves.toEqual([]);
    await expect(
      staticApiTestHooks.fetchJSON<unknown[]>("/same-url.json", {
        on404: "throw",
      })
    ).rejects.toSatisfy((error) => {
      return error instanceof ApiError && error.status === 404;
    });
    expect(fetch).toHaveBeenCalledTimes(1);
  });


  it("getGrammarExamples requests the flat slug path", async () => {
    responses.set(
      "/data/grammar-examples/sae.jsonl",
      '{"id":1,"text_ja":"名前さえ書けばいい。","text_zh":"只要寫名字就好。"}'
    );

    await staticApi.getGrammarExamples?.("sae");

    expect(fetch).toHaveBeenCalledWith("/data/grammar-examples/sae.jsonl");
  });

  it("getGrammarExamples deduplicates concurrent reads for the same slug", async () => {
    responses.set(
      "/data/grammar-examples/sae.jsonl",
      '{"id":1,"text_ja":"名前さえ書けばいい。","text_zh":"只要寫名字就好。"}'
    );

    await Promise.all([
      staticApi.getGrammarExamples?.("sae"),
      staticApi.getGrammarExamples?.("sae"),
    ]);

    expect(fetch).toHaveBeenCalledTimes(1);
  });
});
