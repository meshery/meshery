import React from 'react'
import LinkIcon from '@material-ui/icons/Link';
import { Tooltip, IconButton, TextField,Button, Grid } from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import { URLValidator } from '../utils/URLValidator';

const styles = makeStyles((theme) => ({
  paper : {
    position : 'absolute',
    width : 600,
    backgroundColor : theme.palette.background.paper,
    border : '0px solid #000',
    boxShadow : theme.shadows[5],
    padding : theme.spacing(2, 4, 3),
    top : '50%',
    left : '50%',
    transform : `translate(-50%, -50%)`,
    borderRadius : 10,
  },
  grid : {
    width : '100%'
  }
}));
import GenericModal from "./GenericModal";

const URLUploader = ({ onSubmit }) => {
  const classes = styles();
  const [open, setOpen] = React.useState(false);
  const [input, setInput] = React.useState('');
  const handleOpen = () => {
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
  };
  const validURL = (str) =>  {
    return URLValidator(str);
  }
  const handleError = (input) => {
    console.log(input + ' is not valid url')
  }
  const handleSubmit = () => {
    validURL(input) ? onSubmit(input) : handleError(input);
    handleClose()
  }

  return (
    <>
      <label htmlFor="url-upload-button">

        <Tooltip title="Upload URL">
          <IconButton aria-label="URL-Upload" component="span" onClick={handleOpen}>
            <LinkIcon />
          </IconButton>
        </Tooltip>
        <GenericModal
          open={open}
          handleClose={handleClose}
          Content={
            <div
              className={classes.paper}
            >
              <Grid
                container spacing={2}
              >
                <Grid
                  item xs={12}>
                  <h2 id="simple-modal-title">Import using URL</h2>
                </Grid>
                <Grid
                  item xs={12}>
                  <TextField id="standard-basic" label="Paste URL here" fullWidth onChange={(e) => setInput(e.target.value) }/>
                </Grid>
                <Grid item xs={12} sm={6}>
                  <Button fullWidth variant="contained" onClick={handleClose}>Cancel</Button>
                </Grid>
                <Grid item xs={12} sm={6}>
                  <Button fullWidth variant="contained" color="primary" onClick={() => handleSubmit()}>Insert</Button>
                </Grid>
              </Grid>

            </div>
          }
        />
      </label>
    </>
  )
}

export default URLUploader
