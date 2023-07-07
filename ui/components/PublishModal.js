import React from 'react';
import { Button, Grid } from '@material-ui/core';
import validator from "@rjsf/validator-ajv8";
import {
  Dialog, DialogActions,
  DialogContent,
  DialogTitle
} from '@material-ui/core';
import { Form } from '@rjsf/material-ui';
import useStyles from "./MesheryPatterns/Cards.styles";
import PublicIcon from '@material-ui/icons/Public';
import { publishSchema } from './schemas/connections/publishSchema';

export default function PublishModal(props) {
  const { open, handleClose, resourceType, handlePublish, title } = props;
  const classes = useStyles();

  const data = resourceType?.catalog_data || {};

  return (
    <>
      <Dialog
        open={open}
        onClose={handleClose}>
        <DialogTitle className={classes.dialogTitle}>
         Request To Publish: {title}
        </DialogTitle>
        <DialogContent>
          <Grid container spacing={24} alignItems="center">
            <Form schema={publishSchema} formData={data} validator={validator}
              onSubmit={(data) => {
                handlePublish({
                  id : resourceType.id,
                  catalog_data : data.formData
                })
                handleClose();
              }}
            >
              <Button
                title="Publish"
                variant="contained"
                color="primary"
                type="submit"
                className={classes.testsButton}
              >
                <PublicIcon className={classes.iconPatt} />
                <span className={classes.btnText}> Submit for Approval </span>
              </Button>
            </Form>
          </Grid>

        </DialogContent>
        <DialogActions>

        </DialogActions>

      </Dialog>
    </>
  )
}