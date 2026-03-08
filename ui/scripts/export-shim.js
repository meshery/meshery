const fs = require('fs');

if (!fs.existsSync('out')) {
  console.error('Static export output not found. Run npm run build first.');
  process.exit(1);
}

console.log("Static export is handled by next build via output='export'.");
