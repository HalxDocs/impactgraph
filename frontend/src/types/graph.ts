export type NodeType = "Package" | "Repository" | "Service" | "Developer" | "Team" | "Organization";
export interface GraphNode {
  id: string;
  type: NodeType | string;
  label: string;
  name: string;
  properties: Record<string, any>;
}
export interface GraphEdge {
  id: string;
  source: string;
  target: string;
  type: string;
  properties?: Record<string, any>;
}
export interface Graph {
  nodes: GraphNode[];
  edges: GraphEdge[];
}
export interface SearchResult {
  id: string;
  type: string;
  name: string;
  label: string;
  properties: Record<string, any>;
}
export interface ImpactResponse {
  root: { id: string; name: string; type: string };
  affected: { repositories: number; services: number; teams: number; packages: number };
  maxDepth: number;
  nodes: GraphNode[];
  edges: GraphEdge[];
}
export interface PathResponse {
  exists: boolean;
  hops: number;
  path: GraphNode[];
  edges: GraphEdge[];
  reason?: string;
}
export interface StatsResponse {
  repositories: number;
  packages: number;
  services: number;
  developers: number;
  teams: number;
  organizations: number;
  mostConnected: { id: string; name: string; type: string; count: number }[];
  recent: { source: string; target: string; type: string; detail: string }[];
}
