import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { HugeiconsIcon } from "@hugeicons/react";
import { PackageIcon, GitBranchIcon, WorkflowCircle06Icon, UserGroupIcon, Building02Icon, Analytics01Icon, ArrowRight01Icon } from "@hugeicons/core-free-icons";
import { api } from "../api/client";
import type { StatsResponse } from "../types/graph";

export default function Dashboard() {
  const [stats, setStats] = useState<StatsResponse | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [health, setHealth] = useState<string>("");

  useEffect(() => {
    api.stats().then(setStats).catch((e) => setError("We couldn't reach the database right now. Try again in a moment."));
    api.health().then((h) => setHealth(h.database)).catch(() => setHealth("unknown"));
  }, []);

  if (error) {
    return <div className="mx-auto max-w-6xl p-12 text-center"><div className="mx-auto max-w-md rounded-2xl border border-red-200 dark:border-red-900/50 bg-red-50 dark:bg-red-950/30 p-6 text-sm text-red-700 dark:text-red-300">We couldn't reach the database right now. Try again in a moment.</div></div>;
  }
  if (!stats) return <div className="p-12 flex items-center justify-center gap-2 text-sm text-zinc-500 dark:text-zinc-400"><span className="h-2 w-2 animate-pulse rounded-full bg-zinc-400" />Looking at what's connected…</div>;

  const cards = [
    { label: "Code areas", hint: "Places where code lives", value: stats.repositories, icon: GitBranchIcon },
    { label: "Building blocks", hint: "Databases and libraries", value: stats.packages, icon: PackageIcon },
    { label: "Services", hint: "Things that keep running", value: stats.services, icon: WorkflowCircle06Icon },
    { label: "People", hint: "Who builds it", value: stats.developers, icon: UserGroupIcon },
    { label: "Teams", hint: "Who looks after it", value: stats.teams, icon: Analytics01Icon },
    { label: "Organizations", hint: "Groups of teams", value: stats.organizations, icon: Building02Icon },
  ];

  return (
    <div className="space-y-8">
      {/* One-sentence purpose — visible without scrolling */}
      <div className="overflow-hidden rounded-[28px] border border-black dark:border-white bg-white dark:bg-zinc-900 shadow-sm">
        <div className="bg-black dark:bg-white px-6 py-8 md:px-8 md:py-9 text-white dark:text-black relative overflow-hidden">
          <div className="absolute inset-0 bg-[linear-gradient(to_right,#ffffff0a_1px,transparent_1px),linear-gradient(to_bottom,#ffffff0a_1px,transparent_1px)] dark:bg-[linear-gradient(to_right,#00000008_1px,transparent_1px),linear-gradient(to_bottom,#00000008_1px,transparent_1px)] bg-[size:32px_32px]" />
          <div className="relative">
            <div className="inline-flex items-center gap-2 rounded-full border border-white/15 dark:border-black/10 bg-white/10 dark:bg-black/5 px-3 py-1 text-xs font-medium backdrop-blur"><span className="h-1.5 w-1.5 rounded-full bg-white dark:bg-black animate-pulse" />Connected {health ? "• live" : ""}</div>
            <h1 className="mt-3 text-3xl font-bold tracking-tight md:text-4xl" style={{ fontFamily: "Space Grotesk, sans-serif" }}>See what breaks if you change one thing — before you change it.</h1>
            <p className="mt-3 max-w-2xl text-sm leading-relaxed text-zinc-400 dark:text-zinc-500">Pick any database, building block or service. We’ll show you everything that relies on it, and the chain that connects them — in plain language, with a picture.</p>
            <div className="mt-6 flex flex-wrap items-center gap-3">
              <Link to="/explore" className="inline-flex items-center gap-2 rounded-full bg-white dark:bg-black px-6 py-3 text-sm font-semibold text-black dark:text-white shadow-lg hover:bg-zinc-100 dark:hover:bg-zinc-900 transition">
                Start exploring <HugeiconsIcon icon={ArrowRight01Icon} size={16} strokeWidth={2.5} color="currentColor" />
              </Link>
              <span className="text-xs text-white/60 dark:text-black/60">No setup needed — try “PostgreSQL” on the next page.</span>
            </div>
          </div>
        </div>
        <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 divide-x divide-zinc-100 dark:divide-zinc-800 border-t dark:border-zinc-800 bg-zinc-50 dark:bg-black">
          {cards.map((c) => (
            <div key={c.label} className="p-5 bg-white dark:bg-zinc-900">
              <div className="inline-flex h-9 w-9 items-center justify-center rounded-xl bg-black dark:bg-white text-white dark:text-black"><HugeiconsIcon icon={c.icon} size={16} color="currentColor" strokeWidth={2} /></div>
              <div className="mt-3 text-2xl font-bold tracking-tight text-black dark:text-white" style={{ fontFamily: "Space Grotesk, sans-serif" }}>{c.value}</div>
              <div className="text-xs font-semibold text-black dark:text-white">{c.label}</div>
              <div className="text-[11px] text-zinc-500 dark:text-zinc-400">{c.hint} — what this counts</div>
            </div>
          ))}
        </div>
      </div>

      <div className="grid md:grid-cols-2 gap-6">
        <div className="rounded-2xl border border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900 p-6 shadow-sm">
          <div className="flex items-center justify-between mb-2">
            <h3 className="flex items-center gap-2 text-sm font-semibold text-black dark:text-white"><HugeiconsIcon icon={PackageIcon} size={16} /> Most used building blocks</h3>
            <span className="rounded-full bg-black dark:bg-white px-2.5 py-1 text-xs font-medium text-white dark:text-black">{stats.mostConnected.length}</span>
          </div>
          <p className="text-xs text-zinc-500 dark:text-zinc-400 mb-4">Building blocks that many code areas rely on — change these and a lot could be affected.</p>
          <div className="space-y-2">
            {stats.mostConnected.map((m) => (
              <Link key={m.id} to={`/entities/Package/${m.id}`} className="flex items-center justify-between rounded-xl border border-zinc-100 dark:border-zinc-800 bg-zinc-50 dark:bg-zinc-800/50 px-4 py-3 hover:bg-white dark:hover:bg-zinc-800 hover:border-zinc-200 dark:hover:border-zinc-700 hover:shadow-sm transition group">
                <span className="flex items-center gap-3 text-sm font-medium text-black dark:text-white"><span className="h-7 w-7 rounded-full bg-black dark:bg-white flex items-center justify-center text-white dark:text-black text-xs font-bold">{m.name.slice(0,2).toUpperCase()}</span>{m.name}</span>
                <span className="flex items-center gap-2 text-xs text-zinc-500 dark:text-zinc-400">{m.count} code areas rely on it <HugeiconsIcon icon={ArrowRight01Icon} size={12} className="opacity-0 group-hover:opacity-100 transition" /></span>
              </Link>
            ))}
            {stats.mostConnected.length === 0 && <div className="text-sm text-zinc-400">No data yet — seed the database to see examples.</div>}
          </div>
        </div>
        <div className="rounded-2xl border border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900 p-6 shadow-sm">
          <h3 className="flex items-center gap-2 text-sm font-semibold text-black dark:text-white mb-1"><HugeiconsIcon icon={GitBranchIcon} size={16} /> How things connect — recent examples</h3>
          <p className="text-xs text-zinc-500 dark:text-zinc-400 mb-4">Plain view of how a code area relies on a building block, or powers a service.</p>
          <div className="space-y-2.5">
            {stats.recent.map((r, i) => (
              <div key={i} className="flex items-center gap-2 rounded-xl border border-zinc-100 dark:border-zinc-800 px-3 py-2.5 text-sm bg-zinc-50/70 dark:bg-zinc-800/30">
                <span className="font-medium text-black dark:text-white">{r.source}</span>
                <span className="inline-flex items-center gap-1 rounded-full bg-white dark:bg-zinc-800 border dark:border-zinc-700 px-2 py-0.5 text-xs font-medium text-zinc-600 dark:text-zinc-300">{r.type === "CALLS" ? "talks to" : r.type === "DEPLOYED_AS" ? "runs as" : "relies on"} →</span>
                <span className="font-medium text-black dark:text-white">{r.target}</span>
                {r.detail && <span className="ml-auto text-xs text-zinc-500">({r.detail})</span>}
              </div>
            ))}
            {stats.recent.length === 0 && <div className="text-sm text-zinc-400 dark:text-zinc-500">No recent connections to show. Try exploring PostgreSQL.</div>}
          </div>
        </div>
      </div>
      <p className="text-center text-xs text-zinc-400 dark:text-zinc-500">Tip: On the next page, try searching “PostgreSQL” or “Checkout Service” — then hit the big “See what would be affected” button.</p>
    </div>
  );
}
