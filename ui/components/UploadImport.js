import React, { useEffect } from 'react'
import LinkIcon from '@material-ui/icons/Link';
import { TextField, Button, Grid } from '@material-ui/core';
import { makeStyles, MuiThemeProvider } from '@material-ui/core/styles';
import { createTheme } from '@material-ui/core/styles';
import { URLValidator } from '../utils/URLValidator';
import { Dialog,  DialogActions,
  DialogContent,
  DialogTitle } from '@material-ui/core';





const getMuiTheme = () => createTheme({
  palette : {
    primary : {
      main : "#607d8b"
    },
    secondary : {
      main : "#666666"
    }
  },
  overrides : {
    MuiGrid : {
      input : {
        color : '#607d8b'
      }
    },
  }
})


const styles = makeStyles(() => ({
  upload : {
    paddingLeft : "0.7rem",
    paddingTop : "8px"
  },
  title : {
    textAlign : 'center',
    minWidth : 500,
    padding : '15px',
    color : '#fff',
    backgroundColor : '#607d8b'
  },
  content : {
    padding : "1.5rem"
  },
  heading : {
    color : "#607d8b"
  },


}));



const UploadImport = ({ handleUpload, handleImport, configuration, modalStatus }) => {
  const classes = styles();
  const [open, setOpen] = React.useState(false);
  const [input, setInput] = React.useState();
  const [isError, setIsError] = React.useState(false);


  useEffect(() => {
    if (input) {
      setIsError(!URLValidator(input))
    }
  }, [input])

  useEffect(() => {
    if (modalStatus) {
      handleClose()
    }
  }, [modalStatus])

  const handleOpen = () => {
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
  };

  const handleSubmit = () => {
    handleUpload(input)
    handleClose()
  }

  //   const handleUploader = () => {
  //     handleImport(input)
  //     handleClose()
  //   }

  return (
    <>
      <label htmlFor="url-upload-button">

        <Button aria-label="URL-Upload" data-cy="import-button" variant="contained"
          color="primary" className={classes.button}
          size="large" onClick={handleOpen}>
          <LinkIcon style={{ padding : "1px" }} />
           Import {configuration}
        </Button>

        <Dialog
          open={open}
          handleClose={handleClose}>

          <MuiThemeProvider theme={getMuiTheme()}>
            <DialogTitle className={classes.title}>
              <b id="simple-modal-title" style={{ textAlign : "center" }} >Import {configuration}</b>
            </DialogTitle>
            <DialogContent className={classes.content}>
              <Grid container spacing={24}>
                <Grid item xs={3}>
                  <h4 className={classes.heading} > FROM URL </h4>
                </Grid>
                <Grid item xs={9}>
                  <TextField
                    error={isError}
                    helperText={isError && "Invalid URL"}
                    variant="outlined"
                    label={`URL for ${configuration}`}
                    style={{ width : "100%" }}
                    onChange={(e) => setInput(e.target.value)} />
                </Grid>
              </Grid>
              <hr/>
              <Grid container spacing={24}>
                <Grid item xs={3}>
                  <h4 className={classes.heading}> UPLOAD FILE </h4>
                </Grid>
                <Grid item xs={6}>
                  <TextField
                    variant="outlined"
                    label="Filename"
                    style={{ width : "100%" }}
                  />
                </Grid>
                <Grid item xs={3}>
                  <label htmlFor="upload-button" className={classes.upload}>

                    <Button variant="contained" size="large"  color="primary" aria-label="Upload Button"  component="span" >
                      <input id="upload-button" type="file"  accept=".yaml, .yml" hidden onChange={ handleImport } name="upload-button" data-cy="file-upload-button" />
                        Browse
                    </Button>
                  </label>
                </Grid>
              </Grid>
            </DialogContent>
            <DialogActions>
              <label htmlFor="cancel"  className={classes.cancel}>
                <Button variant="outlined" size="large" color="secondary" onClick={handleClose}>Cancel</Button>
              </label>
              <label htmlFor="URL">  <Button disabled={isError || !input} size="large" id="URL" variant="contained" color="primary" onClick={(e) => handleSubmit(e, handleUpload)}>Import</Button> </label>

            </DialogActions>
          </MuiThemeProvider>
        </Dialog>

      </label>
    </>
  )
}

export default UploadImport
