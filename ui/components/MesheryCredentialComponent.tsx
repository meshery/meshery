import React, { useEffect, useRef, useState } from 'react';
import {
  AddIcon,
  Button,
  Chip,
  CustomColumnVisibilityControl,
  DeleteIcon,
  EditIcon,
  IconButton,
  LockIcon,
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
import {
  buildCredentialFormSchema,
  cleanEmptyPasswordFields,
  getCredentialTypes,
  getSchema,
  getUiSchema,
  sanitizeInitialData,
} from './meshery-mesh-interface/PatternService/helper';
import DefaultError from './general/error-404/index';
import Modal from './shared/Modal/Modal';
import _PromptComponent, { PromptRef } from './general/PromptComponent';
import { CoreConnectionKinds } from '../utils/Enum';
import LoadingScreen from './shared/LoadingState/LoadingComponent';
import { useNotification, useNotificationHandlers } from '../utils/hooks/useNotification';
import { EVENT_TYPES } from '../lib/event-types';
import { updateVisibleColumns } from '../utils/responsive-column';
import { useWindowDimensions } from '../utils/dimension';
import { normalizeStaticImagePath } from '../utils/fallback';
import { resolveCredentialPayload } from '../utils/credentialSecret';
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

const credentialTypes = getCredentialTypes();
const DEFAULT_CREDENTIAL_TYPE = credentialTypes.includes('prometheus')
  ? 'prometheus'
  : credentialTypes[0] || 'kubernetes';

type ActionType = 'create' | 'update' | 'delete' | null;

interface CredentialModalState {
  open: boolean;
  data: Record<string, unknown> | null;
  actionType: ActionType;
  id: string | null;
  originalType: string | null;
  openId: number;
}

interface ColumnMeta {
  index: number;
  label: string;
  sortDirection?: 'asc' | 'desc' | null;
}

export { buildCredentialFormSchema, cleanEmptyPasswordFields, sanitizeInitialData };

const renderSortableHead = (columnMeta: ColumnMeta, sortColumn: (_idx: number) => void) => (
  <CustomTableCell key={columnMeta.index} onClick={() => sortColumn(columnMeta.index)}>
    <TableSortLabel
      active={columnMeta.sortDirection != null}
      direction={columnMeta.sortDirection || 'asc'}
    >
      {columnMeta.label}
    </TableSortLabel>
  </CustomTableCell>
);

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
  const modalOpenIdRef = useRef(0);
  const [formData, setFormData] = useState<Record<string, any>>({});
  const [credModal, setCredModal] = useState<CredentialModalState>({
    open: false,
    data: null,
    actionType: null,
    id: null,
    originalType: null,
    openId: 0,
  });
  const [credentialType, setCredentialType] = useState<string>(DEFAULT_CREDENTIAL_TYPE);
  const [credentialName, setCredentialName] = useState<string | null>(null);
  const { notify } = useNotification();
  const { notifyApiError } = useNotificationHandlers();
  const { width } = useWindowDimensions();
  const theme = useTheme();

  const baseSchema = credentialType ? getSchema(credentialType) : undefined;
  const currentUiSchema = credentialType ? getUiSchema(credentialType) : undefined;
  const isSameTypeEdit =
    credModal.actionType === 'update' &&
    Boolean(credModal.originalType) &&
    credentialType.toLowerCase() === credModal.originalType!.toLowerCase();
  const currentSchema = buildCredentialFormSchema(
    baseSchema as Record<string, unknown> | undefined,
    currentUiSchema as Record<string, unknown> | undefined,
    isSameTypeEdit,
  );

  const handleCreateCredential = (ev: React.MouseEvent): void => {
    ev.stopPropagation();
    const openId = ++modalOpenIdRef.current;
    setCredentialType(DEFAULT_CREDENTIAL_TYPE);
    setCredentialName(null);
    setFormData({});
    setCredModal({
      open: true,
      data: {},
      actionType: 'create',
      id: null,
      originalType: null,
      openId,
    });
  };

  const handleEditCredential = (
    e: React.MouseEvent,
    rowData?: { id: string; name?: string; type?: string; secret?: Record<string, any> },
  ): void => {
    e.stopPropagation();
    if (!rowData?.id) return;

    const openId = ++modalOpenIdRef.current;
    const credType = rowData.type || DEFAULT_CREDENTIAL_TYPE;
    setCredentialType(credType);
    setCredentialName(rowData.name || null);

    const schema = getSchema(credType) as Record<string, any> | undefined;
    const uiSchema = getUiSchema(credType) as Record<string, any> | undefined;
    const resolvedSecret = resolveCredentialPayload(rowData.secret);
    const secretObj =
      resolvedSecret && typeof resolvedSecret === 'object' && !Array.isArray(resolvedSecret)
        ? (resolvedSecret as Record<string, unknown>)
        : {};
    const sanitizedSecret = sanitizeInitialData(
      secretObj,
      schema?.properties?.secret,
      uiSchema?.secret,
    );
    const initialFormData = { name: rowData.name || '', secret: sanitizedSecret };

    setFormData(initialFormData);
    setCredModal({
      open: true,
      data: initialFormData,
      actionType: 'update',
      id: rowData.id,
      originalType: credType,
      openId,
    });
  };

  const schemaChangeHandler = (type: string): void => {
    setCredentialType(type);
    const nextData =
      credModal.actionType === 'update'
        ? {
            name:
              (typeof formData?.name === 'string' ? formData.name : undefined) ??
              (typeof credModal.data?.name === 'string' ? credModal.data.name : undefined),
            secret: {},
          }
        : {};
    setFormData(nextData);
    setCredModal((prev) => ({
      ...prev,
      open: true,
      data: prev.actionType === 'update' ? nextData : {},
    }));
  };

  const _onChange = (data: Record<string, any>): void => {
    if (data?.name !== undefined) {
      setCredentialName(data.name || null);
    } else if (data?.credentialName !== undefined) {
      setCredentialName(data.credentialName || null);
    }
    setFormData(data || {});
  };

  const handleClose = (ev?: React.MouseEvent): void => {
    ev?.stopPropagation?.();
    setCredModal((prev) => ({
      ...prev,
      open: false,
      data: null,
      actionType: null,
      id: null,
      originalType: null,
    }));
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
    const lower = type?.toLowerCase();
    if (lower === 'prometheus') {
      return (
        <CredentialIcon
          src="/static/img/integrations/prometheus_logo_orange_circle.svg"
          alt="Prometheus"
        />
      );
    }
    if (lower === 'grafana') {
      return <CredentialIcon src="/static/img/integrations/grafana_icon.svg" alt="Grafana" />;
    }
    if (lower === 'kubernetes') {
      return (
        <CredentialIcon
          src={
            normalizeStaticImagePath(
              connectionMetadataState?.[CoreConnectionKinds.kubernetes]?.icon,
            ) || '/static/img/integrations/kubernetes.svg'
          }
          alt="Kubernetes"
          aria-hidden="true"
        />
      );
    }
    if (lower === 'github' || lower === 'github-oauth') {
      return <CredentialIcon src="/static/img/extensions/github.svg" alt="GitHub" />;
    }
    const metaIcon = connectionMetadataState?.[type]?.icon;
    if (metaIcon) {
      return (
        <CredentialIcon
          src={normalizeStaticImagePath(metaIcon) || undefined}
          alt={type}
          aria-hidden="true"
        />
      );
    }
    return (
      <LockIcon style={{ width: '20px', height: '20px', fill: theme?.palette?.icon?.default }} />
    );
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
        customHeadRender: renderSortableHead,
      },
    },
    {
      name: 'type',
      label: 'Type',
      options: {
        filter: true,
        sort: false,
        searchable: true,
        customHeadRender: renderSortableHead,
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
        customHeadRender: renderSortableHead,
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
        customHeadRender: renderSortableHead,
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
            | { id: string; name?: string; type?: string; secret?: Record<string, any> }
            | undefined;
          return (
            <ActionContainer>
              <Tooltip key={`edit_credential-${tableMeta.rowIndex}`} title="Edit Credential">
                <IconButton
                  aria-label="edit"
                  onClick={(e) => handleEditCredential(e, rowData)}
                  size="large"
                  permissionKey={Keys.SecurityManagementEditCredential}
                >
                  <EditIcon fill={theme?.palette?.icon?.default} />
                </IconButton>
              </Tooltip>
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
      if (action === 'changePage') setPage(tableState.page);
      if (action === 'changeRowsPerPage') {
        setPageSize(tableState.rowsPerPage);
        setPage(0);
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
    draggableColumns: { enabled: true },
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
      updateProgress({ showProgress: true });
      try {
        await deleteCredential(rowData.id).unwrap();
        notify({ message: `Credential deleted.`, event_type: EVENT_TYPES.SUCCESS });
      } catch (error) {
        handleError(error, 'Failed to delete credentials.');
      } finally {
        updateProgress({ showProgress: false });
      }
    }
  };

  const handleModalSubmit = async (submittedFormData?: Record<string, any>): Promise<void> => {
    const dataToSubmit = submittedFormData || formData;
    const { actionType, id, openId: submitOpenId } = credModal;

    updateProgress({ showProgress: true });
    try {
      const name = dataToSubmit?.name ?? credentialName ?? '';
      const secret = dataToSubmit?.secret ?? {};

      if (actionType === 'create') {
        await createCredential({ name, type: credentialType, secret }).unwrap();
        notify({ message: `"${credentialType}" created.`, event_type: EVENT_TYPES.SUCCESS });
      }

      if (actionType === 'update' && id) {
        const schema = getSchema(credentialType) as Record<string, any> | undefined;
        const uiSchema = getUiSchema(credentialType) as Record<string, any> | undefined;

        // Honor writeOnly / password widgets only when keeping the original credential type.
        const submittedSecret = isSameTypeEdit
          ? cleanEmptyPasswordFields(secret, schema?.properties?.secret, uiSchema?.secret)
          : secret;

        await updateCredential({
          id,
          name,
          type: credentialType,
          secret: submittedSecret,
        }).unwrap();
        notify({ message: `"${credentialType}" updated.`, event_type: EVENT_TYPES.SUCCESS });
      }

      setCredModal((prev) =>
        prev.openId === submitOpenId
          ? {
              open: false,
              data: null,
              actionType: null,
              id: null,
              originalType: null,
              openId: prev.openId,
            }
          : prev,
      );
    } catch (error) {
      const errorMessage =
        actionType === 'create' ? 'Failed to create credentials.' : 'Failed to update credentials.';
      handleError(error, errorMessage);
    } finally {
      updateProgress({ showProgress: false });
    }
  };

  const [tableCols, updateCols] = useState(columns);

  const [columnVisibility, setColumnVisibility] = useState<Record<string, boolean>>(() => {
    const showCols = updateVisibleColumns(colViews, width);
    return Object.fromEntries(columns.map((col) => [col.name, showCols[col.name]]));
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
          <Button
            type="submit"
            variant="contained"
            color="primary"
            size="large"
            onClick={handleCreateCredential}
            style={{
              padding: '0.5rem',
              borderRadius: 5,
              marginRight: '2rem',
            }}
            permissionKey={Keys.SecurityManagementCreateCredential}
            data-cy="btnCreateCredential"
          >
            <AddIcon style={{ width: '1.25rem', height: '1.25rem' }} />
            <span style={{ paddingLeft: '0.25rem' }}>Create</span>
          </Button>
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
        key={credModal.openId}
        open={credModal.open}
        schema={currentSchema}
        uiSchema={currentUiSchema}
        formData={credModal.data}
        initialData={credModal.data}
        title={credModal.actionType === 'update' ? 'Edit Credential' : 'Create Credential'}
        handleClose={handleClose}
        onChange={_onChange}
        schema_array={credentialTypes}
        type={credentialType}
        schemaChangeHandler={schemaChangeHandler}
        handleSubmit={handleModalSubmit}
        submitBtnText="Save"
        submitDisabled={!currentSchema}
        closeOnSubmit={false}
      />
      <_PromptComponent ref={modalRef} />
    </div>
  );
};

export default MesheryCredentialComponent;
