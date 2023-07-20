import React, { useEffect } from 'react'
import { TextField, Button, Grid, NativeSelect } from '@material-ui/core';
import { createTheme, MuiThemeProvider, useTheme, withStyles } from '@material-ui/core/styles';
import { URLValidator } from '../utils/URLValidator';
import {
  Dialog, DialogActions,
  DialogContent,
  DialogTitle
} from '@material-ui/core';
import { promisifiedDataFetch } from '../lib/data-fetch';

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
});

const getDarkMuiTheme = () => createTheme({
  palette : {
    type : "dark",
    primary : {
      main : "#607d8b"
    },
  },
  overrides : {
    MuiGrid : {
      input : {
        color : '#607d8b'
      }
    },
    MuiFormLabel : {
      root : {
        "&$focused" : {
          color : "#00B39F",
        },
      }
    },
  }
});

const styles = (theme) => ({
  upload : {
    paddingLeft : "0.7rem",
    paddingTop : "8px"
  },
  title : {
    textAlign : 'center',
    minWidth : 500,
    padding : '15px',
    color : '#fff',
    backgroundColor : theme.palette.type === "dark" ? "#202020" : '#396679'
  },
  heading : {
    color : theme.palette.type === "dark" ? "#fff" : "#607d8b"
  },
  selectType : {
    color : theme.palette.type === "dark" ? "#fff" : "#607d8b",
    marginRight : "1.2rem"
  },
  button : {
    backgroundColor : theme.palette.type === "dark" ? "#00B39F" : "#607d8b",
    "&:hover" : {
      backgroundColor : theme.palette.type === "dark" ? "#00B39F" : "#607d8b"
    },
    color : "#fff"
  },
});

function UploadImport(props) {
  const { handleUpload, handleUrlUpload, configuration, isApplication, open, handleClose, classes, fetch } = props;
  const [input, setInput] = React.useState();
  const [isError, setIsError] = React.useState(false);
  const [fileType, setFileType] = React.useState();
  const [sourceType, setSourceType] = React.useState();
  const [supportedTypes, setSupportedTypes] = React.useState();
  const theme=useTheme()
  useEffect(() => {
    if (isApplication) {
      (async () => {
        setSupportedTypes(await promisifiedDataFetch("/api/application/types"));
      })()
    }
  }, []);

  const handleFileType = (index) => {
    if (isApplication) {
      setFileType(supportedTypes?.[index]?.supported_extensions);
      setSourceType(supportedTypes?.[index]?.application_type);
    }
  }

  useEffect(() => {
    if (input) {
      setIsError(!URLValidator(input))
    }
  }, [input])

  useEffect(() => {
    if (isApplication) {
      setFileType(supportedTypes?.[0]?.supported_extensions)
      setSourceType(supportedTypes?.[0]?.application_type);
    }
  }, [open])

  const handleSubmit = async() => {
    await handleUrlUpload(input, sourceType)
    handleClose()
  }

  const handleUploader = async(input) => {
    await handleUpload(input, sourceType)
    handleClose()
  }


  return (
    <>
      <label htmlFor="url-upload-button">
        <Dialog
          open={open}
          onClose={handleClose}>

          <MuiThemeProvider theme={theme.palette.type == "dark" ? getDarkMuiTheme() : getMuiTheme()}>
            <DialogTitle className={classes.title}>
              <b id="simple-modal-title" style={{ textAlign : "center" }} >Import {configuration}</b>
            </DialogTitle>
            <DialogContent>
              <Grid container spacing={24} alignItems="center">
                <Grid item xs={3}>
                  <h4 className={classes.heading} > FROM URL </h4>
                </Grid>
                <Grid item xs={9}>
                  <TextField
                    size="small"
                    error={isError}
                    helperText={isError && "Invalid URL"}
                    variant="outlined"
                    label={`URL for ${configuration}`}
                    style={{ width : "100%" }}
                    onChange={(e) => setInput(e.target.value)} />
                </Grid>
              </Grid>
              <hr />
              {
                sourceType !== "Helm Chart" && (
                  <Grid container spacing={24} alignItems="center">
                    <Grid item xs={3}>
                      <h4 className={classes.heading}> UPLOAD FILE </h4>
                    </Grid>
                    <Grid item xs={6}>
                      <TextField
                        size="small"
                        variant="outlined"
                        label="Filename"
                        style={{ width : "100%" }}
                      />
                    </Grid>
                    <Grid item xs={3}>

                      <label htmlFor="upload-button" className={classes.upload}>

                        <Button disabled={sourceType === "Helm Chart"} variant="contained" className={classes.button} aria-label="Upload Button" onChange={sourceType === "Helm Chart" ? null : handleUploader} component="span" >
                          <input id="upload-button" type="file" accept={fileType} disabled={sourceType === "Helm Chart"} hidden name="upload-button" data-cy="file-upload-button" />
                        Browse
                        </Button>
                      </label>
                    </Grid>
                  </Grid>
                )
              }

              <Grid container spacing={24} alignItems="center">
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
                      supportedTypes?.map((type, index) => (
                        <option key={index} value={index}>
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
                <Button variant="outlined" color="secondary" onClick={handleClose}>Cancel</Button>
              </label>
              <label htmlFor="URL">  <Button disabled={isError || !input} id="URL" variant="contained" className={classes.button} onClick={async(e) => {
                await handleSubmit(e, handleUploader);
                fetch?.();
              }}>Import</Button> </label>
            </DialogActions>
          </MuiThemeProvider>
        </Dialog>

      </label>
    </>
  )
}

export default withStyles(styles)(UploadImport);