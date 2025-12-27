module.exports = {
  // Lint & Prettify TS and JS files - only the staged files
  '**/*.(ts|tsx|js|jsx)': (filenames) => [
    `npx eslint --fix --max-warnings=0 ${filenames.map((f) => `'${f.replace(/'/g, "'\\''" )}'`).join(' ')}`,
    `npx prettier --write ${filenames.map((f) => `'${f.replace(/'/g, "'\\''" )}'`).join(' ')}`,
  ],

  // Prettify only Markdown and JSON files
  '**/*.(md|json)': (filenames) => `npx prettier --write ${filenames.map((f) => `'${f.replace(/'/g, "'\\''" )}'`).join(' ')}`,
};
