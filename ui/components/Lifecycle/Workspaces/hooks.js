import { useDeletePatternMutation } from '@/rtk-query/design';
import { usePublishFilterMutation } from '@/rtk-query/filter';
import { useUnassignDesignFromWorkspaceMutation } from '@/rtk-query/workspace';
import { useNotification, useNotificationHandlers } from '@/utils/hooks/useNotification';
import { EVENT_TYPES } from 'lib/event-types';
import _ from 'lodash';
import { useRouter } from 'next/router';

export const useDeletePattern = () => {
  const [deletePatternById] = useDeletePatternMutation();
  const { notify } = useNotification();
  const [unassignDesignFromWorkspace] = useUnassignDesignFromWorkspaceMutation();
  const router = useRouter();

  const deletePattern = async (pattern) => {
    return deletePatternById({
      id: pattern.id,
    })
      .unwrap()
      .then(() => {
        notify({
          message: `"${pattern?.name}" deleted`,
          event_type: EVENT_TYPES.SUCCESS,
        });
        router.push({ pathname: '/catalog/my-designs' });
      })
      .catch((error) =>
        notify({
          message: error?.data?.error || 'Error deleting design',
          event_type: EVENT_TYPES.ERROR,
        }),
      );
  };

  const handleDeleteModal = (pattern, modalRef) => {
    return async () => {
      let response = await modalRef.current.show({
        title: `Delete Catalog item?`,
        subtitle: `Are you sure that you want to delete "${pattern?.name}"?`,
        primaryOption: 'Delete',
        variant: 'error',
      });
      if (response === 'Delete') {
        await deletePattern(pattern);
      }
    };
  };

  const handleWorkspaceDesignDeleteModal = (designId, workspaceId) => {
    unassignDesignFromWorkspace({
      workspaceId,
      designId: designId,
    });
  };

  const handleBulkDeleteModal = async (patterns, modalRef) => {
    let response = await modalRef.current.show({
      title: `Delete Selected Catalog items?`,
      subtitle: `Are you sure that you want to delete ${patterns.length} designs?`,
      primaryOption: 'Delete',
    });
    if (response === 'Delete') {
      await Promise.allSettled(patterns.map((pattern) => deletePattern(pattern)));
    }
  };

  const handleBulkWorkspaceDesignDeleteModal = async (
    patterns,
    modalRef,
    workspaceName,
    workspaceId,
  ) => {
    let response = await modalRef.current.show({
      title: `Delete Selected Designs?`,
      subtitle: `Are you sure that you want to delete ${patterns.length} designs from workspace "${workspaceName}"?`,
      primaryOption: 'Delete',
    });

    if (response === 'Delete') {
      await Promise.allSettled(
        patterns.map((design) =>
          unassignDesignFromWorkspace({
            workspaceId,
            designId: design.id,
          }),
        ),
      );
    }
  };

  return {
    deletePattern,
    handleDeleteModal,
    handleWorkspaceDesignDeleteModal,
    handleBulkDeleteModal,
    handleBulkWorkspaceDesignDeleteModal,
  };
};

export const usePublishPattern = (meshModelModelsData, refetchPatternData) => {
  const [publishCatalogContent] = usePublishFilterMutation();
  const { handleSuccess, handleInfo, handleError } = useNotificationHandlers();

  const handlePublishModal = (pattern) => {
    return pattern;
  };

  const handlePublish = (publishModal, data) => {
    const compatibilityStore = _.uniqBy(meshModelModelsData?.models, (model) =>
      _.toLower(model.display_name),
    )
      ?.filter((model) =>
        data?.compatibility?.some((comp) => _.toLower(comp) === _.toLower(model.display_name)),
      )
      ?.map((model) => model.name);

    const payload = {
      id: publishModal.pattern.id,
      catalog_type: 'pattern',
      catalog_data: {
        ...data,
        type: _.toLower(data?.type),
        compatibility: compatibilityStore,
      },
    };

    return publishCatalogContent({
      type: 'pattern',
      mesheryCatalogRequestBody: payload,
    })
      .unwrap()
      .then((res) => {
        if (_.toLower(res.status) === 'pending') {
          handleInfo(
            `${publishModal?.pattern?.name} design queued for publishing into the catalog. Maintainers notified for review.`,
          );
        } else if (_.toLower(res.status) === 'approved') {
          handleSuccess(`${publishModal?.pattern?.name} published successfully`);
        }
        refetchPatternData();
      })
      .catch((error) => handleError(error.data));
  };

  return {
    handlePublish,
    handlePublishModal,
  };
};
