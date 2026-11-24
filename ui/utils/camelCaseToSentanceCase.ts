export const CamelCaseToSentanceCase = (stringValue) => {
  if (typeof stringValue !== 'string' || !stringValue) {
    return '';
  }
  return stringValue.replace(/([A-Z]+)/g, ' $1').replace(/([A-Z][a-z])/g, ' $1');
};
