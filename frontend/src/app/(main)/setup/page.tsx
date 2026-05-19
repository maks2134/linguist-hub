"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { motion, AnimatePresence } from "framer-motion";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { api } from "@/lib/api";
import { Profile } from "@/types";
import { toast } from "sonner";
import { Globe, GraduationCap, ArrowRight, Check } from "lucide-react";

const NATIVE_LANGUAGES = [
  { code: "Russian", name: "Русский", flag: "🇷🇺" },
  { code: "English", name: "English", flag: "🇺🇸" },
  { code: "Spanish", name: "Español", flag: "🇪🇸" },
  { code: "French", name: "Français", flag: "🇫🇷" },
  { code: "German", name: "Deutsch", flag: "🇩🇪" },
  { code: "Chinese", name: "中文", flag: "🇨🇳" },
];

const TARGET_LANGUAGES = [
  { code: "English", name: "English", flag: "🇺🇸" },
  { code: "Spanish", name: "Spanish", flag: "🇪🇸" },
  { code: "French", name: "French", flag: "🇫🇷" },
  { code: "German", name: "German", flag: "🇩🇪" },
  { code: "Italian", name: "Italian", flag: "🇮🇹" },
  { code: "Japanese", name: "Japanese", flag: "🇯🇵" },
];

export default function SetupPage() {
  const router = useRouter();
  const [step, setStep] = useState(1);
  const [nativeLanguage, setNativeLanguage] = useState("");
  const [targetLanguage, setTargetLanguage] = useState("");
  const [loading, setLoading] = useState(false);

  const handleFinish = async () => {
    if (!nativeLanguage || !targetLanguage) return;
    setLoading(true);
    try {
      const profile = await api.post<Profile>("/api/profile/select", {
        language: targetLanguage,
        native_language: nativeLanguage,
      });
      toast.success("Profile created! Starting onboarding...");
      router.push(`/diagnostic/${profile.id}`);
    } catch (error) {
      toast.error("Failed to setup profile");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-slate-950 flex flex-col items-center justify-center p-4 relative overflow-hidden">
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
                <h1 className="text-4xl md:text-5xl font-black text-white">Choose your native language</h1>
                <p className="text-slate-400 text-lg">We&apos;ll use this to explain complex topics and guide you.</p>
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
                        : "border-white/5 bg-white/5 hover:border-white/20 hover:bg-white/10"
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
                        <span className="text-lg font-bold text-white">{lang.name}</span>
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
                  Continue
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
                   className="inline-flex p-3 rounded-2xl bg-purple-500/20 text-purple-400 mb-4"
                >
                  <GraduationCap size={32} />
                </motion.div>
                <h1 className="text-4xl md:text-5xl font-black text-white">What language do you want to master?</h1>
                <p className="text-slate-400 text-lg">Your AI tutor will build a path specifically for this goal.</p>
              </div>

              <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
                {TARGET_LANGUAGES.map((lang, idx) => (
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
                        : "border-white/5 bg-white/5 hover:border-white/20 hover:bg-white/10"
                      }`}
                    >
                      {targetLanguage === lang.code && (
                        <motion.div 
                          layoutId="active-check-target"
                          className="absolute top-2 right-2 text-purple-400"
                        >
                          <Check size={20} />
                        </motion.div>
                      )}
                      <div className="flex flex-col items-center gap-4">
                        <span className="text-4xl group-hover:scale-125 transition-transform duration-300">
                          {lang.flag}
                        </span>
                        <span className="text-lg font-bold text-white">{lang.name}</span>
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
                  className="h-14 px-12 rounded-full text-lg font-bold bg-purple-600 hover:bg-purple-700 shadow-xl shadow-purple-500/20 group"
                >
                  {loading ? "Initializing..." : "Start Learning"}
                  {!loading && <ArrowRight className="ml-2 group-hover:translate-x-1 transition-transform" />}
                </Button>
                <button 
                  onClick={() => setStep(1)}
                  className="text-slate-500 hover:text-slate-300 transition-colors"
                >
                  Go back
                </button>
              </div>
            </motion.div>
          )}
        </AnimatePresence>
      </div>
    </div>
  );
}
