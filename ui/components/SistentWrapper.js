import React from 'react';
import { SistentThemeProvider } from '@layer5/sistent';
import { useTheme } from '@material-ui/core/styles';

export const UsesSistent = ({ children }) => {
  const theme = useTheme();
  const mode = theme.palette.type;
  if (mode === 'dark') {
    return <SistentThemeProvider initialMode={mode}> {children} </SistentThemeProvider>;
  }
  return <SistentThemeProvider initialMode={mode}> {children} </SistentThemeProvider>;
};
