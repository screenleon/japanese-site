import { useEffect, useRef } from "react";
import { api } from "../api";
import { readKeyIdentifier, type ReadKey } from "../apiTypes";
import { useCapabilities } from "../capabilities";

function normalizeReadKey(key: ReadKey, identifier: string): ReadKey {
  switch (key.type) {
    case "grammar":
      return { type: "grammar", slug: identifier };
    case "vocab":
      return { type: "vocab", headword: identifier };
    case "kanji":
      return { type: "kanji", character: identifier };
  }
}

export function useReadTracking(key: ReadKey | null | undefined) {
  const { progress, loaded, bumpProgress } = useCapabilities();
  const lastTracked = useRef("");
  const normalizedIdentifier = key ? readKeyIdentifier(key).trim() : "";
  const keyType = key?.type;

  useEffect(() => {
    if (!loaded || !progress || !key || !normalizedIdentifier) return;

    const trackingKey = `${key.type}:${normalizedIdentifier}`;
    if (lastTracked.current === trackingKey) return;
    lastTracked.current = trackingKey;

    void api
      .markRead(normalizeReadKey(key, normalizedIdentifier))
      .then(() => bumpProgress(key.type))
      .catch((error) => {
        console.warn("failed to mark content as read", error);
      });
    // Depend on primitive identifier+type rather than the object so call
    // sites that pass freshly-allocated `{type, ...}` literals don't churn
    // the effect on unrelated parent re-renders.
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [bumpProgress, keyType, loaded, normalizedIdentifier, progress]);
}
