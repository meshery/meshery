//@ts-check
import {
  getDesign,
  useDeletePatternFileMutation,
  useGetUserDesignsQuery,
  usePublishPatternMutation,
  useUpdatePatternFileMutation,
} from '@/rtk-query/design';
import { useGetLoggedInUserQuery } from '@/rtk-query/user';
import {
  styled,
  List,
  ListItem,
  ListItemText,
  Divider,
  CircularProgress,
  PROMPT_VARIANTS,
  PromptComponent,
  OutlinedPatternIcon,
  useModal,
  Modal,
  Box,
} from '@layer5/sistent';
import React, { useCallback, useContext, useEffect, useRef, useState } from 'react';
import DesignViewListItem, { DesignViewListItemSkeleton } from './DesignViewListItem';
import useInfiniteScroll, {
  getModelNamesBasedOnDisplayNames,
  handleUpdatePatternVisibility,
} from './hooks';
import MenuComponent from './MenuComponent';
import { MoreVert } from '@mui/icons-material';
import { DesignList, GhostContainer, GhostImage, GhostText, LoadingContainer } from './styles';
import ExportModal from '../ExportModal';
import { updateProgress } from 'lib/store';
import downloadContent from '@/utils/fileDownloader';
import { useNotification } from '@/utils/hooks/useNotification';
import { EVENT_TYPES } from 'lib/event-types';
import { RESOURCE_TYPE } from '@/utils/Enum';
import ShareModal from './ShareModal';
import InfoModal from '../Modals/Information/InfoModal';
import { useGetSchemaQuery } from '@/rtk-query/schema';
import { useGetMeshModelsQuery } from '@/rtk-query/meshModel';
import _ from 'lodash';
import { openDesignInKanvas, useIsKanvasDesignerEnabled } from '@/utils/utils';
import { WorkspaceSwitcherContext } from './WorkspaceSwitcher';
import Router from 'next/router';

const MainDesignsContent = ({ setPage, isLoading, isFetching, designs, hasMore, total_count }) => {
  const { data: currentUser } = useGetLoggedInUserQuery({});
  const [selectedDesign, setSelectedDesign] = useState(null);
  const [shareModal, setShareModal] = useState(false);
  const [infoModal, setInfoModal] = useState({ open: false, userId: '' });
  const [publishCatalog] = usePublishPatternMutation();

  const { notify } = useNotification();
  const modalRef = useRef(true);
  const [deletePatternFile] = useDeletePatternFileMutation();
  const loadNextPage = useCallback(() => {
    if (isLoading || isFetching) return;
    setPage((prevPage) => prevPage + 1);
  }, [isLoading, isFetching]);

  const { loadingRef } = useInfiniteScroll({
    isLoading: isLoading || isFetching,
    hasMore,
    onLoadMore: loadNextPage,
  });
  const [downloadModal, setDownloadModal] = useState({
    open: false,
    content: null,
  });

  const handleDesignDownloadModal = (design) => {
    setDownloadModal((prevState) => ({
      ...prevState,
      open: true,
      content: design,
    }));
  };
  const handleDownloadDialogClose = () => {
    setDownloadModal((prevState) => ({
      ...prevState,
      open: false,
      content: null,
    }));
  };

  const handleDownload = (e, design, source_type, params) => {
    e.stopPropagation();
    updateProgress({ showProgress: true });
    try {
      let id = design.id;
      let name = design.name;
      downloadContent({ id, name, type: 'pattern', source_type, params });
      updateProgress({ showProgress: false });
      notify({ message: `"${name}" design downloaded`, event_type: EVENT_TYPES.INFO });
    } catch (e) {
      console.error(e);
    }
  };

  const handleDelete = async (design) => {
    const response = await modalRef.current.show({
      title: `Delete catalog item?`,
      subtitle: `Are you sure you want to delete ${design?.name}?`,
      primaryOption: 'DELETE',
      variant: PROMPT_VARIANTS.DANGER,
      showInfoIcon:
        "Unpublishing a catolog item removes the item from the public-facing catalog (a public website accessible to anonymous visitors at meshery.io/catalog). The catalog item's visibility will change to either public (or private with a subscription). The ability to for other users to continue to access, edit, clone and collaborate on your content depends upon the assigned visibility level (public or private). Prior collaborators (users with whom you have shared your catalog item) will retain access. However, you can always republish it whenever you want. Remember: unpublished catalog items can still be available to other users if that item is set to public visibility. For detailed information, please refer to the [documentation](https://docs.meshery.io/concepts/designs).",
    });
    if (response === 'DELETE') {
      const selectedPattern = design;
      const { name, id } = selectedPattern;
      deletePatternFile({
        id: id,
      })
        .unwrap()
        .then(() => {
          notify({ message: `"${name}" Design deleted`, event_type: EVENT_TYPES.SUCCESS });
        })
        .catch(() => {
          notify({ message: `Unable to delete "${name}" Design`, event_type: EVENT_TYPES.ERROR });
        });
    }
  };

  const handleShare = async (design) => {
    setShareModal(true);
    setSelectedDesign(design);
  };

  const handleShareClose = () => {
    setShareModal(false);
    setSelectedDesign(null);
  };
  const sistentInfoModal = useModal({
    headerIcon: OutlinedPatternIcon,
  });

  const handleInfoModal = async (design) => {
    const selectedDesignWithPatternFile = await getDesign({
      design_id: design?.id,
    });

    setSelectedDesign(selectedDesignWithPatternFile?.data);

    sistentInfoModal.openModal({
      title: selectedDesign?.name,
    });
    setInfoModal({
      open: true,
      userId: selectedDesignWithPatternFile?.data?.user_id,
    });
  };
  const [formSchema, setFormSchema] = useState({});
  const { data: publishFormSchema } = useGetSchemaQuery(
    { schemaName: 'publish' },
    {
      skip: !infoModal.open,
    },
  );

  const { data: meshModelsData } = useGetMeshModelsQuery(
    {
      params: {
        pagesize: 'all',
        trim: true,
      },
    },
    {
      skip: !infoModal.open,
    },
  );
  const modifyRJSFSchema = (schema, propertyPath, newValue) => {
    const clonedSchema = _.cloneDeep(schema);
    _.set(clonedSchema, propertyPath, newValue);
    return clonedSchema;
  };

  useEffect(() => {
    if (infoModal.open && publishFormSchema && meshModelsData) {
      const modelNames = _.uniq(meshModelsData.models?.map((model) => model.displayName));

      const modifiedSchema = modifyRJSFSchema(
        publishFormSchema.rjsfSchema,
        'properties.compatibility.items.enum',
        modelNames,
      );

      setFormSchema({ rjsfSchema: modifiedSchema, uiSchema: publishFormSchema.uiSchema });
    }
  }, [publishFormSchema, meshModelsData]);

  const handleInfoModalClose = () => {
    setSelectedDesign(null);
    setInfoModal({ open: false, userId: '' });
  };

  const handlePublish = (formData) => {
    const compatibilityStore = _.uniqBy(meshModelsData, (model) => _.toLower(model.displayName))
      ?.filter((model) =>
        formData?.compatibility?.some((comp) => _.toLower(comp) === _.toLower(model.displayName)),
      )
      ?.map((model) => model.name);

    const payload = {
      id: infoModal.selectedDesign?.id,
      catalog_data: {
        ...formData,
        compatibility: compatibilityStore,
        type: _.toLower(formData?.type),
      },
    };
    updateProgress({ showProgress: true });
    publishCatalog({
      publishBody: JSON.stringify(payload),
    })
      .unwrap()
      .then(() => {
        updateProgress({ showProgress: false });
        if (currentUser.role_names.includes('admin')) {
          notify({
            message: `${selectedDesign?.name} Design Published`,
            event_type: EVENT_TYPES.SUCCESS,
          });
        } else {
          notify({
            message:
              'Design queued for publishing into Meshery Catalog. Maintainers notified for review',
            event_type: EVENT_TYPES.SUCCESS,
          });
        }
      })
      .catch(() => {
        updateProgress({ showProgress: false });
        notify({
          message: `Unable to publish ${selectedDesign?.name} Design`,
          event_type: EVENT_TYPES.ERROR,
        });
      });
  };
  const ghostRef = useRef(null);
  const ghostTextNodeRef = useRef(null);
  const [updatePatterns] = useUpdatePatternFileMutation();
  const isKanvasDesignerAvailable = useIsKanvasDesignerEnabled();
  const workspaceSwitcherContext = useContext(WorkspaceSwitcherContext);

  const handleOpenDesignInDesigner = (designId, designName) => {
    if (!isKanvasDesignerAvailable) {
      notify({
        message: 'Kanvas Designer is not available',
        event_type: EVENT_TYPES.ERROR,
      });
      return;
    }
    if (workspaceSwitcherContext?.closeModal) {
      workspaceSwitcherContext.closeModal();
    }

    openDesignInKanvas(designId, designName, Router);
  };
  return (
    <>
      <DesignList data-testid="designs-list-item">
        {(!isFetching || !isLoading) && total_count === 0 && (
          <ListItem>
            <ListItemText primary={`No Designs found`} style={{ textAlign: 'center' }} />
          </ListItem>
        )}

        {total_count !== 0 &&
          designs?.map((design) => {
            const isPublished = design?.visibility === 'published';
            const isOwner = currentUser?.id === design?.user_id;
            const canChangeVisibility = !isPublished && isOwner;

            return (
              <React.Fragment key={design?.id}>
                <DesignViewListItem
                  type="design"
                  selectedItem={design}
                  handleItemClick={() => {
                    handleOpenDesignInDesigner(design?.id, design?.name);
                  }}
                  canChangeVisibility={canChangeVisibility}
                  onVisibilityChange={(value, selectedItem) => {
                    handleUpdatePatternVisibility({
                      value: value,
                      selectedResource: selectedItem,
                      updatePatterns: updatePatterns,
                    });
                  }}
                  MenuComponent={
                    <MenuComponent
                      rowData={design}
                      visibility={design?.visibility}
                      items={[
                        {
                          deleteHandler: () => handleDelete(design),
                          downloadHandler: () => handleDesignDownloadModal(design),
                          shareHandler: () => handleShare(design),
                          infoHandler: () => handleInfoModal(design),
                        },
                      ]}
                    />
                  }
                />
                <Divider light />
              </React.Fragment>
            );
          })}

        <LoadingContainer ref={loadingRef}>
          {isLoading ? (
            Array(10)
              .fill()
              .map((_, index) => <DesignViewListItemSkeleton key={index} />)
          ) : isFetching ? (
            <DesignViewListItemSkeleton />
          ) : (
            <></>
          )}
          {!hasMore && designs?.length > 0 && total_count > 0 && (
            <ListItemText secondary={`No more designs to load`} sx={{ padding: '1rem' }} />
          )}
        </LoadingContainer>
      </DesignList>
      <GhostContainer ref={ghostRef}>
        <GhostImage src="/static/img/service-mesh-pattern.png" height={30} width={30} />
        <GhostText ref={ghostTextNodeRef}></GhostText>
      </GhostContainer>
      <ExportModal
        downloadModal={downloadModal}
        handleDownloadDialogClose={handleDownloadDialogClose}
        handleDesignDownload={handleDownload}
      />
      {shareModal && (
        <ShareModal
          resource={selectedDesign}
          handleClose={handleShareClose}
          type={RESOURCE_TYPE.DESIGN}
        />
      )}
      {infoModal.open && (
        <Modal {...sistentInfoModal}>
          <InfoModal
            handlePublish={handlePublish}
            infoModalOpen={infoModal.open}
            handleInfoModalClose={handleInfoModalClose}
            dataName={'patterns'}
            selectedResource={selectedDesign}
            resourceOwnerID={infoModal.userId}
            currentUser={currentUser}
            formSchema={formSchema}
            meshModels={meshModelsData?.models}
          />
        </Modal>
      )}
      <PromptComponent ref={modalRef} />
    </>
  );
};

export default MainDesignsContent;
