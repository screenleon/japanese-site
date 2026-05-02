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
  return <div data-testid="capabilities">{JSON.stringify(capabilities)}</div>;
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

  it("resolves to {progress, history, loaded:true} on successful fetch", async () => {
    getCapabilities.mockResolvedValueOnce({ progress: true, history: false });

    render(
      <CapabilitiesProvider>
        <Consumer />
      </CapabilitiesProvider>
    );

    await waitFor(() => {
      expect(screen.getByTestId("capabilities")).toHaveTextContent(
        JSON.stringify({ progress: true, history: false, loaded: true })
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
        JSON.stringify({ progress: false, history: false, loaded: true })
      );
    });
  });
});
