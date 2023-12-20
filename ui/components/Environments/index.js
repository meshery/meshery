import { Button, Grid, NoSsr } from '@material-ui/core';
import { Provider, connect } from 'react-redux';
import AddIcon from '@mui/icons-material/Add';
import ChevronLeftIcon from '@mui/icons-material/ChevronLeft';
import ChevronRightIcon from '@mui/icons-material/ChevronRight';

import EnvironmentCard from './environment-card';
import EnvironmentIcon from '../../assets/icons/Environment';
import { useEffect, useRef, useState } from 'react';
import dataFetch from '../../lib/data-fetch';
import { EVENT_TYPES } from '../../lib/event-types';
import { updateProgress } from '../../lib/store';
import { useNotification } from '../../utils/hooks/useNotification';
import useStyles from '../../assets/styles/general/tool.styles';
import SearchBar from '../../utils/custom-search';
import { CreateButtonWrapper, EditButton, TextButton } from './styles';
import theme from '../../themes/app';
import Modal from '../Modal';
import PromptComponent, { PROMPT_VARIANTS } from '../PromptComponent';
import { withRouter } from 'next/router';
import { Box, Pagination, PaginationItem } from '@mui/material';
import { debounce } from 'lodash';
import EmptyState from './empty-state';
import { Delete } from '@material-ui/icons';
import TransferList from './transfer-list/transfer-list';
import GenericModal from './generic-modal';
import ConnectionIcon from '../../assets/icons/Connection';
import { TRANSFER_COMPONET } from '../../utils/Enum';
import {
  useAddConnectionToEnvironmentMutation,
  useRemoveConnectionFromEnvironmentMutation,
} from '../../rtk-query/environments';
import { store } from '../../store';

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

const Environments = ({ organization }) => {
  const [environments, setEnvironments] = useState([]);
  const [page, setPage] = useState(0);
  const [pageSize, setPageSize] = useState(10);
  const [sortOrder /* setSortOrder */] = useState('');
  const [search, setSearch] = useState('');
  const [filter /* setFilter */] = useState('');
  const [totalCount, setTotalCount] = useState();
  const [, /* loading */ setLoading] = useState(false);
  const [orgId, setOrgId] = useState('');
  const [selectedEnvironments, setSelectedEnvironments] = useState([]);
  const [isSearchExpanded, setIsSearchExpanded] = useState(false);
  const [environmentModal, setEnvironmentModal] = useState({
    open: false,
    schema: {},
  });
  const [actionType, setActionType] = useState('');
  const [initialData, setInitialData] = useState({});
  const [editEnvId, setEditEnvId] = useState('');
  const [orgValue, setOrgValue] = useState([]);
  const [orgLabel, setOrgLabel] = useState([]);

  const [assignConnectionModal, setAssignConnectionModal] = useState(false);
  const [connectionAssignEnv, setConnectionAssignEnv] = useState({});
  // const [environmentConnections, setEnvironmentConnections] = useState([]);
  // const [connectionsOfEnvironmentPage, setConnectionsOfEnvironmentPage] = useState(0);
  const [assignedConnections, setAssignedConnections] = useState([]);

  const [connectionsData, setConnectionsData] = useState([]);
  const [connectionsPage, setConnectionsPage] = useState(0);
  const [connectionCount, setConnectionCount] = useState(0);

  const [environmentConnectionsData, setEnvironmentConnectionsData] = useState([]);
  const [environmentConnectionCount, setEnvironmentConnectionCount] = useState(0);
  const [environmentConnectionPage, setEnvironmentConnectionPage] = useState(0);

  const connectionPageSize = 25;

  const modalRef = useRef(null);
  const { notify } = useNotification();
  const StyleClass = useStyles();

  const handleError = (action) => (error) => {
    updateProgress({ showProgress: false });
    notify({
      message: `${action.error_msg}: ${error}`,
      event_type: EVENT_TYPES.ERROR,
      details: error.toString(),
    });
  };

  useEffect(() => {
    setOrgId(organization?.id);
    fetchAvailableOrgs();
  }, [organization]);

  useEffect(() => {
    if (orgId) {
      getEnvironments(page, pageSize, search, sortOrder, filter);
    }
  }, [page, pageSize, search, sortOrder, filter, orgId]);

  const getEnvironments = (page, pageSize, search, sortOrder) => {
    setLoading(true);
    if (!search) search = '';
    dataFetch(
      `/api/environments?orgID=${orgId}&page=${page}&pagesize=${pageSize}&search=${encodeURIComponent(
        search,
      )}&order=${encodeURIComponent(sortOrder)}&filter=${filter}`,
      {
        credentials: 'include',
        method: 'GET',
      },
      (res) => {
        const { environments, page, page_size, total_count } = res;
        setEnvironments(environments);
        setPage(page);
        setPageSize(page_size);
        setTotalCount(total_count);
        setLoading(false);
      },
      handleError(ERROR_MESSAGE.FETCH_ENVIRONMENTS),
    );
  };

  const createEnvironment = (payload) => {
    setLoading(true);
    dataFetch(
      `/api/environments`,
      {
        credentials: 'include',
        method: 'POST',
        body: payload,
      },
      () => {
        notify({ message: 'Environment created', event_type: EVENT_TYPES.SUCCESS });
        getEnvironments(page, pageSize, search, sortOrder, filter);
      },
      handleError(ERROR_MESSAGE.CREATE_ENVIRONMENT),
    );
  };

  const updateEnvironment = (payload) => {
    setLoading(true);
    dataFetch(
      `/api/environments/${editEnvId}`,
      {
        credentials: 'include',
        method: 'PUT',
        body: payload,
      },
      () => {
        notify({ message: 'Environment updated', event_type: EVENT_TYPES.SUCCESS });
        getEnvironments(page, pageSize, search, sortOrder, filter);
      },
      handleError(ERROR_MESSAGE.UPDATE_ENVIRONMENT),
    );
  };

  const deleteEnvironment = (id) => {
    setLoading(true);
    dataFetch(
      `/api/environments/${id}`,
      {
        credentials: 'include',
        method: 'DELETE',
      },
      () => {
        notify({ message: 'Environment deleted', event_type: EVENT_TYPES.SUCCESS });
        getEnvironments(page, pageSize, search, sortOrder, filter);
      },
      handleError(ERROR_MESSAGE.DELETE_ENVIRONMENT),
    );
  };

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

  const fetchSchema = async (actionType) => {
    setLoading(true);
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
          actionType === ACTION_TYPES.CREATE
            ? (uiSchemaOrg['ui:widget'] = 'select')
            : (uiSchemaOrg['ui:widget'] = 'hidden');
          setEnvironmentModal({
            open: true,
            schema: res,
          });
        }
      },
    );
  };

  const getConnections = (page) => {
    setLoading(true);
    dataFetch(
      `/api/integrations/connections?pagesize=${connectionPageSize}&page=${page}`,
      {
        credentials: 'include',
        method: 'GET',
      },
      (res) => {
        /* eslint-disable-next-line no-unsafe-optional-chaining */
        setConnectionsData((prevConnectionsData) => [...prevConnectionsData, ...res?.connections]);
        setConnectionCount(res?.total_count);
        setConnectionsPage(res?.page + 1);
        setLoading(false);
      },
      handleError(ERROR_MESSAGE.FETCH_CONNECTIONS),
    );
  };

  const getEnvironmentConnections = (environmentId, page) => {
    setLoading(true);
    dataFetch(
      `/api/environments/${environmentId}/connections?pagesize=${connectionPageSize}&page=${page}`,
      {
        credentials: 'include',
        method: 'GET',
      },
      (res) => {
        /* eslint-disable-next-line no-unsafe-optional-chaining */
        setEnvironmentConnectionsData(res?.connections);
        setEnvironmentConnectionCount(res?.total_count);
        setEnvironmentConnectionPage(res?.page + 1);
        setLoading(false);
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
        name: '',
        description: '',
        organization: orgId,
      });
      setEditEnvId('');
    }
    fetchSchema(actionType);
  };

  const handleEnvironmentModalClose = () => {
    setEnvironmentModal({
      open: false,
      schema: {},
    });
    setActionType('');
  };

  const handleCreateEnvironment = ({ organization, name, description }) => {
    const payload = JSON.stringify({
      name: name,
      description: description,
      organization_id: organization,
    });
    createEnvironment(payload);
    handleEnvironmentModalClose();
  };

  const handleEditEnvironment = ({ name, description }) => {
    const payload = JSON.stringify({
      name: name,
      description: description,
      organization_id: initialData.organization,
    });
    updateEnvironment(payload);
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
    deleteEnvironment(id);
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

  useEffect(() => {
    const pagesCount = parseInt(Math.ceil(parseInt(connectionCount) / connectionPageSize));
    if (pagesCount > connectionsPage || connectionsPage === 0) {
      getConnections(connectionsPage);
    }
  }, [connectionsPage]);

  useEffect(() => {
    const pagesCount = parseInt(
      Math.ceil(parseInt(environmentConnectionCount) / connectionPageSize),
    );
    if (pagesCount > environmentConnectionPage || environmentConnectionPage === 0) {
      getEnvironmentConnections(connectionAssignEnv.id, environmentConnectionPage);
    }
  }, [connectionAssignEnv, environmentConnectionPage]);

  const handleAssignConnection = () => {
    const originalConnectionsIds = environmentConnectionsData.map((conn) => conn.id);
    const updatedConnectionsIds = assignedConnections.map((conn) => conn.id);

    const addedConnectionsIds = updatedConnectionsIds.filter(
      (id) => !originalConnectionsIds.includes(id),
    );
    const removedConnectionsIds = originalConnectionsIds.filter(
      (id) => !updatedConnectionsIds.includes(id),
    );

    addedConnectionsIds.map((id) => addConnectionToEnvironment(connectionAssignEnv.id, id));

    removedConnectionsIds.map((id) => removeConnectionFromEnvironment(connectionAssignEnv.id, id));
    handleonAssignConnectionModalClose();
  };

  const handleonAssignConnectionModalOpen = (e, environment) => {
    e.stopPropagation();
    setAssignConnectionModal(true);
    setConnectionAssignEnv(environment);
  };

  const handleonAssignConnectionModalClose = () => {
    setAssignConnectionModal(false);
    setConnectionAssignEnv({});
  };

  const handleAssignConnectionData = (updatedAssignedData) => {
    setAssignedConnections(updatedAssignedData);
  };

  return (
    <NoSsr>
      <div className={StyleClass.toolWrapper} style={{ marginBottom: '20px', display: 'flex' }}>
        <CreateButtonWrapper>
          <EditButton
            variant="contained"
            onClick={(e) => handleEnvironmentModalOpen(e, ACTION_TYPES.CREATE)}
          >
            <AddIcon height="24" width="24" fill={theme.palette.secondary.white} />
            <TextButton>Create</TextButton>
          </EditButton>
        </CreateButtonWrapper>
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
        <Box
          sx={{
            width: '100%',
            p: '0.8rem',
            justifyContent: 'space-between',
            marginTop: '0.18rem',
            marginBottom: '1rem',
            backgroundColor: theme.palette.secondary.white,
            borderRadius: '.25rem',
            display: 'flex',
            flexDirection: 'row',
            alignItems: 'center',
          }}
        >
          <span style={{ color: '#000' }}>
            {selectedEnvironments.length > 1
              ? `${selectedEnvironments.length} environments selected`
              : `${selectedEnvironments.length} environment selected`}
          </span>
          <Button>
            <Delete
              style={{ color: 'red', margin: '0 2px' }}
              onClick={handleBulkDeleteEnvironmentConfirm}
              disabled={selectedEnvironments.length > 0 ? false : true}
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
              count={Math.ceil(totalCount / pageSize)}
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
            <EnvironmentIcon height="6rem" width="6rem" fill="#808080" secondaryFill="#979797" />
          }
          message="No environment available"
          pointerLabel="Click “Create” to establish your first environment."
        />
      )}

      {environmentModal.open && (
        <Modal
          open={environmentModal.open}
          schema={environmentModal.schema.rjsfSchema}
          uiSchema={environmentModal.schema.uiSchema}
          handleClose={handleEnvironmentModalClose}
          handleSubmit={
            actionType === ACTION_TYPES.CREATE ? handleCreateEnvironment : handleEditEnvironment
          }
          title={actionType === ACTION_TYPES.CREATE ? 'Create Environment' : 'Edit Environment'}
          submitBtnText={
            actionType === ACTION_TYPES.CREATE ? 'Create Environment' : 'Edit Environment'
          }
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
            transferComponentType={TRANSFER_COMPONET.CHIP}
          />
        }
        action={handleAssignConnection}
        buttonTitle="Save"
        leftHeaderIcon={<EnvironmentIcon height="2rem" width="2rem" fill="white" />}
        helpText="Assign connections to environment"
        maxWidth="md"
      />
      <PromptComponent ref={modalRef} />
    </NoSsr>
  );
};

const mapStateToProps = (state) => {
  const organization = state.get('organization');
  return {
    organization,
  };
};

export default connect(mapStateToProps)(
  withRouter((props) => (
    <Provider store={store}>
      <Environments {...props} />
    </Provider>
  )),
);
