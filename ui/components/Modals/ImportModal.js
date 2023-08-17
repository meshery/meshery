import React, {  useRef } from 'react';
import useStyles from "../MesheryPatterns/Cards.styles";
import RJSFModal from "../Modal";
import { Button } from '@material-ui/core';
import { capitalize } from 'lodash';

// might also expect RJSFWrapperComponent from extensions
export default function ImportModal(props) {
  const { importType, handleSubmit, handleClose, rjsfSchema, uiSchema } = props;
  const classes = useStyles();
  const data = useRef(null);

  return (
    <RJSFModal
      {...props}
      title={`Import ${capitalize(importType)}`}
      onChange={null}
      schema={rjsfSchema}
      formData={{}}
      type={importType}
      uiSchema={uiSchema || {}}
      submitBtnText={`Import ${capitalize(importType)}`}
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