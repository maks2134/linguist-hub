"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { api } from "@/lib/api";
import { StudyPlan } from "@/types";
import { toast } from "sonner";
import { ArrowLeft, BookOpen } from "lucide-react";

export default function PlanPage({
  params,
}: {
  params: Promise<{ profileId: string }>;
}) {
  const router = useRouter();
  const [profileId, setProfileId] = useState<string | null>(null);
  const [plan, setPlan] = useState<StudyPlan | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    params.then((p) => {
      setProfileId(p.profileId);
      loadPlan(p.profileId);
    });
  }, [params]);

  const loadPlan = async (id: string) => {
    try {
      const data = await api.get<StudyPlan>(`/api/ai/plan/${id}`);
      setPlan(data);
    } catch (error) {
      toast.error("Failed to load study plan");
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <p>Loading plan...</p>
      </div>
    );
  }

  if (!plan) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <Card className="max-w-md">
          <CardHeader>
            <CardTitle>No Study Plan</CardTitle>
            <CardDescription>You need to complete the diagnostic first</CardDescription>
          </CardHeader>
          <CardContent>
            <Button onClick={() => router.push("/dashboard")}>
              Back to Dashboard
            </Button>
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 via-white to-purple-50 dark:from-gray-900 dark:via-gray-800 dark:to-gray-900">
      <div className="container mx-auto py-8 px-4 max-w-6xl">
        <div className="mb-6">
          <Button variant="ghost" onClick={() => router.push("/dashboard")}>
            <ArrowLeft className="w-4 h-4 mr-2" />
            Back to Dashboard
          </Button>
        </div>

        <div className="mb-8">
          <h1 className="text-4xl font-bold mb-2">Your Study Plan</h1>
          <p className="text-muted-foreground">
            AI-generated curriculum tailored to your level
          </p>
        </div>

        <div className="grid gap-6 md:grid-cols-2">
          {plan.modules.map((module) => (
            <Card key={module.id} className="hover:shadow-lg transition-shadow">
              <CardHeader>
                <div className="flex items-start justify-between">
                  <div>
                    <CardTitle className="text-xl">
                      Module {module.id}: {module.title}
                    </CardTitle>
                    <CardDescription className="mt-2">
                      {module.description}
                    </CardDescription>
                  </div>
                  <Badge variant="secondary">{module.topics.length} topics</Badge>
                </div>
              </CardHeader>
              <CardContent>
                <div className="space-y-2">
                  <p className="text-sm font-semibold mb-2">Topics:</p>
                  <div className="flex flex-wrap gap-2">
                    {module.topics.map((topic, idx) => (
                      <Button
                        key={idx}
                        variant="outline"
                        size="sm"
                        onClick={() =>
                          router.push(`/tutor/${profileId}?topic=${encodeURIComponent(topic)}`)
                        }
                      >
                        <BookOpen className="w-3 h-3 mr-1" />
                        {topic}
                      </Button>
                    ))}
                  </div>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      </div>
    </div>
  );
}
