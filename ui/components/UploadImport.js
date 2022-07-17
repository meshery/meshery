import React, {useState} from 'react'
import {TextField, Button, Grid,} from "@mui/material";
import LinkIcon from '@mui/icons-material/Link';
import CustomDialog from "@/components/Dialog"

function UploadImport({configuration}) {
  const [open, setOpen] = useState(false);

  const handleOpenModal = () => setOpen(true);
  const handleCloseModal = () => setOpen(false);
  
  return (
    <>
    <Button aria-label="URL-Upload" data-cy="import-button" variant="contained"
          color="primary"
          size="large"
          onClick={handleOpenModal} >
          <LinkIcon style={{ paddingRight : ".35rem" }} />
          Import {configuration}
        </Button>
      <CustomDialog 
        open={open} handleClose={handleCloseModal}
        // handleClose={handleClose}
      Title="qwe"
      Content={
        <>
      <h1>HI</h1>
      </>
      }
      />  
     </>   
  )
}

export default UploadImport


 {/* <Grid container spacing={24}>
        <Grid item xs={3}>
          <h4 > FROM URL </h4>
        </Grid>
        <Grid item xs={9}>
          <TextField
            variant="outlined"
            label={`URL for ${configuration}`}
            style={{ width : "100%" }}
            />
        </Grid>
      </Grid>
      <hr />
      <Grid container spacing={24}>
        <Grid item xs={3}>
          <h4> UPLOAD FILE </h4>
        </Grid>
        <Grid item xs={6}>
          <TextField
            variant="outlined"
            label="Filename"
            sx={{ width : "100%" }}
          />
        </Grid>
        <Grid item xs={3}>
          <label htmlFor="upload-button">
            <Button variant="contained" size="large" color="primary" aria-label="Upload Button" component="span" >
              <input id="upload-button" type="file" accept=".yaml, .yml" hidden name="upload-button" data-cy="file-upload-button" />
              Browse
            </Button>
          </label>
        </Grid>
      </Grid> */}