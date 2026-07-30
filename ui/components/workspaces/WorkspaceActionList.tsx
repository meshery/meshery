import {
  AccessTimeFilledIcon,
  CustomTooltip,
  DeleteIcon,
  EditIcon,
  GroupAddIcon,
  IconButton,
  ListItemIcon,
  Menu,
  MenuItem,
  MoreVertIcon,
  useHasPermission,
  useTheme,
  useWindowDimensions,
} from '@sistent/sistent';
import { useState } from 'react';
import { TableIconsContainer, IconWrapper } from './styles';
import { iconMedium } from 'css/icons.styles';
import { WORKSPACE_ACTION_TYPES } from '.';
import { Keys } from '@meshery/schemas/permissions';

const WorkspaceActionList = ({
  handleTeamsModalOpen,
  handleActivityModalOpen,
  handleWorkspaceModalOpen,
  handleDeleteWorkspaceConfirm,
  workspaceId,
  workspaceName,
  selectedWorkspace,
}) => {
  const { width } = useWindowDimensions();
  const isMobile = width < 1024;
  const theme = useTheme();
  const canEditWorkspace = useHasPermission(Keys.WorkspaceManagementEditWorkspace);
  const canDeleteWorkspace = useHasPermission(Keys.WorkspaceManagementDeleteWorkspace);

  const [anchorEl, setAnchorEl] = useState(null);
  const open = Boolean(anchorEl);

  const handleClick = (event) => {
    event.stopPropagation();
    event.preventDefault();
    setAnchorEl(event.currentTarget);
  };

  const handleClose = (event) => {
    event.stopPropagation();
    event.preventDefault();
    setAnchorEl(null);
  };

  const actionItems = [
    {
      key: 'assign-teams',
      label: 'Assign Teams',
      icon: <GroupAddIcon style={{ color: theme.palette.icon.default, ...iconMedium }} />,
      onClick: (e) => handleTeamsModalOpen(e, workspaceId, workspaceName),
    },
    {
      key: 'recent-activity',
      label: 'Recent Activity',
      icon: <AccessTimeFilledIcon style={{ color: theme.palette.icon.default, ...iconMedium }} />,
      onClick: (e) => handleActivityModalOpen(e, workspaceId, workspaceName),
    },
    {
      key: 'edit-workspace',
      label: 'Edit Workspace',
      icon: <EditIcon style={{ fill: theme.palette.icon.default, ...iconMedium }} />,
      onClick: (e) => handleWorkspaceModalOpen(e, WORKSPACE_ACTION_TYPES.EDIT, selectedWorkspace),
      disabled: !canEditWorkspace,
    },
    {
      key: 'delete-workspace',
      label: 'Delete Workspace',
      icon: <DeleteIcon style={{ fill: theme.palette.icon.default, ...iconMedium }} />,
      onClick: (e) => handleDeleteWorkspaceConfirm(e, selectedWorkspace),
      disabled: !canDeleteWorkspace,
    },
  ];

  return (
    <TableIconsContainer>
      <IconWrapper>
        {isMobile ? (
          <>
            <IconButton aria-label="more" onClick={handleClick}>
              <MoreVertIcon />
            </IconButton>
            <Menu anchorEl={anchorEl} open={open} onClose={handleClose}>
              {actionItems.map(({ key, label, icon, onClick, disabled }) => (
                <MenuItem
                  key={key}
                  onClick={(e) => {
                    onClick(e);
                    handleClose(e);
                  }}
                  disabled={disabled}
                >
                  <ListItemIcon>{icon}</ListItemIcon>
                  {label}
                </MenuItem>
              ))}
            </Menu>
          </>
        ) : (
          <>
            {actionItems.map(({ key, label, icon, onClick, disabled }) => (
              <CustomTooltip title={label} key={key}>
                <IconButton aria-label={key} onClick={(e) => onClick(e)} disabled={disabled}>
                  {icon}
                </IconButton>
              </CustomTooltip>
            ))}
          </>
        )}
      </IconWrapper>
    </TableIconsContainer>
  );
};

export default WorkspaceActionList;
