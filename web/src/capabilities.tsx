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
import type { ReadContentType } from "./apiTypes";

type CapabilitiesState = Capabilities & {
  loaded: boolean;
  progressRevisions: Record<ReadContentType, number>;
  bumpProgress: (type: ReadContentType) => void;
};

const initialProgressRevisions: Record<ReadContentType, number> = {
  grammar: 0,
  vocab: 0,
  kanji: 0,
};

const defaultCapabilities: CapabilitiesState = {
  progress: false,
  history: false,
  quiz: false,
  sentence: false,
  loaded: false,
  progressRevisions: initialProgressRevisions,
  bumpProgress: () => {},
};

const CapabilitiesContext = createContext<CapabilitiesState>(defaultCapabilities);

export function CapabilitiesProvider({ children }: { children: ReactNode }) {
  const [capabilities, setCapabilities] = useState<Capabilities>({
    progress: false,
    history: false,
    quiz: false,
    sentence: false,
  });
  const [loaded, setLoaded] = useState(false);
  const [progressRevisions, setProgressRevisions] =
    useState<Record<ReadContentType, number>>(initialProgressRevisions);
  const bumpProgress = useCallback((type: ReadContentType) => {
    setProgressRevisions((current) => ({
      ...current,
      [type]: current[type] + 1,
    }));
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
        setCapabilities({
          progress: false,
          history: false,
          quiz: false,
          sentence: false,
        });
        setLoaded(true);
      });
    return () => {
      cancelled = true;
    };
  }, []);

  const value = useMemo(
    () => ({ ...capabilities, loaded, progressRevisions, bumpProgress }),
    [bumpProgress, capabilities, loaded, progressRevisions]
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
