import React, { useState, useEffect, useRef, useCallback } from 'react';
import { useTableUrlState } from '@/utils/hooks/useTableUrlState';
import { useColumnVisibilityPreference } from '@/utils/hooks/useColumnVisibilityPreference';
import { NoSsr } from '@sistent/sistent';
import { Publish as PublishIcon } from '@/assets/icons';
import _PromptComponent from '../general/PromptComponent';
import { MesheryFiltersCatalog, VISIBILITY } from '../../utils/Enum';
import ViewSwitch from '../general/ViewSwitch';
import FiltersGrid from './FiltersGrid';
import fetchCatalogFilter from '@/graphql/queries/CatalogFilterQuery';
import { iconMedium } from '../../css/icons.styles';
import { modifyRJSFSchema } from '../../utils/utils';
import { getMeshModels } from '../../api/meshmodel';
import _ from 'lodash';
import { useNotification } from '../../utils/hooks/useNotification';
import { EVENT_TYPES } from '../../lib/event-types';
import {
  CustomColumnVisibilityControl,
  ResponsiveDataTable,
  SearchBar,
  UniversalFilter,
  publishCatalogItemSchema,
  publishCatalogItemUiSchema,
  Button,
  PROMPT_VARIANTS,
} from '@sistent/sistent';
import { updateVisibleColumns } from '../../utils/responsive-column';
import { useWindowDimensions } from '../../utils/dimension';
import InfoModal from '../shared/Modal/Information/InfoModal';
import CAN from '@/utils/can';
import { Keys } from '@meshery/schemas/permissions';
import DefaultError from '../general/error-404/index';
import {
  useGetFiltersQuery,
  useCloneFilterMutation,
  usePublishFilterMutation,
  useUnpublishFilterMutation,
  useDeleteFilterMutation,
  useUpdateFilterFileMutation,
  useUploadFilterFileMutation,
} from '@/rtk-query/filter';
import LoadingScreen from '../shared/LoadingState/LoadingComponent';
import { useGetProviderCapabilitiesQuery } from '@/rtk-query/user';
import { isLocalProvider } from '@/utils/provider';
import { ToolWrapper } from '@/assets/styles/general/tool.styles';
import { useSelector } from 'react-redux';
import { updateProgress } from '@/store/slices/mesheryUi';
import { CreateButton, ViewSwitchButton, BtnText } from './Filters.styled';
import YAMLEditor from './YAMLEditor';
import ImportModal from './ImportModal';
import PublishModal from './PublishModal';
import { ACTION_TYPES, COLUMN_VIEWS } from './Filters.constants';
import { buildFiltersColumns } from './Filters.columns';
import { buildFiltersTableOptions } from './Filters.tableOptions';
import {
  createDeleteFilter,
  createHandleClone,
  createHandleDownload,
  createHandleImportFilter,
  createHandlePublish,
  createHandleSubmit,
  createHandleUnpublishModal,
  createUploadHandler,
} from './Filters.fileActions';
import type { TypeView } from './Filters.types';

function resetSelectedFilter() {
  return { show: false, filter: null };
}

function MesheryFilters() {
  const { tableState, updateTableState } = useTableUrlState({
    tableKey: 'fil',
    defaults: {
      page: 0,
      pageSize: 10,
      sortOrder: '',
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

  const visibilityFilter = tableState.filters.vis || null;

  const [count, setCount] = useState(0);
  const modalRef = useRef<{ show: (_args: any) => Promise<string> } | null>(null);
  const [filters, setFilters] = useState<any[]>([]);
  const [selectedFilter, setSelectedFilter] = useState(resetSelectedFilter());
  const [selectedRowData, setSelectedRowData] = useState<any | null>(null);
  const [canPublishFilter, setCanPublishFilter] = useState(false);
  const [publishSchema, setPublishSchema] = useState<{ rjsfSchema?: any; uiSchema?: any }>({});
  const { width } = useWindowDimensions();
  const [meshModels, setMeshModels] = useState<any[]>([]);
  const { user, catalogVisibility } = useSelector((state) => state.ui);
  const [viewType, setViewType] = useState<TypeView>('grid');

  //hooks
  const { notify } = useNotification();
  const [isSearchExpanded, setIsSearchExpanded] = useState(false);

  const [infoModal, setInfoModal] = useState<{
    open: boolean;
    ownerID: string;
    selectedResource: any;
  }>({
    open: false,
    ownerID: '',
    selectedResource: {},
  });

  const [importModal, setImportModal] = useState<{ open: boolean }>({
    open: false,
  });
  const [publishModal, setPublishModal] = useState<{ open: boolean; filter: any; name?: string }>({
    open: false,
    filter: {},
    name: '',
  });

  const catalogContentRef = useRef<any[]>([]);
  const catalogVisibilityRef = useRef<boolean>(false);
  const [selectedFilters, setSelectedFilters] = useState<{ visibility: string }>(() => ({
    visibility: tableState.filters.vis || 'All',
  }));

  const {
    data: filtersData,
    isLoading: isFiltersLoading,
    refetch: getFilters,
  } = useGetFiltersQuery({
    page: page,
    pagesize: pageSize,
    search: search,
    order: sortOrder,
    visibility: visibilityFilter ? JSON.stringify([visibilityFilter]) : '',
  });

  const { data: capabilitiesData } = useGetProviderCapabilitiesQuery();

  const [cloneFilter] = useCloneFilterMutation();
  const [publishFilter] = usePublishFilterMutation();
  const [unpublishFilter] = useUnpublishFilterMutation();
  const [deleteFilterFile] = useDeleteFilterMutation();
  const [updateFilterFile] = useUpdateFilterFileMutation();
  const [uploadFilterFile] = useUploadFilterFileMutation();

  const handleError = (action: { error_msg: string }) => (error: Error | string) => {
    updateProgress({ showProgress: false });
    notify({
      message: `${action.error_msg}: ${error}`,
      event_type: EVENT_TYPES.ERROR,
      details: error.toString(),
    });
  };

  const handleSetFilters = (filters: any[] | undefined) => {
    if (filters != undefined) {
      if (catalogVisibilityRef.current && catalogContentRef.current?.length > 0) {
        setFilters([
          ...catalogContentRef.current,
          ...filters.filter((content) => content.visibility !== VISIBILITY.PUBLISHED),
        ]);
        return;
      }
      setFilters(filters.filter((content) => content.visibility !== VISIBILITY.PUBLISHED));
    }
  };

  const searchTimeout = useRef<ReturnType<typeof setTimeout> | null>(null);

  const handleUploadImport = () => {
    setImportModal({ open: true });
  };

  const handleUploadImportClose = () => {
    setImportModal({ open: false });
  };

  const handlePublishModal = (ev: React.MouseEvent, filter: any) => {
    if (canPublishFilter) {
      ev.stopPropagation();
      setPublishModal({
        open: true,
        filter: filter,
      });
    }
  };

  const handleInfoModalClose = () => {
    setInfoModal({ open: false });
  };

  const handleInfoModal = (filter: any) => {
    setInfoModal({
      open: true,
      ownerID: filter.userId,
      selectedResource: filter,
    });
  };

  const handlePublishModalClose = () => {
    setPublishModal({
      open: false,
      filter: {},
      name: '',
    });
  };

  const handleApplyFilter = () => {
    updateTableState({
      filters: { vis: selectedFilters.visibility === 'All' ? '' : selectedFilters.visibility },
      page: 0,
    });
  };

  function resetSelectedRowData() {
    return () => {
      setSelectedRowData(null);
    };
  }

  async function showmodal(count: number, name?: string) {
    return await modalRef.current?.show({
      title: `Delete ${count ? count : ''} Filter${count > 1 ? 's' : ''}?`,
      subtitle:
        count > 1
          ? `Are you sure you want to delete these ${count} filters?`
          : `Are you sure you want to delete ${name ? `"${name}"` : 'this filter'}?`,
      primaryOption: 'Delete',
      variant: PROMPT_VARIANTS.DANGER,
    });
  }

  const handleUnpublishModal = createHandleUnpublishModal({
    canPublishFilter,
    modalRef,
    unpublishFilter,
    notify,
    handleError,
  });
  const handlePublish = createHandlePublish({
    meshModels,
    publishModal,
    user,
    publishFilter,
    notify,
    handleError,
  });
  const handleClone = createHandleClone({ cloneFilter, notify, handleError });
  const handleDownload = createHandleDownload({ notify });
  const deleteFilter = createDeleteFilter({ deleteFilterFile, notify, handleError });
  const handleSubmit = createHandleSubmit({
    notify,
    handleError,
    showmodal,
    resetSelectedRowData,
    deleteFilterFile,
    uploadFilterFile,
    updateFilterFile,
  });
  const uploadHandler = createUploadHandler({ handleSubmit });
  const handleImportFilter = createHandleImportFilter({
    notify,
    handleError,
    updateFilterFile,
    getFilters,
  });

  useEffect(() => {
    if (filtersData) {
      const filteredWasmFilters = filtersData.filters.filter((content) => {
        if (visibilityFilter === null || content.visibility === visibilityFilter) {
          return true;
        }
        return false;
      });
      setCount(filtersData.totalCount || 0);
      handleSetFilters(filteredWasmFilters);
      setFilters(filtersData.filters || []);
    }
  }, [filtersData]);

  /**
   * Checking whether users are signed in under a provider that doesn't have
   * publish filter capability and setting the canPublishFilter state accordingly
   */
  useEffect(() => {
    const fetchSchemaData = async () => {
      try {
        const { models } = await getMeshModels();
        const modelNames = _.uniq(models?.map((model) => model.displayName));
        modelNames.sort();

        // Modify the schema using the utility function
        const modifiedSchema = modifyRJSFSchema(
          publishCatalogItemSchema,
          'properties.compatibility.items.enum',
          modelNames,
        );
        setPublishSchema({ rjsfSchema: modifiedSchema, uiSchema: publishCatalogItemUiSchema });
        setMeshModels(models);
      } catch {
        handleError(ACTION_TYPES.SCHEMA_FETCH);
      }
    };

    if (capabilitiesData) {
      // `capabilitiesData.capabilities` is the provider-declared feature list.
      const filtersCatalogCapability = capabilitiesData?.capabilities?.filter(
        (val) => val.feature === MesheryFiltersCatalog,
      );
      if (filtersCatalogCapability?.length > 0) setCanPublishFilter(true);
    }
    fetchSchemaData();
  }, [capabilitiesData]);

  useEffect(() => {
    if (viewType === 'grid') setSearch('');
  }, [viewType]);

  useEffect(() => {
    catalogVisibilityRef.current = catalogVisibility;
    const fetchCatalogFilters = fetchCatalogFilter({
      selector: {
        search: '',
        order: '',
        page: 0,
        pagesize: 0,
      },
    }).subscribe({
      next: (result) => {
        catalogContentRef.current = result?.catalogFilters;
      },
      error: (err) => console.log('There was an error fetching Catalog Filter: ', err),
    });

    return () => {
      fetchCatalogFilters.unsubscribe();
    };
  }, []);

  const columns = buildFiltersColumns({
    filters,
    canPublishFilter,
    handleClone,
    handleDownload,
    handleInfoModal,
    handlePublishModal,
    handleUnpublishModal,
    setSelectedRowData,
    sortOrder,
  });

  const options = buildFiltersTableOptions({
    isLocalProvider: isLocalProvider(capabilitiesData),
    count,
    page,
    pageSize,
    search,
    sortOrder,
    filters,
    columns,
    searchTimeout,
    setPage,
    setPageSize,
    setSearch,
    setSortOrder,
    setSelectedRowData,
    showmodal,
    deleteFilter,
  });

  const [tableCols, updateCols] = useState(columns);

  const responsiveColDefaults = (() => {
    const showCols = updateVisibleColumns(COLUMN_VIEWS, width);
    const initialVisibility: Record<string, boolean> = {};
    columns.forEach((col) => {
      initialVisibility[col.name] = showCols[col.name];
    });
    return initialVisibility;
  })();

  const {
    columnVisibility,
    setColumnVisibilityByUser: setColumnVisibility,
    setColumnVisibilityByResponsive,
  } = useColumnVisibilityPreference('filters', responsiveColDefaults);

  useEffect(() => {
    setColumnVisibilityByResponsive(responsiveColDefaults);
  }, [width, setColumnVisibilityByResponsive]);

  const filter = {
    visibility: {
      name: 'visibility',
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

  if (isFiltersLoading) {
    return (
      <>
        <LoadingScreen animatedIcon="AnimatedFilter" message={`Loading Filters...`} />
      </>
    );
  }

  return (
    <>
      <>
        <NoSsr>
          {CAN(Keys.CatalogManagementViewFilters.id, Keys.CatalogManagementViewFilters.function) ? (
            <>
              {selectedRowData && Object.keys(selectedRowData).length > 0 && (
                <YAMLEditor
                  filter={selectedRowData}
                  onClose={resetSelectedRowData()}
                  onSubmit={handleSubmit}
                />
              )}
              <ToolWrapper>
                {width < 600 && isSearchExpanded ? null : (
                  <div style={{ display: 'flex' }}>
                    {!selectedFilter.show && (filters.length > 0 || viewType === 'table') && (
                      <CreateButton>
                        <Button
                          aria-label="Add Filter"
                          variant="contained"
                          color="primary"
                          size="large"
                          onClick={handleUploadImport}
                          permissionKey={Keys.CatalogManagementImportFilter}
                        >
                          <PublishIcon style={iconMedium} data-cy="import-button" />
                          <BtnText> Import Filters </BtnText>
                        </Button>
                      </CreateButton>
                    )}
                    <ViewSwitchButton style={{ jdisplay: 'flex' }}>
                      {/* <CatalogFilter
                      catalogVisibility={catalogVisibility}
                      handleCatalogVisibility={handleCatalogVisibility}
                      classes={classes}
                    /> */}
                    </ViewSwitchButton>
                  </div>
                )}
                <div style={{ display: 'flex', alignItems: 'center' }}>
                  <SearchBar
                    onSearch={(value) => {
                      setSearch(value);
                    }}
                    expanded={isSearchExpanded}
                    setExpanded={setIsSearchExpanded}
                    placeholder="Search"
                  />
                  <UniversalFilter
                    id="ref"
                    filters={filter}
                    selectedFilters={selectedFilters}
                    setSelectedFilters={setSelectedFilters}
                    handleApplyFilter={handleApplyFilter}
                  />
                  {viewType === 'table' && (
                    <CustomColumnVisibilityControl
                      id="ref"
                      columns={columns}
                      customToolsProps={{ columnVisibility, setColumnVisibility }}
                    />
                  )}

                  {!selectedFilter.show && (
                    <ViewSwitch data-cy="table-view" view={viewType} changeView={setViewType} />
                  )}
                </div>
              </ToolWrapper>
              {!selectedFilter.show && viewType === 'table' && (
                <ResponsiveDataTable
                  data={filters}
                  columns={columns}
                  tableCols={tableCols}
                  updateCols={updateCols}
                  columnVisibility={columnVisibility}
                  options={options}
                />
              )}
              {!selectedFilter.show && viewType === 'grid' && (
                // grid view
                <FiltersGrid
                  filters={filters}
                  handleSubmit={handleSubmit}
                  canPublishFilter={canPublishFilter}
                  handlePublish={handlePublish}
                  handleUnpublishModal={handleUnpublishModal}
                  handleUploadImport={handleUploadImport}
                  handleClone={handleClone}
                  handleDownload={handleDownload}
                  uploadHandler={uploadHandler}
                  setSelectedFilter={setSelectedFilter}
                  selectedFilter={selectedFilter}
                  pages={Math.ceil(count / pageSize)}
                  setPage={setPage}
                  selectedPage={page}
                  publishModal={publishModal}
                  setPublishModal={setPublishModal}
                  publishSchema={publishSchema}
                  fetch={() => getFilters()}
                  handleInfoModal={handleInfoModal}
                />
              )}
              {canPublishFilter &&
                publishModal.open &&
                CAN(
                  Keys.CatalogManagementPublishWasmFilter.id,
                  Keys.CatalogManagementPublishWasmFilter.function,
                ) && (
                  <PublishModal
                    handleClose={handlePublishModalClose}
                    title={publishModal.filter?.name}
                    handleSubmit={handlePublish}
                  />
                )}
              {importModal.open &&
                CAN(
                  Keys.CatalogManagementImportFilter.id,
                  Keys.CatalogManagementImportFilter.function,
                ) && (
                  <ImportModal
                    handleClose={handleUploadImportClose}
                    handleImportFilter={handleImportFilter}
                  />
                )}
              {infoModal.open &&
                CAN(
                  Keys.CatalogManagementDetailsOfWasmFilter.id,
                  Keys.CatalogManagementDetailsOfWasmFilter.function,
                ) && (
                  <InfoModal
                    handlePublish={handlePublish}
                    infoModalOpen={true}
                    handleInfoModalClose={handleInfoModalClose}
                    dataName="filters"
                    selectedResource={infoModal.selectedResource}
                    resourceOwnerID={infoModal.ownerID}
                    currentUser={user}
                    formSchema={publishSchema}
                    meshModels={meshModels}
                    patternFetcher={getFilters}
                  />
                )}
              <_PromptComponent ref={modalRef} />
            </>
          ) : (
            <DefaultError />
          )}
        </NoSsr>
      </>
    </>
  );
}

export default MesheryFilters;
