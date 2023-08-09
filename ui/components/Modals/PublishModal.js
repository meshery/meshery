import React from "react";
import { publish_schema, publish_ui_schema } from "../schemas/publish_schema";
import Modal from "../Modal";
import PublicIcon from '@material-ui/icons/Public';

// This modal is used in MeshMap also
export default function PublishModal(props) {
  const { open, title, handleClose, handleSubmit } = props;

  return (
    <Modal
      open={open}
      schema={publish_schema}
      uiSchema={publish_ui_schema}
      title={title}
      handleClose={handleClose}
      handleSubmit={handleSubmit}
      submitBtnText="Submit for Approval"
      submitBtnIcon={<PublicIcon  style={iconMedium} className={classes.addIcon} data-cy="import-button"/>}
      showInfoIcon={{
        text : "Upon submitting your catalog item, an approval flow will be initiated.",
        link : "https://docs.meshery.io/concepts/catalog",
      }}
    />
  );
}
