import { useEffect, useRef, useState } from 'react';
import { Button, Grid, NoSsr, Typography, Box } from '@material-ui/core';
import { Provider, connect } from 'react-redux';
import ChevronLeftIcon from '@mui/icons-material/ChevronLeft';
import ChevronRightIcon from '@mui/icons-material/ChevronRight';
import { withStyles } from '@material-ui/core/styles';
import { Pagination, PaginationItem } from '@material-ui/lab';
import { withRouter } from 'next/router';
import { debounce } from 'lodash';
import { Delete } from '@material-ui/icons';
import classNames from 'classnames';

import AddIconCircleBorder from '../../../assets/icons/AddIconCircleBorder';
import EnvironmentCard from './environment-card';
import EnvironmentIcon from '../../../assets/icons/Environment';
import dataFetch from '../../../lib/data-fetch';
import { EVENT_TYPES } from '../../../lib/event-types';
import { updateProgress } from '../../../lib/store';
import { useNotification } from '../../../utils/hooks/useNotification';
import useStyles from '../../../assets/styles/general/tool.styles';
import SearchBar from '../../../utils/custom-search';
import Modal from '../../Modal';
import PromptComponent, { PROMPT_VARIANTS } from '../../PromptComponent';
import { EmptyState, TransferList, GenericModal } from '../General';
import ConnectionIcon from '../../../assets/icons/Connection';
import { TRANSFER_COMPONENT } from '../../../utils/Enum';
import {
  useAddConnectionToEnvironmentMutation,
  useRemoveConnectionFromEnvironmentMutation,
  useGetEnvironmentConnectionsQuery,
  useGetEnvironmentsQuery,
  useCreateEnvironmentMutation,
  useUpdateEnvironmentMutation,
  useDeleteEnvironmentMutation,
} from '../../../rtk-query/environments';
import { store } from '../../../store';
import styles from './styles';
import { keys } from '@/utils/permission_constants';
import CAN from '@/utils/can';
import DefaultError from '../../General/error-404/index';

const ERROR_MESSAGE = {
  FETCH_ENVIRONMENTS: {
    name: 'FETCH_ENVIRONMENTS',
    error_msg: 'Failed to fetch environments',
  },
  CREATE_ENVIRONMENT: {
    name: 'CREATE_ENVIRONMENT',
    error_msg: 'Failed to create environment',
  },
  UPDATE_ENVIRONMENT: {
    name: 'UPDATE_ENVIRONMENT',
    error_msg: 'Failed to update environment',
  },
  DELETE_ENVIRONMENT: {
    name: 'DELETE_ENVIRONMENT',
    error_msg: 'Failed to delete environment',
  },
  FETCH_ORGANIZATIONS: {
    name: 'FETCH_ORGANIZATIONS',
    error_msg: 'There was an error fetching available orgs',
  },
  FETCH_CONNECTIONS: {
    name: 'FETCH_CONNECTIONS',
    error_msg: 'There was an error fetching connections',
  },
};

const ACTION_TYPES = {
  CREATE: 'create',
  EDIT: 'edit',
};

const Environments = ({ organization, classes }) => {
  const [environmentModal, setEnvironmentModal] = useState({
    open: false,
    schema: {},
  });
  const [actionType, setActionType] = useState('');
  const [initialData, setInitialData] = useState({});
  const [editEnvId, setEditEnvId] = useState('');
  const [orgValue, setOrgValue] = useState([]);
  const [orgLabel, setOrgLabel] = useState([]);
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
  const StyleClass = useStyles();

  const {
    data: environmentsData,
    // isLoading: isEnvironmentsLoading,
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
      handleError(`Environments Fetching Error: ${environmentsError?.data}`);
    }
    if (isEnvironmentConnectionsError) {
      handleError(
        `Connections of a Environment fetching Error: ${environmentConnectionsError?.data}`,
      );
    }
    if (isConnectionsError) {
      handleError(`Connections fetching Error: ${connectionsError?.data}`);
    }
  }, [
    isEnvironmentsError,
    isEnvironmentConnectionsError,
    isConnectionsError,
    environmentsError,
    environmentConnectionsError,
    connectionsError,
    handleError,
  ]);

  const handleError = (action) => (error) => {
    updateProgress({ showProgress: false });
    notify({
      message: `${action.error_msg}: ${error}`,
      event_type: EVENT_TYPES.ERROR,
      details: error.toString(),
    });
  };

  const handleSuccess = (msg) => {
    updateProgress({ showProgress: false });
    notify({
      message: msg,
      event_type: EVENT_TYPES.SUCCESS,
    });
  };

  useEffect(() => {
    setOrgId(organization?.id);
    fetchAvailableOrgs();
  }, [organization]);

  const fetchAvailableOrgs = async () => {
    dataFetch(
      '/api/identity/orgs',
      {
        method: 'GET',
        credentials: 'include',
      },
      (result) => {
        if (result) {
          const label = result?.organizations.map((option) => option.name);
          const value = result?.organizations.map((option) => option.id);
          setOrgLabel(label);
          setOrgValue(value);
        }
      },
      handleError(ERROR_MESSAGE.FETCH_ORGANIZATIONS),
    );
  };

  const fetchSchema = async () => {
    dataFetch(
      `/api/schema/resource/environment`,
      {
        credentials: 'include',
        method: 'GET',
      },
      (res) => {
        if (res) {
          const rjsfSchemaOrg = res.rjsfSchema?.properties?.organization;
          const uiSchemaOrg = res.uiSchema?.organization;
          rjsfSchemaOrg.enum = orgValue;
          rjsfSchemaOrg.enumNames = orgLabel;
          uiSchemaOrg['ui:widget'] = 'hidden';
          setEnvironmentModal({
            open: true,
            schema: res,
          });
        }
      },
    );
  };

  const [addConnectionToEnvironmentMutator] = useAddConnectionToEnvironmentMutation();
  const [removeConnectionFromEnvMutator] = useRemoveConnectionFromEnvironmentMutation();

  const addConnectionToEnvironment = async (environmentId, connectionId) => {
    addConnectionToEnvironmentMutator({ environmentId, connectionId });
  };

  const removeConnectionFromEnvironment = (environmentId, connectionId) => {
    removeConnectionFromEnvMutator({ environmentId, connectionId });
  };

  const handleEnvironmentModalOpen = (e, actionType, envObject) => {
    e.stopPropagation();
    if (actionType === ACTION_TYPES.EDIT) {
      setActionType(ACTION_TYPES.EDIT);
      setInitialData({
        name: envObject.name,
        description: envObject.description,
        organization: envObject.organization_id,
      });
      setEditEnvId(envObject.id);
    } else {
      setActionType(ACTION_TYPES.CREATE);
      setInitialData({
        name: undefined,
        description: '',
        organization: orgId,
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

  const handleCreateEnvironment = ({ organization, name, description }) => {
    createEnvironment({
      environmentPayload: {
        name: name,
        description: description,
        organization_id: organization,
      },
    })
      .unwrap()
      .then(handleSuccess(`Environment "${name}" created `))
      .catch((error) => handleError(`Environment Create Error: ${error?.data}`));
    handleEnvironmentModalClose();
  };

  const handleEditEnvironment = ({ name, description }) => {
    updateEnvironment({
      environmentId: editEnvId,
      environmentPayload: {
        name: name,
        description: description,
        organization_id: initialData.organization,
      },
    })
      .unwrap()
      .then(handleSuccess(`Environment "${name}" updated`))
      .catch((error) => handleError(`Environment Update Error: ${error?.data}`));
    handleEnvironmentModalClose();
  };

  const handleDeleteEnvironmentConfirm = async (e, environment) => {
    e.stopPropagation();
    let response = await modalRef.current.show({
      title: `Delete Environment ?`,
      subtitle: deleteEnvironmentModalContent(environment.name),
      options: ['DELETE', 'CANCEL'],
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
      .then(handleSuccess(`Environment deleted`))
      .catch((error) => handleError(`Environment Delete Error: ${error?.data}`));
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
      options: ['DELETE', 'CANCEL'],
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

  const handleAssignConnection = () => {
    const { addedConnectionsIds, removedConnectionsIds } =
      getAddedAndRemovedConnection(assignedConnections);

    addedConnectionsIds.map((id) => addConnectionToEnvironment(connectionAssignEnv.id, id));

    removedConnectionsIds.map((id) => removeConnectionFromEnvironment(connectionAssignEnv.id, id));
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
    (addedConnectionsIds.length > 0 || removedConnectionsIds.length) > 0
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
    const pagesCount = parseInt(Math.ceil(parseInt(connections?.total_count) / connectionPageSize));
    if (connectionsPage < pagesCount - 1) {
      setConnectionsPage((prevConnectionsPage) => prevConnectionsPage + 1);
    }
  };

  const handleAssignedPage = () => {
    const pagesCount = parseInt(
      Math.ceil(parseInt(environmentConnections?.total_count) / connectionPageSize),
    );
    if (connectionsOfEnvironmentPage < pagesCount - 1) {
      setConnectionsOfEnvironmentPage(
        (prevConnectionsOfEnvironmentPage) => prevConnectionsOfEnvironmentPage + 1,
      );
    }
  };

  return (
    <NoSsr>
      {CAN(keys.VIEW_ENVIRONMENTS.action, keys.VIEW_ENVIRONMENTS.subject) ? (
        <>
          <div className={StyleClass.toolWrapper} style={{ marginBottom: '20px', display: 'flex' }}>
            <div className={classes.createButtonWrapper}>
              <Button
                type="submit"
                variant="contained"
                color="primary"
                size="large"
                onClick={(e) => handleEnvironmentModalOpen(e, ACTION_TYPES.CREATE)}
                style={{
                  padding: '8px',
                  borderRadius: 5,
                  marginRight: '2rem',
                }}
                disabled={!CAN(keys.CREATE_ENVIRONMENT.action, keys.CREATE_ENVIRONMENT.subject)}
                data-cy="btnResetDatabase"
              >
                <AddIconCircleBorder style={{ width: '20px', height: '20px' }} />
                <Typography
                  style={{
                    paddingLeft: '4px',
                    marginRight: '4px',
                  }}
                >
                  Create
                </Typography>
              </Button>
            </div>
            <SearchBar
              onSearch={(value) => {
                setSearch(value);
              }}
              placeholder="Search connections..."
              expanded={isSearchExpanded}
              setExpanded={setIsSearchExpanded}
            />
          </div>
          {selectedEnvironments.length > 0 && (
            <Box className={classNames(classes.bulkActionWrapper, StyleClass.toolWrapper)}>
              <Typography>
                {selectedEnvironments.length > 1
                  ? `${selectedEnvironments.length} environments selected`
                  : `${selectedEnvironments.length} environment selected`}
              </Typography>
              <Button className={classes.iconButton}>
                <Delete
                  style={{ color: 'red', margin: '0 2px' }}
                  onClick={handleBulkDeleteEnvironmentConfirm}
                  disabled={
                    selectedEnvironments.length > 0
                      ? !CAN(keys.DELETE_ENVIRONMENT.action, keys.DELETE_ENVIRONMENT.subject)
                      : true
                  }
                />
              </Button>
            </Box>
          )}
          {environments.length > 0 ? (
            <>
              <Grid container spacing={2} sx={{ marginTop: '10px' }}>
                {environments.map((environment) => (
                  <Grid item xs={12} md={6} key={environment.id}>
                    <EnvironmentCard
                      classes={classes}
                      environmentDetails={environment}
                      selectedEnvironments={selectedEnvironments}
                      onEdit={(e) => handleEnvironmentModalOpen(e, ACTION_TYPES.EDIT, environment)}
                      onDelete={(e) => handleDeleteEnvironmentConfirm(e, environment)}
                      onSelect={(e) => handleBulkSelect(e, environment.id)}
                      onAssignConnection={(e) => handleonAssignConnectionModalOpen(e, environment)}
                    />
                  </Grid>
                ))}
              </Grid>
              <Grid
                container
                sx={{ padding: '2rem 0' }}
                style={{ marginTop: '20px' }}
                flex
                justifyContent="center"
                spacing={2}
              >
                <Pagination
                  count={Math.ceil(environmentsData?.total_count / pageSize)}
                  page={page + 1}
                  sx={{
                    backgroundColor: 'white',
                    borderRadius: '1rem',
                    padding: '0.5rem',
                  }}
                  onChange={debounce((_, page) => setPage(page - 1), 150)}
                  boundaryCount={3}
                  renderItem={(item) => (
                    <PaginationItem
                      slots={{ previous: ChevronLeftIcon, next: ChevronRightIcon }}
                      {...item}
                    />
                  )}
                />
              </Grid>
            </>
          ) : (
            <EmptyState
              icon={
                <EnvironmentIcon
                  height="6rem"
                  width="6rem"
                  fill="#808080"
                  secondaryFill="#979797"
                />
              }
              message="No environment available"
              pointerLabel="Click “Create” to establish your first environment."
            />
          )}
          {(CAN(keys.CREATE_ENVIRONMENT.action, keys.CREATE_ENVIRONMENT.subject) ||
            CAN(keys.EDIT_ENVIRONMENT.action, keys.EDIT_ENVIRONMENT.subject)) &&
            environmentModal.open && (
              <Modal
                open={environmentModal.open}
                schema={environmentModal.schema.rjsfSchema}
                uiSchema={environmentModal.schema.uiSchema}
                handleClose={handleEnvironmentModalClose}
                handleSubmit={
                  actionType === ACTION_TYPES.CREATE
                    ? handleCreateEnvironment
                    : handleEditEnvironment
                }
                title={
                  actionType === ACTION_TYPES.CREATE ? 'Create Environment' : 'Edit Environment'
                }
                submitBtnText={actionType === ACTION_TYPES.CREATE ? 'Save' : 'Update'}
                initialData={initialData}
              />
            )}
          <GenericModal
            open={assignConnectionModal}
            handleClose={handleonAssignConnectionModalClose}
            title={`${connectionAssignEnv.name} Resources`}
            body={
              <TransferList
                name="Connections"
                assignableData={connectionsData}
                assignedData={handleAssignConnectionData}
                originalAssignedData={environmentConnectionsData}
                emptyStateIconLeft={
                  <ConnectionIcon width="120" primaryFill="#808080" secondaryFill="#979797" />
                }
                emtyStateMessageLeft="No connections available"
                emptyStateIconRight={
                  <ConnectionIcon width="120" primaryFill="#808080" secondaryFill="#979797" />
                }
                emtyStateMessageRight="No connections assigned"
                transferComponentType={TRANSFER_COMPONENT.CHIP}
                assignablePage={handleAssignablePage}
                assignedPage={handleAssignedPage}
                originalLeftCount={connections?.total_count}
                originalRightCount={environmentConnections?.total_count}
              />
            }
            action={handleAssignConnection}
            buttonTitle="Save"
            disabled={disableTranferButton}
            leftHeaderIcon={<EnvironmentIcon height="2rem" width="2rem" fill="white" />}
            helpText="Assign connections to environment"
            maxWidth="md"
          />
          <PromptComponent ref={modalRef} />
        </>
      ) : (
        <DefaultError />
      )}
    </NoSsr>
  );
};

const mapStateToProps = (state) => {
  const organization = state.get('organization');
  return {
    organization,
  };
};

export default withStyles(styles)(
  connect(mapStateToProps)(
    withRouter((props) => (
      <Provider store={store}>
        <Environments {...props} />
      </Provider>
    )),
  ),
);
