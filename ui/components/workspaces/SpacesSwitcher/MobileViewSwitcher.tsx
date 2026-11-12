import OrgOutlinedIcon from '@/assets/icons/OrgOutlinedIcon';
import { WorkspaceModalContext } from '@/utils/context/WorkspaceModalContextProvider';
import { Keys } from '@meshery/schemas/permissions';
import { BottomSheet, Box, Button, Typography, useTheme, WorkspaceIcon } from '@sistent/sistent';
import { iconMedium, iconSmall } from 'css/icons.styles';
import React, { useContext, useState } from 'react';
import WorkspaceSwitcher from './WorkspaceSwitcher';
import { OrgMenu } from './SpaceSwitcher';

const MobileOrgWksSwither_ = ({ organization, router }) => (
  <SwitcherMenu organization={organization} router={router} />
);

export const MobileOrgWksSwither = MobileOrgWksSwither_;

function SectionLabel({ children, icon }: { children: React.ReactNode; icon: React.ReactNode }) {
  const theme = useTheme();

  return (
    <Box
      sx={{
        display: 'flex',
        alignItems: 'center',
        gap: 1,
        mb: 1,
      }}
    >
      <Box sx={{ display: 'flex', flexShrink: 0, alignItems: 'center' }}>{icon}</Box>
      <Typography
        variant="subtitle2"
        component="p"
        sx={{
          m: 0,
          color: theme.palette.text.secondary,
          fontWeight: 500,
          fontSize: '0.8125rem',
        }}
      >
        {children}
      </Typography>
    </Box>
  );
}

type ExpandedField = 'org' | 'workspace' | null;

function SwitcherMenu({ organization: _organization, router: _router }) {
  const theme = useTheme();
  const [open, setOpen] = useState(false);
  const [expandedField, setExpandedField] = useState<ExpandedField>(null);
  const {
    openModal: openWorkspaceModal,
    setCreateNewWorkspaceModalOpen,
    setSelectedWorkspace,
  } = useContext(WorkspaceModalContext);

  const handleClose = () => {
    setOpen(false);
    setExpandedField(null);
  };

  const openWorkspaceExplorer = () => {
    handleClose();
    openWorkspaceModal(true);
  };

  const openCreateWorkspace = () => {
    handleClose();
    setSelectedWorkspace({
      id: 'All Workspaces',
      name: 'All Workspaces',
    });
    setCreateNewWorkspaceModalOpen(true);
    openWorkspaceModal(true);
  };

  return (
    <>
      <Button
        aria-label="contexts"
        className="switcher-icon-button"
        onClick={(e) => {
          e.preventDefault();
          setOpen((prev) => !prev);
        }}
        aria-expanded={open}
        aria-haspopup="dialog"
        sx={{
          minWidth: 'auto',
          flexShrink: 0,
          ml: 0.25,
          px: 0.5,
          py: 0.25,
        }}
      >
        <Box
          sx={{
            display: 'flex',
            alignItems: 'center',
            gap: 0.5,
          }}
        >
          <OrgOutlinedIcon {...iconSmall} fill={theme.palette.common.white} />
          <WorkspaceIcon
            {...iconSmall}
            fill={theme.palette.common.white}
            secondaryFill={theme.palette.common.white}
          />
        </Box>
      </Button>

      <BottomSheet
        open={open}
        onClose={handleClose}
        title="Organization & Workspace"
        maxHeight="85vh"
      >
        <Box
          sx={{
            display: 'flex',
            flexDirection: 'column',
            minHeight: '100%',
            gap: 2.5,
          }}
        >
          <Box>
            <SectionLabel
              icon={<OrgOutlinedIcon {...iconMedium} fill={theme.palette.icon.default} />}
            >
              Organization
            </SectionLabel>
            <OrgMenu
              open={true}
              fromMobileView={true}
              expanded={expandedField === 'org'}
              onExpandedChange={(isExpanded) => setExpandedField(isExpanded ? 'org' : null)}
            />
          </Box>

          <Box>
            <SectionLabel
              icon={
                <WorkspaceIcon
                  {...iconMedium}
                  fill={theme.palette.icon.default}
                  secondaryFill={theme.palette.icon.default}
                />
              }
            >
              Workspace
            </SectionLabel>
            <WorkspaceSwitcher
              open={true}
              fromMobileView={true}
              expanded={expandedField === 'workspace'}
              onExpandedChange={(isExpanded) => setExpandedField(isExpanded ? 'workspace' : null)}
              onMobileSelect={handleClose}
            />
          </Box>

          <Box
            sx={{
              mt: 'auto',
              pt: 2,
              borderTop: `1px solid ${theme.palette.divider}`,
              display: 'flex',
              flexDirection: 'column',
              gap: 1.5,
            }}
          >
            <Button
              variant="contained"
              fullWidth
              onClick={openWorkspaceExplorer}
              permissionKey={Keys.WorkspaceManagementViewWorkspace}
            >
              Explore Workspaces
            </Button>
            <Button
              variant="outlined"
              fullWidth
              onClick={openCreateWorkspace}
              permissionKey={Keys.WorkspaceManagementCreateWorkspace}
            >
              Create Workspace
            </Button>
          </Box>
        </Box>
      </BottomSheet>
    </>
  );
}
