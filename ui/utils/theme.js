export const getTheme = () => {
  const storedTheme = localStorage.getItem('Theme');
  const systemTheme = window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
  const defaultTheme = 'dark';
  if (storedTheme) {
    return storedTheme;
  } else if (systemTheme) {
    return systemTheme;
  } else {
    return defaultTheme;
  }
};

export const setTheme = (theme) => {
  localStorage.setItem('Theme', theme);
};
