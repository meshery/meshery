export function pSBCr(color: string | undefined | null, percent: number): string {
  if (!color) {
    return '#607d8b';
  }

  // handle color in Hex Format
  if (color.startsWith('#')) {
    const R = parseInt(color?.substring(1, 3), 16);
    const G = parseInt(color?.substring(3, 5), 16);
    const B = parseInt(color?.substring(5, 7), 16);

    return RGBBlender(R, G, B, percent);
  }

  // handle color in RGB Format
  const RGB = color.substring(3).replaceAll('(', '').replaceAll(')', '').split(',');
  if (RGB.length >= 3) {
    // rgb with opacity may have length == 4
    return RGBBlender(parseInt(RGB[0]), parseInt(RGB[1]), parseInt(RGB[2]), percent);
  }

  return '#607d8b';
}

function RGBBlender(R: number, G: number, B: number, percent: number): string {
  let newR = Math.floor((R * (100 + percent)) / 100);
  let newG = Math.floor((G * (100 + percent)) / 100);
  let newB = Math.floor((B * (100 + percent)) / 100);

  newR = newR < 255 ? newR : 255;
  newG = newG < 255 ? newG : 255;
  newB = newB < 255 ? newB : 255;

  const RR = newR.toString(16).padStart(2, '0');
  const GG = newG.toString(16).padStart(2, '0');
  const BB = newB.toString(16).padStart(2, '0');

  return '#' + RR + GG + BB;
}
