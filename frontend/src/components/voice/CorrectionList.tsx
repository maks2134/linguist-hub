"use client";

import { VoiceCorrection } from "@/lib/websocket";

interface CorrectionListProps {
  corrections: VoiceCorrection[];
  title: string;
}

export function CorrectionList({ corrections, title }: CorrectionListProps) {
  if (!corrections.length) return null;
  return (
    <div className="rounded-xl border border-amber-200 dark:border-amber-900/50 bg-amber-50/80 dark:bg-amber-950/30 p-4 space-y-2">
      <p className="text-sm font-semibold text-amber-900 dark:text-amber-200">{title}</p>
      <ul className="space-y-2 text-sm">
        {corrections.map((c, i) => (
          <li key={i}>
            <span className="line-through text-muted-foreground">{c.original}</span>
            {" → "}
            <span className="font-medium">{c.suggested}</span>
            {c.reason && <span className="text-muted-foreground"> — {c.reason}</span>}
          </li>
        ))}
      </ul>
    </div>
  );
}
