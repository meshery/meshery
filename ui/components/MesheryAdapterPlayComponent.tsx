import React, { useState, useRef, useEffect } from 'react';
import { Grid, NoSsr } from '@sistent/sistent';
import { useRouter } from 'next/router';
import * as yaml from 'js-yaml';
import { ctxUrl, getK8sClusterIdsFromCtxId } from '../utils/multi-ctx';
import {
  useAdapterOperationMutation,
  useLazyPingAdapterQuery,
  useLazyGetSmiResultsQuery,
} from '../rtk-query/system';
import fetchAvailableAddons from '@/graphql/queries/AddonsStatusQuery';
import fetchAvailableNamespaces from '@/graphql/queries/NamespaceQuery';
import MesheryResultDialog from './MesheryResultDialog';
import ReactSelectWrapper from './ReactSelectWrapper';
import ConfirmationMsg from '@/components/designs/lifecycle/DeployConfirmationModal';
import { ACTIONS } from '../utils/Enum';
import { getModelByName } from '../api/meshmodel';
import { EVENT_TYPES } from '../lib/event-types';
import { useNotification } from '../utils/hooks/useNotification';
import { useSelector } from 'react-redux';
import { updateProgress } from '@/store/slices/mesheryUi';
import {
  AdapterChip,
  AdapterSmWrapper,
  ChipNamespaceContainer,
  InputWrapper,
  PaneSection,
} from './adapter-play-styled';
import AdapterYamlDialog from './adapter-play-yaml-dialog';
import AdapterSmiResultsDialog from './adapter-play-smi-results';
import AdapterCategoryCard from './adapter-play-category-card';
import AdapterAddonSwitches from './adapter-play-addon-switches';

// Re-export AdapterChip for backwards compatibility — the previous
// public surface of this module included this styled component.
export { AdapterChip } from './adapter-play-styled';

interface AdapterOperation {
  key: string;
  value: string;
  category?: number;
}

interface Adapter {
  name: string;
  adapterLocation: string;
  ops?: AdapterOperation[];
}

interface MesheryAdapterPlayComponentProps {
  adapter: Adapter;
  // Optional consumer-supplied props that the original component reads
  // via destructuring without declaring; kept for parity.
  user?: { userId?: string };
}

const MesheryAdapterPlayComponent: React.FC<MesheryAdapterPlayComponentProps> = (props) => {
  const { k8sConfig, selectedK8sContexts } = useSelector((state) => state.ui);
  const { adapter } = props;
  const { notify } = useNotification();
  const [triggerAdapterOp] = useAdapterOperationMutation();
  const [triggerPingAdapter] = useLazyPingAdapterQuery();
  const [triggerGetSmiResults] = useLazyGetSmiResultsQuery();
  const router = useRouter();
  const addIconEles = useRef({});
  const delIconEles = useRef({});

  const [cmEditorValAdd, setCmEditorValAdd] = useState('');
  const [, setCmEditorValAddError] = useState(false);
  const [cmEditorValDel, setCmEditorValDel] = useState('');
  const [, setCmEditorValDelError] = useState(false);
  const [, setSelectionError] = useState(false);
  const [namespace, setNamespace] = useState({ value: 'default', label: 'default' });
  const [namespaceError, setNamespaceError] = useState(false);
  const [customDialogAdd, setCustomDialogAdd] = useState(false);
  const [customDialogDel, setCustomDialogDel] = useState(false);
  const [customDialogSMI, setCustomDialogSMI] = useState(false);
  const [modalOpen, setModalOpen] = useState(false);
  const [addonSwitchGroup, setAddonSwitchGroup] = useState({});
  const [smi_result, setSmiResult] = useState([]);
  const [selectedRowData, setSelectedRowData] = useState(null);
  const [page, setPage] = useState(0);
  const [search, setSearch] = useState('');
  const [sortOrder, setSortOrder] = useState('');
  const [pageSize, setPageSize] = useState(10);
  const [namespaceList, setNamespaceList] = useState([]);
  const [namespaceSubscription, setNamespaceSubscription] = useState(null);
  const [, setActiveContexts] = useState([]);
  // const [deployModalOpen, setDeployModalOpen] = useState(false);
  const [category, setCategory] = useState(0);
  const [selectedOp, setSelectedOp] = useState('');
  const [isDeleteOp, setIsDeleteOp] = useState(false);
  const [operationName, setOperationName] = useState('');
  const [versionList, setVersionList] = useState([]);
  const [version, setVersion] = useState({ label: '', value: '' });
  const [versionError, setVersionError] = useState(false);

  const isYamlValid = (value) => {
    if (!value) {
      return false;
    }

    try {
      yaml.load(value);
      return true;
    } catch (error) {
      return false;
    }
  };

  // Initialize menu state based on adapter operations
  const initMenuState = () => {
    const menuState = {};
    if (adapter && adapter.ops) {
      [0, 1, 2, 3, 4].forEach((i) => {
        menuState[i] = {
          add: false,
          delete: false,
        };
      });
    }
    return menuState;
  };

  const [menuState, setMenuState] = useState(initMenuState());

  // Equivalent to componentDidMount
  useEffect(() => {
    const meshname = mapAdapterNameToMeshName(adapter.name);
    const variables = { type: meshname, k8sClusterIDs: getK8sClusterIds() };

    initSubscription();
    getMeshVersions();

    if (selectedK8sContexts) {
      if (selectedK8sContexts.includes('all')) {
        let active = [];
        k8sConfig.forEach((ctx) => {
          active.push(ctx.contextID);
        });
        setActiveContexts(active);
      } else {
        setActiveContexts(selectedK8sContexts);
      }
    }

    fetchAvailableAddons(variables).subscribe({
      next: (res) => {
        setAddonsState(res);
      },
      error: (err) => console.log('error at addon fetch: ' + err),
    });

    // Cleanup function (componentWillUnmount)
    return () => {
      disposeSubscriptions();
    };
  }, [adapter.name]);

  // Equivalent to componentDidUpdate for selectedK8sContexts
  useEffect(() => {
    if (selectedK8sContexts) {
      disposeSubscriptions();
      initSubscription();
    }
  }, [selectedK8sContexts]);

  // Equivalent to componentDidUpdate for adapter.name
  useEffect(() => {
    getMeshVersions();
  }, [props.adapter.name]);

  const initSubscription = () => {
    disposeSubscriptions();

    const subscription = fetchAvailableNamespaces({
      k8sClusterIDs: getK8sClusterIds(),
    }).subscribe({
      next: (res) => {
        let namespaces = [];
        res?.namespaces?.map((ns) => {
          namespaces.push({
            value: ns?.namespace,
            label: ns?.namespace,
          });
        });
        if (namespaces.length === 0) {
          namespaces.push({
            value: 'default',
            label: 'default',
          });
        }
        namespaces.sort((a, b) => (a.value > b.value ? 1 : -1));
        setNamespaceList(namespaces);
      },
      error: (err) => console.log('error at namespace fetch: ' + err),
    });

    setNamespaceSubscription(subscription);
  };

  const disposeSubscriptions = () => {
    if (namespaceSubscription) {
      namespaceSubscription.unsubscribe();
    }
  };

  const getK8sClusterIds = () => {
    return getK8sClusterIdsFromCtxId(selectedK8sContexts, k8sConfig);
  };

  const mapAdapterNameToMeshName = (name) => {
    if (name?.toLowerCase() === 'istio') return 'ISTIO';
    return 'ALL';
  };

  const getMeshVersions = () => {
    const activeMesh = adapter?.name;
    getModelByName(activeMesh.toLowerCase()).then((res) => {
      let uniqueVersions = [...new Set(res?.models?.map((model) => model?.version))].reverse();
      if (uniqueVersions.length === 0) {
        uniqueVersions = [''];
      }
      let versionList = uniqueVersions.map((version) => ({ value: version, label: version }));
      setVersionList(versionList);
      setVersion(versionList[0]);
    });
  };

  const setAddonsState = (data) => {
    const meshname = adapter.name;
    const localState = {};
    data?.addonsState?.forEach((addon) => {
      if (addon.owner === meshname) {
        const name = addon.name !== 'jaeger-collector' ? addon.name : 'jaeger';
        localState[`${name}-addon`] = true;
      }
    });
    setAddonSwitchGroup(localState);
  };

  const handleNamespaceChange = (newValue) => {
    if (typeof newValue !== 'undefined') {
      setNamespace(newValue);
      setNamespaceError(false);
    } else {
      setNamespaceError(true);
    }
  };

  const handleVersionChange = (newValue) => {
    if (typeof newValue !== 'undefined') {
      setVersion(newValue);
      setVersionError(false);
    } else {
      setVersionError(true);
    }
  };

  // const handleChange =
  //   (name, isDelete = false) =>
  //   (event) => {
  //     if (name === 'selectedOp' && event.target.value !== '') {
  //       if (event.target.value === 'custom') {
  //         if (isDelete) {
  //           if (cmEditorValDel !== '' && cmEditorDelRef.current?.state.lint.marked.length === 0) {
  //             setSelectionError(false);
  //             setCmEditorValDelError(false);
  //           }
  //         } else if (
  //           cmEditorValAdd !== '' &&
  //           cmEditorAddRef.current?.state.lint.marked.length === 0
  //         ) {
  //           setSelectionError(false);
  //           setCmEditorValAddError(false);
  //         }
  //       } else {
  //         setSelectionError(false);
  //       }
  //     }

  //     if (name === 'selectedOp') {
  //       setSelectedOp(event.target.value);
  //     }
  //   };

  const handleModalClose = (isDelete) => () => {
    if (isDelete) {
      setCustomDialogDel(false);
    } else {
      setCustomDialogAdd(false);
    }
  };

  const handleSMIClose = () => {
    setCustomDialogSMI(false);
  };

  const resetSelectedRowData = () => {
    setSelectedRowData(null);
  };

  // const handleModalOpen = (isDelete) => () => {
  //   if (isDelete) {
  //     setCustomDialogDel(true);
  //   } else {
  //     setCustomDialogAdd(true);
  //   }
  // };

  const handleOpen = () => {
    setModalOpen(true);
  };

  const handleClose = () => {
    setModalOpen(false);
  };

  // const handleDeployModalOpen = () => {
  //   setDeployModalOpen(true);
  // };

  // const handleDeployModalClose = () => {
  //   setDeployModalOpen(false);
  // };

  const handleSubmit =
    (cat, op, deleteOp = false) =>
    () => {
      handleOpen();

      const filteredOp = adapter.ops.filter(({ key }) => key === op);
      if (op === '' || typeof filteredOp === 'undefined' || filteredOp.length === 0) {
        setSelectionError(true);
        return;
      }

      if (deleteOp) {
        if (op === 'custom' && !isYamlValid(cmEditorValDel)) {
          setCmEditorValDelError(true);
          setSelectionError(true);
          return;
        }
      } else if (op === 'custom' && !isYamlValid(cmEditorValAdd)) {
        setCmEditorValAddError(true);
        setSelectionError(true);
        return;
      }

      if (namespace.value === '') {
        setNamespaceError(true);
        return;
      }

      if (version?.value === '') {
        setVersionError(true);
        return;
      }

      const opName = op
        .replaceAll('_', ' ')
        .split(' ')
        .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
        .join(' ');

      setCategory(cat);
      setSelectedOp(op);
      setIsDeleteOp(deleteOp);
      setOperationName(opName);
    };

  const submitOp = (cat, op, deleteOp = false) => {
    const data = {
      adapter: adapter.adapterLocation,
      query: op,
      namespace: namespace.value,
      customBody: deleteOp ? cmEditorValDel : cmEditorValAdd,
      deleteOp: deleteOp ? 'on' : '',
      version: version.value,
    };

    const params = Object.keys(data)
      .map((key) => `${encodeURIComponent(key)}=${encodeURIComponent(data[key])}`)
      .join('&');

    updateProgress({ showProgress: true });
    handleClose();

    triggerAdapterOp({
      url: ctxUrl('system/adapter/operation', selectedK8sContexts),
      body: params,
    })
      .unwrap()
      .then((result) => {
        updateProgress({ showProgress: false });

        const newMenuState = { ...menuState };
        newMenuState[cat][deleteOp ? 'delete' : 'add'] = false;
        setMenuState(newMenuState);

        if (deleteOp) {
          setCustomDialogDel(false);
        } else {
          setCustomDialogAdd(false);
        }

        if (typeof result !== 'undefined') {
          notify({ message: 'Operation executing...', event_type: EVENT_TYPES.INFO });
        }
      })
      .catch(handleError(cat, deleteOp, op));
  };

  const handleAdapterClick = (adapterLoc) => async () => {
    updateProgress({ showProgress: true });
    try {
      await triggerPingAdapter(adapterLoc).unwrap();
      updateProgress({ showProgress: false });
      notify({ message: 'Adapter pinged!', event_type: EVENT_TYPES.SUCCESS });
    } catch (err) {
      updateProgress({ showProgress: false });
      handleError('Could not ping adapter.')(err);
    }
  };

  const fetchSMIResults = async (adapterName, page, pageSize, search, sortOrder) => {
    if (typeof search === 'undefined' || search === null) {
      search = '';
    }
    if (typeof sortOrder === 'undefined' || sortOrder === null) {
      sortOrder = '';
    }

    try {
      const result = await triggerGetSmiResults({
        page,
        pagesize: pageSize,
        search,
        order: sortOrder,
      }).unwrap();

      if (typeof result !== 'undefined' && result.results) {
        const results = result.results.filter(
          (val) => val.meshName.toLowerCase() == adapterName.toLowerCase(),
        );
        setSmiResult({ ...result, results: results, totalCount: results.length });
      }
    } catch (error) {
      console.log('Could not fetch SMI results.', error);
    }
  };

  // const handleSMIClick = (adapterName) => () => {
  //   updateProgress({ showProgress: true });
  //   fetchSMIResults(adapterName, page, pageSize, search, sortOrder);
  //   updateProgress({ showProgress: false });
  //   setCustomDialogSMI(true);
  // };

  const handleError = (cat, deleteOp, selectedOp) => {
    return (error) => {
      if (cat && deleteOp !== undefined) {
        const newMenuState = { ...menuState };
        newMenuState[cat][deleteOp ? 'delete' : 'add'] = false;
        setMenuState(newMenuState);

        if (deleteOp) {
          setCustomDialogDel(false);
        } else {
          setCustomDialogAdd(false);
        }
      }

      if (selectedOp) {
        setAddonSwitchGroup({ ...addonSwitchGroup, [selectedOp]: deleteOp });
      }

      updateProgress({ showProgress: false });
      notify({
        message: `Operation submission failed: ${error}`,
        event_type: EVENT_TYPES.ERROR,
        details: error.toString(),
      });
    };
  };

  const addDelHandleClick = (cat: number, isDelete: boolean) => {
    const newMenuState = { ...menuState };
    newMenuState[cat][isDelete ? 'delete' : 'add'] =
      !newMenuState[cat][isDelete ? 'delete' : 'add'];
    setMenuState(newMenuState);

    if (cat === 4) {
      if (isDelete) {
        setCustomDialogDel(!customDialogDel);
      } else {
        setCustomDialogAdd(!customDialogAdd);
      }
    }
  };

  const renderYamlDialog = (cat, isDelete) => {
    const value = isDelete ? cmEditorValDel : cmEditorValAdd;
    return (
      <AdapterYamlDialog
        open={isDelete ? customDialogDel : customDialogAdd}
        isDelete={isDelete}
        adapterName={adapter.name}
        namespace={namespace}
        namespaceError={namespaceError}
        namespaceList={namespaceList}
        onNamespaceChange={handleNamespaceChange}
        version={version}
        versionError={versionError}
        versionList={versionList}
        onVersionChange={handleVersionChange}
        value={value}
        onBeforeChange={(_editor, _data, newValue) => {
          const isValid = isYamlValid(newValue);
          if (isDelete) {
            setCmEditorValDel(newValue);
            if (isValid) {
              setSelectionError(false);
              setCmEditorValDelError(false);
            }
          } else {
            setCmEditorValAdd(newValue);
            if (isValid) {
              setSelectionError(false);
              setCmEditorValAddError(false);
            }
          }
        }}
        onClose={handleModalClose(isDelete)}
        onApply={handleSubmit(cat, 'custom', isDelete)}
      />
    );
  };

  const extractAddonOperations = (addonOpsCat) => {
    return adapter.ops.filter(
      ({ category, value }) => category === addonOpsCat && value?.startsWith('Add-on:'),
    );
  };

  const handleAddonSwitchChange = (name, checked, ops) => {
    setAddonSwitchGroup({ ...addonSwitchGroup, [name]: checked });
    submitOp(ops.category, ops.key, !addonSwitchGroup[ops.key]);
  };

  // Render component
  let adapterName = adapter.name.split(' ').join('').toLowerCase();
  let imageSrc = '/static/img/' + adapterName + '.svg';
  let adapterChip = (
    <AdapterChip
      label={adapter.adapterLocation}
      data-cy="adapter-chip-ping"
      onClick={handleAdapterClick(adapter.adapterLocation)}
      icon={<img src={imageSrc} width={'1.25rem'} />}
      variant="outlined"
    />
  );

  const filteredOps = [];
  if (adapter && adapter.ops && adapter.ops.length > 0) {
    adapter.ops.forEach(({ category }) => {
      if (typeof category === 'undefined') {
        category = 0;
      }
      if (filteredOps.indexOf(category) === -1) {
        filteredOps.push(category);
      }
    });
    filteredOps.sort();
  }

  return (
    <NoSsr>
      {selectedRowData && selectedRowData !== null && Object.keys(selectedRowData).length > 0 && (
        <MesheryResultDialog rowData={selectedRowData} close={resetSelectedRowData} />
      )}
      <React.Fragment>
        <AdapterSmWrapper>
          <Grid container spacing={2} direction="row" alignItems="flex-start">
            {/* SECTION 1 */}
            <Grid item xs={12}>
              <PaneSection>
                <Grid container spacing={4}>
                  <ChipNamespaceContainer
                    container
                    item
                    xs={12}
                    alignItems="flex-start"
                    justify="space-between"
                  >
                    <div>{adapterChip}</div>
                    <InputWrapper>
                      <ReactSelectWrapper
                        label="Namespace"
                        value={namespace}
                        error={namespaceError}
                        options={namespaceList}
                        onChange={handleNamespaceChange}
                      />
                    </InputWrapper>
                    <InputWrapper>
                      <ReactSelectWrapper
                        label="Version"
                        value={version}
                        error={versionError}
                        options={versionList}
                        onChange={handleVersionChange}
                      />
                    </InputWrapper>
                  </ChipNamespaceContainer>
                  <Grid container spacing={1} style={{ margin: '1rem' }}>
                    <Grid
                      container
                      item
                      lg={!extractAddonOperations(2).length ? 12 : 10}
                      xs={12}
                      spacing={2}
                    >
                      {filteredOps.map((val, i) => (
                        <Grid item lg={3} md={4} xs={12} key={`adapter-card-${i}`}>
                          <AdapterCategoryCard
                            cat={typeof val === 'undefined' ? 0 : val}
                            adapterOps={adapter?.ops ?? []}
                            menuState={menuState}
                            addIconEles={addIconEles}
                            delIconEles={delIconEles}
                            onMenuToggle={addDelHandleClick}
                            onMenuItemClick={handleSubmit}
                            renderYamlDialog={renderYamlDialog}
                          />
                        </Grid>
                      ))}
                    </Grid>
                    <Grid container item lg={2} xs={12}>
                      <Grid item xs={12} md={4}>
                        <AdapterAddonSwitches
                          selectedAdapterOps={extractAddonOperations(2)}
                          addonSwitchGroup={addonSwitchGroup}
                          onSwitchChange={handleAddonSwitchChange}
                        />
                      </Grid>
                    </Grid>
                  </Grid>
                </Grid>
              </PaneSection>
            </Grid>
          </Grid>
        </AdapterSmWrapper>
        <ConfirmationMsg
          open={modalOpen}
          handleClose={handleClose}
          submit={{
            deploy: () => submitOp(category, selectedOp, false),
            unDeploy: () => submitOp(category, selectedOp, true),
          }}
          isDelete={isDeleteOp}
          title={operationName}
          tab={isDeleteOp ? ACTIONS.UNDEPLOY : ACTIONS.DEPLOY}
        />
        {customDialogSMI && (
          <AdapterSmiResultsDialog
            open={customDialogSMI}
            onClose={handleSMIClose}
            adapterName={adapter.name}
            smiResult={smi_result}
            page={page}
            pageSize={pageSize}
            search={search}
            sortOrder={sortOrder}
            fetchSMIResults={fetchSMIResults}
            onPageChange={setPage}
            onPageSizeChange={setPageSize}
            onSearchChange={setSearch}
            onSortOrderChange={setSortOrder}
          />
        )}
      </React.Fragment>
    </NoSsr>
  );
};

export default MesheryAdapterPlayComponent;
