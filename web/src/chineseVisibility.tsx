import {
  createContext,
  type ReactNode,
  useCallback,
  useContext,
  useEffect,
  useState,
} from "react";

const STORAGE_KEY = "japanese-site:chineseVisible";

type ChineseVisibility = {
  visible: boolean;
  setVisible: (next: boolean) => void;
  toggle: () => void;
};

const defaultValue: ChineseVisibility = {
  visible: false,
  setVisible: () => {},
  toggle: () => {},
};

const ChineseVisibilityContext = createContext<ChineseVisibility>(defaultValue);

export function ChineseVisibilityProvider({
  children,
  initialVisible,
}: {
  children: ReactNode;
  initialVisible?: boolean;
}) {
  const [visible, setVisibleRaw] = useState<boolean>(() => {
    if (typeof initialVisible === "boolean") return initialVisible;
    if (typeof window === "undefined") return false;
    try {
      return window.localStorage.getItem(STORAGE_KEY) === "true";
    } catch {
      return false;
    }
  });

  useEffect(() => {
    if (typeof window === "undefined") return;
    try {
      window.localStorage.setItem(STORAGE_KEY, visible ? "true" : "false");
    } catch {
      /* localStorage write may fail in private mode; tolerate silently */
    }
  }, [visible]);

  const setVisible = useCallback((next: boolean) => setVisibleRaw(next), []);
  const toggle = useCallback(() => setVisibleRaw((current) => !current), []);

  return (
    <ChineseVisibilityContext.Provider value={{ visible, setVisible, toggle }}>
      {children}
    </ChineseVisibilityContext.Provider>
  );
}

export function useChineseVisibility() {
  return useContext(ChineseVisibilityContext);
}

export function ChineseVisibilityToggle() {
  const { visible, toggle } = useChineseVisibility();
  return (
    <button
      type="button"
      onClick={toggle}
      aria-pressed={visible}
      title={visible ? "中文を非表示にする" : "中文を表示する"}
      className={`rounded-md border px-2.5 py-1 text-xs font-medium transition-colors ${
        visible
          ? "border-sky-300 bg-sky-100 text-sky-900"
          : "border-slate-300 bg-white text-slate-600 hover:bg-slate-50"
      }`}
    >
      中文
    </button>
  );
}

export function IfChinese({ children }: { children: ReactNode }) {
  const { visible } = useChineseVisibility();
  if (!visible) return null;
  return <>{children}</>;
}
