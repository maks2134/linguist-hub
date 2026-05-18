"use client";

import { useEffect, useState, useRef } from "react";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Badge } from "@/components/ui/badge";
import { useWebSocket } from "@/hooks/useWebSocket";
import { ChatMessage } from "@/types";
import { api, ChatHistoryMessage } from "@/lib/api";
import { toast } from "sonner";
import { ArrowLeft, Send, Loader2, Sparkles, BrainCircuit, GraduationCap, CheckCircle2, MessageSquare } from "lucide-react";
import { motion, AnimatePresence } from "framer-motion";

export default function DiagnosticPage({
  params,
}: {
  params: Promise<{ profileId: string }>;
}) {
  const router = useRouter();
  const [profileId, setProfileId] = useState<string | null>(null);
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [input, setInput] = useState("");
  const [loading, setLoading] = useState(false);
  const [assessmentComplete, setAssessmentComplete] = useState(false);
  const [level, setLevel] = useState<string | null>(null);
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
        const history = await api.getChatHistory(profileId, "diagnostic");
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

        if (response.is_assessment && response.level) {
          setAssessmentComplete(true);
          setLevel(response.level);
          toast.success(`Level assessed: ${response.level}`);
        }
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
        type: "diagnostic",
        message: messageContent,
        profile_id: parseInt(profileId),
      });
    } catch (error) {
      toast.error("Failed to send message");
      setLoading(false);
    }
  };

  const handleGeneratePlan = () => {
    if (!profileId || !connected) return;

    try {
      sendMessage({
        type: "generate_plan",
        profile_id: parseInt(profileId),
      });
      toast.success("Study plan generated!");
      setTimeout(() => router.push(`/plan/${profileId}`), 1000);
    } catch (error) {
      toast.error("Failed to generate plan");
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
          <Button variant="ghost" onClick={() => router.push("/dashboard")} className="rounded-full">
            <ArrowLeft className="w-4 h-4 mr-2" />
            Back to Dashboard
          </Button>
          {level && (
            <Badge variant="default" className="text-lg px-6 py-1.5 rounded-full bg-primary shadow-lg shadow-primary/20 border-none">
              <GraduationCap className="w-5 h-5 mr-2" />
              Level: {level}
            </Badge>
          )}
        </div>

        <Card className="flex-1 flex flex-col border-none shadow-2xl bg-white/70 dark:bg-slate-900/70 backdrop-blur-xl overflow-hidden min-h-[600px]">
          <CardHeader className="border-b border-slate-100 dark:border-slate-800 bg-white/50 dark:bg-slate-900/50">
            <div className="flex items-center gap-4">
               <div className="h-12 w-12 rounded-2xl bg-primary/10 flex items-center justify-center text-primary">
                  <BrainCircuit className="h-7 w-7" />
               </div>
               <div>
                  <CardTitle className="text-xl">AI Language Assessment</CardTitle>
                  <p className="text-sm text-muted-foreground flex items-center gap-1">
                    <Sparkles className="h-3 w-3 text-yellow-500" />
                    Chat with our AI to determine your level
                  </p>
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
                      className="text-center py-10"
                    >
                      <div className="max-w-sm mx-auto space-y-4">
                        <div className="h-16 w-16 rounded-full bg-slate-100 dark:bg-slate-800 flex items-center justify-center mx-auto text-slate-400">
                          <MessageSquare className="h-8 w-8" />
                        </div>
                        <h3 className="font-semibold text-lg">Start your assessment</h3>
                        <p className="text-muted-foreground text-sm">
                          Send a message to begin! Say hello and introduce yourself in the language you&apos;re learning.
                        </p>
                      </div>
                    </motion.div>
                  )}
                  {messages.map((msg, idx) => (
                    <motion.div
                      key={idx}
                      initial={{ opacity: 0, y: 10, scale: 0.95 }}
                      animate={{ opacity: 1, y: 0, scale: 1 }}
                      className={`flex ${
                        msg.role === "user" ? "justify-end" : "justify-start"
                      }`}
                    >
                      <div
                        className={`max-w-[85%] rounded-2xl px-5 py-3 shadow-sm ${
                          msg.role === "user"
                            ? "bg-primary text-primary-foreground rounded-tr-none"
                            : "bg-white dark:bg-slate-800 border border-slate-100 dark:border-slate-700 rounded-tl-none"
                        }`}
                      >
                        <p className="whitespace-pre-wrap leading-relaxed">{msg.content}</p>
                      </div>
                    </motion.div>
                  ))}
                  
                  {(streamingContent || (loading && !streamingContent)) && (
                    <motion.div
                      initial={{ opacity: 0, y: 10 }}
                      animate={{ opacity: 1, y: 0 }}
                      className="flex justify-start"
                    >
                      <div className="max-w-[85%] rounded-2xl px-5 py-3 bg-white dark:bg-slate-800 border border-slate-100 dark:border-slate-700 rounded-tl-none shadow-sm">
                        {streamingContent ? (
                          <div className="flex items-end gap-1">
                            <p className="whitespace-pre-wrap leading-relaxed">{streamingContent}</p>
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
              <AnimatePresence mode="wait">
                {!assessmentComplete ? (
                  <motion.div 
                    key="input"
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    exit={{ opacity: 0, y: -20 }}
                    className="relative"
                  >
                    <div className="flex gap-2">
                      <Input
                        value={input}
                        onChange={(e) => setInput(e.target.value)}
                        onKeyPress={(e) => e.key === "Enter" && handleSend()}
                        placeholder="Type your message..."
                        disabled={loading}
                        className="h-12 rounded-full px-6 bg-white dark:bg-slate-800 border-slate-200 dark:border-slate-700 focus-visible:ring-primary shadow-inner"
                      />
                      <Button 
                        onClick={handleSend} 
                        disabled={loading || !input.trim()}
                        className="h-12 w-12 rounded-full p-0 flex-shrink-0 shadow-lg shadow-primary/20"
                      >
                        <Send className="w-5 h-5" />
                      </Button>
                    </div>
                  </motion.div>
                ) : (
                  <motion.div 
                    key="complete"
                    initial={{ opacity: 0, scale: 0.9 }}
                    animate={{ opacity: 1, scale: 1 }}
                    className="space-y-4"
                  >
                    <div className="flex flex-col items-center text-center p-6 rounded-3xl bg-green-500/5 border border-green-500/20">
                      <div className="h-12 w-12 rounded-full bg-green-500 text-white flex items-center justify-center mb-4">
                        <CheckCircle2 className="w-7 h-7" />
                      </div>
                      <h3 className="text-xl font-bold text-green-700 dark:text-green-400">Assessment Complete!</h3>
                      <p className="text-muted-foreground mt-1">
                        We&apos;ve analyzed your skills and determined your level is <strong>{level}</strong>.
                      </p>
                      <Button onClick={handleGeneratePlan} size="lg" className="mt-6 w-full sm:w-auto px-10 rounded-full shadow-lg shadow-primary/20">
                        Generate My Study Plan
                      </Button>
                    </div>
                  </motion.div>
                )}
              </AnimatePresence>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}

