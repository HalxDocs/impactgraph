/// <reference types="vite/client" />
declare module "@hugeicons/core-free-icons";
interface ImportMetaEnv {
  readonly VITE_API_URL: string;
}
interface ImportMeta {
  readonly env: ImportMetaEnv;
}
