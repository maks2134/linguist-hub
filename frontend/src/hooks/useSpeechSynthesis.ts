"use client";

import { useCallback, useRef } from "react";

const SENTENCE_END = /[.!?…]+[\s\n]*/;

export function useSpeechSynthesis(locale: string) {
  const queueRef = useRef<string[]>([]);
  const speakingRef = useRef(false);
  const bufferRef = useRef("");

  const drainQueue = useCallback(() => {
    if (speakingRef.current || queueRef.current.length === 0) return;
    if (typeof window === "undefined" || !window.speechSynthesis) return;

    const next = queueRef.current.shift();
    if (!next) return;

    speakingRef.current = true;
    const utterance = new SpeechSynthesisUtterance(next);
    utterance.lang = locale;
    utterance.onend = () => {
      speakingRef.current = false;
      drainQueue();
    };
    utterance.onerror = () => {
      speakingRef.current = false;
      drainQueue();
    };
    window.speechSynthesis.speak(utterance);
  }, [locale]);

  const enqueueSentence = useCallback(
    (sentence: string) => {
      const trimmed = sentence.trim();
      if (!trimmed) return;
      queueRef.current.push(trimmed);
      drainQueue();
    },
    [drainQueue]
  );

  const feedChunk = useCallback(
    (chunk: string) => {
      bufferRef.current += chunk;
      const parts = bufferRef.current.split(SENTENCE_END);
      if (parts.length <= 1) return;
      for (let i = 0; i < parts.length - 1; i++) {
        enqueueSentence(parts[i]);
      }
      bufferRef.current = parts[parts.length - 1];
    },
    [enqueueSentence]
  );

  const flush = useCallback(() => {
    if (bufferRef.current.trim()) {
      enqueueSentence(bufferRef.current);
      bufferRef.current = "";
    }
  }, [enqueueSentence]);

  const cancel = useCallback(() => {
    bufferRef.current = "";
    queueRef.current = [];
    speakingRef.current = false;
    if (typeof window !== "undefined" && window.speechSynthesis) {
      window.speechSynthesis.cancel();
    }
  }, []);

  return { feedChunk, flush, cancel, enqueueSentence };
}
