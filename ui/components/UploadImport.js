import React, {useState} from 'react'
import {TextField, Button, Grid,   Dialog, DialogActions,
  DialogContent,
  DialogTitle} from "@mui/material";
  import { useTheme } from "@mui/system";
import LinkIcon from '@mui/icons-material/Link';

function UploadImport({configuration}) {
  const theme = useTheme();
  const [open, setOpen] = useState(false);

  const handleOpenModal = () => setOpen(true);
  const handleCloseModal = () => setOpen(false);
  
  return (
    <>
    <Button aria-label="URL-Upload" data-cy="import-button" variant="contained"
          color="primary"
          size="large"
          onClick={handleOpenModal} sx={{marginBottom: theme.spacing(2), marginLeft: theme.spacing(2) }} >
          <LinkIcon sx={{ paddingRight : ".35rem" }} />
          Import {configuration}
        </Button>
      <Dialog 
        open={open} handleClose={handleCloseModal} >

     <DialogTitle> Import {configuration} </DialogTitle>
     <DialogContent>
        <>
       <Grid container sx={{paddingTop: theme.spacing(2)}} >
        <Grid item xs={3}>
          <h4 > FROM URL </h4>
        </Grid>
        <Grid item xs={9}>
          <TextField
            variant="outlined"
            label={`URL for ${configuration}`}
            sx={{ width : "100%" }}
            />
        </Grid>
      </Grid>
      <hr />
      <Grid container spacing={2}>
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
      </Grid>
      </>
      </DialogContent>
      <DialogActions>
              <label htmlFor="cancel">
                <Button variant="outlined" size="large" color="secondary" onClick={handleCloseModal}>Cancel</Button>
              </label>
              <label htmlFor="URL">  <Button size="large" id="URL" variant="contained" color="primary" > Import</Button> </label>

            </DialogActions>
       </Dialog>
     </>   
  )
}

export default UploadImport


 {/* */}