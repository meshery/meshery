import React, { useEffect } from 'react';
import LinkIcon from '@material-ui/icons/Link';
import { Tooltip, IconButton, TextField, Button, Grid } from '@layer5/sistent-components';
import { withStyles } from '@material-ui/core/styles';
import { URLValidator } from '../utils/URLValidator';

// Remove this as this won't be needed in the MUI v5 migration
// import { createTheme } from '@material-ui/core/styles';
// const getMuiTheme = () => createTheme({
//   palette : {
//     primary : {
//       main : "#607d8b"
//     }
//   },
//   overrides : {
//     MuiGrid : {
//       input : {
//         color : '#607d8b'
//       }
//     },
//   }
// })

const styles = (theme) => ({
  paper: {
    position: 'absolute',
    width: 600,
    backgroundColor: theme.palette.background.paper,
    border: '0px solid #000',
    boxShadow: theme.shadows[5],
    padding: theme.spacing(2, 4, 3),
    top: '50%',
    left: '50%',
    transform: `translate(-50%, -50%)`,
    borderRadius: 10,
  },
  grid: {
    width: '100%',
  },
});
import GenericModal from './GenericModal';

const URLUploader = ({ onSubmit, classes }) => {
  const [open, setOpen] = React.useState(false);
  const [input, setInput] = React.useState();
  const [isError, setIsError] = React.useState(false);

  useEffect(() => {
    if (input) {
      setIsError(!URLValidator(input));
    }
  }, [input]);

  const handleOpen = () => {
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
  };

  const handleSubmit = () => {
    onSubmit(input);
    handleClose();
  };

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
            <div className={classes.paper}>
              <Grid container spacing={2}>
                <Grid item xs={12}>
                  <h2 id="simple-modal-title">Import using URL</h2>
                </Grid>
                <Grid item xs={12}>
                  <TextField
                    error={isError}
                    helperText={isError && 'Invalid URL'}
                    variant="outlined"
                    label="Paste URL here"
                    fullWidth
                    onChange={(e) => setInput(e.target.value)}
                  />
                </Grid>
                <Grid item xs={12} sm={6}>
                  <Button fullWidth variant="contained" onClick={handleClose}>
                    Cancel
                  </Button>
                </Grid>
                <Grid item xs={12} sm={6}>
                  <Button
                    disabled={isError || !input}
                    fullWidth
                    variant="contained"
                    color="primary"
                    onClick={() => handleSubmit()}
                  >
                    Import
                  </Button>
                </Grid>
              </Grid>
            </div>
          }
        />
      </label>
    </>
  );
};

export default withStyles(styles)(URLUploader);
