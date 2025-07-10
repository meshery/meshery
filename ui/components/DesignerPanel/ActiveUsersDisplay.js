import React, { useState } from 'react';
import { AvatarGroup, Avatar, CustomTooltip, Badge, Box } from '@sistent/sistent';
import { styled } from '@sistent/sistent';
import DesignerPanel from './DesignerPanel';

const StyledAvatarGroup = styled(AvatarGroup)(({ theme }) => ({
  '& .MuiAvatar-root': {
    width: 32,
    height: 32,
    border: `2px solid ${theme.palette.background.paper}`,
    cursor: 'pointer',
    transition: 'transform 0.2s ease-in-out',
    '&:hover': {
      transform: 'scale(1.1)',
      zIndex: 10,
    },
  },
}));

const StyledBadge = styled(Badge)(({ color }) => ({
  '& .MuiBadge-badge': {
    backgroundColor: color,
    color: color,
    border: '2px solid white',
    '&::after': {
      position: 'absolute',
      top: 0,
      left: 0,
      width: '100%',
      height: '100%',
      borderRadius: '50%',
      animation: 'ripple 1.2s infinite ease-in-out',
      border: '1px solid currentColor',
      content: '""',
    },
  },
  '@keyframes ripple': {
    '0%': {
      transform: 'scale(.8)',
      opacity: 1,
    },
    '100%': {
      transform: 'scale(1.8)',
      opacity: 0,
    },
  },
}));

const ActiveUsersDisplay = ({
  activeUsers = {},
  getUserAccessToken,
  getUserProfile,
  maxDisplayed = 5,
}) => {
  // console.log('activeUsers', activeUsers);
  const [selectedUser, setSelectedUser] = useState(null);
  const [designerPanelOpen, setDesignerPanelOpen] = useState(false);
  console.log('I am userDisplay');
  const activeUsersList = Object.entries(activeUsers)
    .filter(([designId]) => designId !== 'meshery_ui') //->removethe meshery_ui design key, as these users are just active
    .flatMap(([, users]) => users)
    .reduce((acc, user) => {
      if (!acc.find((u) => u.user_id === user.user_id || u.client_id === user.client_id)) {
        acc.push(user);
      }
      return acc;
    }, []);
  // console.log('activeUsersList', activeUsersList);
  const handleUserClick = (user, designId = null) => {
    setSelectedUser({
      ...user,
      currentDesignId: designId,
      currentDesignName: designId ? `Design ${designId}` : null,
    });
    setDesignerPanelOpen(true);
  };

  const handleClosePanel = () => {
    setDesignerPanelOpen(false);
    setSelectedUser(null);
  };

  const getUserDesignInfo = (user) => {
    for (const [designId, users] of Object.entries(activeUsers)) {
      const userInDesign = users.find(
        (u) => u.user_id === user.user_id || u.client_id === user.client_id,
      );
      if (userInDesign) {
        return {
          designId,
          designName: `Design ${designId}`,
        };
      }
    }
    return null;
  };

  if (activeUsersList.length === 0) {
    return null;
  }
  console.log('activeUsers', activeUsers);
  console.log('getUserAccessToken', getUserAccessToken);
  console.log('getUserProfile', getUserProfile);

  return (
    <Box sx={{ display: 'flex', alignItems: 'center' }}>
      <StyledAvatarGroup max={maxDisplayed}>
        {activeUsersList.map((user) => {
          const designInfo = getUserDesignInfo(user);
          return (
            <CustomTooltip
              key={user.client_id || user.user_id}
              title={
                <div>
                  <div>{user.name || 'Unknown User'}</div>
                  {designInfo && (
                    <div style={{ fontSize: '0.8em', opacity: 0.8 }}>
                      Working on: {designInfo.designName}
                    </div>
                  )}
                </div>
              }
              placement="bottom"
            >
              <StyledBadge
                overlap="circular"
                anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
                variant="dot"
                color={user.color}
              >
                <Avatar
                  src={user.avatar_url}
                  alt={user.name}
                  onClick={() => handleUserClick(user, designInfo?.designId)}
                  sx={{
                    border: user.color ? `2px solid ${user.color}` : undefined,
                  }}
                >
                  {!user.avatar_url && (user.name?.[0] || '?')}
                </Avatar>
              </StyledBadge>
            </CustomTooltip>
          );
        })}
      </StyledAvatarGroup>

      <DesignerPanel
        open={designerPanelOpen}
        onClose={handleClosePanel}
        activeUser={selectedUser}
        currentDesignId={selectedUser?.currentDesignId}
        getUserAccessToken={getUserAccessToken}
        getUserProfile={getUserProfile}
      />
    </Box>
  );
};

export default ActiveUsersDisplay;
