import DashboardSwitcherIcon from '@/assets/icons/DashboardSwitcherIcon';
import {
  Box,
  Button,
  ClickAwayListener,
  Grid2,
  Slide,
  styled,
  useMediaQuery,
  useTheme,
  WorkspaceIcon,
} from '@sistent/sistent';
import React, { useState } from 'react';
import { CMenuContainer } from '../Header.styles';
import WorkspaceSwitcher from './WorkspaceSwitcher';
import OrgOutlinedIcon from '@/assets/icons/OrgOutlinedIcon';
import { iconLarge, iconXLarge } from 'css/icons.styles';
import { OrgMenu } from './SpaceSwitcher';

const MobileOrgWksSwither_ = ({ organization, router }) => {
  return (
    <Box sx={{ margin: '0.25rem 0' }}>
      <div
        style={{
          display: 'flex',
          justifyContent: 'flex-start',
          alignItems: 'center',
        }}
      >
        <SwitcherMenu organization={organization} router={router} />
      </div>
    </Box>
  );
};

export const MobileOrgWksSwither = MobileOrgWksSwither_;

const StyledGrid = styled(Grid2)({
  width: '100%',
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'center',
  justifyContent: 'space-between',
});

function SwitcherMenu({ organization, router }) {
  const theme = useTheme();
  const isSmallScreen = useMediaQuery('(max-width:400px)');
  const [anchorEl, setAnchorEl] = useState(false);
  const [showFullContextMenu, setShowFullContextMenu] = useState(false);

  const styleSlider = {
    position: 'absolute',
    zIndex: '-1',
    top: '-8%',
    width: isSmallScreen ? '100%' : '270px',
    transition: 'top 0.4s ease, transform 0.4s ease',
    transform: showFullContextMenu ? `translateY(105%)` : 'translateY(-8%)',
  };

  let open = Boolean(anchorEl);
  if (showFullContextMenu) {
    open = showFullContextMenu;
  }

  return (
    <>
      <div>
        <Button
          aria-label="contexts"
          className="switcher-icon-button"
          onClick={(e) => {
            e.preventDefault();
            setShowFullContextMenu((prev) => !prev);
          }}
          onMouseOver={(e) => {
            e.preventDefault();
            setAnchorEl(true);
          }}
          onMouseLeave={(e) => {
            e.preventDefault();
            setAnchorEl(false);
          }}
          aria-owns={open ? 'menu-list-grow' : undefined}
          aria-haspopup="true"
          style={{
            minWidth: '2.5rem',
            marginLeft: '0.5rem',
          }}
        >
          <DashboardSwitcherIcon height={28} width={28} />
        </Button>

        <Slide
          direction="down"
          style={styleSlider}
          timeout={400}
          in={open}
          mountOnEnter
          unmountOnExit
        >
          <div>
            <ClickAwayListener
              onClickAway={(e) => {
                if (
                  !(e.target.tagName.toLowerCase() === 'body') &&
                  !e.target.classList.contains('switcher-icon-button') &&
                  !(e.target.tagName.toLowerCase() === 'path') &&
                  !(e.target.tagName.toLowerCase() === 'svg') &&
                  !(e.target.tagName.toLowerCase() === 'circle')
                ) {
                  setAnchorEl(false);
                  setShowFullContextMenu(false);
                }
              }}
            >
              <CMenuContainer>
                <Grid2 container spacing={2} alignItems="center" flexDirection={'column'}>
                  <StyledGrid container>
                    <OrgOutlinedIcon {...iconXLarge} fill={theme.palette.common.white} />
                    <OrgMenu open={true} organization={organization} />
                  </StyledGrid>

                  <StyledGrid container>
                    <WorkspaceIcon
                      {...iconLarge}
                      style={{ marginLeft: '0.2rem' }}
                      secondaryFill={theme.palette.common.white}
                      fill={theme.palette.common.white}
                    />
                    <WorkspaceSwitcher open={true} organization={organization} router={router} />
                  </StyledGrid>
                </Grid2>
              </CMenuContainer>
            </ClickAwayListener>
          </div>
        </Slide>
      </div>
    </>
  );
}
