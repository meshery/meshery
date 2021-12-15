import React from "react";
import { withTheme } from "@rjsf/core";
import { Theme as MaterialUITheme } from "@rjsf/material-ui";
import { TextField, Typography } from "@material-ui/core";
import JS4 from "../../../assets/jsonschema/schema-04.json";
import { IconButton } from "@material-ui/core";
import HelpOutlineIcon from "@material-ui/icons/HelpOutline";
import { MuiThemeProvider } from '@material-ui/core/styles';
import { rjsfTheme } from "../../../themes";
import { camelCaseToCapitalize } from "../helpers";
import EnlargedTextTooltip from "./EnlargedTextTooltip";

const Form = withTheme(MaterialUITheme);

function deleteTitleFromJSONSchema(jsonSchema) {
  return { ...jsonSchema, title : "" };
}

function deleteDescriptionFromJSONSchema(jsonSchema) {
  return { ...jsonSchema, description : "" }
}

/**
 * Replace description properties of each items inside
 * form with "help" key
 *
 * @param {Object.<String, Object>} jsonSchema
 * @returns
 */
function deleteInternalDescriptionAndAddHelpDesc(jsonSchema) {
  Object.keys(jsonSchema.properties).forEach(key => {
    // Create a new object-field to store the description inside tooltip
    if (!jsonSchema.properties[key]['help']) {
      jsonSchema.properties[key]['help'] = jsonSchema.properties[key]['description']
    }

    delete jsonSchema.properties[key]['description']
  })
  return jsonSchema
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
  refinedSchema =  hideTitle ?  deleteTitleFromJSONSchema(jsonSchema) : jsonSchema
  refinedSchema = deleteDescriptionFromJSONSchema(refinedSchema)
  refinedSchema = addTitleToPropertiesJSONSchema(refinedSchema)
  refinedSchema = deleteInternalDescriptionAndAddHelpDesc(refinedSchema)
  return refinedSchema
}

function formatString(text){
  if (!text) return null

  // format string for prettified camelCase
  let formattedText = text.replaceAll("IP", "Ip");
  formattedText = camelCaseToCapitalize(formattedText),
  formattedText = formattedText.replaceAll("Ip", "IP")
  return formattedText
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

const CustomInputField = (props) => {
  const name = props?.name || props?.idSchema['$id']?.split('_')[1]
  const prettifiedName = formatString(name) || 'Input'
  return (
    <div key={props.id}>
      <Typography variant="body1" style={{ fontWeight : "bold" }}>{prettifiedName}
        {props.schema?.help && (
          <EnlargedTextTooltip title={props.schema?.help}>
            <IconButton component="span" size="small">
              <HelpOutlineIcon style={{ fontSize : 17 }} />
            </IconButton>
          </EnlargedTextTooltip>
        )}
      </Typography>
      <TextField variant="outlined" size="small" style={{ margin : '10px 0 ' }} autoFocus key={props.id} value={props.value} id={props.id} onChange={e => props?.onChange(e.target.value)} placeholder={`${prettifiedName}`}/>
    </div>
  )
}

const MemoizedCustomInputField = React.memo(CustomInputField)

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
  const fields =  {
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
    <RJSFWrapperComponent {...{ ...props, RJSFWrapperComponent : null, RJSFFormChildComponent : null } }>
      <MuiThemeProvider theme={rjsfTheme}>
        <Form
          schema={getRefinedJsonSchema(jsonSchema, hideTitle)}
          idPrefix={jsonSchema?.title}
          onChange={(e) => {
            setData(e.formData)
          }}
          formData={data}
          fields={fields}
          additionalMetaSchemas={[JS4]}
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
