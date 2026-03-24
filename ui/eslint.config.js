const nextConfig = require('eslint-config-next');
const prettierRecommended = require('eslint-plugin-prettier/recommended');
const unusedImports = require('eslint-plugin-unused-imports');
const globals = require('globals');

// ESLint 10: eslint-config-next's babel-based parser returns a scope manager that
// doesn't implement addGlobals (new ESLint 10 API). Replace it with espree (ESLint's
// built-in parser) for JS/JSX files; the TS entry already uses @typescript-eslint/parser.
const patchedNextConfig = nextConfig.map((cfg) => {
  if (cfg.name === 'next') {
    const { parser: _babelParser, globals: _g, ...restLangOpts } = cfg.languageOptions ?? {};
    return {
      ...cfg,
      languageOptions: {
        ...restLangOpts,
        parserOptions: {
          ...restLangOpts.parserOptions,
          ecmaFeatures: { jsx: true },
          ecmaVersion: 'latest',
          sourceType: 'module',
        },
      },
    };
  }
  return cfg;
});

module.exports = [
  // Global ignores (replaces .eslintignore — not supported in flat config)
  {
    ignores: [
      'node_modules/**',
      'out/**',
      '.next/**',
      'static/**',
      'lib/**',
      'tests/samples/**',
      '**/__generated__/**',
      'playwright-report/**',
      'test-results/**',
    ],
  },

  // Globals via default parser (avoids babel parser / addGlobals incompatibility)
  {
    languageOptions: {
      globals: {
        ...globals.browser,
        ...globals.node,
        Atomics: 'readonly',
        SharedArrayBuffer: 'readonly',
        globalThis: 'readonly',
      },
      parserOptions: {
        ecmaFeatures: { jsx: true },
        ecmaVersion: 'latest',
        sourceType: 'module',
      },
    },
  },

  // Next.js flat config (includes react, react-hooks, @next/next rules)
  ...patchedNextConfig,

  // Prettier integration (flat config format — disables conflicting style rules)
  prettierRecommended,

  // Custom overrides
  {
    plugins: {
      'unused-imports': unusedImports,
    },
    settings: {
      // eslint-plugin-react calls context.getFilename() during 'detect' (removed in ESLint 9+).
      // Provide an explicit version to skip detection entirely.
      react: { version: '19' },
    },
    rules: {
      '@next/next/no-img-element': 'off',
      'react-hooks/rules-of-hooks': 'warn',
      'react-hooks/exhaustive-deps': 'off',
      'jsx-a11y/alt-text': 'off',
      'valid-typeof': 'warn',
      'react/react-in-jsx-scope': 'off',
      'no-undef': 'error',
      'react/jsx-uses-vars': [2],
      'react/jsx-no-undef': 'error',
      'no-console': 0,
      'unused-imports/no-unused-imports': 'error',
      'react/jsx-key': 'warn',
      'no-dupe-keys': 'error',
      'react/prop-types': 'off',
      'prettier/prettier': ['error', { endOfLine: 'lf' }],
    },
  },

  // no-unused-vars: JS/JSX only — TypeScript files should use @typescript-eslint/no-unused-vars
  {
    files: ['**/*.{js,jsx,mjs,cjs}'],
    rules: {
      'no-unused-vars': ['error', { argsIgnorePattern: '^_', varsIgnorePattern: '^_' }],
    },
  },
];
