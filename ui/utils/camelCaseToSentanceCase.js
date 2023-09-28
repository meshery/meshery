export const CamelCaseToSentanceCase = (stringValue) => {
  return stringValue.replace(/([A-Z]+)/g, ' $1').replace(/([A-Z][a-z])/g, ' $1');
};
