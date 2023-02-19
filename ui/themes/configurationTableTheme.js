import { createTheme } from "@material-ui/core";

const configurationTableTheme = () => createTheme({
  overrides : {
    MuiInput : {
      underline : {
        "&:hover:not(.Mui-disabled):before" : {
          borderBottom : "2px solid #222"
        },
        "&:after" : {
          borderBottom : "2px solid #222"
        }
      }
    },
    MUIDataTableSearch : {
      searchIcon : {
        color : "#607d8b",
        marginTop : "7px",
        marginRight : "8px",
      },
      clearIcon : {
        "&:hover" : {
          color : "#607d8b"
        }
      },
    },
    MUIDataTableSelectCell : {
      checkboxRoot : {
        '&$checked' : {
          color : '#3C494F',
        },
      },
    },
    MUIDataTableToolbar : {
      iconActive : {
        color : "#222"
      },
      icon : {
        "&:hover" : {
          color : "#607d8b"
        }
      },
    },
    MUIDataTableBodyCell : {
      root : {
        cursor : "pointer"
      },
    },
    MuiFormLabel : {
      root : {
        "&$focused" : {
          color : "#00B39F",
        },
      }
    },
  }
});

export default configurationTableTheme;