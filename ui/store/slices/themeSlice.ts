import { createSlice, PayloadAction } from '@reduxjs/toolkit';

type ThemeMode = 'light' | 'dark';

const THEME_STORAGE_KEY = 'meshery-theme';

// Read from localStorage on app start
// If nothing saved yet, fall back to 'dark' (Meshery's default)
const getInitialMode = (): ThemeMode => {
  if (typeof window === 'undefined') return 'dark'; // SSR safety for Next.js
  const saved = localStorage.getItem(THEME_STORAGE_KEY);
  if (saved === 'light' || saved === 'dark') return saved;
  return 'dark';
};

const themeSlice = createSlice({
  name: 'theme',
  initialState: {
    mode: getInitialMode(),
  },
  reducers: {
    setThemeMode: (state, action: PayloadAction<ThemeMode>) => {
      state.mode = action.payload;
      // Persist to localStorage every time theme changes
      if (typeof window !== 'undefined') {
        localStorage.setItem(THEME_STORAGE_KEY, action.payload);
      }
    },
    toggleThemeMode: (state) => {
      const newMode = state.mode === 'dark' ? 'light' : 'dark';
      state.mode = newMode;
      if (typeof window !== 'undefined') {
        localStorage.setItem(THEME_STORAGE_KEY, newMode);
      }
    },
  },
});

export const { setThemeMode, toggleThemeMode } = themeSlice.actions;

// Selector — any component can read theme mode from the store
export const selectThemeMode = (state: { theme: { mode: ThemeMode } }) => state.theme.mode;

export default themeSlice.reducer;
