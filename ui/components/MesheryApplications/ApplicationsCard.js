//@ts-check
import React, { useState } from "react";
import {
  Divider, Grid, IconButton, Typography, Tooltip
} from "@material-ui/core";
import DeleteIcon from "@material-ui/icons/Delete";
import Save from "@material-ui/icons/Save";
import Fullscreen from "@material-ui/icons/Fullscreen";
import Moment from "react-moment";
import FlipCard from "../FlipCard";
import { UnControlled as CodeMirror } from "react-codemirror2";
import FullscreenExit from "@material-ui/icons/FullscreenExit";
import UndeployIcon from "../public/static/img/UndeployIcon";
import DoneAllIcon from '@material-ui/icons/DoneAll';
import useStyles from "../MesheryPatterns/Cards.styles";
import YAMLDialog from "../YamlDialog";
import TooltipButton from '../../utils/TooltipButton.js'
import { useTheme } from '@material-ui/core/styles';

const INITIAL_GRID_SIZE = { xl : 4, md : 6, xs : 12 };


function MesheryApplicationCard({
  name,
  id,
  updated_at,
  created_at,
  application_file,
  handleDeploy,
  handleUnDeploy,
  updateHandler,
  deleteHandler,
  setYaml,
  source_type,
  handleAppDownload
}) {

  function genericClickHandler(ev, fn) {
    ev.stopPropagation();
    fn();
  }
  const [gridProps, setGridProps] = useState(INITIAL_GRID_SIZE);
  const [fullScreen, setFullScreen] = useState(false);
  const [showCode, setShowCode] = useState(false);

  const toggleFullScreen = () => {
    setFullScreen(!fullScreen);
  };
  const theme = useTheme();
  const classes = useStyles()

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
        onShow={() => setTimeout(() => setShowCode(currentCodeVisibilty => !currentCodeVisibilty),500)}
      >
        {/* FRONT PART */}
        <>
          <div>
            <div className="helm" style={{ display : "flex", justifyContent : "space-between" }}>
              <Typography variant="h6" component="div">
                {name}
              </Typography>
              <IconButton
                title="click to download"
                onClick={(e) => {
                  e.stopPropagation(); handleAppDownload(id ,source_type, name)
                }}
              >
                <img src={`/static/img/${(source_type).replaceAll(" ", "_").toLowerCase()}${(theme.palette.type)==='dark'?"-light":''}.svg`} width="45px" height="45px" />
              </IconButton>
            </div>
            <div className={classes.lastRunText} >
              <div>
                {updated_at
                  ? (
                    <Typography variant="caption" style={{ fontStyle : "italic", color : `${theme.palette.type === "dark" ? "rgba(255, 255, 255, 0.7)" : "#647881"}` }}>
                      Modified On: <Moment format="LLL">{updated_at}</Moment>
                    </Typography>
                  )
                  : null}
              </div>
            </div>
          </div>
          <div className={classes.bottomPart} >

            <div className={classes.cardButtons} >
              <TooltipButton
                title="Undeploy"
                variant="contained"
                className={classes.undeployButton}
                onClick={(ev) =>
                  genericClickHandler(ev, handleUnDeploy)
                }
              >
                <UndeployIcon fill="#ffffff" className={classes.iconPatt} />
                Undeploy
              </TooltipButton>
              <TooltipButton
                title="deploy"
                variant="contained"
                color="primary"
                onClick={(ev) =>
                  genericClickHandler(ev, handleDeploy)
                }
                className={classes.testsButton}
              >
                <DoneAllIcon className={classes.iconPatt} />
                Deploy
              </TooltipButton>

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
              <Tooltip
                title="Enter Fullscreen" arrow interactive placement="top"
              >
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
              </Tooltip>
            </Grid>
            <Grid item xs={12}
              onClick={(ev) =>
                genericClickHandler(ev, () => { })
              }>

              <Divider variant="fullWidth" light />

              <CodeMirror
                value={showCode && application_file}
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
                <Tooltip
                  title="Save" arrow interactive placement="bottom"
                >
                  <IconButton onClick={(ev) =>
                    genericClickHandler(ev, updateHandler)
                  }>
                    <Save color="primary" />
                  </IconButton>
                </Tooltip>
                {/* Delete Button */}
                <Tooltip
                  title="Delete" arrow interactive placement="bottom"
                >
                  <IconButton onClick={(ev) =>
                    genericClickHandler(ev, deleteHandler)
                  }>
                    <DeleteIcon color="primary" />
                  </IconButton>
                </Tooltip>
              </div>
            </Grid>
          </Grid>
        </>
      </FlipCard >
    </>
  );
}

// @ts-ignore
export default MesheryApplicationCard;
