import React, { useRef } from 'react';
import RJSFModal from '../Modal';
import { Button } from '@layer5/sistent';
import { capitalize } from 'lodash';

// might also expect RJSFWrapperComponent from extensions
// NOTE: Use modal from sistent
export default function ImportModal(props) {
  const { importType, handleSubmit, handleClose, rjsfSchema, uiSchema } = props;
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
        onClick={() => {
          handleClose();
          handleSubmit(data.current);
        }}
      >
        Import {capitalize(importType)}
      </Button>
    </RJSFModal>
  );
}
