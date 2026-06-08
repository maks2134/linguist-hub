"use client";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { LessonScore } from "@/lib/websocket";

interface LessonScoreCardProps {
  score: LessonScore;
  title: string;
  labels: {
    fluency: string;
    grammar: string;
    vocabulary: string;
    pronunciation: string;
  };
}

function ScoreBar({ label, value }: { label: string; value: number }) {
  const pct = Math.round(value * 100);
  return (
    <div>
      <div className="flex justify-between text-sm mb-1">
        <span>{label}</span>
        <span>{pct}%</span>
      </div>
      <div className="h-2 rounded-full bg-slate-200 dark:bg-slate-700 overflow-hidden">
        <div className="h-full bg-primary rounded-full transition-all" style={{ width: `${pct}%` }} />
      </div>
    </div>
  );
}

export function LessonScoreCard({ score, title, labels }: LessonScoreCardProps) {
  return (
    <Card className="border-none shadow-lg bg-white/80 dark:bg-slate-900/80">
      <CardHeader>
        <CardTitle>{title}</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <ScoreBar label={labels.fluency} value={score.fluency} />
        <ScoreBar label={labels.grammar} value={score.grammar} />
        <ScoreBar label={labels.vocabulary} value={score.vocabulary} />
        {score.pronunciation > 0 && (
          <ScoreBar label={labels.pronunciation} value={score.pronunciation} />
        )}
        {score.notes && <p className="text-sm text-muted-foreground pt-2">{score.notes}</p>}
      </CardContent>
    </Card>
  );
}
