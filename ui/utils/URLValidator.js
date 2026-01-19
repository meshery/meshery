export const URLValidator = (url) => {
  const compulsoryProtocolValidUrlPattern = new RegExp(
    '(^(http|https|nats|tcp):\\/\\/)' + // compulsory protocol
      '((([a-z\\d]([a-z\\d-]*[a-z\\d])*)\\.{0,})+[a-z]{0,}|' + // domain name
      'localhost|' +
      '((\\d{1,3}.){3}\\d{1,3}))' + // OR ip (v4) address
      '(\\:\\d+)?(/[-a-z\\d%_.~+]*)*' + // port and path
      '(\\?[;&a-z\\d%_.~+=-]*)?' + // query string
      '(\\#[-a-z\\d_]*)?$',
    'i',
  ); // fragment locator

  return url?.match(compulsoryProtocolValidUrlPattern);
};
