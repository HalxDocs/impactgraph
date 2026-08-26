import { HugeiconsIcon } from "@hugeicons/react";
import type { FC } from "react";

export default function StatCard({ label, value, hint, icon, className = "" }: { label: string; value: number | string; hint?: string; icon?: any; className?: string }) {
  return (
    <div className={`rounded-2xl border border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900 p-5 shadow-sm ${className}`}>
      {icon && <div className="inline-flex h-9 w-9 items-center justify-center rounded-xl bg-black dark:bg-white text-white dark:text-black"><HugeiconsIcon icon={icon} size={16} color="currentColor" /></div>}
      <div className="mt-3 text-2xl font-bold tracking-tight text-black dark:text-white" style={{ fontFamily: "Space Grotesk, sans-serif" }}>{value}</div>
      <div className="text-xs font-semibold text-black dark:text-white">{label}</div>
      {hint && <div className="text-[11px] text-zinc-500 dark:text-zinc-400">{hint}</div>}
    </div>
  );
}
