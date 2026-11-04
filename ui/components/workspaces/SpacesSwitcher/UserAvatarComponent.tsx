import { MESHERY_CLOUD_PROD } from '@/constants/endpoints';
import React from 'react';
import { CustomTooltip, Link, Avatar, Typography, useTheme } from '@sistent/sistent';
import { StyledAvatarContainer, StyledUserDetailsContainer, StyledUpdatedText } from './styles';

const UserAvatarComponent = ({ userData }) => {
  const theme = useTheme();
  return (
    <StyledAvatarContainer
      sx={{
        [theme.breakpoints.down('sm')]: {
          justifyContent: 'center',
        },
      }}
    >
      <CustomTooltip
        title={`${userData?.firstName} ${userData?.lastName || ''}`}
        sx={{ display: { xs: 'block', sm: 'block' } }}
      >
        <div>
          <Link href={`${MESHERY_CLOUD_PROD}/user/${userData?.id}`}>
            <Avatar alt={userData?.firstName} src={userData?.avatarUrl} />
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
            [theme.breakpoints.down('xl')]: {
              width: theme.spacing(12),
              overflow: 'hidden',
              textOverflow: 'ellipsis',
              whiteSpace: 'nowrap',
            },
          }}
        >
          {userData?.firstName}
          <span>{userData?.lastName ? ` ${userData?.lastName}` : ''}</span>
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
