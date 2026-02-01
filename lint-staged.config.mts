import type { Configuration } from "lint-staged";

const config: Configuration = {
  "*.{js,jsx,ts,tsx,json,css,html}": ["pnpm run format", "pnpm run lint"],
  //   "*.java": (_) => ["pnpm --filter @yu-mi-box/api run format:java"],
};

export default config;
