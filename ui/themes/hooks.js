import React, { useEffect } from 'react';
import { useGetUserPrefQuery, useUpdateUserPrefWithContextMutation } from '@/rtk-query/user';
import { useState } from 'react';
import _ from 'lodash/fp';

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
  const mode = data?.remoteProviderPreferences?.theme || systemPref || 'dark';

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
  const { data: userPrefs } = useGetUserPrefQuery();
  const [mode, setMode] = useState(themePref?.data?.mode);

  useEffect(() => {
    setMode(themePref?.data?.mode);
  }, [themePref?.data?.mode]);

  const toggleTheme = () => {
    const newTheme = mode === 'light' ? 'dark' : 'light';
    setMode(newTheme);
    const updated = _.set('remoteProviderPreferences.theme', newTheme, userPrefs);
    console.log('updated', updated);
    handleUpdateUserPref({
      body: updated,
    });
  };

  return <Component mode={mode} toggleTheme={toggleTheme} />;
};
