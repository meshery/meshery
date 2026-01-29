const ignoreList = ['.eslintrc.js', '.lintstagedrc.js', 'package.json', 'package-lock.json'];

module.exports = {
  '**/*.(ts|tsx|js|jsx)': (filenames) => {
    const filtered = filenames.filter((f) => !ignoreList.includes(require('path').basename(f)));
    return filtered.length === 0
      ? []
      : [
          `npx eslint --fix --max-warnings=0 ${filtered.map((f) => `"${f}"`).join(' ')}`,
          `npx biome format --write ${filtered.map((f) => `"${f}"`).join(' ')}`,
        ];
  },
  '**/*.(md|json)': (filenames) =>
    `npx biome format --write ${filenames.map((f) => `"${f}"`).join(' ')}`,
};
