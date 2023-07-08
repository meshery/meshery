export const CHART_COLORS = ['#00B39F' ,'#396679' ,'#83B71E' , '#00382D','#477E96','#F0A303', '#396679']

export const dataToColors = (data) => {
  const columns = data.map((item) => item[0] )
  const colors = { }
  let colorIdx = 0

  columns.forEach((col) => {
    if (colorIdx >= CHART_COLORS.length) {
      colorIdx = 0
    }
    colors[col] = CHART_COLORS[colorIdx]
    colorIdx +=1
  })

  return colors
}
