import React, { useState } from 'react';
import {
  Divider,
  Grid,
  IconButton,
  Typography,
  Tooltip,
  Link,
  Avatar,
  useTheme,
} from '@layer5/sistent';
import DeleteIcon from '@mui/icons-material/Delete';
import Fullscreen from '@mui/icons-material/Fullscreen';
import Save from '@mui/icons-material/Save';
import Moment from 'react-moment';
import FlipCard from '../FlipCard';
import { UnControlled as CodeMirror } from 'react-codemirror2';
import FullscreenExit from '@mui/icons-material/FullscreenExit';
import {
  BottomContainer,
  CardBackGrid,
  CatalogCardButtons,
  UpdateDeleteButtons,
  YamlDialogTitleGrid,
  GridBtnText,
  GridCloneBtnText,
  CardHeaderRight,
  StyledCodeMirrorWrapper,
} from '../MesheryPatterns/Cards.styles';
import YAMLDialog from '../YamlDialog';
import CloneIcon from '../../public/static/img/CloneIcon';
import PublicIcon from '@mui/icons-material/Public';
import TooltipButton from '../../utils/TooltipButton.js';
import { VISIBILITY } from '../../utils/Enum';
import GetAppIcon from '@mui/icons-material/GetApp';
import InfoOutlinedIcon from '@mui/icons-material/InfoOutlined';
import { Provider } from 'react-redux';
import { store } from '../../store';
import { useGetUserByIdQuery } from '../../rtk-query/user';
import { MESHERY_CLOUD_PROD } from '../../constants/endpoints';
import { keys } from '@/utils/permission_constants';
import CAN from '@/utils/can';
import { VisibilityChipMenu } from '@layer5/sistent';
import { VIEW_VISIBILITY } from '../Modals/Information/InfoModal';
import { Public, Lock } from '@mui/icons-material';
import { iconMedium } from 'css/icons.styles';

const INITIAL_GRID_SIZE = { xl: 4, md: 6, xs: 12 };

function FiltersCard_({
  name,
  updated_at,
  created_at,
  filter_resource,
  handleClone,
  handleDownload,
  deleteHandler,
  setYaml,
  description = {},
  visibility,
  handlePublishModal,
  handleUnpublishModal,
  updateHandler,
  canPublishFilter = false,
  handleInfoModal,
  ownerId,
}) {
  const genericClickHandler = (ev, fn) => {
    ev.stopPropagation();
    fn(ev);
  };
  const [gridProps, setGridProps] = useState(INITIAL_GRID_SIZE);
  const [fullScreen, setFullScreen] = useState(false);
  const [showCode, setShowCode] = useState(false);

  const { data: owner } = useGetUserByIdQuery(ownerId || '');

  const toggleFullScreen = () => {
    setFullScreen(!fullScreen);
  };

  const catalogContentKeys = Object.keys(description);
  const catalogContentValues = Object.values(description);
  const theme = useTheme();

  return (
    <>
      {fullScreen && (
        <YAMLDialog
          fullScreen={fullScreen}
          name={name}
          toggleFullScreen={toggleFullScreen}
          config_file={filter_resource}
          setYaml={setYaml}
          deleteHandler={deleteHandler}
          updateHandler={updateHandler}
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
        <>
          <div>
            <div style={{ display: 'flex', justifyContent: 'space-between' }}>
              <Typography variant="h6" component="div">
                {name}
              </Typography>
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
              </div>{' '}
            </div>
            <div style={{ marginRight: '0.5rem' }}>
              <div>
                {updated_at ? (
                  <Typography variant="caption" style={{ fontStyle: 'italic' }}>
                    Modified On: <Moment format="LLL">{updated_at}</Moment>
                  </Typography>
                ) : null}
              </div>
            </div>
          </div>
          <BottomContainer>
            <CatalogCardButtons>
              {canPublishFilter && visibility !== VISIBILITY.PUBLISHED ? (
                <TooltipButton
                  variant="outlined"
                  title="Publish"
                  style={{
                    padding: '6px 9px',
                    borderRadius: '8px',
                  }}
                  onClick={(ev) => genericClickHandler(ev, handlePublishModal)}
                  disabled={!CAN(keys.PUBLISH_WASM_FILTER.action, keys.PUBLISH_WASM_FILTER.subject)}
                >
                  <PublicIcon style={iconMedium} />
                  <> Publish </>
                </TooltipButton>
              ) : (
                <TooltipButton
                  variant="outlined"
                  title="Unpublish"
                  style={{
                    padding: '6px 9px',
                    borderRadius: '8px',
                  }}
                  onClick={(ev) => genericClickHandler(ev, handleUnpublishModal)}
                  disabled={
                    !CAN(keys.UNPUBLISH_WASM_FILTER.action, keys.UNPUBLISH_WASM_FILTER.subject)
                  }
                >
                  <PublicIcon style={iconMedium} />
                  <GridBtnText> Unpublish </GridBtnText>
                </TooltipButton>
              )}
              <TooltipButton
                title="Download"
                variant="contained"
                color="primary"
                onClick={handleDownload}
                style={{
                  padding: '6px 9px',
                  borderRadius: '8px',
                }}
                disabled={
                  !CAN(keys.DOWNLOAD_A_WASM_FILTER.action, keys.DOWNLOAD_A_WASM_FILTER.subject)
                }
              >
                <GetAppIcon fill={theme.palette.background.constant.white} style={iconMedium} />
                <GridBtnText>Download</GridBtnText>
              </TooltipButton>

              {visibility === VISIBILITY.PUBLISHED ? (
                <TooltipButton
                  title="Clone"
                  variant="contained"
                  color="primary"
                  style={{
                    padding: '6px 9px',
                    borderRadius: '8px',
                  }}
                  onClick={(ev) => genericClickHandler(ev, handleClone)}
                  disabled={!CAN(keys.CLONE_WASM_FILTER.action, keys.CLONE_WASM_FILTER.subject)}
                >
                  <CloneIcon fill={theme.palette.background.constant.white} style={iconMedium} />
                  <GridCloneBtnText>Clone</GridCloneBtnText>
                </TooltipButton>
              ) : null}
              <TooltipButton
                title="Filter Information"
                variant="contained"
                color="primary"
                onClick={(ev) => genericClickHandler(ev, handleInfoModal)}
                style={{
                  padding: '6px 9px',
                  borderRadius: '8px',
                }}
                disabled={
                  !CAN(keys.DETAILS_OF_WASM_FILTER.action, keys.DETAILS_OF_WASM_FILTER.subject)
                }
              >
                <InfoOutlinedIcon
                  fill={theme.palette.background.constant.white}
                  style={iconMedium}
                />
                <GridBtnText> Info </GridBtnText>
              </TooltipButton>
            </CatalogCardButtons>
          </BottomContainer>
        </>

        {/* BACK PART */}
        <>
          <CardBackGrid container spacing={1} alignContent="space-between" alignItems="center">
            <YamlDialogTitleGrid item xs={12}>
              <Typography variant="h6">{name}</Typography>
              <CardHeaderRight>
                <Link href={`${MESHERY_CLOUD_PROD}/user/${ownerId}`} target="_blank">
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
              </CardHeaderRight>
            </YamlDialogTitleGrid>
            <Grid item xs={12} onClick={(ev) => genericClickHandler(ev, () => {})}>
              <Divider variant="fullWidth" light />

              {catalogContentKeys.length === 0 ? (
                <StyledCodeMirrorWrapper fullScreen={fullScreen}>
                  <CodeMirror
                    value={showCode && filter_resource}
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
                    <Typography variant="caption" style={{ fontStyle: 'italic' }}>
                      Created at: <Moment format="LLL">{created_at}</Moment>
                    </Typography>
                  ) : null}
                </div>
              </div>
            </Grid>

            <Grid item xs={12}>
              <UpdateDeleteButtons>
                {/* Save button */}
                <Tooltip title="Save" arrow interactive placement="bottom">
                  <IconButton
                    disabled={!CAN(keys.EDIT_WASM_FILTER.action, keys.EDIT_WASM_FILTER.subject)}
                    onClick={(ev) => genericClickHandler(ev, updateHandler)}
                  >
                    <Save fill={theme.palette.icon.default} />
                  </IconButton>
                </Tooltip>

                {/* Delete Button */}
                <Tooltip title="Delete" arrow interactive placement="bottom">
                  <IconButton
                    disabled={!CAN(keys.DELETE_WASM_FILTER.action, keys.DELETE_WASM_FILTER.subject)}
                    onClick={(ev) => genericClickHandler(ev, deleteHandler)}
                  >
                    <DeleteIcon fill={theme.palette.icon.default} />
                  </IconButton>
                </Tooltip>
              </UpdateDeleteButtons>
            </Grid>
          </CardBackGrid>
        </>
      </FlipCard>
    </>
  );
}

export const FiltersCard = (props) => {
  return (
    <Provider store={store}>
      <FiltersCard_ {...props} />
    </Provider>
  );
};

// @ts-ignore
export default FiltersCard;
