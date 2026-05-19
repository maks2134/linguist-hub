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
  progress_percentage: number;
}

export interface Module {
  id: number;
  title: string;
  description: string;
  topics: string[];
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
