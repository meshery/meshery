import { createTheme } from "@mui/material";
import { black, charcoal, colorText, keppel, white, error, warning, success, info, neutral } from "./colors";

export const lightTheme = createTheme({
    palette: {
        mode: "light",
        common: {
            black: black,
            white: white
        },
        primary: {
            main: charcoal.main,
            dark: charcoal.dark,
            light: charcoal.light
        },
        secondary: {
            main: keppel.main,
            dark: keppel.dark,
            light: keppel.light,
            contrastText: white,
        },
        error: {
            main: error.main
        },
        warning: {
            main: warning.main,
            contrastText: white,
        },
        success: {
            main: success.main,
        },
        info: {
            main: info.main,
        }
    }
});

export const darkTheme = createTheme({
    palette: {
        mode: "dark",
        common: {
            black: black,
            white: white,
        },
        primary: {
            main: keppel.main,
            dark: keppel.dark,
            light: keppel.light,
            contrastText: white,
        },
        secondary: {
            main: charcoal.main,
            dark: charcoal.dark,
            light: charcoal.light,
            contrastText: white,
        },
        text: {
            primary: colorText.primary,
            secondary: colorText.secondary,
            disabled: colorText.disabled
        },
        background: {
            default: neutral.main,
            paper: neutral.main
        },
        error: {
            main: error.accent,
        },
        warning: {
            main: warning.accent,
            contrastText: white,
        },
        success: {
            main: success.accent,
        },
        info: {
            main: info.accent,
        }
    }
});