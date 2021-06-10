import React from "react";
import { withTheme } from "@rjsf/core";
import { Theme as MaterialUITheme } from "@rjsf/material-ui";

const Form = withTheme(MaterialUITheme);

function deleteTitleFromJSONSchema(jsonSchema) {
  return { ...jsonSchema, title: "" };
}

function RJSF({ jsonSchema, onChange, hideSubmit, hideTitle }) {
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
      {hideSubmit ? true : null}
    </Form>
  );
}

export default RJSF;
