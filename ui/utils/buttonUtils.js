export const getButtonStateStyles = (condition) => ({
  opacity: condition ? 0.5 : 1,
  pointerEvents: condition ? 'none' : 'auto',
  cursor: condition ? 'not-allowed' : 'pointer',
});
