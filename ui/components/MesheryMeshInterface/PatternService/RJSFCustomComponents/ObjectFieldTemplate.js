import React from 'react';
import Grid from '@material-ui/core/Grid';
import { makeStyles } from '@material-ui/styles';
import { canExpand } from '@rjsf/utils';
import { Box, CssBaseline, IconButton, Typography, useTheme } from '@material-ui/core';
import AddIcon from '../../../../assets/icons/AddIcon';
import { EnlargedTextTooltip } from '../EnlargedTextTooltip';
import HelpOutlineIcon from '../../../../assets/icons/HelpOutlineIcon';
import ExpandMoreIcon from '../../../../assets/icons/ExpandMoreIcon';
import ExpandLessIcon from '../../../../assets/icons/ExpandLessIcon'
import createBreakpoints from '@material-ui/core/styles/createBreakpoints'
const breakpoints = createBreakpoints({});
const useStyles = makeStyles((theme) => ({
  objectFieldGrid : {
    // paddingLeft: "0.6rem",
    padding : ".5rem",
    paddingTop : "0.7rem",
    // margin : ".5rem",
    backgroundColor : theme.palette.type === 'dark' ? "#545454" : "#f4f4f4",
    border : `1px solid  ${theme.palette.type === 'dark' ? "rgba(255, 255, 255, .45)" : "rgba(0, 0, 0, .125)"}`,
    width : "100%",
    margin : "0px"
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
  TopClass : {

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



const ObjectFieldTemplate = ({
  description,
  title,
  properties,
  // required,
  disabled,
  readonly,
  uiSchema,
  idSchema,
  schema,
  formData,
  onAddClick,
}) => {
  const additional = schema?.__additional_property; // check if the object is additional
  const classes = useStyles();
  const theme = useTheme();
  // If the parent type is an `array`, then expand the current object.
  const [show, setShow] = React.useState(schema?.p_type ? true : false);
  properties.forEach((property, index) => {
    if (schema.properties[property.name].type) {
      properties[index].type = schema.properties[property.name].type;
      properties[index].__additional_property =
        schema.properties[property.name]?.__additional_property || false;
    }
  });
  const CustomTitleField = ({ title, id, description, properties }) => {
    return <Box mb={1} mt={1} id={id} >
      <CssBaseline />
      <Grid container justify="flex-start" alignItems="center">
        {canExpand(schema, uiSchema, formData) ? (
          <Grid item={true} onClick={() => {
            if (!show) setShow(true);
          }}>
            <IconButton
              className="object-property-expand"
              onClick={onAddClick(schema)}
              disabled={disabled || readonly}
            >
              <AddIcon width="18px" height="18px" fill="white" style={{ backgroundColor : "#647881", width : "1.25rem", height : "1.25rem", color : "#ffffff", borderRadius : ".2rem" }} />
            </IconButton>
          </Grid>
        ) : (
          Object.keys(properties).length > 0 && (
            <Grid item={true}>
              <IconButton
                className="object-property-expand"
                onClick={() => setShow(!show)}
              >
                {show ? <ExpandLessIcon width="18px" height="18px" fill="gray" /> : <ExpandMoreIcon width="18px" height="18px" fill="gray" />}
              </IconButton>
            </Grid>
          )
        )}

        <Grid item mb={1} mt={1}>
          <Typography variant="body1" className={classes.typography} style={{ fontWeight : "bold", display : "inline" }}>{title.charAt(0).toUpperCase() + title.slice(1)}{" "}
          </Typography>
          {description &&
            <EnlargedTextTooltip title={description}>
              <IconButton disableTouchRipple="true" disableRipple="true" component="span" size="small">
                <HelpOutlineIcon width="14px" height="14px" fill={theme.palette.type === 'dark' ? "white" : "black"} style={{ marginLeft : "4px", verticalAlign : "middle" }} />
              </IconButton>
            </EnlargedTextTooltip>}
        </Grid>


      </Grid>
    </Box>
  };

  const Properties = (<Grid container={true} spacing={2} className={classes.objectFieldGrid} style={Object.keys(properties).length === 0 || schema["$schema"] ? { border : "none" } : null}>
    {properties.map((element, index) => {
      return (
        element.hidden ? (
          element.content
        ) : (
          <Grid
            item={true}
            sm={12}
            lg={
              element.type === "object" ||
                element.type === "array" ||
                element.__additional_property ||
                additional
                ? 12
                : 6
            }
            key={index}
          >
            {element.content}
          </Grid>
        )
      );
    })}
  </Grid>
  )

  const fieldTitle = uiSchema['ui:title'] || title;

  return (
    <>
      {fieldTitle ? (
        <>
          {schema.p_type !== "array" ? (
            <CustomTitleField
              id={`${idSchema.$id}-title`}
              title={additional ? "Value" : fieldTitle}
              description={description}
              properties={properties}
            />
          ) : null
          }
          {Object.keys(properties).length > 0 && show && Properties}
        </>
      ) : Properties}
    </>
  );
};

export default ObjectFieldTemplate;
