import React from 'react';
import { SistentThemeProviderWithoutBaseLine, useTheme } from '@layer5/sistent';

export const UsesSistent = ({ children }) => {
  const theme = useTheme();
  const mode = theme.palette.mode;
  if (mode === 'dark') {
    return (
      <SistentThemeProviderWithoutBaseLine initialMode={mode}>
        {children}
      </SistentThemeProviderWithoutBaseLine>
    );
  }
  return (
    <SistentThemeProviderWithoutBaseLine initialMode={mode}>
      {children}
    </SistentThemeProviderWithoutBaseLine>
  );
};
