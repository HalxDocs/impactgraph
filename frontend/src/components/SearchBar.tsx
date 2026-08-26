import { useEffect, useState } from "react";
import { HugeiconsIcon } from "@hugeicons/react";
import { Search01Icon, Cancel01Icon } from "@hugeicons/core-free-icons";
import { api } from "../api/client";
import type { SearchResult } from "../types/graph";

const typeDot: Record<string, string> = {
  Package: "bg-black dark:bg-white",
  Repository: "bg-zinc-800 dark:bg-zinc-200",
  Service: "bg-zinc-700 dark:bg-zinc-300",
  Developer: "bg-zinc-500 dark:bg-zinc-500",
  Team: "bg-zinc-600 dark:bg-zinc-400",
  Organization: "bg-zinc-900 dark:bg-zinc-100",
};

const typeHelp: Record<string, string> = {
  Package: "building block",
  Repository: "code area",
  Service: "service",
  Developer: "person",
  Team: "team",
  Organization: "organization",
};

export default function SearchBar({ onSelect }: { onSelect: (r: SearchResult) => void }) {
  const [q, setQ] = useState("");
  const [results, setResults] = useState<SearchResult[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [touched, setTouched] = useState(false);

  useEffect(() => {
    if (q.trim().length < 1) {
      setResults([]);
      if (touched && q.trim() === "" ) setError(null);
      return;
    }
    const t = setTimeout(async () => {
      setLoading(true);
      setError(null);
      try {
        const res = await api.search(q);
        setResults(res);
        if (res.length === 0) setError(null);
      } catch (e: any) {
        setError("We couldn't reach the database right now. Try again in a moment.");
      } finally {
        setLoading(false);
      }
    }, 300);
    return () => clearTimeout(t);
  }, [q, touched]);

  const showNoResults = touched && q.trim().length > 0 && !loading && results.length === 0 && !error;

  return (
    <div className="relative">
      <div className="relative group">
        <div className="pointer-events-none absolute left-4 top-1/2 -translate-y-1/2 text-zinc-400 group-focus-within:text-black dark:group-focus-within:text-white transition-colors">
          <HugeiconsIcon icon={Search01Icon} size={18} strokeWidth={2} />
        </div>
        <input
          value={q}
          onChange={(e) => { setQ(e.target.value); setTouched(true); }}
          placeholder="Try searching 'PostgreSQL' or 'Checkout Service'"
          aria-label="Search for a building block or service"
          className="w-full rounded-full border border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900 py-3.5 pl-11 pr-12 text-sm dark:text-zinc-100 shadow-sm outline-none placeholder:text-zinc-400 focus:border-black dark:focus:border-white focus:ring-4 focus:ring-black/5 dark:focus:ring-white/10 transition-all"
        />
        <div className="absolute right-2 top-1/2 -translate-y-1/2 flex items-center gap-2">
          {loading && <span className="text-xs text-zinc-500 dark:text-zinc-400 pr-2">Looking at what's connected…</span>}
          {q && (
            <button onClick={() => { setQ(""); setResults([]); }} aria-label="Clear search" className="inline-flex h-7 w-7 items-center justify-center rounded-full bg-zinc-100 dark:bg-zinc-800 text-zinc-500 dark:text-zinc-400 hover:bg-black hover:text-white dark:hover:bg-white dark:hover:text-black transition">
              <HugeiconsIcon icon={Cancel01Icon} size={12} strokeWidth={2.5} />
            </button>
          )}
        </div>
      </div>
      <p className="mt-2 text-xs text-zinc-500 dark:text-zinc-400">Shows building blocks, code areas, services and teams — not raw IDs. Example searches: <span className="font-medium text-black dark:text-white">PostgreSQL</span>, <span className="font-medium text-black dark:text-white">Redis</span>, <span className="font-medium text-black dark:text-white">Checkout Service</span></p>
      {error && <div className="mt-3 rounded-xl border border-red-200 dark:border-red-900/50 bg-red-50 dark:bg-red-950/30 px-4 py-3 text-sm text-red-700 dark:text-red-300">{error}</div>}
      {showNoResults && (
        <div className="mt-3 rounded-2xl border border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900 p-4 shadow-sm">
          <p className="text-sm font-medium text-black dark:text-white">We couldn't find anything matching '{q}' — try a different name</p>
          <p className="mt-1 text-xs text-zinc-500 dark:text-zinc-400">Try: PostgreSQL, Redis, React, Payments, Checkout Service, Notification Service, Infrastructure Team</p>
        </div>
      )}
      {results.length > 0 && (
        <div className="absolute z-20 mt-3 max-h-80 w-full overflow-hidden rounded-2xl border border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900 shadow-xl shadow-zinc-900/10">
          <div className="px-3 py-2 text-[11px] font-semibold tracking-widest text-zinc-400 uppercase">Results • {results.length} — building blocks and services that match</div>
          <div className="max-h-72 overflow-auto pb-2">
            {results.map((r) => (
              <button
                key={r.id}
                onClick={() => {
                  onSelect(r);
                  setQ(r.name);
                  setResults([]);
                }}
                className="flex w-full items-center gap-3 px-3 py-2.5 text-left hover:bg-zinc-50 dark:hover:bg-zinc-800 transition group"
              >
                <span className={`h-8 w-8 shrink-0 rounded-full ${typeDot[r.type] || "bg-zinc-400"} flex items-center justify-center text-white dark:text-black text-[10px] font-bold border border-black/10 dark:border-white/10`}>{r.name.slice(0,2).toUpperCase()}</span>
                <div className="min-w-0 flex-1">
                  <div className="text-sm font-medium text-black dark:text-white group-hover:text-black dark:group-hover:text-white truncate">{r.name} <span className="text-xs font-normal text-zinc-500">— {typeHelp[r.type] || r.type.toLowerCase()}</span></div>
                  <div className="text-xs text-zinc-500 dark:text-zinc-400 truncate">{r.id}</div>
                </div>
                <span className="shrink-0 rounded-full border dark:border-zinc-700 bg-white dark:bg-zinc-800 px-2.5 py-1 text-xs font-medium text-zinc-600 dark:text-zinc-300">{typeHelp[r.type] || r.type}</span>
              </button>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
