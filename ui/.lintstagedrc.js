module.exports = {
  // Lint & Prettify only staged JS/JSX/TS/TSX files
  '**/*.(ts|tsx|js|jsx)': (filenames) => [
    `npx eslint --fix --max-warnings=0 ${filenames.join(' ')}`,
    `npx prettier --write ${filenames.join(' ')}`,
  ],

  // Prettify only Markdown and JSON files
  '**/*.(md|json)': (filenames) => `npx prettier --write ${filenames.join(' ')}`,
};