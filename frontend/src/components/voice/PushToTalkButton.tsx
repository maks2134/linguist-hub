"use client";

import { Mic, MicOff, Loader2 } from "lucide-react";
import { Button } from "@/components/ui/button";

interface PushToTalkButtonProps {
  isRecording: boolean;
  disabled?: boolean;
  onPressStart: () => void;
  onPressEnd: () => void;
  label: string;
  recordingLabel: string;
}

export function PushToTalkButton({
  isRecording,
  disabled,
  onPressStart,
  onPressEnd,
  label,
  recordingLabel,
}: PushToTalkButtonProps) {
  return (
    <Button
      type="button"
      variant={isRecording ? "destructive" : "default"}
      disabled={disabled}
      className="h-16 w-16 rounded-full p-0 shadow-xl"
      onMouseDown={onPressStart}
      onMouseUp={onPressEnd}
      onMouseLeave={isRecording ? onPressEnd : undefined}
      onTouchStart={(e) => {
        e.preventDefault();
        onPressStart();
      }}
      onTouchEnd={(e) => {
        e.preventDefault();
        onPressEnd();
      }}
    >
      {disabled ? (
        <Loader2 className="h-7 w-7 animate-spin" />
      ) : isRecording ? (
        <MicOff className="h-7 w-7" />
      ) : (
        <Mic className="h-7 w-7" />
      )}
      <span className="sr-only">{isRecording ? recordingLabel : label}</span>
    </Button>
  );
}
