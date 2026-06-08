export type TTSStrategy = "web-speech" | "server";

export function getTTSStrategy(): TTSStrategy {
  return "web-speech";
}
