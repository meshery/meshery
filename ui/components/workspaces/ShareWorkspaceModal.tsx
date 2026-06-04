/**
 * Workspace share modal wrapper.
 *
 * Thin adapter that wires Meshery's RTK Query hooks (resource access actors,
 * owner data, visibility mutations) into Sistent's `ShareModal` (aliased here
 * as `CatalogShare`). The dialog itself is rendered by Sistent, so this
 * wrapper does not compose a project-side modal primitive — doing so would
 * double-wrap the dialog and break the sharing behaviour.
 */
import { FC } from 'react';
import { ShareModal as CatalogShare } from '@sistent/sistent';
import { useUpdatePatternFileMutation } from '@/rtk-query/design';
import {
  useCreateAndRevokeResourceAccessRecordMutation,
  useGetAccessActorsInfoOfResourceQuery,
} from '@/rtk-query/resource';
import {
  useGetAllUsersQuery,
  useGetLoggedInUserQuery,
  useGetUserProfileSummaryByIdQuery,
} from '@/rtk-query/user';
import { useUpdateViewVisibilityMutation } from '@/rtk-query/view';
import { useNotification } from '@/utils/hooks/useNotification';
import { getShareableResourceRoute } from './SpacesSwitcher/hooks';
import { JsonParse } from '@/utils/utils';
import { MESHERY_CLOUD_PROD } from '@/constants/endpoints';
import { RESOURCE_TYPE } from '@/utils/Enum';

type Resource = {
  id?: string;
  name?: string;
  userId?: string;
  catalogData?: unknown;
  patternFile?: string;
};

type ShareWorkspaceModalInternalProps = {
  selectedResource: Resource | Resource[];
  dataName: string;
  handleShareModalClose: () => void;
};

export const ShareWorkspaceModal_: FC<ShareWorkspaceModalInternalProps> = ({
  selectedResource,
  dataName,
  handleShareModalClose,
}) => {
  const firstSelectedResource = Array.isArray(selectedResource)
    ? selectedResource[0]
    : selectedResource;
  const resourceType = dataName === 'design' ? 'pattern' : dataName;

  const { data: ownerData, isSuccess: isOwnerDataFetched } = useGetUserProfileSummaryByIdQuery(
    { id: firstSelectedResource?.userId },
    { skip: !firstSelectedResource?.userId },
  );

  const { data: accessActorsInfoOfResource, isSuccess: accessActorsFetched } =
    useGetAccessActorsInfoOfResourceQuery({
      resourceId: firstSelectedResource?.id,
      resourceType: resourceType,
      actorType: 'users',
    });

  const { notify } = useNotification();
  const { data: currentUser } = useGetLoggedInUserQuery();
  const [updatePatterns] = useUpdatePatternFileMutation();
  const [handleResourceShare] = useCreateAndRevokeResourceAccessRecordMutation();
  const [updateView] = useUpdateViewVisibilityMutation();

  const fetchAccessActors = () => {
    if (accessActorsFetched && isOwnerDataFetched) {
      return [ownerData, ...accessActorsInfoOfResource.users];
    }
    return [];
  };

  const shareableLink = getShareableResourceRoute(
    dataName,
    firstSelectedResource?.id,
    firstSelectedResource?.name,
  );

  const handleUpdatePatternVisibility = async (value: string) => {
    const res = await updatePatterns({
      body: {
        id: firstSelectedResource?.id,
        name: firstSelectedResource?.name,
        catalogData: firstSelectedResource?.catalogData,
        designFile: JsonParse(firstSelectedResource?.patternFile),
        visibility: value,
      },
    });
    return {
      error: res?.error?.error,
    };
  };

  const handleUpdateViewVisibilityCall = async (value: string) => {
    const res = await updateView({
      id: firstSelectedResource?.id,
      body: {
        visibility: value,
      },
    });
    return {
      error: res.error?.error,
    };
  };

  const handleUpdateVisibility = (updatedVisibility: string) => {
    if (dataName === RESOURCE_TYPE.DESIGN) {
      return handleUpdatePatternVisibility(updatedVisibility);
    }
    if (dataName === RESOURCE_TYPE.VIEW) {
      return handleUpdateViewVisibilityCall(updatedVisibility);
    }
    return Promise.resolve({ error: undefined });
  };

  return (
    <CatalogShare
      selectedResource={selectedResource}
      dataName={dataName}
      resourceAccessMutator={handleResourceShare}
      notify={notify}
      handleShareModalClose={handleShareModalClose}
      currentUser={currentUser}
      shareableLink={shareableLink}
      useGetAllUsersQuery={useGetAllUsersQuery}
      ownerData={{ ...ownerData, userId: ownerData?.id }}
      fetchAccessActors={fetchAccessActors}
      hostURL={MESHERY_CLOUD_PROD}
      handleUpdateVisibility={handleUpdateVisibility}
      accessActorsInfoOfResource={accessActorsInfoOfResource?.users}
    />
  );
};

export interface ShareWorkspaceModalProps {
  resource: Resource | Resource[];
  type: string;
  handleClose: () => void;
}

const ShareWorkspaceModal: FC<ShareWorkspaceModalProps> = ({ resource, type, handleClose }) => {
  return (
    <ShareWorkspaceModal_
      handleShareModalClose={handleClose}
      selectedResource={resource}
      dataName={type}
    />
  );
};

export default ShareWorkspaceModal;
