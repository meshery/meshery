import React from 'react'
import {TextField, Button, Grid } from "@mui/material";
import LinkIcon from '@mui/icons-material/Link';

function UploadImport({configuration}) {
  return (
    <Button aria-label="URL-Upload" data-cy="import-button" variant="contained"
          color="primary"
          size="large" >
          <LinkIcon style={{ paddingRight : ".35rem" }} />
          Import {configuration}
        </Button>
  )
}

export default UploadImport