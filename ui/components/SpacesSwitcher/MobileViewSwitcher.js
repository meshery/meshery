import DashboardSwitcherIcon from '@/assets/icons/DashboardSwitcherIcon';
import {
  Box,
  Button,
  ClickAwayListener,
  Grid2,
  Slide,
  useTheme,
  WorkspaceIcon,
} from '@sistent/sistent';
import React, { useEffect, useState } from 'react';
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

function SwitcherMenu({ organization, router }) {
  const theme = useTheme();
  const [anchorEl, setAnchorEl] = useState(false);
  const [showFullContextMenu, setShowFullContextMenu] = useState(false);
  const [transformProperty, setTransformProperty] = useState(100);

  const styleSlider = {
    position: 'absolute',
    zIndex: '-1',
    bottom: '-170%',
    transform: showFullContextMenu ? `translateY(${transformProperty}%)` : 'translateY(0)',
  };

  let open = Boolean(anchorEl);
  if (showFullContextMenu) {
    open = showFullContextMenu;
  }

  useEffect(() => {
    setTransformProperty((prev) => prev + 2 * 3.125);
  }, []);

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
          aria-owns={open ? 'menu-list-grow' : undefined}
          aria-haspopup="true"
          style={{
            minWidth: '2.5rem',
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
                console.log('click away', e.target == 'body', e.target);
                if (
                  !(e.target.tagName.toLowerCase() === 'body') &&
                  !e.target.classList.contains('switcher-icon-button') &&
                  !e.target.classList.contains('MuiIconButton-root') &&
                  !e.target.classList.contains('MuiButtonBase-root') &&
                  !e.target.classList.contains('MuiSvgIcon-root') &&
                  !e.target.classList.contains('MuiTouchRipple-root') &&
                  !e.target.classList.contains('MuiTypography-root') &&
                  !e.target.classList.contains('MuiBox-root') &&
                  !e.target.classList.contains('MuiGrid2-root') &&
                  !e.target.classList.contains('MuiPaper-root') &&
                  !e.target.classList.contains('MuiClickAwayListener-root') &&
                  !e.target.classList.contains('MuiMenuItem-root') &&
                  !e.target.classList.contains('MuiListItem-root') &&
                  !e.target.classList.contains('MuiListItemButton-root') &&
                  !e.target.classList.contains('MuiListItemIcon-root') &&
                  !e.target.classList.contains('MuiListItemText-root') &&
                  !e.target.classList.contains('MuiList-root') &&
                  !e.target.classList.contains('MuiMenu-root') &&
                  !e.target.classList.contains('MuiPopover-root') &&
                  !e.target.classList.contains('MuiPaper-root') &&
                  !e.target.classList.contains('MuiSelect-select')
                ) {
                  setAnchorEl(false);
                  setShowFullContextMenu(false);
                }
              }}
            >
              <CMenuContainer>
                <Grid2
                  container
                  spacing={2}
                  alignItems="center"
                  sx={{
                    marginTop: '1rem',
                    backgroundColor: theme.palette.background.card,
                    borderRadius: '3px',
                    padding: '1rem',
                    boxShadow: '20px #979797',
                    transition: 'linear .2s',
                    transitionProperty: 'height',
                  }}
                  flexDirection={'column'}
                >
                  <Grid2
                    container
                    sx={{
                      display: 'flex',
                      flexDirection: 'row',
                      alignItems: 'center',
                      gap: '2rem',
                    }}
                  >
                    <OrgOutlinedIcon {...iconXLarge} fill={theme.palette.common.white} />
                    <OrgMenu open={true} organization={organization} />
                  </Grid2>
                  <Grid2
                    container
                    sx={{
                      display: 'flex',
                      flexDirection: 'row',
                      alignItems: 'center',
                      gap: '2.4rem',
                    }}
                  >
                    <WorkspaceIcon
                      {...iconLarge}
                      secondaryFill={theme.palette.common.white}
                      fill={theme.palette.common.white}
                    />
                    <WorkspaceSwitcher open={true} organization={organization} router={router} />
                  </Grid2>
                </Grid2>
              </CMenuContainer>
            </ClickAwayListener>
          </div>
        </Slide>
      </div>
    </>
  );
}
