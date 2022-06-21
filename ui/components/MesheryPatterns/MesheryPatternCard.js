//@ts-check
import React, { useState } from "react";
import {
  Avatar, Dialog, DialogActions, DialogContent, DialogTitle, Divider, Grid, IconButton, Tooltip, Typography
} from "@material-ui/core";
import Button from "@material-ui/core/Button";
import DeleteIcon from "@material-ui/icons/Delete";
import Save from "@material-ui/icons/Save";
import Fullscreen from "@material-ui/icons/Fullscreen";
import Moment from "react-moment";
import FlipCard from "../FlipCard";
import { makeStyles } from "@material-ui/core/styles";
import { UnControlled as CodeMirror } from "react-codemirror2";
import FullscreenExit from "@material-ui/icons/FullscreenExit";
import UndeployIcon from "../../public/static/img/UndeployIcon";
import DoneAllIcon from '@material-ui/icons/DoneAll';

const INITIAL_GRID_SIZE = { xl : 4, md : 6, xs : 12 };

const useStyles= makeStyles(() => ({
  cardButtons : {
    display : "flex",
    justifyContent : "flex-end",
    alignItems : "center",
    marginTop : "50px"
  },
  testsButton : {
    marginRight : "0.5rem"
  },
  perfResultsContainer : {
    marginTop : "0.5rem"
  },
  backGrid : {
    marginBottom : "0.25rem",
    minHeight : "6rem",
    position : "relative"
  },
  updateDeleteButtons : {
    width : "fit-content",
    margin : "10 0 0 auto",
    position : "absolute",
    right : 0,
    bottom : 0,
  },
  yamlDialogTitle : {
    display : "flex",
    alignItems : "center"
  },
  yamlDialogTitleText : {
    flexGrow : 1
  },
  fullScreenCodeMirror : {
    height : '100%',
    width : '100%',
    '& .CodeMirror' : {
      minHeight : "300px",
      height : '100%',
      width : '100%'
    }
  },
  maximizeButton : {
    width : "fit-content",
    margin : "0 0 0 auto",
    position : "absolute",
    right : 0,
    top : 0
  },
  noOfResultsContainer : {
    margin : "0 0 1rem",
    '& div' : {
      display : "flex",
      alignItems : "center"
    },
  },
  bottomPart : {
    display : "flex",
    justifyContent : "flex-end",
    alignItems : "center",
  },
  lastRunText : {
    marginRight : "0.5rem"

  },
  iconPatt : {
    width : "24px",
    height : "24px",
    marginRight : "5px"
  },
  undeployButton : {
    backgroundColor : "#B32700",
    color : "#ffffff"
  }
}))

//Full screen Dialog: Similar to the dialog on the table view, with few modification on parameters
const YAMLDialog = ({
  fullScreen,
  name,
  toggleFullScreen,
  pattern_file,
  setYaml,
  updateHandler,
  deleteHandler
}) => {
  const classes = useStyles()
  return (
    <Dialog aria-labelledby="pattern-dialog-title" open maxWidth="md" fullScreen={fullScreen} fullWidth={!fullScreen}>
      <DialogTitle disableTypography id="pattern-dialog-title" className={classes.yamlDialogTitle}>
        <Typography variant="h6" className={classes.yamlDialogTitleText}>
          {name}
        </Typography>
        <IconButton
          onClick={toggleFullScreen}>
          {fullScreen ? <FullscreenExit /> : <Fullscreen />}
        </IconButton>
      </DialogTitle>
      <Divider variant="fullWidth" light />
      <DialogContent>
        <CodeMirror
          value={pattern_file}
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
        />
      </DialogContent>
      <Divider variant="fullWidth" light />
      <DialogActions>
        <Tooltip title="Update Pattern">
          <IconButton
            aria-label="Update"
            color="primary"
            onClick={updateHandler}
          >
            <Save />
          </IconButton>
        </Tooltip>
        <Tooltip title="Delete Pattern">
          <IconButton
            aria-label="Delete"
            color="primary"
            onClick={deleteHandler}
          >
            <DeleteIcon />
          </IconButton>
        </Tooltip>
      </DialogActions>
    </Dialog>
  )
}

function MesheryPatternCard({
  name,
  updated_at,
  created_at,
  pattern_file,
  handleDeploy,
  handleUnDeploy,
  updateHandler,
  deleteHandler,
  setSelectedPatterns,
  setYaml,
  // requestFullSize,
  // requestSizeRestore,
}) {

  function genericClickHandler(ev, fn) {
    ev.stopPropagation();
    fn();
  }
  const [gridProps, setGridProps] = useState(INITIAL_GRID_SIZE);
  const [fullScreen, setFullScreen] = useState(false);

  const toggleFullScreen = () => {
    setFullScreen(!fullScreen);
  };

  const classes=useStyles()

  return (
    <>
      {fullScreen &&
        <YAMLDialog
          fullScreen={fullScreen}
          name={name}
          toggleFullScreen={toggleFullScreen}
          pattern_file={pattern_file}
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
            <div className={classes.lastRunText} >
              <div>
                {updated_at
                  ? (
                    <Typography color="primary" variant="caption" style={{ fontStyle : "italic" }}>
                  Modified On: <Moment format="LLL">{updated_at}</Moment>
                    </Typography>
                  )
                  : null}
              </div>
            </div>
          </div>
          <div className={classes.bottomPart} >

            <div className={classes.cardButtons} >
              <Button
                variant="contained"
                onClick={(ev) =>
                  genericClickHandler(ev, setSelectedPatterns)
                }
                className={classes.testsButton}
              >
                <Avatar src="/static/img/pattwhite.svg" className={classes.iconPatt} imgProps={{ height : "16px", width : "16px" }} />
              Design
              </Button>

              <Button
                variant="contained"
                color="primary"
                onClick={(ev) =>
                  genericClickHandler(ev, handleDeploy)
                }
                className={classes.testsButton}
              >
                <DoneAllIcon className={classes.iconPatt}/>
              Deploy
              </Button>

              <Button
                variant="contained"
                className={classes.undeployButton}
                onClick={(ev) =>
                  genericClickHandler(ev, handleUnDeploy)
                }
              >
                <UndeployIcon fill="#ffffff" className={classes.iconPatt} />
              Undeploy
              </Button>
            </div>
          </div>
        </>

        {/* BACK PART */}
        <>
          <Grid className={classes.backGrid}
            container
            spacing={1}
            alignContent="space-between"
            alignItems="center"
          >
            <Grid item xs={8} className={classes.yamlDialogTitle}>
              <Typography variant="h6" className={classes.yamlDialogTitleText}>
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
                className={classes.maximizeButton}
              >
                {fullScreen ? <FullscreenExit /> : <Fullscreen />}
              </IconButton>
            </Grid>
            <Grid item xs={12}
              onClick={(ev) =>
                genericClickHandler(ev, () => {})
              }>

              <Divider variant="fullWidth" light />

              <CodeMirror
                value={pattern_file}
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
              />
            </Grid>

            <Grid item xs={8}>
              <div className={classes.lastRunText} >
                <div>
                  {created_at
                    ? (
                      <Typography color="primary" variant="caption" style={{ fontStyle : "italic" }}>
                  Created at: <Moment format="LLL">{created_at}</Moment>
                      </Typography>
                    )
                    : null}
                </div>
              </div>
            </Grid>

            <Grid item xs={12}>
              <div className={classes.updateDeleteButtons} >

                {/* Save button */}
                <IconButton onClick={(ev) =>
                  genericClickHandler(ev,updateHandler)
                }>
                  <Save color="primary" />
                </IconButton>

                {/* Delete Button */}
                <IconButton onClick={(ev) =>
                  genericClickHandler(ev,deleteHandler)
                }>
                  <DeleteIcon color="primary" />
                </IconButton>
              </div>
            </Grid>
          </Grid>
        </>
      </FlipCard >
    </>
  );
}

// @ts-ignore
export default MesheryPatternCard;
