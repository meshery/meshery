import { useEffect, useState } from 'react';
import { useGetUserPrefQuery, useUpdateUserPrefWithContextMutation } from '@/rtk-query/user';
import _ from 'lodash/fp';
import ProviderStoreWrapper from '@/store/ProviderStoreWrapper';
import { useMediaQuery } from '@sistent/sistent';

const THEME_STORAGE_KEY = 'meshery-theme';

type ThemeMode = 'light' | 'dark';

const isThemeMode = (value: string | null): value is ThemeMode =>
  value === 'light' || value === 'dark';

const getStoredTheme = (): ThemeMode | null => {
  try {
    const value = localStorage.getItem(THEME_STORAGE_KEY);
    return isThemeMode(value) ? value : null;
  } catch {
    return null;
  }
};

const setStoredTheme = (theme: ThemeMode): void => {
  try {
    localStorage.setItem(THEME_STORAGE_KEY, theme);
  } catch {
    // Ignore storage errors (e.g., storage disabled/unavailable).
  }
};

// Default to dark on the server and first client render, matching the
// pre-hydration UI; resolves to the real system preference after mount.
export const useGetSystemTheme = (): ThemeMode => {
  const prefersDark = useMediaQuery('(prefers-color-scheme: dark)', { defaultMatches: true });
  return prefersDark ? 'dark' : 'light';
};

export const useThemePreference = () => {
  const { data, isLoading, ...res } = useGetUserPrefQuery();
  const systemPref = useGetSystemTheme();
  const [storedMode, setStoredMode] = useState<ThemeMode | null>(getStoredTheme);

  useEffect(() => {
    const handler = () => setStoredMode(getStoredTheme());
    window.addEventListener('theme-change', handler);
    return () => window.removeEventListener('theme-change', handler);
  }, []);

  const remoteThemeValue = data?.remoteProviderPreferences?.theme ?? null;
  const remoteMode = isThemeMode(remoteThemeValue) ? remoteThemeValue : null;
  const mode = isLoading ? storedMode || systemPref : remoteMode || storedMode || systemPref;

  useEffect(() => {
    const remoteTheme = data?.remoteProviderPreferences?.theme;
    if (!isLoading && isThemeMode(remoteTheme)) {
      setStoredTheme(remoteTheme);
      setStoredMode(remoteTheme);
    }
  }, [isLoading, data?.remoteProviderPreferences?.theme]);

  return {
    data: { mode },
    isLoading,
    setStoredMode,
    ...res,
  };
};

const ThemeTogglerCore_ = ({ Component }) => {
  const themePref = useThemePreference();
  const [handleUpdateUserPref] = useUpdateUserPrefWithContextMutation();
  const { data: userPrefs, isLoading } = useGetUserPrefQuery();

  const mode = themePref?.data?.mode;
  const { setStoredMode } = themePref;

  const toggleTheme = () => {
    if (isLoading) return;
    const newTheme = mode === 'light' ? 'dark' : 'light';
    setStoredTheme(newTheme);
    setStoredMode(newTheme);
    window.dispatchEvent(new Event('theme-change'));

    const isRemoteProvider = !!userPrefs?.remoteProviderPreferences;
    if (isRemoteProvider) {
      const updated = _.set('remoteProviderPreferences.theme', newTheme, userPrefs);
      handleUpdateUserPref({ body: updated });
    }
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
