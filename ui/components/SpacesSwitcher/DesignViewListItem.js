//@ts-check
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
} from '@layer5/sistent';
import { Lock, Public } from '@mui/icons-material';
import { VIEW_VISIBILITY } from '../Modals/Information/InfoModal';
import {
  StyledActionsContainer,
  StyledAvatarContainer,
  StyledListItem,
  StyledListItemText,
  StyledTextContainer,
  StyledUpdatedText,
  StyledUserDetailsContainer,
  StyledUserInfoContainer,
  StyledVisibilityContainer,
} from './styles';
import React from 'react';

const DesignViewListItem = ({
  selectedItem,
  handleItemClick,
  MenuComponent,
  onVisibilityChange,
  canChangeVisibility,
}) => {
  const { data: userData, isLoading: isUserLoading } = useGetUserProfileSummaryByIdQuery({
    id: selectedItem.user_id,
  });

  return (
    <>
      <StyledListItem
        data-testid={`designs-tr-${selectedItem.id}`}
        key={selectedItem.id}
        onClick={handleItemClick}
      >
        <StyledTextContainer>
          <StyledListItemText
            primary={selectedItem.name ?? 'Untitled'}
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
          <Skeleton animation="wave" height={24} width="80%" />
          <Skeleton animation="wave" height={16} width="40%" />
        </StyledTextContainer>

        <StyledUserInfoContainer>
          <AvatarSkeleton />
        </StyledUserInfoContainer>

        <StyledVisibilityContainer>
          <Skeleton animation="wave" height={32} width="70%" />
        </StyledVisibilityContainer>

        <StyledActionsContainer>
          {Array(4)
            .fill()
            .map((_, index) => (
              <Skeleton key={index} variant="circular" animation="wave" width={24} height={24} />
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
        height={32}
        width={32}
        style={{ minWidth: '32px' }}
      />
      <StyledUserDetailsContainer style={{ width: '70%' }}>
        <Skeleton animation="wave" height={20} width="60%" />
        <Skeleton animation="wave" height={16} width="80%" />
      </StyledUserDetailsContainer>
    </StyledAvatarContainer>
  );
};
