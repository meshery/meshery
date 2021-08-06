import React from "react";
import { withTheme } from "@rjsf/core";
import { Theme as MaterialUITheme } from "@rjsf/material-ui";
import { Button, IconButton } from "@material-ui/core";
import { Add } from "@material-ui/icons";
import JS4 from "../../../assets/jsonschema/schema-04.json";

const Form = withTheme(MaterialUITheme);

function deleteTitleFromJSONSchema(jsonSchema) {
  return { ...jsonSchema, title: "" };
}

function RJSFButton({ handler, text }) {
  return (
    <Button variant="contained" color="primary" style={{ marginRight: "0.5rem" }} onClick={handler}>
      {text}
    </Button>
  );
}

function ArrayFieldTemplate(props) {
  return (
    <div>
      {props.items.map(element => element.children)}
      {props.canAdd && <IconButton style={{ float: "right", margin: 0 }} onClick={props.onAddClick}><Add /></IconButton>}
    </div>
  );
}

const uiSchema = {
  replicas: {
    "ui:widget" : "range"
  }
};

// function CustomFieldTemplate(props) {
//   const {id, classNames, label, help, required, description, errors, children} = props;
//   console.log(label, props);
//   return (
//     <div className={classNames}>
//       <label htmlFor={id}>{label}{required ? "*" : null}</label>
//       {description}
//       {children}
//       {errors}
//       {help}
//     </div>
//   );
// }


function RJSF({ jsonSchema, onChange, hideSubmit, hideTitle, onSubmit, onDelete, isMeshery }) {
  const [data, setData] = React.useState();

  React.useEffect(() => {
    onChange?.(data);
  }, [data]);

  return (
    <>
      {isMeshery ? (
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
        <Form
          schema={hideTitle ? deleteTitleFromJSONSchema(jsonSchema) : jsonSchema}
          idPrefix={jsonSchema?.title}
          onChange={(e) => setData(e.formData)}
          formData={data}
          liveValidate
          showErrorList={false}
          additionalMetaSchemas={[JS4]}
          ArrayFieldTemplate={ArrayFieldTemplate}
          uiSchema={uiSchema}
          // noHtml5Validate
        >
          {hideSubmit ? true : <RJSFButton handler={onDelete} text="Delete" />}
        </Form>
      )}
    </>
  );
}

export default RJSF;
