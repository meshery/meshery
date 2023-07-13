import React from 'react';
import filterSchema from "../../assets/jsonschema/filterImport.json";
import applicationSchema from "../../assets/jsonschema/applicationImport.json";
import applicationUISchema from "../../assets/jsonschema/uiSchemaApplication.json";
import useStyles from "../MesheryPatterns/Cards.styles";
import RJSFModal from "../Modal";
import { Button } from '@material-ui/core';
import { capitalize } from 'lodash';

const schemaMap = {
  filter: filterSchema,
  application: applicationSchema,
  design: {}
}

const uiSchema = {
  application: applicationUISchema
}

export default function ImportModal(props) {
  const { importType, open, handleClose, handleImport } = props;
  const schema = schemaMap[importType];
  const classes = useStyles();

  const data = null;
  return (
    <RJSFModal
      {...props}
      title={`Import ${capitalize(importType)}`}
      onChange={(e) => console.log("oncgan", e)}
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
          handlePublishModalClose();
          handlePublish(payload);
        }}
      >
        Import {capitalize(importType)}
      </Button>
    </RJSFModal>
  )
}