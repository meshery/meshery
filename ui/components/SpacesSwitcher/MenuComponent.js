import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import { DARK_BLUE_GRAY, IconButton } from '@layer5/sistent';
import { CustomTooltip, styled, Menu, MenuItem } from '@layer5/sistent';
import { MoreVert, Reply } from '@mui/icons-material';
import DeleteIcon from '@mui/icons-material/Delete';
import GetAppIcon from '@mui/icons-material/GetApp';
import InfoIcon from '@mui/icons-material/Info';
import { useMediaQuery, useTheme } from '@mui/material';
import { iconMedium } from 'css/icons.styles';
import React from 'react';

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

export const GeorgeMenu = ({ options = [] }) => {
  const [anchorEl, setAnchorEl] = React.useState(null);
  const open = Boolean(anchorEl);
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('lg'));

  const handleClick = (event) => {
    event.stopPropagation();
    setAnchorEl(event.currentTarget);
  };

  const handleClose = (event) => {
    event.stopPropagation();
    event.preventDefault();
    setAnchorEl(null);
  };

  const renderDirectIcons = () => {
    return (
      <div style={{ display: 'flex' }}>
        {options.map((option) => (
          <CustomTooltip key={option.title} title={option.title}>
            <IconButton
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

  if (!isMobile) {
    return renderDirectIcons();
  }

  return (
    <div>
      <div style={{ display: 'flex', alignItems: 'center' }} onClick={handleClick}>
        <CustomTooltip title="Single-click Actions">
          <MoreVert
            style={{
              cursor: 'pointer',
              fontSize: '1rem',
              ...iconMedium,
            }}
          />
        </CustomTooltip>
      </div>
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
        open={open}
        onClose={handleClose}
        style={{
          borderRadius: '3px',
        }}
      >
        {options.map((option) => (
          <StyledMenuDiv key={option.key || option.title}>
            <StyledMenuDiv>
              <CustomTooltip key={option.title} title={option.title}>
                <StyledMenuItem
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
    </div>
  );
};

const MenuComponent = ({ items, visibility, rowData = null }) => {
  const [anchorEl, setAnchorEl] = React.useState(null);
  const open = Boolean(anchorEl);
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('lg'));

  const handleClick = (event) => {
    event.stopPropagation();
    setAnchorEl(event.currentTarget);
  };

  const handleClose = (event) => {
    event.stopPropagation();
    event.preventDefault();
    setAnchorEl(null);
  };

  const handleAction = (action, event) => {
    event.stopPropagation();
    if (action === 'export' && items[0]?.downloadHandler) {
      items[0].downloadHandler();
    } else if (action === 'delete' && items[0]?.deleteHandler) {
      items[0].deleteHandler();
    } else if (action === 'share' && items[0]?.shareHandler) {
      items[0].shareHandler();
    } else if (action === 'info' && items[0]?.infoHandler) {
      items[0].infoHandler();
    }

    if (!isMobile) handleClose(event);
  };

  const menuItems = [
    {
      key: 'export_design',
      title: 'Export Design',
      icon: GetAppIcon,
      action: 'export',
      iconProps: { style: { ...iconMedium } },
      visible: CAN(keys.DOWNLOAD_A_DESIGN.action, keys.DOWNLOAD_A_DESIGN.subject),
    },
    {
      key: 'delete',
      title: 'Delete',
      icon: DeleteIcon,
      action: 'delete',
      iconProps: { style: { ...iconMedium } },
      visible: CAN(keys.DELETE_A_DESIGN.action, keys.DELETE_A_DESIGN.subject),
    },
    {
      key: 'share',
      title: 'Share',
      icon: Reply,
      action: 'share',
      iconProps: { style: { ...iconMedium, transform: 'scaleX(-1)' } },
      visible:
        visibility !== 'published' && CAN(keys.SHARE_DESIGN.action, keys.SHARE_DESIGN.subject),
    },
    {
      key: 'info',
      title: 'Info',
      icon: InfoIcon,
      action: 'info',
      iconProps: { style: { ...iconMedium, transform: 'scaleX(-1)' } },
      visible: true,
    },
  ];

  const renderDirectIcons = () => {
    return (
      <div data-testid={`designs-tr-icons-${rowData?.id}`} style={{ display: 'flex' }}>
        {menuItems
          .filter((item) => item.visible)
          .map((item) => (
            <CustomTooltip key={item.key} title={item.title}>
              <IconButton onClick={(event) => handleAction(item.action, event)}>
                <item.icon style={{ ...(item.iconProps?.style || {}) }} />
              </IconButton>
            </CustomTooltip>
          ))}
      </div>
    );
  };

  const renderMenuItems = (option) => {
    return menuItems
      .filter((item) => item.visible)
      .map((item) => (
        <CustomTooltip key={item.key} title={item.title}>
          <StyledMenuItem
            data-testid={`designs-tr-menu-li-${item.key}-${rowData?.id}`}
            key={item.key}
            onClick={(event) => {
              switch (item.action) {
                case 'export':
                  option.downloadHandler();
                  break;
                case 'delete':
                  option.deleteHandler();
                  break;
                case 'share':
                  option.shareHandler();
                  break;
                case 'info':
                  option.infoHandler();
                  break;
                default:
                  break;
              }
              handleClose(event);
            }}
          >
            <item.icon {...(item.iconProps || {})} />
          </StyledMenuItem>
        </CustomTooltip>
      ));
  };

  if (!isMobile) {
    return renderDirectIcons();
  }

  return (
    <div data-testid={`designs-tr-menu-${rowData?.id}`}>
      <CustomTooltip title="Quick actions">
        <div>
          <IconButton onClick={handleClick}>
            <MoreVert style={{ fontSize: '1rem', ...iconMedium }} />
          </IconButton>
        </div>
      </CustomTooltip>
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
        open={open}
        onClose={handleClose}
        style={{ borderRadius: '3px' }}
      >
        {items.map((option) => (
          <StyledMenuDiv
            key={option.key || 'menu-item'}
            data-testid={`design-tr-menu-list-${rowData?.id}`}
          >
            <StyledMenuDiv>{renderMenuItems(option)}</StyledMenuDiv>
          </StyledMenuDiv>
        ))}
      </Menu>
    </div>
  );
};

export default MenuComponent;
