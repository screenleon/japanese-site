import { render, screen, waitFor } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { api } from "./api";
import { CapabilitiesProvider, useCapabilities } from "./capabilities";

vi.mock("./api", () => ({
  api: {
    getCapabilities: vi.fn(),
  },
}));

const getCapabilities = vi.mocked(api.getCapabilities);

function Consumer() {
  const capabilities = useCapabilities();
  const { bumpProgress, ...serializable } = capabilities;
  return (
    <div data-testid="capabilities">
      {JSON.stringify({ ...serializable, bumpProgress: typeof bumpProgress })}
    </div>
  );
}

describe("CapabilitiesProvider", () => {
  let warnSpy: ReturnType<typeof vi.spyOn>;

  beforeEach(() => {
    warnSpy = vi.spyOn(console, "warn").mockImplementation(() => {});
  });

  afterEach(() => {
    vi.clearAllMocks();
    warnSpy.mockRestore();
  });

  it("provides default null-mode state outside a provider", () => {
    render(<Consumer />);

    expect(screen.getByTestId("capabilities")).toHaveTextContent(
      JSON.stringify({
        progress: false,
        history: false,
        quiz: false,
        sentence: false,
        kokugo: false,
        loaded: false,
        progressRevisions: { grammar: 0, vocab: 0, kanji: 0 },
        bumpProgress: "function",
      })
    );
  });

  it("resolves to {progress, history, quiz, sentence, kokugo, loaded:true} on successful fetch", async () => {
    getCapabilities.mockResolvedValueOnce({
      progress: true,
      history: false,
      quiz: true,
      sentence: true,
      kokugo: true,
    });

    render(
      <CapabilitiesProvider>
        <Consumer />
      </CapabilitiesProvider>
    );

    await waitFor(() => {
      expect(screen.getByTestId("capabilities")).toHaveTextContent(
        JSON.stringify({
          progress: true,
          history: false,
          quiz: true,
          sentence: true,
          kokugo: true,
          loaded: true,
          progressRevisions: { grammar: 0, vocab: 0, kanji: 0 },
          bumpProgress: "function",
        })
      );
    });
  });

  it("falls back to null-mode on fetch error", async () => {
    getCapabilities.mockRejectedValueOnce(new Error("network failed"));

    render(
      <CapabilitiesProvider>
        <Consumer />
      </CapabilitiesProvider>
    );

    await waitFor(() => {
      expect(screen.getByTestId("capabilities")).toHaveTextContent(
        JSON.stringify({
          progress: false,
          history: false,
          quiz: false,
          sentence: false,
          kokugo: false,
          loaded: true,
          progressRevisions: { grammar: 0, vocab: 0, kanji: 0 },
          bumpProgress: "function",
        })
      );
    });
  });
});
