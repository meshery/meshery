import js from '@eslint/js';
import { FlatCompat } from '@eslint/eslintrc';

const compat = new FlatCompat({
  recommendedConfig: js.configs.recommended,
});

export default [
  {
    ignores: [
      'node_modules/**',
      'out/**',
      '.next/**',
      'package.json',
      'package-lock.json',
      'nodemon.json'
    ]
  },
  js.configs.recommended,
  ...compat.extends(
    'next/core-web-vitals',
    'plugin:react/recommended',
    'plugin:react-hooks/recommended'
  ),
  {
    files: ['**/*.{js,jsx,ts,tsx}'],
    languageOptions: {
      ecmaVersion: 'latest',
      sourceType: 'module',
      parser: (await import('@babel/eslint-parser')).default,
      parserOptions: {
        ecmaFeatures: {
          jsx: true,
        },
        requireConfigFile: false,
        babelOptions: {
          presets: ['@babel/preset-react'],
        },
      },
    },
    settings: {
      react: {
        version: 'detect',
      },
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
      'react/jsx-filename-extension': [
        'warn',
        {
          extensions: ['.js', '.jsx'],
        },
      ],
      'react/prop-types': 'off',
    },
  },
];
