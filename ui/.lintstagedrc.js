module.exports = {
  // Lint & Prettify TS and JS files
  '**/*.{ts,tsx,js,jsx}': ['npm run lint:fix', 'npm run format'],

  // Prettify only Markdown and JSON files, excluding generated/large model artifacts.
  // Note: `ignore` is not a recognized lint-staged config key (it maps globs to
  // commands only), so filtering must happen inside the task function itself.
  '**/*.{md,json}': (filenames) => {
    const filtered = filenames.filter((f) => !f.includes('/server/meshmodels/'));
    return filtered.length ? `npx prettier --write ${filtered.join(' ')}` : [];
  },
};
