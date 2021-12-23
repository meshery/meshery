import { MuiThemeProvider } from '@material-ui/core/styles';
import { withTheme } from "@rjsf/core";
import { Theme as MaterialUITheme } from "@rjsf/material-ui";
import React from "react";
import JS4 from "../../../assets/jsonschema/schema-04.json";
import { rjsfTheme } from "../../../themes";
import { formatString } from "../helpers";
import ArrayFieldTemplate from "./RJSFCustomComponents/ArrayFieldTemlate";
import MemoizedCustomInputField from "./RJSFCustomComponents/CustomInputField";
import CustomObjFieldTemplate from "./RJSFCustomComponents/ObjectFieldTemplate";
import { isEqualArr } from "../../../utils/utils"

const Form = withTheme(MaterialUITheme);

function deleteTitleFromJSONSchema(jsonSchema) {
  return { ...jsonSchema, title : "" };
}

function deleteDescriptionFromJSONSchema(jsonSchema) {
  return { ...jsonSchema, description : "" }
}

/**
 * remove top-level title, top-level description and
 * replace internal description with "help" key for
 * tooltip description
 *
 * @param {Object.<String, Object>} jsonSchema
 * @returns
 */
function getRefinedJsonSchema(jsonSchema, hideTitle = true) {
  let refinedSchema;
  refinedSchema = hideTitle ? deleteTitleFromJSONSchema(jsonSchema) : jsonSchema
  refinedSchema = deleteDescriptionFromJSONSchema(refinedSchema)
  refinedSchema = addTitleToPropertiesJSONSchema(refinedSchema)
  return refinedSchema
}

function uiSchema(jsonSchema) {
  let uiJsonSchema = {}
  // TODO: needs to do recursively for deep fields
  Object.keys(jsonSchema.properties).map(key => {
    if (
      isEqualArr(jsonSchema.properties[key].type, ["string", "null"], false)
      || isEqualArr(jsonSchema.properties[key].type, ["integer", "null"], false)
    ) {
      console.log("into")
      uiJsonSchema[key] = {
        'ui:description' : ' '
      }
    }
  })

  return uiJsonSchema
}

function addTitleToPropertiesJSONSchema(jsonSchema) {
  const newProperties = jsonSchema?.properties

  if (newProperties && typeof newProperties === 'object') {
    Object.keys(newProperties).map(key => {
      if (Object.prototype.hasOwnProperty.call(newProperties, key)) {
        let defaultValue;
        let types = []
        if (!Array.isArray(newProperties[key].type) && Object.prototype.hasOwnProperty.call(newProperties[key], 'type')) {
          types.push(newProperties[key].type)
        } else {
          types.push(...newProperties[key].type)
        }
        if (types.includes('null')) {
          defaultValue = null
        } else if (types.includes('integer')) {
          defaultValue = 0
        } else if (types.includes('string')) {
          defaultValue = ''
        } else if (types.includes('array')) {
          defaultValue = []
        }
        newProperties[key] = {
          ...newProperties[key],
          title : formatString(key),
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

// function RJSFButton({ handler, text, ...restParams }) {
//   return (
//     <Button variant="contained" color="primary" style={{ marginRight : "0.5rem" }} onClick={handler} {...restParams}>
//       {text}
//     </Button>
//   );
// }

function RJSF(props) {
  const {
    formData,
    jsonSchema,
    onChange,
    hideTitle,
    RJSFWrapperComponent = React.Fragment,
    RJSFFormChildComponent = React.Fragment, // eslint-disable-line no-unused-vars
    //.. temporarily ignoring till handler is attached successfully
  } = props;

  // define new string field
  const fields = {
    StringField : ({ idSchema, formData, ...props }) => <MemoizedCustomInputField id={idSchema['$id']} value={formData} idSchema={idSchema} {...props} />
  }

  const [data, setData] = React.useState(prev => ({ ...formData, ...prev }));

  React.useEffect(() => {
    // Apply debouncing mechanism for the state propagation
    const timer = setTimeout(() => {
      onChange?.(data);
    }, 300);

    return () => clearTimeout(timer);
  }, [data]);

  return (
    <RJSFWrapperComponent {...{ ...props, RJSFWrapperComponent : null, RJSFFormChildComponent : null }}>
      <MuiThemeProvider theme={rjsfTheme}>
        <Form
          schema={getRefinedJsonSchema(jsonSchema, hideTitle)}
          idPrefix={jsonSchema?.title}
          onChange={(e) => {
            setData(e.formData)
          }}
          formData={data}
          fields={fields}
          ArrayFieldTemplate={ArrayFieldTemplate}
          ObjectFieldTemplate={CustomObjFieldTemplate}
          additionalMetaSchemas={[JS4]}
          uiSchema={uiSchema(jsonSchema)}
        >
          {/* {hideSubmit ? true : <RJSFButton handler={onSubmit} text="Submit" {...restparams} />}
        {hideSubmit ? true : <RJSFButton handler={onDelete} text="Delete" />} */}
          {/* <RJSFFormChildComponent /> */}
          <></> {/* temporary change for functionality */}
        </Form>
      </MuiThemeProvider>
    </RJSFWrapperComponent>
  );
}

export default RJSF;
