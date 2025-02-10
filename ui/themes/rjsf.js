import { createTheme } from '@mui/material/styles';
import createBreakpoints from '@mui/system/createTheme/createBreakpoints';

const breakpoints = createBreakpoints({});

export const rjsfTheme = createTheme({
  palette: {
    primary: {
      main: '#607d8b',
    },
  },
  typography: {
    fontFamily: 'Qanelas Soft, sans-serif',
    fontSize: 13,
  },
  components: {
    MuiTextField: {
      defaultProps: {
        variant: 'outlined',
        margin: 'dense',
      },
      styleOverrides: {
        root: {
          width: 'calc(100% - 4px)',
        },
      },
    },
    MuiCheckbox: {
      defaultProps: {
        color: 'primary',
      },
      styleOverrides: {
        root: {
          marginLeft: '4px',
          '& > *:nth-child(1)': {
            backgroundColor: '#ffffff',
            width: '1rem',
            height: '1rem',
          },
        },
      },
    },
    MuiMenu: {
      defaultProps: {
        variant: 'outlined',
      },
    },
    MuiOutlinedInput: {
      styleOverrides: {
        root: {
          backgroundColor: '#ffffff',
        },
      },
    },
    MuiButton: {
      styleOverrides: {
        textSecondary: {
          color: '#00b39f',
          '&:hover': { color: '#00b39f' },
        },
      },
    },
    MuiInputLabel: {
      styleOverrides: {
        root: {
          whiteSpace: 'nowrap',
          overflow: 'hidden',
          textOverflow: 'ellipsis',
          maxWidth: '60%',
          height: '100%',
          '&:hover': {
            overflow: 'visible',
          },
          '&.Mui-focused': {
            padding: '0.2rem',
          },
        },
        shrink: {
          maxWidth: '100%',
        },
      },
    },
    MuiFormControl: {
      styleOverrides: {
        root: {
          marginTop: '0.3rem',
        },
      },
    },
    MuiFormControlLabel: {
      styleOverrides: {
        root: {
          textTransform: 'capitalize',
          padding: '3px',
        },
      },
    },
    MuiBox: {
      styleOverrides: {
        root: {
          marginTop: 0,
        },
      },
    },
    MuiDivider: {
      styleOverrides: {
        root: {
          height: '0.5px',
        },
      },
    },
    MuiFormLabel: {
      styleOverrides: {
        root: {
          color: '#333',
          fontSize: '0.8rem',
          textTransform: 'capitalize',
        },
      },
    },
    MuiTypography: {
      styleOverrides: {
        body1: {
          fontSize: '0.8rem',
        },
        h5: {
          textTransform: 'capitalize',
          fontSize: '1.1rem',
          fontWeight: 'bold',
        },
        subtitle2: {
          fontSize: '0.8rem',
          fontStyle: 'italic',
          color: '#565656',
        },
      },
    },
    MuiInputBase: {
      styleOverrides: {
        root: {
          fontSize: '0.8rem',
          minRows: 5,
        },
      },
    },
    MuiTooltip: {
      styleOverrides: {
        tooltip: {
          backgroundColor: '#3C494F',
          color: '#fff',
          opacity: '100%',
          fontSize: '14px',
          borderRadius: '0.9375rem',
          padding: '0.9rem',
          zIndex: 99999999999,
        },
        popper: {
          zIndex: 99999,
        },
      },
    },
    MuiAccordionSummary: {
      styleOverrides: {
        root: {
          backgroundColor: 'rgba(242,242,242)',
          borderBottom: '1px solid rgba(0, 0, 0, .125)',
          marginBottom: -1,
          maxHeight: '1.5rem',
          '&$expanded': {
            minHeight: 56,
          },
        },
        content: {
          justifyContent: 'space-between',
          '&$expanded': {
            margin: '12px 0',
            justifyContent: 'space-between',
          },
        },
      },
    },
    MuiAccordionDetails: {
      styleOverrides: {
        root: {
          padding: 16,
        },
      },
    },
    MuiAccordion: {
      styleOverrides: {
        root: {
          border: '1px solid rgba(0, 0, 0, .125)',
          boxShadow: 'none',
          '&:not(:last-child)': {
            borderBottom: 0,
          },
          '&:before': {
            display: 'none',
          },
          '&$expanded': {
            margin: 'auto',
          },
          backgroundColor: 'rgba(242,242,242)',
        },
      },
    },
    MuiGrid: {
      styleOverrides: {
        root: {
          '& > *': {
            border: 'none !important',
          },
          marginTop: '0.1rem !important',
          overflow: 'hidden',
          alignSelf: 'center',
          textOverflow: 'ellipsis',
          '&:hover': {
            overflow: 'visible',
          },
          [breakpoints.up('lg')]: {
            '& > *:nth-child(2)': {
              '& > *:nth-child(1)': {
                '& > *:nth-child(2)': {
                  justifyContent: 'space-around',
                },
              },
            },
          },
        },
      },
    },
    MuiIconButton: {
      styleOverrides: {
        root: {
          cursor: 'pointer',
        },
        sizeSmall: {
          padding: '1px',
        },
      },
    },
    MuiPaper: {
      styleOverrides: {
        elevation0: {
          backgroundColor: 'inherit',
        },
        elevation2: {
          boxShadow: 'none',
        },
      },
    },
  },
});

export const darkRjsfTheme = createTheme({
  palette: {
    mode: 'dark',
    primary: {
      main: '#00B39F',
    },
  },
  typography: {
    fontFamily: 'Qanelas Soft, sans-serif',
    fontSize: 13,
  },
  components: {
    // Most components mirror the light theme, with dark mode adjustments
    MuiTextField: {
      defaultProps: {
        variant: 'outlined',
        margin: 'dense',
      },
      styleOverrides: {
        root: {
          width: 'calc(100% - 4px)',
        },
      },
    },
    MuiSelect: {
      defaultProps: {
        variant: 'outlined',
        margin: 'dense',
      },
      styleOverrides: {
        root: {
          width: 'calc(100% - 4px)',
        },
      },
    },
    MuiOutlinedInput: {
      defaultProps: {
        variant: 'outlined',
        margin: 'dense',
      },
    },
    MuiCheckbox: {
      defaultProps: {
        color: 'primary',
      },
      styleOverrides: {
        root: {
          marginLeft: '4px',
          '& > *:nth-child(1)': {
            backgroundColor: '#303030',
            width: '1rem',
            height: '1rem',
          },
        },
      },
    },
    MuiAccordionSummary: {
      styleOverrides: {
        root: {
          backgroundColor: '#303030',
          borderBottom: '1px solid rgba(255, 255, 255, .125)',
        },
      },
    },
    MuiAccordion: {
      styleOverrides: {
        root: {
          border: '1px solid rgba(255, 255, 255, .125)',
          backgroundColor: '#303030',
        },
      },
    },
    MuiFormLabel: {
      styleOverrides: {
        root: {
          color: '#FFF',
          fontSize: '0.8rem',
          textTransform: 'capitalize',
        },
      },
    },
    MuiInputLabel: {
      styleOverrides: {
        root: {
          whiteSpace: 'nowrap',
          overflow: 'hidden',
          textOverflow: 'ellipsis',
          maxWidth: '60%',
          height: '100%',
          '&:hover': {
            overflow: 'visible',
          },
          '&.Mui-focused': {
            padding: '0.2rem',
          },
        },
        shrink: {
          maxWidth: '100%',
        },
      },
    },
    MuiFormControl: {
      styleOverrides: {
        root: {
          marginTop: '0.3rem',
        },
      },
    },
    MuiFormControlLabel: {
      styleOverrides: {
        root: {
          textTransform: 'capitalize',
          padding: '3px',
        },
      },
    },
    MuiBox: {
      styleOverrides: {
        root: {
          marginTop: 0,
        },
      },
    },
    MuiDivider: {
      styleOverrides: {
        root: {
          height: '0.5px',
        },
      },
    },
    MuiTypography: {
      styleOverrides: {
        body1: {
          fontSize: '0.8rem',
        },
        h5: {
          textTransform: 'capitalize',
          fontSize: '1.1rem',
          fontWeight: 'bold',
        },
        subtitle2: {
          fontSize: '0.8rem',
          fontStyle: 'italic',
          color: '#565656',
        },
      },
    },
    MuiInputBase: {
      styleOverrides: {
        root: {
          fontSize: '0.8rem',
          minRows: 5,
        },
      },
    },
    MuiTooltip: {
      styleOverrides: {
        tooltip: {
          backgroundColor: '#3C494F',
          color: '#fff',
          opacity: '100%',
          fontSize: '14px',
          borderRadius: '0.9375rem',
          padding: '0.9rem',
          zIndex: 99999999999,
        },
        popper: {
          zIndex: 99999,
        },
      },
    },
    MuiGrid: {
      styleOverrides: {
        root: {
          '& > *': {
            border: 'none !important',
          },
          marginTop: '0.1rem !important',
          overflow: 'hidden',
          alignSelf: 'center',
          textOverflow: 'ellipsis',
          '&:hover': {
            overflow: 'visible',
          },
          background: 'none',
          [breakpoints.up('lg')]: {
            '& > *:nth-child(2)': {
              '& > *:nth-child(1)': {
                '& > *:nth-child(2)': {
                  justifyContent: 'space-around',
                },
              },
            },
          },
        },
      },
    },
    MuiIconButton: {
      styleOverrides: {
        root: {
          cursor: 'pointer',
        },
        sizeSmall: {
          padding: '1px',
        },
      },
    },
    MuiPaper: {
      styleOverrides: {
        elevation0: {
          backgroundColor: 'inherit',
        },
        elevation2: {
          boxShadow: 'none',
        },
      },
    },
  },
});

export default darkRjsfTheme;
