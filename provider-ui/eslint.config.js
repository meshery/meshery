const js = require('@eslint/js');
const next = require('eslint-config-next');
const globals = require('globals');


// ESLint 10: eslint-config-next's babel-based parser returns a scope manager that
// doesn't implement addGlobals (new ESLint 10 API). Replace it with espree (ESLint's
// built-in parser) for JS/JSX files; the TS entry already uses @typescript-eslint/parser.
const patchedNextConfig = next.map((cfg) => {
  if (cfg.name === 'next') {
    const { parser, globals, ...restLangOpts } = cfg.languageOptions ?? {};
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

const config = [
  {
    ignores: [
      'node_modules/**',
      'out/**',
      '.next/**',
      'package.json',
      'package-lock.json',
      'nodemon.json',
    ],
  },
  js.configs.recommended,
  ...patchedNextConfig,
  {
    files: ['**/*.{js,jsx,ts,tsx}'],
    languageOptions: {
      ecmaVersion: 'latest',
      sourceType: 'module',
      globals: {
        ...globals.browser,
        ...globals.node,
      },
      parserOptions: {
        ecmaFeatures: {
          jsx: true,
        },
      },
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
      'array-bracket-spacing': ['error', 'never'],
      'comma-style': ['error'],
      'jsx-quotes': ['error', 'prefer-double'],
      'block-scoped-var': 'error',
      'keyword-spacing': 'error',
      'no-trailing-spaces': 'error',
      'object-curly-spacing': ['error', 'always'],
      'arrow-spacing': [
        'error',
        {
          before: true,
          after: true,
        },
      ],
      'key-spacing': [
        'error',
        {
          beforeColon: false,
          afterColon: true,
        },
      ],
      'block-spacing': 'error',
      'brace-style': ['error', '1tbs'],
      indent: [
        'error',
        2,
        {
          FunctionExpression: { parameters: 'first' },
          FunctionDeclaration: { parameters: 'first' },
          MemberExpression: 1,
          SwitchCase: 1,
          outerIIFEBody: 0,
          VariableDeclarator: { var: 2, let: 2, const: 3 },
          ignoredNodes: ['TemplateLiteral'],
        },
      ],
      'react/react-in-jsx-scope': 'off',
      'no-undef': 'error',
      'react/jsx-uses-vars': ['error'],
      'react/jsx-no-undef': 'error',
      'no-console': 'off',
      'no-unused-vars': 'error',
      'react/jsx-key': 'warn',
      'no-dupe-keys': 'error',
      "react-hooks/immutability": "off",
      'react/prop-types': 'off',
    },
  },

  {
    files: ["eslint.config.js"],
    rules: {
      "no-unused-vars": "off",
    },
  },
];
module.exports = config;