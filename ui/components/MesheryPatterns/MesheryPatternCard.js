import React, { useState } from 'react';
import { Avatar, Divider, Grid, IconButton, Typography, Link, useTheme } from '@layer5/sistent';
import { CustomTooltip, VisibilityChipMenu } from '@layer5/sistent';
import DeleteIcon from '@mui/icons-material/Delete';
import Save from '@mui/icons-material/Save';
import Fullscreen from '@mui/icons-material/Fullscreen';
import Moment from 'react-moment';
import GetAppIcon from '@mui/icons-material/GetApp';
import FlipCard from '../FlipCard';
import { UnControlled as CodeMirror } from 'react-codemirror2';
import FullscreenExit from '@mui/icons-material/FullscreenExit';
import UndeployIcon from '../../public/static/img/UndeployIcon';
import DoneAllIcon from '@mui/icons-material/DoneAll';
import {
  BottomContainer,
  CatalogCardButtons,
  UpdateDeleteButtons,
  CardBackGrid,
  YamlDialogTitleGrid,
  CardHeaderRight,
  GridBtnText,
  GridCloneBtnText,
  StyledCodeMirrorWrapper,
} from './Cards.styles';
import YAMLDialog from '../YamlDialog';
import PublicIcon from '@mui/icons-material/Public';
import TooltipButton from '@/utils/TooltipButton';
import CloneIcon from '../../public/static/img/CloneIcon';
import { useRouter } from 'next/router';
import { Edit, Lock, Public } from '@mui/icons-material';
import InfoOutlinedIcon from '@mui/icons-material/InfoOutlined';
import { MESHERY_CLOUD_PROD } from '../../constants/endpoints';
import { useGetUserByIdQuery } from '../../rtk-query/user';
import { Provider } from 'react-redux';
import { store } from '../../store';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import ActionButton from './ActionButton';
import DryRunIcon from '@/assets/icons/DryRunIcon';
import CheckIcon from '@/assets/icons/CheckIcon';
import { VISIBILITY } from '@/utils/Enum';
import PatternIcon from '@/assets/icons/Pattern';
import { iconLarge, iconMedium } from 'css/icons.styles';
import { UsesSistent } from '../SistentWrapper';
import { VIEW_VISIBILITY } from '../Modals/Information/InfoModal';
const INITIAL_GRID_SIZE = { xl: 4, md: 6, xs: 12 };

function MesheryPatternCard_({
  id,
  name,
  updated_at,
  created_at,
  pattern_file,
  handleVerify,
  handleDryRun,
  handleUnpublishModal,
  handleDeploy,
  handleUnDeploy,
  handleDownload,
  updateHandler,
  deleteHandler,
  handleClone,
  setSelectedPatterns,
  setYaml,
  description = {},
  visibility,
  user,
  pattern,
  handleInfoModal,
  hideVisibility = false,
  isReadOnly = false,
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
  const theme = useTheme();

  const editInConfigurator = () => {
    router.push('/configuration/designs/configurator?design_id=' + id);
  };
  const isOwner = user?.user_id == pattern?.user_id;
  const userCanEdit = CAN(keys.EDIT_DESIGN.action, keys.EDIT_DESIGN.subject) || isOwner;

  const formatPatternFile = (file) => {
    try {
      const jsonData = JSON.parse(file);
      return JSON.stringify(jsonData, null, 1);
    } catch (err) {
      return file;
    }
  };

  const formatted_pattern_file = formatPatternFile(pattern_file);
  return (
    <UsesSistent>
      {fullScreen && (
        <YAMLDialog
          fullScreen={fullScreen}
          name={name}
          toggleFullScreen={toggleFullScreen}
          config_file={formatted_pattern_file}
          setYaml={setYaml}
          updateHandler={updateHandler}
          deleteHandler={deleteHandler}
          type={'pattern'}
          isReadOnly={isReadOnly}
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
                  width: '100%',
                }}
                variant="h6"
                component="div"
              >
                {name}
              </Typography>
              {hideVisibility ? (
                <PatternIcon {...iconLarge} color={true} />
              ) : (
                <div>
                  <VisibilityChipMenu
                    value={visibility}
                    onChange={() => {}}
                    enabled={false}
                    options={[
                      [VIEW_VISIBILITY.PUBLIC, Public],
                      [VIEW_VISIBILITY.PRIVATE, Lock],
                    ]}
                  />
                </div>
              )}
            </div>
            <div style={{ marginRight: '0.5rem' }}>
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
          <BottomContainer>
            <CatalogCardButtons>
              {visibility === VISIBILITY.PUBLISHED && (
                <TooltipButton
                  variant="outlined"
                  title="Unpublish"
                  style={{
                    padding: '6px 9px',
                    borderRadius: '8px',
                  }}
                  onClick={(ev) => genericClickHandler(ev, handleUnpublishModal)}
                  disabled={!CAN(keys.UNPUBLISH_DESIGN.action, keys.UNPUBLISH_DESIGN.subject)}
                >
                  <PublicIcon style={iconMedium} />
                  <GridBtnText> Unpublish </GridBtnText>
                </TooltipButton>
              )}
              <ActionButton
                defaultActionClick={(e) => genericClickHandler(e, handleVerify)}
                options={[
                  {
                    label: 'Validate',
                    icon: <CheckIcon style={iconMedium} />,
                    onClick: (e) => genericClickHandler(e, handleVerify),
                    disabled: !CAN(keys.VALIDATE_DESIGN.action, keys.VALIDATE_DESIGN.subject),
                  },
                  {
                    label: 'Dry Run',
                    icon: <DryRunIcon style={iconMedium} />,
                    onClick: (e) => genericClickHandler(e, handleDryRun),
                    disabled: !CAN(keys.VALIDATE_DESIGN.action, keys.VALIDATE_DESIGN.subject),
                  },
                  {
                    label: 'Deploy',
                    icon: <DoneAllIcon style={iconMedium} />,
                    onClick: (e) => genericClickHandler(e, handleDeploy),
                    disabled: !CAN(keys.DEPLOY_DESIGN.action, keys.DEPLOY_DESIGN.subject),
                  },
                  {
                    label: 'Undeploy',
                    icon: <UndeployIcon fill={'currentColor'} style={iconMedium} />,
                    onClick: (e) => genericClickHandler(e, handleUnDeploy),
                    disabled: !CAN(keys.DEPLOY_DESIGN.action, keys.DEPLOY_DESIGN.subject),
                  },
                ]}
              />
              <TooltipButton
                title="Download"
                variant="contained"
                color="primary"
                onClick={handleDownload}
                style={{
                  padding: '6px 9px',
                  borderRadius: '8px',
                }}
              >
                <GetAppIcon
                  fill={theme.palette.background.constant.white}
                  data-cy="download-button"
                />
                <GridBtnText> Download </GridBtnText>
              </TooltipButton>
              {visibility === VISIBILITY.PRIVATE && userCanEdit ? (
                <TooltipButton
                  title="Design"
                  variant="contained"
                  color="primary"
                  onClick={(ev) => genericClickHandler(ev, setSelectedPatterns)}
                  style={{
                    padding: '6px 9px',
                    borderRadius: '8px',
                  }}
                  disabled={!CAN(keys.EDIT_DESIGN.action, keys.EDIT_DESIGN.subject)}
                >
                  <img
                    src="/static/img/pattern_trans.svg"
                    style={{ borderRadius: '50%', ...iconMedium }}
                    // imgProps={{ height: '16px', width: '16px' }}
                  />
                  <GridBtnText> Design </GridBtnText>
                </TooltipButton>
              ) : (
                <TooltipButton
                  title="Clone"
                  variant="contained"
                  color="primary"
                  onClick={(ev) => genericClickHandler(ev, handleClone)}
                  style={{
                    padding: '6px 9px',
                    borderRadius: '8px',
                  }}
                  disabled={!CAN(keys.CLONE_DESIGN.action, keys.CLONE_DESIGN.subject)}
                >
                  <CloneIcon fill={theme.palette.background.constant.white} style={iconMedium} />
                  <GridCloneBtnText> Clone </GridCloneBtnText>
                </TooltipButton>
              )}

              {userCanEdit && (
                <TooltipButton
                  title="Edit In Configurator"
                  variant="contained"
                  color="primary"
                  onClick={(ev) => genericClickHandler(ev, editInConfigurator)}
                  disabled={!CAN(keys.EDIT_DESIGN.action, keys.EDIT_DESIGN.subject)}
                  style={{
                    padding: '6px 9px',
                    borderRadius: '8px',
                  }}
                >
                  <Edit style={{ fill: theme.palette.background.constant.white, ...iconMedium }} />
                  <GridCloneBtnText> Edit </GridCloneBtnText>
                </TooltipButton>
              )}
              <TooltipButton
                title="Pattern Information"
                variant="contained"
                color="primary"
                onClick={(ev) => genericClickHandler(ev, handleInfoModal)}
                style={{
                  padding: '6px 9px',
                  borderRadius: '8px',
                }}
                disabled={!CAN(keys.DETAILS_OF_DESIGN.action, keys.DETAILS_OF_DESIGN.subject)}
              >
                <InfoOutlinedIcon
                  style={{ fill: theme.palette.background.constant.white, ...iconMedium }}
                />
                <GridBtnText> Info </GridBtnText>
              </TooltipButton>
            </CatalogCardButtons>
          </BottomContainer>
        </div>

        {/* BACK PART */}
        <>
          <CardBackGrid container spacing={1} alignContent="space-between" alignItems="center">
            <YamlDialogTitleGrid item xs={12}>
              <Typography variant="h6">{name}</Typography>
              <CardHeaderRight>
                <Link href={`${MESHERY_CLOUD_PROD}/user/${pattern?.user_id}`} target="_blank">
                  <Avatar alt="profile-avatar" src={owner?.avatar_url} />
                </Link>
                <CustomTooltip title="Enter Fullscreen" arrow interactive placement="top">
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
                </CustomTooltip>
              </CardHeaderRight>
            </YamlDialogTitleGrid>
            <Grid item xs={12} onClick={(ev) => genericClickHandler(ev, () => {})}>
              <Divider variant="fullWidth" light />
              {catalogContentKeys.length === 0 ? (
                <StyledCodeMirrorWrapper fullScreen={fullScreen}>
                  <CodeMirror
                    value={showCode && formatted_pattern_file}
                    options={{
                      theme: 'material',
                      lineNumbers: true,
                      lineWrapping: true,
                      gutters: ['CodeMirror-lint-markers'],
                      // @ts-ignore
                      lint: true,
                      mode: 'text/x-yaml',
                      readOnly: isReadOnly,
                    }}
                    onChange={(_, data, val) => setYaml(val)}
                  />
                </StyledCodeMirrorWrapper>
              ) : (
                catalogContentKeys.map((title, index) => (
                  <>
                    <Typography variant="h6">{title}</Typography>
                    <Typography variant="body2">{catalogContentValues[index]}</Typography>
                  </>
                ))
              )}
            </Grid>

            <Grid item xs={8}>
              <div style={{ marginRight: '0.5rem' }}>
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
              {isReadOnly ? null : (
                <UpdateDeleteButtons>
                  {/* Save button */}
                  <CustomTooltip title="Save" arrow interactive placement="bottom">
                    <IconButton
                      disabled={!CAN(keys.EDIT_DESIGN.action, keys.EDIT_DESIGN.subject)}
                      onClick={(ev) => genericClickHandler(ev, updateHandler)}
                    >
                      <Save fill={theme.palette.background.constant.white} />
                    </IconButton>
                  </CustomTooltip>

                  {/* Delete Button */}
                  <CustomTooltip title="Delete" arrow interactive placement="bottom">
                    <IconButton
                      disabled={!CAN(keys.DELETE_A_DESIGN.action, keys.DELETE_A_DESIGN.subject)}
                      onClick={(ev) => genericClickHandler(ev, deleteHandler)}
                    >
                      <DeleteIcon fill={theme.palette.background.constant.white} />
                    </IconButton>
                  </CustomTooltip>
                </UpdateDeleteButtons>
              )}
            </Grid>
          </CardBackGrid>
        </>
      </FlipCard>
    </UsesSistent>
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
