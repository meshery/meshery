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
} from '@layer5/sistent';
import { Lock, Public } from '@mui/icons-material';

import moment from 'moment';
import { VIEW_VISIBILITY } from '../Modals/Information/InfoModal';

const StyledListItem = styled(ListItem)(({ theme }) => ({
  display: 'grid',
  gridTemplateColumns: '0 1fr auto',
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
  width: '100%',
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

const AvatarContainer = styled('div')(({ transform, clipPath }) => ({
  position: 'absolute',
  bottom: '-1.65rem',
  left: '2rem',
  padding: '0rem',
  transition: 'all 0.4s ease',
  transform: transform,
  clipPath: clipPath,
  zIndex: 1,
}));

const StyledAvatar = ({ borderColor }) => ({
  height: '23px',
  width: '23px',
  border: borderColor ? `2px solid ${borderColor}` : undefined,
});

const DesignViewListItem = ({
  selectedItem,
  handleItemClick,
  //   onRowClick,
  //   ghostTextNodeRef,
  //   ghostRef,
  MenuComponent,
  onVisibilityChange,
  canChangeVisibility,
  //   users,
}) => {
  //   const dispatchCmdToEditor = useDispatchCmdToEditor();
  const formatTooltipDateTime = (date) => {
    return moment(date).format('MMM DD, YYYY HH:MM');
  };
  const getDateTime = (date) => {
    return `updated ${moment(date).fromNow()}`;
  };
  const { data: userData, isLoading: isUserLoading } = useGetUserProfileSummaryByIdQuery({
    id: selectedItem.user_id,
  });

  return (
    <>
      <StyledListItem
        data-testid={`designs-tr-${selectedItem.id}`}
        key={selectedItem.id}
        draggable={true}
        onClick={handleItemClick}
        // onDragStart={(e) => {
        //   if (selectedItem?.visibility === 'published') {
        //     onRowClick?.(selectedItem, true);
        //     return;
        //   }
        //   ghostTextNodeRef.current.innerHTML = selectedItem.name;

        //   dispatchCmdToEditor(
        //     uiCmds.startDraggingDesign({
        //       id: selectedItem.id,
        //     }),
        //   );
        //   e.dataTransfer.setDragImage(ghostRef.current, 0, 0);
        //   e.dataTransfer.setData('text/plain', `pattern/${selectedItem.id}`);
        // }}
        // onDragEnd={() => {
        //   dispatchCmdToEditor(uiCmds.stopDragging());
        // }}
        // onMouseEnter={() => {
        //   const menuComp = document.getElementById(`menu-${selectedItem.id}`);
        //   const avatarHolder = document.getElementById(`avatar-${selectedItem.id}`);
        //   if (menuComp) {
        //     menuComp.style.opacity = '1';
        //     menuComp.style.visibility = 'visible';
        //     menuComp.style.pointerEvents = 'auto';
        //   }
        //   if (avatarHolder) {
        //     avatarHolder.style.clipPath = 'polygon(0 0, 100% 0, 100% 100%, 0 100%)';
        //     avatarHolder.style.transform = 'translate(0px,-100%)';
        //   }
        // }}
        // onMouseLeave={() => {
        //   const menuComp = document.getElementById(`menu-${selectedItem.id}`);
        //   const avatarHolder = document.getElementById(`avatar-${selectedItem.id}`);
        //   if (menuComp) {
        //     menuComp.style.opacity = '0';
        //     menuComp.style.visibility = 'hidden';
        //     menuComp.style.pointerEvents = 'none';
        //   }
        //   if (avatarHolder) {
        //     avatarHolder.style.transform = 'translate(0px,-12px)';
        //     avatarHolder.style.clipPath = 'polygon(0 0, 100% 0, 100% 40%, 0 40%)';
        //   }
        // }}
      >
        <MainMenuComponent id={`menu-${selectedItem.id}`} className="menu-component">
          {/* {MenuComponent(selectedItem)} */}
          {MenuComponent}
        </MainMenuComponent>

        <TextContainer>
          <StyledListItemText
            primary={selectedItem.name ?? 'Untitled'}
            primaryTypographyProps={{ fontSize: '0.9rem' }}
            secondary={
              <CustomTooltip
                variant="small"
                title={formatTooltipDateTime(selectedItem.updated_at)}
                placement="bottom"
              >
                <UpdatedText>{getDateTime(selectedItem.updated_at)}</UpdatedText>
              </CustomTooltip>
            }
          />
        </TextContainer>

        <div style={{ display: 'flex', gap: '1rem' }}>
          {isUserLoading ? (
            <Skeleton
              animation="wave"
              variant="circular"
              height={32}
              width={32}
              style={{ minWidth: '32px' }}
            />
          ) : (
            <CustomTooltip title={userData?.first_name + ' ' + userData?.last_name}>
              <div>
                <Link href={`https://${MESHERY_CLOUD_PROD}/user/${userData?.id}`}>
                  <Avatar alt={userData?.first_name} src={userData?.avatar_url} />
                </Link>
              </div>
            </CustomTooltip>
          )}

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

        {/* {users && (
          <AvatarContainer
            id={`avatar-${selectedItem.id}`}
            transform="translate(0px, -12px)"
            clipPath="polygon(0 0, 100% 0, 100% 33%, 0 33%)"
          >
            <AvatarGroup
              max={3}
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
              {users.map((user) => (
                <CustomTooltip key={user.client_id} title={user.name}>
                  <StyledAvatar
                    borderColor={user.color}
                    key={user.client_id}
                    alt={user.name}
                    src={user.avatar_url}
                    imgProps={{ referrerPolicy: 'no-referrer' }}
                  />
                </CustomTooltip>
              ))}
            </AvatarGroup>
          </AvatarContainer>
        )} */}
      </StyledListItem>
      <Divider light />
    </>
  );
};

export default DesignViewListItem;
