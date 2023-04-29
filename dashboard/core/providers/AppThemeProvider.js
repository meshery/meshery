import { useDispatch, useSelector } from "react-redux";
import createEmotionCache from "../createEmotionCache";
import IconButton from "@mui/material/IconButton";
import Head from "next/head";
import { Brightness2Outlined, WbSunnyOutlined } from "@mui/icons-material";
import { CacheProvider } from "@emotion/react";
import { CssBaseline } from "@mui/material";
import { ThemeProvider } from "@mui/material/styles";
import Paper from "@mui/material/Paper";
import { darkTheme, lightTheme } from "../styles/themes";
import { toggleTheme } from "../store/theme/themeSlice";

const clientSideEmotionCache = createEmotionCache();

export function AppThemeProvider({ children, emotionCache = clientSideEmotionCache }) {
    const theme = useSelector((state) => state.theme);

    const dispatch = useDispatch();

    const ToggleSwitch = () => {
        return (
            <div
                style={{
                    position: "absolute",
                    top: "10px",
                    right: "10px",
                }}
            >
                <IconButton
                    onClick={() => dispatch(toggleTheme())}
                    aria-label="toggle theme"
                >
                    {theme.darkTheme ? <WbSunnyOutlined /> : <Brightness2Outlined />}
                </IconButton>
            </div>
        );
    };

    return (
        <CacheProvider value={emotionCache}>
        <Head>
            <meta name="viewport" content="initial-scale=1, width=device-width" />
        </Head>
        <ThemeProvider theme={theme.darkTheme ? darkTheme : lightTheme}>
            <CssBaseline />
            <Paper
                style={{
                    minHeight: "100vh",
                    borderRadius: "0",
                    display: "flex",
                    flexDirection: "column",
                    alignItems: "center",
                    justifyContent: "center",
                }}
            >
                <ToggleSwitch />
                {children}
            </Paper>
        </ThemeProvider>
    </CacheProvider>
    )
}