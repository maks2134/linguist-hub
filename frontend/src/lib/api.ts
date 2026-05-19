const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

interface ApiResponse<T> {
  success: boolean;
  data?: T;
  error?: {
    code: string;
    message: string;
  };
}

export interface ChatHistoryMessage {
  id: number;
  user_id: number;
  profile_id: number;
  role: "user" | "assistant";
  context: string;
  content: string;
  created_at: string;
}

class ApiClient {
  private baseUrl: string;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }

  private getToken(): string | null {
    if (typeof window === "undefined") return null;
    return localStorage.getItem("token");
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const token = this.getToken();
    const headers: Record<string, string> = {
      "Content-Type": "application/json",
      ...(options.headers as Record<string, string>),
    };

    if (token) {
      headers["Authorization"] = `Bearer ${token}`;
    }

    const response = await fetch(`${this.baseUrl}${endpoint}`, {
      ...options,
      headers,
    });

    const data: ApiResponse<T> = await response.json();

    if (!response.ok || !data.success) {
      if (response.status === 401) {
        if (typeof window !== "undefined") {
          localStorage.removeItem("token");
          window.location.href = "/login";
        }
      }
      throw new Error(data.error?.message || "API request failed");
    }

    return data.data as T;
  }

  async get<T>(endpoint: string): Promise<T> {
    return this.request<T>(endpoint, { method: "GET" });
  }

  async post<T>(endpoint: string, body?: unknown): Promise<T> {
    return this.request<T>(endpoint, {
      method: "POST",
      body: body ? JSON.stringify(body) : undefined,
    });
  }

  async put<T>(endpoint: string, body?: unknown): Promise<T> {
    return this.request<T>(endpoint, {
      method: "PUT",
      body: body ? JSON.stringify(body) : undefined,
    });
  }

  async delete<T>(endpoint: string): Promise<T> {
    return this.request<T>(endpoint, { method: "DELETE" });
  }

  async getChatHistory(
    profileId: string,
    context: "diagnostic" | "tutor" | "onboarding",
    limit?: number,
    offset?: number
  ): Promise<ChatHistoryMessage[]> {
    const params = new URLSearchParams({ context });
    if (limit) params.append("limit", limit.toString());
    if (offset) params.append("offset", offset.toString());
    return this.get<ChatHistoryMessage[]>(`/api/chat/history/${profileId}?${params.toString()}`);
  }

  async updateNativeLanguage(profileId: number, nativeLanguage: string) {
    return this.put<any>(`/api/profile/${profileId}/native-language`, { native_language: nativeLanguage });
  }

  async onboardingChat(profileId: string, message: string) {
    return this.post<any>(`/api/ai/onboarding/${profileId}`, { message });
  }

  async updateProfile(nickname: string) {
    return this.put<any>("/api/auth/profile", { nickname });
  }

  async changePassword(old_password: string, new_password: string) {
    return this.put<any>("/api/auth/password", { old_password, new_password });
  }
}

export const api = new ApiClient(API_URL);
