import { useEffect, useRef, useState, useCallback } from "react";
import { WebSocketClient, WebSocketMessage, WebSocketResponse } from "@/lib/websocket";

export function useWebSocket() {
  const [connected, setConnected] = useState(false);
  const [currentMessage, setCurrentMessage] = useState("");
  const wsRef = useRef<WebSocketClient | null>(null);

  useEffect(() => {
    const ws = new WebSocketClient();
    wsRef.current = ws;

    ws.connect()
      .then(() => setConnected(true))
      .catch((err) => console.error("WebSocket connection failed:", err));

    return () => {
      ws.disconnect();
    };
  }, []);

  const sendMessage = useCallback((message: WebSocketMessage) => {
    if (!wsRef.current?.isConnected()) {
      throw new Error("WebSocket not connected");
    }
    wsRef.current.send(message);
  }, []);

  const onMessage = useCallback((handler: (response: WebSocketResponse) => void) => {
    if (!wsRef.current) return () => {};
    return wsRef.current.onMessage(handler);
  }, []);

  return {
    connected,
    sendMessage,
    onMessage,
    currentMessage,
  };
}
