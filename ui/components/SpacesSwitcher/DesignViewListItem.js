import { MESHERY_CLOUD_PROD } from '@/constants/endpoints';
import { useGetUserProfileSummaryByIdQuery } from '@/rtk-query/user';
import {
  Divider,
  CustomTooltip,
  Skeleton,
  Link,
  Avatar,
  VisibilityChipMenu,
  getRelativeTime,
  getFullFormattedTime,
  Typography,
  DesignIcon,
  ViewIcon,
  useTheme,
} from '@layer5/sistent';
import { Lock, Public } from '@mui/icons-material';
import { VIEW_VISIBILITY } from '../Modals/Information/InfoModal';
import {
  StyledActionsContainer,
  StyledAvatarContainer,
  StyledListIcon,
  StyledListItem,
  StyledListItemText,
  StyledTextContainer,
  StyledUpdatedText,
  StyledUserDetailsContainer,
  StyledUserInfoContainer,
  StyledVisibilityContainer,
} from './styles';
import React from 'react';
import { iconMedium } from 'css/icons.styles';
import { RESOURCE_TYPE } from '@/utils/Enum';

const DesignViewListItem = ({
  selectedItem,
  handleItemClick,
  MenuComponent,
  onVisibilityChange,
  canChangeVisibility,
  type = RESOURCE_TYPE.DESIGN,
}) => {
  const { data: userData, isLoading: isUserLoading } = useGetUserProfileSummaryByIdQuery({
    id: selectedItem.user_id,
  });
  const theme = useTheme();

  return (
    <>
      <StyledListItem
        data-testid={`designs-tr-${selectedItem.id}`}
        key={selectedItem.id}
        onClick={handleItemClick}
      >
        <StyledTextContainer>
          <StyledListIcon>
            {type === RESOURCE_TYPE.DESIGN ? (
              <DesignIcon />
            ) : (
              <ViewIcon {...iconMedium} fill={theme.palette.icon.brand} />
            )}
          </StyledListIcon>
          <StyledListItemText
            primary={selectedItem.name || ''}
            primaryTypographyProps={{ fontSize: '0.9rem' }}
            secondary={
              <CustomTooltip
                variant="small"
                title={getFullFormattedTime(selectedItem.updated_at)}
                placement="bottom"
              >
                <StyledUpdatedText>{getRelativeTime(selectedItem.updated_at)}</StyledUpdatedText>
              </CustomTooltip>
            }
          />
        </StyledTextContainer>

        <StyledUserInfoContainer>
          {isUserLoading ? (
            <AvatarSkeleton />
          ) : (
            <StyledAvatarContainer>
              <Link href={`https://${MESHERY_CLOUD_PROD}/user/${userData?.id}`}>
                <Avatar alt={userData?.first_name} src={userData?.avatar_url} />
              </Link>
              <StyledUserDetailsContainer>
                <Typography variant="body2">
                  {userData?.first_name}
                  {userData?.last_name ? ` ${userData?.last_name}` : ''}
                </Typography>
                <StyledUpdatedText variant="subtitle1">{userData?.email}</StyledUpdatedText>
              </StyledUserDetailsContainer>
            </StyledAvatarContainer>
          )}
        </StyledUserInfoContainer>

        <StyledVisibilityContainer>
          <VisibilityChipMenu
            value={selectedItem?.visibility}
            onChange={(value) => onVisibilityChange(value, selectedItem)}
            enabled={canChangeVisibility}
            options={[
              [VIEW_VISIBILITY.PUBLIC, Public],
              [VIEW_VISIBILITY.PRIVATE, Lock],
            ]}
          />
        </StyledVisibilityContainer>

        <StyledActionsContainer>{MenuComponent}</StyledActionsContainer>
      </StyledListItem>
      <Divider light />
    </>
  );
};

export default DesignViewListItem;

export const DesignViewListItemSkeleton = () => {
  return (
    <>
      <StyledListItem>
        <StyledTextContainer>
          <Skeleton variant="circular" animation="wave" width={24} height={24} />
          <div style={{ width: '100%', paddingLeft: '1rem' }}>
            <Skeleton animation="wave" height={24} width="80%" />
            <Skeleton animation="wave" height={16} width="40%" />
          </div>
        </StyledTextContainer>

        <StyledUserInfoContainer>
          <AvatarSkeleton />
        </StyledUserInfoContainer>

        <StyledVisibilityContainer>
          <Skeleton animation="wave" height={32} width="70%" />
        </StyledVisibilityContainer>

        <StyledActionsContainer style={{ gap: '0.75rem' }}>
          {Array(4)
            .fill()
            .map((_, index) => (
              <Skeleton key={index} variant="circular" animation="wave" {...iconMedium} />
            ))}
        </StyledActionsContainer>
      </StyledListItem>
      <Divider light />
    </>
  );
};

const AvatarSkeleton = () => {
  return (
    <StyledAvatarContainer>
      <Skeleton
        animation="wave"
        variant="circular"
        height={40}
        width={40}
        style={{ minWidth: '32px' }}
      />
      <StyledUserDetailsContainer style={{ width: '70%' }}>
        <Skeleton animation="wave" height={20} width="60%" />
        <Skeleton animation="wave" height={16} width="80%" />
      </StyledUserDetailsContainer>
    </StyledAvatarContainer>
  );
};
