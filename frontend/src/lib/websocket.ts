const WS_URL = process.env.NEXT_PUBLIC_WS_URL || "ws://localhost:8080";

export interface VoiceCorrection {
  original: string;
  suggested: string;
  reason: string;
}

export interface LessonScore {
  fluency: number;
  grammar: number;
  vocabulary: number;
  pronunciation: number;
  notes: string;
}

export interface WebSocketMessage {
  type:
    | "diagnostic"
    | "tutor"
    | "generate_plan"
    | "onboarding"
    | "voice_tutor_turn"
    | "voice_tutor_audio"
    | "voice_lesson_start"
    | "voice_lesson_turn"
    | "voice_lesson_abandon";
  message?: string;
  topic?: string;
  profile_id: number;
  locale?: string;
  lesson_id?: number;
  module_id?: number;
  audio_b64?: string;
  audio_mime?: string;
  seq?: number;
  final?: boolean;
}

export interface WebSocketResponse {
  type: string;
  content?: string;
  chunk?: string;
  user_text?: string;
  level?: string;
  weaknesses?: string[];
  is_assessment?: boolean;
  phase?: string;
  goals?: string[];
  is_completed?: boolean;
  plan?: string;
  error?: string;
  done?: boolean;
  corrections?: VoiceCorrection[];
  lesson_id?: number;
  stage?: string;
  stage_index?: number;
  stage_total?: number;
  prompt?: string;
  target_phrases?: string[];
  hit_phrases?: string[];
  score?: LessonScore;
}

export type MessageHandler = (response: WebSocketResponse) => void;

export class WebSocketClient {
  private ws: WebSocket | null = null;
  private url: string;
  private token: string | null = null;
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;
  private reconnectDelay = 1000;
  private messageHandlers: Set<MessageHandler> = new Set();

  constructor(url: string = WS_URL) {
    this.url = url;
  }

  connect(): Promise<void> {
    return new Promise((resolve, reject) => {
      if (typeof window === "undefined") {
        reject(new Error("WebSocket only available in browser"));
        return;
      }

      this.token = localStorage.getItem("token");
      if (!this.token) {
        reject(new Error("No authentication token found"));
        return;
      }

      // Token via query (browser WS cannot set Authorization).
      // Backend also accepts Sec-WebSocket-Protocol: bearer,<jwt> when negotiated.
      const wsUrl = `${this.url}/api/ws?token=${encodeURIComponent(this.token)}`;
      this.ws = new WebSocket(wsUrl);

      this.ws.onopen = () => {
        console.log("WebSocket connected");
        this.reconnectAttempts = 0;
        resolve();
      };

      this.ws.onmessage = (event) => {
        try {
          const response: WebSocketResponse = JSON.parse(event.data);
          this.messageHandlers.forEach((handler) => handler(response));
        } catch (error) {
          console.error("Failed to parse WebSocket message:", error);
        }
      };

      this.ws.onerror = (error) => {
        console.error("WebSocket error:", error);
        reject(error);
      };

      this.ws.onclose = () => {
        console.log("WebSocket closed");
        this.attemptReconnect();
      };
    });
  }

  private attemptReconnect() {
    if (this.reconnectAttempts >= this.maxReconnectAttempts) {
      console.error("Max reconnection attempts reached");
      return;
    }

    this.reconnectAttempts++;
    const delay = this.reconnectDelay * this.reconnectAttempts;

    console.log(`Reconnecting in ${delay}ms (attempt ${this.reconnectAttempts})`);
    
    setTimeout(() => {
      this.connect().catch((err) => {
        console.error("Reconnection failed:", err);
      });
    }, delay);
  }

  send(message: WebSocketMessage): void {
    if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
      throw new Error("WebSocket is not connected");
    }

    this.ws.send(JSON.stringify(message));
  }

  onMessage(handler: MessageHandler): () => void {
    this.messageHandlers.add(handler);
    return () => this.messageHandlers.delete(handler);
  }

  disconnect(): void {
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
    this.messageHandlers.clear();
  }

  isConnected(): boolean {
    return this.ws !== null && this.ws.readyState === WebSocket.OPEN;
  }
}
