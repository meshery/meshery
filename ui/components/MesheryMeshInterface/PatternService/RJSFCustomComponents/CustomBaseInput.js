import React from "react";
import { IconButton, InputAdornment, makeStyles, TextField, useTheme } from "@material-ui/core";
import { EnlargedTextTooltip, EnlargedTextErrorTooltip } from "../EnlargedTextTooltip";
import HelpOutlineIcon from "../../../../assets/icons/HelpOutlineIcon";
import ErrorOutlineIcon from "../../../../assets/icons/ErrorOutlineIcon";
import createBreakpoints from '@material-ui/core/styles/createBreakpoints'
const breakpoints = createBreakpoints({});
const BaseInput = (props) => {
  const additional = props.schema?.__additional_property; // check if the field is additional
  const name = (additional ? "Value" : props.label) // || props.id?.split('_')[-1].trim()
  const focused = props.options?.focused // true for datetime-local
  const prettifiedName = name || 'Enter a value'
  const style = {
    display : "flex",
    alignItems : "center",
  }

  const useStyles = makeStyles(() => ({
    root : {
      "& .MuiFormLabel-root" : {
        textTransform : "capitalize"
      },
      typography : {
        fontFamily : "inherit",
        fontSize : 13,

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
      },


      "& .MuiButton" : {
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
      "& .MuiGrid-root" : {
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


  }));
  const classes = useStyles();
  const theme = useTheme();
  return (
    <>
      <div key={props.id} style={style}>
        <TextField
          variant={additional ? "standard" : "outlined"}
          size="small"
          className={classes.root}
          focused={focused}
          type={props.options?.inputType}
          key={props.id}
          style={{ width : "100%" }}
          value={additional && props?.value === "New Value" ? "" : props?.value} // remove the default value i.e. New Value for additionalFields
          id={props.id}
          margin="dense"
          error={props.rawErrors?.length > 0}
          onChange={e => props?.onChange(e.target.value === "" ? props.options.emptyValue : e.target.value)}
          label={`${prettifiedName}`}
          InputProps={{
            style : { padding : "0px 0px 0px 0px", },
            endAdornment : (<InputAdornment position="start">
              {props.rawErrors?.length > 0 && (
                <EnlargedTextErrorTooltip title={
                  <div>
                    {props.rawErrors?.map((error, index) => (
                      <div key={index}>{error}</div>
                    ))}
                  </div>

                }>
                  <IconButton component="span" size="small">
                    <ErrorOutlineIcon width="14px" height="14px" fill="#B32700" style={{ verticalAlign : "middle" }} />
                  </IconButton>
                </EnlargedTextErrorTooltip>
              )}
              {props.schema?.description && (
                <EnlargedTextTooltip title={props.schema?.description}>
                  <IconButton component="span" size="small">
                    <HelpOutlineIcon width="14px" height="14px" fill={theme.palette.type === 'dark' ? "white" : "black"} style={{ verticalAlign : "middle" }} />
                  </IconButton>
                </EnlargedTextTooltip>
              )}
            </InputAdornment>),
          }} />
      </div>
    </>
  )
}

export default BaseInput;
