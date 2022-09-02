//@ts-check
import React, { useState } from "react";
import {
  Divider, Grid, IconButton, Typography
} from "@material-ui/core";
import Button from "@material-ui/core/Button";
import DeleteIcon from "@material-ui/icons/Delete";
import Fullscreen from "@material-ui/icons/Fullscreen";
import Moment from "react-moment";
import FlipCard from "../FlipCard";
import { UnControlled as CodeMirror } from "react-codemirror2";
import FullscreenExit from "@material-ui/icons/FullscreenExit";
import DoneAllIcon from '@material-ui/icons/DoneAll';
import useStyles from "../MesheryPatterns/Cards.styles";
import YAMLDialog from "../YamlDialog";
import UndeployIcon from "../../public/static/img/UndeployIcon";

const INITIAL_GRID_SIZE = { xl : 4, md : 6, xs : 12 };

function FiltersCard({
  name,
  updated_at,
  created_at,
  filter_file,
  handleDeploy,
  handleUndeploy,
  handleClone,
  deleteHandler,
  setYaml,
  description={},
  visibility
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

  const catalogContentKeys = Object.keys(description);
  const catalogContentValues = Object.values(description);
  const classes=useStyles()

  return (
    <>
      {fullScreen &&
        <YAMLDialog
          fullScreen={fullScreen}
          name={name}
          toggleFullScreen={toggleFullScreen}
          config_file={filter_file}
          setYaml={setYaml}
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
              <Typography variant="h6" component="div">
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
                variant="contained"
                className={classes.undeployButton}
                onClick={(ev) =>
                  genericClickHandler(ev, handleUndeploy)
                }
              >
                <UndeployIcon fill="#ffffff" className={classes.iconPatt} />
                Undeploy
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
              {visibility === "public" ? <Button
                variant="contained"
                color="primary"
                onClick={(ev) =>
                  genericClickHandler(ev, handleClone)
                }>
                <img src="/static/img/clone-white.svg" className={classes.iconPatt} />
                  Clone
              </Button> : null }
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

              { catalogContentKeys.length === 0 ?
                <CodeMirror
                  value={showCode && filter_file}
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
              { visibility === "private" ? <div className={classes.deleteButton} >
                <IconButton onClick={(ev) =>
                  genericClickHandler(ev,deleteHandler)
                }>
                  <DeleteIcon color="primary" />
                </IconButton>
              </div> : null}
            </Grid>
          </Grid>
        </>
      </FlipCard >
    </>
  );
}

// @ts-ignore
export default FiltersCard;
