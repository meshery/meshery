import React from 'react';
import { useTheme } from '@material-ui/core/styles';
import { SistentThemeProviderWithoutBaseLine } from '@layer5/sistent';

export const UsesSistent = ({ children }) => {
  const theme = useTheme();
  const mode = theme.palette.type;
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
