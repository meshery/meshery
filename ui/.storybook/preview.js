import React, { useMemo } from "react";
import { CssBaseline, ThemeProvider } from "@mui/material";
import { lightTheme, darkTheme } from "../stories/themes/app";
import { withThemeFromJSXProvider } from "@storybook/addon-styling";

export const parameters = {
  actions: { argTypesRegex: "^on[A-Z].*" },
  controls: {
    matchers: {
      color: /(background|color)$/i,
      date: /Date$/,
    },
  },
}

export const decorators = [
  withThemeFromJSXProvider({
    themes: {
      light: lightTheme,
      dark: darkTheme
    },
    defaultTheme: "light",
    Provider: ThemeProvider,
    GlobalStyles: CssBaseline
  }),
];

export const globalTypes = {
  theme: {
    name: "Theme",
    title: "Theme",
    description: "Theme for your components",
    defaultValue: "light",
    toolbar: {
      icon: "paintbrush",
      dynamicTitle: true,
      items: [
        { value: "light", left: "☀️", title: "Light mode" },
        { value: "dark", left: "🌙", title: "Dark mode" },
      ],
    },
  },
};