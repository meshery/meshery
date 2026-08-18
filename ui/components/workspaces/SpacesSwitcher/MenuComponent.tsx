import {
  BottomSheet,
  CustomTooltip,
  DARK_BLUE_GRAY,
  IconButton,
  ListItemIcon,
  Menu,
  MenuItem,
  MenuList,
  MoreVertIcon as MoreVert,
  Typography,
  styled,
  useMediaQuery,
  useTheme,
} from '@sistent/sistent';
import { iconMedium } from 'css/icons.styles';
import { useState } from 'react';

const StyledMenuItem = styled(MenuItem)({
  paddingLeft: '.5rem',
  paddingRight: '.5rem',
  paddingTop: '.65rem',
  paddingBottom: '.65rem',
  margin: '0px',
  height: '100%',
});

const StyledMenuDiv = styled('div')(({ theme }) => ({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  flexDirection: 'row',
  height: '50%',
  position: 'relative',
  padding: '0',
  gap: '0.4rem',
  backgroundColor: theme.palette.mode == 'light' ? theme.palette.background.paper : DARK_BLUE_GRAY,
}));

export const MenuComponent = ({ options = [] }) => {
  const [anchorEl, setAnchorEl] = useState(null);
  const [sheetOpen, setSheetOpen] = useState(false);
  const theme = useTheme();
  const isCollapsed = useMediaQuery(theme.breakpoints.down('xl'));
  const isSmallScreen = useMediaQuery(theme.breakpoints.down('sm'));

  const handleClick = (event) => {
    event.stopPropagation();
    if (isSmallScreen) {
      setSheetOpen(true);
      return;
    }
    setAnchorEl(event.currentTarget);
  };

  const handleClose = (event?) => {
    event?.stopPropagation?.();
    event?.preventDefault?.();
    setAnchorEl(null);
    setSheetOpen(false);
  };

  const renderDirectIcons = () => {
    return (
      <div style={{ display: 'flex', gap: '0' }}>
        {options.map((option) => (
          <CustomTooltip key={option.title} title={option.title}>
            <IconButton
              sx={{
                padding: '0.15rem',
              }}
              disabled={option.disabled}
              style={{ cursor: 'pointer' }}
              onClick={(event) => {
                event.stopPropagation();
                option.handler(event);
              }}
            >
              {option.icon}
            </IconButton>
          </CustomTooltip>
        ))}
      </div>
    );
  };

  if (!isCollapsed) {
    return renderDirectIcons();
  }

  return (
    <div>
      <div style={{ display: 'flex', alignItems: 'center' }} onClick={handleClick}>
        <CustomTooltip title="Quick Actions">
          <MoreVert
            style={{
              cursor: 'pointer',
              fontSize: '1rem',
              ...iconMedium,
            }}
          />
        </CustomTooltip>
      </div>
      {isSmallScreen ? (
        <BottomSheet open={sheetOpen} onClose={handleClose} title="Actions">
          <MenuList disablePadding>
            {options.map((option) => (
              <MenuItem
                disabled={option.disabled}
                key={option.title}
                onClick={(event) => {
                  event.stopPropagation();
                  option.handler(event);
                  handleClose();
                }}
              >
                <ListItemIcon>{option.icon}</ListItemIcon>
                <Typography variant="body1">{option.title}</Typography>
              </MenuItem>
            ))}
          </MenuList>
        </BottomSheet>
      ) : (
        <Menu
          id="long-menu"
          anchorEl={anchorEl}
          anchorOrigin={{
            vertical: 'top',
            horizontal: 'left',
          }}
          MenuListProps={{
            style: {
              padding: 0,
              display: 'flex',
            },
          }}
          transformOrigin={{
            vertical: 'bottom',
            horizontal: 'left',
          }}
          open={Boolean(anchorEl)}
          onClose={handleClose}
          style={{
            borderRadius: '3px',
            zIndex: 9999999999,
          }}
        >
          {options.map((option) => (
            <StyledMenuDiv key={option.key || option.title}>
              <StyledMenuDiv>
                <CustomTooltip key={option.title} title={option.title}>
                  <StyledMenuItem
                    disabled={option.disabled}
                    key={option.title}
                    onClick={(event) => {
                      event.stopPropagation();
                      option.handler(event);
                      handleClose(event);
                    }}
                  >
                    {option.icon}
                  </StyledMenuItem>
                </CustomTooltip>
              </StyledMenuDiv>
            </StyledMenuDiv>
          ))}
        </Menu>
      )}
    </div>
  );
};
