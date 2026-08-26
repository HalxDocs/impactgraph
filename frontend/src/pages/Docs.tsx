import { Link } from "react-router-dom";
import { HugeiconsIcon } from "@hugeicons/react";
import { BookOpen02Icon, FlashIcon, GitBranchIcon, InformationCircleIcon, ArrowRight01Icon, HelpCircleIcon, Search01Icon } from "@hugeicons/core-free-icons";

export default function Docs() {
  return (
    <div className="space-y-8 max-w-4xl mx-auto">
      <div className="flex items-center gap-2 text-xs text-zinc-500 dark:text-zinc-400">
        <Link to="/" className="hover:text-black dark:hover:text-white">Dashboard</Link>
        <span>›</span>
        <span className="rounded-full bg-black dark:bg-white px-2.5 py-1 text-white dark:text-black font-medium">Docs</span>
      </div>

      <div className="rounded-[28px] border border-black dark:border-white bg-black dark:bg-white p-8 md:p-10 text-white dark:text-black">
        <div className="inline-flex items-center gap-2 rounded-full border border-white/15 dark:border-black/10 bg-white/10 dark:bg-black/5 px-3 py-1 text-xs font-medium"><HugeiconsIcon icon={BookOpen02Icon} size={12} color="currentColor" /> Documentation — plain language</div>
        <h1 className="mt-4 text-3xl md:text-4xl font-bold tracking-tight" style={{ fontFamily: "Space Grotesk, sans-serif" }}>How ImpactGraph works — no jargon</h1>
        <p className="mt-3 max-w-2xl text-sm leading-relaxed text-white/70 dark:text-black/60">A non-technical guide. If you can read this page, you can use the app without explanation.</p>
        <div className="mt-6 flex flex-wrap gap-2 text-xs">
          <a href="#impact" className="rounded-full bg-white dark:bg-black px-4 py-2 font-medium text-black dark:text-white">Impact — what breaks?</a>
          <a href="#path" className="rounded-full border border-white/20 dark:border-black/15 bg-transparent px-4 py-2 font-medium text-white dark:text-black">Path — how two things connect</a>
          <a href="#use" className="rounded-full border border-white/20 dark:border-black/15 bg-transparent px-4 py-2 font-medium text-white dark:text-black">How to use</a>
        </div>
      </div>

      <div className="grid md:grid-cols-3 gap-4">
        <div className="rounded-2xl border border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900 p-5">
          <HugeiconsIcon icon={Search01Icon} size={18} className="text-black dark:text-white" />
          <h3 className="mt-2 text-sm font-semibold text-black dark:text-white">1. Search</h3>
          <p className="text-xs text-zinc-500 dark:text-zinc-400">Try “PostgreSQL” or “Checkout Service” — building blocks, code areas, services, teams.</p>
        </div>
        <div className="rounded-2xl border-2 border-black dark:border-white bg-white dark:bg-zinc-900 p-5">
          <HugeiconsIcon icon={FlashIcon} size={18} className="text-black dark:text-white" />
          <h3 className="mt-2 text-sm font-semibold text-black dark:text-white">2. See what would be affected</h3>
          <p className="text-xs text-zinc-500 dark:text-zinc-400">Big black button — highlights the whole chain in the picture.</p>
        </div>
        <div className="rounded-2xl border border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900 p-5">
          <HugeiconsIcon icon={GitBranchIcon} size={18} className="text-black dark:text-white" />
          <h3 className="mt-2 text-sm font-semibold text-black dark:text-white">3. Or find a connection</h3>
          <p className="text-xs text-zinc-500 dark:text-zinc-400">Pick two items — see the shortest chain between them, in plain words.</p>
        </div>
      </div>

      <section id="impact" className="rounded-2xl border border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900 overflow-hidden">
        <div className="border-b dark:border-zinc-800 bg-zinc-50/70 dark:bg-zinc-900/50 px-6 py-5 flex items-start gap-3">
          <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-black dark:bg-white text-white dark:text-black"><HugeiconsIcon icon={FlashIcon} size={16} color="currentColor" /></div>
          <div>
            <h2 className="text-lg font-bold text-black dark:text-white" style={{ fontFamily: "Space Grotesk, sans-serif" }}>Impact — “What breaks if I change this?”</h2>
            <p className="text-xs text-zinc-500 dark:text-zinc-400">Plain: Pick one building block. We show everything that relies on it. Technical: reverse DEPENDS_ON → DEPLOYED_AS → CALLS</p>
          </div>
        </div>
        <div className="p-6 space-y-4 text-sm leading-relaxed text-zinc-700 dark:text-zinc-300">
          <div className="rounded-xl bg-zinc-50 dark:bg-zinc-800/50 border dark:border-zinc-700 p-4">
            <p className="font-semibold text-black dark:text-white">In one sentence:</p>
            <p>If <span className="font-medium text-black dark:text-white">PostgreSQL</span> breaks, every code area that relies on it, every service that runs that code, and every service that talks to those services could stop working. We count them and show the longest chain — e.g. <span className="font-medium">4 steps away from PostgreSQL, worst case.</span></p>
          </div>
          <ul className="list-disc list-inside space-y-1 text-zinc-600 dark:text-zinc-400">
            <li><span className="font-medium text-black dark:text-white">Code areas: 10</span> — places where code lives that rely on it → “4 code areas could need fixes.”</li>
            <li><span className="font-medium text-black dark:text-white">Services: 12</span> — things that keep running → “12 services could stop working correctly.”</li>
            <li><span className="font-medium text-black dark:text-white">Teams: 6</span> — people who look after it → “6 teams would need to be told.”</li>
            <li><span className="font-medium text-black dark:text-white">Longest chain: 4 steps</span> — how far the chain reaches, worst case.</li>
          </ul>
          <p className="flex items-center gap-2 text-xs text-zinc-500 dark:text-zinc-400"><HugeiconsIcon icon={InformationCircleIcon} size={14} /> Highlighted black line in the picture is the affected chain. Left is the starting building block, right is what relies on it. Rest muted to 35% so cause→effect is glanceable.</p>
          <details className="rounded-xl border dark:border-zinc-700 bg-white dark:bg-zinc-800 p-4">
            <summary className="cursor-pointer text-xs font-semibold text-black dark:text-white flex items-center gap-2"><HugeiconsIcon icon={HelpCircleIcon} size={12} /> How we find this (for engineers) — Cypher</summary>
            <pre className="mt-3 overflow-auto rounded-lg bg-black dark:bg-zinc-900 text-white dark:text-zinc-100 p-3 text-xs mono">MATCH (p:Package &#123;id:$id&#125;){"\n"}MATCH path = (p)&lt;-[:DEPENDS_ON*0..3]-(:Repository)-[:DEPLOYED_AS*0..1]-&gt;(:Service)-[:CALLS*0..4]-&gt;(target:Service){"\n"}RETURN path, length(path) as depth</pre>
            <p className="mt-2 text-xs text-zinc-500">SQL would need recursive CTEs + UNIONs; Cypher is one traversal — the “why graph” story.</p>
          </details>
        </div>
      </section>

      <section id="path" className="rounded-2xl border border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900 overflow-hidden">
        <div className="border-b dark:border-zinc-800 bg-zinc-50/70 dark:bg-zinc-900/50 px-6 py-5 flex items-start gap-3">
          <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-black dark:bg-white text-white dark:text-black"><HugeiconsIcon icon={GitBranchIcon} size={16} color="currentColor" /></div>
          <div>
            <h2 className="text-lg font-bold text-black dark:text-white" style={{ fontFamily: "Space Grotesk, sans-serif" }}>Path — “How do these two connect?”</h2>
            <p className="text-xs text-zinc-500 dark:text-zinc-400">Plain: Pick any two items, we find the shortest chain. Human verbs: relies on → runs as → talks to</p>
          </div>
        </div>
        <div className="p-6 space-y-4 text-sm leading-relaxed text-zinc-700 dark:text-zinc-300">
          <div className="rounded-xl bg-zinc-50 dark:bg-zinc-800/50 border dark:border-zinc-700 p-4">
            <p className="font-semibold text-black dark:text-white">In one sentence:</p>
            <p>Choose <span className="font-medium">PostgreSQL</span> and <span className="font-medium">Notification Service</span> — we find <span className="font-medium">PostgreSQL → billing-service (relies on) → Payment Service (runs as) → Notification Service (talks to)</span> — narrated as “PostgreSQL connects to billing-service, which connects to Payment Service — that's the chain.”</p>
          </div>
          <p className="text-xs text-zinc-500 dark:text-zinc-400">In <span className="font-medium text-black dark:text-white">Explore</span> sidebar → “How two things connect” → pick Starting + Connected → <span className="font-medium">Find connection</span>. Picture updates, highlighted black chain, pills show `building block → code area → service`.</p>
          <details className="rounded-xl border dark:border-zinc-700 bg-white dark:bg-zinc-800 p-4">
            <summary className="cursor-pointer text-xs font-semibold text-black dark:text-white flex items-center gap-2"><HugeiconsIcon icon={HelpCircleIcon} size={12} /> Cypher behind it</summary>
            <pre className="mt-3 overflow-auto rounded-lg bg-black dark:bg-zinc-900 text-white dark:text-zinc-100 p-3 text-xs mono">MATCH (a &#123;id:$from&#125;), (b &#123;id:$to&#125;), path = shortestPath((a)-[:DEPENDS_ON|DEPLOYED_AS|CALLS|OWNS|BELONGS_TO|CONTRIBUTED_TO|MAINTAINS*..8]-(b)) RETURN path</pre>
            <p className="mt-2 text-xs text-zinc-500">If no path → plain “No connection found — try a different pair.” No raw error JSON.</p>
          </details>
        </div>
      </section>

      <section id="use" className="rounded-2xl border border-black dark:border-white bg-black dark:bg-white p-6 text-white dark:text-black">
        <h2 className="text-lg font-bold" style={{ fontFamily: "Space Grotesk, sans-serif" }}>Try it now — 30 seconds</h2>
        <ol className="mt-3 list-decimal list-inside text-sm space-y-1 text-white/80 dark:text-black/70">
          <li>Go to <Link to="/explore" className="underline font-medium text-white dark:text-black">Explore</Link> → type <span className="font-medium">PostgreSQL</span> → pick the building block</li>
          <li>Hit the big <span className="font-medium">See what would be affected</span> — read the sentences, not just numbers</li>
          <li>Scroll to “How two things connect” → pick <span className="font-medium">PostgreSQL</span> → <span className="font-medium">Notification Service</span> → Find connection</li>
        </ol>
        <Link to="/explore" className="mt-4 inline-flex items-center gap-2 rounded-full bg-white dark:bg-black px-6 py-2.5 text-sm font-semibold text-black dark:text-white">Go to Explore <HugeiconsIcon icon={ArrowRight01Icon} size={14} color="currentColor" /></Link>
      </section>
    </div>
  );
}
