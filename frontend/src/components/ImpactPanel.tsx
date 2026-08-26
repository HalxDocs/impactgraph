import { useState } from "react";
import { HugeiconsIcon } from "@hugeicons/react";
import { Analytics01Icon, GitBranchIcon, FlashIcon, WorkflowCircle06Icon, InformationCircleIcon, HelpCircleIcon, ArrowRight01Icon } from "@hugeicons/core-free-icons";
import type { ImpactResponse, PathResponse } from "../types/graph";

function severity(affected: ImpactResponse["affected"]) {
  const total = affected.services + affected.repositories;
  if (total >= 12) return { label: "High impact", desc: "Many services could be affected — worth a careful check before changing." };
  if (total >= 6) return { label: "Medium impact", desc: "A few services could be affected." };
  if (total > 0) return { label: "Low impact", desc: "Only a small part of the system is affected." };
  return { label: "No connected impact found", desc: "Nothing else relies on this in the current map." };
}

export function ImpactPanel({ impact }: { impact: ImpactResponse | null }) {
  const [showHelp, setShowHelp] = useState(false);
  if (!impact) return null;
  const sev = severity(impact.affected);
  const narrated = (() => {
    const r = impact.root.name;
    if (impact.affected.services > 0) return `${r} feeds into code that powers services — if ${r} breaks, those services could stop working correctly.`;
    if (impact.affected.repositories > 0) return `${r} is relied on by code areas — if ${r} changes, that code may need updating.`;
    return `No services were found that rely on ${r} in this map.`;
  })();

  return (
    <div className="overflow-hidden rounded-2xl border border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900 shadow-sm">
      <div className="flex items-center justify-between border-b dark:border-zinc-800 bg-zinc-50/50 dark:bg-zinc-900/50 px-5 py-4">
        <div className="flex items-center gap-3">
          <div className="flex h-8 w-8 items-center justify-center rounded-xl bg-black dark:bg-white text-white dark:text-black"><HugeiconsIcon icon={FlashIcon} size={16} color="currentColor" /></div>
          <div>
            <h3 className="text-sm font-semibold text-black dark:text-white" style={{ fontFamily: "Space Grotesk, sans-serif" }}>What would be affected if {impact.root.name} changes?</h3>
            <p className="text-xs text-zinc-500 dark:text-zinc-400">Shows the chain starting from {impact.root.name} — a {impact.root.type.toLowerCase() === "package" ? "building block" : impact.root.type.toLowerCase()}.</p>
          </div>
        </div>
        <span className="hidden sm:inline-flex items-center gap-1.5 rounded-full bg-black dark:bg-white px-3 py-1 text-xs font-semibold text-white dark:text-black">{sev.label}</span>
      </div>
      <div className="px-5 pt-4">
        <p className="text-sm text-zinc-700 dark:text-zinc-300">{narrated}</p>
        <p className="mt-1 text-xs text-zinc-500 dark:text-zinc-400">{sev.desc}</p>
        <button onClick={() => setShowHelp(!showHelp)} className="mt-3 inline-flex items-center gap-1.5 rounded-full border border-zinc-200 dark:border-zinc-700 px-3 py-1 text-xs font-medium text-zinc-600 dark:text-zinc-300 hover:bg-zinc-50 dark:hover:bg-zinc-800">
          <HugeiconsIcon icon={HelpCircleIcon} size={12} /> How we find this (plain) <HugeiconsIcon icon={ArrowRight01Icon} size={10} className={`transition ${showHelp ? "rotate-90" : ""}`} />
        </button>
        {showHelp && (
          <div className="mt-2 rounded-xl border border-zinc-200 dark:border-zinc-800 bg-zinc-50 dark:bg-zinc-800/50 p-3 text-xs leading-relaxed text-zinc-600 dark:text-zinc-300">
            We start from <span className="font-medium text-black dark:text-white">{impact.root.name}</span> → find code areas that <span className="font-medium">rely on</span> it → find services that <span className="font-medium">run</span> that code → follow services that <span className="font-medium">talk to</span> each other. The black highlighted line in the picture is that chain. “Steps” counts how far it reaches, worst case — where we show <span className="font-medium">{impact.maxDepth} steps away</span>. Uses one Cypher traversal; see README “Docs — Impact Analysis” for the exact query.
          </div>
        )}
      </div>
      <div className="grid grid-cols-2 md:grid-cols-4 gap-3 p-4">
        <div className="rounded-2xl border dark:border-zinc-700 bg-white dark:bg-zinc-800 p-4">
          <div className="inline-flex h-7 w-7 items-center justify-center rounded-xl bg-black dark:bg-white text-white dark:text-black mb-3"><HugeiconsIcon icon={GitBranchIcon} size={14} color="currentColor" /></div>
          <div className="text-2xl font-bold tracking-tight text-black dark:text-white">{impact.affected.repositories}</div>
          <div className="text-xs font-semibold text-black dark:text-white">Code areas</div>
          <div className="text-[11px] text-zinc-500 dark:text-zinc-400">Places where code lives</div>
          <div className="mt-1 text-xs text-zinc-600 dark:text-zinc-300">If {impact.root.name} breaks, {impact.affected.repositories} code areas could need fixes.</div>
        </div>
        <div className="rounded-2xl border dark:border-zinc-700 bg-white dark:bg-zinc-800 p-4">
          <div className="inline-flex h-7 w-7 items-center justify-center rounded-xl bg-zinc-800 dark:bg-zinc-700 text-white mb-3"><HugeiconsIcon icon={WorkflowCircle06Icon} size={14} color="white" /></div>
          <div className="text-2xl font-bold tracking-tight text-black dark:text-white">{impact.affected.services}</div>
          <div className="text-xs font-semibold text-black dark:text-white">Services</div>
          <div className="text-[11px] text-zinc-500 dark:text-zinc-400">Things that keep running for users</div>
          <div className="mt-1 text-xs text-zinc-600 dark:text-zinc-300">If {impact.root.name} breaks, {impact.affected.services} services could stop working correctly.</div>
        </div>
        <div className="rounded-2xl border dark:border-zinc-700 bg-white dark:bg-zinc-800 p-4">
          <div className="inline-flex h-7 w-7 items-center justify-center rounded-xl bg-zinc-700 dark:bg-zinc-600 text-white mb-3"><HugeiconsIcon icon={Analytics01Icon} size={14} color="white" /></div>
          <div className="text-2xl font-bold tracking-tight text-black dark:text-white">{impact.affected.teams}</div>
          <div className="text-xs font-semibold text-black dark:text-white">Teams</div>
          <div className="text-[11px] text-zinc-500 dark:text-zinc-400">People who look after it</div>
          <div className="mt-1 text-xs text-zinc-600 dark:text-zinc-300">{impact.affected.teams} teams would need to be told.</div>
        </div>
        <div className="rounded-2xl border border-black dark:border-white bg-black dark:bg-white p-4 text-white dark:text-black">
          <div className="inline-flex h-7 w-7 items-center justify-center rounded-xl bg-white dark:bg-black text-black dark:text-white mb-3"><HugeiconsIcon icon={FlashIcon} size={14} color="currentColor" /></div>
          <div className="text-2xl font-bold tracking-tight">{impact.maxDepth}<span className="text-sm font-normal text-zinc-400 dark:text-zinc-600"> steps</span></div>
          <div className="text-xs font-semibold">Longest chain</div>
          <div className="text-[11px] text-zinc-400 dark:text-zinc-600">Steps away, worst case</div>
          <div className="mt-1 text-xs text-zinc-300 dark:text-zinc-700">{impact.maxDepth} steps away from {impact.root.name}, worst case — how far the chain could reach.</div>
        </div>
      </div>
      <div className="flex items-center gap-2 border-t dark:border-zinc-800 bg-zinc-50 dark:bg-zinc-900 px-5 py-3 text-xs text-zinc-700 dark:text-zinc-300"><HugeiconsIcon icon={InformationCircleIcon} size={14} /> Highlighted line in the picture above shows the affected chain. Left is the starting building block, right is what relies on it.</div>
    </div>
  );
}

export function PathPanel({ pathResp }: { pathResp: PathResponse | null }) {
  const [showHelp, setShowHelp] = useState(false);
  if (!pathResp) return null;
  if (!pathResp.exists) return <div className="rounded-2xl border border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900 p-6 text-sm text-zinc-600 dark:text-zinc-400 flex items-center gap-2"><HugeiconsIcon icon={InformationCircleIcon} size={16} />No connection found between the chosen items — try a different pair. Docs: we check <span className="font-medium">relies on → runs as → talks to</span> chain.</div>;
  const chain = pathResp.path.map(n => n.name).join(" → shows a link to → ");
  const narrated = pathResp.path.length <= 2 ? `${pathResp.path[0].name} is directly connected to ${pathResp.path[pathResp.path.length-1].name}.` : `${pathResp.path[0].name} connects to ${pathResp.path.slice(1).map(p=>p.name).join(", which connects to ")} — that's the chain.`;
  return (
    <div className="overflow-hidden rounded-2xl border border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900 shadow-sm">
      <div className="border-b dark:border-zinc-800 bg-zinc-50/50 dark:bg-zinc-900/50 px-5 py-4 flex items-center gap-3">
        <div className="flex h-8 w-8 items-center justify-center rounded-xl bg-black dark:bg-white text-white dark:text-black"><HugeiconsIcon icon={GitBranchIcon} size={16} color="currentColor" /></div>
        <div>
          <h3 className="text-sm font-semibold text-black dark:text-white">Why is this affected? — the chain in plain language</h3>
          <p className="text-xs text-zinc-500 dark:text-zinc-400">{pathResp.hops} steps between the two items • follow the line to see cause → effect</p>
        </div>
      </div>
      <div className="p-5">
        <p className="text-sm text-zinc-700 dark:text-zinc-300 mb-3">{narrated}</p>
        <div className="flex flex-wrap items-center gap-2">
          {pathResp.path.map((n, i) => (
            <span key={n.id} className="flex items-center gap-2">
              <span className="inline-flex items-center gap-1.5 rounded-full bg-black dark:bg-white px-3 py-1.5 text-xs font-semibold text-white dark:text-black shadow-sm">{n.name}<span className="text-white/60 dark:text-black/60 text-[10px]">{n.type === "Package" ? "building block" : n.type === "Repository" ? "code area" : n.type === "Service" ? "service" : n.type.toLowerCase()}</span></span>
              {i < pathResp.path.length - 1 && <span className="text-zinc-300 dark:text-zinc-600">→</span>}
            </span>
          ))}
        </div>
        <div className="mt-3 inline-flex items-center gap-2 rounded-full border dark:border-zinc-700 bg-zinc-50 dark:bg-zinc-800 px-3 py-1.5 text-xs font-medium text-zinc-600 dark:text-zinc-300"><HugeiconsIcon icon={WorkflowCircle06Icon} size={12} />Chain uses: {pathResp.edges.map((e) => ({"DEPENDS_ON":"relies on","DEPLOYED_AS":"runs as","CALLS":"talks to","OWNS":"owned by","BELONGS_TO":"belongs to"} as Record<string,string>)[e.type] || e.type.toLowerCase()).join(" → ")}</div>
        <button onClick={() => setShowHelp(!showHelp)} className="mt-3 inline-flex items-center gap-1.5 rounded-full border border-zinc-200 dark:border-zinc-700 px-3 py-1 text-xs font-medium text-zinc-600 dark:text-zinc-300">
          <HugeiconsIcon icon={HelpCircleIcon} size={12} /> How this chain is found <HugeiconsIcon icon={ArrowRight01Icon} size={10} className={showHelp ? "rotate-90" : ""} />
        </button>
        {showHelp && <div className="mt-2 rounded-xl border border-zinc-200 dark:border-zinc-800 bg-zinc-50 dark:bg-zinc-800/50 p-3 text-xs leading-relaxed text-zinc-600 dark:text-zinc-300">We look for the shortest chain between the two items you picked, checking <span className="font-medium">relies on → runs as → talks to</span> connections (never the raw codes). In the database this is <span className="font-mono text-[11px] bg-white dark:bg-zinc-900 border rounded px-1">shortestPath((a)-[:...*..6]-(b))</span>. See README “Docs — Dependency Path” for the exact Cypher and why the picture is the same chain highlighted in black.</div>}
      </div>
    </div>
  );
}
