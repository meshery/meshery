import { Button } from '@material-ui/core';
import PublicIcon from '@material-ui/icons/Public';
import React from 'react';
import useStyles from "../MesheryPatterns/Cards.styles";
import Modal from './Modal';
import { publish_schema } from './schemas/publish_schema';

export default function PublishModal(props) {
  const { open, title, handleClose, handlePublish, formData, onChange, payload } = props;
  const classes = useStyles();

  return (
    <Modal open={open} schema={publish_schema} title={title} onChange={onChange} handleClose={handleClose} formData={formData}>
      <Button
        title="Publish"
        variant="contained"
        color="primary"
        className={classes.testsButton}
        onClick={() => {
          handleClose();
          handlePublish(payload)
        }}
      >
        <PublicIcon className={classes.iconPatt} />
        <span className={classes.btnText}> Submit For Approval </span>
      </Button>
    </Modal>
  )
}