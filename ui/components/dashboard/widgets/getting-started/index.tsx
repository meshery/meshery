import { useState } from 'react';
import { iconMedium } from 'css/icons.styles';
import { useTheme, ActionButtonCard, GetStartedModal, GetStartedIcon } from '@sistent/sistent';
import {
  useGetLoggedInUserQuery,
  useGetUserByIdQuery,
  useHandleUserInviteMutation,
  useLazyGetTeamsQuery,
  useUpdateUserPrefMutation,
} from '@/rtk-query/user';
import { stepsData } from './data';
import { useNotificationHandlers } from '@/utils/hooks/useNotification';
import { useGetUserOrgRolesQuery } from '@/rtk-query/orgRoles';
import { useGetOrgsQuery } from '@/rtk-query/organization';
import CAN from '@/utils/can';
import { Keys } from '@meshery/schemas/permissions';
import { useSelector } from 'react-redux';
import type { RootState } from '../../../../store';

const GetStarted = (props: { iconsProps?: object }) => {
  const [openModal, setOpenModal] = useState(false);
  const theme = useTheme();
  const { data: currentUser } = useGetLoggedInUserQuery();
  const { data: profileData } = useGetUserByIdQuery(currentUser?.id, {
    skip: !currentUser?.id,
  });
  const { organization: currentOrg } = useSelector((state: RootState) => state.ui);
  const org_id = currentOrg?.id;
  return (
    <>
      <ActionButtonCard
        title="GETTING STARTED"
        description="New here? Follow along these guided tasks to help you get the most of your account."
        onClick={() => setOpenModal(true)}
        profileData={profileData}
        btnTitle="Start"
        icon={
          <GetStartedIcon {...props.iconsProps} {...iconMedium} fill={theme.palette.icon.default} />
        }
        showProgress={true}
        completedSteps={profileData?.preferences?.remoteProviderPreferences?.getstarted || []}
        totalSteps={stepsData.length}
      />

      <GetStartedModal
        open={openModal}
        handleClose={() => setOpenModal(false)}
        handleOpen={() => setOpenModal(true)}
        stepsData={stepsData}
        profileData={profileData}
        useUpdateUserPrefMutation={useUpdateUserPrefMutation}
        currentOrgId={org_id}
        useGetOrgsQuery={useGetOrgsQuery}
        useGetUserOrgRolesQuery={useGetUserOrgRolesQuery}
        useHandleUserInviteMutation={useHandleUserInviteMutation}
        useNotificationHandlers={useNotificationHandlers}
        isAssignUserRolesAllowed={CAN(
          Keys.IdentityAccessManagementAssignUserRoles.id,
          Keys.IdentityAccessManagementAssignUserRoles.function,
        )}
        useLazyGetTeamsQuery={useLazyGetTeamsQuery}
        embedDesignPath="/static/img/getting-started/embedded-design-edge-stack.js"
        isFromMeshery={true}
      />
    </>
  );
};

export default GetStarted;
