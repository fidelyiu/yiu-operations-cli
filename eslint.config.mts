import js from "@eslint/js";
import globals from "globals";
import tseslint from "typescript-eslint";
import pluginReact from "eslint-plugin-react";
import { defineConfig, globalIgnores } from "eslint/config";
import eslintPluginPrettierRecommended from "eslint-plugin-prettier/recommended";
import importPulgin from "eslint-plugin-import";
import eslintConfigPrettier from "eslint-config-prettier";

export default defineConfig([
  globalIgnores(["**/.docusaurus/", "**/build/"]),
  {
    files: ["**/*.{js,mjs,cjs,ts,mts,cts,jsx,tsx}"],
    plugins: { js },
    extends: ["js/recommended"],
    ignores: ["**/.docusaurus"],
    languageOptions: { globals: { ...globals.browser, ...globals.node } },
  },
  tseslint.configs.recommended,
  pluginReact.configs.flat.recommended,
  importPulgin.flatConfigs.typescript,
  eslintConfigPrettier,
  eslintPluginPrettierRecommended,
  {
    settings: {
      react: {
        version: "19.0",
      },
    },
    rules: {
      "react/react-in-jsx-scope": "off",
    },
  },
]);
