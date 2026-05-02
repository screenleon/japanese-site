import { useEffect, useRef } from "react";
import { api, type ReadKey } from "../api";
import { useCapabilities } from "../capabilities";

function readKeyIdentifier(key: ReadKey) {
  switch (key.type) {
    case "grammar":
      return key.slug;
    case "vocab":
      return key.headword;
    case "kanji":
      return key.character;
  }
}

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

  useEffect(() => {
    if (!loaded || !progress || !key || !normalizedIdentifier) return;

    const trackingKey = `${key.type}:${normalizedIdentifier}`;
    if (lastTracked.current === trackingKey) return;
    lastTracked.current = trackingKey;

    void api
      .markRead(normalizeReadKey(key, normalizedIdentifier))
      .then(bumpProgress)
      .catch((error) => {
        console.warn("failed to mark content as read", error);
      });
  }, [bumpProgress, key, loaded, normalizedIdentifier, progress]);
}
