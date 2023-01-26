import React, { useEffect, useState } from 'react'
import { styled } from "@mui/material/styles"
import {
  Button,
  Grid,
  IconButton,
  Paper,
  TextField,
  Tooltip,
} from "@mui/material";
import LinkIcon from "@mui/icons-material/Link";
import GenericModal from "./GenericModal";
import { URLValidator } from '../utils/URLValidator';

const PaperContent = styled(Paper)(({ theme }) => ({
  position : "absolute",
  width : 600,
  backgroundColor : theme.palette.background.paper,
  border : "0px solid #000",
  boxShadow : theme.shadows[5],
  padding : theme.spacing(2, 4, 3),
  top : "50%",
  left : "50%",
  transform : `translate(-50%, -50%)`,
  borderRadius : 10,
}));

const URLUploader = ({ onSubmit }) => {
  const [open, setOpen] = useState(false);
  const [input, setInput] = useState();
  const [isError, setIsError] = useState(false);

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
    <div>
      <label htmlFor="url-upload-button">
        <Tooltip title="Upload URL">
          <IconButton
            aria-label="URL-Upload"
            component="span"
            onClick={handleOpen}
          >
            <LinkIcon />
          </IconButton>
        </Tooltip>
        <GenericModal
          open={open}
          handleClose={handleClose}
          Content={
            <PaperContent>
              <Grid container spacing={2}>
                <Grid item sx={12}>
                  <h2 id="simple-modal-title">Import using URL</h2>
                </Grid>
                <Grid item xs={12}>
                  <TextField
                    error={isError}
                    helperText={isError && "Invalid URL"}
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
                    {" "}
                    Import
                  </Button>
                </Grid>
              </Grid>
            </PaperContent>
          }
        />
      </label>
    </div>
  );
}


export default URLUploader;
