import { createTheme } from '@mui/material/styles'
import { blueGrey } from '@mui/material/colors'

// call createTheme to extend the settings
const defaultTheme = createTheme({})

// Create a theme instance.
const theme = createTheme({
  palette : {
    primary : {
      main : blueGrey[500]
    },
    warning : {
      main : '#F0D053'
    }
  },
  breakpoints : {
    values : {
      xs : 0,
      sm : 600,
      md : 960,
      lg : 1280,
      xl : 1920
    }
  },
  components : {
    MuiButton : {
      styleOverrides : {
        root : {
          textTransform : "none"
        }
      }
    },
    MuiDialogTitle : {
      styleOverrides : {
        root : {
          margin : 0,
          padding : defaultTheme.spacing(2)
        }
      }
    },
    MuiDialogContent : {
      styleOverrides : {
        root : {
          padding : defaultTheme.spacing(2)
        }
      }
    },
    MuiDialogActions : {
      styleOverrides : {
        root : {
          padding : defaultTheme.spacing(1)
        }
      }
    },
    MuiPaper : {
      styleOverrides : {
        root : {
          margin : 'auto',
          overflow : 'hidden',
          maxWidth : "90%",
        }
      }
    },
    MuiIcon : {
      styleOverrides : {
        root : {
          fontSize : 20
        }
      }
    },
    MuiList : {
      styleOverrides : {
        root : {
          width : 'fit-content'
        }
      }
    }
  }
})

export default theme
