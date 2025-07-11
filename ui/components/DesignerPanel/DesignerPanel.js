import React from 'react';
import {
  Modal,
  ModalBody,
  Avatar,
  Typography,
  Button,
  Box,
  CustomTooltip,
  useTheme,
  CircularProgress,
  List,
  ListItem,
  ListItemAvatar,
  ListItemText,
  Divider,
  OpenInNewIcon,
  PersonIcon,
  DesignIcon,
} from '@sistent/sistent';
import { useGetUserProfileSummaryByIdQuery } from '@/rtk-query/user';
import { useGetDesignQuery, useGetUserDesignsQuery } from '@/rtk-query/design';
import { MESHERY_CLOUD_PROD } from '@/constants/endpoints';
import { openDesignInKanvas, useIsKanvasDesignerEnabled } from '@/utils/utils';
import { useRouter } from 'next/router';
import { styled } from '@sistent/sistent';

const StyledModalBody = styled(ModalBody)(({ theme }) => ({
  padding: theme.spacing(3),
  minWidth: '500px',
}));

const UserInfoSection = styled(Box)(({ theme }) => ({
  display: 'flex',
  alignItems: 'center',
  gap: theme.spacing(2),
  marginBottom: theme.spacing(3),
  padding: theme.spacing(2),
  borderRadius: theme.spacing(1),
  backgroundColor: theme.palette.background.card,
}));

const ActionButton = styled(Button)(({ theme }) => ({
  marginRight: theme.spacing(1),
  minWidth: '120px',
}));

const CurrentDesignSection = styled(Box)(({ theme }) => ({
  marginTop: theme.spacing(2),
  padding: theme.spacing(2),
  borderRadius: theme.spacing(1),
  backgroundColor: theme.palette.background.secondary,
}));

const DesignItem = styled(ListItem)(({ theme }) => ({
  borderRadius: theme.spacing(1),
  marginBottom: theme.spacing(1),
  backgroundColor: theme.palette.background.card,
  '&:hover': {
    backgroundColor: theme.palette.action.hover,
    cursor: 'pointer',
  },
}));

const DesignerPanel = ({ open, onClose, activeUser, currentDesignId }) => {
  const theme = useTheme();
  const router = useRouter();
  const isKanvasDesignerAvailable = useIsKanvasDesignerEnabled();

  const { data: userProfile, isLoading: isUserLoading } = useGetUserProfileSummaryByIdQuery(
    { id: activeUser?.user_id },
    { skip: !activeUser?.user_id },
  );
  const { data: designdata } = useGetDesignQuery(
    {
      design_id: currentDesignId,
    },
    { skip: !currentDesignId },
  );

  const { data: userDesigns } = useGetUserDesignsQuery(
    {
      user_id: activeUser?.user_id,
      page: 0,
      pagesize: 5,
      order: 'updated_at desc',
    },
    { skip: !activeUser?.user_id },
  );

  const handleGoToProfile = () => {
    if (userProfile?.id) {
      window.open(`${MESHERY_CLOUD_PROD}/user/${userProfile.id}`, '_blank');
    }
    onClose();
  };

  const handleOpenCurrentDesign = () => {
    if (currentDesignId && activeUser?.currentDesignName) {
      if (!isKanvasDesignerAvailable) {
        router.push(`/configuration/designs/configurator?design_id=${currentDesignId}`);
      } else {
        openDesignInKanvas(currentDesignId, activeUser.currentDesignName, router);
      }
      onClose();
    }
  };

  const handleOpenDesign = (design) => {
    if (!isKanvasDesignerAvailable) {
      router.push(`/configuration/designs/configurator?design_id=${design.id}`);
    } else {
      openDesignInKanvas(design.id, design.name, router);
    }
    onClose();
  };

  if (!activeUser) return null;

  return (
    <Modal open={open} closeModal={onClose} maxWidth="md" title="Designer Activity">
      <StyledModalBody>
        {isUserLoading ? (
          <Box display="flex" justifyContent="center" alignItems="center" minHeight="200px">
            <CircularProgress />
          </Box>
        ) : (
          <>
            {/* User Information Section */}
            <UserInfoSection>
              <Avatar
                src={userProfile?.avatar_url || activeUser?.avatar_url}
                sx={{ width: 64, height: 64 }}
              >
                {!userProfile?.avatar_url && !activeUser?.avatar_url && <PersonIcon />}
              </Avatar>
              <Box flex={1}>
                <Typography variant="h6" gutterBottom>
                  {userProfile?.first_name || activeUser?.name || 'Unknown User'}
                  {userProfile?.last_name && ` ${userProfile.last_name}`}
                </Typography>
                <Typography variant="body2" color="text.secondary" gutterBottom>
                  {userProfile?.email}
                </Typography>
              </Box>
              <ActionButton
                variant="outlined"
                startIcon={<PersonIcon />}
                onClick={handleGoToProfile}
              >
                View Profile
              </ActionButton>
            </UserInfoSection>

            {/* Current Design Section */}
            {designdata && (
              <CurrentDesignSection>
                <Typography variant="subtitle1" gutterBottom>
                  Currently Working On:
                </Typography>
                <DesignItem onClick={handleOpenCurrentDesign}>
                  <ListItemAvatar>
                    <Avatar sx={{ bgcolor: theme.palette.primary.main }}>
                      <DesignIcon />
                    </Avatar>
                  </ListItemAvatar>
                  <ListItemText
                    primary={designdata.name}
                    secondary={`Updated: ${new Date(designdata.updated_at).toLocaleDateString()}`}
                  />
                  <ActionButton
                    variant="contained"
                    startIcon={<DesignIcon />}
                    onClick={handleOpenCurrentDesign}
                  >
                    Collaborate
                  </ActionButton>
                </DesignItem>
              </CurrentDesignSection>
            )}

            {/* Recent Designs Section */}
            {userDesigns?.patterns && userDesigns.patterns.length > 0 && (
              <Box mt={3}>
                <Typography variant="subtitle1" gutterBottom>
                  Recent Designs:
                </Typography>
                <List>
                  {userDesigns.patterns.slice(0, 3).map((design) => (
                    <React.Fragment key={design.id}>
                      <DesignItem onClick={() => handleOpenDesign(design)}>
                        <ListItemAvatar>
                          <Avatar sx={{ bgcolor: theme.palette.secondary.main }}>
                            <DesignIcon />
                          </Avatar>
                        </ListItemAvatar>
                        <ListItemText
                          primary={design.name}
                          secondary={`Updated: ${new Date(design.updated_at).toLocaleDateString()}`}
                        />
                        <CustomTooltip title="Open in Designer">
                          <OpenInNewIcon />
                        </CustomTooltip>
                      </DesignItem>
                      <Divider />
                    </React.Fragment>
                  ))}
                </List>
              </Box>
            )}
          </>
        )}
      </StyledModalBody>
    </Modal>
  );
};

export default DesignerPanel;
