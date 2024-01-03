import React, { useState } from 'react';
import { Avatar, Divider, Grid, IconButton, Typography, Tooltip, Link } from '@material-ui/core';
import DeleteIcon from '@material-ui/icons/Delete';
import Save from '@material-ui/icons/Save';
import Fullscreen from '@material-ui/icons/Fullscreen';
import Moment from 'react-moment';
import FlipCard from '../FlipCard';
import { UnControlled as CodeMirror } from 'react-codemirror2';
import FullscreenExit from '@material-ui/icons/FullscreenExit';
import UndeployIcon from '../../public/static/img/UndeployIcon';
import DoneAllIcon from '@material-ui/icons/DoneAll';
import DoneIcon from '@material-ui/icons/Done';
import useStyles from './Cards.styles';
import YAMLDialog from '../YamlDialog';
import PublicIcon from '@material-ui/icons/Public';
import TooltipButton from '@/utils/TooltipButton';
import CloneIcon from '../../public/static/img/CloneIcon';
import { VISIBILITY } from '@/utils/Enum';
import { useTheme } from '@material-ui/core/styles';
import { useRouter } from 'next/router';
import { Edit } from '@material-ui/icons';
import InfoOutlinedIcon from '@material-ui/icons/InfoOutlined';
import { MESHERY_CLOUD_PROD } from '../../constants/endpoints';
import { useGetUserByIdQuery } from '../../rtk-query/user';
import { Provider } from 'react-redux';
import { store } from '../../store';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';

const INITIAL_GRID_SIZE = { xl: 4, md: 6, xs: 12 };

function MesheryPatternCard_({
  id,
  name,
  updated_at,
  created_at,
  pattern_file,
  handleVerify,
  handlePublishModal,
  handleUnpublishModal,
  handleDeploy,
  handleUnDeploy,
  updateHandler,
  deleteHandler,
  handleClone,
  setSelectedPatterns,
  setYaml,
  description = {},
  visibility,
  canPublishPattern = false,
  user,
  pattern,
  handleInfoModal,
}) {
  const router = useRouter();

  const genericClickHandler = (ev, fn) => {
    ev.stopPropagation();
    fn(ev);
  };
  const [gridProps, setGridProps] = useState(INITIAL_GRID_SIZE);
  const [fullScreen, setFullScreen] = useState(false);
  const [showCode, setShowCode] = useState(false);

  const toggleFullScreen = () => {
    setFullScreen(!fullScreen);
  };

  const { data: owner } = useGetUserByIdQuery(pattern.user_id || '');

  const catalogContentKeys = Object.keys(description);
  const catalogContentValues = Object.values(description);
  const classes = useStyles();
  const theme = useTheme();

  const editInConfigurator = () => {
    router.push('/configuration/designs/configurator?design_id=' + id);
  };
  const userCanEdit =
    CAN(keys.EDIT_DESIGN.action, keys.EDIT_DESIGN.subject) || user?.user_id == pattern?.user_id; // allow if owner

  return (
    <>
      {fullScreen && (
        <YAMLDialog
          fullScreen={fullScreen}
          name={name}
          toggleFullScreen={toggleFullScreen}
          config_file={pattern_file}
          setYaml={setYaml}
          updateHandler={updateHandler}
          deleteHandler={deleteHandler}
          type={'pattern'}
        />
      )}
      <FlipCard
        onClick={() => {
          console.log(gridProps);
          setGridProps(INITIAL_GRID_SIZE);
        }}
        duration={600}
        onShow={() =>
          setTimeout(() => setShowCode((currentCodeVisibilty) => !currentCodeVisibilty), 500)
        }
      >
        {/* FRONT PART */}
        <div>
          <div>
            <div style={{ height: 'max', display: 'flex', justifyContent: 'space-between' }}>
              <Typography
                style={{
                  whiteSpace: 'nowrap',
                  overflow: 'hidden',
                  textOverflow: 'ellipsis',
                  width: '20rem',
                }}
                variant="h6"
                component="div"
              >
                {name}
              </Typography>
              <img className={classes.img} src={`/static/img/${visibility}.svg`} />
            </div>
            <div className={classes.lastRunText}>
              <div>
                {updated_at ? (
                  <Typography
                    variant="caption"
                    style={{
                      fontStyle: 'italic',
                      color: `${
                        theme.palette.type === 'dark' ? 'rgba(255, 255, 255, 0.7)' : '#647881'
                      }`,
                    }}
                  >
                    Modified On: <Moment format="LLL">{updated_at}</Moment>
                  </Typography>
                ) : null}
              </div>
            </div>
          </div>
          <div className={classes.bottomPart}>
            <div className={classes.cardButtons}>
              {canPublishPattern && visibility !== VISIBILITY.PUBLISHED ? (
                <TooltipButton
                  variant="contained"
                  title="Publish"
                  className={classes.testsButton}
                  onClick={(ev) => genericClickHandler(ev, handlePublishModal)}
                  disabled={!CAN(keys.PUBLISH_DESIGN.action, keys.PUBLISH_DESIGN.subject)}
                >
                  <PublicIcon className={classes.iconPatt} />
                  <span className={classes.btnText}> Publish </span>
                </TooltipButton>
              ) : (
                <TooltipButton
                  variant="contained"
                  title="Unpublish"
                  className={classes.testsButton}
                  onClick={(ev) => genericClickHandler(ev, handleUnpublishModal)}
                  disabled={!CAN(keys.UNPUBLISH_DESIGN.action, keys.UNPUBLISH_DESIGN.subject)}
                >
                  <PublicIcon className={classes.iconPatt} />
                  <span className={classes.btnText}> Unpublish </span>
                </TooltipButton>
              )}

              <TooltipButton
                title="Valildate"
                variant="contained"
                className={classes.testsButton}
                onClick={(e) => genericClickHandler(e, handleVerify)}
                disabled={!CAN(keys.VALIDATE_DESIGN.action, keys.VALIDATE_DESIGN.subject)}
              >
                <DoneIcon className={classes.iconPatt} />
                <span className={classes.btnText}> Validate </span>
              </TooltipButton>

              <TooltipButton
                title="Deploy"
                variant="contained"
                onClick={(ev) => genericClickHandler(ev, handleDeploy)}
                className={classes.testsButton}
                disabled={!CAN(keys.DEPLOY_DESIGN.action, keys.DEPLOY_DESIGN.subject)}
              >
                <DoneAllIcon className={classes.iconPatt} />
                <span className={classes.btnText}>Deploy</span>
              </TooltipButton>
              <TooltipButton
                title="Undeploy"
                variant="contained"
                className={classes.undeployButton}
                onClick={(ev) => genericClickHandler(ev, handleUnDeploy)}
                disabled={!CAN(keys.DEPLOY_DESIGN.action, keys.DEPLOY_DESIGN.subject)} // use undeploy keys after it get seeded
              >
                <UndeployIcon fill="#ffffff" className={classes.iconPatt} />
                <span className={classes.btnText}>Undeploy</span>
              </TooltipButton>

              {visibility === VISIBILITY.PRIVATE ? (
                <TooltipButton
                  title="Design"
                  variant="contained"
                  color="primary"
                  onClick={(ev) => genericClickHandler(ev, setSelectedPatterns)}
                  className={classes.testsButton}
                  disabled={!CAN(keys.EDIT_DESIGN.action, keys.EDIT_DESIGN.subject)}
                >
                  <Avatar
                    src="/static/img/pattern_trans.svg"
                    className={classes.iconPatt}
                    imgProps={{ height: '16px', width: '16px' }}
                  />
                  <span className={classes.btnText}> Design </span>
                </TooltipButton>
              ) : (
                <TooltipButton
                  title="Clone"
                  variant="contained"
                  color="primary"
                  onClick={(ev) => genericClickHandler(ev, handleClone)}
                  className={classes.testsButton}
                  // disabled={!CAN(keys.CLONE_DESIGN.action, keys.CLONE_DESIGN.subject)} // TODO: uncomment when key get seeded
                >
                  <CloneIcon fill="#ffffff" className={classes.iconPatt} />
                  <span className={classes.cloneBtnText}> Clone </span>
                </TooltipButton>
              )}

              {userCanEdit && (
                <TooltipButton
                  title="Edit In Configurator"
                  variant="contained"
                  color="primary"
                  onClick={(ev) => genericClickHandler(ev, editInConfigurator)}
                  className={classes.testsButton}
                >
                  <Edit style={{ fill: '#fff' }} className={classes.iconPatt} />
                  <span className={classes.cloneBtnText}> Edit </span>
                </TooltipButton>
              )}
              <TooltipButton
                title="Pattern Information"
                variant="contained"
                color="primary"
                onClick={(ev) => genericClickHandler(ev, handleInfoModal)}
                className={classes.testsButton}
                disabled={!CAN(keys.DETAILS_OF_DESIGN.action, keys.DETAILS_OF_DESIGN.subject)}
              >
                <InfoOutlinedIcon style={{ fill: '#fff' }} className={classes.iconPatt} />
                <span className={classes.btnText}> Info </span>
              </TooltipButton>
            </div>
          </div>
        </div>

        {/* BACK PART */}
        <>
          <Grid
            className={classes.backGrid}
            container
            spacing={1}
            alignContent="space-between"
            alignItems="center"
          >
            <Grid item xs={12} className={classes.yamlDialogTitle}>
              <Typography variant="h6" className={classes.yamlDialogTitleText}>
                {name}
              </Typography>
              <div className={classes.cardHeaderRight}>
                <Link href={`${MESHERY_CLOUD_PROD}/user/${pattern?.user_id}`} target="_blank">
                  <Avatar alt="profile-avatar" src={owner?.avatar_url} />
                </Link>
                <Tooltip title="Enter Fullscreen" arrow interactive placement="top">
                  <IconButton
                    onClick={(ev) =>
                      genericClickHandler(ev, () => {
                        {
                          toggleFullScreen();
                        }
                      })
                    }
                  >
                    {fullScreen ? <FullscreenExit /> : <Fullscreen />}
                  </IconButton>
                </Tooltip>
              </div>
            </Grid>
            <Grid item xs={12} onClick={(ev) => genericClickHandler(ev, () => {})}>
              <Divider variant="fullWidth" light />
              {catalogContentKeys.length === 0 ? (
                <CodeMirror
                  value={showCode && pattern_file}
                  className={fullScreen ? classes.fullScreenCodeMirror : ''}
                  options={{
                    theme: 'material',
                    lineNumbers: true,
                    lineWrapping: true,
                    gutters: ['CodeMirror-lint-markers'],
                    // @ts-ignore
                    lint: true,
                    mode: 'text/x-yaml',
                  }}
                  onChange={(_, data, val) => setYaml(val)}
                />
              ) : (
                catalogContentKeys.map((title, index) => (
                  <>
                    <Typography variant="h6" className={classes.yamlDialogTitleText}>
                      {title}
                    </Typography>
                    <Typography variant="body2">{catalogContentValues[index]}</Typography>
                  </>
                ))
              )}
            </Grid>

            <Grid item xs={8}>
              <div className={classes.lastRunText}>
                <div>
                  {created_at ? (
                    <Typography
                      variant="caption"
                      style={{
                        fontStyle: 'italic',
                        color: `${
                          theme.palette.type === 'dark' ? 'rgba(255, 255, 255, 0.7)' : '#647881'
                        }`,
                      }}
                    >
                      Created at: <Moment format="LLL">{created_at}</Moment>
                    </Typography>
                  ) : null}
                </div>
              </div>
            </Grid>
            <Grid item xs={12}>
              <div className={classes.updateDeleteButtons}>
                {/* Save button */}
                <Tooltip title="Save" arrow interactive placement="bottom">
                  <IconButton
                    disabled={!CAN(keys.EDIT_DESIGN.action, keys.EDIT_DESIGN.subject)}
                    onClick={(ev) => genericClickHandler(ev, updateHandler)}
                  >
                    <Save color="primary" />
                  </IconButton>
                </Tooltip>

                {/* Delete Button */}
                <Tooltip title="Delete" arrow interactive placement="bottom">
                  <IconButton
                    disabled={!CAN(keys.DELETE_A_DESIGN.action, keys.DELETE_A_DESIGN.subject)}
                    onClick={(ev) => genericClickHandler(ev, deleteHandler)}
                  >
                    <DeleteIcon color="primary" />
                  </IconButton>
                </Tooltip>
              </div>
            </Grid>
          </Grid>
        </>
      </FlipCard>
    </>
  );
}

export const MesheryPatternCard = (props) => {
  return (
    <Provider store={store}>
      <MesheryPatternCard_ {...props} />
    </Provider>
  );
};

// @ts-ignore
export default MesheryPatternCard;
