import React, { useEffect } from 'react'
import LinkIcon from '@material-ui/icons/Link';
import { TextField, Button, Grid, NativeSelect } from '@material-ui/core';
import { makeStyles, MuiThemeProvider } from '@material-ui/core/styles';
import { createTheme } from '@material-ui/core/styles';
import { URLValidator } from '../utils/URLValidator';
import {
  Dialog, DialogActions,
  DialogContent,
  DialogTitle
} from '@material-ui/core';

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
  selectType : {
    color : "#607d8b",
    marginRight : "1.2rem"
  }


}));

const UploadImport = ({ handleUpload, handleImport, configuration, isApplication, supportedTypes }) => {
  const classes = styles();
  const [open, setOpen] = React.useState(false);
  const [input, setInput] = React.useState();
  const [isError, setIsError] = React.useState(false);
  const [fileType, setFileType] = React.useState();
  const [sourceType, setSourceType] = React.useState();

  const handleFileType = (index) => {
    setFileType(supportedTypes[index]?.supported_extensions);
    setSourceType(supportedTypes[index]?.application_type);
  }
  useEffect(() => {
    if (input) {
      setIsError(!URLValidator(input))
    }
  }, [input])

  useEffect(() => {
    setFileType(supportedTypes[0]?.supported_extensions)
    setSourceType(supportedTypes[0]?.application_type);
  },[open])
  const handleOpen = () => {
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
  };

  const handleSubmit = () => {
    handleUpload(input, sourceType)
    handleClose()
  }

  const handleUploader = (input) => {
    handleImport(input, sourceType)
    handleClose()
  }


  return (
    <>
      <label htmlFor="url-upload-button">

        <Button aria-label="URL-Upload" data-cy="import-button" variant="contained"
          color="primary" className={classes.button}
          size="large" onClick={handleOpen}>
          <LinkIcon style={{ paddingRight : ".35rem" }} />
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
              <hr />
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

                    <Button disabled={sourceType==="Helm Chart"} variant="contained" size="large" color="primary" aria-label="Upload Button" onChange={sourceType==="Helm Chart" ? null : handleUploader} component="span" >
                      <input id="upload-button" type="file" accept={fileType} disabled={sourceType==="Helm Chart"} hidden  name="upload-button" data-cy="file-upload-button" />
                      Browse
                    </Button>
                  </label>
                </Grid>
              </Grid>
              <Grid container spacing={24}>
                {
                  isApplication && 
                <h4 className={classes.selectType}>SELECT TYPE </h4>
                }
                {isApplication &&
                  <>
                    <NativeSelect
                      defaultValue={0}
                      onChange={(e) => handleFileType(e.target.value)}
                      inputProps={{
                        name : 'name',
                        id : 'uncontrolled-native',
                      }}
                    >
                      {
                        supportedTypes.map((type, index) => (
                          <option value={index}>
                            {type.application_type}
                          </option>
                        ))
                      }
                    </NativeSelect>
                  </>
                }
              </Grid>
            </DialogContent>
            <DialogActions>
              <label htmlFor="cancel" className={classes.cancel}>
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
