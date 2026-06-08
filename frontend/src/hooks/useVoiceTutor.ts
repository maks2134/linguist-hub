"use client";

import { useCallback, useEffect, useState } from "react";
import { useWebSocket } from "@/hooks/useWebSocket";
import { useSpeechSynthesis } from "@/hooks/useSpeechSynthesis";
import { useVoiceCapture } from "@/hooks/useVoiceCapture";
import { api, ChatHistoryMessage } from "@/lib/api";
import { VoiceCorrection, WebSocketResponse } from "@/lib/websocket";
import { ChatMessage } from "@/types";
import { toast } from "sonner";

export function useVoiceTutor(profileId: string | null, topic: string, locale: string) {
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [loading, setLoading] = useState(false);
  const [streamingContent, setStreamingContent] = useState("");
  const [corrections, setCorrections] = useState<VoiceCorrection[]>([]);
  const { connected, sendMessage, onMessage } = useWebSocket();
  const { feedChunk, flush, cancel } = useSpeechSynthesis(locale);
  const capture = useVoiceCapture(locale);

  useEffect(() => {
    if (!profileId) return;
    Promise.all([
      api.getChatHistory(profileId, "tutor"),
      api.getChatHistory(profileId, "voice_tutor"),
    ])
      .then(([tutor, voice]) => {
        const merged = [...tutor, ...voice].sort(
          (a, b) => new Date(a.created_at).getTime() - new Date(b.created_at).getTime()
        );
        setMessages(merged.map((m) => ({ role: m.role, content: m.content })));
      })
      .catch(() => {});
  }, [profileId]);

  useEffect(() => {
    if (!connected) return;
    return onMessage((response: WebSocketResponse) => {
      if (response.type !== "voice_tutor_chunk" && response.type !== "voice_tutor_final") return;
      if (response.error) {
        toast.error(response.error);
        setLoading(false);
        cancel();
        return;
      }
      if (response.chunk) {
        setStreamingContent((prev) => prev + response.chunk);
        feedChunk(response.chunk!);
      }
      if (response.done && response.type === "voice_tutor_final") {
        flush();
        const content = response.content || "";
        setMessages((prev) => [...prev, { role: "assistant", content }]);
        if (response.user_text) {
          setMessages((prev) => {
            const withoutDup = prev.filter(
              (m, i) => !(i === prev.length - 1 && m.role === "user" && m.content === response.user_text)
            );
            return withoutDup;
          });
        }
        setCorrections(response.corrections || []);
        setStreamingContent("");
        setLoading(false);
      }
    });
  }, [connected, onMessage, feedChunk, flush, cancel]);

  const sendVoiceTurn = useCallback(
    (text: string) => {
      if (!profileId || !connected || !text.trim()) return;
      setMessages((prev) => [...prev, { role: "user", content: text }]);
      setLoading(true);
      setStreamingContent("");
      setCorrections([]);
      cancel();
      try {
        sendMessage({
          type: "voice_tutor_turn",
          message: text,
          topic,
          profile_id: parseInt(profileId, 10),
          locale,
        });
      } catch {
        toast.error("Failed to send voice message");
        setLoading(false);
      }
    },
    [profileId, connected, topic, locale, sendMessage, cancel]
  );

  const handlePushToTalkEnd = useCallback(() => {
    capture.stop();
    const text = capture.finalText || capture.partialText;
    capture.reset();
    if (text.trim()) sendVoiceTurn(text);
  }, [capture, sendVoiceTurn]);

  return {
    messages,
    loading,
    streamingContent,
    corrections,
    connected,
    capture,
    sendVoiceTurn,
    handlePushToTalkEnd,
  };
}
