import React, {  useRef } from 'react';
import filterSchema from "../../assets/jsonschema/filterImport.json";
import applicationSchema from "../../assets/jsonschema/applicationImport.json";
import applicationUISchema from "../../assets/jsonschema/uiSchemaApplication.json";
import designSchema from "../../assets/jsonschema/designImport.json"
import useStyles from "../MesheryPatterns/Cards.styles";
import RJSFModal from "../Modal";
import { Button } from '@material-ui/core';
import { capitalize } from 'lodash';

const schemaMap = {
  filter : filterSchema,
  application : applicationSchema,
  design : designSchema
}

const uiSchema = {
  application : applicationUISchema
}

export default function ImportModal(props) {
  const { importType, handleSubmit, handleClose } = props;
  const schema = schemaMap[importType];
  const classes = useStyles();
  const data = useRef(null);

  return (
    <RJSFModal
      {...props}
      title={`Import ${capitalize(importType)}`}
      // onChange={()=> console}
      onChange={(newData) => data.current = newData}
      schema={schema}
      formData={{}}
      type={importType}
      onSubmit={(d) => console.log("onsubmit called", d)}
      uiSchema={uiSchema?.[importType]}
    >
      <Button
        fullWidth
        title="Publish"
        variant="contained"
        color="primary"
        className={classes.testsButton}
        onClick={() => {
          handleClose();
          handleSubmit(data.current);
        }}
      >
        Import {capitalize(importType)}
      </Button>
    </RJSFModal>
  )
}