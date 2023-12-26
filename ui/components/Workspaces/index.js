import { Button, Grid, NoSsr, Typography, withStyles } from '@material-ui/core';
import { Provider, connect } from 'react-redux';
import { withRouter } from 'next/router';
import { Pagination, PaginationItem } from '@material-ui/lab';
import ChevronLeftIcon from '@mui/icons-material/ChevronLeft';
import ChevronRightIcon from '@mui/icons-material/ChevronRight';

import { store } from '../../store';
import WorkspaceIcon from '../../assets/icons/Workspace';
import EmptyState from '../Environments/empty-state';
import useStyles from '../../assets/styles/general/tool.styles';
import styles from '../Environments/styles';
import SearchBar from '../../utils/custom-search';
import AddIconCircleBorder from '../../assets/icons/AddIconCircleBorder';
import { useEffect, useRef, useState } from 'react';
import {
  useCreateWorkspaceMutation,
  useDeleteWorkspaceMutation,
  useGetWorkspacesQuery,
  useUpdateWorkspaceMutation,
} from '../../rtk-query/workspace';
import dataFetch from '../../lib/data-fetch';
import { updateProgress } from '../../lib/store';
import { useNotification } from '../../utils/hooks/useNotification';
import WorkspaceCard from './workspace-card';
import Modal from '../Modal';
import PromptComponent, { PROMPT_VARIANTS } from '../PromptComponent';
import { debounce } from 'lodash';
import { EVENT_TYPES } from '../../lib/event-types';

const ERROR_MESSAGE = {
  FETCH_ORGANIZATIONS: {
    name: 'FETCH_ORGANIZATIONS',
    error_msg: 'There was an error fetching available orgs',
  },
};

const ACTION_TYPES = {
  CREATE: 'create',
  EDIT: 'edit',
};

const Workspaces = ({ organization, classes }) => {
  const [workspaceModal, setWorkspaceModal] = useState({
    open: false,
    schema: {},
  });
  const [page, setPage] = useState(0);
  const [pageSize /*setPageSize*/] = useState(10);
  const [sortOrder /*setSortOrder*/] = useState('');
  const [search, setSearch] = useState('');

  const [orgId, setOrgId] = useState('');
  const [orgValue, setOrgValue] = useState([]);
  const [orgLabel, setOrgLabel] = useState([]);
  const [isSearchExpanded, setIsSearchExpanded] = useState(false);
  const [actionType, setActionType] = useState('');
  const [initialData, setInitialData] = useState({});
  const [editWorkspaceId, setEditWorkspaceId] = useState('');

  const ref = useRef(null);
  const { notify } = useNotification();
  const StyleClass = useStyles();

  const {
    data: workspacesData,
    // isLoading: isWorkspacesLoading,
    isError: isWorkspacesError,
    error: workspacesError,
  } = useGetWorkspacesQuery(
    {
      page: page,
      pagesize: pageSize,
      search: search,
      order: sortOrder,
      orgId: orgId,
    },
    {
      skip: !orgId ? true : false,
    },
  );

  const [createWorkspace] = useCreateWorkspaceMutation();

  const [updateWorkspace] = useUpdateWorkspaceMutation();

  const [deleteWorkspace] = useDeleteWorkspaceMutation();

  // const loading = isWorkspacesLoading;

  const workspaces = workspacesData?.workspaces ? workspacesData.workspaces : [];

  const handleCreateWorkspace = ({ organization, name, description }) => {
    createWorkspace({
      workspacePayload: {
        name: name,
        description: description,
        organization_id: organization,
      },
    })
      .unwrap()
      .then(() => handleSuccess(`Workspace "${name}" created `))
      .catch((error) => handleError(`Workspace Create Error: ${error?.data}`));
    handleWorkspaceModalClose();
  };

  const handleEditWorkspace = ({ organization, name, description }) => {
    updateWorkspace({
      workspaceId: editWorkspaceId,
      workspacePayload: {
        name: name,
        description: description,
        organization_id: organization,
      },
    })
      .unwrap()
      .then(() => handleSuccess(`Workspace "${name}" updated`))
      .catch((error) => handleError(`Workspace Update Error: ${error?.data}`));
    handleWorkspaceModalClose();
  };

  const handleDeleteWorkspace = (id) => {
    deleteWorkspace({
      workspaceId: id,
    })
      .unwrap()
      .then(() => handleSuccess(`Workspace deleted`))
      .catch((error) => handleError(`Workspace Delete Error: ${error?.data}`));
  };

  useEffect(() => {
    if (isWorkspacesError) {
      handleError(`Workspaces Fetching Error: ${workspacesError?.data}`);
    }
  }, [isWorkspacesError, workspacesError, handleError]);

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

  const fetchSchema = async (actionType) => {
    dataFetch(
      `/api/schema/resource/workspace`,
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
          setWorkspaceModal({
            open: true,
            schema: res,
          });
        }
      },
    );
  };

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

  const handleWorkspaceModalOpen = (e, actionType, workspaceObject) => {
    e.stopPropagation();
    if (actionType === ACTION_TYPES.EDIT) {
      setActionType(ACTION_TYPES.EDIT);
      setInitialData({
        name: workspaceObject.name,
        description: workspaceObject.description,
        organization: workspaceObject.organization_id,
      });
      setEditWorkspaceId(workspaceObject.id);
    } else {
      setActionType(ACTION_TYPES.CREATE);
      setInitialData({
        name: '',
        description: '',
        organization: orgId,
      });
      setEditWorkspaceId('');
    }
    fetchSchema(actionType);
  };

  const handleWorkspaceModalClose = () => {
    setWorkspaceModal({
      open: false,
      schema: {},
    });
    setActionType('');
  };

  const handleDeleteWorkspaceConfirm = async (e, workspace) => {
    e.stopPropagation();
    let response = await ref.current.show({
      title: `Delete workspace ?`,
      subtitle: deleteWorkspaceModalContent(workspace.name),
      options: ['DELETE', 'CANCEL'],
      variant: PROMPT_VARIANTS.DANGER,
    });
    if (response === 'DELETE') {
      handleDeleteWorkspace(workspace.id);
    }
  };

  const deleteWorkspaceModalContent = (workspace) => (
    <>
      <p>Are you sure you want to delete this workspace? (This action is irreversible)</p>
      <p>
        Workspace Name:
        <i>
          <b>{workspace}</b>
        </i>
      </p>
    </>
  );

  return (
    <NoSsr>
      <div className={StyleClass.toolWrapper} style={{ marginBottom: '20px', display: 'flex' }}>
        <div className={classes.createButtonWrapper}>
          <Button
            type="submit"
            variant="contained"
            color="primary"
            size="large"
            onClick={(e) => handleWorkspaceModalOpen(e, ACTION_TYPES.CREATE)}
            style={{
              padding: '8px',
              borderRadius: 5,
              marginRight: '2rem',
            }}
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
      {/* {selectedEnvironments.length > 0 && (
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
              disabled={selectedEnvironments.length > 0 ? false : true}
            />
          </Button>
        </Box>
      )} */}
      {workspaces.length > 0 ? (
        <>
          <Grid container spacing={2} sx={{ marginTop: '10px' }}>
            {workspaces.map((workspace) => (
              <Grid item xs={12} md={6} key={workspace.id}>
                <WorkspaceCard
                  workspaceDetails={workspace}
                  onEdit={(e) => handleWorkspaceModalOpen(e, ACTION_TYPES.EDIT, workspace)}
                  onDelete={(e) => handleDeleteWorkspaceConfirm(e, workspace)}
                  // onSelect={e => handleBulkSelect(e, workspace.id)}
                  // selectedWorkspaces={selectedWorkspaces}
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
              count={Math.ceil(workspacesData?.total_count / pageSize)}
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
          icon={<WorkspaceIcon height="6rem" width="6rem" fill="#808080" />}
          message="No workspace available"
          pointerLabel="Click “Create” to establish your first workspace."
        />
      )}
      {workspaceModal.open && (
        <Modal
          open={workspaceModal.open}
          schema={workspaceModal.schema.rjsfSchema}
          uiSchema={workspaceModal.schema.uiSchema}
          handleClose={handleWorkspaceModalClose}
          handleSubmit={
            actionType === ACTION_TYPES.CREATE ? handleCreateWorkspace : handleEditWorkspace
          }
          title={actionType === ACTION_TYPES.CREATE ? 'Create Workspace' : 'Edit Workspace'}
          submitBtnText={actionType === ACTION_TYPES.CREATE ? 'Create Workspace' : 'Edit Workspace'}
          initialData={initialData}
        />
      )}
      <PromptComponent ref={ref} />
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
        <Workspaces {...props} />
      </Provider>
    )),
  ),
);
