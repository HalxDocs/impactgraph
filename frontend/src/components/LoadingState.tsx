export default function LoadingState({ message = "Loading…" }: { message?: string }) {
  return (
    <div className="flex items-center justify-center gap-2 rounded-2xl border border-dashed border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900 p-8 text-sm text-zinc-500 dark:text-zinc-400">
      <span className="h-2 w-2 animate-pulse rounded-full bg-black dark:bg-white" />
      {message}
    </div>
  );
}
export function EmptyState({ title, desc }: { title: string; desc?: string }) {
  return (
    <div className="rounded-2xl border border-dashed border-zinc-200 dark:border-zinc-800 bg-zinc-50 dark:bg-zinc-900 p-8 text-center">
      <div className="text-sm font-medium text-black dark:text-white">{title}</div>
      {desc && <div className="mt-1 text-xs text-zinc-500 dark:text-zinc-400">{desc}</div>}
    </div>
  );
}
