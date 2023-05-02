import { createTheme } from "@mui/material";
import { blueGrey } from "@mui/material/colors";

const createCustomTheme = (paletteType) => {
  const commonPalette = {
    primary : {
      main : blueGrey[500],
      dark : blueGrey[700],
      light : blueGrey[400]
    },
    secondary : {
      main : "#00D3A9",
      dark : "#00B39F",
      darker : "#3C494F"
    },
    neutral : {
      main : '#64748B',
      contrastText : '#FFF',
    }
  };

  const palette = paletteType === "dark"
    ? {
      mode : "dark",
      ...commonPalette,
      secondary : {
        ...commonPalette.secondary
      },
    }
    : {
      mode : "light",
      ...commonPalette,
      secondary : {
        ...commonPalette.secondary,
      }
    }

  const theme = createTheme({
    palette,
    typography : {
      h5 : {
        fontWeight : "bolder",
        fontSize : 26,
        letterSpacing : 0.5
      },
    },
    shape : { borderRadius : 8 },
    breakpoints : {
      values : {
        xs : 0,
        sm : 600,
        md : 960,
        lg : 1280,
        xl : 1920,
      },
    },
    components : {
      MuiAppBar : {
        styleOverrides : {
          colorPrimary : {
            backgroundColor : blueGrey[500],
          },
          colorSecondary : {
            backgroundColor : blueGrey[600],
          }
        }
      }
    }
  })

  const overrides = {
    //
  }

  theme.overrides = {
    ...overrides,
    MuiCssBaseline : {
      '@global' : {
        body : {
          backgroundColor : paletteType === 'dark' ? '#303030' : '#eaeff1',
          color : paletteType === 'dark' ? '#000000' : '#FFFFFF'
        },
        h5 : {
          color : paletteType === 'dark' ? '#000000' : '#FFFFFF'
        },
        p : {
          color : paletteType === 'dark' ? '#000000' : '#FFFFFF'
        }
      }
    }
  }

  return theme
}

const darkTheme = createCustomTheme('dark');
const lightTheme = createCustomTheme('light');

export { darkTheme, lightTheme };