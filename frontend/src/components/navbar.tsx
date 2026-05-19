"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useTheme } from "next-themes";
import { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { LogOut, Languages, User, Sun, Moon, Monitor } from "lucide-react";
import { motion, AnimatePresence } from "framer-motion";

export function Navbar() {
  const router = useRouter();
  const { theme, setTheme } = useTheme();
  const [themeOpen, setThemeOpen] = useState(false);
  const [mounted, setMounted] = useState(false);

  useEffect(() => setMounted(true), []);

  const themes = [
    { key: "light", label: "Light", icon: Sun },
    { key: "dark", label: "Dark", icon: Moon },
    { key: "system", label: "System", icon: Monitor },
  ];

  const handleLogout = () => {
    localStorage.removeItem("token");
    router.push("/login");
  };

  return (
    <motion.nav 
      initial={{ y: -20, opacity: 0 }}
      animate={{ y: 0, opacity: 1 }}
      className="sticky top-0 z-50 w-full border-b bg-background/80 backdrop-blur-xl supports-[backdrop-filter]:bg-background/60"
    >
      <div className="container flex h-16 items-center justify-between px-4">
        <div className="flex items-center gap-2">
          <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-primary text-primary-foreground shadow-lg shadow-primary/20">
            <Languages className="h-6 w-6" />
          </div>
          <Link href="/dashboard" className="flex items-center space-x-2">
            <span className="text-xl font-bold tracking-tight bg-gradient-to-r from-primary to-purple-600 bg-clip-text text-transparent">
              LinguistAI
            </span>
          </Link>
        </div>

        <div className="flex items-center gap-4">
          <Button 
            variant="ghost" 
            size="icon" 
            className="rounded-full"
            onClick={() => router.push("/profile")}
          >
            <User className="h-5 w-5" />
          </Button>
          <div className="relative">
            <Button
              variant="ghost"
              size="icon"
              className="rounded-full"
              onClick={() => setThemeOpen((v) => !v)}
              aria-label="Toggle theme"
            >
              {mounted ? (
                theme === "dark" ? <Moon className="h-5 w-5" /> :
                theme === "light" ? <Sun className="h-5 w-5" /> :
                <Monitor className="h-5 w-5" />
              ) : <Monitor className="h-5 w-5" />}
            </Button>
            <AnimatePresence>
              {themeOpen && (
                <motion.div
                  initial={{ opacity: 0, scale: 0.95, y: -4 }}
                  animate={{ opacity: 1, scale: 1, y: 0 }}
                  exit={{ opacity: 0, scale: 0.95, y: -4 }}
                  transition={{ duration: 0.15 }}
                  className="absolute right-0 top-12 z-50 w-36 rounded-2xl border bg-background/95 backdrop-blur-xl shadow-xl overflow-hidden"
                  onMouseLeave={() => setThemeOpen(false)}
                >
                  {themes.map(({ key, label, icon: Icon }) => (
                    <button
                      key={key}
                      onClick={() => { setTheme(key); setThemeOpen(false); }}
                      className={`w-full flex items-center gap-3 px-4 py-2.5 text-sm transition-colors hover:bg-accent ${theme === key ? "text-primary font-semibold" : "text-muted-foreground"}`}
                    >
                      <Icon className="h-4 w-4 flex-shrink-0" />
                      {label}
                    </button>
                  ))}
                </motion.div>
              )}
            </AnimatePresence>
          </div>
          <div className="h-6 w-px bg-border mx-2" />
          <Button 
            variant="ghost" 
            size="sm" 
            onClick={handleLogout}
            className="text-muted-foreground hover:text-destructive hover:bg-destructive/10 transition-colors"
          >
            <LogOut className="mr-2 h-4 w-4" />
            Logout
          </Button>
        </div>
      </div>
    </motion.nav>
  );
}
