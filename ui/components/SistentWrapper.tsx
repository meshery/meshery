import type { ReactNode } from 'react';
import { SistentThemeProviderWithoutBaseLine, useTheme } from '@sistent/sistent';

interface UsesSistentProps {
  children: ReactNode;
}

export const UsesSistent = ({ children }: UsesSistentProps) => {
  const theme = useTheme();
  const mode = theme.palette.mode;
  return (
    <SistentThemeProviderWithoutBaseLine initialMode={mode}>
      {children}
    </SistentThemeProviderWithoutBaseLine>
  );
};
