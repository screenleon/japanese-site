import {
  createContext,
  type ReactNode,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
} from "react";
import { api, type Capabilities } from "./api";

type CapabilitiesState = Capabilities & {
  loaded: boolean;
  progressRevision: number;
  bumpProgress: () => void;
};

const defaultCapabilities: CapabilitiesState = {
  progress: false,
  history: false,
  loaded: false,
  progressRevision: 0,
  bumpProgress: () => {},
};

const CapabilitiesContext = createContext<CapabilitiesState>(defaultCapabilities);

export function CapabilitiesProvider({ children }: { children: ReactNode }) {
  const [capabilities, setCapabilities] = useState<Capabilities>({
    progress: false,
    history: false,
  });
  const [loaded, setLoaded] = useState(false);
  const [progressRevision, setProgressRevision] = useState(0);
  const bumpProgress = useCallback(() => {
    setProgressRevision((current) => current + 1);
  }, []);

  useEffect(() => {
    let cancelled = false;
    api
      .getCapabilities()
      .then((next) => {
        if (cancelled) return;
        setCapabilities(next);
        setLoaded(true);
      })
      .catch((error) => {
        console.warn("failed to load capabilities", error);
        if (cancelled) return;
        setCapabilities({ progress: false, history: false });
        setLoaded(true);
      });
    return () => {
      cancelled = true;
    };
  }, []);

  const value = useMemo(
    () => ({ ...capabilities, loaded, progressRevision, bumpProgress }),
    [bumpProgress, capabilities, loaded, progressRevision]
  );

  return (
    <CapabilitiesContext.Provider value={value}>
      {children}
    </CapabilitiesContext.Provider>
  );
}

export function useCapabilities() {
  return useContext(CapabilitiesContext);
}
