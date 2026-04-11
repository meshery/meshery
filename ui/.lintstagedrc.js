module.exports = {
  // Lint & Prettify only staged JS/JSX/TS/TSX files
  '**/*.(ts|tsx|js|jsx)': (filenames) => [
    `eslint --fix --max-warnings=0 --no-warn-ignored ${filenames.map((f) => `"${f}"`).join(' ')}`,
    `npx prettier --write ${filenames.map((f) => `"${f}"`).join(' ')}`,
  ],

  // Prettify only Markdown and JSON files
  '**/*.(md|json)': (filenames) => `npx prettier --write ${filenames.join(' ')}`,
};
