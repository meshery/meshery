import React, {useState} from 'react'
import {Grid, Typography, Button, Box, IconButton, Divider,} from "@mui/material"
import { styled } from "@mui/material/styles";
// import moment from "moment";
import DoneAllIcon from '@mui/icons-material/DoneAll';
import DeleteIcon from '@mui/icons-material/Delete';
import SaveIcon from '@mui/icons-material/Save';
import UndeployIcon from "../../public/static/img/UndeployIcon";
import FullscreenIcon from '@mui/icons-material/Fullscreen';
import FullscreenExitIcon from '@mui/icons-material/FullscreenExit';
import FlipCard from "../FlipCard";

const INITIAL_GRID_SIZE = { xl : 4, md : 6, xs : 12 };

const DeployButtonsWrapper = styled(Box)(({ theme }) => ({
display : "flex",
justifyContent : "flex-end",
gap: theme.spacing(1),
alignItems : "center",
marginTop : "50px"
}))

function MesheryApplicationCard({  
  name,
  updated_at,
  created_at,
  application_file,
  handleDeploy,
  handleUnDeploy,
  updateHandler,
  deleteHandler,
  setYaml}) {

    const [gridProps, setGridProps] = useState(INITIAL_GRID_SIZE);
    const [fullScreen, setFullScreen] = useState(false);

    const toggleFullScreen = () => {
        setFullScreen(!fullScreen);
      };
   
      function genericClickHandler(ev, fn) {
        ev.stopPropagation();
        fn();
      }

  return (
    <>
         {fullScreen &&
        <YAMLDialog
          fullScreen={fullScreen}
          name={name}
          toggleFullScreen={toggleFullScreen}
          config_file={application_file}
          setYaml={setYaml}
          updateHandler={updateHandler}
          deleteHandler={deleteHandler}
        />
      }
    <FlipCard
    onClick={() => {
      console.log(gridProps)
      setGridProps(INITIAL_GRID_SIZE)
    }}
    duration={600}
  >
    {/* FRONT PART */}
    <>
      <div>
        <Typography variant="h6" component="div">
          {name}
        </Typography>
        <div>
          <div>
            {updated_at
              ? (
                <Typography color="primary" variant="caption" style={{ fontStyle : "italic" }}>
                  Modified On: 
                </Typography>
              )
              : null}
          </div>
        </div>
      </div>
      <div>

        <DeployButtonsWrapper>
          <Button
            variant="contained"
            color="primary"
            onClick={(ev) =>
              genericClickHandler(ev, handleDeploy)
            }
          >
            <DoneAllIcon />
            Deploy
          </Button>

          <Button
            variant="unDeploy"
            onClick={(ev) =>
              genericClickHandler(ev, handleUnDeploy)
            }
          >
            <UndeployIcon fill="#ffffff" />
            <span>Undeploy</span>
          </Button>
        </DeployButtonsWrapper>
      </div>
    </>

    {/* BACK PART */}
    <>
      <Grid
        container
        spacing={1}
        alignContent="space-between"
        alignItems="center"
      >
        <Grid item xs={8}>
          <Typography variant="h6">
            {name}
          </Typography>
          <IconButton
            onClick={(ev) =>
              genericClickHandler(ev, () => {
                {
                  toggleFullScreen()
                }
              })
            }
          >
            {fullScreen ? <FullscreenIcon /> : <FullscreenExitIcon />}
          </IconButton>
        </Grid>
        <Grid item xs={12}
          onClick={(ev) =>
            genericClickHandler(ev, () => { })
          }
          >

          <Divider variant="fullWidth" light />

        </Grid>

        <Grid item xs={8}>
          <div>
            <div>
              {created_at
                ? (
                  <Typography color="primary" variant="caption" style={{ fontStyle : "italic" }}>
                    Created at: 
                  </Typography>
                )
                : null}
            </div>
          </div>
        </Grid>

        <Grid item xs={12}>
          <div>

            {/* Save button */}
            <IconButton 
            onClick={(ev) =>
              genericClickHandler(ev, updateHandler)
            }
            >
              <SaveIcon color="primary" />
            </IconButton>

            {/* Delete Button */}
            <IconButton
             onClick={(ev) =>
              genericClickHandler(ev, deleteHandler)
            }
            >
              <DeleteIcon color="primary" />
            </IconButton>
          </div>
        </Grid>
      </Grid>
    </>
  </FlipCard >
  </>
  )
}

export default MesheryApplicationCard