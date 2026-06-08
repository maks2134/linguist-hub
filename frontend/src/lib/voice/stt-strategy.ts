import { getCapabilities } from "@/lib/capabilities";

export type STTStrategy = "web-speech" | "server";

export async function getSTTStrategy(): Promise<STTStrategy> {
  const caps = await getCapabilities();
  return caps.server_stt ? "server" : "web-speech";
}
