// @ts-nocheck
import React, { useState } from 'react';
import {
  Avatar,
  Divider,
  IconButton,
  Typography,
  Link,
  useTheme,
  DeleteIcon,
  GetAppIcon,
  DoneAllIcon,
  Box,
} from '@sistent/sistent';
import { CustomTooltip, VisibilityChipMenu } from '@sistent/sistent';
import { Save, Fullscreen, FullscreenExit, Public as PublicIcon } from '@/components/icons';
import Moment from 'react-moment';
import { UnControlled as CodeMirror } from '../CodeMirror';
import UndeployIcon from '../../public/static/img/UndeployIcon';
import {
  CatalogCardButtons,
  UpdateDeleteButtons,
  CardHeaderRight,
  GridBtnText,
  GridCloneBtnText,
  StyledCodeMirrorWrapper,
} from './Cards.styles';
import ResourceCard from '../shared/Card/ResourceCard';
import YAMLDialog from '../shared/Modal/YamlDialog';
import TooltipButton from '@/utils/TooltipButton';
import CloneIcon from '../../public/static/img/CloneIcon';
import { useRouter } from 'next/router';
import { Edit, Lock, Public, InfoOutlined as InfoOutlinedIcon } from '@/components/icons';
import { MESHERY_CLOUD_PROD } from '../../constants/endpoints';
import { useGetUserByIdQuery } from '../../rtk-query/user';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import ActionButton from './ActionButton';
import DryRunIcon from '@/assets/icons/DryRunIcon';
import CheckIcon from '@/assets/icons/CheckIcon';
import { VISIBILITY } from '@/utils/Enum';
import PatternIcon from '@/assets/icons/Pattern';
import { iconLarge, iconMedium } from '@/css/icons.styles';
import { VIEW_VISIBILITY } from '@/components/shared/Modal/InfoModal';
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
    } catch {
      return file;
    }
  };

  const formatted_pattern_file = formatPatternFile(pattern_file);
  return (
    <div data-testid="meshery-pattern-card-item">
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
      <ResourceCard
        title={name}
        isFlippable={true}
        onFlip={() => {
          setGridProps(INITIAL_GRID_SIZE);
          setShowCode(true); // Ensure code is shown when flipped
        }}
        headerRight={
          hideVisibility ? (
            <PatternIcon {...iconLarge} color={true} data-testid="pattern-icon" />
          ) : (
            <div data-testid="visibility-chip-menu">
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
          )
        }
        subtitle={
          updated_at ? (
            <Typography
              variant="caption"
              sx={{
                fontStyle: 'italic',
                color: (theme) =>
                  theme.palette.mode === 'dark' ? 'rgba(255, 255, 255, 0.7)' : '#647881',
              }}
              data-testid="pattern-card-modified-on"
            >
              Modified On: <Moment format="LLL">{updated_at}</Moment>
            </Typography>
          ) : null
        }
        footer={
          <CatalogCardButtons data-testid="pattern-card-actions">
            {visibility === VISIBILITY.PUBLISHED && (
              <TooltipButton
                variant="outlined"
                title="Unpublish"
                sx={{
                  padding: '6px 9px',
                  borderRadius: '8px',
                }}
                onClick={(ev) => genericClickHandler(ev, handleUnpublishModal)}
                disabled={!CAN(keys.UNPUBLISH_DESIGN.action, keys.UNPUBLISH_DESIGN.subject)}
                data-testid="pattern-btn-unpublish"
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
                  'data-testid': 'pattern-btn-validate',
                },
                {
                  label: 'Dry Run',
                  icon: <DryRunIcon style={iconMedium} />,
                  onClick: (e) => genericClickHandler(e, handleDryRun),
                  disabled: !CAN(keys.VALIDATE_DESIGN.action, keys.VALIDATE_DESIGN.subject),
                  'data-testid': 'pattern-btn-dryrun',
                },
                {
                  label: 'Deploy',
                  icon: <DoneAllIcon style={iconMedium} />,
                  onClick: (e) => genericClickHandler(e, handleDeploy),
                  disabled: !CAN(keys.DEPLOY_DESIGN.action, keys.DEPLOY_DESIGN.subject),
                  'data-testid': 'pattern-btn-deploy',
                },
                {
                  label: 'Undeploy',
                  icon: <UndeployIcon fill={'currentColor'} style={iconMedium} />,
                  onClick: (e) => genericClickHandler(e, handleUnDeploy),
                  disabled: !CAN(keys.DEPLOY_DESIGN.action, keys.DEPLOY_DESIGN.subject),
                  'data-testid': 'pattern-btn-undeploy',
                },
              ]}
              data-testid="pattern-btn-action-dropdown"
            />
            <TooltipButton
              title="Download"
              variant="contained"
              color="primary"
              onClick={handleDownload}
              sx={{
                padding: '6px 9px',
                borderRadius: '8px',
              }}
              data-testid="pattern-btn-download"
            >
              <GetAppIcon fill={theme.palette.common.white} data-cy="download-button" />
              <GridBtnText> Download </GridBtnText>
            </TooltipButton>
            {visibility === VISIBILITY.PRIVATE && userCanEdit ? (
              <TooltipButton
                title="Design"
                variant="contained"
                color="primary"
                onClick={(ev) => genericClickHandler(ev, setSelectedPatterns)}
                sx={{
                  padding: '6px 9px',
                  borderRadius: '8px',
                }}
                disabled={!CAN(keys.EDIT_DESIGN.action, keys.EDIT_DESIGN.subject)}
                data-testid="pattern-btn-design"
              >
                <img
                  src="/static/img/pattern_trans.svg"
                  style={{ borderRadius: '50%', ...iconMedium }}
                />
                <GridBtnText> Design </GridBtnText>
              </TooltipButton>
            ) : (
              <TooltipButton
                title="Clone"
                variant="contained"
                color="primary"
                onClick={(ev) => genericClickHandler(ev, handleClone)}
                sx={{
                  padding: '6px 9px',
                  borderRadius: '8px',
                }}
                disabled={!CAN(keys.CLONE_DESIGN.action, keys.CLONE_DESIGN.subject)}
                data-testid="pattern-btn-clone"
              >
              <CloneIcon fill={theme.palette.common.white} style={iconMedium} />
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
                sx={{ padding: '6px 9px', borderRadius: '8px' }}
                data-testid="pattern-btn-edit"
              >
              <Edit style={{ fill: theme.palette.common.white, ...iconMedium }} />
                <GridCloneBtnText> Edit </GridCloneBtnText>
              </TooltipButton>
            )}
            <TooltipButton
              title="Pattern Information"
              variant="contained"
              color="primary"
              onClick={(ev) => genericClickHandler(ev, handleInfoModal)}
              sx={{
                padding: '6px 9px',
                borderRadius: '8px',
              }}
              disabled={!CAN(keys.DETAILS_OF_DESIGN.action, keys.DETAILS_OF_DESIGN.subject)}
              data-testid="pattern-btn-info"
            >
              <InfoOutlinedIcon style={{ fill: theme.palette.common.white, ...iconMedium }} />
              <GridBtnText> Info </GridBtnText>
            </TooltipButton>
          </CatalogCardButtons>
        }
        backHeaderRight={
          <CardHeaderRight>
            <Link href={`${MESHERY_CLOUD_PROD}/user/${pattern?.user_id}`} target="_blank">
              <Avatar alt="profile-avatar" src={owner?.avatar_url} />
            </Link>
            <CustomTooltip title="Enter Fullscreen" arrow interactive placement="top">
              <IconButton
                onClick={(ev) =>
                  genericClickHandler(ev, () => {
                    toggleFullScreen();
                  })
                }
              >
                {fullScreen ? <FullscreenExit /> : <Fullscreen />}
              </IconButton>
            </CustomTooltip>
          </CardHeaderRight>
        }
        backSubtitle={
          created_at ? (
            <Typography
              variant="caption"
              sx={{
                fontStyle: 'italic',
                color: (theme) =>
                  theme.palette.mode === 'dark' ? 'rgba(255, 255, 255, 0.7)' : '#647881',
              }}
              data-testid="pattern-card-created-at"
            >
              Created at: <Moment format="LLL">{created_at}</Moment>
            </Typography>
          ) : null
        }
        backFooter={
          !isReadOnly && (
            <UpdateDeleteButtons data-testid="pattern-card-save-delete-buttons">
              <CustomTooltip title="Save" arrow interactive placement="bottom">
                <IconButton
                  disabled={!CAN(keys.EDIT_DESIGN.action, keys.EDIT_DESIGN.subject)}
                  onClick={(ev) => genericClickHandler(ev, updateHandler)}
                  data-testid="pattern-btn-save"
                >
                  <Save fill={theme.palette.common.white} />
                </IconButton>
              </CustomTooltip>

              <CustomTooltip title="Delete" arrow interactive placement="bottom">
                <IconButton
                  disabled={!CAN(keys.DELETE_A_DESIGN.action, keys.DELETE_A_DESIGN.subject)}
                  onClick={(ev) => genericClickHandler(ev, deleteHandler)}
                  data-testid="pattern-btn-delete"
                >
                  <DeleteIcon fill={theme.palette.common.white} />
                </IconButton>
              </CustomTooltip>
            </UpdateDeleteButtons>
          )
        }
        backContent={
          <Box
            sx={{ flexGrow: 1 }}
            onClick={(ev) => genericClickHandler(ev, () => {})}
            data-testid="pattern-card-back"
          >
            <Divider sx={{ mb: 1 }} light />
            {catalogContentKeys.length === 0 ? (
              <StyledCodeMirrorWrapper fullScreen={fullScreen}>
                <CodeMirror
                  value={showCode && formatted_pattern_file}
                  options={{
                    theme: 'material',
                    lineNumbers: true,
                    lineWrapping: true,
                    gutters: ['CodeMirror-lint-markers'],
                    lint: true,
                    mode: 'text/x-yaml',
                    readOnly: isReadOnly,
                  }}
                  onChange={(_, data, val) => setYaml(val)}
                />
              </StyledCodeMirrorWrapper>
            ) : (
              catalogContentKeys.map((title, index) => (
                <Box key={index} mb={1}>
                  <Typography variant="h6">{title}</Typography>
                  <Typography variant="body2">{catalogContentValues[index]}</Typography>
                </Box>
              ))
            )}
          </Box>
        }
      />
    </div>
  );
}

export const MesheryPatternCard = (props) => {
  return <MesheryPatternCard_ {...props} />;
};

export default MesheryPatternCard;
