import { Button, Grid, NoSsr, Typography, withStyles } from '@material-ui/core';
import { Provider, connect } from 'react-redux';
import { withRouter } from 'next/router';
import { Pagination, PaginationItem } from '@material-ui/lab';
import ChevronLeftIcon from '@mui/icons-material/ChevronLeft';
import ChevronRightIcon from '@mui/icons-material/ChevronRight';
import DesignsIcon from '../../assets/icons/DesignIcon';

import { store } from '../../store';
import WorkspaceIcon from '../../assets/icons/Workspace';
import EmptyState from '../Environments/empty-state';
import useStyles from '../../assets/styles/general/tool.styles';
import styles from '../Environments/styles';
import SearchBar from '../../utils/custom-search';
import AddIconCircleBorder from '../../assets/icons/AddIconCircleBorder';
import { useEffect, useRef, useState } from 'react';
import {
  useAssignDesignToWorkspaceMutation,
  useAssignEnvironmentToWorkspaceMutation,
  useCreateWorkspaceMutation,
  useDeleteWorkspaceMutation,
  useGetDesignsOfWorkspaceQuery,
  useGetEnvironmentsOfWorkspaceQuery,
  useGetWorkspacesQuery,
  useUnassignDesignFromWorkspaceMutation,
  useUnassignEnvironmentFromWorkspaceMutation,
  useUpdateWorkspaceMutation,
} from '../../rtk-query/workspace';
import { useGetEnvironmentsQuery } from '../../rtk-query/environments';
import { useGetPatternsQuery } from '../../rtk-query/design';
import dataFetch from '../../lib/data-fetch';
import { updateProgress } from '../../lib/store';
import { useNotification } from '../../utils/hooks/useNotification';
import WorkspaceCard from './workspace-card';
import Modal from '../Modal';
import PromptComponent, { PROMPT_VARIANTS } from '../PromptComponent';
import { debounce } from 'lodash';
import { EVENT_TYPES } from '../../lib/event-types';
import EnvironmentIcon from '../../assets/icons/Environment';
import GenericModal from '../Environments/generic-modal';
import TransferList from '../Environments/transfer-list/transfer-list';

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

  const [assignEnvironmentModal, setAssignEnvironmentModal] = useState(false);
  const [environmentAssignWorkspace, setEnvironmentAssignWorkspace] = useState({});
  const [environmentsData, setEnvironmentsData] = useState([]);
  const [environmentsPage, setEnvironmentsPage] = useState(0);
  const [environmentsPageSize /*setEnvironmentsPageSize*/] = useState(25);
  const [workspaceEnvironmentsData, setWorkspaceEnvironmentsData] = useState([]);
  const [skipEnvironments, setSkipEnvironments] = useState(true);
  const [assignedEnvironments, setAssignedEnvironments] = useState([]);
  const [environmentsOfWorkspacePage, setEnvironmentsOfWorkspacePage] = useState(0);
  const [assignDesignModal, setAssignDesignModal] = useState(false);
  const [designAssignWorkspace, setDesignAssignWorkspace] = useState({});
  const [designsData, setDesignsData] = useState([]);
  const [workspaceDesignsData, setWorkspaceDesignsData] = useState([]);
  const [assignedDesigns, setAssignedDesigns] = useState([]);
  const [skipDesigns, setSkipDesigns] = useState(true);
  const [designsOfWorkspacePage, setDesignsOfWorkspacePage] = useState(0);
  const [designsPage, setDesignsPage] = useState(0);
  const [designsPageSize, /*setDesignssPageSize*/] = useState(25);

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

  const {
    data: environments,
    // isLoading: isEnvironmentsLoading,
    isError: isEnvironmentsError,
    error: environmentsError,
  } = useGetEnvironmentsQuery(
    {
      orgId: environmentAssignWorkspace.organization_id,
      page: environmentsPage,
      pagesize: environmentsPageSize,
    },
    {
      skip: skipEnvironments,
    },
  );

  const {
    data: environmentsOfWorkspace,
    // isLoading: isEnvironmentsOfWorkspaceLoading,
    isError: isEnvironmentsOfWorkspaceError,
    error: environmentsOfWorkspaceError,
  } = useGetEnvironmentsOfWorkspaceQuery(
    {
      workspaceId: environmentAssignWorkspace.id,
      page: environmentsOfWorkspacePage,
      pagesize: environmentsPageSize,
    },
    {
      skip: skipEnvironments,
    },
  );

  const [assignEnvironmentToWorkspace] = useAssignEnvironmentToWorkspaceMutation();

  const [unassignEnvironmentFromWorkspace] = useUnassignEnvironmentFromWorkspaceMutation();

  const { data: designs } = useGetPatternsQuery({
    page: designsPage,
    pagesize: designsPageSize
  },
    {
      skip: skipDesigns
    }
  );

  const {
    data: designsOfWorkspace,
    // isLoading: isDesignsOfWorkspaceLoading,
    isError: isDesignsOfWorkspaceError,
    error: designsOfWorkspaceError
  } = useGetDesignsOfWorkspaceQuery({
    workspaceId: designAssignWorkspace.id,
    page: designsOfWorkspacePage,
    pagesize: designsPageSize
  },
    {
      skip: skipDesigns
    }
  );

  const [assignDesignToWorkspace] = useAssignDesignToWorkspaceMutation();

  const [unassignDesignFromWorkspace] = useUnassignDesignFromWorkspaceMutation();


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
    if (isEnvironmentsError) {
      handleError(`Environments Fetching Error: ${environmentsError?.data}`);
    }
    if (isEnvironmentsOfWorkspaceError) {
      handleError(
        `Environments of Workspace Fetching Error: ${environmentsOfWorkspaceError?.data}`,
      );
    }
    if (isDesignsOfWorkspaceError) {
      handleError(
        `Designs of Workspace Fetching Error: ${designsOfWorkspaceError?.data}`
      );
    }
  }, [
    isWorkspacesError,
    workspacesError,
    isEnvironmentsError,
    environmentsError,
    isEnvironmentsOfWorkspaceError,
    environmentsOfWorkspaceError,
    isDesignsOfWorkspaceError,
    designsOfWorkspaceError,
    handleError,
  ]);

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

  useEffect(() => {
    const pagesCount = parseInt(
      Math.ceil(parseInt(environments?.total_count) / environmentsPageSize),
    );
    if (environments) {
      /* eslint-disable-next-line no-unsafe-optional-chaining */
      setEnvironmentsData((prevData) => [...prevData, ...environments?.environments]);
      if (environments?.total_count && environmentsPage < pagesCount - 1) {
        setEnvironmentsPage((prevEnvironmentsPage) => prevEnvironmentsPage + 1);
      }
    }
  }, [environments, environmentsPage, environmentsPageSize]);

  useEffect(() => {
    const pagesCount = 0;
    parseInt(Math.ceil(parseInt(environmentsOfWorkspace?.total_count) / environmentsPageSize));
    if (environmentsOfWorkspace) {
      /* eslint-disable-next-line no-unsafe-optional-chaining */
      setWorkspaceEnvironmentsData((prevWorkspaceEnvironmentsData) => [
        ...prevWorkspaceEnvironmentsData,
        ...environmentsOfWorkspace?.environments,
      ]);
      if (environmentsOfWorkspace?.total_count && environmentsOfWorkspacePage < pagesCount - 1) {
        setEnvironmentsOfWorkspacePage(
          (prevEnvironmentsOfWorkspacePage) => prevEnvironmentsOfWorkspacePage + 1,
        );
      }
    }
  }, [environmentsOfWorkspace, environmentsOfWorkspacePage, environmentsPageSize]);

  const handleAssignEnvironmentModalClose = () => {
    setAssignEnvironmentModal(false);
    setSkipEnvironments(true);
  };

  const handleAssignEnvironmentModalOpen = (e, workspace) => {
    e.stopPropagation();
    setAssignEnvironmentModal(true);
    if (environmentAssignWorkspace.id !== workspace.id) {
      setWorkspaceEnvironmentsData([]);
    }
    setEnvironmentAssignWorkspace(workspace);
    setSkipEnvironments(false);
  };

  const handleAssignEnvironmentsData = (updatedAssignedData) => {
    setAssignedEnvironments(updatedAssignedData);
  };

  const handleAssignEnvironments = () => {
    const originalEnvironmentsIds = workspaceEnvironmentsData.map((environment) => environment.id);
    const updatedEnvironmentsIds = assignedEnvironments.map((environment) => environment.id);

    const addedEnvironmentsIds = updatedEnvironmentsIds.filter(
      (id) => !originalEnvironmentsIds.includes(id),
    );
    const removedTeamsIds = originalEnvironmentsIds.filter(
      (id) => !updatedEnvironmentsIds.includes(id),
    );

    addedEnvironmentsIds.map((id) =>
      assignEnvironmentToWorkspace({
        workspaceId: environmentAssignWorkspace.id,
        environmentId: id,
      }).unwrap(),
    );

    removedTeamsIds.map((id) =>
      unassignEnvironmentFromWorkspace({
        workspaceId: environmentAssignWorkspace.id,
        environmentId: id,
      }).unwrap(),
    );
    setWorkspaceEnvironmentsData([]);
    handleAssignEnvironmentModalClose();
  };

  const handleAssignDesignModalClose = () => {
    setAssignDesignModal(false);
    setDesignAssignWorkspace({});
    setSkipDesigns(true);
  };

  const handleAssignDesignModalOpen = (e, workspace) => {
    e.stopPropagation();
    setAssignDesignModal(true);
    setDesignAssignWorkspace(workspace);
    setSkipDesigns(false);
  };

  const handleAssignDesignsData = updatedAssignedData => {
    setAssignedDesigns(updatedAssignedData);
  };

  const handleAssignDesigns = () => {
    const originalDesignsIds = workspaceDesignsData.map(design => design.id);
    const updatedDesignsIds = assignedDesigns.map(design => design.id);

    const addedDesignsIds = updatedDesignsIds.filter(
      id => !originalDesignsIds.includes(id)
    );
    const removedDesignsIds = originalDesignsIds.filter(
      id => !updatedDesignsIds.includes(id)
    );

    addedDesignsIds.map(id =>
      assignDesignToWorkspace({
        workspaceId: designAssignWorkspace.id,
        designId: id
      }).unwrap()
    );

    removedDesignsIds.map(id =>
      unassignDesignFromWorkspace({
        workspaceId: designAssignWorkspace.id,
        designId: id
      }).unwrap()
    );
    handleAssignDesignModalClose();
  };

  useEffect(() => {
    const pagesCount = parseInt(
      Math.ceil(parseInt(designs?.total_count) / designsPageSize)
    );
    if (designs) {
      /* eslint-disable-next-line no-unsafe-optional-chaining */
      setDesignsData(prevData => [...prevData, ...designs?.patterns]);
        if (designs?.total_count && designsPage < pagesCount - 1) {
          setDesignsPage(prevDesignsPage => prevDesignsPage + 1);
        }
    }
  }, [designs, designsPage, designsPageSize]);

  useEffect(() => {
    const pagesCount = parseInt(
      Math.ceil(parseInt(designsOfWorkspace?.total_count) / designsPageSize)
    );
    if (designsOfWorkspace) {
      /* eslint-disable-next-line no-unsafe-optional-chaining */
      setWorkspaceDesignsData(prevData => [...prevData, ...designsOfWorkspace?.designs]);
      if (designsOfWorkspace?.total_count && designsOfWorkspacePage < pagesCount - 1) {
        setDesignsOfWorkspacePage(
          prevDesignsOfWorkspacePage => prevDesignsOfWorkspacePage + 1
        );
      }
    }
  }, [designsOfWorkspace, designsOfWorkspacePage, designsPageSize]);

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
                  onAssignEnvironment={(e) => handleAssignEnvironmentModalOpen(e, workspace)}
                  onAssignDesign={e => handleAssignDesignModalOpen(e, workspace)}
                  classes={classes}
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
      <GenericModal
        open={assignEnvironmentModal}
        handleClose={handleAssignEnvironmentModalClose}
        title={`Assign Environments to ${environmentAssignWorkspace.name}`}
        body={
          <TransferList
            name="Environments"
            assignableData={environmentsData}
            assignedData={handleAssignEnvironmentsData}
            originalAssignedData={workspaceEnvironmentsData}
            emptyStateIconLeft={
              <EnvironmentIcon height="5rem" width="5rem" fill="#808080" secondaryFill="#979797" />
            }
            emtyStateMessageLeft="No environments available"
            emptyStateIconRight={
              <EnvironmentIcon height="5rem" width="5rem" fill="#808080" secondaryFill="#979797" />
            }
            emtyStateMessageRight="No environments assigned"
          />
        }
        action={handleAssignEnvironments}
        buttonTitle="Save"
        leftHeaderIcon={<EnvironmentIcon height="2rem" width="2rem" fill="white" />}
        helpText="Assign environment to workspace"
        maxWidth="md"
      />
      <GenericModal
        open={assignDesignModal}
        handleClose={handleAssignDesignModalClose}
        title={`Assign Designs to ${designAssignWorkspace.name}`}
        body={
          <TransferList
            name="Designs"
            assignableData={designsData}
            assignedData={handleAssignDesignsData}
            originalAssignedData={workspaceDesignsData}
            emptyStateIconLeft={
              <DesignsIcon
                height="5rem"
                width="5rem"
              />
            }
            emtyStateMessageLeft="No designs available"
            emptyStateIconRight={
              <DesignsIcon
                height="5rem"
                width="5rem"
              />
            }
            emtyStateMessageRight="No designs assigned"
          />
        }
        action={handleAssignDesigns}
        buttonTitle="Save"
        leftHeaderIcon={
          <DesignsIcon
            height="2rem"
            width="2rem"
            fill="#ffffff"
          />
        }
        helpText="Assign designs to workspace"
        maxWidth="md"
      />
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
