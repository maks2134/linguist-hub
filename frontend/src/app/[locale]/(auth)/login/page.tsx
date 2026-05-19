"use client";

import { useState } from "react";
import { useRouter, Link } from "@/i18n/routing";
import { useTranslations } from "next-intl";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { api } from "@/lib/api";
import { AuthResponse } from "@/types";
import { toast } from "sonner";

export default function LoginPage() {
  const router = useRouter();
  const t = useTranslations("Auth");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);

    try {
      const response = await api.post<AuthResponse>("/api/auth/login", {
        email,
        password,
      });

      localStorage.setItem("token", response.token);
      toast.success(t("loginSuccess"));
      router.push("/dashboard");
    } catch (error) {
      toast.error(error instanceof Error ? error.message : t("loginFailed"));
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex min-h-screen items-center justify-center bg-gradient-to-br from-blue-50 via-white to-purple-50 dark:from-gray-900 dark:via-gray-800 dark:to-gray-900 p-4">
      <Card className="w-full max-w-md">
        <CardHeader className="space-y-1">
          <CardTitle className="text-2xl font-bold">{t("welcomeBack")}</CardTitle>
          <CardDescription>
            {t("loginDesc")}
          </CardDescription>
        </CardHeader>
        <form onSubmit={handleSubmit}>
          <CardContent className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="email">{t("email")}</Label>
              <Input
                id="email"
                type="email"
                placeholder="m@example.com"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="password">{t("password")}</Label>
              <Input
                id="password"
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
              />
            </div>
          </CardContent>
          <CardFooter className="flex flex-col space-y-4">
            <Button type="submit" className="w-full" disabled={loading}>
              {loading ? t("loggingIn") : t("logIn")}
            </Button>
            <p className="text-sm text-center text-muted-foreground">
              {t("dontHaveAccount")}{" "}
              <Link href="/register" className="text-primary hover:underline">
                {t("signUp")}
              </Link>
            </p>
          </CardFooter>
        </form>
      </Card>
    </div>
  );
}
