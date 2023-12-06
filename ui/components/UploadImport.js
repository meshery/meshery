import React, { useEffect } from 'react';
import { TextField, Button, Grid, NativeSelect, Divider, Typography } from '@material-ui/core';
import { createTheme, MuiThemeProvider, useTheme, withStyles } from '@material-ui/core/styles';
import { URLValidator } from '../utils/URLValidator';
import { Dialog, DialogActions, DialogContent, DialogTitle } from '@material-ui/core';
import { promisifiedDataFetch } from '../lib/data-fetch';

const getMuiTheme = () =>
  createTheme({
    palette: {
      primary: {
        main: '#607d8b',
      },
      secondary: {
        main: '#666666',
      },
    },
    overrides: {
      MuiGrid: {
        input: {
          color: '#607d8b',
        },
      },
    },
  });

const getDarkMuiTheme = () =>
  createTheme({
    palette: {
      type: 'dark',
      primary: {
        main: '#607d8b',
      },
    },
    overrides: {
      MuiGrid: {
        input: {
          color: '#607d8b',
        },
      },
      MuiFormLabel: {
        root: {
          '&$focused': {
            color: '#00B39F',
          },
        },
      },
    },
  });

const styles = (theme) => ({
  upload: {
    paddingLeft: '0.7rem',
    paddingTop: '8px',
  },
  title: {
    textAlign: 'center',
    minWidth: 500,
    padding: '15px',
    color: '#fff',
    backgroundColor: theme.palette.type === 'dark' ? '#202020' : '#396679',
  },
  heading: {
    color: theme.palette.type === 'dark' ? '#fff' : '#607d8b',
  },
  selectType: {
    color: theme.palette.type === 'dark' ? '#fff' : '#607d8b',
    marginRight: '1.2rem',
  },
  button: {
    backgroundColor: theme.palette.type === 'dark' ? '#00B39F' : '#607d8b',
    '&:hover': {
      backgroundColor: theme.palette.type === 'dark' ? '#00B39F' : '#607d8b',
    },
    color: '#fff',
  },
});

function UploadImport(props) {
  const {
    handleUpload,
    handleUrlUpload,
    configuration,
    isFilter,
    open,
    handleClose,
    classes,
    fetch,
  } = props; // hack, to remove, please............
  const [input, setInput] = React.useState();
  const [name, setName] = React.useState('');
  const [config, setConfig] = React.useState('');
  const [isError, setIsError] = React.useState(false);
  const [fileType, setFileType] = React.useState();
  const [sourceType, setSourceType] = React.useState();
  const [supportedTypes, setSupportedTypes] = React.useState();
  const theme = useTheme();
  const isDesign = configuration === 'patterns';
  useEffect(() => {
    if (isDesign) {
      (async () => {
        setSupportedTypes(await promisifiedDataFetch('/api/pattern/types'));
      })();
    }
  }, []);

  const handleFileType = (index) => {
    if (isDesign) {
      setFileType(supportedTypes?.[index]?.supported_extensions);
      setSourceType(supportedTypes?.[index]?.design_type);
    }
  };

  useEffect(() => {
    if (input) {
      setIsError(!URLValidator(input));
    }
  }, [input]);

  useEffect(() => {
    if (isDesign) {
      setFileType(supportedTypes?.[0]?.supported_extensions);
      setSourceType(supportedTypes?.[0]?.design_type);
    }
  }, [open]);

  const handleSubmit = async () => {
    await handleUrlUpload(input, sourceType, { name, config });
    handleClose();
  };

  const handleUploader = async (input) => {
    await handleUpload(input, sourceType, { name, config });
    fetch?.();
    handleClose();
  };

  return (
    <>
      <label htmlFor="url-upload-button">
        <Dialog open={open} onClose={handleClose}>
          <MuiThemeProvider
            theme={theme.palette.type == 'dark' ? getDarkMuiTheme() : getMuiTheme()}
          >
            <DialogTitle className={classes.title}>
              <b id="simple-modal-title" style={{ textAlign: 'center' }}>
                Import {configuration}
              </b>
            </DialogTitle>
            <DialogContent>
              <Grid container spacing={24} alignItems="center">
                <Grid item xs={3}>
                  <h4 className={classes.heading}>Name</h4>
                </Grid>
                <Grid item xs={9}>
                  <TextField
                    required
                    size="small"
                    variant="outlined"
                    label="Name"
                    style={{ width: '100%' }}
                    onChange={(e) => setName(e.target.value)}
                  />
                </Grid>
                {isFilter && (
                  <>
                    <Grid container spacing={24} alignItems="center">
                      <Grid item xs={3}>
                        <h4 className={classes.heading}>WASM Config</h4>
                      </Grid>
                      <Grid item xs={9}>
                        <TextField
                          placeholder={
                            'typed_config:\r\n  "@type": type.googleapis.com/envoy.extensions.filters.http.wasm.v3.Wasm\r\n  config:\r\n    name: example-filter\r\n    rootId: my_root_id\r\n    vmConfig:\r\n      code:\r\n        local:\r\n          filename: /var/local/lib/wasm-filters/example-filter.wasm\r\n      runtime: envoy.wasm.runtime.v8\r\n      vmId: example-filter\r\n      allow_precompiled: true\r\nname: envoy.filters.http.wasm'
                          }
                          multiline
                          required
                          minRows={4}
                          size="small"
                          variant="outlined"
                          label="WASM Filter Config"
                          style={{ width: '100%' }}
                          onChange={(e) => setConfig(e.target.value)}
                        />
                      </Grid>
                    </Grid>
                    <br />
                  </>
                )}
                <Grid item xs={12}>
                  <Divider style={{ margin: '8px 0px' }} />
                </Grid>
                <Grid item xs={3}>
                  <h4 className={classes.heading}> FROM URL </h4>
                </Grid>
                <Grid item xs={9}>
                  <TextField
                    size="small"
                    error={isError}
                    helperText={isError && 'Invalid URL'}
                    variant="outlined"
                    label={`URL for ${configuration}`}
                    style={{ width: '100%' }}
                    onChange={(e) => setInput(e.target.value)}
                  />
                </Grid>
              </Grid>
              {isFilter && (
                <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
                  <Typography variant="caption">---OR---</Typography>
                </div>
              )}
              {sourceType !== 'Helm Chart' && (
                <Grid container spacing={24} alignItems="center">
                  <Grid item xs={3}>
                    <h4 className={classes.heading}> UPLOAD FILE </h4>
                  </Grid>
                  <Grid item xs={6}>
                    <TextField
                      size="small"
                      variant="outlined"
                      label="Filename"
                      style={{ width: '100%' }}
                    />
                  </Grid>
                  <Grid item xs={3}>
                    <label htmlFor="upload-button" className={classes.upload}>
                      <Button
                        disabled={sourceType === 'Helm Chart'}
                        variant="contained"
                        className={classes.button}
                        aria-label="Upload Button"
                        onChange={sourceType === 'Helm Chart' ? null : handleUploader}
                        component="span"
                      >
                        <input
                          id="upload-button"
                          type="file"
                          accept={fileType}
                          disabled={sourceType === 'Helm Chart'}
                          hidden
                          name="upload-button"
                          data-cy="file-upload-button"
                        />
                        Browse
                      </Button>
                    </label>
                  </Grid>
                </Grid>
              )}

              <Grid container spacing={24} alignItems="center">
                {isDesign && <h4 className={classes.selectType}>SELECT TYPE </h4>}
                {isDesign && (
                  <>
                    <NativeSelect
                      defaultValue={0}
                      onChange={(e) => handleFileType(e.target.value)}
                      inputProps={{
                        name: 'name',
                        id: 'uncontrolled-native',
                      }}
                    >
                      {supportedTypes?.map((type, index) => (
                        <option key={index} value={index}>
                          {type.design_type}
                        </option>
                      ))}
                    </NativeSelect>
                  </>
                )}
              </Grid>
            </DialogContent>
            <DialogActions>
              <label htmlFor="cancel" className={classes.cancel}>
                <Button variant="outlined" color="secondary" onClick={handleClose}>
                  Cancel
                </Button>
              </label>
              <label htmlFor="URL">
                {' '}
                <Button
                  id="URL"
                  disabled={isError || !input}
                  variant="contained"
                  className={classes.button}
                  onClick={async (e) => {
                    await handleSubmit(e, handleUploader);
                    fetch?.();
                  }}
                >
                  Import
                </Button>{' '}
              </label>
            </DialogActions>
          </MuiThemeProvider>
        </Dialog>
      </label>
    </>
  );
}

export default withStyles(styles)(UploadImport);
