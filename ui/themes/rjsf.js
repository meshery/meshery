
import { createTheme } from '@material-ui/core/styles';


export const rjsfTheme = createTheme({
  palette : {
    primary : {
      main : '#607d8b',
    },
  },
  typography : {
    fontSize : 13,
  },
  props : {
    MuiTextField : {
      variant : 'outlined',
      margin : 'dense',
    },
    MuiCheckbox : {
      color : 'primary',
    },
    MuiMenu : {
      variant : "outlined",
    }
  },
  overrides : {
    MuiButton : {
      textSecondary : {
        color : "#00b39f",
        "&:hover" : "00b39f"
      }
    },
    MuiBox : {
      root : {
        marginTop : 0
      }
    },
    MuiDivider : {
      root : {
        height : "0.5px"
      }
    },
    MuiFormLabel : {
      root : {
        color : "#333",
        fontSize : "0.8rem",
        textTransform : "capitalize"
      }
    },
    MuiTypography : {
      body1 : {
        fontSize : '0.8rem'
      },
      h5 : {
        textTransform : 'capitalize',
        fontSize : '1.1rem',
        fontWeight : "bold"
      },
      subtitle2 : {
        fontSize : '0.8rem',
        fontStyle : "italic",
        color : "#565656"
      },
    },
    MuiInputBase : {
      root : {
        fontSize : "0.8rem" // same as title
      }
    },
    MuiGrid : {
      root : {
        "& > *" : {
          border : 'none !important'
        },
        marginTop : '0.1rem !important',
      },
    },
    MuiPaper : {
      elevation2 : {
        boxShadow : "none"
      }
    },
    "spacing-xs-2" : {
      padding : 0,
      "& > *" : {
        paddingTop : "0 !important",
        paddingBottom : "0 !important"
      }
    }
  }
});