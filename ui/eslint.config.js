const js = require('@eslint/js');
const nextConfig = require('eslint-config-next');
const unusedImports = require('eslint-plugin-unused-imports');
const prettier = require('eslint-plugin-prettier');
const prettierConfig = require('eslint-config-prettier');

module.exports = [
  {
    ignores: [
      'node_modules/**',
      'out/**',
      '.next/**',
      'static/**',
      'package.json',
      'package-lock.json',
      'nodemon.json',
      'playwright-report/**',
      'test-results/**',
      'lib/**',
      'tests/samples/**',
      '**/__generated__/**',
      '../server/**',
    ],
  },
  js.configs.recommended,
  ...nextConfig,
  {
    files: ['**/*.{js,jsx,ts,tsx}'],
    languageOptions: {
      ecmaVersion: 'latest',
      sourceType: 'module',
      parserOptions: {
        ecmaFeatures: {
          jsx: true,
        },
      },
      globals: {
        Atomics: 'readonly',
        SharedArrayBuffer: 'readonly',
        globalThis: 'readonly',
      },
    },
    settings: {
      react: {
        version: 'detect',
      },
    },
    plugins: {
      'unused-imports': unusedImports,
      prettier: prettier,
    },
    rules: {
      '@next/next/no-img-element': 'off',
      'react-hooks/rules-of-hooks': 'warn',
      'react-hooks/exhaustive-deps': 'off',
      'jsx-a11y/alt-text': 'off',
      'valid-typeof': 'warn',
      'no-trailing-spaces': 'error',
      'block-spacing': 'error',
      'brace-style': ['error', '1tbs'],
      'react/react-in-jsx-scope': 'off',
      'no-undef': 'error',
      'react/jsx-uses-vars': ['error'],
      'react/jsx-no-undef': 'error',
      'no-console': 'off',
      'no-unused-vars': 'error',
      'unused-imports/no-unused-imports': 'error',
      'react/jsx-key': 'warn',
      'no-dupe-keys': 'error',
      'react/jsx-filename-extension': [
        'warn',
        {
          extensions: ['.js', '.jsx'],
        },
      ],
      'react/prop-types': 'off',
      'prettier/prettier': ['error', { endOfLine: 'lf' }],
    },
  },
  prettierConfig,
];
