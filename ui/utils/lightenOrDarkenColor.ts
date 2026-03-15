export function pSBCr(color: string | undefined, percent: number): string {
  if (!color) {
    return '#607d8b';
  }

  // handle color in Hex Format
  if (color.startsWith('#')) {
    const R = parseInt(color.substring(1, 3), 16);
    const G = parseInt(color.substring(3, 5), 16);
    const B = parseInt(color.substring(5, 7), 16);

    return RGBBlender(R, G, B, percent);
  }

  // handle color in RGB Format
  const RGB = color.substring(3).replaceAll('(', '').replaceAll(')', '').split(',');
  if (RGB.length >= 3) {
    // rgb with opacity may have length == 4
    return RGBBlender(Number(RGB[0]), Number(RGB[1]), Number(RGB[2]), percent);
  }

  return '#607d8b';
}

function RGBBlender(R: number, G: number, B: number, percent: number): string {
  R = Math.round((R * (100 + percent)) / 100);
  G = Math.round((G * (100 + percent)) / 100);
  B = Math.round((B * (100 + percent)) / 100);

  R = R < 255 ? R : 255;
  G = G < 255 ? G : 255;
  B = B < 255 ? B : 255;

  let RR = R.toString(16).length == 1 ? '0' + R.toString(16) : R.toString(16);
  let GG = G.toString(16).length == 1 ? '0' + G.toString(16) : G.toString(16);
  let BB = B.toString(16).length == 1 ? '0' + B.toString(16) : B.toString(16);

  return '#' + RR + GG + BB;
}
