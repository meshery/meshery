import { createTheme } from "@material-ui/core";

function configurationTableTheme ( style = {} ) {
  const theme = Object.assign({
    shadows : ["none"],
    overrides : {
      MUIDataTable : {
      },
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
            color : '#607d8b',
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
    }
  }, style)
  return createTheme(theme);
}

function configurationTableThemeDark ( style = {} )  {
  const theme = Object.assign({
    shadows : ["none"],
    palette : {
      type : "dark",
    },
    overrides : {
      MuiPaper : { root : { backgroundColor : '#363636' } },
      MuiFormLabel : {
        root : {
          "&$focused" : {
            color : "#00B39F",
          },
        }
      },
    }
  }, style)
  return createTheme(theme);
}

export { configurationTableTheme, configurationTableThemeDark };