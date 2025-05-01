//@ts-check
import { useGetLoggedInUserQuery } from '@/rtk-query/user';
import { CloneIcon } from '@layer5/sistent';
import { CustomTooltip, styled, Menu, MenuItem } from '@layer5/sistent';
import { Public as PublicIcon, Reply } from '@mui/icons-material';
import DeleteIcon from '@mui/icons-material/Delete';
import GetAppIcon from '@mui/icons-material/GetApp';
import InfoIcon from '@mui/icons-material/Info';
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
          horizontal: 'right',
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
          <StyledMenuDiv key={option.key}>
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
          horizontal: 'right',
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
          <StyledMenuDiv key={CustomTooltip} data-testid={`design-tr-menu-list-${rowData?.id}`}>
            <StyledMenuDiv>
              {visibility === 'published' &&
                currentUser?.user_details?.role_names?.includes('admin') && (
                  <CustomTooltip key="unpublish" title="Unpublish">
                    <StyledMenuItem
                      data-testid={`designs-tr-menu-li-unpublish-${rowData?.id}`}
                      key="unpublish"
                      onClick={(event) => {
                        option.unPublishHandler();
                        handleClose(event);
                      }}
                    >
                      <PublicIcon styles={{ ...iconMedium }} />
                    </StyledMenuItem>
                  </CustomTooltip>
                )}

              {visibility === 'published' && (
                <CustomTooltip key="clone" title="Clone">
                  <StyledMenuItem
                    data-testid={`designs-tr-menu-li-clone-${rowData?.id}`}
                    key="clone"
                    onClick={(event) => {
                      option.cloneHandler();
                      handleClose(event);
                    }}
                  >
                    <CloneIcon fill="#eee" style={{ ...iconSmall }} />
                  </StyledMenuItem>
                </CustomTooltip>
              )}
              <CustomTooltip key="export_design" title="Export Design">
                <StyledMenuItem
                  data-testid={`designs-tr-menu-li-export-${rowData?.id}`}
                  key="export_design"
                  onClick={(event) => {
                    handleClose(event);
                  }}
                >
                  <GetAppIcon fill="#eee" style={{ ...iconMedium }} />
                </StyledMenuItem>
              </CustomTooltip>

              <CustomTooltip key="delete" title="Delete">
                <StyledMenuItem
                  data-testid={`designs-tr-menu-li-delete-${rowData?.id}`}
                  key="delete"
                  onClick={(event) => {
                    handleClose(event);
                  }}
                >
                  <DeleteIcon fill="#eee" style={{ ...iconMedium }} />
                </StyledMenuItem>
              </CustomTooltip>

              {visibility !== 'published' && (
                <CustomTooltip key="share" title="Share">
                  <StyledMenuItem
                    data-testid={`designs-tr-menu-li-share-${rowData?.id}`}
                    key="Share"
                    onClick={(event) => {
                      handleClose(event);
                    }}
                  >
                    <Reply style={{ ...iconMedium, transform: 'scaleX(-1)', color: '#eee' }} />
                  </StyledMenuItem>
                </CustomTooltip>
              )}

              <CustomTooltip key="info" title="Info">
                <StyledMenuItem
                  data-testid={`designs-tr-menu-li-info-${rowData?.id}`}
                  key="Info"
                  onClick={(event) => {
                    handleClose(event);
                  }}
                >
                  <InfoIcon style={{ ...iconMedium, transform: 'scaleX(-1)', color: '#eee' }} />
                </StyledMenuItem>
              </CustomTooltip>
            </StyledMenuDiv>
          </StyledMenuDiv>
        ))}
      </Menu>
    </div>
  );
};

export default MenuComponent;
