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
import { toast } from "sonner";
import { ArrowLeft, Send, Loader2 } from "lucide-react";

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
    <div className="min-h-screen bg-gradient-to-br from-blue-50 via-white to-purple-50 dark:from-gray-900 dark:via-gray-800 dark:to-gray-900">
      <div className="container mx-auto py-8 px-4 max-w-4xl">
        <div className="mb-6 flex items-center justify-between">
          <Button variant="ghost" onClick={() => router.push("/dashboard")}>
            <ArrowLeft className="w-4 h-4 mr-2" />
            Back to Dashboard
          </Button>
          {level && (
            <Badge variant="default" className="text-lg px-4 py-2">
              Level: {level}
            </Badge>
          )}
        </div>

        <Card className="h-[calc(100vh-12rem)]">
          <CardHeader>
            <CardTitle>AI Language Assessment</CardTitle>
            <p className="text-sm text-muted-foreground">
              Chat with the AI to determine your language level
            </p>
          </CardHeader>
          <CardContent className="flex flex-col h-[calc(100%-5rem)]">
            <ScrollArea className="flex-1 pr-4">
              <div className="space-y-4">
                {messages.map((msg, idx) => (
                  <div
                    key={idx}
                    className={`flex ${
                      msg.role === "user" ? "justify-end" : "justify-start"
                    }`}
                  >
                    <div
                      className={`max-w-[80%] rounded-lg px-4 py-2 ${
                        msg.role === "user"
                          ? "bg-primary text-primary-foreground"
                          : "bg-muted"
                      }`}
                    >
                      <p className="whitespace-pre-wrap">{msg.content}</p>
                    </div>
                  </div>
                ))}
                {streamingContent && (
                  <div className="flex justify-start">
                    <div className="bg-muted rounded-lg px-4 py-2">
                      <p className="whitespace-pre-wrap">{streamingContent}</p>
                      <span className="inline-block w-2 h-4 bg-current animate-pulse ml-1" />
                    </div>
                  </div>
                )}
                {loading && !streamingContent && (
                  <div className="flex justify-start">
                    <div className="bg-muted rounded-lg px-4 py-2">
                      <Loader2 className="w-4 h-4 animate-spin" />
                    </div>
                  </div>
                )}
                <div ref={scrollRef} />
              </div>
            </ScrollArea>

            {!assessmentComplete ? (
              <div className="flex gap-2 mt-4">
                <Input
                  value={input}
                  onChange={(e) => setInput(e.target.value)}
                  onKeyPress={(e) => e.key === "Enter" && handleSend()}
                  placeholder="Type your message..."
                  disabled={loading}
                />
                <Button onClick={handleSend} disabled={loading || !input.trim()}>
                  <Send className="w-4 h-4" />
                </Button>
              </div>
            ) : (
              <div className="mt-4 space-y-2">
                <p className="text-center text-sm text-muted-foreground">
                  Assessment complete! Your level: <strong>{level}</strong>
                </p>
                <Button onClick={handleGeneratePlan} className="w-full">
                  Generate Study Plan
                </Button>
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
