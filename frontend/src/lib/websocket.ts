const WS_URL = process.env.NEXT_PUBLIC_WS_URL || "ws://localhost:8080";

export interface WebSocketMessage {
  type: "diagnostic" | "tutor" | "generate_plan" | "onboarding";
  message?: string;
  topic?: string;
  profile_id: number;
}

export interface WebSocketResponse {
  type: string;
  content?: string;
  chunk?: string;
  level?: string;
  weaknesses?: string[];
  is_assessment?: boolean;
  phase?: string;
  goals?: string[];
  is_completed?: boolean;
  plan?: string;
  error?: string;
  done?: boolean;
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
