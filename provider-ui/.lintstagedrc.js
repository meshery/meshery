module.exports = {
  '**/*.{js,ts,jsx,tsx}': (filenames) =>
    `eslint --fix --cache --max-warnings=0 ${filenames.map((f) => `"${f}"`).join(' ')}`,
};
