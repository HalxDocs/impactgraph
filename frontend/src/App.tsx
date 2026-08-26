import { BrowserRouter, Routes, Route, NavLink, Link } from "react-router-dom";
import { useState } from "react";
import { HugeiconsIcon } from "@hugeicons/react";
import {
  DashboardSquare02Icon,
  Search01Icon,
  GitBranchIcon,
  WorkflowCircle06Icon,
  BookOpen02Icon,
  ArrowRight01Icon,
  Menu01Icon,
  Cancel01Icon,
  FlashIcon,
  Sun02Icon,
  Moon02Icon,
} from "@hugeicons/core-free-icons";
import { useTheme } from "./hooks/useTheme";
import Dashboard from "./pages/Dashboard";
import Explore from "./pages/Explore";
import Entity from "./pages/Entity";
import Docs from "./pages/Docs";

function Navbar() {
  const [open, setOpen] = useState(false);
  const { theme, toggle } = useTheme();
  const pillLink = ({ isActive }: { isActive: boolean }) =>
    `flex items-center gap-1.5 rounded-full px-3.5 py-1.5 text-[13px] font-medium transition-all ${
      isActive ? "bg-black text-white dark:bg-white dark:text-black shadow-sm" : "text-zinc-500 dark:text-zinc-400 hover:text-black dark:hover:text-white hover:bg-zinc-100 dark:hover:bg-zinc-800"
    }`;

  return (
    <header className="sticky top-0 z-40 border-b border-zinc-200 dark:border-zinc-800 bg-white/95 dark:bg-black/95 backdrop-blur-xl">
      <div className="mx-auto flex h-[64px] max-w-[1280px] items-center justify-between px-4 md:px-6">
        <Link to="/" className="flex items-center gap-2.5">
          <div className="flex h-8 w-8 items-center justify-center rounded-xl bg-black dark:bg-white text-white dark:text-black">
            <HugeiconsIcon icon={FlashIcon} size={18} color="currentColor" strokeWidth={2} />
          </div>
          <span className="text-[15px] font-semibold tracking-tight text-black dark:text-white" style={{ fontFamily: "Space Grotesk, sans-serif" }}>ImpactGraph</span>
          <span className="hidden lg:inline text-xs font-medium text-zinc-500">• Dependency Impact</span>
        </Link>

        <nav className="hidden md:flex items-center gap-1 rounded-full border border-zinc-200 dark:border-zinc-800 bg-zinc-100 dark:bg-zinc-900 p-1">
          <NavLink to="/" className={pillLink}>
            <HugeiconsIcon icon={DashboardSquare02Icon} size={14} strokeWidth={2} />
            Dashboard
          </NavLink>
          <NavLink to="/explore" className={pillLink}>
            <HugeiconsIcon icon={Search01Icon} size={14} strokeWidth={2} />
            Explore
          </NavLink>
          <NavLink to="/explore" onClick={() => setTimeout(() => document.getElementById("find-connection")?.scrollIntoView({ behavior: "smooth" }), 100)} className="flex items-center gap-1.5 rounded-full px-3.5 py-1.5 text-[13px] font-medium text-zinc-500 dark:text-zinc-400 hover:text-black dark:hover:text-white hover:bg-zinc-100 dark:hover:bg-zinc-800 transition">
            <HugeiconsIcon icon={GitBranchIcon} size={14} strokeWidth={2} />
            Paths
          </NavLink>
          <NavLink to="/docs" className={({ isActive }) => `flex items-center gap-1.5 rounded-full px-3.5 py-1.5 text-[13px] font-medium transition-all ${isActive ? "bg-black text-white dark:bg-white dark:text-black shadow-sm" : "text-zinc-500 dark:text-zinc-400 hover:text-black dark:hover:text-white hover:bg-zinc-100 dark:hover:bg-zinc-800"}`}>
            <HugeiconsIcon icon={BookOpen02Icon} size={14} strokeWidth={2} />
            Docs
          </NavLink>
        </nav>

        <div className="hidden md:flex items-center gap-3">
          <button
            onClick={toggle}
            aria-label="Toggle theme"
            className="inline-flex h-9 w-9 items-center justify-center rounded-full border border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900 text-zinc-600 dark:text-zinc-400 hover:bg-black hover:text-white dark:hover:bg-white dark:hover:text-black transition"
            title={theme === "dark" ? "Switch to light" : "Switch to dark"}
          >
            <HugeiconsIcon icon={theme === "dark" ? Sun02Icon : Moon02Icon} size={16} strokeWidth={2} />
          </button>
          <a href="#" className="text-sm font-medium text-zinc-500 dark:text-zinc-400 hover:text-black dark:hover:text-white transition">Sign in</a>
          <Link
            to="/explore"
            className="inline-flex items-center gap-1.5 rounded-full bg-black dark:bg-white px-5 py-2 text-sm font-semibold text-white dark:text-black shadow hover:bg-zinc-800 dark:hover:bg-zinc-100 transition"
          >
            Get started
            <HugeiconsIcon icon={ArrowRight01Icon} size={14} strokeWidth={2.5} color="currentColor" />
          </Link>
        </div>

        <div className="flex items-center gap-2 md:hidden">
          <button
            onClick={toggle}
            aria-label="Toggle theme"
            className="inline-flex h-9 w-9 items-center justify-center rounded-full border border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900 text-zinc-600 dark:text-zinc-400"
          >
            <HugeiconsIcon icon={theme === "dark" ? Sun02Icon : Moon02Icon} size={16} />
          </button>
          <button onClick={() => setOpen(!open)} className="inline-flex h-9 w-9 items-center justify-center rounded-full border border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900 text-black dark:text-white">
            <HugeiconsIcon icon={open ? Cancel01Icon : Menu01Icon} size={18} />
          </button>
        </div>
      </div>

      {open && (
        <div className="md:hidden border-t border-zinc-200 dark:border-zinc-800 bg-white dark:bg-black px-4 py-4">
          <nav className="flex flex-col gap-2">
            <NavLink to="/" onClick={() => setOpen(false)} className={({ isActive }) => `flex items-center gap-2 rounded-xl px-4 py-3 text-sm font-medium ${isActive ? "bg-black text-white dark:bg-white dark:text-black" : "text-zinc-600 dark:text-zinc-300 bg-zinc-50 dark:bg-zinc-900"}`}>
              <HugeiconsIcon icon={DashboardSquare02Icon} size={16} /> Dashboard
            </NavLink>
            <NavLink to="/explore" onClick={() => setOpen(false)} className={({ isActive }) => `flex items-center gap-2 rounded-xl px-4 py-3 text-sm font-medium ${isActive ? "bg-black text-white dark:bg-white dark:text-black" : "text-zinc-600 dark:text-zinc-300 bg-zinc-50 dark:bg-zinc-900"}`}>
              <HugeiconsIcon icon={Search01Icon} size={16} /> Explore
            </NavLink>
            <NavLink to="/docs" onClick={() => setOpen(false)} className={({ isActive }) => `flex items-center gap-2 rounded-xl px-4 py-3 text-sm font-medium ${isActive ? "bg-black text-white dark:bg-white dark:text-black" : "text-zinc-600 dark:text-zinc-300 bg-zinc-50 dark:bg-zinc-900"}`}>
              <HugeiconsIcon icon={BookOpen02Icon} size={16} /> Docs
            </NavLink>
            <Link to="/explore" onClick={() => setOpen(false)} className="flex items-center justify-center gap-1.5 rounded-full bg-black dark:bg-white px-5 py-3 text-sm font-semibold text-white dark:text-black">
              Get started <HugeiconsIcon icon={ArrowRight01Icon} size={14} color="currentColor" />
            </Link>
          </nav>
        </div>
      )}
    </header>
  );
}

function Layout({ children }: { children: React.ReactNode }) {
  return (
    <div className="min-h-screen bg-white dark:bg-black transition-colors">
      <Navbar />
      <main className="mx-auto max-w-6xl px-4 md:px-6 py-8">{children}</main>
      <footer className="mt-12 border-t border-zinc-200 dark:border-zinc-800 bg-white dark:bg-black">
        <div className="mx-auto max-w-6xl px-4 md:px-6 py-6 flex flex-col md:flex-row items-center justify-between gap-3 text-xs text-zinc-500 dark:text-zinc-400">
          <span className="flex items-center gap-2"><span className="h-6 w-6 rounded-lg bg-black dark:bg-white text-white dark:text-black flex items-center justify-center text-[10px] font-bold">IG</span> If this dependency changes, what else could be affected?</span>
          <span className="flex items-center gap-3">React → Go (Gin) → CognoDB • ImpactGraph • <span className="rounded-full bg-black dark:bg-white px-2.5 py-1 text-white dark:text-black font-medium text-[11px]">Wexa Take-Home</span></span>
        </div>
      </footer>
    </div>
  );
}

export default function App() {
  return (
    <BrowserRouter>
      <Layout>
        <Routes>
          <Route path="/" element={<Dashboard />} />
          <Route path="/explore" element={<Explore />} />
          <Route path="/docs" element={<Docs />} />
          <Route path="/entity/:type/:id" element={<Entity />} />
          <Route path="/entities/:type/:id" element={<Entity />} />
        </Routes>
      </Layout>
    </BrowserRouter>
  );
}
