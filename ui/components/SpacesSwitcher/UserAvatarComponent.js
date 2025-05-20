import { MESHERY_CLOUD_PROD } from '@/constants/endpoints';
import React from 'react';
import { CustomTooltip, Link, Avatar, Typography } from '@layer5/sistent';
import { StyledAvatarContainer, StyledUserDetailsContainer, StyledUpdatedText } from './styles';

const UserAvatarComponent = ({ userData }) => {
  return (
    <StyledAvatarContainer>
      <CustomTooltip
        title={`${userData?.first_name} ${userData?.last_name || ''}`}
        sx={{ display: { xs: 'block', sm: 'block' } }}
      >
        <div>
          <Link href={`${MESHERY_CLOUD_PROD}/user/${userData?.id}`}>
            <Avatar alt={userData?.first_name} src={userData?.avatar_url} />
          </Link>
        </div>
      </CustomTooltip>

      <StyledUserDetailsContainer
        sx={{
          display: {
            xs: 'none',
            sm: 'flex',
            md: 'flex',
          },
          flexDirection: 'column',
        }}
      >
        <Typography
          variant="body2"
          noWrap
          sx={{
            '& span': {
              display: { sm: 'none', md: 'inline' },
            },
          }}
        >
          {userData?.first_name}
          <span>{userData?.last_name ? ` ${userData?.last_name}` : ''}</span>
        </Typography>
        <StyledUpdatedText
          variant="subtitle1"
          sx={{
            display: { xs: 'none', sm: 'none', md: 'block' },
          }}
        >
          {userData?.email}
        </StyledUpdatedText>
      </StyledUserDetailsContainer>
    </StyledAvatarContainer>
  );
};

export default UserAvatarComponent;
