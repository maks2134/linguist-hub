"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { api } from "@/lib/api";
import { Profile } from "@/types";
import { LANGUAGES } from "@/lib/constants";
import { toast } from "sonner";
import { Plus, Sparkles, BookOpen, MessageSquare, Target, ChevronRight, Loader2 } from "lucide-react";
import { motion, AnimatePresence } from "framer-motion";

export default function DashboardPage() {
  const router = useRouter();
  const [profiles, setProfiles] = useState<Profile[]>([]);
  const [loading, setLoading] = useState(true);
  const [showLanguageSelector, setShowLanguageSelector] = useState(false);

  useEffect(() => {
    loadProfiles();
  }, []);

  const loadProfiles = async () => {
    try {
      const data = await api.get<Profile[]>("/api/profile/");
      setProfiles(data);
    } catch (error) {
      toast.error("Failed to load profiles");
    } finally {
      setLoading(false);
    }
  };

  const handleSelectLanguage = async (languageCode: string) => {
    try {
      const profile = await api.post<Profile>("/api/profile/select", {
        language: languageCode,
      });
      toast.success("Language profile created!");
      router.push(`/diagnostic/${profile.id}`);
    } catch (error) {
      toast.error("Failed to create profile");
    }
  };

  const getLevelColor = (level: string) => {
    const l = level.toUpperCase();
    if (l.startsWith("C")) return "bg-purple-500/10 text-purple-600 border-purple-500/20";
    if (l.startsWith("B")) return "bg-blue-500/10 text-blue-600 border-blue-500/20";
    if (l.startsWith("A")) return "bg-green-500/10 text-green-600 border-green-500/20";
    return "bg-slate-500/10 text-slate-600 border-slate-500/20";
  };

  if (loading) {
    return (
      <div className="flex h-[calc(100vh-4rem)] items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    );
  }

  return (
    <div className="min-h-full bg-[radial-gradient(ellipse_at_top,_var(--tw-gradient-stops))] from-indigo-50 via-white to-white dark:from-slate-900 dark:via-slate-950 dark:to-slate-950 px-4 py-8">
      <div className="container mx-auto max-w-6xl">
        <header className="mb-10">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5 }}
          >
            <h1 className="text-4xl font-extrabold tracking-tight lg:text-5xl mb-2 bg-gradient-to-r from-slate-900 to-slate-700 dark:from-white dark:to-slate-400 bg-clip-text text-transparent">
              Welcome back
            </h1>
            <p className="text-lg text-muted-foreground flex items-center gap-2">
              Select your path and continue your linguistic journey <Sparkles className="h-4 w-4 text-yellow-500" />
            </p>
          </motion.div>
        </header>

        <AnimatePresence mode="wait">
          {!showLanguageSelector ? (
            <motion.div
              key="profiles"
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              className="space-y-8"
            >
              {profiles.length > 0 ? (
                <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
                  {profiles.map((profile, index) => {
                    const lang = LANGUAGES.find((l) => l.code.toLowerCase() === profile.target_language.toLowerCase());
                    return (
                      <motion.div
                        key={profile.id}
                        initial={{ opacity: 0, y: 20 }}
                        animate={{ opacity: 1, y: 0 }}
                        transition={{ delay: index * 0.1 }}
                      >
                        <Card
                          className="group relative overflow-hidden border-none shadow-md hover:shadow-xl transition-all duration-300 bg-white/50 dark:bg-slate-900/50 backdrop-blur-sm ring-1 ring-slate-200 dark:ring-slate-800"
                        >
                          <div className="absolute top-0 right-0 p-4">
                             <Badge variant="outline" className={`${getLevelColor(profile.current_level)} border font-semibold`}>
                               {profile.current_level === "unknown" ? "?" : profile.current_level}
                             </Badge>
                          </div>
                          
                          <CardHeader className="pb-2">
                            <div className="text-4xl mb-4 group-hover:scale-110 transition-transform duration-300 w-fit">
                              {lang?.flag}
                            </div>
                            <CardTitle className="text-2xl font-bold">{lang?.name}</CardTitle>
                            <CardDescription className="flex items-center gap-2">
                              <Target className="h-3 w-3" />
                              Progress: {profile.progress_percentage.toFixed(0)}%
                            </CardDescription>
                          </CardHeader>
                          
                          <CardContent className="space-y-6">
                            <div className="space-y-2">
                              <div className="w-full bg-slate-100 dark:bg-slate-800 h-2 rounded-full overflow-hidden">
                                <motion.div
                                  initial={{ width: 0 }}
                                  animate={{ width: `${profile.progress_percentage}%` }}
                                  className="bg-gradient-to-r from-primary to-purple-500 h-full rounded-full"
                                />
                              </div>
                            </div>

                            <div className="grid grid-cols-2 gap-2 pt-2">
                              <Button 
                                variant="secondary" 
                                size="sm" 
                                className="w-full bg-slate-100 dark:bg-slate-800 hover:bg-slate-200 dark:hover:bg-slate-700"
                                onClick={() => router.push(`/tutor/${profile.id}`)}
                              >
                                <MessageSquare className="h-4 w-4 mr-2 text-primary" />
                                Tutor
                              </Button>
                              <Button 
                                variant="default" 
                                size="sm" 
                                className="w-full shadow-lg shadow-primary/20"
                                onClick={() => {
                                  if (profile.current_level === "unknown") {
                                    router.push(`/diagnostic/${profile.id}`);
                                  } else {
                                    router.push(`/plan/${profile.id}`);
                                  }
                                }}
                              >
                                {profile.current_level === "unknown" ? "Start" : "Continue"}
                                <ChevronRight className="h-4 w-4 ml-1" />
                              </Button>
                            </div>
                          </CardContent>
                        </Card>
                      </motion.div>
                    );
                  })}
                  
                  {/* Add New Card Slot */}
                  <motion.div
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    transition={{ delay: profiles.length * 0.1 }}
                  >
                    <button
                      onClick={() => setShowLanguageSelector(true)}
                      className="w-full h-full min-h-[250px] group flex flex-col items-center justify-center gap-4 rounded-xl border-2 border-dashed border-slate-200 dark:border-slate-800 hover:border-primary/50 hover:bg-primary/[0.02] transition-all duration-300"
                    >
                      <div className="flex h-12 w-12 items-center justify-center rounded-full bg-slate-50 dark:bg-slate-900 group-hover:bg-primary group-hover:text-primary-foreground transition-colors duration-300">
                        <Plus className="h-6 w-6" />
                      </div>
                      <div className="text-center">
                        <p className="font-semibold text-slate-900 dark:text-white">Add new language</p>
                        <p className="text-sm text-muted-foreground">Expand your horizons</p>
                      </div>
                    </button>
                  </motion.div>
                </div>
              ) : (
                <div className="text-center py-20 bg-white/50 dark:bg-slate-900/50 rounded-3xl border border-slate-200 dark:border-slate-800 backdrop-blur-sm">
                   <div className="flex h-20 w-20 items-center justify-center rounded-2xl bg-primary/10 text-primary mx-auto mb-6">
                     <BookOpen className="h-10 w-10" />
                   </div>
                   <h2 className="text-2xl font-bold mb-2">No language profiles yet</h2>
                   <p className="text-muted-foreground mb-8 max-w-sm mx-auto">
                     Pick a language you want to master and let our AI guide you through the process.
                   </p>
                   <Button size="lg" onClick={() => setShowLanguageSelector(true)} className="rounded-full px-8">
                     <Plus className="h-5 w-5 mr-2" />
                     Choose my first language
                   </Button>
                </div>
              )}
            </motion.div>
          ) : (
            <motion.div
              key="selector"
              initial={{ opacity: 0, scale: 0.95 }}
              animate={{ opacity: 1, scale: 1 }}
              exit={{ opacity: 0, scale: 0.95 }}
              className="max-w-4xl mx-auto"
            >
              <Card className="border-none shadow-2xl bg-white/80 dark:bg-slate-900/80 backdrop-blur-md">
                <CardHeader className="text-center border-b border-slate-100 dark:border-slate-800 pb-8">
                  <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-primary/10 text-primary mx-auto mb-4">
                    <Plus className="h-6 w-6" />
                  </div>
                  <CardTitle className="text-3xl font-bold">What&apos;s next?</CardTitle>
                  <CardDescription className="text-lg">Choose a language you&apos;d like to explore today</CardDescription>
                </CardHeader>
                <CardContent className="p-8">
                  <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-4">
                    {LANGUAGES.map((lang, index) => (
                      <motion.button
                        key={lang.code}
                        initial={{ opacity: 0, y: 10 }}
                        animate={{ opacity: 1, y: 0 }}
                        transition={{ delay: index * 0.05 }}
                        whileHover={{ y: -5 }}
                        whileTap={{ scale: 0.95 }}
                        onClick={() => handleSelectLanguage(lang.code)}
                        className="flex flex-col items-center gap-3 p-4 rounded-2xl bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 hover:border-primary hover:shadow-lg hover:shadow-primary/5 transition-all"
                      >
                        <span className="text-4xl">{lang.flag}</span>
                        <span className="font-medium text-sm">{lang.name}</span>
                      </motion.button>
                    ))}
                  </div>
                  <div className="flex justify-center mt-10">
                    <Button
                      variant="ghost"
                      size="lg"
                      className="rounded-full text-muted-foreground hover:text-foreground"
                      onClick={() => setShowLanguageSelector(false)}
                    >
                      Back to dashboard
                    </Button>
                  </div>
                </CardContent>
              </Card>
            </motion.div>
          )}
        </AnimatePresence>
      </div>
    </div>
  );
}

