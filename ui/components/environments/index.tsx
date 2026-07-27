import { useEffect, useRef, useState } from 'react';
import {
  ChevronLeftIcon,
  ChevronRightIcon,
  DeleteIcon,
  NoSsr,
  Pagination,
  PaginationItem,
  useHasPermission,
  useTheme,
} from '@sistent/sistent';
import { withRouter } from 'next/router';
import { debounce } from 'lodash';
import { CreateButtonWrapper, BulkActionWrapper } from './styles';
import { ToolWrapper } from '@/assets/styles/general/tool.styles';
import AddIconCircleBorder from '../../assets/icons/AddIconCircleBorder';
import EnvironmentCard from './environment-card';
import EnvironmentIcon from '../../assets/icons/Environment';
import { EVENT_TYPES } from '../../lib/event-types';
import { useNotification } from '../../utils/hooks/useNotification';
import { formatApiError } from '../../utils/helpers/meshkitError';
import { RJSFModalWrapper } from '../shared/Modal/Modal';
import _PromptComponent from '../general/PromptComponent';
import { EmptyState } from '../lifecycle/general';
import {
  Modal as SisitentModal,
  ModalBody,
  TransferList,
  ModalFooter,
  PrimaryActionButtons,
  createAndEditEnvironmentSchema,
  createAndEditEnvironmentUiSchema,
  ErrorBoundary,
  Button,
  Grid2,
  Typography,
  SearchBar,
  PROMPT_VARIANTS,
} from '@sistent/sistent';
import ConnectionIcon from '../../assets/icons/Connection';
import { TRANSFER_COMPONENT } from '../../utils/Enum';
import {
  useAddConnectionToEnvironmentMutation,
  useRemoveConnectionFromEnvironmentMutation,
  useGetEnvironmentConnectionsQuery,
  useGetEnvironmentsQuery,
  useCreateEnvironmentMutation,
  useUpdateEnvironmentMutation,
  useDeleteEnvironmentMutation,
} from '../../rtk-query/environments';
import { Keys } from '@meshery/schemas/permissions';
import DefaultError from '../general/error-404/index';
import { useSelector } from 'react-redux';
import { updateProgress } from '@/store/slices/mesheryUi';

const ACTION_TYPES = {
  CREATE: 'create',
  EDIT: 'edit',
};

const Environments = () => {
  const canViewEnv = useHasPermission(Keys.WorkspaceManagementViewEnvironment);
  const canCreateEnv = useHasPermission(Keys.WorkspaceManagementCreateEnvironment);
  const canEditEnv = useHasPermission(Keys.WorkspaceManagementEditEnvironment);
  const canAssignConnToEnv = useHasPermission(
    Keys.WorkspaceManagementAssignConnectionsToEnvironment,
  );
  const canRemoveConnFromEnv = useHasPermission(
    Keys.WorkspaceManagementRemoveConnectionsFromEnvironments,
  );
  const theme = useTheme();
  const { organization } = useSelector((state) => state.ui);
  const [environmentModal, setEnvironmentModal] = useState({
    open: false,
    schema: {},
  });
  const [actionType, setActionType] = useState('');
  const [initialData, setInitialData] = useState({});
  const [editEnvId, setEditEnvId] = useState('');
  const [page, setPage] = useState(0);
  const [search, setSearch] = useState('');
  const [orgId, setOrgId] = useState('');
  const [selectedEnvironments, setSelectedEnvironments] = useState([]);
  const [assignConnectionModal, setAssignConnectionModal] = useState(false);
  const [connectionAssignEnv, setConnectionAssignEnv] = useState({});
  const [assignedConnections, setAssignedConnections] = useState([]);
  const [connectionsData, setConnectionsData] = useState([]);
  const [connectionsPage, setConnectionsPage] = useState(0);
  const [environmentConnectionsData, setEnvironmentConnectionsData] = useState([]);
  const [connectionsOfEnvironmentPage, setConnectionsOfEnvironmentPage] = useState(0);
  const [skip, setSkip] = useState(true);
  const [isSearchExpanded, setIsSearchExpanded] = useState(false);
  const [disableTranferButton, setDisableTranferButton] = useState(true);

  const pageSize = 10;
  const connectionPageSize = 25;

  const modalRef = useRef(null);
  const { notify } = useNotification();

  const {
    data: environmentsData,
    isError: isEnvironmentsError,
    error: environmentsError,
  } = useGetEnvironmentsQuery(
    {
      search: search,
      page: page,
      pagesize: pageSize,
      orgId: orgId,
    },
    {
      skip: !orgId ? true : false,
    },
  );

  const [createEnvironment] = useCreateEnvironmentMutation();

  const [updateEnvironment] = useUpdateEnvironmentMutation();

  const [deleteEnvironment] = useDeleteEnvironmentMutation();

  const {
    data: connections,
    isError: isConnectionsError,
    error: connectionsError,
  } = useGetEnvironmentConnectionsQuery(
    {
      environmentId: connectionAssignEnv.id,
      page: connectionsData.length === 0 ? 0 : connectionsPage,
      pagesize: connectionPageSize,
      filter: '{"assigned":false}',
    },
    {
      skip,
    },
  );

  const {
    data: environmentConnections,
    isError: isEnvironmentConnectionsError,
    error: environmentConnectionsError,
  } = useGetEnvironmentConnectionsQuery(
    {
      environmentId: connectionAssignEnv.id,
      page: environmentConnectionsData.length === 0 ? 0 : connectionsOfEnvironmentPage,
      pagesize: connectionPageSize,
    },
    {
      skip,
    },
  );

  const environments = environmentsData?.environments ? environmentsData.environments : [];
  const connectionsDataRtk = connections?.connections ? connections.connections : [];
  const environmentConnectionsDataRtk = environmentConnections?.connections
    ? environmentConnections.connections
    : [];

  useEffect(() => {
    setConnectionsData((prevData) => [...prevData, ...connectionsDataRtk]);
  }, [connections]);

  useEffect(() => {
    setEnvironmentConnectionsData((prevData) => [...prevData, ...environmentConnectionsDataRtk]);
  }, [environmentConnections]);

  useEffect(() => {
    if (isEnvironmentsError) {
      handleError({ error_msg: 'Unable to fetch environments' })(environmentsError);
    }
    if (isEnvironmentConnectionsError) {
      handleError({ error_msg: "Unable to fetch an environment's connections" })(
        environmentConnectionsError,
      );
    }
    if (isConnectionsError) {
      handleError({ error_msg: 'Unable to fetch connections' })(connectionsError);
    }
  }, [
    isEnvironmentsError,
    isEnvironmentConnectionsError,
    isConnectionsError,
    environmentsError,
    environmentConnectionsError,
    connectionsError,
  ]);

  // Curried so a call site can name the failing operation once and hand the
  // raw error to the returned callback: `handleError({ error_msg })(error)`.
  // `formatApiError` consumes the MeshKit envelope the server now sends (code
  // and suggested remediation) and renders it as the markdown that `notify`
  // displays through BasicMarkdown; `error_msg` is the fallback title used only
  // when the response carries no envelope. The raw error must be passed
  // through - do not pre-flatten it to a string, or the envelope is lost.
  function handleError(action) {
    return (error) => {
      updateProgress({ showProgress: false });
      const { message } = formatApiError(error, action?.error_msg);
      notify({
        message,
        event_type: EVENT_TYPES.ERROR,
      });
    };
  }

  const handleSuccess = (msg) => {
    updateProgress({ showProgress: false });
    notify({
      message: msg,
      event_type: EVENT_TYPES.SUCCESS,
    });
  };

  useEffect(() => {
    setOrgId(organization?.id);
  }, [organization]);

  const fetchSchema = () => {
    // Organization is derived from the user's active session and hidden by the
    // canonical form UI schema (organizationId -> "ui:widget": "hidden" in
    // meshery/schemas environment/forms/createOrEdit.ui.json). Its value is
    // seeded into the form via initialData, so no per-render schema patching is
    // required here.
    setEnvironmentModal({
      open: true,
      schema: {
        schema: createAndEditEnvironmentSchema,
        uischema: createAndEditEnvironmentUiSchema,
      },
    });
  };

  const [addConnectionToEnvironmentMutator] = useAddConnectionToEnvironmentMutation();
  const [removeConnectionFromEnvMutator] = useRemoveConnectionFromEnvironmentMutation();

  // Both mutators previously fired and discarded the returned promise, so a
  // rejected assignment left the transfer list looking like it had succeeded.
  const addConnectionToEnvironment = (environmentId, connectionId) =>
    addConnectionToEnvironmentMutator({ environmentId, connectionId })
      .unwrap()
      .catch(handleError({ error_msg: 'Unable to assign connection to environment' }));

  const removeConnectionFromEnvironment = (environmentId, connectionId) =>
    removeConnectionFromEnvMutator({ environmentId, connectionId })
      .unwrap()
      .catch(handleError({ error_msg: 'Unable to remove connection from environment' }));

  const handleEnvironmentModalOpen = (e, actionType, envObject) => {
    e.stopPropagation();
    if (actionType === ACTION_TYPES.EDIT) {
      setActionType(ACTION_TYPES.EDIT);
      setInitialData({
        name: envObject.name,
        description: envObject.description,
        organizationId: envObject.organizationId,
      });
      setEditEnvId(envObject.id);
    } else {
      setActionType(ACTION_TYPES.CREATE);
      setInitialData({
        name: undefined,
        description: '',
        organizationId: orgId,
      });
      setEditEnvId('');
    }
    fetchSchema();
  };

  const handleEnvironmentModalClose = () => {
    setEnvironmentModal({
      open: false,
      schema: {},
    });
    setActionType('');
  };

  const handleCreateEnvironment = ({ organizationId, name, description }) => {
    createEnvironment({
      environmentPayload: {
        name: name,
        description: description,
        organizationId: organizationId,
      },
    })
      .unwrap()
      .then(() => {
        handleSuccess(`Environment "${name}" created`);
        handleEnvironmentModalClose();
      })
      .catch(handleError({ error_msg: `Unable to create environment "${name}"` }));
  };

  const handleEditEnvironment = ({ name, description }) => {
    updateEnvironment({
      environmentId: editEnvId,
      environmentPayload: {
        name: name,
        description: description,
        organizationId: initialData.organizationId,
      },
    })
      .unwrap()
      .then(() => {
        handleSuccess(`Environment "${name}" updated`);
        handleEnvironmentModalClose();
      })
      .catch(handleError({ error_msg: `Unable to update environment "${name}"` }));
  };

  const handleDeleteEnvironmentConfirm = async (e, environment) => {
    e.stopPropagation();
    let response = await modalRef.current.show({
      title: `Delete "${environment.name}" environment?`,
      subtitle: deleteEnvironmentModalContent(environment.name),
      primaryOption: 'DELETE',
      showInfoIcon: `Deleting an environment does not delete any resources (e.g. connections) currently contained with the environment.
      Resources that belong to others environments will continue to belong to those other environments.
      Learn more about the behavior of [lifecycle of environments and their resources](https://docs.meshery.io/concepts/logical/environments) in Meshery Docs.`,
      variant: PROMPT_VARIANTS.DANGER,
    });
    if (response === 'DELETE') {
      handleDeleteEnvironment(environment.id);
    }
  };

  const handleDeleteEnvironment = (id) => {
    deleteEnvironment({
      environmentId: id,
    })
      .unwrap()
      .then(() => handleSuccess(`Environment deleted`))
      .catch(handleError({ error_msg: 'Unable to delete environment' }));
  };

  const deleteEnvironmentModalContent = (environment) => (
    <>
      <p>Are you sure you want to delete this environment? (This action is irreversible)</p>
      <p>
        `Environment Name: `
        <i>
          <b>{environment}</b>
        </i>
      </p>
    </>
  );

  const handleBulkSelect = (e, id) => {
    const isChecked = e.target.checked;
    if (isChecked) {
      setSelectedEnvironments([...selectedEnvironments, id]);
    } else {
      const newSelectedEnv = selectedEnvironments.filter((env) => env !== id);
      setSelectedEnvironments(newSelectedEnv);
    }
  };

  const handleBulkDeleteEnvironmentConfirm = async (e) => {
    e.stopPropagation();
    let response = await modalRef.current.show({
      title: `Delete Environment(s) ?`,
      subtitle: `Do you want to delete ${selectedEnvironments.length} environment(s) ?`,
      primaryOption: 'DELETE',
      variant: PROMPT_VARIANTS.DANGER,
    });
    if (response === 'DELETE') {
      handleBulkDeleteEnv();
    }
  };

  const handleBulkDeleteEnv = () => {
    selectedEnvironments.map((envId) => {
      handleDeleteEnvironment(envId);
    });
    setSelectedEnvironments([]);
  };

  const handleAssignConnection = async () => {
    const { addedConnectionsIds, removedConnectionsIds } =
      getAddedAndRemovedConnection(assignedConnections);

    // Await every assignment/removal before closing so a rejected mutation
    // surfaces its error toast while the modal is still open, rather than the
    // modal vanishing as if the change had succeeded. allSettled keeps a single
    // failure from aborting the rest; each rejection is reported by the
    // mutators' own .catch above.
    await Promise.allSettled([
      ...addedConnectionsIds.map((id) => addConnectionToEnvironment(connectionAssignEnv.id, id)),
      ...removedConnectionsIds.map((id) =>
        removeConnectionFromEnvironment(connectionAssignEnv.id, id),
      ),
    ]);
    setEnvironmentConnectionsData([]);
    setConnectionsData([]);
    handleonAssignConnectionModalClose();
  };

  const handleonAssignConnectionModalOpen = (e, environment) => {
    e.stopPropagation();
    setAssignConnectionModal(true);
    if (connectionAssignEnv.id !== environment.id) {
      setConnectionsData([]);
      setEnvironmentConnectionsData([]);
    }
    setConnectionAssignEnv(environment);
    setSkip(false);
  };

  const handleonAssignConnectionModalClose = () => {
    setAssignConnectionModal(false);
    setSkip(true);
  };

  const handleAssignConnectionData = (updatedAssignedData) => {
    const { addedConnectionsIds, removedConnectionsIds } =
      getAddedAndRemovedConnection(updatedAssignedData);
    (addedConnectionsIds.length > 0 || removedConnectionsIds.length) > 0 &&
    (canAssignConnToEnv || canRemoveConnFromEnv)
      ? setDisableTranferButton(false)
      : setDisableTranferButton(true);

    setAssignedConnections(updatedAssignedData);
  };

  const getAddedAndRemovedConnection = (allAssignedConnections) => {
    const originalConnectionsIds = environmentConnectionsData.map((conn) => conn.id);
    const updatedConnectionsIds = allAssignedConnections.map((conn) => conn.id);

    const addedConnectionsIds = updatedConnectionsIds.filter(
      (id) => !originalConnectionsIds.includes(id),
    );
    const removedConnectionsIds = originalConnectionsIds.filter(
      (id) => !updatedConnectionsIds.includes(id),
    );

    return {
      addedConnectionsIds,
      removedConnectionsIds,
    };
  };

  const handleAssignablePage = () => {
    const pagesCount = Math.ceil(Number(connections?.totalCount ?? 0) / connectionPageSize);
    if (connectionsPage < pagesCount - 1) {
      setConnectionsPage((prevConnectionsPage) => prevConnectionsPage + 1);
    }
  };

  const handleAssignedPage = () => {
    const pagesCount = Math.ceil(
      Number(environmentConnections?.totalCount ?? 0) / connectionPageSize,
    );
    if (connectionsOfEnvironmentPage < pagesCount - 1) {
      setConnectionsOfEnvironmentPage(
        (prevConnectionsOfEnvironmentPage) => prevConnectionsOfEnvironmentPage + 1,
      );
    }
  };

  return (
    <NoSsr>
      {canViewEnv ? (
        <>
          <ToolWrapper>
            <CreateButtonWrapper style={{ marginRight: '2rem' }}>
              <Button
                type="submit"
                variant="contained"
                color="primary"
                size="large"
                onClick={(e) => handleEnvironmentModalOpen(e, ACTION_TYPES.CREATE)}
                sx={{
                  padding: '8px',
                  borderRadius: '5px',
                }}
                permissionKey={Keys.WorkspaceManagementCreateEnvironment}
                data-cy="btnResetDatabase"
              >
                <AddIconCircleBorder sx={{ width: '20px', height: '20px' }} />
                <Typography
                  sx={{
                    paddingLeft: '4px',
                    marginRight: '4px',
                  }}
                >
                  Create
                </Typography>
              </Button>
            </CreateButtonWrapper>
            <SearchBar
              onSearch={(value) => {
                setSearch(value);
              }}
              placeholder="Search Environments..."
              expanded={isSearchExpanded}
              setExpanded={setIsSearchExpanded}
            />
          </ToolWrapper>
          {selectedEnvironments.length > 0 && (
            <BulkActionWrapper>
              <Typography>
                {selectedEnvironments.length > 1
                  ? `${selectedEnvironments.length} environments selected`
                  : `${selectedEnvironments.length} environment selected`}
              </Typography>
              <Button
                onClick={handleBulkDeleteEnvironmentConfirm}
                permissionKey={Keys.WorkspaceManagementDeleteEnvironment}
              >
                <DeleteIcon fill="red" style={{ margin: '0 2px' }} />
              </Button>
            </BulkActionWrapper>
          )}
          {environments.length > 0 ? (
            <>
              <Grid2 container spacing={2} sx={{ marginTop: '10px' }} size="grow">
                {environments.map((environment) => (
                  <Grid2 key={environment.id} size={{ xs: 12, md: 6 }}>
                    <EnvironmentCard
                      // classes={classes}
                      environmentDetails={environment}
                      selectedEnvironments={selectedEnvironments}
                      onEdit={(e) => handleEnvironmentModalOpen(e, ACTION_TYPES.EDIT, environment)}
                      onDelete={(e) => handleDeleteEnvironmentConfirm(e, environment)}
                      onSelect={(e) => handleBulkSelect(e, environment.id)}
                      onAssignConnection={(e) => handleonAssignConnectionModalOpen(e, environment)}
                    />
                  </Grid2>
                ))}
              </Grid2>
              <Grid2
                container
                sx={{ padding: '2rem 0', marginTop: '20px' }}
                flex
                justifyContent="center"
                spacing={2}
                size="grow"
              >
                <Pagination
                  count={Math.ceil(environmentsData?.totalCount / pageSize)}
                  page={page + 1}
                  onChange={debounce((_, page) => setPage(page - 1), 150)}
                  boundaryCount={3}
                  renderItem={(item) => (
                    <PaginationItem
                      slots={{ previous: ChevronLeftIcon, next: ChevronRightIcon }}
                      {...item}
                    />
                  )}
                />
              </Grid2>
            </>
          ) : (
            <EmptyState
              icon={
                <EnvironmentIcon
                  height="6rem"
                  width="6rem"
                  fill={theme.palette.icon.default}
                  secondaryFill={theme.palette.icon.secondary}
                />
              }
              message="No environment available"
              pointerLabel="Click “Create” to establish your first environment."
            />
          )}
          {(canCreateEnv || canEditEnv) && environmentModal.open && (
            <SisitentModal
              open={environmentModal.open}
              closeModal={handleEnvironmentModalClose}
              title={actionType === ACTION_TYPES.CREATE ? 'Create Environment' : 'Edit Environment'}
            >
              <RJSFModalWrapper
                schema={environmentModal.schema.schema}
                uiSchema={environmentModal.schema.uischema}
                handleSubmit={
                  actionType === ACTION_TYPES.CREATE
                    ? handleCreateEnvironment
                    : handleEditEnvironment
                }
                submitBtnText={actionType === ACTION_TYPES.CREATE ? 'Save' : 'Update'}
                initialData={initialData}
                handleClose={handleEnvironmentModalClose}
              />
            </SisitentModal>
          )}
          <SisitentModal
            open={assignConnectionModal}
            closeModal={handleonAssignConnectionModalClose}
            title={`${connectionAssignEnv.name} Resources`}
            headerIcon={<EnvironmentIcon height="2rem" width="2rem" fill="white" />}
            maxWidth="md"
          >
            <ModalBody>
              <TransferList
                name="Connections"
                assignableData={connectionsData}
                assignedData={handleAssignConnectionData}
                originalAssignedData={environmentConnectionsData}
                emptyStateIconLeft={
                  <ConnectionIcon
                    width="120"
                    primaryFill={theme.palette.icon.default}
                    secondaryFill={theme.palette.icon.secondary}
                  />
                }
                emtyStateMessageLeft="No connections available"
                emptyStateIconRight={
                  <ConnectionIcon
                    width="120"
                    primaryFill={theme.palette.icon.default}
                    secondaryFill={theme.palette.icon.secondary}
                  />
                }
                emtyStateMessageRight="No connections assigned"
                transferComponentType={TRANSFER_COMPONENT.CHIP}
                assignablePage={handleAssignablePage}
                assignedPage={handleAssignedPage}
                originalLeftCount={connections?.totalCount}
                originalRightCount={environmentConnections?.totalCount}
                leftPermission={canRemoveConnFromEnv}
                rightPermission={canAssignConnToEnv}
              />
            </ModalBody>
            <ModalFooter variant="filled" helpText="Assign connections to environment">
              <PrimaryActionButtons
                primaryText="Save"
                secondaryText="Cancel"
                primaryButtonProps={{
                  onClick: handleAssignConnection,
                  disabled: disableTranferButton,
                }}
                secondaryButtonProps={{
                  onClick: handleonAssignConnectionModalClose,
                }}
              />
            </ModalFooter>
          </SisitentModal>
          <_PromptComponent ref={modalRef} />
        </>
      ) : (
        <DefaultError />
      )}
    </NoSsr>
  );
};

const EnvironmentsPageWithErrorBoundary = (props) => {
  return (
    <NoSsr>
      <ErrorBoundary>
        <Environments {...props} />
      </ErrorBoundary>
    </NoSsr>
  );
};

export default withRouter(EnvironmentsPageWithErrorBoundary);
