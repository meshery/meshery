import React from "react";
import { withTheme } from "@rjsf/core";
import { Theme as MaterialUITheme } from "@rjsf/material-ui";
import { Button } from "@material-ui/core";
import { MuiThemeProvider, createTheme } from '@material-ui/core/styles';
import JS4 from "../../../assets/jsonschema/schema-04.json";
import ArrayFieldTemplate from "./RJSFCustomComponents/ArrayFieldTemlate"

const Form = withTheme(MaterialUITheme);

const muiTheme = createTheme({
  palette : {
    primary : {
      main : '#00b39f',
    },
  },
  props : {
    MuiTextField : {
      variant : 'outlined',
      margin : 'dense',
    },
  },
});

function deleteTitleFromJSONSchema(jsonSchema) {
  return { ...jsonSchema, title : "" };
}

function RJSFButton({ handler, text }) {
  return (
    <Button variant="contained" color="primary" style={{ marginRight : "0.5rem" }} onClick={handler}>
      {text}
    </Button>
  );
}


const uiSchema = {
  replicas : {
    "ui:widget" : "range"
  }
};

function RJSF({ jsonSchema, onChange, hideSubmit, hideTitle, onSubmit, onDelete, renderAsTooltip }) {

  const [data, setData] = React.useState();

  React.useEffect(() => {
    onChange?.(data);
  }, [data]);

  return (
    <>
      {!renderAsTooltip ? (
        <Form
          schema={hideTitle ? deleteTitleFromJSONSchema(jsonSchema) : jsonSchema}
          idPrefix={jsonSchema?.title}
          onChange={(e) => setData(e.formData)}
          formData={data}
          liveValidate
          additionalMetaSchemas={[JS4]}
          // noHtml5Validate
        >
          {hideSubmit ? true : <RJSFButton handler={onSubmit} text="Submit" />}
          {hideSubmit ? true : <RJSFButton handler={onDelete} text="Delete" />}
        </Form>) : (
        <MuiThemeProvider theme={muiTheme}>
          <Form
            schema={hideTitle ? deleteTitleFromJSONSchema(jsonSchema) : jsonSchema}
            idPrefix={jsonSchema?.title}
            onChange={(e) => setData(e.formData)}
            formData={data}
            liveValidate
            showErrorList={false}
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
