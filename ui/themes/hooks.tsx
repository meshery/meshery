import React, { useEffect } from 'react';
import { useGetUserPrefQuery, useUpdateUserPrefWithContextMutation } from '@/rtk-query/user';
import _ from 'lodash/fp';
import ProviderStoreWrapper from '@/store/ProviderStoreWrapper';
import { useDispatch, useSelector } from 'react-redux';
import { selectThemeMode, setThemeMode } from '@/store/slices/themeSlice';

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
  const storedMode = useSelector(selectThemeMode);

  // Priority order:
  // 1. localStorage/RTK store (if user has explicitly set a preference)
  // 2. Remote provider server preference
  // 3. OS system preference
  // 4. Fallback to dark
  const mode = storedMode || data?.remoteProviderPreferences?.theme || systemPref || 'dark';

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
  const { data: userPrefs } = useGetUserPrefQuery();
  const dispatch = useDispatch();

  const mode = themePref?.data?.mode;

  const toggleTheme = () => {
    const newTheme = mode === 'light' ? 'dark' : 'light';

    // 1. Persist to RTK store + localStorage (survives refresh, works for all providers)
    dispatch(setThemeMode(newTheme));

    // 2. Also save to server for remote provider users
    const updated = _.set('remoteProviderPreferences.theme', newTheme, userPrefs);
    handleUpdateUserPref({
      body: updated,
    });
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
