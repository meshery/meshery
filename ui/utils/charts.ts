export const isValidColumnName = (name) => {
  return name !== '' && name !== ' ' && name != undefined && name != null;
};

export const CHART_COLORS = [
  '#14232A', // Gunmetal
  // '#213A45',
  '#2E5261',
  // '#294957',
  '#3B697D', // Paynes' Gray
  // '#396679',
  '#4A839C',
  // '#477E96', Teal Blue
  '#5996B1',
  // '#639CB5',
  '#74A8BE',
  // '#8Bb2C6',
  '#90B9CB',
  // '#AACCBD8', // Columbia Blue
  '#CBDEE6',
  // '#EEF4F7'
];

export const dataToColors = (data) => {
  const columns = data.map((item) => item[0]);
  const colors = {};
  let colorIdx = 0;

  columns.forEach((col) => {
    if (colorIdx >= CHART_COLORS.length) {
      colorIdx = 0;
    }
    colors[col] = CHART_COLORS[colorIdx];
    colorIdx += 1;
  });

  return colors;
};
