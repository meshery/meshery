import {
  AccessTimeFilledIcon,
  BottomSheet,
  CustomTooltip,
  DeleteIcon,
  EditIcon,
  GroupAddIcon,
  IconButton,
  ListItemIcon,
  MenuItem,
  MenuList,
  MoreVertIcon,
  useMediaQuery,
  useTheme,
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
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'));

  const [open, setOpen] = useState(false);

  const handleClick = (event) => {
    event.stopPropagation();
    event.preventDefault();
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
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
      permissionKey: Keys.WorkspaceManagementEditWorkspace,
    },
    {
      key: 'delete-workspace',
      label: 'Delete Workspace',
      icon: <DeleteIcon style={{ fill: theme.palette.icon.default, ...iconMedium }} />,
      onClick: (e) => handleDeleteWorkspaceConfirm(e, selectedWorkspace),
      permissionKey: Keys.WorkspaceManagementDeleteWorkspace,
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
            <BottomSheet open={open} onClose={handleClose} title="Workspace Actions">
              <MenuList disablePadding>
                {actionItems.map(({ key, label, icon, onClick, permissionKey }) => (
                  <MenuItem
                    key={key}
                    onClick={(e) => {
                      onClick(e);
                      handleClose();
                    }}
                    permissionKey={permissionKey}
                  >
                    <ListItemIcon>{icon}</ListItemIcon>
                    {label}
                  </MenuItem>
                ))}
              </MenuList>
            </BottomSheet>
          </>
        ) : (
          <>
            {actionItems.map(({ key, label, icon, onClick, permissionKey }) => (
              <CustomTooltip title={label} key={key}>
                <IconButton
                  aria-label={key}
                  onClick={(e) => onClick(e)}
                  permissionKey={permissionKey}
                >
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
