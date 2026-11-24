import React, { useEffect, useRef, useState } from 'react';
import {
  Chip,
  CustomColumnVisibilityControl,
  DeleteIcon,
  IconButton,
  PROMPT_VARIANTS,
  ResponsiveDataTable,
  styled,
  TableCell,
  TableSortLabel,
  Tooltip,
  useHasPermission,
  useTheme,
} from '@sistent/sistent';
import { Keys } from '@meshery/schemas/permissions';
import { useGetUserCredentialsQuery } from '@meshery/schemas/mesheryApi';
import DefaultError from './general/error-404/index';
import Modal from './shared/Modal/Modal';
import _PromptComponent, { PromptRef } from './general/PromptComponent';
import { CON_OPS, CoreConnectionKinds } from '../utils/Enum';
import LoadingScreen from './shared/LoadingState/LoadingComponent';
import { useNotification, useNotificationHandlers } from '../utils/hooks/useNotification';
import { EVENT_TYPES } from '../lib/event-types';
import { updateVisibleColumns } from '../utils/responsive-column';
import { useWindowDimensions } from '../utils/dimension';
import { normalizeStaticImagePath } from '../utils/fallback';
import { ToolWrapper } from '@/assets/styles/general/tool.styles';
import {
  useCreateCredentialMutation,
  useDeleteCredentialMutation,
  useUpdateCredentialMutation,
} from '@/rtk-query/credentials';
import { useSelector } from 'react-redux';
import { updateProgress } from '@/store/slices/mesheryUi';
import type { RootState } from '../store';
import type { MUIDataTableColumn, MUIDataTableMeta } from '@sistent/mui-datatables';

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
  const canViewCredentials = useHasPermission(Keys.SecurityManagementViewCredentials);
  // The table pages on the server: without page/pageSize the endpoint returns
  // only its default first page (10), and credentials beyond it were unreachable.
  const [page, setPage] = useState(0);
  const [pageSize, setPageSize] = useState(10);
  const { data: credentialsData, isLoading } = useGetUserCredentialsQuery(
    { page, pageSize },
    { skip: !canViewCredentials },
  );

  // Deleting the last credential on the last page leaves that page empty, so
  // step back to the new last page once the server reports the smaller total.
  useEffect(() => {
    if (!credentialsData || credentialsData.page !== page) {
      return;
    }
    const lastPage = Math.max(0, Math.ceil((credentialsData.totalCount ?? 0) / pageSize) - 1);
    if (page > lastPage) {
      setPage(lastPage);
    }
  }, [credentialsData, page, pageSize]);
  const [createCredential] = useCreateCredentialMutation();
  const [updateCredential] = useUpdateCredentialMutation();
  const [deleteCredential] = useDeleteCredentialMutation();
  const { connectionMetadataState } = useSelector((state: RootState) => state.ui);

  const modalRef = useRef<PromptRef>(null);
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
  const { notifyApiError } = useNotificationHandlers();
  const { width } = useWindowDimensions();
  const theme = useTheme();

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

  /**
   * Surface a credential mutation error. Delegates to `notifyApiError` so the
   * structured MeshKit envelope (when present) renders the server-supplied
   * message and remediation guidance; otherwise falls back to the supplied
   * description string.
   */
  const handleError = (error: unknown, fallbackMessage: string): void => {
    updateProgress({ showProgress: false });
    notifyApiError(error, fallbackMessage);
  };

  const getCredentialsIcon = (type: string): React.ReactNode => {
    switch (type) {
      case 'prometheus':
        return <CredentialIcon src="/static/img/integrations/prometheus_logo_orange_circle.svg" />;
      case 'grafana':
        return <CredentialIcon src="/static/img/integrations/grafana_icon.svg" />;
      case 'kubernetes':
        return (
          <CredentialIcon
            src={
              normalizeStaticImagePath(
                connectionMetadataState?.[CoreConnectionKinds.kubernetes]?.icon,
              ) || undefined
            }
            alt=""
            aria-hidden="true"
          />
        );
      default:
        return null;
    }
  };

  const colViews: [string, string][] = [
    ['name', 'xs'],
    ['type', 'l'],
    ['createdAt', 'xl'],
    ['updatedAt', 'xl'],
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
      // Named after the camelCase wire fields so the cell gets a value;
      // Sistent's ResponsiveDataTable formats createdAt/updatedAt itself.
      name: 'createdAt',
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
      },
    },
    {
      name: 'updatedAt',
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
          const rowData = credentials[tableMeta.rowIndex] as
            | { id: string; name?: string; type?: string }
            | undefined;
          return (
            <ActionContainer>
              <Tooltip key={`delete_credential-${tableMeta.rowIndex}`} title="Delete Credential">
                <IconButton
                  aria-label="delete"
                  onClick={(e) => handleDeleteCredentialConfirm(e, rowData)}
                  size="large"
                  permissionKey={Keys.SecurityManagementDeleteCredential}
                >
                  <DeleteIcon fill={theme?.palette?.icon?.default} />
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
    serverSide: true,
    count: credentialsData?.totalCount ?? 0,
    page,
    rowsPerPage: pageSize,
    onTableChange: (action: string, tableState: { page: number; rowsPerPage: number }) => {
      switch (action) {
        case 'changePage':
          setPage(tableState.page);
          break;
        case 'changeRowsPerPage':
          setPageSize(tableState.rowsPerPage);
          setPage(0);
          break;
        default:
          break;
      }
    },
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

  const handleDeleteCredentialConfirm = async (
    e: React.MouseEvent,
    rowData?: { id: string; name?: string; type?: string },
  ): Promise<void> => {
    e.stopPropagation();
    if (!rowData?.id || !modalRef.current) {
      return;
    }

    const credentialName = rowData.name || rowData.type || 'this credential';

    const response = await modalRef.current.show({
      title: `Delete Credential?`,
      subtitle: `Are you sure you want to delete "${credentialName}"? (This action is irreversible and may break dependent connections)`,
      primaryOption: 'DELETE',
      variant: PROMPT_VARIANTS.DANGER,
    });

    if (response === 'DELETE') {
      handleSubmit({ type: CON_OPS.DELETE, id: rowData.id });
    }
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
      handleError(error, errorMessage);
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

  if (!canViewCredentials) {
    return <DefaultError permissionKey={Keys.SecurityManagementViewCredentials} />;
  }

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
        // A prop, not an option: ResponsiveDataTable overwrites
        // options.rowsPerPageOptions with this prop's default of 10/25/50/100,
        // and the credentials endpoint serves at most 25 per page.
        rowsPerPageOptions={[10, 20, 25]}
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
      <_PromptComponent ref={modalRef} />
    </div>
  );
};

export default MesheryCredentialComponent;
