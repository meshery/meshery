import React, { useEffect } from 'react'
import { TextField, Button, Grid, NativeSelect } from '@mui/material';
import { styled } from "@mui/material/styles"
import { URLValidator } from '../utils/URLValidator';
import {
  Dialog, DialogActions,
  DialogContent,
  DialogTitle
} from '@mui/material';
import { promisifiedDataFetch } from '../lib/data-fetch';

const LabelUpload = styled("label")(() => ({
  paddingLeft : "0.7rem",
  paddingTop : "8px",
}));

const DialogTitleTitle = styled(DialogTitle)(() => ({
  textAlign : "center",
  minWidth : 500,
  padding : "15px",
  color : "#fff",
  backgroundColor : "#396679",
}));

const H4Heading = styled("h4")(() => ({
  color : "#607d8b",
}));

const H4SelectType = styled("h4")(() => ({
  color : "#607d8b",
  marginRight : "1.2rem",
}));


function UploadImport(props) {
  const {
    handleUpload,
    handleUrlUpload,
    configuration,
    isApplication,
    open,
    handleClose,
    fetch,
  } = props;
  const [input, setInput] = React.useState();
  const [isError, setIsError] = React.useState(false);
  const [fileType, setFileType] = React.useState();
  const [sourceType, setSourceType] = React.useState();
  const [supportedTypes, setSupportedTypes] = React.useState();

  useEffect(() => {
    if (isApplication) {
      (async () => {
        setSupportedTypes(await promisifiedDataFetch("/api/application/types"));
      })();
    }
  }, []);

  const handleFileType = (index) => {
    if (isApplication) {
      setFileType(supportedTypes?.[index]?.supported_extensions);
      setSourceType(supportedTypes?.[index]?.application_type);
    }
  };

  useEffect(() => {
    if (input) {
      setIsError(!URLValidator(input));
    }
  }, [input]);

  useEffect(() => {
    if (isApplication) {
      setFileType(supportedTypes?.[0]?.supported_extensions);
      setSourceType(supportedTypes?.[0]?.application_type);
    }
  }, [open]);

  const handleSubmit = async () => {
    await handleUrlUpload(input, sourceType);
    handleClose();
  };

  const handleUploader = async (input) => {
    await handleUpload(input, sourceType);
    handleClose();
  };

  return (
    <>
      <label htmlFor="url-upload-button">
        <Dialog open={open} onClose={handleClose}>
          <DialogTitleTitle>
            <b id="simple-modal-title" style={{ textAlign : "center" }}>
                Import {configuration}
            </b>
          </DialogTitleTitle>
          <DialogContent>
            <Grid container spacing={24} alignItems="center">
              <Grid item xs={3}>
                <H4Heading> FROM URL </H4Heading>
              </Grid>
              <Grid item xs={9}>
                <TextField
                  size="small"
                  error={isError}
                  helperText={isError && "Invalid URL"}
                  variant="outlined"
                  label={`URL for ${configuration}`}
                  style={{ width : "100%" }}
                  onChange={(e) => setInput(e.target.value)}
                />
              </Grid>
            </Grid>
            <hr />
            {sourceType !== "Helm Chart" && (
              <Grid container spacing={24} alignItems="center">
                <Grid item xs={3}>
                  <H4Heading> UPLOAD FILE </H4Heading>
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
                  <LabelUpload htmlFor="upload-button">
                    <Button
                      disabled={sourceType === "Helm Chart"}
                      variant="contained"
                      color="primary"
                      aria-label="Upload Button"
                      onChange={
                        sourceType === "Helm Chart" ? null : handleUploader
                      }
                      component="span"
                    >
                      <input
                        id="upload-button"
                        type="file"
                        accept={fileType}
                        disabled={sourceType === "Helm Chart"}
                        hidden
                        name="upload-button"
                        data-cy="file-upload-button"
                      />
                        Browse
                    </Button>
                  </LabelUpload>
                </Grid>
              </Grid>
            )}

            <Grid container spacing={24} alignItems="center">
              {isApplication && <H4SelectType>SELECT TYPE </H4SelectType>}
              {isApplication && (
                <>
                  <NativeSelect
                    defaultValue={0}
                    onChange={(e) => handleFileType(e.target.value)}
                    inputProps={{
                      name : "name",
                      id : "uncontrolled-native",
                    }}
                  >
                    {supportedTypes?.map((type, index) => (
                      <option key={index} value={index}>
                        {type.application_type}
                      </option>
                    ))}
                  </NativeSelect>
                </>
              )}
            </Grid>
          </DialogContent>
          <DialogActions>
            <label htmlFor="cancel">
              <Button
                variant="outlined"
                color="secondary"
                onClick={handleClose}
              >
                  Cancel
              </Button>
            </label>
            <label htmlFor="URL">
              {" "}
              <Button
                disabled={isError || !input}
                id="URL"
                variant="contained"
                color="primary"
                onClick={async (e) => {
                  await handleSubmit(e, handleUploader);
                  fetch();
                }}
              >
                  Import
              </Button>{" "}
            </label>
          </DialogActions>
        </Dialog>
      </label>
    </>
  );
}

export default UploadImport;
