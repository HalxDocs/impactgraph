import { useEffect, useState } from "react";
import { HugeiconsIcon } from "@hugeicons/react";
import { Search01Icon, FlashIcon, GitBranchIcon, ArrowRight01Icon, InformationCircleIcon, WorkflowCircle06Icon, SidebarLeftIcon, Cancel01Icon } from "@hugeicons/core-free-icons";
import SearchBar from "../components/SearchBar";
import GraphViewer from "../components/GraphViewer";
import { ImpactPanel, PathPanel } from "../components/ImpactPanel";
import { api } from "../api/client";
import type { Graph, ImpactResponse, PathResponse, SearchResult } from "../types/graph";

export default function Explore() {
  const [selected, setSelected] = useState<SearchResult | null>(null);
  const [graph, setGraph] = useState<Graph>({ nodes: [], edges: [] });
  const [impact, setImpact] = useState<ImpactResponse | null>(null);
  const [pathResp, setPathResp] = useState<PathResponse | null>(null);
  const [graphError, setGraphError] = useState<string | null>(null);
  const [from, setFrom] = useState("");
  const [to, setTo] = useState("");
  const [allNodes, setAllNodes] = useState<{ id: string; name: string; type: string }[]>([]);
  const [loadingGraph, setLoadingGraph] = useState(false);
  const [sidebarOpen, setSidebarOpen] = useState(true);

  useEffect(() => {
    api.nodes().then((ns) => setAllNodes(ns.map((n) => ({ id: n.id, name: n.name, type: n.type })))).catch(() => {});
  }, []);

  const onSelect = async (r: SearchResult) => {
    setSelected(r);
    setImpact(null);
    setPathResp(null);
    setGraphError(null);
    setFrom(r.id);
    setLoadingGraph(true);
    try {
      const g = await api.entityGraph(r.type, r.id);
      setGraph(g);
    } catch (e: any) {
      setGraphError("We couldn't reach the database right now. Try again in a moment.");
      setGraph({ nodes: [], edges: [] });
    } finally { setLoadingGraph(false); }
  };

  const analyzeImpact = async () => {
    if (!selected) return;
    setLoadingGraph(true);
    try {
      const res = await api.impact(selected.id);
      setImpact(res);
      setGraph({ nodes: res.nodes, edges: res.edges });
    } catch (e: any) {
      setGraphError("We couldn't reach the database right now. Try again in a moment.");
    } finally { setLoadingGraph(false); }
  };

  const findConnection = async () => {
    if (!from || !to) return;
    setLoadingGraph(true);
    try {
      const res = await api.path(from, to);
      setPathResp(res);
      if (res.exists) {
        const nodes = res.path;
        const edges = res.edges.map((e, i) => ({ ...e, id: `${e.source}-${e.target}-${i}` }));
        setGraph({ nodes: nodes as any, edges: edges as any });
      }
    } catch (e: any) {
      setGraphError("We couldn't reach the database right now. Try again in a moment.");
    } finally { setLoadingGraph(false); }
  };

  const highlightNodes = impact ? new Set(impact.nodes.map((n) => n.id)) : pathResp?.exists ? new Set(pathResp.path.map((n) => n.id)) : undefined;
  const highlightEdges = impact ? new Set(impact.edges.map((e) => e.id)) : pathResp?.exists ? new Set(pathResp.edges.map((e) => e.id)) : undefined;

  return (
    <div className="space-y-4">
      {/* Breadcrumbs + toggle */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2 text-xs text-zinc-500 dark:text-zinc-400">
          <span className="rounded-full bg-zinc-100 dark:bg-zinc-800 px-2.5 py-1 font-medium text-black dark:text-white">Explore</span>
          <span className="hidden sm:inline">Pick something to see what relies on it</span>
          {selected && <span className="hidden md:flex items-center gap-1">→ {selected.name} {impact && "→ What would be affected"}</span>}
        </div>
        <button onClick={() => setSidebarOpen(!sidebarOpen)} className="inline-flex items-center gap-1.5 rounded-full border border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900 px-3 py-1.5 text-xs font-medium text-zinc-700 dark:text-zinc-300 hover:bg-zinc-50 dark:hover:bg-zinc-800">
          <HugeiconsIcon icon={sidebarOpen ? Cancel01Icon : SidebarLeftIcon} size={14} />{sidebarOpen ? "Hide controls" : "Show controls"}
        </button>
      </div>

      <div className="flex flex-col lg:flex-row gap-6 items-start">
        {/* Sidebar */}
        {sidebarOpen && (
          <aside className="w-full lg:w-[380px] shrink-0 space-y-4 lg:sticky lg:top-[80px] lg:max-h-[calc(100vh-96px)] lg:overflow-y-auto lg:pr-2">
            <div className="rounded-2xl border border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900 p-5 shadow-sm">
              <div className="flex items-center gap-2 text-xs font-semibold tracking-widest text-zinc-400 dark:text-zinc-500 uppercase"><HugeiconsIcon icon={Search01Icon} size={14} /> Find something</div>
              <h2 className="mt-1 text-sm font-semibold text-black dark:text-white" style={{ fontFamily: "Space Grotesk, sans-serif" }}>Search building blocks and services</h2>
              <p className="text-xs text-zinc-500 dark:text-zinc-400">Shows your building blocks, code areas and services as a picture. Choose one.</p>
              <div className="mt-4"><SearchBar onSelect={onSelect} /></div>
              {!selected && !loadingGraph && graph.nodes.length === 0 && (
                <div className="mt-4 rounded-xl border border-zinc-200 dark:border-zinc-800 bg-zinc-50 dark:bg-zinc-800/30 p-4">
                  <p className="text-sm font-medium text-black dark:text-white">What to do first:</p>
                  <ol className="mt-2 list-decimal list-inside text-sm text-zinc-600 dark:text-zinc-300 space-y-1">
                    <li>Type <span className="font-medium text-black dark:text-white">PostgreSQL</span> or <span className="font-medium text-black dark:text-white">Checkout Service</span></li>
                    <li>Pick the result — picture updates on the right</li>
                    <li>Hit <span className="font-medium text-black dark:text-white">See what would be affected</span></li>
                  </ol>
                </div>
              )}
            </div>

            {selected && (
              <div className="rounded-2xl border-2 border-black dark:border-white bg-white dark:bg-zinc-900 p-4 shadow-sm">
                <div className="flex items-center gap-3">
                  <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-black dark:bg-white text-white dark:text-black font-bold text-sm">{selected.name.slice(0,2).toUpperCase()}</div>
                  <div className="min-w-0 flex-1">
                    <div className="text-sm font-semibold text-black dark:text-white truncate">{selected.name}</div>
                    <div className="text-xs text-zinc-500 dark:text-zinc-400">{selected.type === "Package" ? "building block" : selected.type === "Repository" ? "code area" : selected.type.toLowerCase()} • {selected.id}</div>
                  </div>
                </div>
                <p className="mt-2 text-xs text-zinc-600 dark:text-zinc-300">Starting point — see everything that would be affected if it changes.</p>
                {!impact ? (
                  <button onClick={analyzeImpact} className="mt-3 w-full inline-flex items-center justify-center gap-2 rounded-full bg-black dark:bg-white px-5 py-3 text-sm font-semibold text-white dark:text-black shadow hover:bg-zinc-800 dark:hover:bg-zinc-100 transition">
                    <HugeiconsIcon icon={FlashIcon} size={16} color="currentColor" /> See what would be affected
                  </button>
                ) : (
                  <div className="mt-3 rounded-xl bg-zinc-900 dark:bg-white text-white dark:text-black px-3 py-2 text-xs font-medium text-center">Showing impact — highlighted chain on the right</div>
                )}
              </div>
            )}

            <div id="find-connection" className="rounded-2xl border border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900 p-5 shadow-sm">
              <div className="flex items-center gap-2">
                <div className="flex h-7 w-7 items-center justify-center rounded-xl bg-black dark:bg-white text-white dark:text-black"><HugeiconsIcon icon={GitBranchIcon} size={14} color="currentColor" /></div>
                <h3 className="text-sm font-semibold text-black dark:text-white">How two things connect</h3>
              </div>
              <p className="mt-1 text-xs text-zinc-500 dark:text-zinc-400">“Talks to” means one service asks another for help. Pick two items and find the shortest chain.</p>
              <div className="mt-4 space-y-3">
                <div className="relative">
                  <HugeiconsIcon icon={WorkflowCircle06Icon} size={14} className="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-zinc-400" />
                  <select value={from} onChange={(e) => setFrom(e.target.value)} className="w-full rounded-full border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-800 py-2.5 pl-9 pr-4 text-sm dark:text-zinc-100 outline-none focus:border-black dark:focus:border-white">
                    <option value="">Starting item…</option>
                    {allNodes.map((n) => (<option key={n.id} value={n.id}>{n.name}</option>))}
                  </select>
                </div>
                <div className="relative">
                  <HugeiconsIcon icon={WorkflowCircle06Icon} size={14} className="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-zinc-400" />
                  <select value={to} onChange={(e) => setTo(e.target.value)} className="w-full rounded-full border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-800 py-2.5 pl-9 pr-4 text-sm dark:text-zinc-100 outline-none focus:border-black dark:focus:border-white">
                    <option value="">Connected item…</option>
                    {allNodes.map((n) => (<option key={n.id} value={n.id}>{n.name}</option>))}
                  </select>
                </div>
                <button onClick={findConnection} className="w-full inline-flex items-center justify-center gap-1.5 rounded-full bg-black dark:bg-white px-5 py-2.5 text-sm font-semibold text-white dark:text-black hover:bg-zinc-800 dark:hover:bg-zinc-100 transition">
                  Find connection <HugeiconsIcon icon={ArrowRight01Icon} size={14} color="currentColor" />
                </button>
                <p className="text-xs text-zinc-500 dark:text-zinc-400">Example: <span className="font-medium text-black dark:text-white">PostgreSQL</span> → <span className="font-medium text-black dark:text-white">Notification Service</span></p>
              </div>
            </div>
          </aside>
        )}

        {/* Main */}
        <div className="flex-1 min-w-0 space-y-4">
          {loadingGraph && <div className="rounded-2xl border border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900 p-4 text-sm text-zinc-600 dark:text-zinc-300">Looking at what's connected to {selected?.name || "your choice"}… drawing the picture.</div>}
          {graphError && <div className="rounded-2xl border border-red-200 dark:border-red-900/50 bg-red-50 dark:bg-red-950/30 px-4 py-3 text-sm text-red-700 dark:text-red-300 flex items-center gap-2"><HugeiconsIcon icon={InformationCircleIcon} size={16} />We couldn't reach the database right now. Try again in a moment.</div>}
          <div id="impact-anchor" />
          {impact && <ImpactPanel impact={impact} />}
          {pathResp && <PathPanel pathResp={pathResp} />}
          <GraphViewer graph={graph} highlightNodes={highlightNodes} highlightEdges={highlightEdges} onNodeClick={(n) => onSelect({ id: n.id, type: n.type, name: n.name, label: n.label, properties: n.properties })} selectedId={selected?.id} />
          {!sidebarOpen && graph.nodes.length > 0 && <p className="text-center text-xs text-zinc-400">Controls hidden — <button onClick={() => setSidebarOpen(true)} className="underline font-medium text-black dark:text-white">show controls</button> to search again.</p>}
        </div>
      </div>
    </div>
  );
}
