import React from "react";
import { withTheme } from "@rjsf/core";
import { Theme as MaterialUITheme } from "@rjsf/material-ui";
import { Button } from "@material-ui/core";

const Form = withTheme(MaterialUITheme);

function deleteTitleFromJSONSchema(jsonSchema) {
  return { ...jsonSchema, title: "" };
}

function RJSFButton({ handler, text }) {
  return (
    <Button variant="contained" color="primary" style={{ marginRight: "0.5rem" }} onClick={() => handler?.()} >
      {text}
    </Button>
  );
}

function RJSF({ jsonSchema, onChange, hideSubmit, hideTitle, onSubmit, onDelete }) {
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
      {hideSubmit ? true : <RJSFButton handler={onSubmit} text="Submit" />}
      {hideSubmit ? true : <RJSFButton handler={onDelete} text="Delete" />}
    </Form>
  );
}

export default RJSF;
