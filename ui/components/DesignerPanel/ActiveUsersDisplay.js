import React, { useState } from 'react';
import { AvatarGroup, Avatar, CustomTooltip, Badge, Box, Button } from '@sistent/sistent';
import { styled } from '@sistent/sistent';
import DesignerPanel from './DesignerPanel';
import { useGetDesignQuery } from '@/rtk-query/design';
import { useRouter } from 'next/router';
import { openDesignInKanvas, useIsKanvasDesignerEnabled } from '@/utils/utils';

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

const UserAvatar = ({ user, onUserClick }) => {
  const router = useRouter();
  const isKanvasDesignerAvailable = useIsKanvasDesignerEnabled();

  const { data: designData } = useGetDesignQuery(
    { design_id: user.designId },
    { skip: !user.designId },
  );

  const designName = designData?.name || `Design ${user.designId}`;

  const handleOpenDesign = (e) => {
    e.stopPropagation(); // Prevent triggering the avatar click
    if (user.designId && designData) {
      if (!isKanvasDesignerAvailable) {
        router.push(`/configuration/designs/configurator?design_id=${user.designId}`);
      } else {
        openDesignInKanvas(user.designId, designData.name, router);
      }
    }
  };

  return (
    <CustomTooltip
      key={user.client_id || user.user_id}
      title={
        <div style={{ padding: '8px' }}>
          <div style={{ marginBottom: '4px' }}>{user.name || 'Unknown User'}</div>
          {user.designId && (
            <>
              <div style={{ fontSize: '0.8em', opacity: 0.8, marginBottom: '8px' }}>
                Working on: {designName}
              </div>
              <Button
                size="small"
                variant="outlined"
                onClick={handleOpenDesign}
                style={{
                  fontSize: '0.7em',
                  padding: '2px 8px',
                  minWidth: 'auto',
                }}
              >
                Open Design
              </Button>
            </>
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
          onClick={() => onUserClick(user, user.designId)}
        >
          {!user.avatar_url && (user.name?.[0] || '?')}
        </Avatar>
      </StyledBadge>
    </CustomTooltip>
  );
};

const ActiveUsersDisplay = ({
  activeUsers = {},
  getUserAccessToken,
  getUserProfile,
  maxDisplayed = 5,
}) => {
  const [selectedUser, setSelectedUser] = useState(null);
  const [designerPanelOpen, setDesignerPanelOpen] = useState(false);

  const activeUsersList = Object.entries(activeUsers)
    .filter(([designId]) => designId !== 'meshery_ui') //->removethe meshery_ui design key, as these users are just active
    .flatMap(([designId, users]) => users.map((user) => ({ ...user, designId })))
    .reduce((acc, user) => {
      if (!acc.find((u) => u.user_id === user.user_id || u.client_id === user.client_id)) {
        acc.push(user);
      }
      return acc;
    }, []);

  console.log('activeUsersList', activeUsersList);
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

  if (activeUsersList.length === 0) {
    return null;
  }

  return (
    <Box sx={{ display: 'flex', alignItems: 'center' }}>
      <StyledAvatarGroup max={maxDisplayed}>
        {activeUsersList.map((user) => (
          <UserAvatar
            key={user.client_id || user.user_id}
            user={user}
            onUserClick={handleUserClick}
          />
        ))}
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
