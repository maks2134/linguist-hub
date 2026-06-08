"use client";

import { useEffect, useState, useRef } from "react";
import { useRouter } from "@/i18n/routing";
import { useTranslations } from "next-intl";
import { useSearchParams } from "next/navigation";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { ScrollArea } from "@/components/ui/scroll-area";
import { ArrowLeft, Loader2 } from "lucide-react";
import { useVoiceLesson } from "@/hooks/useVoiceLesson";
import { PushToTalkButton } from "@/components/voice/PushToTalkButton";
import { TranscriptBubble } from "@/components/voice/TranscriptBubble";
import { LessonStageHeader } from "@/components/voice/LessonStageHeader";
import { LessonScoreCard } from "@/components/voice/LessonScoreCard";
import { api } from "@/lib/api";
import { toast } from "sonner";

export default function VoiceLessonPage({
  params,
}: {
  params: Promise<{ profileId: string }>;
}) {
  const router = useRouter();
  const t = useTranslations("Lesson");
  const searchParams = useSearchParams();
  const topic = searchParams.get("topic") || t("defaultTopic");
  const moduleId = parseInt(searchParams.get("moduleId") || "0", 10);
  const [profileId, setProfileId] = useState<string | null>(null);
  const [locale, setLocale] = useState("en-US");
  const scrollRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    params.then((p) => setProfileId(p.profileId));
  }, [params]);

  useEffect(() => {
    if (!profileId) return;
    api.getVoiceLocales(profileId).then((loc) => setLocale(loc.primary)).catch(() => {});
  }, [profileId]);

  const {
    lesson,
    userLines,
    loading,
    streamingPrompt,
    connected,
    capture,
    handlePushToTalkEnd,
    startLesson,
  } = useVoiceLesson(profileId, moduleId, topic, locale, true);

  useEffect(() => {
    if (!capture.supported) toast.error(t("speechNotSupported"));
  }, [capture.supported, t]);

  useEffect(() => {
    scrollRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [userLines, streamingPrompt, lesson?.prompt]);

  const promptText = streamingPrompt || lesson?.prompt || "";

  return (
    <div className="min-h-full bg-[radial-gradient(ellipse_at_top,_var(--tw-gradient-stops))] from-violet-50 via-white to-white dark:from-slate-900 dark:via-slate-950 dark:to-slate-950 py-8 px-4">
      <div className="container mx-auto max-w-4xl">
        <Button variant="ghost" onClick={() => router.back()} className="rounded-full mb-6">
          <ArrowLeft className="w-4 h-4 mr-2" />
          {t("back")}
        </Button>

        <Card className="border-none shadow-2xl bg-white/70 dark:bg-slate-900/70 backdrop-blur-xl min-h-[600px] flex flex-col">
          <CardHeader>
            <CardTitle>{t("title")}</CardTitle>
            <p className="text-sm text-muted-foreground">{topic}</p>
          </CardHeader>
          <CardContent className="flex-1 flex flex-col p-0 overflow-hidden">
            {lesson && (
              <div className="px-6 pt-2">
                <LessonStageHeader
                  stage={lesson.stage}
                  stageIndex={lesson.stageIndex}
                  stageTotal={lesson.stageTotal}
                  targetPhrases={lesson.targetPhrases}
                  hitPhrases={lesson.hitPhrases}
                  labels={{ stage: t("stage"), targets: t("targetPhrases") }}
                />
              </div>
            )}
            <ScrollArea className="flex-1 px-6 py-4">
              <div className="space-y-4">
                {promptText && <TranscriptBubble role="assistant" content={promptText} />}
                {userLines.map((line, i) => (
                  <TranscriptBubble key={i} role="user" content={line} />
                ))}
                {capture.displayText && loading && (
                  <TranscriptBubble role="user" content={capture.displayText} />
                )}
                {lesson?.score && (
                  <LessonScoreCard
                    score={lesson.score}
                    title={t("lessonComplete")}
                    labels={{
                      fluency: t("fluency"),
                      grammar: t("grammar"),
                      vocabulary: t("vocabulary"),
                      pronunciation: t("pronunciation"),
                    }}
                  />
                )}
                <div ref={scrollRef} />
              </div>
            </ScrollArea>
            {lesson?.stage !== "completed" && (
              <div className="p-8 flex flex-col items-center gap-3 border-t border-slate-100 dark:border-slate-800">
                {!lesson && !loading && (
                  <Button onClick={startLesson} disabled={!connected}>
                    {t("startLesson")}
                  </Button>
                )}
                <PushToTalkButton
                  isRecording={capture.isRecording}
                  disabled={loading || !connected || !capture.supported || !lesson}
                  onPressStart={capture.start}
                  onPressEnd={handlePushToTalkEnd}
                  label={t("holdToTalk")}
                  recordingLabel={t("releaseToSend")}
                />
                {loading && <Loader2 className="h-5 w-5 animate-spin text-primary" />}
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
