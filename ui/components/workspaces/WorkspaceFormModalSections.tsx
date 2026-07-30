/**
 * Drawer rail sections used by the workspace navigation modal: a nav-item
 * row plus the workspaces list (an "All Workspaces" entry and one entry per
 * workspace in the active organisation).
 */
import React, { useContext, FC, useRef } from 'react';
import {
  CustomTooltip,
  ListItem,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  WorkspaceIcon,
} from '@sistent/sistent';
import { useTheme } from '@/theme';
import { iconSmall } from 'css/icons.styles';
import { WorkspaceModalContext } from '@/utils/context/WorkspaceModalContextProvider';
import { useIsTextTruncated } from '@/utils/hooks/useIsTextTruncated';

export type NavConfigItem = {
  id: string;
  label: string;
  icon: React.ReactNode;
  content: React.ReactNode;
  enabled?: boolean;
};

type NavItemProps = {
  item: NavConfigItem;
  open: boolean;
  selectedId: string;
  onSelect: (id: string) => void;
};

export const NavItem: FC<NavItemProps> = ({ item, open, selectedId, onSelect }) => {
  const { setMultiSelectedContent } = useContext(WorkspaceModalContext);
  return (
    <CustomTooltip title={item.label} disableHoverListener={open} placement="right">
      <ListItem disablePadding sx={{ display: 'block' }}>
        <ListItemButton
          selected={selectedId === item.id}
          onClick={() => {
            setMultiSelectedContent([]);
            onSelect(item.id);
          }}
          sx={{
            minHeight: 48,
            px: 2.5,
            justifyContent: open ? 'initial' : 'center',
          }}
        >
          <ListItemIcon
            sx={{
              minWidth: 0,
              justifyContent: 'center',
              mr: open ? 3 : 'auto',
            }}
          >
            {item.icon}
          </ListItemIcon>
          <ListItemText primary={item.label} sx={{ opacity: open ? 1 : 0 }} />
        </ListItemButton>
      </ListItem>
    </CustomTooltip>
  );
};

type WorkspaceItem = {
  id: string;
  name: string;
  icon: React.ReactNode;
};

type WorkspacesSectionProps = {
  open: boolean;
  selectedId: string;
  onSelect: (id: string) => void;
  workspacesData?: { workspaces?: Array<{ id: string; name: string }> };
  isLoading: boolean;
};

type WorkspaceListItemProps = {
  workspace: WorkspaceItem;
  open: boolean;
  selected: boolean;
  onSelect: () => void;
  backgroundColor: string | undefined;
};

const WorkspaceListItem: FC<WorkspaceListItemProps> = ({
  workspace,
  open,
  selected,
  onSelect,
  backgroundColor,
}) => {
  const primaryRef = useRef<HTMLSpanElement | null>(null);
  const isTruncated = useIsTextTruncated(primaryRef, [open, workspace.name]);
  const showTooltip = !open || isTruncated;

  return (
    <CustomTooltip
      title={workspace.name}
      disableHoverListener={!showTooltip}
      placement="right"
      key={workspace.id}
    >
      <ListItem disablePadding sx={{ display: 'block', backgroundColor }}>
        <ListItemButton
          selected={selected}
          onClick={onSelect}
          sx={{
            minHeight: 48,
            px: 2.5,
            pl: open ? '2.5rem' : undefined,
            justifyContent: open ? 'initial' : 'center',
          }}
        >
          <ListItemIcon
            sx={{
              minWidth: 0,
              justifyContent: 'center',
              mr: open ? 3 : 'auto',
            }}
          >
            {workspace.icon}
          </ListItemIcon>
          <ListItemText
            primary={workspace.name}
            sx={{ opacity: open ? 1 : 0 }}
            slotProps={{ primary: { ref: primaryRef } }}
          />
        </ListItemButton>
      </ListItem>
    </CustomTooltip>
  );
};

export const WorkspacesSection: FC<WorkspacesSectionProps> = ({
  open,
  selectedId,
  onSelect,
  workspacesData,
  isLoading,
}) => {
  const theme = useTheme();

  const handleWorkspacesClick = () => {
    console.log('hey');
    onSelect('All Workspaces');
  };

  const workspaces: WorkspaceItem[] | undefined = workspacesData?.workspaces?.map((workspace) => ({
    id: workspace.id,
    name: workspace.name,
    icon: (
      <WorkspaceIcon
        fill={theme.palette.icon.default}
        secondaryFill={theme.palette.icon.default}
        {...iconSmall}
      />
    ),
  }));
  const { setMultiSelectedContent } = useContext(WorkspaceModalContext);
  return (
    <>
      <CustomTooltip title={'All Workspaces'} disableHoverListener={open} placement="right">
        <ListItem disablePadding sx={{ display: 'block' }}>
          <ListItemButton
            selected={selectedId === 'All Workspaces'}
            onClick={handleWorkspacesClick}
            sx={{
              minHeight: 48,
              px: 2.5,
              justifyContent: open ? 'initial' : 'center',
            }}
          >
            <ListItemIcon
              sx={{
                minWidth: 0,
                justifyContent: 'center',
                mr: open ? 3 : 'auto',
              }}
            >
              <WorkspaceIcon
                fill={theme.palette.icon.default}
                secondaryFill={theme.palette.icon.default}
                {...iconSmall}
              />
            </ListItemIcon>
            <ListItemText primary="All Workspaces" sx={{ opacity: open ? 1 : 0 }} />
          </ListItemButton>
        </ListItem>
      </CustomTooltip>
      {isLoading ? (
        <ListItem sx={{ pl: 4 }}>
          <ListItemText primary="Loading..." />
        </ListItem>
      ) : (
        workspaces &&
        workspaces.map((workspace) => (
          <WorkspaceListItem
            key={workspace.id}
            workspace={workspace}
            open={open}
            selected={selectedId === workspace.id}
            onSelect={() => {
              setMultiSelectedContent([]);
              onSelect(workspace.id);
            }}
            backgroundColor={theme.palette.background.secondary}
          />
        ))
      )}
    </>
  );
};
