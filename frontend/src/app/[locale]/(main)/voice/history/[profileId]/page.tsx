"use client";

import { useEffect, useState } from "react";
import { useRouter } from "@/i18n/routing";
import { useTranslations } from "next-intl";
import { Button } from "@/components/ui/button";
import { api } from "@/lib/api";
import { VoiceLessonHistory } from "@/types";
import { LessonHistoryCard } from "@/components/voice/LessonHistoryCard";
import { ArrowLeft, History, Loader2 } from "lucide-react";
import { toast } from "sonner";
import { motion } from "framer-motion";

export default function VoiceHistoryPage({
  params,
}: {
  params: Promise<{ profileId: string }>;
}) {
  const router = useRouter();
  const t = useTranslations("VoiceHistory");
  const tLesson = useTranslations("Lesson");
  const tToast = useTranslations("Toast");
  const [profileId, setProfileId] = useState<string | null>(null);
  const [lessons, setLessons] = useState<VoiceLessonHistory[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    params.then((p) => {
      setProfileId(p.profileId);
      loadLessons(p.profileId);
    });
  }, [params]);

  const loadLessons = async (id: string) => {
    try {
      const data = await api.listVoiceLessons(id);
      setLessons(data);
    } catch {
      toast.error(tToast("loadHistoryFailed"));
    } finally {
      setLoading(false);
    }
  };

  const scoreLabels = {
    fluency: tLesson("fluency"),
    grammar: tLesson("grammar"),
    vocabulary: tLesson("vocabulary"),
    pronunciation: tLesson("pronunciation"),
  };

  if (loading) {
    return (
      <div className="flex h-[calc(100vh-4rem)] items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    );
  }

  return (
    <div className="min-h-full bg-[radial-gradient(ellipse_at_top,_var(--tw-gradient-stops))] from-indigo-50 via-white to-white dark:from-slate-900 dark:via-slate-950 dark:to-slate-950 py-10 px-4">
      <div className="container mx-auto max-w-3xl">
        <div className="mb-8 flex items-center justify-between">
          <Button variant="ghost" onClick={() => router.push("/dashboard")} className="rounded-full">
            <ArrowLeft className="w-4 h-4 mr-2" />
            {t("backToDashboard")}
          </Button>
          <div className="flex items-center gap-2 text-sm text-muted-foreground">
            <History className="h-4 w-4" />
            {t("title")}
          </div>
        </div>

        <header className="mb-8">
          <h1 className="text-3xl font-bold tracking-tight">{t("title")}</h1>
          <p className="text-muted-foreground mt-2">{t("subtitle")}</p>
        </header>

        {lessons.length === 0 ? (
          <div className="text-center py-16 text-muted-foreground">{t("empty")}</div>
        ) : (
          <div className="space-y-4">
            {lessons.map((lesson, index) => (
              <motion.div
                key={lesson.lesson_id}
                initial={{ opacity: 0, y: 16 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: index * 0.05 }}
              >
                <LessonHistoryCard
                  lesson={lesson}
                  labels={{
                    completed: t("statusCompleted"),
                    abandoned: t("statusAbandoned"),
                    inProgress: t("statusInProgress"),
                    module: t("module"),
                    scoreTitle: t("scoreTitle"),
                    ...scoreLabels,
                  }}
                />
              </motion.div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
