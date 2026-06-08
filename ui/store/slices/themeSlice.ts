import { createSlice, PayloadAction } from '@reduxjs/toolkit';

type ThemeMode = 'light' | 'dark';

const THEME_STORAGE_KEY = 'meshery-theme';

const getInitialMode = (): ThemeMode | null => {
  if (typeof window === 'undefined') return null;
  try {
    const saved = localStorage.getItem(THEME_STORAGE_KEY);
    if (saved === 'light' || saved === 'dark') return saved;
  } catch (e) {
    console.error('Failed to read theme from localStorage', e);
  }
  return null;
};

const themeSlice = createSlice({
  name: 'theme',
  initialState: {
    mode: getInitialMode(),
  } as { mode: ThemeMode | null },
  reducers: {
    setThemeMode: (state, action: PayloadAction<ThemeMode>) => {
      state.mode = action.payload;
      if (typeof window !== 'undefined') {
        try {
          localStorage.setItem(THEME_STORAGE_KEY, action.payload);
        } catch (e) {
          console.error('Failed to write theme to localStorage', e);
        }
      }
    },
  },
});

export const { setThemeMode } = themeSlice.actions;

export const selectThemeMode = (state: { theme: { mode: ThemeMode | null } }) => state.theme.mode;

export default themeSlice.reducer;
