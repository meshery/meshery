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
import DesignViewListItem from '../SpacesSwitcher/DesignViewListItem';
import { RESOURCE_TYPE } from '@/utils/Enum';

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
    <Modal open={open} closeModal={onClose} maxWidth="lg" title="Designer Activity">
      <StyledModalBody>
        {isUserLoading ? (
          <Box display="flex" justifyContent="center" alignItems="center" minHeight="200px">
            <CircularProgress />
          </Box>
        ) : (
          <>
            {/* Current Design Section */}
            {designdata && (
              <Box>
                <Typography variant="subtitle1" gutterBottom>
                  Currently Working On:
                </Typography>
                <DesignViewListItem
                  type={RESOURCE_TYPE.DESIGN}
                  selectedItem={designdata}
                  handleItemClick={handleOpenCurrentDesign}
                  showWorkspaceName={false}
                  showOrganizationName={false}
                  isMultiSelectMode={false}
                />
              </Box>
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
                      <DesignViewListItem
                        type={RESOURCE_TYPE.DESIGN}
                        selectedItem={design}
                        handleItemClick={() => handleOpenDesign(design)}
                        showWorkspaceName={false}
                        showOrganizationName={false}
                        isMultiSelectMode={false}
                      />
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
