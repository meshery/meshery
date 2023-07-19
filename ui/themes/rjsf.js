
import { createTheme } from '@material-ui/core/styles';
import createBreakpoints from '@material-ui/core/styles/createBreakpoints'
const breakpoints = createBreakpoints({});
export const rjsfTheme = createTheme({
  palette : {
    primary : {
      main : '#607d8b',
    },
  },
  typography : {
    fontFamily : "inherit",
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
    MuiOutlinedInput : {
      root : {
        backgroundColor : '#ffffff',
      },
    },
    MuiButton : {
      textSecondary : {
        color : "#00b39f",
        "&:hover" : "00b39f"
      }
    },
    MuiTextField : {
      root : {
        width : "calc(100% - 4px)",
      }
    },
    MuiInputLabel : {
      root : {
        whiteSpace : "nowrap",
        overflow : "hidden",
        textOverflow : "ellipsis",
        maxWidth : "60%",
        height : "100%",
        '&:hover' : {
          overflow : "visible",
        }
      },
      shrink : {
        maxWidth : "100%",
      }
    },
    MuiFormControl : {
      root : {
        marginTop : "0.3rem"
      }
    },
    MuiFormControlLabel : {
      root : {
        textTransform : 'capitalize',
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
        fontSize : '0.8rem',
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
        backgroundColor : "rgba(242,242,242)",
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
        },
        backgroundColor : "rgba(242,242,242)",
      },
    },
    MuiGrid : {
      root : {
        "& > *" : {
          border : 'none !important'
        },
        marginTop : '0.1rem !important',
        overflow : "hidden",
        alignSelf : "center",
        textOverflow : "ellipsis",
        '&:hover' : {
          overflow : "visible",
        },
        // To scale the grid items on a particular screen size
        [breakpoints.up('lg')] : {
          "& > *:nth-child(2)" : {
            "& > *:nth-child(1)" : {
              "& > *:nth-child(2)" : {
                justifyContent : 'space-around',
              }
            },
          }
        },
      },
    },
    MuiCheckbox : {
      // checkboxes white background
      root : {
        marginLeft : "4px",
        "& > *:nth-child(1)" : {
          backgroundColor : "#ffffff",
          width : "1rem",
          height : "1rem"
        },
      }
    },
    MuiIconButton : {
      root : {
        cursor : "pointer"
      },
      sizeSmall : {
        padding : "1px"
      }
    },
    MuiPaper : {
      elevation0 : {
        backgroundColor : "inherit",
      },
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

const darkRjsfTheme = createTheme({
  palette : {
    type : "dark",
    primary : {
      main : '#607d8b',

    },
  },
  typography : {
    fontFamily : "inherit",
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
    MuiOutlinedInput : {
      root : {
        backgroundColor : '#303030',
      },
    },
    MuiButton : {
      textSecondary : {
        color : "#00b39f",
        "&:hover" : "00b39f"
      }
    },
    MuiTextField : {
      root : {
        width : "calc(100% - 4px)",
      }
    },
    MuiFormControlLabel : {
      root : {
        textTransform : 'capitalize',
      }
    },
    MuiInputLabel : {
      root : {
        whiteSpace : "nowrap",
        overflow : "hidden",
        textOverflow : "ellipsis",
        maxWidth : "60%",
        height : "100%",
        '&:hover' : {
          overflow : "visible",
        }
      },
      shrink : {
        maxWidth : "100%",
      }
    },
    MuiFormControl : {
      root : {
        marginTop : "0.3rem"
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
        color : "#FFF",
        fontSize : "0.8rem",
        textTransform : "capitalize"
      }
    },
    MuiTypography : {
      body1 : {
        fontSize : '0.8rem',
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
        backgroundColor : "#303030",
        borderBottom : "1px solid rgba(255, 255, 255, .125)",
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
        border : "1px solid rgba(255, 255, 255, .125)",
        boxShadow : "none",
        "&:not(:last-child)" : {
          borderBottom : 0
        },
        "&:before" : {
          display : "none"
        },
        "&$expanded" : {
          margin : "auto"
        },
        backgroundColor : "#303030",
      },
    },
    MuiGrid : {
      root : {
        "& > *" : {
          border : 'none !important'
        },
        marginTop : '0.1rem !important',
        overflow : "hidden",
        alignSelf : "center",
        textOverflow : "ellipsis",
        '&:hover' : {
          overflow : "visible",
        },
        // To scale the grid items on a particular screen size
        [breakpoints.up('lg')] : {
          "& > *:nth-child(2)" : {
            "& > *:nth-child(1)" : {
              "& > *:nth-child(2)" : {
                justifyContent : 'space-around',
              }
            },
          }
        },
      },
    },
    MuiCheckbox : {
      // checkboxes white background
      root : {
        marginLeft : "4px",
        "& > *:nth-child(1)" : {
          backgroundColor : "#303030",
          width : "1rem",
          height : "1rem"
        },
      }
    },
    MuiIconButton : {
      root : {
        cursor : "default"
      },
      sizeSmall : {
        padding : "1px"
      }
    },
    MuiPaper : {
      elevation0 : {
        backgroundColor : "inherit",
      },
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
export default darkRjsfTheme;