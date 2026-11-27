import { useEffect, useState } from 'react';
import { useGetUserPrefQuery, useUpdateUserPrefWithContextMutation } from '@/rtk-query/user';
import ProviderStoreWrapper from '@/store/ProviderStoreWrapper';
import { useMediaQuery } from '@sistent/sistent';

export const useThemePreference = () => {
  const { data, ...res } = useGetUserPrefQuery();
  // Default to dark on the server and first client render, matching the
  // pre-hydration UI; resolves to the real system preference after mount.
  const prefersDark = useMediaQuery('(prefers-color-scheme: dark)', { defaultMatches: true });
  const mode = data?.remoteProviderPreferences?.theme || (prefersDark ? 'dark' : 'light');

  return {
    data: {
      mode,
    },
    ...res,
  };
};

const ThemeTogglerCore_ = ({ Component }) => {
  const themePref = useThemePreference();
  const [handleUpdateUserPref] = useUpdateUserPrefWithContextMutation();
  const [mode, setMode] = useState(themePref?.data?.mode);

  useEffect(() => {
    setMode(themePref?.data?.mode);
  }, [themePref?.data?.mode]);

  const toggleTheme = () => {
    const newTheme = mode === 'light' ? 'dark' : 'light';
    setMode(newTheme);
    // Send only the changed theme field. The server merges the partial body
    // onto the existing preferences, so we never read - and therefore never
    // persist - optimistic-but-unconfirmed values from an in-flight update.
    handleUpdateUserPref({
      body: { remoteProviderPreferences: { theme: newTheme } },
    });
  };

  return <Component mode={mode} toggleTheme={toggleTheme} />;
};

// The single ProviderStoreWrapper here is what gives ThemeTogglerCore_'s RTK
// Query hooks their store context; ThemeTogglerCore is mounted from extension
// surfaces that render outside the app's Redux provider tree.
export const ThemeTogglerCore = (props) => {
  return (
    <ProviderStoreWrapper>
      <ThemeTogglerCore_ {...props} />
    </ProviderStoreWrapper>
  );
};
