export const CamelCaseToSentanceCase = (stringValue: string): string => {
  return stringValue.replace(/([A-Z]+)/g, ' $1').replace(/([A-Z][a-z])/g, ' $1');
};
