import React from "react";
import { withTheme } from "@rjsf/core";
import { Theme as MaterialUITheme } from "@rjsf/material-ui";
import { Button } from "@material-ui/core";

const Form = withTheme(MaterialUITheme);

function deleteTitleFromJSONSchema(jsonSchema) {
  return { ...jsonSchema, title: "" };
}

function RJSFSubmit({ onSubmit }) {
  return (
    <Button variant="contained" color="primary" onClick={() => onSubmit?.()} >
      Submit
    </Button>
  );
}

function RJSF({ jsonSchema, onChange, hideSubmit, hideTitle, onSubmit }) {
  const [data, setData] = React.useState();

  React.useEffect(() => {
    onChange?.(data);
  }, [data]);

  return (
    <Form
      schema={hideTitle ? deleteTitleFromJSONSchema(jsonSchema) : jsonSchema}
      idPrefix={jsonSchema?.title}
      onChange={(e) => setData(e.formData)}
      formData={data}
      liveValidate
      // noHtml5Validate
    >
      {hideSubmit ? true : <RJSFSubmit onSubmit={onSubmit} />}
    </Form>
  );
}

export default RJSF;
