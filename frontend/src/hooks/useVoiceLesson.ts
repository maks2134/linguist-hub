"use client";

import { useCallback, useEffect, useState } from "react";
import { useWebSocket } from "@/hooks/useWebSocket";
import { useSpeechSynthesis } from "@/hooks/useSpeechSynthesis";
import { useVoiceCapture } from "@/hooks/useVoiceCapture";
import { LessonScore, WebSocketResponse } from "@/lib/websocket";
import { toast } from "sonner";

export interface LessonState {
  lessonId: number;
  stage: string;
  stageIndex: number;
  stageTotal: number;
  prompt: string;
  targetPhrases: string[];
  hitPhrases: string[];
  score?: LessonScore;
}

export function useVoiceLesson(
  profileId: string | null,
  moduleId: number,
  topic: string,
  locale: string,
  autoStart: boolean
) {
  const [lesson, setLesson] = useState<LessonState | null>(null);
  const [userLines, setUserLines] = useState<string[]>([]);
  const [loading, setLoading] = useState(false);
  const [streamingPrompt, setStreamingPrompt] = useState("");
  const { connected, sendMessage, onMessage } = useWebSocket();
  const { feedChunk, flush, cancel } = useSpeechSynthesis(locale);
  const capture = useVoiceCapture(locale);

  const applyStageResponse = useCallback((response: WebSocketResponse) => {
    if (!response.lesson_id) return;
    setLesson({
      lessonId: response.lesson_id,
      stage: response.stage || "warmup",
      stageIndex: response.stage_index ?? 0,
      stageTotal: response.stage_total ?? 4,
      prompt: response.prompt || response.content || "",
      targetPhrases: response.target_phrases || [],
      hitPhrases: response.hit_phrases || [],
      score: response.score,
    });
  }, []);

  useEffect(() => {
    if (!connected) return;
    return onMessage((response: WebSocketResponse) => {
      if (response.type !== "voice_lesson_stage" && response.type !== "voice_lesson_complete") {
        return;
      }
      if (response.error) {
        toast.error(response.error);
        setLoading(false);
        cancel();
        return;
      }
      if (response.chunk) {
        setStreamingPrompt((prev) => prev + response.chunk);
        feedChunk(response.chunk!);
      }
      if (response.done) {
        flush();
        setStreamingPrompt("");
        setLoading(false);
        if (response.type === "voice_lesson_complete") {
          applyStageResponse(response);
        } else if (response.type === "voice_lesson_stage") {
          applyStageResponse(response);
        }
      }
    });
  }, [connected, onMessage, feedChunk, flush, cancel, applyStageResponse]);

  const startLesson = useCallback(() => {
    if (!profileId || !connected) return;
    setLoading(true);
    cancel();
    sendMessage({
      type: "voice_lesson_start",
      profile_id: parseInt(profileId, 10),
      module_id: moduleId,
      topic,
      locale,
    });
  }, [profileId, connected, moduleId, topic, locale, sendMessage, cancel]);

  useEffect(() => {
    if (autoStart && connected && profileId && !lesson) {
      startLesson();
    }
  }, [autoStart, connected, profileId, lesson, startLesson]);

  const sendTurn = useCallback(
    (text: string) => {
      if (!lesson || !connected) return;
      setUserLines((prev) => [...prev, text]);
      setLoading(true);
      cancel();
      sendMessage({
        type: "voice_lesson_turn",
        lesson_id: lesson.lessonId,
        message: text,
        profile_id: parseInt(profileId!, 10),
        locale,
      });
    },
    [lesson, connected, profileId, locale, sendMessage, cancel]
  );

  const handlePushToTalkEnd = useCallback(() => {
    capture.stop();
    const text = capture.finalText || capture.partialText;
    capture.reset();
    if (text.trim()) sendTurn(text);
  }, [capture, sendTurn]);

  return {
    lesson,
    userLines,
    loading,
    streamingPrompt,
    connected,
    capture,
    startLesson,
    sendTurn,
    handlePushToTalkEnd,
  };
}
