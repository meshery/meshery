import React, { useEffect, useState } from "react";
import { withTheme } from "@rjsf/core";
import { Theme as MaterialUITheme } from "@rjsf/material-ui";
import { TextField, Typography } from "@material-ui/core";
import JS4 from "../../../assets/jsonschema/schema-04.json";
import { Tooltip, IconButton } from "@material-ui/core";
import HelpOutlineIcon from "@material-ui/icons/HelpOutline";

const Form = withTheme(MaterialUITheme);

function deleteTitleFromJSONSchema(jsonSchema) {
  return { ...jsonSchema, title : "" };
}

function deleteDescriptionFromJSONSchema(jsonSchema) {
  return { ...jsonSchema, description : "" };
}

function formatString(text) {
  if (!text) return null;

  // format string for prettified camelCase
  return text.replaceAll("IP", "Ip");
}

function camelCaseToCapitalize(text) {
  if (!text) return null;

  const result = text.replace(/([A-Z])/g, " $1");

  return result.charAt(0).toUpperCase() + result.slice(1);
}

function addTitleToPropertiesJSONSchema(jsonSchema) {
  const newProperties = jsonSchema?.properties;

  if (newProperties && typeof newProperties === 'object') {
    Object.keys(newProperties).map(key => {
      if (Object.prototype.hasOwnProperty.call(newProperties, key)) {
        let defaultValue;
        let types = [];
        if (!Array.isArray(newProperties[key].type) && Object.prototype.hasOwnProperty.call(newProperties[key], 'type')) {
          types.push(newProperties[key].type);
        } else {
          types.push(...newProperties[key].type);
        }
        if (types.includes('null')) {
          defaultValue = null;
        } else if (types.includes('integer')) {
          defaultValue = 0;
        } else if (types.includes('string')) {
          defaultValue = '';
        } else if (types.includes('array')) {
          defaultValue = [];
        }
        newProperties[key] = {
          ...newProperties[key],
          title : camelCaseToCapitalize(formatString(key)),
          default : defaultValue
        };
        // if (typeof newProperties[key] === 'object' && Object.prototype.hasOwnProperty.call(newProperties[key], 'properties')){
        //   newProperties[key] = {
        //     ...newProperties[key],
        //     properties : addTitleToPropertiesJSONSchema(newProperties[key])
        //   }
        // }
      }

    });

    return { ...jsonSchema, properties : newProperties };
  }
  return undefined;
}

const CustomInputField = (props) => {
  const name = props?.name || props?.idSchema['$id']?.split('_')[1];
  const prettifiedName = camelCaseToCapitalize(formatString(name)) || 'Input';
  return (
    <div key={props.id}>
      <Typography variant="body1" style={{ fontWeight : "bold" }}>{prettifiedName}
        {props.schema?.description && (
          <Tooltip title={props.schema?.description}>
            <IconButton component="span" size="small">
              <HelpOutlineIcon style={{ fontSize : 17 }} />
            </IconButton>
          </Tooltip>
        )}
      </Typography>
      <TextField variant="outlined" size="small" style={{ margin : '10px 0 ' }} autoFocus key={props.id} value={props.value} id={props.id} onChange={e => props?.onChange(e.target.value)} placeholder={`${prettifiedName}`} />
    </div>
  );
};

const MemoizedCustomInputField = React.memo(CustomInputField);

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
    RJSFFormChildComponent = React.Fragment,
  } = props;

  // define new string field
  const fields =  {
    StringField : ({ idSchema, formData, ...props }) => <MemoizedCustomInputField id={idSchema['$id']} value={formData} idSchema={idSchema} {...props} />
  }

  const [data, setData] = React.useState(prev => ({ ...formData, ...prev }));


  return (
    <RJSFWrapperComponent {...{ ...props, RJSFWrapperComponent : null, RJSFFormChildComponent : null } }>
      <Form
        schema={hideTitle ? deleteTitleFromJSONSchema(deleteDescriptionFromJSONSchema(addTitleToPropertiesJSONSchema(jsonSchema))) : deleteDescriptionFromJSONSchema(addTitleToPropertiesJSONSchema(jsonSchema))}
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
        <RJSFFormChildComponent />
      </Form>
    </RJSFWrapperComponent>
  );
}

export default RJSF;
