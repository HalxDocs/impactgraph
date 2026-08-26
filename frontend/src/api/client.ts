const BASE = import.meta.env.VITE_API_URL || "";

async function fetchJSON<T>(url: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${BASE}${url}`, {
    headers: { "Content-Type": "application/json" },
    ...init,
  });
  if (!res.ok) {
    const body = await res.json().catch(() => ({}));
    const msg = body.error || `HTTP ${res.status}`;
    const err: any = new Error(msg);
    err.status = res.status;
    throw err;
  }
  return res.json();
}

export const api = {
  health: () => fetchJSON<{ status: string; database: string }>("/api/health"),
  search: (q: string) => fetchJSON<import("../types/graph").SearchResult[]>(`/api/search?q=${encodeURIComponent(q)}`),
  stats: () => fetchJSON<import("../types/graph").StatsResponse>("/api/stats"),
  nodes: () => fetchJSON<import("../types/graph").GraphNode[]>("/api/nodes"),
  entity: (type: string, id: string) => fetchJSON<any>(`/api/entities/${type}/${id}`),
  entityGraph: (type: string, id: string) => fetchJSON<import("../types/graph").Graph>(`/api/entities/${type}/${id}/graph`),
  impact: (id: string) => fetchJSON<import("../types/graph").ImpactResponse>(`/api/impact/package/${id}`),
  path: (from: string, to: string) => fetchJSON<import("../types/graph").PathResponse>(`/api/path?from=${encodeURIComponent(from)}&to=${encodeURIComponent(to)}`),
};
