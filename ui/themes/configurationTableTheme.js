import { createTheme } from '@mui/material/styles';

function configurationTableTheme(style = {}) {
  const theme = createTheme({
    shadows: ['none'],
    components: {
      MuiInput: {
        styleOverrides: {
          underline: {
            '&:hover:not(.Mui-disabled):before': {
              borderBottom: '2px solid #222',
            },
            '&:after': {
              borderBottom: '2px solid #222',
            },
          },
        },
      },
      MUIDataTableSearch: {
        styleOverrides: {
          searchIcon: {
            color: '#607d8b',
            marginTop: '7px',
            marginRight: '8px',
          },
          clearIcon: {
            '&:hover': {
              color: '#607d8b',
            },
          },
        },
      },
      MUIDataTableSelectCell: {
        styleOverrides: {
          checkboxRoot: {
            '&$checked': {
              color: '#607d8b',
            },
          },
        },
      },
      MUIDataTableToolbar: {
        styleOverrides: {
          iconActive: {
            color: '#222',
          },
          icon: {
            '&:hover': {
              color: '#607d8b',
            },
          },
        },
      },
    },
    ...style,
  });
  return theme;
}

function configurationTableThemeDark(style = {}) {
  const theme = createTheme({
    shadows: ['none'],
    palette: {
      mode: 'dark', // Use mode instead of type for dark theme
    },
    components: {
      MuiPaper: {
        styleOverrides: {
          root: {
            backgroundColor: '#363636',
          },
        },
      },
      MuiFormLabel: {
        styleOverrides: {
          root: {
            '&$focused': {
              color: '#00B39F',
            },
          },
        },
      },
    },
    ...style,
  });
  return theme;
}

export { configurationTableTheme, configurationTableThemeDark };
