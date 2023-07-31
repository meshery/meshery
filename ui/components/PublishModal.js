import React from 'react';
import { publish_schema, publish_ui_schema } from './schemas/publish_schema';
import Modal from './Modal';

export default function PublishModal(props) {
  const { open, title, handleClose, formData, onChange } = props;

  return (
    <Modal open={open} schema={publish_schema} uiSchema={publish_ui_schema} title={title} onChange={onChange} handleClose={handleClose} formData={formData} showInfoIcon={{ text : "Upon submitting your catalog item, an approval flow will be initiated.", link : "https://docs.meshery.io/concepts/catalog" }}/>
  )
}