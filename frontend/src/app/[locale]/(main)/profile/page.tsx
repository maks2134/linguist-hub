"use client";

import { useState, useEffect } from "react";
import { useRouter } from "@/i18n/routing";
import { useTranslations } from "next-intl";
import { motion } from "framer-motion";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardDescription, CardHeader, CardTitle, CardFooter } from "@/components/ui/card";
import { api } from "@/lib/api";
import { User, Profile } from "@/types";
import { toast } from "sonner";
import { User as UserIcon, Lock, Languages, Save, KeyRound, ArrowLeft, ChevronDown, Check } from "lucide-react";

const NATIVE_LANGUAGES = [
  { code: "Russian", name: "Русский", flag: "🇷🇺" },
  { code: "English", name: "English", flag: "🇬🇧" },
  { code: "Spanish", name: "Español", flag: "🇪🇸" },
  { code: "French", name: "Français", flag: "🇫🇷" },
  { code: "German", name: "Deutsch", flag: "🇩🇪" },
  { code: "Chinese", name: "中文", flag: "🇨🇳" },
  { code: "Italian", name: "Italiano", flag: "🇮🇹" },
  { code: "Japanese", name: "日本語", flag: "🇯🇵" },
  { code: "Portuguese", name: "Português", flag: "🇵🇹" },
  { code: "Arabic", name: "العربية", flag: "🇸🇦" },
];

export default function ProfilePage() {
  const router = useRouter();
  const t = useTranslations("Profile");
  const [user, setUser] = useState<User | null>(null);
  const [profiles, setProfiles] = useState<Profile[]>([]);
  const [nickname, setNickname] = useState("");
  const [oldPassword, setOldPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [activeTab, setActiveTab] = useState<"account" | "security">("account");
  const [nativeLangOpen, setNativeLangOpen] = useState<number | null>(null);
  const [savingNative, setSavingNative] = useState<number | null>(null);

  const handleUpdateNativeLanguage = async (profileId: number, lang: string) => {
    setSavingNative(profileId);
    try {
      const updated = await api.updateNativeLanguage(profileId, lang);
      setProfiles((prev) => prev.map((p) => p.id === profileId ? { ...p, native_language: lang } : p));
      toast.success(t("nativeUpdated"));
    } catch {
      toast.error("Failed to update native language");
    } finally {
      setSavingNative(null);
      setNativeLangOpen(null);
    }
  };

  useEffect(() => {
    const loadData = async () => {
      try {
        const userData = await api.get<User>("/api/auth/me");
        setUser(userData);
        setNickname(userData.nickname || "");

        const profilesData = await api.get<Profile[]>("/api/profile/");
        setProfiles(profilesData);
      } catch (error) {
        toast.error("Failed to load profile data");
      }
    };
    loadData();
  }, []);

  const handleUpdateProfile = async () => {
    if (!nickname) return;
    setLoading(true);
    try {
      await api.updateProfile(nickname);
      toast.success(t("profileUpdated"));
    } catch (error) {
      toast.error("Failed to update profile");
    } finally {
      setLoading(false);
    }
  };

  const handleChangePassword = async () => {
    if (newPassword !== confirmPassword) {
      toast.error(t("passwordsMismatch"));
      return;
    }
    setLoading(true);
    try {
      await api.changePassword(oldPassword, newPassword);
      toast.success(t("passwordChanged"));
      setOldPassword("");
      setNewPassword("");
      setConfirmPassword("");
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Failed to change password");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-full bg-slate-50 dark:bg-slate-950 py-12 px-4">
      <div className="container mx-auto max-w-4xl">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          className="space-y-8"
        >
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-4">
               <Button variant="ghost" onClick={() => router.back()} className="rounded-full h-10 w-10 p-0">
                 <ArrowLeft size={20} />
               </Button>
               <h1 className="text-3xl font-bold">{t("title")}</h1>
            </div>
          </div>

          <div className="w-full">
            <div className="flex gap-2 max-w-md mb-8 p-1 rounded-xl bg-slate-100 dark:bg-slate-800">
              {(["account", "security"] as const).map((tab) => (
                <button
                  key={tab}
                  onClick={() => setActiveTab(tab)}
                  className={`flex-1 flex items-center justify-center gap-2 py-2 px-4 rounded-lg text-sm font-medium transition-all ${
                    activeTab === tab
                      ? "bg-white dark:bg-slate-700 shadow text-foreground"
                      : "text-muted-foreground hover:text-foreground"
                  }`}
                >
                  {tab === "account" ? <UserIcon size={16} /> : <Lock size={16} />}
                  {tab === "account" ? t("tabAccount") : t("tabSecurity")}
                </button>
              ))}
            </div>

            {activeTab === "account" && (
              <Card className="border-none shadow-xl bg-white/70 dark:bg-slate-900/70 backdrop-blur-xl">
                <CardHeader>
                  <CardTitle>{t("publicProfile")}</CardTitle>
                  <CardDescription>{t("publicProfileDesc")}</CardDescription>
                </CardHeader>
                <CardContent className="space-y-6">
                  <div className="space-y-2">
                    <Label htmlFor="nickname">{t("nickname")}</Label>
                    <div className="relative">
                      <UserIcon className="absolute left-3 top-3 h-4 w-4 text-muted-foreground" />
                      <Input
                        id="nickname"
                        placeholder={t("nicknamePlaceholder")}
                        className="pl-10"
                        value={nickname}
                        onChange={(e) => setNickname(e.target.value)}
                      />
                    </div>
                  </div>

                  <div className="space-y-4">
                    <Label>{t("learningLanguages")}</Label>
                    <div className="grid gap-4">
                      {profiles.map((profile) => (
                        <div 
                          key={profile.id}
                          className="p-4 rounded-xl bg-slate-100 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 space-y-3"
                        >
                          <div className="flex items-center justify-between">
                            <div className="flex items-center gap-3">
                              <div className="h-10 w-10 rounded-full bg-primary/10 flex items-center justify-center text-primary">
                                <Languages size={20} />
                              </div>
                              <div>
                                <p className="font-bold capitalize">{profile.target_language}</p>
                                <p className="text-xs text-muted-foreground">{t("level")}: {profile.current_level}</p>
                              </div>
                            </div>
                          </div>

                          <div className="relative">
                            <p className="text-xs font-medium text-muted-foreground mb-1">{t("nativeLanguageLabel")}</p>
                            <button
                              onClick={() => setNativeLangOpen(nativeLangOpen === profile.id ? null : profile.id)}
                              className="w-full flex items-center justify-between px-3 py-2 rounded-lg border border-slate-200 dark:border-slate-600 bg-white dark:bg-slate-700 text-sm hover:border-primary transition-colors"
                            >
                              <span>
                                {NATIVE_LANGUAGES.find(l => l.code === profile.native_language)?.flag ?? "🌐"}{" "}
                                {profile.native_language}
                              </span>
                              <ChevronDown size={14} className={`transition-transform ${nativeLangOpen === profile.id ? "rotate-180" : ""}`} />
                            </button>

                            {nativeLangOpen === profile.id && (
                              <motion.div
                                initial={{ opacity: 0, y: -4 }}
                                animate={{ opacity: 1, y: 0 }}
                                className="absolute left-0 right-0 top-full mt-1 z-20 rounded-xl border border-slate-200 dark:border-slate-600 bg-white dark:bg-slate-800 shadow-xl overflow-hidden"
                              >
                                {NATIVE_LANGUAGES.map((lang) => (
                                  <button
                                    key={lang.code}
                                    disabled={savingNative === profile.id}
                                    onClick={() => handleUpdateNativeLanguage(profile.id, lang.code)}
                                    className={`w-full flex items-center gap-3 px-4 py-2.5 text-sm transition-colors hover:bg-slate-100 dark:hover:bg-slate-700 ${
                                      profile.native_language === lang.code ? "text-primary font-semibold" : "text-foreground"
                                    }`}
                                  >
                                    <span className="text-base">{lang.flag}</span>
                                    <span>{lang.name}</span>
                                    {profile.native_language === lang.code && <Check size={14} className="ml-auto" />}
                                  </button>
                                ))}
                              </motion.div>
                            )}
                          </div>
                        </div>
                      ))}
                    </div>
                  </div>
                </CardContent>
                <CardFooter>
                  <Button onClick={handleUpdateProfile} disabled={loading || !nickname}>
                    <Save className="mr-2 h-4 w-4" />
                    {t("saveChanges")}
                  </Button>
                </CardFooter>
              </Card>
            )}

            {activeTab === "security" && (
              <Card className="border-none shadow-xl bg-white/70 dark:bg-slate-900/70 backdrop-blur-xl">
                <CardHeader>
                  <CardTitle>{t("password")}</CardTitle>
                  <CardDescription>{t("passwordDesc")}</CardDescription>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div className="space-y-2">
                    <Label htmlFor="old-password">{t("currentPassword")}</Label>
                    <Input
                      id="old-password"
                      type="password"
                      value={oldPassword}
                      onChange={(e) => setOldPassword(e.target.value)}
                    />
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="new-password">{t("newPassword")}</Label>
                    <Input
                      id="new-password"
                      type="password"
                      value={newPassword}
                      onChange={(e) => setNewPassword(e.target.value)}
                    />
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="confirm-password">{t("confirmNewPassword")}</Label>
                    <Input
                      id="confirm-password"
                      type="password"
                      value={confirmPassword}
                      onChange={(e) => setConfirmPassword(e.target.value)}
                    />
                  </div>
                </CardContent>
                <CardFooter>
                  <Button onClick={handleChangePassword} disabled={loading || !oldPassword || !newPassword}>
                    <KeyRound className="mr-2 h-4 w-4" />
                    {t("updatePassword")}
                  </Button>
                </CardFooter>
              </Card>
            )}
          </div>
        </motion.div>
      </div>
    </div>
  );
}

function Badge({ children, variant, className }: any) {
  return (
    <span className={`px-2 py-1 rounded text-xs font-medium ${
      variant === "outline" ? "border border-slate-200 dark:border-slate-700" : "bg-primary text-white"
    } ${className}`}>
      {children}
    </span>
  );
}
