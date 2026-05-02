import { useEffect, useRef } from "react";
import { api, type ReadContentType } from "../api";

export function useReadTracking(contentType: ReadContentType, slug?: string | null) {
  const lastTracked = useRef("");
  const normalizedSlug = slug?.trim() || "";

  useEffect(() => {
    if (!normalizedSlug) return;

    const key = `${contentType}:${normalizedSlug}`;
    if (lastTracked.current === key) return;
    lastTracked.current = key;

    void api.markRead(contentType, normalizedSlug).catch((error) => {
      console.warn("failed to mark content as read", error);
    });
  }, [contentType, normalizedSlug]);
}
