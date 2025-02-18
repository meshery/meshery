import { useState } from 'react';
import { iconMedium } from 'css/icons.styles';
import { useTheme, ActionButtonCard, GetStartedModal, GetStartedIcon } from '@layer5/sistent';
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
import { keys } from '@/utils/permission_constants';
import { useLegacySelector } from 'lib/store';

const GetStarted = (props) => {
  const [openModal, setOpenModal] = useState(false);
  const theme = useTheme();
  const { data: currentUser } = useGetLoggedInUserQuery();
  const { data: profileData } = useGetUserByIdQuery(currentUser?.id);
  const currentOrg = useLegacySelector((state) => state.get('organization'));
  const { id: org_id } = currentOrg;
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
          keys.ASSIGN_USER_ROLES.action,
          keys.ASSIGN_USER_ROLES.subject,
        )}
        useLazyGetTeamsQuery={useLazyGetTeamsQuery}
        embedDesignPath="/static/img/getting-started/embedded-design-edge-stack.js"
        isFromMeshery={true}
      />
    </>
  );
};

export default GetStarted;
