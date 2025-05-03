import InfoIcon from '@/assets/icons/InfoIcon';
import {
  useAssignEnvironmentToWorkspaceMutation,
  useGetEnvironmentsOfWorkspaceQuery,
  useGetWorkspacesQuery,
  useUnassignEnvironmentFromWorkspaceMutation,
} from '@/rtk-query/workspace';
import CAN from '@/utils/can';
import { useNotificationHandlers } from '@/utils/hooks/useNotification';
import { keys } from '@/utils/permission_constants';
import { getColumnValue } from '@/utils/utils';
import {
  AuthorCell,
  Box,
  CustomTooltip,
  Grid,
  IconButton,
  ResponsiveDataTable,
  TableCell,
  Typography,
  updateVisibleColumns,
  useTheme,
  useWindowDimensions,
  WorkspaceEnvironmentSelection,
  WorkspaceIcon,
  Slide,
} from '@layer5/sistent';
import { useLegacySelector } from 'lib/store';
import { useEffect, useState } from 'react';
import { iconSmall } from 'css/icons.styles';
import WorkSpaceContentDataTable from './WorkSpaceContentDataTable';
import WorkspaceActionList from './WorkspaceActionList';

const WorkspaceDataTable = ({
  handleWorkspaceModalOpen,
  handleTeamsModalOpen,
  handleActivityModalOpen,
  handleDeleteWorkspaceConfirm,
  columnVisibility,
  selectedWorkspace,
  handleRowClick,
  setColumnVisibility,
  search,
  viewType,
}) => {
  let colViews = [
    ['id', 'na'],
    ['name', 'xs'],
    ['owner_email', 'na'],
    ['description', 'm'],
    ['owner', 'xs'],
    ['actions', 'xs'],
    ['environments', 'm'],
    ['updated_at', 'l'],
    ['created_at', 'na'],
  ];

  const [page, setPage] = useState(0);
  const [pageSize, setPageSize] = useState(10);
  const [sortOrder, setSortOrder] = useState('updated_at desc');
  const org_id = useLegacySelector((state) => {
    return typeof state?.get === 'function'
      ? state.get('organization')?.id
      : state?.organization?.id || '';
  });

  const theme = useTheme();

  const { data: workspaces } = useGetWorkspacesQuery({
    page: page,
    pagesize: pageSize,
    search: search,
    order: sortOrder,
    orgId: org_id,
  });

  const workspacesData = workspaces?.workspaces ? workspaces.workspaces : [];

  const columns = [
    {
      name: 'id',
      label: 'ID',
      filter: false,
    },
    {
      name: 'owner_email',
      label: 'Owner Email',
      filter: false,
    },
    {
      name: 'owner_id',
      label: 'Owner Id',
      filter: false,
      options: {
        display: 'excluded',
      },
    },
    {
      name: 'org_name',
      label: 'Org Name',
      filter: false,
      options: {
        display: 'excluded',
      },
    },
    {
      name: 'name',
      label: 'Name',
      options: {
        filter: false,
        sort: true,
        customBodyRender: (value) => {
          return (
            <Box display="flex" gap={'0.5rem'}>
              <WorkspaceIcon {...iconSmall} fill={theme.palette.icon.default} />
              {value}
            </Box>
          );
        },
      },
    },
    {
      name: 'description',
      label: 'Description',
    },
    {
      name: 'owner_avatar',
      label: 'Owner Avatar',
      options: {
        display: 'excluded',
      },
    },
    {
      name: 'owner',
      label: 'Owner',
      options: {
        filter: false,
        sort: true,
        searchable: true,
        customBodyRender: (value, tableMeta) => {
          if (!value) {
            return <span>{getColumnValue(tableMeta.rowData, 'org_name', columns)}</span>;
          }
          return (
            <AuthorCell
              userId={getColumnValue(tableMeta.rowData, 'owner_id', columns)}
              firstName={value}
              avatarUrl={getColumnValue(tableMeta.rowData, 'owner_avatar', columns)}
            />
          );
        },
      },
    },
    {
      name: 'environments',
      label: 'Environments',
      options: {
        sort: false,
        sortThirdClickReset: true,
        customHeadRender: function CustomHead({ ...column }) {
          return (
            <>
              <TableCell>
                <Grid style={{ display: 'flex' }}>
                  <Grid style={{ display: 'flex', alignItems: 'center' }}>
                    <Typography>
                      <b>{column.label}</b>
                    </Typography>
                    <CustomTooltip
                      title={`Meshery Environments allow you to logically group related Connections and their associated Credentials. [Learn more](https://docs.meshery.io/concepts/logical/environments)`}
                    >
                      <Typography style={{ display: 'flex', marginLeft: '5px' }} variant="span">
                        <IconButton disableRipple={true} disableFocusRipple={true}>
                          <InfoIcon
                            style={{
                              cursor: 'pointer',
                              height: 20,
                              width: 20,
                            }}
                            onClick={(e) => {
                              e.stopPropagation();
                            }}
                          />
                        </IconButton>
                      </Typography>
                    </CustomTooltip>
                  </Grid>
                </Grid>
              </TableCell>
            </>
          );
        },
        customBodyRender: (value, tableMeta) => {
          const workspaceId = getColumnValue(tableMeta.rowData, 'id', columns);
          return (
            <WorkspaceEnvironmentSelection
              workspaceId={workspaceId}
              useAssignEnvironmentToWorkspaceMutation={useAssignEnvironmentToWorkspaceMutation}
              useGetEnvironmentsOfWorkspaceQuery={useGetEnvironmentsOfWorkspaceQuery}
              useUnassignEnvironmentFromWorkspaceMutation={
                useUnassignEnvironmentFromWorkspaceMutation
              }
              useNotificationHandlers={useNotificationHandlers}
              isAssignedEnvironmentAllowed={CAN(
                keys.ASSIGN_ENVIRONMENT_TO_WORKSPACE.action,
                keys.ASSIGN_ENVIRONMENT_TO_WORKSPACE.subject,
              )}
            />
          );
        },
      },
    },
    {
      name: 'created_at',
      label: 'Created At',
      options: {
        filter: false,
        sort: true,
        searchable: true,
      },
    },
    {
      name: 'updated_at',
      label: 'Updated At',
      options: {
        filter: false,
        sort: true,
        searchable: true,
      },
    },
    {
      name: 'actions',
      label: 'Actions',
      options: {
        sort: false,
        customBodyRender: (value, tableMeta) => {
          const workspaceId = getColumnValue(tableMeta.rowData, 'id', columns);
          const workspaceName = getColumnValue(tableMeta.rowData, 'name', columns);

          return (
            <WorkspaceActionList
              handleActivityModalOpen={handleActivityModalOpen}
              handleTeamsModalOpen={handleTeamsModalOpen}
              handleWorkspaceModalOpen={handleWorkspaceModalOpen}
              handleDeleteWorkspaceConfirm={handleDeleteWorkspaceConfirm}
              workspaceId={workspaceId}
              workspaceName={workspaceName}
              selectedWorkspace={workspacesData[tableMeta.rowIndex]}
            />
          );
        },
      },
    },
  ];

  // Window dimensions for responsive column visibility
  const { width } = useWindowDimensions();

  useEffect(() => {
    let showCols = updateVisibleColumns(colViews, width);
    const initialVisibility = {};
    columns.forEach((col) => {
      initialVisibility[col.name] = showCols[col.name];
    });
    setColumnVisibility(initialVisibility);
  }, [width]);

  const options = {
    filter: false,
    viewColumns: false,
    filterType: 'dropdown',
    selectableRows: 'none',
    responsive: 'standard',
    onRowClick: handleRowClick,
    count: workspaces?.total_count,
    rowsPerPage: pageSize,
    page,
    print: false,
    download: false,
    elevation: 0,
    serverSide: true,
    search: false,
    textLabels: {
      selectedRows: {
        text: 'user(s) selected',
      },
    },
    sortOrder: {
      name: sortOrder.split(' ')[0],
      direction: sortOrder.split(' ')[1],
    },
    onTableChange: (action, tableState) => {
      const sortInfo = tableState.announceText ? tableState.announceText.split(' : ') : [];
      let order = '';
      if (tableState.activeColumn) {
        order = `${columns[tableState.activeColumn].name} desc`;
      }

      switch (action) {
        case 'changePage':
          if (tableState.selectedRows.data.length) {
            tableState.selectedRows = {
              data: [],
              lookup: {},
            };
          }

          setPage(tableState.page);
          break;
        case 'changeRowsPerPage':
          setPageSize(tableState.rowsPerPage);
          break;

        case 'sort':
          if (sortInfo.length == 2) {
            if (sortInfo[1] === 'ascending') {
              order = `${columns[tableState.activeColumn].name} asc`;
            } else {
              order = `${columns[tableState.activeColumn].name} desc`;
            }
          }
          if (order !== sortOrder) {
            setSortOrder(order);
          }
          break;
        default:
          break;
      }
    },
    setTableProps: () => ({
      style: {
        cursor: 'pointer',
      },
    }),
  };

  const [tableCols, updateCols] = useState(columns);

  return (
    <div key={`list-view-${viewType}`}>
      <Slide direction="left" in={selectedWorkspace.id ? true : false}>
        <div
          style={{
            marginTop: '1rem',
          }}
        >
          {selectedWorkspace?.id && (
            <WorkSpaceContentDataTable
              workspaceId={selectedWorkspace?.id}
              workspaceName={selectedWorkspace?.name}
            />
          )}
        </div>
      </Slide>
      <Slide direction="right" in={!selectedWorkspace.id ? true : false}>
        <div
          style={{
            marginTop: '1rem',
          }}
        >
          {!selectedWorkspace?.id && (
            <ResponsiveDataTable
              columns={columns}
              data={workspacesData}
              options={options}
              columnVisibility={columnVisibility}
              tableCols={tableCols}
              updateCols={updateCols}
            />
          )}
        </div>
      </Slide>
    </div>
  );
};

export default WorkspaceDataTable;
