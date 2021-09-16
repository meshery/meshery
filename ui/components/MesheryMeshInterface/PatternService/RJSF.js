import React from "react";
import { withTheme } from "@rjsf/core";
import { Theme as MaterialUITheme } from "@rjsf/material-ui";
import { Button, Grid } from "@material-ui/core";
import { MuiThemeProvider, createTheme } from '@material-ui/core/styles';
import JS4 from "../../../assets/jsonschema/schema-04.json";
import ArrayFieldTemplate from "./RJSFCustomComponents/ArrayFieldTemlate"
import CustomObjectFieldTemplate from "./RJSFCustomComponents/ObjectFieldTemplate"

const Form = withTheme(MaterialUITheme);

const muiTheme = createTheme({
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
  },
  overrides : {
    MuiButton : {
      textSecondary : {
        color : '#00b39f',
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
        height : '0.5px'
      }
    },
    MuiFormLabel : {
      root : {
        color : "#333",
        fontSize : '0.8rem',
        textTransform : 'capitalize',
      }
    },
    MuiTypography : {
      body1 : {
        fontSize : '0.8rem'
      },
      h5 : {
        textTransform : 'capitalize',
        fontSize : '1.2rem',
      }
    },
    MuiGrid : {
      root : {
        "& > *" : {
          border : 'none !important'
        },
        marginTop : '0.2rem !important',
      },
      'spacing-xs-2' : {
        padding : 0,
        '& > *' : {
          paddingTop : '0 !important',
          paddingBottom : '0 !important'
        }

      }
    }
  }
});

function deleteTitleFromJSONSchema(jsonSchema) {
  return { ...jsonSchema, title : "" };
}

function deleteDescriptionFromJSONSchema(jsonSchema) {
  return { ...jsonSchema, description : "" };
}

function camelCaseToCapitalize(text){
  const result = text.replace(/([A-Z])/g, " $1");

  return result.charAt(0).toUpperCase() + result.slice(1);
}

function addTitleToPropertiesJSONSchema(jsonSchema) {
  const newProperties = jsonSchema?.properties

  if (newProperties && typeof newProperties === 'object'){
    Object.keys(newProperties).map(key => {
      if (Object.prototype.hasOwnProperty.call(newProperties, key)){
        let defaultValue;
        let types = []
        if (!Array.isArray(newProperties[key].type) && Object.prototype.hasOwnProperty.call(newProperties[key], 'type')){
          types.push(newProperties[key].type)
        } else {
          types.push(...newProperties[key].type)
        }
        if (types.includes('null')){
          defaultValue = null
        } else if (types.includes('integer')){
          defaultValue = 0
        } else if (types.includes('string')){
          defaultValue = ''
        } else if (types.includes('array')){
          defaultValue = []
        }
        newProperties[key] = {
          ...newProperties[key],
          title : camelCaseToCapitalize(key),
          default : defaultValue
        }
        // if (typeof newProperties[key] === 'object' && Object.prototype.hasOwnProperty.call(newProperties[key], 'properties')){
        //   newProperties[key] = {
        //     ...newProperties[key],
        //     properties : addTitleToPropertiesJSONSchema(newProperties[key])
        //   }
        // }
      }

    })

    return { ...jsonSchema, properties : newProperties };
  }
  return undefined
}

function RJSFButton({ handler, text, ...restParams }) {
  return (
    <Button variant="contained" color="primary" style={{ marginRight : "0.5rem" }} onClick={handler} {...restParams}>
      {text}
    </Button>
  );
}



function RJSF({ formData, jsonSchema, onChange, hideSubmit, hideTitle, onSubmit, onDelete, renderAsTooltip, ...restparams }) {
  const uiSchema = {
    replicas : {
      "ui:widget" : "range"
    },
    "ui:order" : [
      "name",
      "namespace",
      "*",
    ],

  };

  <Grid container xs={12}>
    <Grid item xs={6}>
    </Grid>
    <Grid item xs={6}>
    </Grid>
  </Grid>

  const [data, setData] = React.useState({ ...formData });

  React.useEffect(() => {
    onChange?.(data);
  }, [data]);

  return (
    <>
      {!renderAsTooltip ? (
        <Form
          schema={hideTitle ? deleteTitleFromJSONSchema(deleteDescriptionFromJSONSchema(addTitleToPropertiesJSONSchema(jsonSchema))) : deleteDescriptionFromJSONSchema(addTitleToPropertiesJSONSchema(jsonSchema))}
          idPrefix={jsonSchema?.title}
          onChange={(e) => setData(e.formData)}
          formData={data}
          liveValidate
          additionalMetaSchemas={[JS4]}
        // noHtml5Validate
        >
          {hideSubmit ? true : <RJSFButton handler={onSubmit} text="Submit" {...restparams} />}
          {hideSubmit ? true : <RJSFButton handler={onDelete} text="Delete" />}
        </Form>) : (
        <MuiThemeProvider theme={muiTheme}>
          <Form
            schema={hideTitle ? deleteTitleFromJSONSchema(deleteDescriptionFromJSONSchema(jsonSchema)) : deleteDescriptionFromJSONSchema(jsonSchema)}
            idPrefix={jsonSchema?.title}
            onChange={(e) => setData(e.formData)}
            formData={data}
            showErrorList={false}
            ObjectFieldTemplate={CustomObjectFieldTemplate}
            ArrayFieldTemplate={ArrayFieldTemplate}
            additionalMetaSchemas={[JS4]}
            uiSchema={uiSchema}
          // noHtml5Validate
          >
            <button style={{ opacity : '0' }} />
          </Form>
        </MuiThemeProvider>
      )}
    </>
  );
}

export default RJSF;
