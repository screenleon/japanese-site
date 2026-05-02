import { createContext, type ReactNode, useContext, useEffect, useState } from "react";
import { api, type Capabilities } from "./api";

type CapabilitiesState = Capabilities & {
  loaded: boolean;
};

const defaultCapabilities: CapabilitiesState = {
  progress: false,
  history: false,
  loaded: false,
};

const CapabilitiesContext = createContext<CapabilitiesState>(defaultCapabilities);

export function CapabilitiesProvider({ children }: { children: ReactNode }) {
  const [capabilities, setCapabilities] = useState<CapabilitiesState>(defaultCapabilities);

  useEffect(() => {
    let cancelled = false;
    api
      .getCapabilities()
      .then((next) => {
        if (cancelled) return;
        setCapabilities({ ...next, loaded: true });
      })
      .catch((error) => {
        console.warn("failed to load capabilities", error);
        if (cancelled) return;
        setCapabilities({ progress: false, history: false, loaded: true });
      });
    return () => {
      cancelled = true;
    };
  }, []);

  return (
    <CapabilitiesContext.Provider value={capabilities}>
      {children}
    </CapabilitiesContext.Provider>
  );
}

export function useCapabilities() {
  return useContext(CapabilitiesContext);
}
