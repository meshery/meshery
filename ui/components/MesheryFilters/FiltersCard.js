//@ts-check
import React, { useState } from "react";
import {
  Divider, Grid, IconButton, Typography, Tooltip
} from "@material-ui/core";
import DeleteIcon from "@material-ui/icons/Delete";
import Fullscreen from "@material-ui/icons/Fullscreen";
import Save from "@material-ui/icons/Save";
import Moment from "react-moment";
import FlipCard from "../FlipCard";
import { UnControlled as CodeMirror } from "react-codemirror2";
import FullscreenExit from "@material-ui/icons/FullscreenExit";
import useStyles from "../MesheryPatterns/Cards.styles";
import YAMLDialog from "../YamlDialog";
import CloneIcon from "../../public/static/img/CloneIcon";
import PublicIcon from '@material-ui/icons/Public';
import TooltipButton from '../../utils/TooltipButton.js'
import { VISIBILITY } from "../../utils/Enum";

const INITIAL_GRID_SIZE = { xl : 4, md : 6, xs : 12 };

function FiltersCard({
  name,
  updated_at,
  created_at,
  filter_resource,
  handleClone,
  deleteHandler,
  setYaml,
  description={},
  visibility,
  handlePublishModal,
  handleUnpublishModal,
  updateHandler,
  canPublishFilter = false
}) {

  const genericClickHandler = (ev, fn) => {
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
  const classes=useStyles()

  return (
    <>
      {fullScreen &&
        <YAMLDialog
          fullScreen={fullScreen}
          name={name}
          toggleFullScreen={toggleFullScreen}
          config_file={filter_resource}
          setYaml={setYaml}
          deleteHandler={deleteHandler}
          updateHandler={updateHandler}
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
              {canPublishFilter &&
                (visibility !== VISIBILITY.PUBLISHED)
                ?
                (<TooltipButton
                  variant="contained"
                  title="Publish"
                  className={classes.testsButton}
                  onClick={(ev) => genericClickHandler(ev, handlePublishModal)}
                >
                  <PublicIcon className={classes.iconPatt} />
                  <span className={classes.btnText}> Publish </span>
                </TooltipButton>)
                : (
                  <TooltipButton
                    variant="contained"
                    title="Unpublish"
                    className={classes.testsButton}
                    onClick={(ev) => genericClickHandler(ev, handleUnpublishModal)}
                  >
                    <PublicIcon className={classes.iconPatt} />
                    <span className={classes.btnText}> Unpublish </span>
                  </TooltipButton>
                )
              }

              {visibility === VISIBILITY.PUBLISHED ?
                <TooltipButton
                  title="Clone"
                  variant="contained"
                  color="primary"
                  onClick={(ev) =>
                    genericClickHandler(ev, handleClone)
                  }
                  className={classes.testsButton}
                >
                  <CloneIcon fill="#ffffff" className={classes.iconPatt} />
                  <span className={classes.cloneBtnText}> Clone </span>
                </TooltipButton> : null }
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
                title="Enter Fullscreen"
                arrow interactive
                placement="top"
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
                genericClickHandler(ev, () => {})
              }>

              <Divider variant="fullWidth" light />

              { catalogContentKeys.length === 0 ?
                <CodeMirror
                  value={showCode && filter_resource}
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
export default FiltersCard;