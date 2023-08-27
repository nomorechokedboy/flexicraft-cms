/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_GITHUB_OAUTH_CLIENT_ID: string;
  readonly VITE_GITHUB_OAUTH_CLIENT_SECRET: string;
  readonly VITE_GITHUB_OAUTH_ENDPOINT: string;
  readonly VITE_GITHUB_OAUTH_REDIRECT: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
