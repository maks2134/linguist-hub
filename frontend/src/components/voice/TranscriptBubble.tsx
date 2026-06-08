"use client";

import { UserCircle2, Sparkles } from "lucide-react";

interface TranscriptBubbleProps {
  role: "user" | "assistant";
  content: string;
}

export function TranscriptBubble({ role, content }: TranscriptBubbleProps) {
  return (
    <div
      className={`flex items-start gap-3 ${role === "user" ? "flex-row-reverse" : "flex-row"}`}
    >
      <div
        className={`flex-shrink-0 h-8 w-8 rounded-full flex items-center justify-center ${
          role === "user" ? "bg-slate-200 dark:bg-slate-700" : "bg-primary/10 text-primary"
        }`}
      >
        {role === "user" ? <UserCircle2 className="h-5 w-5" /> : <Sparkles className="h-4 w-4" />}
      </div>
      <div
        className={`max-w-[80%] rounded-2xl px-5 py-3 shadow-sm ${
          role === "user"
            ? "bg-primary text-primary-foreground rounded-tr-none"
            : "bg-white dark:bg-slate-800 border border-slate-100 dark:border-slate-700 rounded-tl-none"
        }`}
      >
        <p className="whitespace-pre-wrap leading-relaxed text-[15px]">{content}</p>
      </div>
    </div>
  );
}
