import {
  OutlinedPatternIcon,
  publishCatalogItemSchema,
  publishCatalogItemUiSchema,
  ResponsiveDataTable,
  useHasPermission,
} from '@sistent/sistent';
import { NoSsr } from '@sistent/sistent';
import React, { useCallback, useEffect, useRef, useState } from 'react';
import MesheryPatternGrid from './MesheryPatternGridView';
import _PromptComponent from '../../general/PromptComponent';
import LoadingScreen from '../../shared/LoadingState/LoadingComponent';
import { MesheryPatternsCatalog, VISIBILITY } from '../../../utils/Enum';
import { useRouter } from 'next/router';
import { useNotification } from '../../../utils/hooks/useNotification';
import _ from 'lodash';
import { getMeshModels } from '../../../api/meshmodel';
import { modifyRJSFSchema } from '../../../utils/utils';
import { updateVisibleColumns } from '../../../utils/responsive-column';
import { useWindowDimensions } from '../../../utils/dimension';
import { useTableUrlState } from '@/utils/hooks/useTableUrlState';
import { useColumnVisibilityPreference } from '@/utils/hooks/useColumnVisibilityPreference';
import InfoModal from '../../shared/Modal/Information/InfoModal';
import DefaultError from '../../general/error-404/index';

import { Keys } from '@meshery/schemas/permissions';
import { canEditDesign } from './design-permissions';
import ExportDesignModal from '../export/ExportDesignModal';
import { useModal, Modal as SistentModal } from '@sistent/sistent';
import PatternIcon from '@/assets/icons/Pattern';
import { useActorRef } from '@xstate/react';
import { designValidationMachine } from 'machines/validator/designValidator';
import {
  useClonePatternMutation,
  useDeletePatternFileMutation,
  useDeletePatternMutation,
  useDeployPatternMutation,
  useGetPatternsQuery,
  useImportPatternMutation,
  usePublishPatternMutation,
  useUndeployPatternMutation,
  useUnpublishPatternMutation,
  useUpdatePatternFileMutation,
  useUploadPatternFileMutation,
  useEvaluateRelationshipsMutation,
} from '@/rtk-query/design';
// import { useGetUserPrefQuery } from '@/rtk-query/user';
import { useGetProviderCapabilitiesQuery } from '@/rtk-query/user';
import { isLocalProvider } from '@/utils/provider';
import { useSelector } from 'react-redux';
import { ACTION_TYPES, resetSelectedPattern } from './MesheryPatterns.constants';
import YAMLEditor from './YAMLEditor';
import { ImportDesignModal } from '../ImportDesignModal';
import { PublishDesignModal } from '../PublishDesignModal';
import MesheryPatternsToolbar from './MesheryPatternsToolbar';
import {
  buildPatternColumns,
  buildPatternsTableOptions,
  PATTERN_COL_VIEWS,
} from './MesheryPatterns.columns';
import { buildDesignLifecycleHandlers } from './design-lifecycle-handlers';
import { createPatternsActions } from './patterns-actions';

function MesheryPatterns({
  disableCreateImportDesignButton = false,
  disableUniversalFilter = false,
  hideVisibility = false,
  initialFilters = { visibility: 'All' },
  pageTitle = 'Designs',
  arePatternsReadOnly = false,
}) {
  const canViewDesigns = useHasPermission(Keys.CatalogManagementViewDesigns);
  const canViewCatalog = useHasPermission(Keys.CatalogManagementViewCatalog);
  const canViewPage = canViewDesigns || (pageTitle === 'Catalog' && canViewCatalog);
  const canViewDesignDetails = useHasPermission(Keys.CatalogManagementDetailsOfDesign);
  const canPublishDesign = useHasPermission(Keys.CatalogManagementPublishDesign);
  const canImportDesign = useHasPermission(Keys.CatalogManagementImportDesign);
  const canEditDesignPermission = useHasPermission(Keys.CatalogManagementEditDesign);
  const canCloneDesign = useHasPermission(Keys.CatalogManagementCloneDesign);
  const canValidateDesign = useHasPermission(Keys.CatalogManagementValidateDesign);
  const canEvaluateRelationships = useHasPermission(Keys.CatalogManagementEvaluateRelationships);
  const canUndeployDesign = useHasPermission(Keys.CatalogManagementUndeployDesign);
  const canDeployDesign = useHasPermission(Keys.CatalogManagementDeployDesign);
  const canDownloadDesign = useHasPermission(Keys.CatalogManagementDownloadADesign);
  const canUnpublishDesign = useHasPermission(Keys.CatalogManagementUnpublishDesign);
  const router = useRouter();

  const { tableState, updateTableState } = useTableUrlState({
    tableKey: 'des',
    defaults: {
      page: 0,
      pageSize: 10,
      sortOrder: 'updated_at desc',
      search: '',
      filters: { vis: '' },
    },
  });

  const { page, pageSize, sortOrder, search } = tableState;
  const setPage = useCallback((p) => updateTableState({ page: p }), [updateTableState]);
  const setPageSize = useCallback((ps) => updateTableState({ pageSize: ps }), [updateTableState]);
  const setSortOrder = useCallback((so) => updateTableState({ sortOrder: so }), [updateTableState]);
  const setSearch = useCallback(
    (s) => updateTableState({ search: s, page: 0 }),
    [updateTableState],
  );

  // visibilityFilter is persisted in URL — no separate useState needed.
  const visibilityFilter = tableState.filters.vis || null;

  const [count, setCount] = useState(0);
  const modalRef = useRef();
  const [patterns, setPatterns] = useState([]);
  const [selectedRowData, setSelectedRowData] = useState(null);
  const [selectedPattern, setSelectedPattern] = useState(resetSelectedPattern());
  const [meshModels, setMeshModels] = useState([]);
  const [selectedFilters, setSelectedFilters] = useState(initialFilters);
  const [canPublishPattern, setCanPublishPattern] = useState(false);
  const [publishSchema, setPublishSchema] = useState({});
  const [infoModal, setInfoModal] = useState({
    open: false,
    ownerID: '',
    selectedResource: {},
  });
  const [isSearchExpanded, setIsSearchExpanded] = useState(false);
  const { view } = router.query;
  const [viewType, setViewType] = useState(view === 'table' ? 'table' : 'grid');
  const { notify } = useNotification();
  const { selectedK8sContexts, catalogVisibility, user } = useSelector((state) => state.ui);
  const [deployPatternMutation] = useDeployPatternMutation();
  const [undeployPatternMutation] = useUndeployPatternMutation();
  const {
    data: patternsData,
    isLoading: ispatternsLoading,
    refetch: getPatterns,
  } = useGetPatternsQuery({
    page: page,
    pagesize: pageSize,
    search: search,
    order: sortOrder,
    visibility: visibilityFilter ? JSON.stringify([visibilityFilter]) : '',
    populate: 'pattern_file',
  });
  const [clonePattern] = useClonePatternMutation();
  const [publishCatalog] = usePublishPatternMutation();
  const [unpublishCatalog] = useUnpublishPatternMutation();
  const [deletePattern] = useDeletePatternMutation();
  const [importPattern] = useImportPatternMutation();
  const [updatePattern] = useUpdatePatternFileMutation();
  const [uploadPatternFile] = useUploadPatternFileMutation();
  const [deletePatternFile] = useDeletePatternFileMutation();
  const [evaluateRelationships] = useEvaluateRelationshipsMutation();

  useEffect(() => {
    if (patternsData) {
      const filteredPatterns = patternsData.patterns.filter((content) => {
        if (visibilityFilter === null || content.visibility === visibilityFilter) {
          return true;
        }
        return false;
      });
      setCount(patternsData.totalCount || 0);
      handleSetPatterns(filteredPatterns);
      setPatterns(patternsData.patterns || []);
    }
  }, [patternsData]);

  const [importModal, setImportModal] = useState({
    open: false,
  });
  const [publishModal, setPublishModal] = useState({
    open: false,
    pattern: {},
    name: '',
  });

  const [downloadModal, setDownloadModal] = useState({
    open: false,
    content: null,
  });

  const designValidationActorRef = useActorRef(designValidationMachine);

  const designLifecycleModal = useModal({
    headerIcon: <PatternIcon fill="#fff" height={'2rem'} width={'2rem'} />,
  });
  const sistentInfoModal = useModal({
    headerIcon: OutlinedPatternIcon,
  });

  const handleDownloadDialogClose = () => {
    setDownloadModal((prevState) => ({
      ...prevState,
      open: false,
      content: null,
    }));
  };

  const handleDesignDownloadModal = (e, pattern) => {
    e.stopPropagation();
    setDownloadModal((prevState) => ({
      ...prevState,
      open: true,
      content: pattern,
    }));
  };

  const {
    handleError,
    resetSelectedRowData,
    handleDeploy,
    handleUndeploy,
    handleUploadImport,
    handleUploadImportClose,
    handleInfoModalClose,
    handleInfoModal,
    handleUnpublishModal,
    handlePublishModalClose,
    handlePublish,
    handleClone,
    handleSubmit,
    handleImportDesign,
    deletePatterns,
    handleDownload,
    handleEvaluateRelationship,
    showModal,
  } = createPatternsActions({
    clonePattern,
    publishCatalog,
    unpublishCatalog,
    deletePattern,
    deletePatternFile,
    importPattern,
    updatePattern,
    uploadPatternFile,
    deployPatternMutation,
    undeployPatternMutation,
    evaluateRelationships,
    modalRef,
    meshModels,
    infoModal,
    publishModal,
    user,
    setImportModal,
    setPublishModal,
    setSelectedRowData,
    setInfoModal,
    notify,
    sistentInfoModal,
    getPatterns,
  });

  // const [loading, stillLoading] = useState(true);
  const { width } = useWindowDimensions();

  const catalogVisibilityRef = useRef(false);
  const catalogContentRef = useRef();

  /**
   * Checking whether users are signed in under a provider that doesn't have
   * publish pattern capability and setting the canPublishPattern state accordingly
   */
  const { data: capabilitiesData } = useGetProviderCapabilitiesQuery();

  useEffect(() => {
    if (capabilitiesData) {
      // `capabilitiesData.capabilities` is the provider-declared feature list.
      const patternsCatalogueCapability = capabilitiesData?.capabilities?.filter(
        (val) => val.feature === MesheryPatternsCatalog,
      );
      if (patternsCatalogueCapability?.length > 0) setCanPublishPattern(true);
    }
  }, []);

  const searchTimeout = useRef(null);
  /**
   * fetch patterns when the page loads
   */
  useEffect(() => {
    document.body.style.overflowX = 'hidden';
    return () => {
      document.body.style.overflowX = 'auto';
    };
  }, []);

  useEffect(() => {
    if (viewType === 'grid') {
      setSearch('');
    }
    router.push(
      {
        pathname: router.pathname,
        query: { ...router.query, view: viewType },
      },
      undefined,
      { shallow: true },
    );
  }, [viewType]);

  useEffect(() => {
    const fetchMeshModels = async () => {
      try {
        const { models } = await getMeshModels();
        const modelNames = _.uniqBy(
          models?.map((model) => {
            if (model.displayName && model.displayName !== '') {
              return model.displayName;
            }
          }),
          _.toLower,
        );
        modelNames.sort();

        // Modify the schema using the utility function
        const modifiedSchema = modifyRJSFSchema(
          publishCatalogItemSchema,
          'properties.compatibility.items.enum',
          modelNames,
        );

        setPublishSchema({ rjsfSchema: modifiedSchema, uiSchema: publishCatalogItemUiSchema });
        setMeshModels(models);
      } catch (err) {
        console.error(err);
        handleError(ACTION_TYPES.SCHEMA_FETCH);
      }

      catalogVisibilityRef.current = catalogVisibility;
    };

    void fetchMeshModels();

    // Meshery's REST API already returns catalog items with `published`
    // visibility, so no separate catalog-pattern fetch is needed here.
  }, []);

  // useEffect(() => {
  //   handleSetPatterns(patterns);
  // }, [catalogVisibility]);

  const handleSetPatterns = (patterns) => {
    if (catalogVisibilityRef.current && catalogContentRef.current?.length > 0) {
      setPatterns([
        ...(catalogContentRef.current || []),
        ...(patterns?.filter((content) => content.visibility !== VISIBILITY.PUBLISHED) || []),
      ]);
      return;
    }

    setPatterns(patterns?.filter((content) => content.visibility !== VISIBILITY.PUBLISHED) || []);
  };

  const { openDeployModal, openUndeployModal, openDryRunModal, openValidateModal } =
    buildDesignLifecycleHandlers({
      designLifecycleModal,
      designValidationActorRef,
      selectedK8sContexts,
      handleDeploy,
      handleUndeploy,
    });

  const userCanEdit = (pattern) => canEditDesign(user, pattern, canEditDesignPermission);

  const handleOpenInConfigurator = (id) => {
    router.push('/configuration/designs/configurator?design_id=' + id);
  };

  const columns = buildPatternColumns({
    patterns,
    handlers: {
      handleOpenInConfigurator,
      handleClone,
      openValidateModal,
      openDryRunModal,
      openUndeployModal,
      openDeployModal,
      handleDesignDownloadModal,
      handleInfoModal,
      handleUnpublishModal,
      handleEvaluateRelationship,
      userCanEdit,
    },
    permissions: {
      editDesign: canEditDesignPermission,
      cloneDesign: canCloneDesign,
      validateDesign: canValidateDesign,
      evaluateRelationships: canEvaluateRelationships,
      undeployDesign: canUndeployDesign,
      deployDesign: canDeployDesign,
      downloadDesign: canDownloadDesign,
      detailsOfDesign: canViewDesignDetails,
      publishDesign: canPublishDesign,
      unpublishDesign: canUnpublishDesign,
    },
  });

  columns.forEach((column, idx) => {
    if (column.name === sortOrder.split(' ')[0]) {
      columns[idx].options.sortDirection = sortOrder.split(' ')[1];
    }
  });

  const [tableCols, updateCols] = useState(columns);

  const responsiveColDefaults = (() => {
    const showCols = updateVisibleColumns(PATTERN_COL_VIEWS, width);
    const initialVisibility = {};
    columns.forEach((col) => {
      if (!(hideVisibility && col.name === 'visibility')) {
        initialVisibility[col.name] = showCols[col.name];
      }
    });
    return initialVisibility;
  })();

  const { columnVisibility, setColumnVisibilityByUser, setColumnVisibilityByResponsive } =
    useColumnVisibilityPreference('designs', responsiveColDefaults);

  useEffect(() => {
    setColumnVisibilityByResponsive(responsiveColDefaults);
  }, [width, setColumnVisibilityByResponsive]);

  const options = buildPatternsTableOptions({
    patterns,
    columns,
    count,
    pageSize,
    page,
    search,
    sortOrder,
    isLocalProvider: isLocalProvider(capabilitiesData),
    searchTimeout,
    setPage,
    setPageSize,
    setSearch,
    setSortOrder,
    setSelectedRowData,
    deletePatterns,
    showModal,
  });

  if (ispatternsLoading) {
    return (
      <>
        <LoadingScreen animatedIcon="AnimatedMeshPattern" message={`Loading ${pageTitle}...`} />
      </>
    );
  }

  const filter = {
    visibility: {
      name: 'Visibility',
      //if catalog content is enabled, then show all filters including published otherwise only show public and private filters
      options: catalogVisibility
        ? [
            { label: 'Public', value: 'public' },
            { label: 'Private', value: 'private' },
            { label: 'Published', value: 'published' },
          ]
        : [
            { label: 'Public', value: 'public' },
            { label: 'Private', value: 'private' },
          ],
    },
  };

  const handleApplyFilter = () => {
    updateTableState({
      filters: { vis: selectedFilters.visibility === 'All' ? '' : selectedFilters.visibility },
      page: 0,
    });
  };

  return (
    <>
      <NoSsr>
        {canViewPage ? (
          <>
            {selectedRowData && Object.keys(selectedRowData).length > 0 && (
              <YAMLEditor
                pattern={selectedRowData}
                onClose={resetSelectedRowData()}
                onSubmit={handleSubmit}
                isReadOnly={arePatternsReadOnly}
              />
            )}
            <MesheryPatternsToolbar
              width={width}
              isSearchExpanded={isSearchExpanded}
              setIsSearchExpanded={setIsSearchExpanded}
              selectedPattern={selectedPattern}
              patterns={patterns}
              viewType={viewType}
              setViewType={setViewType}
              disableCreateImportDesignButton={disableCreateImportDesignButton}
              disableUniversalFilter={disableUniversalFilter}
              pageTitle={pageTitle}
              router={router}
              handleUploadImport={handleUploadImport}
              setSearch={setSearch}
              filter={filter}
              selectedFilters={selectedFilters}
              setSelectedFilters={setSelectedFilters}
              handleApplyFilter={handleApplyFilter}
              columns={columns}
              columnVisibility={columnVisibility}
              setColumnVisibility={setColumnVisibilityByUser}
            />
            {!selectedPattern.show && viewType === 'table' && (
              <>
                {/* <StyledRow> */}
                <ResponsiveDataTable
                  data={patterns}
                  columns={columns}
                  // @ts-ignore
                  options={options}
                  tableCols={tableCols}
                  updateCols={updateCols}
                  columnVisibility={columnVisibility}
                />
                {/* </StyledRow> */}

                <ExportDesignModal
                  downloadModal={downloadModal}
                  handleDownloadDialogClose={handleDownloadDialogClose}
                  handleDesignDownload={handleDownload}
                />
              </>
            )}
            {!selectedPattern.show && viewType === 'grid' && (
              // grid vieww
              <MesheryPatternGrid
                selectedK8sContexts={selectedK8sContexts}
                canPublishPattern={canPublishPattern}
                patterns={patterns}
                handlePublish={handlePublish}
                handleUnpublishModal={handleUnpublishModal}
                handleClone={handleClone}
                supportedTypes="null"
                handleSubmit={handleSubmit}
                setSelectedPattern={setSelectedPattern}
                selectedPattern={selectedPattern}
                pages={Math.ceil(count / pageSize)}
                setPage={setPage}
                selectedPage={page}
                patternErrors={[]}
                publishModal={publishModal}
                setPublishModal={setPublishModal}
                publishSchema={publishSchema}
                user={user}
                fetch={() => getPatterns()}
                handleInfoModal={handleInfoModal}
                handleEvaluateRelationship={handleEvaluateRelationship}
                openUndeployModal={openUndeployModal}
                openValidationModal={openValidateModal}
                openDryRunModal={openDryRunModal}
                openDeployModal={openDeployModal}
                hideVisibility={hideVisibility}
                arePatternsReadOnly={arePatternsReadOnly}
              />
            )}

            <SistentModal maxWidth="sm" {...designLifecycleModal}></SistentModal>
            <SistentModal {...sistentInfoModal}>
              {canViewDesignDetails && infoModal.open && (
                <InfoModal
                  infoModalOpen={true}
                  handleInfoModalClose={handleInfoModalClose}
                  selectedResource={infoModal.selectedResource}
                  resourceOwnerID={infoModal.ownerID}
                  patternFetcher={getPatterns}
                />
              )}
            </SistentModal>

            {canPublishPattern && publishModal.open && canPublishDesign && (
              <PublishDesignModal
                publishFormSchema={publishSchema}
                handleClose={handlePublishModalClose}
                title={publishModal.pattern?.name || ''}
                handleSubmit={handlePublish}
              />
            )}
            {importModal.open && canImportDesign && (
              <ImportDesignModal
                handleClose={handleUploadImportClose}
                handleImportDesign={handleImportDesign}
              />
            )}
            <_PromptComponent ref={modalRef} />
          </>
        ) : (
          <DefaultError />
        )}
      </NoSsr>
    </>
  );
}

// @ts-ignore
export default MesheryPatterns;
