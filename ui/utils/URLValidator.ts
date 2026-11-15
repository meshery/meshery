export const URLValidator = (url?: string | null): boolean => {
  if (!url) return false;

  // Pattern covering general scheme, host/ip, port, and path/query/fragment
  const compulsoryProtocolValidUrlPattern = new RegExp(
    '^(http|https|nats|tcp):\\/\\/' + // compulsory protocol
      '(' +
      '([a-z\\d]([a-z\\d-]{0,61}[a-z\\d])?\\.)*[a-z\\d]([a-z\\d-]{0,61}[a-z\\d])?|' + // domain name
      'localhost|' +
      '(\\d{1,3}\\.){3}\\d{1,3}' + // IPv4
      ')' +
      '(:\\d+)?' + // port (optional)
      '(?:[/?#][^\\s]*)?$', // path, query, fragment
    'i',
  );

  if (!compulsoryProtocolValidUrlPattern.test(url)) {
    return false;
  }

  try {
    // Replace custom protocols like nats:// or tcp:// with http:// for native URL parsing
    const parsedUrl = new URL(url.replace(/^(nats|tcp):/i, 'http:'));

    // Guard against port 0 (native URL parser accepts port 0, but network ports must be >= 1)
    if (parsedUrl.port) {
      const portNum = Number(parsedUrl.port);
      if (portNum < 1) return false;
    }
  } catch {
    return false;
  }

  return true;
};
