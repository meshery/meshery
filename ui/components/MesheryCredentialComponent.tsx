import React, { useState } from 'react';
import {
  Chip,
  IconButton,
  Tooltip,
  TableCell,
  TableSortLabel,
  styled,
  ResponsiveDataTable,
  CustomColumnVisibilityControl,
} from '@sistent/sistent';
import Modal from './Modal';
import { CONNECTION_KINDS, CON_OPS } from '../utils/Enum';
import DeleteIcon from '@mui/icons-material/Delete';
import Moment from 'react-moment';
import LoadingScreen from './LoadingComponents/LoadingComponent';
import { useNotification } from '../utils/hooks/useNotification';
import { EVENT_TYPES } from '../lib/event-types';
import { updateVisibleColumns } from '../utils/responsive-column';
import { useWindowDimensions } from '../utils/dimension';
import { ToolWrapper } from '@/assets/styles/general/tool.styles';
import {
  useCreateCredentialMutation,
  useDeleteCredentialMutation,
  useGetCredentialsQuery,
  useUpdateCredentialMutation,
} from '@/rtk-query/credentials';
import { useSelector } from 'react-redux';
import { updateProgress } from '@/store/slices/mesheryUi';
import type { RootState } from '@/store/store';
import type { MUIDataTableColumn, MUIDataTableMeta } from 'mui-datatables';

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

type CredentialType = 'prometheus' | 'grafana' | 'kubernetes';
type ActionType = 'create' | 'update' | 'delete' | null;

interface CredentialModalState {
  open: boolean;
  data: Record<string, unknown> | null;
  actionType: ActionType;
  id: string | null;
}

interface CredentialFormData {
  credentialName?: string;
  [key: string]: unknown;
}

interface HandleSubmitPayload {
  type: string;
  id?: string;
}

interface ColumnMeta {
  index: number;
  label: string;
  sortDirection?: 'asc' | 'desc' | null;
}

const MesheryCredentialComponent: React.FC = () => {
  const { data: credentialsData, isLoading } = useGetCredentialsQuery();
  const [createCredential] = useCreateCredentialMutation();
  const [updateCredential] = useUpdateCredentialMutation();
  const [deleteCredential] = useDeleteCredentialMutation();
  const { connectionMetadataState } = useSelector((state: RootState) => state.ui);

  const [formData, setFormData] = useState<CredentialFormData>({});
  const [credModal, setCredModal] = useState<CredentialModalState>({
    open: false,
    data: null,
    actionType: null,
    id: null,
  });
  const [credentialType, setCredentialType] = useState<CredentialType>(
    schema_array[0] as CredentialType,
  );
  const [credentialName, setCredentialName] = useState<string | null>(null);
  const { notify } = useNotification();
  const { width } = useWindowDimensions();

  const schemaChangeHandler = (type: CredentialType): void => {
    setCredentialType(type);
    setCredModal((prev) => ({
      ...prev,
      open: true,
      data: null,
    }));
  };

  const _onChange = (formData: CredentialFormData): void => {
    setCredentialName(formData?.credentialName || null);
    setFormData(formData);
  };

  const handleClose = (ev: React.MouseEvent): void => {
    ev.stopPropagation();
    setCredModal({
      open: false,
      data: null,
      actionType: null,
      id: null,
    });
  };

  const handleError = (error_msg: string): void => {
    updateProgress({ showProgress: false });
    notify({
      message: `${error_msg}`,
      event_type: EVENT_TYPES.ERROR,
      details: error_msg.toString(),
    });
  };

  const getCredentialsIcon = (type: string): React.ReactNode => {
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

  const colViews: [string, string][] = [
    ['name', 'xs'],
    ['type', 'l'],
    ['created_at', 'xl'],
    ['updated_at', 'xl'],
    ['actions', 'xs'],
  ];

  const columns: MUIDataTableColumn[] = [
    {
      name: 'name',
      label: 'Name',
      options: {
        filter: true,
        sort: false,
        searchable: true,
        customHeadRender: function CustomHead(
          columnMeta: ColumnMeta,
          sortColumn: (_idx: number) => void,
        ) {
          return (
            <CustomTableCell key={columnMeta.index} onClick={() => sortColumn(columnMeta.index)}>
              <TableSortLabel
                active={columnMeta.sortDirection != null}
                direction={columnMeta.sortDirection || 'asc'}
              >
                {columnMeta.label}
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
        customHeadRender: function CustomHead(
          columnMeta: ColumnMeta,
          sortColumn: (_idx: number) => void,
        ) {
          return (
            <CustomTableCell key={columnMeta.index} onClick={() => sortColumn(columnMeta.index)}>
              <TableSortLabel
                active={columnMeta.sortDirection != null}
                direction={columnMeta.sortDirection || 'asc'}
              >
                {columnMeta.label}
              </TableSortLabel>
            </CustomTableCell>
          );
        },
        customBodyRender: function CustomBody(_: unknown, tableMeta: MUIDataTableMeta) {
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
        customHeadRender: function CustomHead(
          columnMeta: ColumnMeta,
          sortColumn: (_idx: number) => void,
        ) {
          return (
            <CustomTableCell key={columnMeta.index} onClick={() => sortColumn(columnMeta.index)}>
              <TableSortLabel
                active={columnMeta.sortDirection != null}
                direction={columnMeta.sortDirection || 'asc'}
              >
                {columnMeta.label}
              </TableSortLabel>
            </CustomTableCell>
          );
        },
        customBodyRender: function CustomBody(value: string) {
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
        customHeadRender: function CustomHead(
          columnMeta: ColumnMeta,
          sortColumn: (_idx: number) => void,
        ) {
          return (
            <CustomTableCell key={columnMeta.index} onClick={() => sortColumn(columnMeta.index)}>
              <TableSortLabel
                active={columnMeta.sortDirection != null}
                direction={columnMeta.sortDirection || 'asc'}
              >
                {columnMeta.label}
              </TableSortLabel>
            </CustomTableCell>
          );
        },
        customBodyRender: function CustomBody(value: string) {
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
        customHeadRender: function CustomHead(columnMeta: ColumnMeta) {
          return <CustomTableCell key={columnMeta.index}>{columnMeta.label}</CustomTableCell>;
        },
        customBodyRender: (_: unknown, tableMeta: MUIDataTableMeta) => {
          const credentials = credentialsData?.credentials || [];
          const rowData = credentials[tableMeta.rowIndex] as { id: string } | undefined;
          return (
            <ActionContainer>
              <Tooltip key={`delete_credential-${tableMeta.rowIndex}`} title="Delete Credential">
                <IconButton
                  aria-label="delete"
                  onClick={() =>
                    handleSubmit({ type: CON_OPS.DELETE, id: rowData?.id || undefined })
                  }
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
    filterType: 'textField' as const,
    responsive: 'standard' as const,
    print: false,
    search: false,
    viewColumns: false,
    download: false,
    selectToolbarPlacement: 'none' as const,
    selectableRows: 'none' as const,
    elevation: 0,
    draggableColumns: {
      enabled: true,
    },
  };

  // control the entire submit
  const handleSubmit = async ({ id, type }: HandleSubmitPayload): Promise<void> => {
    updateProgress({ showProgress: true });

    try {
      if (type === CON_OPS.DELETE && id) {
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

      if (type === CON_OPS.UPDATE && id) {
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

  const [columnVisibility, setColumnVisibility] = useState<Record<string, boolean>>(() => {
    const showCols = updateVisibleColumns(colViews, width);
    // Initialize column visibility based on the original columns' visibility
    const initialVisibility: Record<string, boolean> = {};
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
        </div>
        <div>
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

export default MesheryCredentialComponent;
