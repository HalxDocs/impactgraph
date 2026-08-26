import { useMemo, useState } from "react";
import { HugeiconsIcon } from "@hugeicons/react";
import { FlowIcon, ZoomInAreaIcon } from "@hugeicons/core-free-icons";
import { humanRelation } from "../utils/humanize";
import type { Graph, GraphNode } from "../types/graph";

const typeColor: Record<string, string> = {
  Package: "#000000",
  Repository: "#18181b",
  Service: "#3f3f46",
  Developer: "#71717a",
  Team: "#a1a1aa",
  Organization: "#52525b",
};
const typeBg: Record<string, string> = {
  Package: "bg-black",
  Repository: "bg-zinc-900",
  Service: "bg-zinc-700",
  Developer: "bg-zinc-500",
  Team: "bg-zinc-400",
  Organization: "bg-zinc-600",
};

function layout(nodes: GraphNode[], edges: { source: string; target: string }[], width = 800, height = 460) {
  const pos: Record<string, { x: number; y: number }> = {};
  if (nodes.length === 0) return pos;
  const adj: Record<string, string[]> = {};
  nodes.forEach((n) => (adj[n.id] = []));
  edges.forEach((e) => {
    if (adj[e.source]) adj[e.source].push(e.target);
    if (adj[e.target]) adj[e.target].push(e.source);
  });
  const level: Record<string, number> = {};
  const queue: string[] = [nodes[0].id];
  level[nodes[0].id] = 0;
  const visited = new Set<string>([nodes[0].id]);
  let maxLevel = 0;
  while (queue.length) {
    const cur = queue.shift()!;
    const lv = level[cur];
    maxLevel = Math.max(maxLevel, lv);
    for (const nb of adj[cur] || []) {
      if (!visited.has(nb)) {
        visited.add(nb);
        level[nb] = lv + 1;
        queue.push(nb);
      }
    }
  }
  const groups: Record<number, string[]> = {};
  nodes.forEach((n) => {
    const lv = level[n.id] ?? maxLevel + 1;
    if (!groups[lv]) groups[lv] = [];
    groups[lv].push(n.id);
  });
  const levels = Object.keys(groups).map(Number).sort((a, b) => a - b);
  const colW = width / (levels.length + 1);
  levels.forEach((lv, idx) => {
    const ids = groups[lv];
    const rowH = height / (ids.length + 1);
    ids.forEach((id, j) => {
      pos[id] = { x: colW * (idx + 1), y: rowH * (j + 1) + 18 };
    });
  });
  if (levels.length === 1 && nodes.length > 5) {
    const cx = width / 2, cy = height / 2, r = Math.min(width, height) * 0.34;
    nodes.forEach((n, i) => {
      const ang = (2 * Math.PI * i) / nodes.length;
      pos[n.id] = { x: cx + r * Math.cos(ang), y: cy + r * Math.sin(ang) };
    });
  }
  return pos;
}

export default function GraphViewer({
  graph,
  highlightNodes,
  highlightEdges,
  onNodeClick,
  selectedId,
}: {
  graph: Graph;
  highlightNodes?: Set<string>;
  highlightEdges?: Set<string>;
  onNodeClick?: (n: GraphNode) => void;
  selectedId?: string;
}) {
  const [hover, setHover] = useState<string | null>(null);
  const pos = useMemo(() => layout(graph.nodes, graph.edges), [graph]);

  if (graph.nodes.length === 0) {
    return (
      <div className="flex h-[520px] flex-col items-center justify-center rounded-2xl border border-dashed border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900 p-8 text-center">
        <div className="flex h-12 w-12 items-center justify-center rounded-2xl bg-zinc-50 dark:bg-zinc-800 border dark:border-zinc-700"><HugeiconsIcon icon={FlowIcon} size={20} className="text-zinc-400" /></div>
        <p className="mt-3 text-sm font-medium text-black dark:text-white">No connection picture yet</p>
        <p className="text-xs text-zinc-500 dark:text-zinc-400 max-w-sm">Search above — try <span className="font-medium text-black dark:text-white">PostgreSQL</span> or <span className="font-medium text-black dark:text-white">Checkout Service</span> — then choose an item. We’ll draw how everything connects, with the affected chain highlighted.</p>
        <p className="mt-2 text-[11px] text-zinc-400">Tip: the big “See what would be affected” button shows the full impact.</p>
      </div>
    );
  }

  return (
    <div className="overflow-hidden rounded-2xl border border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900 shadow-sm">
      <div className="flex items-center justify-between border-b dark:border-zinc-800 bg-zinc-50/70 dark:bg-zinc-900/50 px-4 py-3">
        <div className="flex items-center gap-2 text-xs font-medium text-zinc-600 dark:text-zinc-400">
          <HugeiconsIcon icon={FlowIcon} size={14} />
          Picture: {graph.nodes.length} items • {graph.edges.length} connections
          {highlightNodes && highlightNodes.size > 0 && <span className="rounded-full bg-black dark:bg-white px-2 py-0.5 text-xs font-semibold text-white dark:text-black">{highlightNodes.size} highlighted — what would be affected</span>}
        </div>
        <div className="hidden sm:flex items-center gap-1.5 text-[11px]">
          {Object.entries(typeColor).map(([k, c]) => (
            <span key={k} className="inline-flex items-center gap-1.5 rounded-full border dark:border-zinc-700 bg-white dark:bg-zinc-800 px-2.5 py-1 font-medium text-zinc-600 dark:text-zinc-300 shadow-sm">
              <span className={`h-2 w-2 rounded-full ${typeBg[k]}`} style={{ background: c }} />{k}
            </span>
          ))}
        </div>
      </div>
      <p className="px-4 py-2 text-xs text-zinc-500 dark:text-zinc-400 border-b dark:border-zinc-800 bg-zinc-50 dark:bg-zinc-900">Highlighted line = affected chain. Left is the starting building block, right is what relies on it — follow the arrow to see cause → effect without reading.</p>
      <div className="relative">
        <svg viewBox="0 0 800 480" className="w-full h-[480px] bg-[radial-gradient(ellipse_at_top,_var(--tw-gradient-stops))] from-zinc-50 dark:from-zinc-900 via-white dark:via-zinc-900 to-white dark:to-zinc-900">
          <defs>
            <marker id="arrow" viewBox="0 0 10 10" refX={8} refY={5} markerWidth={6} markerHeight={6} orient="auto-start-reverse">
              <path d="M 0 0 L 10 5 L 0 10 z" fill="#a1a1aa" />
            </marker>
            <marker id="arrow-hl" viewBox="0 0 10 10" refX={8} refY={5} markerWidth={6} markerHeight={6} orient="auto-start-reverse">
              <path d="M 0 0 L 10 5 L 0 10 z" fill="#000000" />
            </marker>
          </defs>
          {graph.edges.map((e) => {
            const a = pos[e.source], b = pos[e.target];
            if (!a || !b) return null;
            const isHL = highlightEdges?.has(e.id) || (highlightNodes?.has(e.source) && highlightNodes?.has(e.target));
            const hl = isHL || (highlightEdges && Array.from(highlightEdges).some((id) => id.includes(e.source) && id.includes(e.target)));
            const label = humanRelation[e.type] || e.type.toLowerCase();
            return (
              <g key={e.id + e.source + e.target}>
                <line x1={a.x} y1={a.y} x2={b.x} y2={b.y} stroke={hl ? "#000000" : "#e4e4e7"} strokeWidth={hl ? 3.2 : 1.6} markerEnd={hl ? "url(#arrow-hl)" : "url(#arrow)"} opacity={hl ? 1 : 0.35} className={hl ? "" : "dark:opacity-20"} />
                <text x={(a.x + b.x) / 2} y={(a.y + b.y) / 2 - 8} fontSize={8.5} fill={hl ? "#000000" : "#a1a1aa"} fontWeight={hl ? 700 : 400} textAnchor="middle" className="select-none" style={{ fontFamily: "JetBrains Mono, monospace" }}>{label}</text>
              </g>
            );
          })}
          {graph.nodes.map((n) => {
            const p = pos[n.id];
            if (!p) return null;
            const isHL = highlightNodes?.has(n.id);
            const isSelected = selectedId === n.id;
            const color = typeColor[n.type] || "#71717a";
            const r = isSelected ? 23 : isHL ? 21 : 17;
            const dimmed = highlightNodes && highlightNodes.size > 0 && !isHL && !isSelected;
            return (
              <g key={n.id} transform={`translate(${p.x},${p.y})`} onMouseEnter={() => setHover(n.id)} onMouseLeave={() => setHover(null)} onClick={() => onNodeClick?.(n)} className="cursor-pointer" opacity={dimmed ? 0.35 : 1}>
                <circle r={r + 4} fill={color} opacity={0.14} className={isHL || isSelected ? "opacity-100" : "opacity-0"} />
                <circle r={r} fill={color} stroke={isHL ? "#000000" : "white"} strokeWidth={isHL || isSelected ? 3 : 2.5} />
                <text textAnchor="middle" dy={4.5} fontSize={9} fontWeight={800} fill="white" className="pointer-events-none select-none">{n.name.slice(0, 2).toUpperCase()}</text>
                <text textAnchor="middle" dy={r + 15} fontSize={10.5} fontWeight={650} fill={isHL ? "#000000" : "#18181b"} className="select-none dark:fill-zinc-100" style={{ fontFamily: "Instrument Sans, sans-serif" }}>{n.name.length > 18 ? n.name.slice(0, 18) + "…" : n.name}</text>
                <text textAnchor="middle" dy={r + 27} fontSize={8.5} fontWeight={500} fill="#a1a1aa" className="select-none uppercase tracking-widest">{n.type}</text>
                {hover === n.id && (
                  <g>
                    <rect x={-72} y={-54} width={144} height={26} rx={10} fill="#18181b" opacity={0.96} />
                    <text textAnchor="middle" dy={-37} fontSize={8} fill="white" fontWeight={500} style={{ fontFamily: "JetBrains Mono, monospace" }}>{n.id}</text>
                  </g>
                )}
              </g>
            );
          })}
        </svg>
        <div className="pointer-events-none absolute left-3 top-3 flex items-center gap-1.5 rounded-full bg-white/90 dark:bg-zinc-800/90 px-3 py-1.5 text-xs font-medium text-zinc-600 dark:text-zinc-300 shadow-sm border dark:border-zinc-700 backdrop-blur">
          <HugeiconsIcon icon={ZoomInAreaIcon} size={12} /> Click any item to see what it is first
        </div>
      </div>
    </div>
  );
}
