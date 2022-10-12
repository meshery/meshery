//@ts-check
import React, { useState } from "react";
import {
  Avatar, Divider, Grid, IconButton, Typography, Tooltip
} from "@material-ui/core";
import Button from "@material-ui/core/Button";
import DeleteIcon from "@material-ui/icons/Delete";
import Save from "@material-ui/icons/Save";
import Fullscreen from "@material-ui/icons/Fullscreen";
import Moment from "react-moment";
import FlipCard from "../FlipCard";
import { UnControlled as CodeMirror } from "react-codemirror2";
import FullscreenExit from "@material-ui/icons/FullscreenExit";
import UndeployIcon from "../../public/static/img/UndeployIcon";
import DoneAllIcon from '@material-ui/icons/DoneAll';
import DoneIcon from '@material-ui/icons/Done';
import useStyles from "./Cards.styles";
import YAMLDialog from "../YamlDialog";

const INITIAL_GRID_SIZE = { xl : 4, md : 6, xs : 12 };


function MesheryPatternCard({
  name,
  updated_at,
  created_at,
  pattern_file,
  handleVerify,
  handleDeploy,
  handleUnDeploy,
  updateHandler,
  deleteHandler,
  handleClone,
  setSelectedPatterns,
  setYaml,
  description={},
  visibility
}) {

  function genericClickHandler(ev, fn) {
    ev.stopPropagation();
    fn(ev);
  }
  const [gridProps, setGridProps] = useState(INITIAL_GRID_SIZE);
  const [fullScreen, setFullScreen] = useState(false);
  const [showCode, setShowCode] = useState(false);


  const toggleFullScreen = () => {
    setFullScreen(!fullScreen);
  };

  const catalogContentKeys = Object.keys(description);
  const catalogContentValues = Object.values(description);
  const classes = useStyles()

  return (
    <>
      {fullScreen &&
        <YAMLDialog
          fullScreen={fullScreen}
          name={name}
          toggleFullScreen={toggleFullScreen}
          config_file={pattern_file}
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
            <div style={{ display : "flex", justifyContent : "space-between" }}>
              <Typography style={{ overflow : "hidden", textOverflow : "ellipsis", width : '20rem' }} variant="h6" component="div">
                {name}
              </Typography>
              <img  className={classes.img} src={`/static/img/${visibility}.svg`} />
            </div>
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
                title="Verify"
                variant="contained"
                className={classes.testsButton}
                onClick={
                  (e) => genericClickHandler(e, handleVerify)
                }
              >
                <DoneIcon className={classes.iconPatt} />
                <span className={classes.btnText}> Verify </span>
              </Button>

              <Button
                title="Deploy"
                variant="contained"
                onClick={(ev) =>
                  genericClickHandler(ev, handleDeploy)
                }
                className={classes.testsButton}
              >
                <DoneAllIcon className={classes.iconPatt} />
                <span className={classes.btnText}>Deploy</span>
              </Button>

              <Button
                title="Undeploy"
                variant="contained"
                className={classes.undeployButton}
                onClick={(ev) =>
                  genericClickHandler(ev, handleUnDeploy)
                }
              >
                <UndeployIcon fill="#ffffff" className={classes.iconPatt} />
                <span className={classes.btnText}>Undeploy</span>
              </Button>

              { visibility === "private" ?  <Button
                title="Design"
                variant="contained"
                color="primary"
                onClick={(ev) =>
                  genericClickHandler(ev, setSelectedPatterns)
                }
                className={classes.testsButton}
              >
                <Avatar src="/static/img/pattern_trans.svg" className={classes.iconPatt} imgProps={{ height : "16px", width : "16px" }} />
                <span className={classes.btnText}> Design </span>
              </Button> : <Button
                title="Clone"
                variant="contained"
                color="primary"
                onClick={(ev) =>
                  genericClickHandler(ev, handleClone)
                }
                className={classes.testsButton}
              >
                <Avatar src="/static/img/clone-white.svg" className={classes.iconPatt} />
                <span className={classes.btnText}> Clone </span>
              </Button>  }

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
              { catalogContentKeys.length === 0 ?
                <CodeMirror
                  value={showCode && pattern_file}
                  className={fullScreen ? classes.fullScreenCodeMirror : ""}
                  options={{
                    theme : "material",
                    lineNumbers : true,
                    lineWrapping : true,
                    gutters : ["CodeMirror-lint-markers"],
                    // @ts-ignore
                    lint : true,
                    mode : "text/x-yaml",
                  }}
                  onChange={(_, data, val) => setYaml(val)}
                />
                :
                catalogContentKeys.map((title, index) => (
                  <>
                    <Typography variant="h6" className={classes.yamlDialogTitleText}>
                      {title}
                    </Typography>
                    <Typography variant="body2">
                      {catalogContentValues[index]}
                    </Typography>
                  </>
                ))
              }
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
              {visibility === "private" ?
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
                </div> : null}
            </Grid>
          </Grid>
        </>
      </FlipCard >
    </>
  );
}

// @ts-ignore
export default MesheryPatternCard;
