import { makeStyles, useTheme } from '@material-ui/core/styles';
import Form, { withTheme } from "@rjsf/core";
import { Theme as MaterialUITheme, } from "@rjsf/material-ui";
import React, { useEffect } from "react";
import JS4 from "../../../assets/jsonschema/schema-04.json";
import { rjsfTheme } from "../../../themes";
import { recursiveCleanObject } from "../helpers";
import MesheryArrayFieldTemplate from "./RJSFCustomComponents/ArrayFieldTemlate";
import MesheryCustomObjFieldTemplate from "./RJSFCustomComponents/ObjectFieldTemplate";
import MesheryWrapIfAdditionalTemplate from './RJSFCustomComponents/WrapIfAdditionalTemplate';
import { customizeValidator } from "@rjsf/validator-ajv6";
import _ from "lodash"
import CustomTextWidget from './RJSFCustomComponents/CustomTextWidget';
import CustomDateTimeWidget from './RJSFCustomComponents/CustomDateTimeWidget';
import createBreakpoints from '@material-ui/core/styles/createBreakpoints'
const breakpoints = createBreakpoints({});
const useStyles = makeStyles((theme) => ({
  root : {
    "& .MuiFormLabel-root" : {
      "& .Mui-error" : {
        color : "#f44336"
      },
      textTransform : "capitalize",
      fontSize : "0.8rem",

    },
    "& .MuiPaper-root" : {
      backgroundColor : "inherit",
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
    "& .MuiOutlinedInput-root" : {
      position : "relative",
      borderRadius : "4px",
      backgroundColor : theme.palette.type === 'dark' ? "#606060" : "#ffffff",

    },
    "& .MuiButton" : {
      textSecondary : {
        color : "#00b39f",
        "&:hover" : "00b39f"
      }
    },
    "& .MuiTextField-root" : {

      width : "calc(100% - 4px)",

    },
    "& .MuiInputLabel-root" : {
      whiteSpace : "nowrap",
      overflow : "hidden",
      textOverflow : "ellipsis",
      maxWidth : "60%",
      height : "100%",
      '&:hover' : {
        overflow : "visible",
      },

      shrink : {
        maxWidth : "100%",
      }
    },
    "& .MuiFormControl-root" : {
      marginTop : "0.3rem"
    },
    "& .MuiBox-root" : {
      marginTop : 0
    },
    "& .MuiDivider-root" : {

      height : "0.5px"

    },

    "& .MuiTypography-body1" : {

      fontSize : '0.8rem',
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
    "& .MuiInputBase-root" : {
      fontSize : "0.8rem" // same as title
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
    "& .MuiAccordionSummary-root" : {

      // border: "5px solid red",
      backgroundColor : theme.palette.type === 'dark' ? "inherit" : "rgba(242,242,242)",
      borderBottom : `1px solid  ${theme.palette.type === 'dark' ? "rgba(255, 255, 255, .45)" : "rgba(0, 0, 0, .125)"}`,
      marginBottom : -1,
      maxHeight : "1.5rem",
      "&$expanded" : {
        minHeight : 56
      }
    },
    "& .MuiAccordionSummary-content" : {
      justifyContent : "space-between",

    },

    MuiAccordionDetails : {
      root : {
        padding : 16
      }
    },
    "& .MuiAccordion-root" : {
      border : `1px solid  ${theme.palette.type === 'dark' ? "rgba(255, 255, 255, .45)" : "rgba(0, 0, 0, .125)"}`,
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
      backgroundColor : theme.palette.type === 'dark' ? "inherit" : "rgba(242,242,242)",

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
    "& .MuiCheckbox-root" : {
      // checkboxes white background

      marginLeft : "4px",
      "& > *:nth-child(1)" : {
        backgroundColor : "#ffffff",
        width : "1rem",
        height : "1rem"
      }
    },
    "& .MuiIconButton-root" : {

      cursor : "default",
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
/*eslint-disable */
class RJSFOverridenComponent extends Form {
  constructor(props) {
    try {
      super(props)
      let oldValidate = this.validate;
      this.validate = (
        formData,
        schema,
      ) => {
        let fixedFormData = recursiveCleanObject(_.cloneDeep(formData));
        return oldValidate.call(this, fixedFormData, schema);
      }
    } catch (e) {
      console.error("An RJSF error occurred", e)
    }
  }
}
/*eslint-enable */

// This is Patched change to include customised Forms
const MuiRJSFForm = withTheme(MaterialUITheme, RJSFOverridenComponent);
const validator = customizeValidator({ additionalMetaSchemas : [JS4] });

/**
 * The Custom RJSF Form that accepts custom fields from the extension
 * or seed it's own default
 * Adding a new custom component:
 * 1. Pass the new prop from the Meshery Extension
 * 2. Extract from props in the RJSFForm Component
 * @param {*} props
 * @returns
 */
function RJSFForm(props) {
  const {
    schema,
    jsonSchema,
    data,
    onChange,
    isLoading,
    ArrayFieldTemplate = MesheryArrayFieldTemplate,
    ObjectFieldTemplate = MesheryCustomObjFieldTemplate,
    WrapIfAdditionalTemplate = MesheryWrapIfAdditionalTemplate,
    LoadingComponent,
    ErrorList,
    // prop should be present in order for the cloned element to override this property
    transformErrors
  } = props;
  const templates = {
    ArrayFieldTemplate,
    ObjectFieldTemplate,
    WrapIfAdditionalTemplate,
  }

  useEffect(() => {
    const extensionTooltipPortal = document.getElementById("extension-tooltip-portal");
    if (extensionTooltipPortal) {
      rjsfTheme.props.MuiMenu.container = extensionTooltipPortal;
    }
    rjsfTheme.zIndex.modal = 99999;
  }, [])

  if (isLoading && LoadingComponent) {
    return <LoadingComponent />
  }
  const theme = useTheme();
  const classes = useStyles();
  console.log("hello", theme)
  return (

    <div>
      <MuiRJSFForm
        schema={schema.rjsfSchema}
        idPrefix={jsonSchema?.title}
        onChange={onChange}
        formData={data}
        className={classes.root}
        validator={validator}
        templates={templates}
        uiSchema={schema.uiSchema}
        widgets={{
          TextWidget : CustomTextWidget,
          // Custom components to be added here
          DateTimeWidget : CustomDateTimeWidget,
          // SelectWidget: CustomSelectWidget,
          // CheckboxWidget: CustomBooleanWidget,
        }}
        liveValidate
        showErrorList={false}
        noHtml5Validate
        ErrorList={ErrorList}
        transformErrors={transformErrors}
      >
        {/* {hideSubmit ? true : <RJSFButton handler={onSubmit} text="Submit" {...restparams} />}
{hideSubmit ? true : <RJSFButton handler={onDelete} text="Delete" />} */}
        {/* <RJSFFormChildComponent /> */}
        <></> {/* temporary change for functionality */}
      </MuiRJSFForm>
    </div>
  )
}

export default RJSFForm;
