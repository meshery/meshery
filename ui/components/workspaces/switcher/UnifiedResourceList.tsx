import React, { useCallback, useContext, useRef, useState } from 'react';
import { useRouter } from 'next/router';
import { useSelector } from 'react-redux';
import {
  ListItem,
  ListItemText,
  Divider,
  PromptComponent,
  useTheme,
  useRoomActivity,
  ExportIcon,
  DeleteIcon,
  ShareIcon,
  InfoIcon,
  WorkspaceContentMoveModal,
} from '@sistent/sistent';
import DesignViewListItem, { DesignViewListItemSkeleton } from './DesignViewListItem';
import useInfiniteScroll, {
  handleUpdateDesignVisibility,
  handleUpdateViewVisibility,
  useContentDelete,
  useContentDownload,
} from '@/components/workspaces/switcher/hooks';
import { MenuComponent } from './MenuComponent';
import { RESOURCE_TYPE } from '@/utils/Enum';
import { DesignList, LoadingContainer, GhostContainer, GhostImage, GhostText } from './styles';
import { useUpdatePatternFileMutation } from '@/rtk-query/design';
import { useUpdateViewVisibilityMutation } from '@/rtk-query/view';
import ShareModal from '@/components/workspaces/ShareWorkspaceModal';
import InfoModal from '@/components/shared/Modal/InfoModal';
import { ViewInfoModal } from '@/components/workspaces/ViewInfoModal';
import { openDesignInKanvas, openViewInKanvas, useIsOperatorEnabled } from '@/utils/utils';
import { useNotification } from '@/utils/hooks/useNotification';
import { EVENT_TYPES } from '@/lib/event-types';
import MoveFileIcon from '@/assets/icons/MoveFileIcon';
import { WorkspaceModalContext } from '@/utils/context/WorkspaceModalContextProvider';
import {
  useAssignDesignToWorkspaceMutation,
  useAssignViewToWorkspaceMutation,
  useGetWorkspacesQuery,
} from '@/rtk-query/workspace';
import { getUserAccessToken, getUserProfile, useGetLoggedInUserQuery } from '@/rtk-query/user';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';

interface UnifiedResourceListProps {
  type: RESOURCE_TYPE.DESIGN | RESOURCE_TYPE.VIEW;
  items: any[];
  page: number;
  setPage: (page: number) => void;
  isLoading: boolean;
  isFetching: boolean;
  hasMore: boolean;
  total_count: number;
  workspace?: any;
  refetch: () => void;
  isMultiSelectMode?: boolean;
  showWorkspaceName?: boolean;
  showOrganizationName?: boolean;
}

/**
 * Unified resource list component for both Designs and Views.
 * Replaces MainDesignsContent and MainViewsContent.
 */
const UnifiedResourceList: React.FC<UnifiedResourceListProps> = ({
  type,
  items,
  page,
  setPage,
  isLoading,
  isFetching,
  hasMore,
  total_count,
  workspace,
  refetch,
  isMultiSelectMode = false,
  showWorkspaceName = true,
  showOrganizationName = true,
}) => {
  const isDesign = type === RESOURCE_TYPE.DESIGN;
  const { data: currentUser } = useGetLoggedInUserQuery({});
  const [shareModal, setShareModal] = useState(false);
  const [infoModal, setInfoModal] = useState(false);
  const [moveModal, setMoveModal] = useState(false);
  const [selectedItem, setSelectedItem] = useState<any>(null);

  const [updatePattern] = useUpdatePatternFileMutation();
  const [updateView] = useUpdateViewVisibilityMutation();

  const handleOpenShareModal = (item: any) => {
    setSelectedItem(item);
    setShareModal(true);
  };

  const handleOpenInfoModal = (item: any) => {
    setSelectedItem(item);
    setInfoModal(true);
  };

  const handleMoveModal = (item: any) => {
    setSelectedItem(item);
    setMoveModal(true);
  };

  const modalRef = useRef<any>(null);
  const { handleDelete } = useContentDelete(modalRef);

  const loadNextPage = useCallback(() => {
    if (isLoading || isFetching) return;
    setPage(page + 1);
  }, [isLoading, isFetching, page, setPage]);

  const { loadingRef } = useInfiniteScroll({
    isLoading: isLoading || isFetching,
    hasMore,
    onLoadMore: loadNextPage,
  });

  const ghostRef = useRef(null);
  const ghostTextNodeRef = useRef(null);
  const theme = useTheme();
  const { handleDesignDownload, handleViewDownload } = useContentDownload();

  const ACTIONS = {
    EXPORT: {
      id: 'EXPORT',
      title: `Export ${isDesign ? 'Design' : 'View'}`,
      icon: <ExportIcon fill={theme.palette.icon.default} />,
      handler: (item: any) => (isDesign ? handleDesignDownload(item) : handleViewDownload(item)),
      enabled: () => true,
    },
    MOVE: {
      id: 'MOVE',
      title: `Move ${isDesign ? 'Design' : 'View'}`,
      icon: <MoveFileIcon fill={theme.palette.icon.default} />,
      enabled: () =>
        isDesign
          ? CAN(keys.ASSIGN_DESIGNS_TO_WORKSPACE.action, keys.ASSIGN_DESIGNS_TO_WORKSPACE.subject)
          : CAN(keys.ASSIGN_VIEWS_TO_WORKSPACE.action, keys.ASSIGN_VIEWS_TO_WORKSPACE.subject),
    },
    INFO: {
      id: 'INFO',
      title: `${isDesign ? 'Design' : 'View'} Info`,
      icon: <InfoIcon fill={theme.palette.icon.default} />,
      enabled: () => true,
    },
    SHARE: {
      id: 'SHARE',
      title: `Share ${isDesign ? 'Design' : 'View'}`,
      icon: <ShareIcon fill={theme.palette.icon.default} />,
      enabled: () => true,
    },
    DELETE: {
      id: 'DELETE',
      title: `Delete ${isDesign ? 'Design' : 'View'}`,
      icon: <DeleteIcon fill={theme.palette.icon.default} />,
      enabled: ({ item, userId }: any) =>
        (isDesign
          ? CAN(keys.DELETE_A_PATTERN.action, keys.DELETE_A_PATTERN.subject)
          : CAN(keys.DELETE_VIEW.action, keys.DELETE_VIEW.subject)) && item.user_id === userId,
    },
  };

  const getMenuOptions = (item: any) => {
    const options = [
      { ...ACTIONS.EXPORT, handler: () => ACTIONS.EXPORT.handler(item) },
      { ...ACTIONS.SHARE, handler: () => handleOpenShareModal(item) },
      { ...ACTIONS.INFO, handler: () => handleOpenInfoModal(item) },
      { ...ACTIONS.DELETE, handler: () => handleDelete([item], type, refetch) },
    ];
    if (workspace) {
      options.unshift({ ...ACTIONS.MOVE, handler: () => handleMoveModal(item) });
    }
    return options.filter((opt) => opt.enabled({ item, userId: currentUser?.id }));
  };

  const isKanvasDesignerAvailable = useIsOperatorEnabled();
  const workspaceSwitcherContext = useContext(WorkspaceModalContext);
  const { notify } = useNotification();
  const router = useRouter();

  const handleOpenInKanvas = (id: string, name: string) => {
    if (!isKanvasDesignerAvailable) {
      notify({ message: 'Kanvas Designer is not available', event_type: EVENT_TYPES.ERROR });
      return;
    }
    workspaceSwitcherContext?.closeModal?.();
    isDesign ? openDesignInKanvas(id, name, router) : openViewInKanvas(id, name, router);
  };

  const isInitialFetch = isFetching && page === 0;
  const isEmpty = total_count === 0;
  const shouldRenderItems = !isEmpty && !isInitialFetch;

  const { capabilitiesRegistry, organization: currentOrganization } = useSelector(
    (state: any) => state.ui,
  );
  const providerUrl = capabilitiesRegistry?.provider_url;
  const [activeUsers] = useRoomActivity({
    provider_url: providerUrl,
    getUserAccessToken,
    getUserProfile,
  });

  const [assignDesignToWorkspace] = useAssignDesignToWorkspaceMutation();
  const [assignViewToWorkspace] = useAssignViewToWorkspaceMutation();

  return (
    <>
      <DesignList data-testid="unified-resource-list">
        {shouldRenderItems &&
          items?.map((item) => {
            const isPublished = item?.visibility === 'published';
            const isOwner = currentUser?.id === item?.user_id;
            const canChangeVisibility = !isPublished && isOwner;

            return (
              <React.Fragment key={item?.id}>
                <DesignViewListItem
                  showOrganizationName={showOrganizationName}
                  showWorkspaceName={showWorkspaceName}
                  type={isDesign ? 'design' : 'view'}
                  selectedItem={item}
                  handleItemClick={() => handleOpenInKanvas(item?.id, item?.name)}
                  canChangeVisibility={canChangeVisibility}
                  onVisibilityChange={async (value, selected) => {
                    isDesign
                      ? await handleUpdateDesignVisibility({
                          value,
                          selectedResource: selected,
                          updatePattern,
                        })
                      : await handleUpdateViewVisibility({
                          value,
                          selectedResource: selected,
                          updateView,
                        });
                    refetch();
                  }}
                  MenuComponent={<MenuComponent options={getMenuOptions(item)} />}
                  activeUsers={activeUsers?.[item?.id]}
                  isMultiSelectMode={isMultiSelectMode}
                />
                <Divider light />
              </React.Fragment>
            );
          })}
        <LoadingContainer ref={loadingRef}>
          {isLoading || isInitialFetch ? (
            Array(10)
              .fill(0)
              .map((_, i) => (
                <DesignViewListItemSkeleton key={i} isMultiSelectMode={isMultiSelectMode} />
              ))
          ) : isFetching ? (
            <DesignViewListItemSkeleton isMultiSelectMode={isMultiSelectMode} />
          ) : null}
          {!hasMore && !isLoading && !isFetching && items?.length > 0 && !isEmpty && (
            <ListItemText
              secondary={`No more ${isDesign ? 'designs' : 'views'} to load`}
              style={{ padding: '1rem' }}
            />
          )}
        </LoadingContainer>

        {!isLoading && isEmpty && (
          <ListItem>
            <ListItemText
              primary={`No ${isDesign ? 'designs' : 'views'} found`}
              style={{ textAlign: 'center' }}
            />
          </ListItem>
        )}
      </DesignList>

      <GhostContainer ref={ghostRef}>
        <GhostImage src="/static/img/service-mesh-pattern.png" height={30} width={30} />
        <GhostText ref={ghostTextNodeRef}></GhostText>
      </GhostContainer>

      {shareModal && (
        <ShareModal resource={selectedItem} handleClose={() => setShareModal(false)} type={type} />
      )}

      {infoModal &&
        (isDesign ? (
          <InfoModal
            open={infoModal}
            closeModal={() => setInfoModal(false)}
            selectedResource={selectedItem}
          />
        ) : (
          <ViewInfoModal
            open={infoModal}
            closeModal={() => setInfoModal(false)}
            view_id={selectedItem?.id}
            view_name={selectedItem?.name}
            metadata={selectedItem?.metadata}
          />
        ))}

      {moveModal && (
        <WorkspaceContentMoveModal
          currentWorkspace={workspace}
          selectedContent={selectedItem}
          setWorkspaceContentMoveModal={setMoveModal}
          type={type}
          workspaceContentMoveModal={moveModal}
          refetch={refetch}
          useGetWorkspacesQuery={useGetWorkspacesQuery}
          WorkspaceModalContext={WorkspaceModalContext}
          assignDesignToWorkspace={assignDesignToWorkspace}
          assignViewToWorkspace={assignViewToWorkspace}
          isCreateWorkspaceAllowed={CAN(
            keys.CREATE_WORKSPACE.action,
            keys.CREATE_WORKSPACE.subject,
          )}
          isMoveDesignAllowed={CAN(
            keys.ASSIGN_DESIGNS_TO_WORKSPACE.action,
            keys.ASSIGN_DESIGNS_TO_WORKSPACE.subject,
          )}
          isMoveViewAllowed={CAN(
            keys.ASSIGN_VIEWS_TO_WORKSPACE.action,
            keys.ASSIGN_VIEWS_TO_WORKSPACE.subject,
          )}
          currentOrgId={currentOrganization?.id}
          notify={notify}
          router={router}
        />
      )}
      <PromptComponent ref={modalRef} />
    </>
  );
};

export default UnifiedResourceList;
