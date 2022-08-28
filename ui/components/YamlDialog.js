import React, {useState} from 'react'
import { AppBar, Box, Button, Grid, Dialog, DialogTitle, Divider, DialogContent, DialogActions, Tooltip, Hidden, IconButton, Toolbar, Typography, TableCell, TableSortLabel } from "@mui/material";
import SaveIcon from '@mui/icons-material/Save';
import DeleteIcon from '@mui/icons-material/Delete';
import FullscreenIcon from '@mui/icons-material/Fullscreen';
import FullscreenExitIcon from '@mui/icons-material/FullscreenExit';

function TooltipIcon({ children, onClick, title }) {
    return (
      <Tooltip title={title} placement="top" arrow interactive >
        <IconButton onClick={onClick}>
          {children}
        </IconButton>
      </Tooltip>
    );
  }

function YAMLEditor({ application, onClose, onSubmit }) {
    const [yaml, setYaml] = useState("");
    const [fullScreen, setFullScreen] = useState(false);
  
    const toggleFullScreen = () => {
      setFullScreen(!fullScreen);
    };
  
    return (
      <>
      <Dialog onClose={onClose} aria-labelledby="application-dialog-title" open maxWidth="md" fullScreen={fullScreen} fullWidth={!fullScreen}>
        <DialogTitle disableTypography id="application-dialog-title" >
          <Typography variant="h6" >
            {application.name} 
          </Typography>
          <TooltipIcon
            title={fullScreen ? "Exit Fullscreen" : "Enter Fullscreen"}
            onClick={toggleFullScreen}>
            {fullScreen ? <FullscreenExitIcon /> : <FullscreenIcon />}
          </TooltipIcon>
        </DialogTitle>
        <Divider variant="fullWidth" light />
        <DialogContent>
          {/* <CodeMirror
            value={application.application_file}
            className={fullScreen ? classes.fullScreenCodeMirror : ""}
            options={{
              theme : "material",
              lineNumbers : true,
              lineWrapping : true,
              gutters : ["CodeMirror-lint-markers"],
              lint : true,
              mode : "text/x-yaml",
            }}
            onChange={(_, data, val) => setYaml(val)}
          /> */}
         <h1> H! </h1>
        </DialogContent>
        <Divider variant="fullWidth" light />
        <DialogActions>
          <Tooltip title="Update Application">
            <IconButton
              aria-label="Update"
              color="primary"
              onClick={() => onSubmit({
                data : yaml, id : application.id, name : application.name, type : FILE_OPS.UPDATE
              })}
            >
              <SaveIcon />
            </IconButton>
          </Tooltip>
          <Tooltip title="Delete Application">
            <IconButton
              aria-label="Delete"
              color="primary"
              onClick={() => onSubmit({
                data : yaml,
                id : application.id,
                name : application.name,
                type : FILE_OPS.DELETE
              })}
            >
              <DeleteIcon />
            </IconButton>
          </Tooltip>
        </DialogActions>
      </Dialog>
      </>
    );
  }

  export default YAMLEditor;