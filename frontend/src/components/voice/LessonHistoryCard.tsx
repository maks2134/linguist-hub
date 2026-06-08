"use client";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { VoiceLessonHistory } from "@/types";
import { LessonScoreCard } from "@/components/voice/LessonScoreCard";
import { CheckCircle2, Clock, XCircle } from "lucide-react";

interface LessonHistoryCardProps {
  lesson: VoiceLessonHistory;
  labels: {
    completed: string;
    abandoned: string;
    inProgress: string;
    module: string;
    scoreTitle: string;
    fluency: string;
    grammar: string;
    vocabulary: string;
    pronunciation: string;
  };
}

function statusIcon(status: string) {
  switch (status) {
    case "completed":
      return <CheckCircle2 className="h-4 w-4 text-green-500" />;
    case "abandoned":
      return <XCircle className="h-4 w-4 text-slate-400" />;
    default:
      return <Clock className="h-4 w-4 text-amber-500" />;
  }
}

function statusLabel(status: string, labels: LessonHistoryCardProps["labels"]) {
  switch (status) {
    case "completed":
      return labels.completed;
    case "abandoned":
      return labels.abandoned;
    default:
      return labels.inProgress;
  }
}

export function LessonHistoryCard({ lesson, labels }: LessonHistoryCardProps) {
  const started = new Date(lesson.started_at).toLocaleString();
  const completed = lesson.completed_at
    ? new Date(lesson.completed_at).toLocaleString()
    : null;

  return (
    <Card className="border-none shadow-md bg-white/70 dark:bg-slate-900/70 backdrop-blur-md">
      <CardHeader className="pb-3">
        <div className="flex items-start justify-between gap-4">
          <div>
            <CardTitle className="text-lg">{lesson.topic}</CardTitle>
            <p className="text-sm text-muted-foreground mt-1">
              {labels.module} {lesson.module_id} · {started}
              {completed ? ` → ${completed}` : ""}
            </p>
          </div>
          <Badge variant="outline" className="flex items-center gap-1 shrink-0">
            {statusIcon(lesson.status)}
            {statusLabel(lesson.status, labels)}
          </Badge>
        </div>
      </CardHeader>
      {lesson.score && lesson.status === "completed" && (
        <CardContent>
          <LessonScoreCard
            score={lesson.score}
            title={labels.scoreTitle}
            labels={{
              fluency: labels.fluency,
              grammar: labels.grammar,
              vocabulary: labels.vocabulary,
              pronunciation: labels.pronunciation,
            }}
          />
        </CardContent>
      )}
    </Card>
  );
}
