// TODO: import this from Sistent
const getUnit8ArrayForDesign = (design) => {
  // converts a string to a Uint8Array
  const uint8Array = Uint8Array.from(design, (char) => char.charCodeAt(0));

  return Array.from(uint8Array);
};


const pattern_file = getUnit8ArrayForDesign(process.env.pattern_file);
console.dir(pattern_file, {'maxArrayLength': null}); // necessary to override the default maxArrayLength of 100


