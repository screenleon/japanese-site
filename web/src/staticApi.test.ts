import { beforeEach, describe, expect, it, vi } from "vitest";
import { ApiError, isApiError } from "./api";
import { resetStaticApiCachesForTest, staticApi } from "./staticApi";

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

    await expect(staticApi.randomGrammar("N1")).rejects.toSatisfy((error) =>
      isApiError(error, "not_found")
    );
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

  it("rejects unsupported quiz and sentence methods in static mode", async () => {
    await expect(staticApi.randomSentence()).rejects.toSatisfy((error) =>
      isApiError(error, "unsupported_in_static_mode")
    );
    await expect(staticApi.nextQuestion()).rejects.toSatisfy((error) =>
      isApiError(error, "unsupported_in_static_mode")
    );
  });
});
