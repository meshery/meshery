import React, { useEffect, useState } from 'react';
import { useGetUserPrefQuery, useUpdateUserPrefWithContextMutation } from '@/rtk-query/user';
import _ from 'lodash/fp';
import ProviderStoreWrapper from '@/store/ProviderStoreWrapper';

const THEME_STORAGE_KEY = 'meshery-theme';

const getStoredTheme = (): string | null => {
  try {
    return localStorage.getItem(THEME_STORAGE_KEY);
  } catch {
    return null;
  }
};

const setStoredTheme = (theme: string): void => {
  try {
    localStorage.setItem(THEME_STORAGE_KEY, theme);
  } catch {
    // Ignore storage errors (e.g., storage disabled/unavailable).
  }
};

export const useGetSystemTheme = () => {
  const [theme, setTheme] = React.useState('dark');
  useEffect(() => {
    const systemPref = window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
    setTheme(systemPref);
  }, []);
  return theme;
};

export const useThemePreference = () => {
  const { data, isLoading, ...res } = useGetUserPrefQuery();
  const systemPref = useGetSystemTheme();
  const [storedMode, setStoredMode] = useState<string | null>(getStoredTheme);

  useEffect(() => {
    const handler = () => setStoredMode(getStoredTheme());
    window.addEventListener('theme-change', handler);
    return () => window.removeEventListener('theme-change', handler);
  }, []);

  const mode = isLoading
    ? storedMode || systemPref || 'dark'
    : data?.remoteProviderPreferences?.theme || storedMode || systemPref || 'dark';

  useEffect(() => {
    if (!isLoading && data?.remoteProviderPreferences?.theme) {
      setStoredTheme(data.remoteProviderPreferences.theme);
      setStoredMode(data.remoteProviderPreferences.theme);
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
  const { data: userPrefs } = useGetUserPrefQuery();

  const mode = themePref?.data?.mode;
  const { setStoredMode } = themePref;

  const toggleTheme = () => {
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

  return (
    <ProviderStoreWrapper>
      <Component mode={mode} toggleTheme={toggleTheme} />
    </ProviderStoreWrapper>
  );
};

export const ThemeTogglerCore = (props) => {
  return (
    <ProviderStoreWrapper>
      <ThemeTogglerCore_ {...props} />
    </ProviderStoreWrapper>
  );
};
