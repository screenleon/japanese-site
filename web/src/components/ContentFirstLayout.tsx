import type { ReactNode } from "react";

export type ContentFirstLayoutProps = {
  content: ReactNode;
  directory: ReactNode;
  /** Mobile toggle label, e.g. from `formatDirectorySummary(n)`. */
  directorySummary: string;
  /**
   * Whether the directory panel is expanded on viewports below `lg`.
   * Owned by the parent tab (UI-001). Ignored visually at `lg+` (always shown).
   */
  directoryOpen: boolean;
  onDirectoryOpenChange: (open: boolean) => void;
  contentTestId?: string;
  directoryTestId?: string;
};

/** Default mobile directory toggle copy for reference tabs. */
export function formatDirectorySummary(count?: number): string {
  if (count === undefined || count < 0) return "本級目錄";
  return `本級目錄 (${count})`;
}

/**
 * Shared reference-tab chrome: study **content first**, directory second on
 * small viewports; side-by-side (directory | content) from `lg` up.
 *
 * ## Reuse contract (JS-143+)
 *
 * | Reuse | Do |
 * |-------|-----|
 * | Grammar / Vocab / future Kanji browse | Yes — wrap focus pane + level directory |
 * | Quiz / Kokugo multi-step flows | No — different IA |
 * | Open state | Parent tab `useState` only; close on item select |
 * | Directory body | Free-form (list, search+list); layout does not own data |
 *
 * Deep-links (JS-144) set the active entry in the parent tab; this component
 * only orders chrome. Do not put fetch or navigation inside this file.
 */
export function ContentFirstLayout({
  content,
  directory,
  directorySummary,
  directoryOpen,
  onDirectoryOpenChange,
  contentTestId = "study-content",
  directoryTestId = "study-directory",
}: ContentFirstLayoutProps) {
  const panelId = `${directoryTestId}-panel`;

  return (
    <div
      data-testid="content-first-layout"
      className="flex flex-col gap-6 lg:grid lg:grid-cols-[240px_minmax(0,1fr)] lg:items-start"
    >
      <div
        data-testid={contentTestId}
        className="order-1 min-w-0 lg:col-start-2 lg:row-start-1"
      >
        {content}
      </div>

      <div
        data-testid={directoryTestId}
        className="order-2 min-w-0 lg:col-start-1 lg:row-start-1"
      >
        <button
          type="button"
          className={
            "mb-2 flex w-full items-center justify-between gap-2 rounded-md border border-slate-200 " +
            "bg-white px-4 py-3 text-left text-sm font-medium text-slate-700 hover:bg-slate-50 lg:hidden"
          }
          aria-expanded={directoryOpen}
          aria-controls={panelId}
          onClick={() => onDirectoryOpenChange(!directoryOpen)}
        >
          <span>{directorySummary}</span>
          <span className="text-xs font-normal text-slate-400" aria-hidden>
            {directoryOpen ? "收合" : "展開"}
          </span>
        </button>

        <div
          id={panelId}
          className={directoryOpen ? "block" : "hidden lg:block"}
        >
          {directory}
        </div>
      </div>
    </div>
  );
}
