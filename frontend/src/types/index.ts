export interface User {
  id: number;
  email: string;
  nickname?: string;
}

export interface AuthResponse {
  token: string;
  user: User;
}

export interface Profile {
  id: number;
  user_id: number;
  target_language: string;
  native_language: string;
  current_level: string;
  onboarding_phase: string;
  goals?: string;
  weaknesses?: string[];
  progress_percentage: number;
}

export interface Topic {
  name: string;
  completed: boolean;
  locked: boolean;
}

export interface Module {
  id: number;
  title: string;
  description: string;
  topics: Topic[];
  locked?: boolean;
}

export interface StudyPlan {
  id: number;
  profile_id: number;
  modules: Module[];
}

export interface ChatMessage {
  role: "user" | "assistant";
  content: string;
}

export interface ChatResponse {
  message: string;
  level?: string;
  weaknesses?: string[];
  is_assessment?: boolean;
}

export interface LessonScore {
  fluency: number;
  grammar: number;
  vocabulary: number;
  pronunciation: number;
  notes: string;
}

export interface VoiceLessonHistory {
  lesson_id: number;
  module_id: number;
  topic: string;
  status: string;
  score?: LessonScore;
  started_at: string;
  completed_at?: string;
}
