//@ts-check
import { useGetLoggedInUserQuery } from '@/rtk-query/user';
import { CloneIcon } from '@layer5/sistent';
import { CustomTooltip, styled, Menu, MenuItem } from '@layer5/sistent';
import { Public as PublicIcon, Reply } from '@mui/icons-material';
import DeleteIcon from '@mui/icons-material/Delete';
import GetAppIcon from '@mui/icons-material/GetApp';
import InfoIcon from '@mui/icons-material/Info';
import { useMediaQuery, useTheme } from '@mui/material';
import { iconMedium, iconSmall } from 'css/icons.styles';
import React from 'react';

const StyledMenuItem = styled(MenuItem)({
  paddingLeft: '.5rem',
  paddingRight: '.5rem',
  paddingTop: '.65rem',
  paddingBottom: '.65rem',
  color: '#eee',
  margin: '0px',
  height: '100%',
  '&:hover': {
    backgroundColor: '#36474f',
  },
});

const StyledMenuDiv = styled('div')({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  flexDirection: 'row',
  height: '50%',
  position: 'relative',
  padding: '0',
  gap: '0.4rem',
  backgroundColor: '#263238',
});

export const GeorgeMenu = ({ triggerIcon, options = [] }) => {
  // States.
  const [anchorEl, setAnchorEl] = React.useState(null);
  const open = Boolean(anchorEl);
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('lg'));

  // Handlers.
  const handleClick = (event) => {
    event.stopPropagation();
    setAnchorEl(event.currentTarget);
  };

  const handleClose = (event) => {
    event.stopPropagation();
    event.preventDefault();
    setAnchorEl(null);
  };

  const WrapperIcon = triggerIcon;

  // Renders direct icons for non-mobile view
  const renderDirectIcons = () => {
    return (
      <div style={{ display: 'flex', gap: '0.5rem' }}>
        {options.map((option) => (
          <CustomTooltip key={option.title} title={option.title}>
            <div
              style={{ cursor: 'pointer' }}
              onClick={(event) => {
                option.handler(event);
              }}
            >
              {option.icon && (
                <option.icon fill="#eeeeee" style={{ ...iconMedium }} {...iconMedium} />
              )}
            </div>
          </CustomTooltip>
        ))}
      </div>
    );
  };

  // Return direct icons for non-mobile/larger screens
  if (!isMobile) {
    return renderDirectIcons();
  }

  // Default menu for mobile screens
  return (
    <div>
      <div style={{ display: 'flex', alignItems: 'center' }} onClick={handleClick}>
        <CustomTooltip title="Single-click Actions">
          <WrapperIcon
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
                    option.handler(event);
                    handleClose(event);
                  }}
                >
                  {option.icon && (
                    <option.icon fill="#eeeeee" style={{ ...iconMedium }} {...iconMedium} />
                  )}
                </StyledMenuItem>
              </CustomTooltip>
            </StyledMenuDiv>
          </StyledMenuDiv>
        ))}
      </Menu>
    </div>
  );
};

const MenuComponent = ({ iconType, items, visibility, rowData = null }) => {
  // States.
  const [anchorEl, setAnchorEl] = React.useState(null);
  const open = Boolean(anchorEl);
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('lg'));

  // Handlers.
  const handleClick = (event) => {
    event.stopPropagation();
    setAnchorEl(event.currentTarget);
  };

  const handleClose = (event) => {
    event.stopPropagation();
    event.preventDefault();
    setAnchorEl(null);
  };

  const WrapperIcon = iconType;
  const { data: currentUser } = useGetLoggedInUserQuery();

  // Function to handle specific actions
  const handleAction = (action, event) => {
    event.stopPropagation();
    if (action === 'unpublish' && items[0]?.unPublishHandler) {
      items[0].unPublishHandler();
    } else if (action === 'clone' && items[0]?.cloneHandler) {
      items[0].cloneHandler();
    }
    if (!isMobile) handleClose(event);
  };

  // Common configuration for menu items
  const menuItems = [
    {
      key: 'unpublish',
      title: 'Unpublish',
      icon: PublicIcon,
      action: 'unpublish',
      visible:
        visibility === 'published' && currentUser?.user_details?.role_names?.includes('admin'),
    },
    {
      key: 'clone',
      title: 'Clone',
      icon: CloneIcon,
      action: 'clone',
      iconProps: { fill: '#eee', style: { ...iconSmall } },
      visible: visibility === 'published',
    },
    {
      key: 'export_design',
      title: 'Export Design',
      icon: GetAppIcon,
      action: 'export',
      iconProps: { fill: '#eee', style: { ...iconMedium } },
      visible: true,
    },
    {
      key: 'delete',
      title: 'Delete',
      icon: DeleteIcon,
      action: 'delete',
      iconProps: { fill: '#eee', style: { ...iconMedium } },
      visible: true,
    },
    {
      key: 'share',
      title: 'Share',
      icon: Reply,
      action: 'share',
      iconProps: { style: { ...iconMedium, transform: 'scaleX(-1)', color: '#eee' } },
      visible: visibility !== 'published',
    },
    {
      key: 'info',
      title: 'Info',
      icon: InfoIcon,
      action: 'info',
      iconProps: { style: { ...iconMedium, transform: 'scaleX(-1)', color: '#eee' } },
      visible: true,
    },
  ];

  // Renders direct icons for mobile view
  const renderDirectIcons = () => {
    return (
      <div
        data-testid={`designs-tr-icons-${rowData?.id}`}
        style={{ display: 'flex', gap: '0.5rem' }}
      >
        {menuItems
          .filter((item) => item.visible)
          .map((item) => (
            <CustomTooltip key={item.key} title={item.title}>
              {item.key === 'clone' ? (
                <div
                  style={{ cursor: 'pointer' }}
                  onClick={(event) => handleAction(item.action, event)}
                >
                  <item.icon {...(item.iconProps || {})} />
                </div>
              ) : (
                <item.icon
                  {...(item.iconProps || {})}
                  style={{ ...(item.iconProps?.style || {}), cursor: 'pointer' }}
                  onClick={(event) => handleAction(item.action, event)}
                />
              )}
            </CustomTooltip>
          ))}
      </div>
    );
  };

  // Render menu items for the dropdown
  const renderMenuItems = (option) => {
    return menuItems
      .filter((item) => item.visible)
      .map((item) => (
        <CustomTooltip key={item.key} title={item.title}>
          <StyledMenuItem
            data-testid={`designs-tr-menu-li-${item.key}-${rowData?.id}`}
            key={item.key}
            onClick={(event) => {
              if (item.action === 'unpublish') option.unPublishHandler();
              else if (item.action === 'clone') option.cloneHandler();
              handleClose(event);
            }}
          >
            <item.icon {...(item.iconProps || {})} />
          </StyledMenuItem>
        </CustomTooltip>
      ));
  };

  // Return direct icons for non-mobile/larger screens
  if (!isMobile) {
    return renderDirectIcons();
  }

  // Default menu for mobile screens
  return (
    <div data-testid={`designs-tr-menu-${rowData?.id}`}>
      <div style={{ display: 'flex', alignItems: 'center' }} onClick={handleClick}>
        <CustomTooltip title="Quick actions">
          <WrapperIcon style={{ fontSize: '1rem', ...iconMedium }} />
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
