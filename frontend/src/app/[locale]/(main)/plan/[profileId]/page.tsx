"use client";

import { useEffect, useState } from "react";
import { useRouter } from "@/i18n/routing";
import { useTranslations } from "next-intl";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { api } from "@/lib/api";
import { StudyPlan } from "@/types";
import { toast } from "sonner";
import { ArrowLeft, BookOpen, Map, Sparkles, ChevronRight, GraduationCap, Loader2 } from "lucide-react";
import { motion } from "framer-motion";

export default function PlanPage({
  params,
}: {
  params: Promise<{ profileId: string }>;
}) {
  const router = useRouter();
  const t = useTranslations("Plan");
  const tToast = useTranslations("Toast");
  const [profileId, setProfileId] = useState<string | null>(null);
  const [plan, setPlan] = useState<StudyPlan | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    params.then((p) => {
      setProfileId(p.profileId);
      loadPlan(p.profileId);
    });
  }, [params]);

  const loadPlan = async (id: string) => {
    try {
      const data = await api.get<StudyPlan>(`/api/ai/plan/${id}`);
      setPlan(data);
    } catch (error) {
      toast.error(tToast("loadPlanFailed"));
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="flex h-[calc(100vh-4rem)] items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    );
  }

  if (!plan) {
    return (
      <div className="min-h-full flex items-center justify-center p-4">
        <motion.div
          initial={{ opacity: 0, scale: 0.95 }}
          animate={{ opacity: 1, scale: 1 }}
        >
          <Card className="max-w-md border-none shadow-2xl bg-white/80 dark:bg-slate-900/80 backdrop-blur-xl">
            <CardHeader className="text-center">
              <div className="h-16 w-16 rounded-2xl bg-slate-100 dark:bg-slate-800 flex items-center justify-center mx-auto mb-4">
                <Map className="h-8 w-8 text-slate-400" />
              </div>
              <CardTitle className="text-2xl font-bold">{t("noPlanTitle")}</CardTitle>
              <CardDescription className="text-base">
                {t("noPlanDesc")}
              </CardDescription>
            </CardHeader>
            <CardContent className="flex justify-center pb-8">
              <Button onClick={() => router.push("/dashboard")} className="rounded-full px-8">
                {t("goToDashboard")}
              </Button>
            </CardContent>
          </Card>
        </motion.div>
      </div>
    );
  }

  return (
    <div className="min-h-full bg-[radial-gradient(ellipse_at_top,_var(--tw-gradient-stops))] from-indigo-50 via-white to-white dark:from-slate-900 dark:via-slate-950 dark:to-slate-950 py-10 px-4">
      <div className="container mx-auto max-w-6xl">
        <div className="mb-10 flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <motion.div
            initial={{ opacity: 0, x: -20 }}
            animate={{ opacity: 1, x: 0 }}
          >
            <Button variant="ghost" onClick={() => router.push("/dashboard")} className="rounded-full mb-4 sm:mb-0">
              <ArrowLeft className="w-4 h-4 mr-2" />
              {t("backToDashboard")}
            </Button>
          </motion.div>
          <motion.div
             initial={{ opacity: 0, scale: 0.9 }}
             animate={{ opacity: 1, scale: 1 }}
             className="flex items-center gap-2 px-6 py-2 rounded-full bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 shadow-sm"
          >
             <GraduationCap className="h-5 w-5 text-primary" />
             <span className="font-semibold text-sm">{t("personalizedRoadmap")}</span>
          </motion.div>
        </div>

        <header className="mb-12">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.1 }}
          >
            <div className="flex items-center gap-3 mb-2">
               <h1 className="text-4xl font-extrabold tracking-tight lg:text-5xl bg-gradient-to-r from-slate-900 to-slate-700 dark:from-white dark:to-slate-400 bg-clip-text text-transparent">
                 {t("learningJourney")}
               </h1>
               <div className="h-10 w-10 rounded-xl bg-yellow-500/10 flex items-center justify-center">
                  <Sparkles className="h-6 w-6 text-yellow-500" />
               </div>
            </div>
            <p className="text-lg text-muted-foreground">
              {t("learningJourneyDesc")}
            </p>
          </motion.div>
        </header>

        <div className="grid gap-8 relative">
          {/* Vertical connection line for roadmap feel */}
          <div className="absolute left-8 top-10 bottom-10 w-1 bg-gradient-to-b from-primary/20 via-primary/10 to-transparent hidden lg:block" />

          {plan.modules.map((module, index) => (
            <motion.div
              key={module.id}
              initial={{ opacity: 0, y: 30 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: index * 0.1 }}
              className="relative lg:pl-20"
            >
              {/* Dot for roadmap */}
              <div className="absolute left-[28px] top-8 h-3 w-3 rounded-full bg-primary ring-4 ring-primary/20 hidden lg:block" />
              
              <Card className="border-none shadow-md hover:shadow-xl transition-all duration-300 bg-white/60 dark:bg-slate-900/60 backdrop-blur-md overflow-hidden ring-1 ring-slate-100 dark:ring-slate-800">
                <CardHeader className="pb-4">
                  <div className="flex flex-col md:flex-row md:items-start justify-between gap-4">
                    <div className="space-y-1">
                      <div className="flex items-center gap-3">
                        <Badge variant="outline" className="bg-primary/5 text-primary border-primary/20 font-bold px-3">
                          {t("module")} {module.id}
                        </Badge>
                        <span className="text-sm text-muted-foreground font-medium flex items-center gap-1">
                           <BookOpen className="h-3 w-3" />
                           {module.topics.length} {t("topics")}
                        </span>
                      </div>
                      <CardTitle className="text-2xl font-bold tracking-tight pt-2">
                        {module.title}
                      </CardTitle>
                      <CardDescription className="text-base pt-1 max-w-2xl">
                        {module.description}
                      </CardDescription>
                    </div>
                  </div>
                </CardHeader>
                <CardContent>
                  <div className="space-y-4">
                    <div className="flex flex-wrap gap-2 pt-2">
                      {module.topics.map((topic, idx) => (
                        <motion.button
                          key={idx}
                          whileHover={{ y: -2 }}
                          whileTap={{ scale: 0.98 }}
                          onClick={() =>
                            router.push(`/tutor/${profileId}?topic=${encodeURIComponent(topic)}`)
                          }
                          className="flex items-center gap-3 px-4 py-2.5 rounded-xl bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 hover:border-primary/50 hover:shadow-lg hover:shadow-primary/5 transition-all group"
                        >
                          <div className="h-2 w-2 rounded-full bg-slate-300 group-hover:bg-primary transition-colors" />
                          <span className="text-sm font-medium pr-1">{topic}</span>
                          <ChevronRight className="h-4 w-4 text-slate-400 group-hover:text-primary transition-colors" />
                        </motion.button>
                      ))}
                    </div>
                  </div>
                </CardContent>
              </Card>
            </motion.div>
          ))}
        </div>
      </div>
    </div>
  );
}

