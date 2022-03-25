import { MuiThemeProvider } from '@material-ui/core/styles';
import { withTheme } from "@rjsf/core";
import { Theme as MaterialUITheme } from "@rjsf/material-ui";
import React from "react";
import JS4 from "../../../assets/jsonschema/schema-04.json";
import { rjsfTheme } from "../../../themes";
import handleError from '../../ErrorHandling';
import { buildUiSchema } from "../helpers";
import { getRefinedJsonSchema } from "./helper";
import MesheryArrayFieldTemplate from "./RJSFCustomComponents/ArrayFieldTemlate";
import MemoizedCustomInputField from "./RJSFCustomComponents/CustomInputField";
import MesheryCustomObjFieldTemplate from "./RJSFCustomComponents/ObjectFieldTemplate";

const Form = withTheme(MaterialUITheme);

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

  const errorHandler = handleError();

  // define new string field
  const fields = {
    StringField : ({ idSchema, formData, ...props }) => <MemoizedCustomInputField id={idSchema['$id']} value={formData} idSchema={idSchema} {...props} />
  }

  const [data, setData] = React.useState(prev => ({ ...formData, ...prev }));
  const [schema, setSchema] = React.useState({ rjsfSchema : {}, uiSchema : {} })

  React.useEffect(() => {
    // Apply debouncing mechanism for the state propagation
    const timer = setTimeout(() => {
      onChange?.(data);
    }, 300);

    return () => clearTimeout(timer);
  }, [data]);

  React.useEffect(() => {
    const rjsfSchema = getRefinedJsonSchema(jsonSchema, hideTitle, errorHandler)
    const uiSchema = buildUiSchema(rjsfSchema)
    setSchema({ rjsfSchema, uiSchema })
  }, [jsonSchema]) // to reduce heavy lifting on every react render

  return (
    <RJSFWrapperComponent {...props}>
      <RJSFForm
        schema={schema}
        data={data}
        onChange={(e) => {
          setData(e.formData)
        }}
        jsonSchema={jsonSchema}
        fields={fields}
      />
    </RJSFWrapperComponent>
  );
}

export default RJSF;

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
    fields,
    ArrayFieldTemplate = MesheryArrayFieldTemplate,
    ObjectFieldTemplate = MesheryCustomObjFieldTemplate,
  } = props;

  return (
    <MuiThemeProvider theme={rjsfTheme}>
      <Form
        schema={schema.rjsfSchema}
        idPrefix={jsonSchema?.title}
        onChange={onChange}
        formData={data}
        fields={fields}
        ArrayFieldTemplate={ArrayFieldTemplate}
        ObjectFieldTemplate={ObjectFieldTemplate}
        additionalMetaSchemas={[JS4]}
        uiSchema={schema.uiSchema}
        liveValidate
        showErrorList={false}
        noHtml5Validate
      >
        {/* {hideSubmit ? true : <RJSFButton handler={onSubmit} text="Submit" {...restparams} />}
{hideSubmit ? true : <RJSFButton handler={onDelete} text="Delete" />} */}
        {/* <RJSFFormChildComponent /> */}
        <></> {/* temporary change for functionality */}
      </Form>

    </MuiThemeProvider>
  )
}
