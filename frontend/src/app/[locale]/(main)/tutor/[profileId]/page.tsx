"use client";

import { useEffect, useState, useRef } from "react";
import { useRouter } from "@/i18n/routing";
import { useTranslations } from "next-intl";
import { useSearchParams } from "next/navigation";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { ScrollArea } from "@/components/ui/scroll-area";
import { useWebSocket } from "@/hooks/useWebSocket";
import { ChatMessage } from "@/types";
import { api, ChatHistoryMessage } from "@/lib/api";
import { toast } from "sonner";
import { ArrowLeft, Send, Loader2, Sparkles, Languages, MessageSquare, BookOpen, UserCircle2 } from "lucide-react";
import { motion, AnimatePresence } from "framer-motion";

export default function TutorPage({
  params,
}: {
  params: Promise<{ profileId: string }>;
}) {
  const router = useRouter();
  const t = useTranslations("Tutor");
  const tToast = useTranslations("Toast");
  const searchParams = useSearchParams();
  const topic = searchParams.get("topic") || t("defaultTopic");
  const [profileId, setProfileId] = useState<string | null>(null);
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [input, setInput] = useState("");
  const [loading, setLoading] = useState(false);
  const [streamingContent, setStreamingContent] = useState("");
  const { connected, sendMessage, onMessage } = useWebSocket();
  const scrollRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    params.then((p) => setProfileId(p.profileId));
  }, [params]);

  useEffect(() => {
    if (!profileId) return;

    const loadHistory = async () => {
      try {
        const history = await api.getChatHistory(profileId, "tutor");
        const formattedMessages: ChatMessage[] = history.map((msg: ChatHistoryMessage) => ({
          role: msg.role,
          content: msg.content,
        }));
        setMessages(formattedMessages);
      } catch (error) {
        console.error("Failed to load chat history:", error);
      }
    };

    loadHistory();
  }, [profileId]);

  useEffect(() => {
    if (!connected) return;

    const unsubscribe = onMessage((response) => {
      if (response.error) {
        toast.error(response.error);
        setLoading(false);
        return;
      }

      if (response.chunk) {
        setStreamingContent((prev) => prev + response.chunk);
      }

      if (response.done) {
        const assistantMessage: ChatMessage = {
          role: "assistant",
          content: response.content || streamingContent,
        };
        setMessages((prev) => [...prev, assistantMessage]);
        setStreamingContent("");
        setLoading(false);
      }
    });

    return () => unsubscribe();
  }, [connected, onMessage, streamingContent]);

  const handleSend = () => {
    if (!input.trim() || !profileId || !connected) return;

    const userMessage: ChatMessage = {
      role: "user",
      content: input,
    };

    setMessages((prev) => [...prev, userMessage]);
    const messageContent = input;
    setInput("");
    setLoading(true);
    setStreamingContent("");

    try {
      sendMessage({
        type: "tutor",
        message: messageContent,
        topic: topic,
        profile_id: parseInt(profileId),
      });
    } catch (error) {
      toast.error(tToast("sendMessageFailed"));
      setLoading(false);
    }
  };

  useEffect(() => {
    if (scrollRef.current) {
      scrollRef.current.scrollIntoView({ behavior: "smooth" });
    }
  }, [messages, streamingContent]);

  return (
    <div className="min-h-full bg-[radial-gradient(ellipse_at_top,_var(--tw-gradient-stops))] from-indigo-50 via-white to-white dark:from-slate-900 dark:via-slate-950 dark:to-slate-950 py-8 px-4">
      <div className="container mx-auto max-w-4xl h-full flex flex-col">
        <div className="mb-6 flex items-center justify-between">
          <Button variant="ghost" onClick={() => router.back()} className="rounded-full">
            <ArrowLeft className="w-4 h-4 mr-2" />
            {t("back")}
          </Button>
          <div className="flex items-center gap-2">
             <div className="h-8 px-4 rounded-full bg-slate-100 dark:bg-slate-800 flex items-center gap-2 text-sm font-medium">
                <BookOpen className="h-4 w-4 text-primary" />
                {t("topic")}: {topic}
             </div>
          </div>
        </div>

        <Card className="flex-1 flex flex-col border-none shadow-2xl bg-white/70 dark:bg-slate-900/70 backdrop-blur-xl overflow-hidden min-h-[600px]">
          <CardHeader className="border-b border-slate-100 dark:border-slate-800 bg-white/50 dark:bg-slate-900/50 py-4">
            <div className="flex items-center gap-4">
               <div className="h-10 w-10 rounded-xl bg-primary/10 flex items-center justify-center text-primary">
                  <Languages className="h-6 w-6" />
               </div>
               <div>
                  <CardTitle className="text-lg">{t("aiTutor")}</CardTitle>
                  <div className="flex items-center gap-2">
                     <span className="flex h-2 w-2 rounded-full bg-green-500 animate-pulse" />
                     <span className="text-xs text-muted-foreground uppercase tracking-wider font-semibold">{t("activeSession")}</span>
                  </div>
               </div>
            </div>
          </CardHeader>
          
          <CardContent className="flex-1 flex flex-col p-0 overflow-hidden relative">
            <ScrollArea className="flex-1 px-6 py-6">
              <div className="space-y-6">
                <AnimatePresence initial={false}>
                  {messages.length === 0 && (
                    <motion.div 
                      initial={{ opacity: 0, y: 10 }}
                      animate={{ opacity: 1, y: 0 }}
                      className="text-center py-20"
                    >
                      <div className="max-w-sm mx-auto space-y-4">
                        <div className="h-20 w-20 rounded-3xl bg-primary/5 flex items-center justify-center mx-auto text-primary">
                          <MessageSquare className="h-10 w-10" />
                        </div>
                        <h3 className="font-bold text-xl tracking-tight">{t("readyToPractice")}</h3>
                        <p className="text-muted-foreground">
                          {t("readyToPracticeDesc", { topic })}
                        </p>
                      </div>
                    </motion.div>
                  )}
                  {messages.map((msg, idx) => (
                    <motion.div
                      key={idx}
                      initial={{ opacity: 0, y: 10, scale: 0.98 }}
                      animate={{ opacity: 1, y: 0, scale: 1 }}
                      className={`flex items-start gap-3 ${
                        msg.role === "user" ? "flex-row-reverse" : "flex-row"
                      }`}
                    >
                      <div className={`flex-shrink-0 h-8 w-8 rounded-full flex items-center justify-center ${
                         msg.role === "user" ? "bg-slate-200 dark:bg-slate-700" : "bg-primary/10 text-primary"
                      }`}>
                         {msg.role === "user" ? <UserCircle2 className="h-5 w-5" /> : <Sparkles className="h-4 w-4" />}
                      </div>
                      <div
                        className={`max-w-[80%] rounded-2xl px-5 py-3 shadow-sm ${
                          msg.role === "user"
                            ? "bg-primary text-primary-foreground rounded-tr-none"
                            : "bg-white dark:bg-slate-800 border border-slate-100 dark:border-slate-700 rounded-tl-none"
                        }`}
                      >
                        <p className="whitespace-pre-wrap leading-relaxed text-[15px]">{msg.content}</p>
                      </div>
                    </motion.div>
                  ))}
                  
                  {(streamingContent || (loading && !streamingContent)) && (
                    <motion.div
                      initial={{ opacity: 0, y: 10 }}
                      animate={{ opacity: 1, y: 0 }}
                      className="flex items-start gap-3"
                    >
                      <div className="flex-shrink-0 h-8 w-8 rounded-full bg-primary/10 text-primary flex items-center justify-center">
                         <Sparkles className="h-4 w-4" />
                      </div>
                      <div className="max-w-[80%] rounded-2xl px-5 py-3 bg-white dark:bg-slate-800 border border-slate-100 dark:border-slate-700 rounded-tl-none shadow-sm">
                        {streamingContent ? (
                          <div className="flex items-end gap-1">
                            <p className="whitespace-pre-wrap leading-relaxed text-[15px]">{streamingContent}</p>
                            <span className="inline-block w-1.5 h-4 bg-primary animate-pulse rounded-full mb-1" />
                          </div>
                        ) : (
                          <div className="flex gap-1.5 py-1 px-1">
                            <span className="w-2 h-2 bg-primary/40 rounded-full animate-bounce" style={{ animationDelay: '0ms' }} />
                            <span className="w-2 h-2 bg-primary/40 rounded-full animate-bounce" style={{ animationDelay: '150ms' }} />
                            <span className="w-2 h-2 bg-primary/40 rounded-full animate-bounce" style={{ animationDelay: '300ms' }} />
                          </div>
                        )}
                      </div>
                    </motion.div>
                  )}
                </AnimatePresence>
                <div ref={scrollRef} className="h-4" />
              </div>
            </ScrollArea>

            <div className="p-6 bg-gradient-to-t from-white/80 via-white/80 to-transparent dark:from-slate-900/80 dark:via-slate-900/80 backdrop-blur-sm">
              <div className="relative group">
                <div className="absolute -inset-1 bg-gradient-to-r from-primary/20 to-purple-500/20 rounded-full blur opacity-25 group-focus-within:opacity-100 transition duration-1000 group-focus-within:duration-200"></div>
                <div className="relative flex gap-2">
                  <Input
                    value={input}
                    onChange={(e) => setInput(e.target.value)}
                    onKeyPress={(e) => e.key === "Enter" && handleSend()}
                    placeholder={t("messagePlaceholder")}
                    disabled={loading}
                    className="h-14 rounded-full px-8 bg-white dark:bg-slate-800 border-slate-200 dark:border-slate-700 focus-visible:ring-primary shadow-xl border-none outline-none ring-1 ring-slate-200 dark:ring-slate-800"
                  />
                  <Button 
                    onClick={handleSend} 
                    disabled={loading || !input.trim()}
                    className="h-14 w-14 rounded-full p-0 flex-shrink-0 shadow-xl shadow-primary/20"
                  >
                    <Send className="w-6 h-6" />
                  </Button>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}

