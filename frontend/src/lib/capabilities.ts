const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export interface Capabilities {
  server_stt: boolean;
}

let cached: Capabilities | null = null;

export async function getCapabilities(): Promise<Capabilities> {
  if (cached) return cached;

  const response = await fetch(`${API_URL}/api/capabilities`);
  const body = await response.json();
  cached = (body.data ?? body) as Capabilities;
  return cached;
}

export function resetCapabilitiesCache() {
  cached = null;
}
