const { ESLint } = require('eslint');

const quoteFilenames = (filenames) => filenames.map((filename) => `"${filename}"`).join(' ');

async function getLintCommands(filenames) {
  const eslint = new ESLint();
  const filesToLint = [];

  for (const filename of filenames) {
    if (!(await eslint.isPathIgnored(filename))) {
      filesToLint.push(filename);
    }
  }

  const commands = [];

  if (filesToLint.length > 0) {
    commands.push(`npx eslint --fix --max-warnings=0 ${quoteFilenames(filesToLint)}`);
  }

  commands.push(`npx prettier --write ${quoteFilenames(filenames)}`);

  return commands;
}

module.exports = {
  // Lint & Prettify only staged JS/JSX/TS/TSX files
  '**/*.(ts|tsx|js|jsx)': getLintCommands,

  // Prettify only Markdown and JSON files
  '**/*.(md|json)': (filenames) => `npx prettier --write ${filenames.join(' ')}`,
};
