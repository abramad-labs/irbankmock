import { dirname } from "path";
import { fileURLToPath } from "url";
import { FlatCompat } from "@eslint/eslintrc";

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

const compat = new FlatCompat({
  baseDirectory: __dirname,
});

const eslintConfig = [
  ...compat.extends("next/core-web-vitals", "next/typescript"),
  ...compat.config({
    plugins: "baseui",
    rules: {
      'baseui/deprecated-theme-api': "warn",
      'baseui/deprecated-component-api': "warn",
      'baseui/no-deep-imports': "warn",
    }
  })
];

export default eslintConfig;
