import React, { useState, useContext } from 'react';
import {
  Button,
  styled,
  Typography,
  useTheme,
  WorkspaceIcon,
  CustomTooltip,
  useMediaQuery,
} from '@sistent/sistent';
import { NoSsr } from '@sistent/sistent';
import { useRouter } from 'next/router';
import OrgOutlinedIcon from '@/assets/icons/OrgOutlinedIcon';
import { iconLarge, iconXLarge } from '@/css/icons.styles';
import { useDynamicComponent } from '@/utils/context/dynamicContext';
import WorkspaceSwitcher from './WorkspaceSwitcher';

import { OrgMenu } from './OrgMenu';
import { useSelector } from 'react-redux';
import { useGetSelectedOrganization } from '@/rtk-query/user';
import { MobileOrgWksSwither } from './MobileViewSwitcher';
import WorkspaceModal from '@/components/workspaces/WorkspaceFormModal';
import { WorkspaceModalContext } from '@/utils/context/WorkspaceModalContextProvider';

export const StyledHeader = styled(Typography)(({ theme }) => ({
  paddingLeft: theme.spacing(1),
  fontSize: '1.65rem',

  [theme.breakpoints.down('sm')]: {
    fontSize: '1.25rem',
    maxWidth: '7rem',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    whiteSpace: 'nowrap',
  },
  color: theme.palette.common.white,
}));
export const StyledBetaHeader = styled('sup')(({ theme }) => ({
  color: theme.palette.background.constant.white,
  fontWeight: '300',
  fontSize: '0.8125rem',
}));

const StyledSwitcher = styled('div')(({ theme }) => ({
  display: 'flex',
  flexDirection: 'row',
  justifyContent: 'center',
  alignItems: 'center',
  fontSize: '1.5rem',
  userSelect: 'none',
  transition: 'width 2s ease-in',
  color: theme.palette.common.white,
  gap: '0.5rem 0rem',
}));

function DefaultHeader({ title, isBeta }) {
  return (
    <StyledHeader variant="h5" data-cy="headerPageTitle">
      {title}
      {isBeta ? <StyledBetaHeader>BETA</StyledBetaHeader> : ''}
    </StyledHeader>
  );
}

function OrganizationAndWorkSpaceSwitcher() {
  const [orgOpen, setOrgOpen] = useState(false);
  const [workspaceOpen, setWorkspaceOpen] = useState(false);

  const { DynamicComponent } = useDynamicComponent();
  const theme = useTheme();
  const isSmallScreen = useMediaQuery(theme.breakpoints.down('md'));
  const router = useRouter();
  const { organization } = useSelector((state) => state.ui);
  const { isBeta } = useSelector((state) => state.ui.page);
  const { title } = useSelector((state) => state.ui.page);
  const { selectedOrganization } = useGetSelectedOrganization();

  //->using the wksp cntxt
  const { open: workspaceModal, closeModal: closeWorkspaceModal } =
    useContext(WorkspaceModalContext);

  if (!selectedOrganization) return null;

  return (
    <NoSsr>
      <StyledSwitcher>
        {isSmallScreen && (
          <>
            <MobileOrgWksSwither organization={organization} router={router} />/
          </>
        )}
        {!isSmallScreen && (
          <>
            <CustomTooltip title={'Organization'}>
              <div>
                <Button
                  onClick={() => {
                    setOrgOpen(!orgOpen);
                  }}
                  style={{ marginRight: orgOpen ? '1rem' : '0' }}
                >
                  <OrgOutlinedIcon {...iconXLarge} fill={theme.palette.common.white} />
                </Button>
              </div>
            </CustomTooltip>
            <OrgMenu open={orgOpen} organization={organization} />/
            <CustomTooltip title={'Workspaces'}>
              <div>
                <Button
                  onClick={() => {
                    setWorkspaceOpen(!workspaceOpen);
                  }}
                  style={{ marginRight: workspaceOpen ? '1rem' : '0' }}
                >
                  <WorkspaceIcon
                    {...iconLarge}
                    secondaryFill={theme.palette.common.white}
                    fill={theme.palette.common.white}
                  />
                </Button>
              </div>
            </CustomTooltip>
            <WorkspaceSwitcher open={workspaceOpen} organization={organization} router={router} />/
          </>
        )}
        <div id="meshery-dynamic-header" style={{ marginLeft: DynamicComponent ? '0' : '' }} />
        {!DynamicComponent && <DefaultHeader title={title} isBeta={isBeta} />}
        <WorkspaceModal workspaceModal={workspaceModal} closeWorkspaceModal={closeWorkspaceModal} />
      </StyledSwitcher>
    </NoSsr>
  );
}

export default OrganizationAndWorkSpaceSwitcher;
