import React from 'react';
import Modal from '../Modal';
import PublicIcon from '@material-ui/icons/Public';
import { PublishSchema, PublishUiSchema } from '@layer5/sistent';

// This modal is used in MeshMap also
export default function PublishModal(props) {
  const { open, title, handleClose, handleSubmit } = props;
  return (
    <Modal
      open={open}
      schema={PublishSchema}
      uiSchema={PublishUiSchema}
      title={title}
      handleClose={handleClose}
      handleSubmit={handleSubmit}
      submitBtnText="Submit for Approval"
      submitBtnIcon={<PublicIcon data-cy="import-button" />}
      showInfoIcon={{
        text: 'Upon submitting your catalog item, an approval flow will be initiated.',
        link: 'https://docs.meshery.io/concepts/catalog',
      }}
    />
  );
}
