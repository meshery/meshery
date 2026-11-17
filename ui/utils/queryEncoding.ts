export const urlEncodeArrayParam = (key: string, array: unknown) => {
  if (!array) {
    return '';
  }
  if (typeof array === 'string') {
    return `${encodeURIComponent(key)}=${encodeURIComponent(array)}`;
  }
  if (!Array.isArray(array) || array.length === 0) {
    return '';
  }
  return array
    .map((item) => `${encodeURIComponent(key)}=${encodeURIComponent(String(item))}`)
    .join('&');
};

export const urlEncodeParams = (params) => {
  const urlSearchParams = new URLSearchParams();

  Object.entries(params).forEach(([key, value]) => {
    if (value === null || value === undefined) {
      return;
    }
    if (Array.isArray(value)) {
      value.forEach((val) => urlSearchParams.append(key, val));
    } else {
      urlSearchParams.append(key, value);
    }
  });

  return urlSearchParams.toString();
};
