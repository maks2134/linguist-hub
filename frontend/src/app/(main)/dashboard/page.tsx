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
import { Plus, LogOut } from "lucide-react";

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

  const handleLogout = () => {
    localStorage.removeItem("token");
    router.push("/login");
  };

  const getLevelColor = (level: string) => {
    if (level === "unknown") return "secondary";
    if (level.startsWith("A")) return "default";
    if (level.startsWith("B")) return "default";
    return "default";
  };

  if (loading) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <p>Loading...</p>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 via-white to-purple-50 dark:from-gray-900 dark:via-gray-800 dark:to-gray-900">
      <div className="container mx-auto py-8 px-4">
        <div className="flex items-center justify-between mb-8">
          <div>
            <h1 className="text-4xl font-bold">LinguistAI</h1>
            <p className="text-muted-foreground mt-1">Your AI language learning hub</p>
          </div>
          <Button variant="outline" onClick={handleLogout}>
            <LogOut className="w-4 h-4 mr-2" />
            Logout
          </Button>
        </div>

        {!showLanguageSelector && (
          <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3 mb-6">
            {profiles.map((profile) => {
              const lang = LANGUAGES.find((l) => l.code === profile.target_language);
              return (
                <Card
                  key={profile.id}
                  className="cursor-pointer hover:shadow-lg transition-shadow"
                  onClick={() => {
                    if (profile.current_level === "unknown") {
                      router.push(`/diagnostic/${profile.id}`);
                    } else {
                      router.push(`/plan/${profile.id}`);
                    }
                  }}
                >
                  <CardHeader>
                    <CardTitle className="flex items-center justify-between">
                      <span>
                        {lang?.flag} {lang?.name}
                      </span>
                      <Badge variant={getLevelColor(profile.current_level)}>
                        {profile.current_level}
                      </Badge>
                    </CardTitle>
                    <CardDescription>
                      Progress: {profile.progress_percentage.toFixed(0)}%
                    </CardDescription>
                  </CardHeader>
                  <CardContent>
                    <div className="w-full bg-secondary h-2 rounded-full overflow-hidden">
                      <div
                        className="bg-primary h-full transition-all"
                        style={{ width: `${profile.progress_percentage}%` }}
                      />
                    </div>
                  </CardContent>
                </Card>
              );
            })}
          </div>
        )}

        {!showLanguageSelector ? (
          <Button
            onClick={() => setShowLanguageSelector(true)}
            size="lg"
            className="w-full md:w-auto"
          >
            <Plus className="w-4 h-4 mr-2" />
            Add New Language
          </Button>
        ) : (
          <Card>
            <CardHeader>
              <CardTitle>Select a Language</CardTitle>
              <CardDescription>Choose a language you want to learn</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-3">
                {LANGUAGES.map((lang) => (
                  <Button
                    key={lang.code}
                    variant="outline"
                    className="h-20 flex flex-col items-center justify-center"
                    onClick={() => handleSelectLanguage(lang.code)}
                  >
                    <span className="text-3xl mb-1">{lang.flag}</span>
                    <span className="text-sm">{lang.name}</span>
                  </Button>
                ))}
              </div>
              <Button
                variant="ghost"
                className="mt-4"
                onClick={() => setShowLanguageSelector(false)}
              >
                Cancel
              </Button>
            </CardContent>
          </Card>
        )}
      </div>
    </div>
  );
}
