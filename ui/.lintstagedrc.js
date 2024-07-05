module.exports = {
  // Lint & Prettify TS and JS files
  '*.{js,jsx,ts,tsx}': ['eslint --fix', 'prettier --write'],

  // Prettify only Markdown and JSON files
  '**/*.(md|json)': (filenames) => `npx prettier --write ${filenames.join(' ')}`,
};
