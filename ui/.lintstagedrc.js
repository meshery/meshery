module.exports = {
  // Lint & Prettify TS and JS files
  '**/*.(ts|tsx|js|jsx)': (filenames) => [`npm run lint`, `npm run format`],

  // Prettify only Markdown and JSON files
  '**/*.(md|json)': (filenames) => `npx prettier --write ${filenames.join(' ')}`,
};
