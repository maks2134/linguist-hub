"use client";

import { useCallback, useEffect, useRef, useState } from "react";
import { getSpeechRecognitionCtor, isSpeechRecognitionSupported } from "@/lib/voice/codec";
import { getSTTStrategy, STTStrategy } from "@/lib/voice/stt-strategy";
import { api } from "@/lib/api";

export function useVoiceCapture(locale: string) {
  const [isRecording, setIsRecording] = useState(false);
  const [partialText, setPartialText] = useState("");
  const [finalText, setFinalText] = useState("");
  const [supported, setSupported] = useState(false);
  const [sttStrategy, setSttStrategy] = useState<STTStrategy>("web-speech");
  const recognitionRef = useRef<SpeechRecognition | null>(null);
  const mediaRecorderRef = useRef<MediaRecorder | null>(null);
  const mediaStreamRef = useRef<MediaStream | null>(null);
  const chunksRef = useRef<Blob[]>([]);

  useEffect(() => {
    getSTTStrategy().then((strategy) => {
      setSttStrategy(strategy);
      if (strategy === "server") {
        setSupported(
          typeof navigator !== "undefined" && !!navigator.mediaDevices?.getUserMedia
        );
      } else {
        setSupported(isSpeechRecognitionSupported());
      }
    });
  }, []);

  const stop = useCallback(() => {
    recognitionRef.current?.stop();
    if (mediaRecorderRef.current?.state === "recording") {
      mediaRecorderRef.current.stop();
    } else {
      setIsRecording(false);
    }
  }, []);

  const startWebSpeech = useCallback(() => {
    const Ctor = getSpeechRecognitionCtor();
    if (!Ctor) return;

    setPartialText("");
    setFinalText("");

    const recognition = new Ctor();
    recognition.lang = locale;
    recognition.interimResults = true;
    recognition.continuous = false;
    recognition.maxAlternatives = 1;

    recognition.onresult = (event: SpeechRecognitionEvent) => {
      let interim = "";
      let final = "";
      for (let i = event.resultIndex; i < event.results.length; i++) {
        const t = event.results[i][0].transcript;
        if (event.results[i].isFinal) final += t;
        else interim += t;
      }
      if (interim) setPartialText(interim);
      if (final) {
        setFinalText((prev) => (prev ? `${prev} ${final}` : final).trim());
        setPartialText("");
      }
    };

    recognition.onend = () => setIsRecording(false);
    recognition.onerror = () => setIsRecording(false);

    recognitionRef.current = recognition;
    recognition.start();
    setIsRecording(true);
  }, [locale]);

  const startServerSTT = useCallback(async () => {
    setPartialText("");
    setFinalText("");

    const stream = await navigator.mediaDevices.getUserMedia({ audio: true });
    mediaStreamRef.current = stream;

    const recorder = new MediaRecorder(stream);
    chunksRef.current = [];

    recorder.ondataavailable = (event) => {
      if (event.data.size > 0) {
        chunksRef.current.push(event.data);
      }
    };

    recorder.onstop = async () => {
      mediaStreamRef.current?.getTracks().forEach((track) => track.stop());
      mediaStreamRef.current = null;
      setIsRecording(false);

      if (chunksRef.current.length === 0) return;

      const blob = new Blob(chunksRef.current, {
        type: recorder.mimeType || "audio/webm",
      });

      try {
        const result = await api.transcribeAudio(blob, locale);
        setFinalText(result.text.trim());
      } catch {
        setFinalText("");
      }
    };

    mediaRecorderRef.current = recorder;
    recorder.start();
    setIsRecording(true);
  }, [locale]);

  const start = useCallback(() => {
    if (sttStrategy === "server") {
      void startServerSTT();
    } else {
      startWebSpeech();
    }
  }, [sttStrategy, startServerSTT, startWebSpeech]);

  useEffect(() => {
    return () => {
      recognitionRef.current?.abort();
      if (mediaRecorderRef.current?.state === "recording") {
        mediaRecorderRef.current.stop();
      }
      mediaStreamRef.current?.getTracks().forEach((track) => track.stop());
    };
  }, []);

  const reset = useCallback(() => {
    setPartialText("");
    setFinalText("");
  }, []);

  return {
    supported,
    sttStrategy,
    isRecording,
    partialText,
    finalText,
    displayText: finalText || partialText,
    start,
    stop,
    reset,
  };
}
