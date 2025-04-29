import {
  Chip,
  IconButton,
  Tooltip,
  TableCell,
  TableSortLabel,
  styled,
  ResponsiveDataTable,
} from '@layer5/sistent';
import React, { useState } from 'react';
import Modal from './Modal';
import { CONNECTION_KINDS, CON_OPS } from '../utils/Enum';
import DeleteIcon from '@mui/icons-material/Delete';
import Moment from 'react-moment';
import LoadingScreen from './LoadingComponents/LoadingComponent';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { updateProgress } from '../lib/store';
import { useNotification } from '../utils/hooks/useNotification';
import { EVENT_TYPES } from '../lib/event-types';
import { updateVisibleColumns } from '../utils/responsive-column';
import { useWindowDimensions } from '../utils/dimension';
import { CustomColumnVisibilityControl } from '@layer5/sistent';
import { ToolWrapper } from '@/assets/styles/general/tool.styles';
import {
  useCreateCredentialMutation,
  useDeleteCredentialMutation,
  useGetCredentialsQuery,
  useUpdateCredentialMutation,
} from '@/rtk-query/credentials';

const CredentialIcon = styled('img')({
  width: '24px',
  height: '24px',
});

const ActionContainer = styled('div')({
  display: 'flex',
  alignItems: 'center',
});

const CustomTableCell = styled(TableCell)({
  '& .MuiTableSortLabel-root': {
    fontWeight: 'bold',
  },
});

const schema_array = ['prometheus', 'grafana', 'kubernetes'];

const MesheryCredentialComponent = ({ updateProgress, connectionMetadataState }) => {
  const { data: credentialsData, isLoading } = useGetCredentialsQuery();
  const [createCredential] = useCreateCredentialMutation();
  const [updateCredential] = useUpdateCredentialMutation();
  const [deleteCredential] = useDeleteCredentialMutation();
  const [formData, setFormData] = useState({});
  const [credModal, setCredModal] = useState({
    open: false,
    data: null,
    actionType: null,
    id: null,
  });
  const [credentialType, setCredentialType] = useState(schema_array[0]);
  const [credentialName, setCredentialName] = useState(null);
  const { notify } = useNotification();
  const { width } = useWindowDimensions();

  const schemaChangeHandler = (type) => {
    setCredentialType(type);
    setCredModal((prev) => ({
      ...prev,
      open: true,
      data: null,
    }));
  };

  const _onChange = (formData) => {
    setCredentialName(formData?.credentialName);
    setFormData(formData);
  };

  const handleClose = (ev) => {
    ev.stopPropagation();
    setCredModal({
      open: false,
      data: null,
      actionType: null,
      id: null,
    });
  };

  const handleError = (error_msg) => {
    updateProgress({ showProgress: false });
    notify({
      message: `${error_msg}`,
      event_type: EVENT_TYPES.ERROR,
      details: error_msg.toString(),
    });
  };

  const getCredentialsIcon = (type) => {
    switch (type) {
      case 'prometheus':
        return <CredentialIcon src="/static/img/prometheus_logo_orange_circle.svg" />;
      case 'grafana':
        return <CredentialIcon src="/static/img/grafana_icon.svg" />;
      case 'kubernetes':
        return (
          <CredentialIcon
            src={
              connectionMetadataState
                ? connectionMetadataState[CONNECTION_KINDS.KUBERNETES]?.icon
                : ''
            }
          />
        );
      default:
        return null;
    }
  };

  let colViews = [
    ['name', 'xs'],
    ['type', 'l'],
    ['created_at', 'xl'],
    ['updated_at', 'xl'],
    ['actions', 'xs'],
  ];

  const columns = [
    {
      name: 'name',
      label: 'Name',
      options: {
        filter: true,
        sort: false,
        searchable: true,
        customHeadRender: function CustomHead({ index, ...column }, sortColumn) {
          return (
            <CustomTableCell key={index} onClick={() => sortColumn(index)}>
              <TableSortLabel
                active={column.sortDirection != null}
                direction={column.sortDirection || 'asc'}
              >
                {column.label}
              </TableSortLabel>
            </CustomTableCell>
          );
        },
      },
    },
    {
      name: 'type',
      label: 'Type',
      options: {
        filter: true,
        sort: false,
        searchable: true,
        customHeadRender: function CustomHead({ index, ...column }, sortColumn) {
          return (
            <CustomTableCell key={index} onClick={() => sortColumn(index)}>
              <TableSortLabel
                active={column.sortDirection != null}
                direction={column.sortDirection || 'asc'}
              >
                {column.label}
              </TableSortLabel>
            </CustomTableCell>
          );
        },
        customBodyRender: function CustomBody(_, tableMeta) {
          return (
            <Tooltip title={tableMeta.rowData[1]}>
              <Chip
                label={tableMeta.rowData[1]}
                variant="outlined"
                icon={getCredentialsIcon(tableMeta.rowData[1])}
              />
            </Tooltip>
          );
        },
      },
    },
    {
      name: 'created_at',
      label: 'Creation Date',
      options: {
        filter: true,
        sort: false,
        searchable: true,
        sortDescFirst: true,
        customHeadRender: function CustomHead({ index, ...column }, sortColumn) {
          return (
            <CustomTableCell key={index} onClick={() => sortColumn(index)}>
              <TableSortLabel
                active={column.sortDirection != null}
                direction={column.sortDirection || 'asc'}
              >
                {column.label}
              </TableSortLabel>
            </CustomTableCell>
          );
        },
        customBodyRender: function CustomBody(value) {
          return <Moment format="LLLL">{value}</Moment>;
        },
      },
    },
    {
      name: 'updated_at',
      label: 'Updation Date',
      options: {
        filter: true,
        sort: false,
        searchable: true,
        sortDescFirst: true,
        customHeadRender: function CustomHead({ index, ...column }, sortColumn) {
          return (
            <CustomTableCell key={index} onClick={() => sortColumn(index)}>
              <TableSortLabel
                active={column.sortDirection != null}
                direction={column.sortDirection || 'asc'}
              >
                {column.label}
              </TableSortLabel>
            </CustomTableCell>
          );
        },
        customBodyRender: function CustomBody(value) {
          return <Moment format="LLLL">{value}</Moment>;
        },
      },
    },
    {
      name: 'actions',
      label: 'Actions',
      options: {
        filter: false,
        sort: false,
        searchable: false,
        customHeadRender: function CustomHead({ index, ...column }) {
          return <CustomTableCell key={index}>{column.label}</CustomTableCell>;
        },
        customBodyRender: (_, tableMeta) => {
          const rowData = (credentialsData?.credentials || [])[tableMeta.rowIndex];
          return (
            <ActionContainer>
              <Tooltip key={`delete_credential-${tableMeta.rowIndex}`} title="Delete Credential">
                <IconButton
                  aria-label="delete"
                  onClick={() => handleSubmit({ type: CON_OPS.DELETE, id: rowData['id'] })}
                  size="large"
                >
                  <DeleteIcon />
                </IconButton>
              </Tooltip>
            </ActionContainer>
          );
        },
      },
    },
  ];
  const options = {
    filter: false,
    rowsPerPageOptions: [10, 20, 25],
    filterType: 'textField',
    responsive: 'standard',
    print: false,
    search: false,
    viewColumns: false,
    download: false,
    selectToolbarPlacement: 'none',
    selectableRows: 'none',
    elevation: 0,
    draggableColumns: {
      enabled: true,
    },
  };

  // control the entire submit
  const handleSubmit = async ({ id, type }) => {
    updateProgress({ showProgress: true });

    try {
      if (type === CON_OPS.DELETE) {
        await deleteCredential(id).unwrap();
        notify({ message: `Credential deleted.`, event_type: EVENT_TYPES.SUCCESS });
      }

      if (type === CON_OPS.CREATE) {
        const data = {
          name: credentialName,
          type: credentialType,
          secret: formData,
        };
        await createCredential(data).unwrap();
        notify({ message: `"${credentialType}" created.`, event_type: EVENT_TYPES.SUCCESS });
      }

      if (type === CON_OPS.UPDATE) {
        const data = {
          id: id,
          name: credentialName,
          type: credentialType,
          secret: formData,
        };
        await updateCredential(data).unwrap();
        notify({ message: `"${credentialType}" updated.`, event_type: EVENT_TYPES.SUCCESS });
      }

      // Close modal if needed
      if (credModal.open) {
        setCredModal({
          open: false,
          data: null,
          actionType: null,
          id: null,
        });
      }
    } catch (error) {
      const errorMessage =
        type === CON_OPS.DELETE
          ? 'Failed to delete credentials.'
          : type === CON_OPS.CREATE
            ? 'Failed to create credentials.'
            : 'Failed to update credentials.';
      handleError(errorMessage);
    } finally {
      updateProgress({ showProgress: false });
    }
  };

  const [tableCols, updateCols] = useState(columns);

  const [columnVisibility, setColumnVisibility] = useState(() => {
    let showCols = updateVisibleColumns(colViews, width);
    // Initialize column visibility based on the original columns' visibility
    const initialVisibility = {};
    columns.forEach((col) => {
      initialVisibility[col.name] = showCols[col.name];
    });
    return initialVisibility;
  });

  const customInlineStyle = {
    marginBottom: '0.5rem',
    marginTop: '1rem',
  };

  if (isLoading) {
    return <LoadingScreen animatedIcon="AnimatedMeshery" message="Loading Credentials" />;
  }
  return (
    <div style={{ display: 'table', tableLayout: 'fixed', width: '100%' }}>
      <ToolWrapper style={customInlineStyle}>
        <div>
          {/* TODO: Uncomment this when schema spec is ready to support various credential */}
          {/* <Button
            type="submit"
            variant="contained"
            color="primary"
            size="large"
            onClick={(ev) => handleOpen(ev)(null, 'create', null)}
            style={{
              padding: '0.5rem',
              borderRadius: 5,
              marginRight: '2rem',
            }}
            data-cy="btnResetDatabase"
          >
            <AddIconCircleBorder style={{ width: '1.25rem' }} />
            <Typography
              style={{
                paddingLeft: '0.25rem',
                marginRight: '0.25rem',
              }}
            >
              Create
            </Typography>
          </Button> */}
        </div>
        <div>
          {/* <SearchBar
          onSearch={(value) => {

          } */}
          <CustomColumnVisibilityControl
            id="ref"
            columns={columns}
            customToolsProps={{ columnVisibility, setColumnVisibility }}
          />
        </div>
      </ToolWrapper>
      <ResponsiveDataTable
        columns={columns}
        data={credentialsData?.credentials || []}
        options={options}
        tableCols={tableCols}
        updateCols={updateCols}
        columnVisibility={columnVisibility}
      />

      <Modal
        open={credModal.open}
        formData={credModal.data}
        title="Credentials"
        handleClose={handleClose}
        onChange={_onChange}
        schema_array={schema_array}
        type={credentialType}
        schemaChangeHandler={schemaChangeHandler}
        handleSubmit={handleSubmit}
        payload={{ type: credModal.actionType, id: credModal.id }}
        submitBtnText="Save"
      />
    </div>
  );
};

const mapDispatchToProps = (dispatch) => ({
  updateProgress: bindActionCreators(updateProgress, dispatch),
});

const mapStateToProps = (state) => {
  return {
    connectionMetadataState: state.get('connectionMetadataState'),
  };
};

export default connect(mapStateToProps, mapDispatchToProps)(MesheryCredentialComponent);
