import { HugeiconsIcon } from "@hugeicons/react";
import { CubeIcon } from "@hugeicons/core-free-icons";

export default function EntityCard({ type, name, id, description }: { type: string; name: string; id: string; description?: string }) {
  return (
    <div className="rounded-2xl border border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900 p-5 shadow-sm">
      <div className="inline-flex items-center gap-1.5 rounded-full border dark:border-zinc-700 bg-zinc-100 dark:bg-zinc-800 px-2.5 py-1 text-xs font-medium text-zinc-600 dark:text-zinc-300"><HugeiconsIcon icon={CubeIcon} size={12} />{type}</div>
      <div className="mt-2 text-sm font-semibold text-black dark:text-white">{name}</div>
      <div className="text-xs font-mono text-zinc-500 dark:text-zinc-400 truncate">{id}</div>
      {description && <div className="mt-1 text-xs text-zinc-600 dark:text-zinc-300">{description}</div>}
    </div>
  );
}
