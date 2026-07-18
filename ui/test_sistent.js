const sistent = require('@sistent/sistent');
const exportsToCheck = [
  'Checkbox',
  'Box',
  'CustomTooltip',
  'TextField',
  'ClickAwayListener',
  'IconButton',
  'Slide',
  'Grid2',
  'Hidden',
  'NoSsr',
  'useTheme',
  'useMediaQuery',
  'SearchIcon',
  'SettingsIcon',
  'FilterAllIcon',
];

exportsToCheck.forEach((exp) => {
  if (sistent[exp] === undefined) {
    console.log('MISSING EXPORT: ' + exp);
  } else {
    console.log('FOUND: ' + exp);
  }
});
