module.exports = {
  settings: {
    react: {
      version: require("./package.json").dependencies.react,
    },
  },
  env: {
    browser: true,
    es2021: true,
    node: true,
  },
  extends: [
    "eslint:recommended",
    "next",
    "plugin:react/recommended",
    "plugin:@next/next/recommended",
  ],
  overrides: [],
  globals: {
    Atomics: "readonly",
    SharedArrayBuffer: "readonly",
    globalThis: "readonly",
  },
  parser: "@babel/eslint-parser",
  parserOptions: {
    ecmaFeatures: {
      jsx: true,
    },
    ecmaVersion: "latest",
    requireConfigFIle: false,
    sourceType: "module",
  },
  plugins: ["react"],
  rules: {
    "@next/next/no-img-element": "off",
    "react-hooks/rules-of-hooks": "warn",
    "react-hooks/exhaustive-deps": "off",
    "jsx-a11y/alt-text": "off",
    "valid-typeof": "warn",
    "array-bracket-spacing": ["error", "never"],
    "comma-style": ["error"],
    "jsx-quotes": ["error", "prefer-double"],
    "block-scoped-var": "error",
    "keyword-spacing": "error",
    "no-trailing-spaces": "error",
    "object-curly-spacing": ["error", "always"],
    "arrow-spacing": [
      "error",
      {
        before: true,
        after: true,
      },
    ],
    "key-spacing": [
      "error",
      {
        beforeColon: false,
        afterColon: true,
      },
    ],
    "block-spacing": "error",
    "brace-style": ["error", "1tbs"],
    indent: [
      "error",
      2,
      {
        FunctionExpression: { parameters: "first" },
        FunctionDeclaration: { parameters: "first" },
        MemberExpression: 1,
        SwitchCase: 1,
        outerIIFEBody: 0,
        VariableDeclarator: { var: 2, let: 2, const: 3 },
        ignoredNodes: ["TemplateLiteral"],
      },
    ],
    "react/react-in-jsx-scope": "off",
    "no-undef": "error",
    "react/jsx-uses-vars": [2],
    "react/jsx-no-undef": "error",
    "no-console": 0,
    "no-unused-vars": "error",
    "react/jsx-key": "warn",
    "no-dupe-keys": "error",
    "react/jsx-filename-extension": [
      1,
      {
        extensions: [".js", ".jsx"],
      },
    ],
    "react/prop-types": "off",
  },
};
