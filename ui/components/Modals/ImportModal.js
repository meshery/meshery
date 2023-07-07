import React from 'react';
import { useEffect } from 'react';
import { Button, Grid } from '@material-ui/core';
// import { createTheme } from '@material-ui/core/styles';
import validator from "@rjsf/validator-ajv8";
import {
  Dialog, DialogActions,
  DialogContent,
  DialogTitle
} from '@material-ui/core';
import { Form } from '@rjsf/material-ui';
import useStyles from "../MesheryPatterns/Cards.styles";
import PublicIcon from '@material-ui/icons/Public';
import filterSchema from "../../assets/jsonschema/filterImport.json"

const schemaMap = {
  filter: filterSchema,
  application: {},
  design: {}
}

export default function ImportModal(props) {
  const { importType, open, handleClose, handleImport } = props;
  const classes = useStyles();
  const schema = schemaMap[importType];

  const [data, setData] = React.useState(null);

  return (
    <>
      <Dialog
        open={open}
        onClose={handleClose}>
        <DialogTitle className={classes.dialogTitle}>
          Import your {importType}
        </DialogTitle>
        <DialogContent>
          <Grid container spacing={24} alignItems="center">
            <Form schema={schema} formData={data} validator={validator}
              onSubmit={(data) => {
                handleClose();
                handleImport(data)
              }}
            >
              <Button
                title="Import"
                variant="contained"
                color="primary"
                type="submit"
                className={classes.testsButton}
              >
                <PublicIcon className={classes.iconPatt} />
                <span className={classes.btnText}> Import {importType} </span>
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