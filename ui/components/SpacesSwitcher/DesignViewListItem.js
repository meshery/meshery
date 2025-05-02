//@ts-check
import { MESHERY_CLOUD_PROD } from '@/constants/endpoints';
import { useGetUserProfileSummaryByIdQuery } from '@/rtk-query/user';
import {
  Divider,
  ListItem,
  ListItemText,
  styled,
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

const StyledListItem = styled(ListItem)(({ theme }) => ({
  display: 'flex',
  alignItems: 'center',
  gap: '1.5rem',
  marginBlock: '0',
  paddingBlock: '6px',
  position: 'relative',
  cursor: 'grab',
  '&:hover': {
    backgroundColor: theme.palette.background.hover,
    '& .menu-component': {
      opacity: 1,
      visibility: 'visible',
    },
  },
}));

const MainMenuComponent = styled('div')({
  opacity: 0,
  visibility: 'hidden',
  transition: 'all 0.2s ease-in-out',
  cursor: 'pointer',
});

const TextContainer = styled('div')({
  width: '45%',
  whiteSpace: 'nowrap',
  overflow: 'hidden',
  textOverflow: 'ellipsis',
});

const StyledListItemText = styled(ListItemText)({
  cursor: 'grab',
  width: '100%',
  whiteSpace: 'nowrap',
  overflow: 'hidden',
  textOverflow: 'ellipsis',
  margin: '0',
});

const UpdatedText = styled('p')({
  margin: '0',
  fontSize: '0.8rem',
  fontStyle: 'italic',
  color: '#647881',
  cursor: 'grab',
});

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
        <TextContainer>
          <StyledListItemText
            primary={selectedItem.name ?? 'Untitled'}
            primaryTypographyProps={{ fontSize: '0.9rem' }}
            secondary={
              <CustomTooltip
                variant="small"
                title={getFullFormattedTime(selectedItem.updated_at)}
                placement="bottom"
              >
                <UpdatedText>{getRelativeTime(selectedItem.updated_at)}</UpdatedText>
              </CustomTooltip>
            }
          />
        </TextContainer>

        <div style={{ display: 'flex', gap: '1rem', width: '30%' }}>
          {isUserLoading ? (
            <div style={{ width: '100%' }}>
              <Skeleton
                animation="wave"
                variant="circular"
                height={32}
                width={32}
                style={{ minWidth: '32px' }}
              />
            </div>
          ) : (
            <div style={{ display: 'flex', alignItems: 'center', width: '100%' }}>
              <Link href={`https://${MESHERY_CLOUD_PROD}/user/${userData?.id}`}>
                <Avatar alt={userData?.first_name} src={userData?.avatar_url} />
              </Link>
              <div
                style={{
                  display: 'flex',
                  alignItems: 'start',
                  flexDirection: 'column',
                  marginLeft: '1rem',
                  gap: '0.1rem',
                }}
              >
                <Typography variant="body2">
                  {userData?.first_name}
                  {userData?.last_name ? ` ${userData?.last_name}` : ''}
                </Typography>
                <UpdatedText variant="subtitle1">{userData?.email}</UpdatedText>
              </div>
            </div>
          )}
        </div>
        <div style={{ width: '10%' }}>
          <VisibilityChipMenu
            value={selectedItem?.visibility}
            onChange={(value) => onVisibilityChange(value, selectedItem)}
            enabled={canChangeVisibility}
            options={[
              [VIEW_VISIBILITY.PUBLIC, Public],
              [VIEW_VISIBILITY.PRIVATE, Lock],
            ]}
          />
        </div>
        <div style={{ width: '15%' }}>
          {/* <MainMenuComponent id={`menu-${selectedItem.id}`} className="menu-component"> */}
          {MenuComponent}
          {/* </MainMenuComponent> */}
        </div>
      </StyledListItem>
      <Divider light />
    </>
  );
};

export default DesignViewListItem;
