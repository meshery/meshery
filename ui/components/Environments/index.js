import { Grid, NoSsr } from '@material-ui/core';
import { withStyles } from '@material-ui/core/styles';
import { connect } from 'react-redux';
import AddIcon from '@mui/icons-material/Add';

import EnvironmentCard from './environment-card';
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
import PromptComponent from '../PromptComponent';
import { withRouter } from 'next/router';

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
  const [, /* totalCount */ setTotalCount] = useState();
  const [, /* loading */ setLoading] = useState(false);
  const [orgId, setOrgId] = useState('');
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
      (res) => {
        notify({ message: 'Environment created', event_type: EVENT_TYPES.SUCCESS });
        console.log('🚀 ~ file: index.js:128 ~ createEnvironment ~ res:', res);
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
      (res) => {
        notify({ message: 'Environment updated', event_type: EVENT_TYPES.SUCCESS });
        console.log('🚀 ~ file: index.js:146 ~ updateEnvironment ~ res:', res);
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
      (res) => {
        notify({ message: 'Environment deleted', event_type: EVENT_TYPES.SUCCESS });
        console.log('🚀 ~ file: index.js:107 ~ createEnvironment ~ res:', res);
        getEnvironments();
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
        Environment Name:
        <i>
          <b>{environment}</b>
        </i>
      </p>
    </>
  );

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
      <Grid container spacing={2} sx={{ marginTop: '10px' }}>
        {environments.map((environment) => (
          <Grid item xs={12} md={6} key={environment.id}>
            <EnvironmentCard
              environmentDetails={environment}
              // selectedEnvironments={selectedEnvironments}
              onEdit={(e) => handleEnvironmentModalOpen(e, ACTION_TYPES.EDIT, environment)}
              onDelete={(e) => handleDeleteEnvironmentConfirm(e, environment)}
              // onSelect={e => handleBulkSelect(e, environment.id)}
              // onAssignConnection={e =>
              //   handleonAssignConnectionModalOpen(e, environment)
              // }
            />
          </Grid>
        ))}
      </Grid>

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

export default withStyles()(connect(mapStateToProps)(withRouter(Environments)));
