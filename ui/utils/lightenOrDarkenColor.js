
export function pSBCr(color, percent) {
  if (!color) {
    return "#607d8b";
  }

  // handle color in Hex Format
  if (color.startsWith("#")) {
    let R = parseInt(color?.substring(1, 3), 16);
    let G = parseInt(color?.substring(3, 5), 16);
    let B = parseInt(color?.substring(5, 7), 16);

    return RGBBlender(R, G, B, percent);
  }

  // handle color in RGB Format
  const RGB = color.substring(3).replaceAll("(", "").replaceAll(")", "").split(",");
  if (RGB.length >= 3) {
    // rgb with opacity may have length == 4
    return RGBBlender(RGB[0], RGB[1], RGB[2], percent);
  }
}

function RGBBlender(R, G, B, percent) {
  R = parseInt((R * (100 + percent)) / 100);
  G = parseInt((G * (100 + percent)) / 100);
  B = parseInt((B * (100 + percent)) / 100);

  R = R < 255 ? R : 255;
  G = G < 255 ? G : 255;
  B = B < 255 ? B : 255;

  let RR = R.toString(16).length == 1 ? "0" + R.toString(16) : R.toString(16);
  let GG = G.toString(16).length == 1 ? "0" + G.toString(16) : G.toString(16);
  let BB = B.toString(16).length == 1 ? "0" + B.toString(16) : B.toString(16);

  return "#" + RR + GG + BB;
}
