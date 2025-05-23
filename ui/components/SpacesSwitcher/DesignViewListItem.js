import { useGetUserProfileSummaryByIdQuery } from '@/rtk-query/user';
import {
  Divider,
  CustomTooltip,
  Skeleton,
  VisibilityChipMenu,
  getRelativeTime,
  getFullFormattedTime,
  AvatarGroup,
  FormControlLabel,
  Checkbox,
  FormGroup,
} from '@layer5/sistent';
import { Lock, Public } from '@mui/icons-material';
import { VIEW_VISIBILITY } from '../Modals/Information/InfoModal';
import {
  StyledActionsContainer,
  StyledAvatarContainer,
  StyledListIcon,
  StyledListItem,
  StyledListItemText,
  StyledSmallAvatar,
  StyledSmallAvatarContainer,
  StyledTextContainer,
  StyledUpdatedText,
  StyledUserDetailsContainer,
  StyledUserInfoContainer,
  StyledVisibilityContainer,
} from './styles';
import React, { useContext } from 'react';
import { iconMedium } from 'css/icons.styles';
import { RESOURCE_TYPE } from '@/utils/Enum';
import UserAvatarComponent from './UserAvatarComponent';
import { WorkspaceModalContext } from '@/utils/context/WorkspaceModalContextProvider';
import { Grid } from '@layer5/sistent';
import { useGetIconBasedOnMode } from './hooks';

const DesignViewListItem = ({
  selectedItem,
  handleItemClick,
  MenuComponent,
  onVisibilityChange,
  canChangeVisibility,
  type = RESOURCE_TYPE.DESIGN,
  activeUsers = [],
  isMultiSelectMode = false,
}) => {
  const { data: userData, isLoading: isUserLoading } = useGetUserProfileSummaryByIdQuery({
    id: selectedItem.user_id,
  });
  const { multiSelectedContent, setMultiSelectedContent } = useContext(WorkspaceModalContext);
  return (
    <>
      <StyledListItem
        data-testid={`${type}-tr-${selectedItem.id}`}
        key={selectedItem.id}
        onClick={handleItemClick}
        onMouseEnter={() => {
          const avatarHolder = document.getElementById(`${type}-avatar-${selectedItem.id}`);
          if (avatarHolder) {
            avatarHolder.style.clipPath = 'polygon(0 0, 100% 0, 100% 100%, 0 100%)';
            avatarHolder.style.transform = 'translate(0px,-100%)';
          }
        }}
        onMouseLeave={() => {
          const avatarHolder = document.getElementById(`${type}-avatar-${selectedItem.id}`);
          if (avatarHolder) {
            avatarHolder.style.transform = 'translate(0px,-12px)';
            avatarHolder.style.clipPath = 'polygon(0 0, 100% 0, 100% 40%, 0 40%)';
          }
        }}
      >
        {isMultiSelectMode && (
          <Grid item xs={0.3} md={0.25} zeroMinWidth>
            <FormGroup>
              <FormControlLabel
                control={
                  <Checkbox
                    checked={multiSelectedContent.some((item) => item.id === selectedItem.id)}
                    onChange={(e) => {
                      e.stopPropagation();
                      if (e.target.checked) {
                        setMultiSelectedContent((prev) => [...prev, selectedItem]);
                      } else {
                        setMultiSelectedContent((prev) =>
                          prev.filter((item) => item.id !== selectedItem.id),
                        );
                      }
                    }}
                    onClick={(e) => e.stopPropagation()}
                  />
                }
              />
            </FormGroup>
          </Grid>
        )}

        <StyledTextContainer item xs={6} md={5} lg={5}>
          <StyledAvatarContainer>
            <StyledListIcon>{useGetIconBasedOnMode({ mode: type })}</StyledListIcon>
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
          </StyledAvatarContainer>
        </StyledTextContainer>
        <StyledUserInfoContainer item xs={4} md={4} lg={4}>
          {isUserLoading ? <AvatarSkeleton /> : <UserAvatarComponent userData={userData} />}
        </StyledUserInfoContainer>
        <StyledVisibilityContainer
          item
          md={2}
          lg={1}
          sx={{ display: { xs: 'none', sm: 'none', md: 'block' }, minWidth: 0 }}
        >
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
        <StyledActionsContainer item xs={1} md={1} lg={2} zeroMinWidth>
          {MenuComponent}
        </StyledActionsContainer>
        {activeUsers && (
          <StyledSmallAvatarContainer
            id={`${type}-avatar-${selectedItem.id}`}
            transform="translate(0px, -12px)"
            clipPath="polygon(0 0, 100% 0, 100% 33%, 0 33%)"
          >
            <AvatarGroup
              max={7}
              className="root"
              componentsProps={{
                additionalAvatar: {
                  style: {
                    height: '23px',
                    width: '23px',
                  },
                },
              }}
            >
              {activeUsers.map((user) => (
                <CustomTooltip key={user.client_id} title={user.name}>
                  <StyledSmallAvatar
                    borderColor={user.color}
                    key={user.client_id}
                    alt={user.name}
                    src={user.avatar_url}
                    imgProps={{ referrerPolicy: 'no-referrer' }}
                  />
                </CustomTooltip>
              ))}
            </AvatarGroup>
          </StyledSmallAvatarContainer>
        )}
      </StyledListItem>
      <Divider light />
    </>
  );
};

export default DesignViewListItem;

export const DesignViewListItemSkeleton = ({ isMultiSelectMode = false }) => {
  return (
    <>
      <StyledListItem>
        {isMultiSelectMode && (
          <Grid item xs={0.3} md={0.25} zeroMinWidth>
            <Skeleton variant="rectangular" animation="wave" {...iconMedium} />
          </Grid>
        )}
        <StyledTextContainer item xs={6} md={5} lg={5}>
          <StyledAvatarContainer>
            <Skeleton variant="circular" animation="wave" width={24} height={24} />
            <div style={{ width: '100%', paddingLeft: '1rem' }}>
              <Skeleton animation="wave" height={24} width="80%" />
              <Skeleton animation="wave" height={16} width="40%" />
            </div>
          </StyledAvatarContainer>
        </StyledTextContainer>

        <StyledUserInfoContainer item xs={4} md={4} lg={4}>
          <AvatarSkeleton />
        </StyledUserInfoContainer>

        <StyledVisibilityContainer
          item
          md={2}
          lg={1}
          sx={{ display: { xs: 'none', sm: 'none', md: 'block' }, minWidth: 0 }}
        >
          <Skeleton animation="wave" height={32} width="70%" />
        </StyledVisibilityContainer>

        <StyledActionsContainer
          item
          xs={2}
          sm={3}
          md={1}
          lg={2}
          zeroMinWidth
          sx={{
            display: 'flex',
            alignItems: 'center',
            gap: '1rem',
            paddingLeft: '0.5rem',
          }}
        >
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
