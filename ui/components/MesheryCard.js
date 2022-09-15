import React, { useState } from "react";
import {
  Avatar,
  Divider,
  Grid,
  IconButton,
  Typography,
  Box,
  Button
} from "@mui/material";
import { styled } from "@mui/material/styles";
import DoneAllIcon from "@mui/icons-material/DoneAll";
import RemoveDoneIcon from "@mui/icons-material/RemoveDone";
import FullscreenIcon from "@mui/icons-material/Fullscreen";
import FullscreenExitIcon from "@mui/icons-material/FullscreenExit";
import DeleteIcon from "@mui/icons-material/Delete";
import SaveIcon from "@mui/icons-material/Save";
import FlipCard from "./FlipCard"
const ActionContainer = styled(Box)(({ theme }) => ({
  display: "flex",
  gap:theme.spacing(1),
  justifyContent:"flex-end"
}));
const INITIAL_GRID_SIZE = { xl : 4, md : 6, xs : 12 };
const MesheryCard = ({
  name,
  updated_at,
  created_at,
  file,
  configurations,
  handleDeploy,
  handleUnDeploy,
  updateHandler,
  deleteHandler,
  setSelectedPatterns,
  setYaml
}) => {
  const [fullscreen, setFullScreen] = useState(false);
  const [gridProps, setGridProps] = useState(INITIAL_GRID_SIZE);

  const genericClickHandler = (ev, fn) => {
    ev.stopPropagation();
    fn();
  };

  const toggleFullScreen = () => setFullScreen((prev) => !prev);
  return (
    <FlipCard
    onClick={() => {
      setGridProps(INITIAL_GRID_SIZE)
    }}
    
    >
      {/* Font Part */}
      <>
        <div>
          <Typography variant="h6" component="div">
            {name}
          </Typography>
          <div>
            {updated_at ? (
              <Typography
                color="primary"
                variant="caption"
                style={{ fontStyle: "italic" }}
              >
                Modified On: {updated_at}
              </Typography>
            ) : null}
          </div>
        </div>
        <div style={{ marginTop: "50px" }}>
          <ActionContainer>
            {setSelectedPatterns?
                <Button
                
                variant="containedSecondory"
                onClick={(ev) => genericClickHandler(ev, setSelectedPatterns)}
                >
                <Avatar
                    src="/static/img/patternwhite.svg"
                    sx={{ width: 27, height: 27, marginRight: 1 }}
                />
                Design
                </Button>
            :null}
            <Button
              variant="contained"
              onClick={(ev) => genericClickHandler(ev, handleDeploy)}
              
            >
              <DoneAllIcon style={{marginRight:"5px"}}/>
              Deploy
            </Button>
            <Button
              variant="unDeploy"
              onClick={(ev) => genericClickHandler(ev, handleUnDeploy)}
            >
              <RemoveDoneIcon style={{marginRight:"5px"}}/> Undeploy
            </Button>
          </ActionContainer>
        </div>
      </>
      {/* Back Part */}
      <>
      <Grid
        conatainer
        spacing={1}
        alignContent="space-between"
        alignItems="center"
      >
        <Grid item container justifyContent="space-between">
          <Typography variant="h6" component="div">
            {name}
          </Typography>
          <IconButton
            onClick={(ev) => genericClickHandler(ev, toggleFullScreen)}
          >
            {!fullscreen ? <FullscreenIcon /> : <FullscreenExitIcon />}
          </IconButton>
        </Grid>
        <Grid item>
          <Divider variant="fullWidth" light />
          {/* Code Editor */}
        </Grid>
        <Grid item container justifyContent="space-between" alignItem="center">
          {created_at ? (
            <Typography
              variant="caption"
              component="div"
              style={{ fontStyle: "italic", paddingTop: "13px" }}
            >
              Created At: {created_at}
            </Typography>
          ) : null}
          <ActionContainer>
            {/* Save button */}
            {configurations!=="filter"?
                <IconButton
                onClick={(ev) => genericClickHandler(ev, updateHandler)}
              >
                <SaveIcon />
              </IconButton>
            :null}
            

            {/* Delete Button */}
            
            <IconButton
              onClick={(ev) => genericClickHandler(ev, deleteHandler)}
            >
              <DeleteIcon />
            </IconButton>
          </ActionContainer>
        </Grid>
        
      </Grid>
      </>
    </FlipCard>
  );
};

export default MesheryCard;