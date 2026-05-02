import { beforeEach, describe, expect, it, vi } from "vitest";
import { isApiError } from "./api";
import { staticApi } from "./staticApi";

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

  it("randomGrammar uses the level index then fetches detail", async () => {
    responses.set("/data/grammar/N1/_index.json", JSON.stringify(["aru-majiki"]));
    responses.set(
      "/data/grammar/N1/aru-majiki.json",
      JSON.stringify({
        slug: "aru-majiki",
        title_ja: "〜まじき",
        title_zh: "〜まじき",
        jlpt_level: "N1",
        explanation_zh: "不該有的",
      })
    );

    await expect(staticApi.randomGrammar("N1")).resolves.toMatchObject({
      slug: "aru-majiki",
      jlpt_level: "N1",
    });
    expect(fetch).toHaveBeenNthCalledWith(1, "/data/grammar/N1/_index.json");
    expect(fetch).toHaveBeenNthCalledWith(2, "/data/grammar/N1/aru-majiki.json");
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

  it("rejects unsupported quiz and sentence methods in static mode", async () => {
    await expect(staticApi.randomSentence()).rejects.toSatisfy((error) =>
      isApiError(error, "unsupported_in_static_mode")
    );
    await expect(staticApi.nextQuestion()).rejects.toSatisfy((error) =>
      isApiError(error, "unsupported_in_static_mode")
    );
  });
});
