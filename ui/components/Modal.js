import React from 'react'
import { Button, Grid, IconButton } from '@material-ui/core';
import {
  Dialog, DialogActions,
  DialogContent,
  DialogTitle
} from '@material-ui/core';
import useStyles from "./MesheryPatterns/Cards.styles";
import PublicIcon from '@material-ui/icons/Public';
import CloseIcon from '@material-ui/icons/Close';
import RJSFWrapper from './MesheryMeshInterface/PatternService/RJSF_wrapper';

function Modal(props) {
  const { open, handleClose, pattern, handlePublish, schema } = props;
  const classes = useStyles();

  const [data, setData] = React.useState(null)
  const [payload, setPayload] = React.useState({
    "id" : pattern?.id,
    "catalog_data" : pattern?.catalog_data
  })
  React.useEffect(() => {
    setData(pattern.catalog_data)
  }, [pattern])
  React.useEffect(() => {
    setPayload({
      "id" : pattern?.id,
      "catalog_data" : data
    })
    console.log(payload)
  }, [data])

  return (
    <>
      <Dialog
        open={open}
        onClose={handleClose}>
        <DialogTitle>
          <div className={classes.publishTitle}>

            <b id="simple-modal-title" style={{ textAlign : "center" }} > {pattern?.name}</b>
            <IconButton aria-label="close" className={classes.closeButton} onClick={handleClose}>
              <CloseIcon />
            </IconButton>
          </div>
        </DialogTitle>
        <DialogContent>
          <Grid container spacing={24} alignItems="center">
            <RJSFWrapper
              formData={data}
              jsonSchema={schema}
              onChange={(e) => {
                console.log(e);
              }}
              hideTitle={true}
            />
          </Grid>

        </DialogContent>
        <DialogActions>
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
            <span className={classes.btnText}> Publish </span>
          </Button>
        </DialogActions>

      </Dialog>
    </>
  )
}

export default Modal;