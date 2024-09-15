import React, { useEffect } from 'react';
import { useGetUserPrefQuery, useUpdateUserPrefWithContextMutation } from '@/rtk-query/user';
import { useState } from 'react';

export const useGetSystemTheme = () => {
  const [theme, setTheme] = React.useState('dark');
  useEffect(() => {
    const systemPref = window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
    setTheme(systemPref);
  }, []);
  return theme;
};

export const useThemePreference = () => {
  const { data, ...res } = useGetUserPrefQuery();
  const systemPref = useGetSystemTheme();
  const mode = data?.theme || systemPref || 'dark';

  return {
    data: {
      mode,
    },
    ...res,
  };
};

export const ThemeTogglerCore = ({ Component }) => {
  const themePref = useThemePreference();
  const [handleUpdateUserPref] = useUpdateUserPrefWithContextMutation();
  const [mode, setMode] = useState(themePref?.data?.mode);

  useEffect(() => {
    setMode(themePref?.data?.mode);
  }, [themePref?.data?.mode]);

  const toggleTheme = () => {
    const newTheme = mode === 'light' ? 'dark' : 'light';
    setMode(newTheme);
    handleUpdateUserPref('theme', newTheme);
  };

  return <Component mode={mode} toggleTheme={toggleTheme} />;
};
