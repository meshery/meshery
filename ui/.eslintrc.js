module.exports = {
  env: {
    browser: true,
    es6: true,
    node: true, // tells the parser that we are using nodejs
  },

  settings: {
    react: {
      version: require('./package.json').dependencies.react,
    },
  },
  extends: [
    'eslint:recommended',
    'plugin:react/recommended',
    'plugin:cypress/recommended',
    'next',
    'plugin:prettier/recommended',
  ],
  globals: {
    Atomics: 'readonly',
    SharedArrayBuffer: 'readonly',
    globalThis: 'readonly',
  },
  parser: '@babel/eslint-parser',
  parserOptions: {
    ecmaFeatures: {
      jsx: true,
    },
    ecmaVersion: 2018,
    sourceType: 'module',
  },
  plugins: ['react', 'cypress', 'prettier'],
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
    'react/jsx-uses-vars': [2],
    'react/jsx-no-undef': 'error',
    'no-console': 0,
    'no-unused-vars': 'error',
    'react/jsx-key': 'warn',
    'no-dupe-keys': 'error',
    'react/jsx-filename-extension': [
      1,
      {
        extensions: ['.js', '.jsx'],
      },
    ],
    'react/prop-types': 'off',
    'prettier/prettier': ['error', { endOfLine: 'lf' }],
  },
};
