module.exports = {
  env: {
    browser: true,
    es6: true,
    node: true,
  },
  settings: {
    react: {
      version: require('./package.json').dependencies.react,
    },
  },
  parser: '@typescript-eslint/parser',
  parserOptions: {
    ecmaFeatures: { jsx: true },
    ecmaVersion: 2018,
    sourceType: 'module',
    project: './tsconfig.json', // optional, for stricter TS rules
  },
  plugins: ['react', 'prettier', 'unused-imports', '@typescript-eslint'],
  extends: [
    'eslint:recommended',
    'plugin:react/recommended',
    'plugin:@typescript-eslint/recommended',
    'next',
    'plugin:prettier/recommended',
  ],
  rules: {
    // Common rules for all files
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
    'react/jsx-uses-vars': [2],
    'react/jsx-no-undef': 'error',
    'no-console': 0,
    'no-unused-vars': 'off', // handled per filetype
    'unused-imports/no-unused-imports': 'error',
    'react/jsx-key': 'warn',
    'no-dupe-keys': 'error',
    'react/jsx-filename-extension': [
      1,
      {
        extensions: ['.js', '.jsx', '.ts', '.tsx'],
      },
    ],
    'react/prop-types': 'off',
    'prettier/prettier': ['error', { endOfLine: 'lf' }],
  },
  overrides: [
    // TypeScript files
    {
      files: ['*.ts', '*.tsx'],
      rules: {
        '@typescript-eslint/no-unused-vars': 'error',
        '@typescript-eslint/no-require-imports': 'error',
        '@typescript-eslint/ban-ts-comment': 'warn',
      },
    },
    // JavaScript files
    {
      files: ['*.js', '*.jsx'],
      rules: {
        '@typescript-eslint/no-unused-vars': 'off',
        '@typescript-eslint/no-require-imports': 'off',
        '@typescript-eslint/ban-ts-comment': 'off',
        'no-unused-vars': 'error',
        'no-unused-expressions': 'off',
        '@typescript-eslint/no-unused-expressions': 'off',
      },
    },
  ],
};