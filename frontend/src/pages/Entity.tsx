import { useEffect, useState } from "react";
import { useParams, Link } from "react-router-dom";
import { HugeiconsIcon } from "@hugeicons/react";
import { ArrowLeft01Icon, FlashIcon, CubeIcon, InformationCircleIcon } from "@hugeicons/core-free-icons";
import { plainTypeLabel } from "../utils/humanize";
import { api } from "../api/client";
import GraphViewer from "../components/GraphViewer";
import { ImpactPanel } from "../components/ImpactPanel";
import type { Graph, ImpactResponse } from "../types/graph";

function plainSentence(entity: any): string {
  const name = entity.name;
  const type = entity.type;
  const stats = entity.stats;
  if (type === "Package") return `${name} is a building block — a database or library — relied on by ${stats?.usedByRepos ?? 0} code areas and ${stats?.usedByServices ?? 0} services. If ${name} breaks, those could stop working correctly.`;
  if (type === "Service") return `${name} is a running service that users rely on — it ${stats?.usedByServices ? `is connected to other services` : `handles one part of the product`}.`;
  if (type === "Repository") return `${name} is a code area — where the code for one part lives — used to run services.`;
  if (type === "Team") return `${name} is a team that looks after code and services.`;
  if (type === "Developer") return `${name} is a person who helps build and look after code.`;
  if (type === "Organization") return `${name} is a group that brings teams together.`;
  return `${name} is a ${plainTypeLabel(type)} in the system.`;
}

export default function Entity() {
  const { type, id } = useParams<{ type: string; id: string }>();
  const [entity, setEntity] = useState<any>(null);
  const [graph, setGraph] = useState<Graph>({ nodes: [], edges: [] });
  const [impact, setImpact] = useState<ImpactResponse | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!type || !id) return;
    setError(null);
    setImpact(null);
    api.entity(type, id).then(setEntity).catch((e) => setError("We couldn't reach the database right now. Try again in a moment."));
    api.entityGraph(type, id).then(setGraph).catch((e) => setError("We couldn't reach the database right now. Try again in a moment."));
  }, [type, id]);

  const analyze = async () => {
    if (!id) return;
    try {
      const res = await api.impact(id);
      setImpact(res);
      setGraph({ nodes: res.nodes, edges: res.edges });
    } catch (e: any) {
      setError("We couldn't reach the database right now. Try again in a moment.");
    }
  };

  if (error) return <div className="mx-auto max-w-6xl p-8"><div className="rounded-2xl border border-red-200 dark:border-red-900/50 bg-red-50 dark:bg-red-950/30 p-6 text-sm text-red-700 dark:text-red-300">{error} <Link to="/" className="underline font-medium">Back to Dashboard</Link></div></div>;
  if (!entity) return <div className="p-12 flex items-center justify-center gap-2 text-sm text-zinc-500 dark:text-zinc-400">Looking at what's connected to this item…</div>;

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-center gap-2 text-xs text-zinc-500 dark:text-zinc-400">
        <Link to="/" className="hover:text-black dark:hover:text-white">Dashboard</Link>
        <span>›</span>
        <Link to="/explore" className="hover:text-black dark:hover:text-white">Explore</Link>
        <span>›</span>
        <span className="rounded-full bg-zinc-100 dark:bg-zinc-800 px-2.5 py-1 text-black dark:text-white font-medium">Viewing: {entity.name} {impact ? "→ What would be affected" : "→ Details"}</span>
      </div>
      <Link to="/explore" className="inline-flex items-center gap-1.5 text-sm font-medium text-zinc-600 dark:text-zinc-400 hover:text-black dark:hover:text-white transition"><HugeiconsIcon icon={ArrowLeft01Icon} size={14} /> Back to Explore — search and picture</Link>

      <div className="overflow-hidden rounded-2xl border border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900 shadow-sm">
        <div className="bg-black dark:bg-white px-6 py-7 text-white dark:text-black relative overflow-hidden">
          <div className="absolute inset-0 bg-[linear-gradient(to_right,#ffffff0a_1px,transparent_1px),linear-gradient(to_bottom,#ffffff0a_1px,transparent_1px)] dark:bg-[linear-gradient(to_right,#00000008_1px,transparent_1px),linear-gradient(to_bottom,#00000008_1px,transparent_1px)] bg-[size:32px_32px]" />
          <div className="relative">
            <div className="inline-flex items-center gap-2 rounded-full border border-white/15 dark:border-black/10 bg-white/10 dark:bg-black/5 px-3 py-1 text-xs font-medium backdrop-blur"><HugeiconsIcon icon={CubeIcon} size={12} color="currentColor" />{plainTypeLabel(entity.type)} — {entity.type}</div>
            <h1 className="mt-3 text-3xl font-bold tracking-tight" style={{ fontFamily: "Space Grotesk, sans-serif" }}>{entity.name}</h1>
            <p className="mt-2 max-w-2xl text-sm leading-relaxed text-white/80 dark:text-black/70">{plainSentence(entity)}</p>
            <p className="mt-1 text-xs text-white/60 dark:text-black/50 font-mono">{entity.id} • Below are the technical details, if you need them.</p>
          </div>
        </div>
        <div className="px-6 py-4 bg-zinc-50/50 dark:bg-zinc-900/50 border-t dark:border-zinc-800">
          <p className="text-xs font-semibold tracking-widest text-zinc-400 dark:text-zinc-500 uppercase mb-2">Technical details — secondary</p>
          <div className="flex flex-wrap gap-2">
            {Object.entries(entity.properties || {}).map(([k, v]) => (
              <span key={k} className="inline-flex items-center gap-1.5 rounded-full border dark:border-zinc-700 bg-white dark:bg-zinc-800 px-3 py-1.5 text-xs shadow-sm dark:text-zinc-200"><span className="text-zinc-500 dark:text-zinc-400">{k}</span><span className="font-medium text-black dark:text-white">{String(v)}</span></span>
            ))}
          </div>
          {entity.stats && (
            <p className="mt-3 text-xs text-zinc-500 dark:text-zinc-400">Relied on by <span className="font-medium text-black dark:text-white">{entity.stats.usedByRepos} code areas</span>, <span className="font-medium text-black dark:text-white">{entity.stats.usedByServices} services</span>, <span className="font-medium text-black dark:text-white">{entity.stats.usedByTeams} teams</span>. <span className="text-[11px]">“Relies on” = depends on • “Runs as” = deployed as • “Talks to” = calls</span></p>
          )}
          <div className="mt-3 flex items-center gap-2 text-xs text-zinc-500 dark:text-zinc-400"><HugeiconsIcon icon={InformationCircleIcon} size={12} /> Picture below shows what’s directly connected — one step away.</div>
        </div>
        <div className="px-6 pb-4 bg-white dark:bg-zinc-900 flex flex-wrap gap-3">
          {entity.type === "Package" && (
            <button onClick={analyze} className="inline-flex items-center gap-2 rounded-full bg-black dark:bg-white px-6 py-3 text-sm font-semibold text-white dark:text-black shadow hover:bg-zinc-800 dark:hover:bg-zinc-100 transition">
              <HugeiconsIcon icon={FlashIcon} size={16} color="currentColor" /> See what would be affected if this changes
            </button>
          )}
          <Link to="/explore" className="inline-flex items-center gap-2 rounded-full border dark:border-zinc-700 bg-white dark:bg-zinc-800 px-5 py-3 text-sm font-medium text-black dark:text-white">Back to picture</Link>
        </div>
      </div>

      {impact && <ImpactPanel impact={impact} />}
      <div>
        <p className="mb-2 text-xs font-semibold tracking-widest text-zinc-400 dark:text-zinc-500 uppercase">Connection picture — highlighted line shows what would be affected</p>
        <GraphViewer graph={graph} highlightNodes={impact ? new Set(impact.nodes.map((n) => n.id)) : undefined} />
      </div>
    </div>
  );
}
