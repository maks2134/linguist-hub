"use client";

import { Badge } from "@/components/ui/badge";

interface LessonStageHeaderProps {
  stage: string;
  stageIndex: number;
  stageTotal: number;
  targetPhrases: string[];
  hitPhrases: string[];
  labels: {
    stage: string;
    targets: string;
  };
}

export function LessonStageHeader({
  stage,
  stageIndex,
  stageTotal,
  targetPhrases,
  hitPhrases,
  labels,
}: LessonStageHeaderProps) {
  return (
    <div className="space-y-3 pb-4 border-b border-slate-100 dark:border-slate-800">
      <div className="flex items-center gap-2 flex-wrap">
        <Badge variant="outline" className="capitalize">
          {labels.stage}: {stage.replace("_", " ")}
        </Badge>
        <span className="text-xs text-muted-foreground">
          {stageIndex + 1} / {stageTotal}
        </span>
      </div>
      {targetPhrases.length > 0 && (
        <div>
          <p className="text-xs font-semibold text-muted-foreground mb-2">{labels.targets}</p>
          <div className="flex flex-wrap gap-2">
            {targetPhrases.map((phrase) => {
              const hit = hitPhrases.includes(phrase);
              return (
                <Badge
                  key={phrase}
                  variant={hit ? "default" : "secondary"}
                  className={hit ? "bg-green-600" : ""}
                >
                  {phrase}
                </Badge>
              );
            })}
          </div>
        </div>
      )}
    </div>
  );
}
