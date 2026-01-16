module.exports = {
  // Lint & Prettify TS and JS files
  '**/*.{ts,tsx,js,jsx}': ['npm run lint:fix', 'npm run format'],

  // Ignore generated/large model artifacts
  ignore: ['../server/meshmodels/**'],
};
