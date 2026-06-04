"use client";

import { useState, useEffect } from "react";
import { useRouter } from "@/i18n/routing";
import { useTranslations } from "next-intl";
import { motion, AnimatePresence } from "framer-motion";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { api } from "@/lib/api";
import { Profile } from "@/types";
import { toast } from "sonner";
import { Globe, GraduationCap, ArrowRight, Check, Loader2 } from "lucide-react";
import { LEARNING_LANGUAGES, NATIVE_LANGUAGES, type LangCode } from "@/lib/language";

export default function SetupPage() {
  const router = useRouter();
  const t = useTranslations("Setup");
  const [mounted, setMounted] = useState(false);
  const [step, setStep] = useState(1);
  const [nativeLanguage, setNativeLanguage] = useState<LangCode | "">("");
  const [targetLanguage, setTargetLanguage] = useState<LangCode | "">("");
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    setMounted(true);
  }, []);

  if (!mounted) {
    return (
      <div className="min-h-screen bg-background flex items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    );
  }

  const handleFinish = async () => {
    if (!nativeLanguage || !targetLanguage) return;
    setLoading(true);
    try {
      const profile = await api.post<Profile>("/api/profile/select", {
        language: targetLanguage,
        native_language: nativeLanguage,
      });
      toast.success(t("profileCreated"));
      router.push(`/diagnostic/${profile.id}`);
    } catch (error) {
      toast.error(t("profileFailed"));
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-background flex flex-col items-center justify-center p-4 relative overflow-hidden text-foreground">
      {/* Background blobs */}
      <div className="absolute top-[-10%] left-[-10%] w-[50%] h-[50%] bg-primary/10 rounded-full blur-[120px] animate-pulse" />
      <div className="absolute bottom-[-10%] right-[-10%] w-[50%] h-[50%] bg-purple-500/10 rounded-full blur-[120px] animate-pulse" />

      <div className="max-w-4xl w-full relative z-10">
        <AnimatePresence mode="wait">
          {step === 1 && (
            <motion.div
              key="step1"
              initial={{ opacity: 0, x: 20 }}
              animate={{ opacity: 1, x: 0 }}
              exit={{ opacity: 0, x: -20 }}
              className="space-y-8"
            >
              <div className="text-center space-y-4">
                <motion.div 
                  initial={{ scale: 0 }}
                  animate={{ scale: 1 }}
                  className="inline-flex p-3 rounded-2xl bg-primary/20 text-primary mb-4"
                >
                  <Globe size={32} />
                </motion.div>
                <h1 className="text-4xl md:text-5xl font-black text-foreground">{t("chooseNativeTitle")}</h1>
                <p className="text-muted-foreground text-lg">{t("chooseNativeDesc")}</p>
              </div>

              <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
                {NATIVE_LANGUAGES.map((lang, idx) => (
                  <motion.div
                    key={lang.code}
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: idx * 0.05 }}
                  >
                    <Card
                      onClick={() => setNativeLanguage(lang.code)}
                      className={`p-6 cursor-pointer transition-all duration-300 group border-2 relative overflow-hidden ${
                        nativeLanguage === lang.code 
                        ? "border-primary bg-primary/10 scale-105 shadow-[0_0_30px_-10px_rgba(59,130,246,0.5)]" 
                        : "border-border bg-card/50 hover:border-primary/30 hover:bg-accent"
                      }`}
                    >
                      {nativeLanguage === lang.code && (
                        <motion.div 
                          layoutId="active-check"
                          className="absolute top-2 right-2 text-primary"
                        >
                          <Check size={20} />
                        </motion.div>
                      )}
                      <div className="flex flex-col items-center gap-4">
                        <span className="text-4xl group-hover:scale-125 transition-transform duration-300">
                          {lang.flag}
                        </span>
                        <span className="text-lg font-bold text-foreground">{lang.endonym}</span>
                      </div>
                    </Card>
                  </motion.div>
                ))}
              </div>

              <div className="flex justify-center pt-8">
                <Button
                  size="lg"
                  disabled={!nativeLanguage}
                  onClick={() => setStep(2)}
                  className="h-14 px-12 rounded-full text-lg font-bold group"
                >
                  {t("continue")}
                  <ArrowRight className="ml-2 group-hover:translate-x-1 transition-transform" />
                </Button>
              </div>
            </motion.div>
          )}

          {step === 2 && (
            <motion.div
              key="step2"
              initial={{ opacity: 0, x: 20 }}
              animate={{ opacity: 1, x: 0 }}
              exit={{ opacity: 0, x: -20 }}
              className="space-y-8"
            >
              <div className="text-center space-y-4">
                <motion.div 
                   initial={{ scale: 0 }}
                   animate={{ scale: 1 }}
                   className="inline-flex p-3 rounded-2xl bg-purple-500/20 text-purple-600 dark:text-purple-400 mb-4"
                >
                  <GraduationCap size={32} />
                </motion.div>
                <h1 className="text-4xl md:text-5xl font-black text-foreground">{t("chooseLearningTitle")}</h1>
                <p className="text-muted-foreground text-lg">{t("chooseLearningDesc")}</p>
              </div>

              <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
                {LEARNING_LANGUAGES.map((lang, idx) => (
                  <motion.div
                    key={lang.code}
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: idx * 0.05 }}
                  >
                    <Card
                      onClick={() => setTargetLanguage(lang.code)}
                      className={`p-6 cursor-pointer transition-all duration-300 group border-2 relative overflow-hidden ${
                        targetLanguage === lang.code 
                        ? "border-purple-500 bg-purple-500/10 scale-105 shadow-[0_0_30px_-10px_rgba(168,85,247,0.5)]" 
                        : "border-border bg-card/50 hover:border-purple-500/30 hover:bg-accent"
                      }`}
                    >
                      {targetLanguage === lang.code && (
                        <motion.div 
                          layoutId="active-check-target"
                          className="absolute top-2 right-2 text-purple-600 dark:text-purple-400"
                        >
                          <Check size={20} />
                        </motion.div>
                      )}
                      <div className="flex flex-col items-center gap-4">
                        <span className="text-4xl group-hover:scale-125 transition-transform duration-300">
                          {lang.flag}
                        </span>
                        <span className="text-lg font-bold text-foreground">{lang.endonym}</span>
                      </div>
                    </Card>
                  </motion.div>
                ))}
              </div>

              <div className="flex flex-col items-center gap-4 pt-8">
                <Button
                  size="lg"
                  disabled={!targetLanguage || loading}
                  onClick={handleFinish}
                  className="h-14 px-12 rounded-full text-lg font-bold bg-purple-600 hover:bg-purple-700 shadow-xl shadow-purple-500/20 group text-white"
                >
                  {loading ? t("initializing") : t("startLearning")}
                  {!loading && <ArrowRight className="ml-2 group-hover:translate-x-1 transition-transform" />}
                </Button>
                <button 
                  onClick={() => setStep(1)}
                  className="text-muted-foreground hover:text-foreground transition-colors"
                >
                  {t("goBack")}
                </button>
              </div>
            </motion.div>
          )}
        </AnimatePresence>
      </div>
    </div>
  );
}
