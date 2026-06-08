"use client";

import { useEffect, useState, useRef } from "react";
import { useRouter } from "@/i18n/routing";
import { useTranslations } from "next-intl";
import { useSearchParams } from "next/navigation";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { ScrollArea } from "@/components/ui/scroll-area";
import { ArrowLeft, Loader2, Mic, Volume2 } from "lucide-react";
import { motion, AnimatePresence } from "framer-motion";
import { useVoiceTutor } from "@/hooks/useVoiceTutor";
import { PushToTalkButton } from "@/components/voice/PushToTalkButton";
import { TranscriptBubble } from "@/components/voice/TranscriptBubble";
import { CorrectionList } from "@/components/voice/CorrectionList";
import { api } from "@/lib/api";
import { toast } from "sonner";

export default function VoiceTutorPage({
  params,
}: {
  params: Promise<{ profileId: string }>;
}) {
  const router = useRouter();
  const t = useTranslations("Voice");
  const searchParams = useSearchParams();
  const topic = searchParams.get("topic") || t("defaultTopic");
  const [profileId, setProfileId] = useState<string | null>(null);
  const [locale, setLocale] = useState("en-US");
  const scrollRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    params.then((p) => setProfileId(p.profileId));
  }, [params]);

  useEffect(() => {
    if (!profileId) return;
    api
      .getVoiceLocales(profileId)
      .then((loc) => setLocale(loc.primary))
      .catch(() => {});
  }, [profileId]);

  const {
    messages,
    loading,
    streamingContent,
    corrections,
    connected,
    capture,
    handlePushToTalkEnd,
  } = useVoiceTutor(profileId, topic, locale);

  useEffect(() => {
    if (!capture.supported) {
      toast.error(t("speechNotSupported"));
    }
  }, [capture.supported, t]);

  useEffect(() => {
    scrollRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages, streamingContent, capture.displayText]);

  return (
    <div className="min-h-full bg-[radial-gradient(ellipse_at_top,_var(--tw-gradient-stops))] from-indigo-50 via-white to-white dark:from-slate-900 dark:via-slate-950 dark:to-slate-950 py-8 px-4">
      <div className="container mx-auto max-w-4xl">
        <div className="mb-6 flex items-center justify-between">
          <Button variant="ghost" onClick={() => router.back()} className="rounded-full">
            <ArrowLeft className="w-4 h-4 mr-2" />
            {t("back")}
          </Button>
          <div className="flex items-center gap-2 text-sm text-muted-foreground">
            <Volume2 className="h-4 w-4" />
            {topic}
          </div>
        </div>

        <Card className="border-none shadow-2xl bg-white/70 dark:bg-slate-900/70 backdrop-blur-xl min-h-[600px] flex flex-col">
          <CardHeader className="border-b border-slate-100 dark:border-slate-800">
            <CardTitle className="flex items-center gap-2">
              <Mic className="h-5 w-5 text-primary" />
              {t("title")}
            </CardTitle>
            <p className="text-sm text-muted-foreground">{t("subtitle")}</p>
          </CardHeader>
          <CardContent className="flex-1 flex flex-col p-0 overflow-hidden">
            <ScrollArea className="flex-1 px-6 py-6">
              <div className="space-y-4">
                <AnimatePresence>
                  {messages.map((msg, idx) => (
                    <motion.div key={idx} initial={{ opacity: 0 }} animate={{ opacity: 1 }}>
                      <TranscriptBubble role={msg.role} content={msg.content} />
                    </motion.div>
                  ))}
                  {capture.displayText && loading && (
                    <TranscriptBubble role="user" content={capture.displayText} />
                  )}
                  {streamingContent && (
                    <TranscriptBubble role="assistant" content={streamingContent} />
                  )}
                </AnimatePresence>
                <CorrectionList corrections={corrections} title={t("corrections")} />
                <div ref={scrollRef} />
              </div>
            </ScrollArea>
            <div className="p-8 flex flex-col items-center gap-3 border-t border-slate-100 dark:border-slate-800">
              <PushToTalkButton
                isRecording={capture.isRecording}
                disabled={loading || !connected || !capture.supported}
                onPressStart={capture.start}
                onPressEnd={handlePushToTalkEnd}
                label={t("holdToTalk")}
                recordingLabel={t("releaseToSend")}
              />
              <p className="text-xs text-muted-foreground text-center">{t("pttHint")}</p>
              {loading && <Loader2 className="h-5 w-5 animate-spin text-primary" />}
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
