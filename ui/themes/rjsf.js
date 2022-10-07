
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
    MuiTextField : {
      root : {
        marginRight : "0.5rem"
      }
    },
    MuiInputLabel : {
      root : {
        whiteSpace : "nowrap",
        // overflow : "hidden",
        textOverflow : "ellipsis",
        maxWidth : "75%",
        height : "100%",
        '&:hover' : {
          overflow : "visible",
        }
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
    MuiTooltip : {
      tooltip : {
        backgroundColor : " #3C494F",
        color : "#fff",
        opacity : "100%",
        fontSize : "14px",
        borderRadius : "0.9375rem",
        padding : "0.9rem",
        zIndex : "99999999999"
      },
      popper : {
        zIndex : "99999 !important"
      }
    },
    MuiAccordionSummary : {
      root : {
        // border: "5px solid red",
        backgroundColor : "rgba(0, 0, 0, .03)",
        borderBottom : "1px solid rgba(0, 0, 0, .125)",
        marginBottom : -1,
        maxHeight : "1.5rem",
        "&$expanded" : {
          minHeight : 56
        }
      },
      content : {
        // border: "5px solid green",
        justifyContent : "space-between",
        "&$expanded" : {
          margin : "12px 0",
          justifyContent : "space-between"
        }
      },
    },
    MuiAccordionDetails : {
      root : {
        padding : 16
      }
    },
    MuiAccordion : {
      root : {
        border : "1px solid rgba(0, 0, 0, .125)",
        boxShadow : "none",
        "&:not(:last-child)" : {
          borderBottom : 0
        },
        "&:before" : {
          display : "none"
        },
        "&$expanded" : {
          margin : "auto"
        }
      },
    },
    MuiGrid : {
      root : {
        "& > *" : {
          border : 'none !important'
        },
        marginTop : '0.1rem !important',
        overflow : "hidden",
        textOverflow : "ellipsis",
        whiteSpace : "nowrap",
        '&:hover' : {
          overflow : "visible",
          whiteSpace : "normal",
        }
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