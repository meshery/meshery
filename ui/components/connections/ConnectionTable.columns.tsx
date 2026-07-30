import React, { useMemo } from 'react';
import {
  Box,
  IconButton,
  Grid2,
  TableCell,
  InfoOutlinedIcon,
  MoreVertIcon,
  CustomTooltip,
  getRelativeTime,
  getFullFormattedTime,
} from '@sistent/sistent';
import { FormatId } from '../data-formatter';
import { iconMedium } from '../../css/icons.styles';
import { CONNECTION_KINDS } from '../../utils/Enum';
import { TooltipWrappedConnectionChip } from './ConnectionChip';
import { ConnectionStatusSelect } from './ConnectionStatusSelect';
import { DefaultTableCell, SortableTableCell } from './common';
import { getColumnValue } from '../../utils/utils';
import MultiSelectWrapper from '../multi-select-wrapper';
import CAN from '@/utils/can';
import { Keys } from '@meshery/schemas/permissions';
import { CustomTextTooltip } from '../meshery-mesh-interface/PatternService/CustomTextTooltip';
import { getFallbackImageBasedOnKind, normalizeStaticImagePath } from '@/utils/fallback';
import type { ConnectionTransitionMap } from './ConnectionTable.constants';
import type { EnvironmentOption, RowData } from './ConnectionTable.types';

/** Shared header info button for column tooltips (environments, status, timestamps). */
const ColumnInfoIcon = () => (
  <IconButton
    disableRipple={true}
    disableFocusRipple={true}
    onClick={(event) => {
      event.stopPropagation();
    }}
  >
    <InfoOutlinedIcon
      style={{
        cursor: 'pointer',
        height: 20,
        width: 20,
      }}
    />
  </IconButton>
);

/**
 * Relative time cell for connection timestamps. Guards empty values, unparsable
 * dates, and Go's zero-time sentinel (`0001-01-01T00:00:00Z`; no omitempty on
 * schemas v1beta3 CreatedAt/UpdatedAt).
 *
 * Same content as Sistent FormattedTime (relative text + full datetime tooltip),
 * but the tooltip target is an inline shrink-wrap so MUI anchors over the text
 * instead of the full table-cell width. FormattedTime uses a block-level div,
 * which makes the popup sit far from left-aligned values like "an hour ago".
 */
const renderTimestampCell = (value: unknown) => {
  if (value == null || value === '') {
    return <span>-</span>;
  }
  const dateStr = String(value);
  const parsed = value instanceof Date ? value : new Date(dateStr);
  if (Number.isNaN(parsed.getTime()) || parsed.getUTCFullYear() <= 1) {
    return <span>-</span>;
  }
  return (
    <CustomTooltip title={getFullFormattedTime(dateStr)} disableInteractive>
      <span data-testid="formatted-time" style={{ display: 'inline-block' }}>
        {getRelativeTime(dateStr)}
      </span>
    </CustomTooltip>
  );
};

type UseConnectionColumnsArgs = {
  url: string;
  envUrl: string;
  environmentOptions: Array<{ label: string; value: string }>;
  isEnvironmentsSuccess: boolean;
  updatingConnection: React.MutableRefObject<boolean>;
  handleDeleteConnection: (connectionId: string) => void | Promise<void>;
  handleEnvironmentSelect: (
    connectionId: string,
    connName: string,
    assignedEnvironments: EnvironmentOption[],
    selectedEnvironments: EnvironmentOption[],
    unSelectedEnvironments: EnvironmentOption[],
  ) => void | Promise<void>;
  handleStatusChange: (
    status: string,
    connectionId: string,
    connectionKind: string,
    connectionStatus: string,
  ) => void | Promise<void>;
  handleActionMenuOpen: (event: any, tableMeta: RowData) => void;
  ping: (name: string, server: string, id: string) => void;
  pingGrafana: (connectionID: string, name?: string) => void;
  pingPrometheus: (connectionID: string, name?: string) => void;
  // Per-kind connection state machine, keyed by connection kind. Sourced from
  // the connection definitions' `transitionMap` (see `_app.tsx`).
  transitionMapByKind: Record<string, ConnectionTransitionMap | undefined> | null;
};

export const useConnectionColumns = ({
  url,
  envUrl,
  environmentOptions,
  isEnvironmentsSuccess,
  updatingConnection,
  handleDeleteConnection,
  handleEnvironmentSelect,
  handleStatusChange,
  handleActionMenuOpen,
  ping,
  pingGrafana,
  pingPrometheus,
  transitionMapByKind,
}: UseConnectionColumnsArgs) => {
  return useMemo(() => {
    const nextColumns = [
      {
        name: 'id',
        label: 'ID',
        options: {
          display: false,
        },
      },
      {
        name: 'metadata.server_location',
        label: 'Server Location',
        options: {
          display: false,
        },
      },
      {
        // The wire metadata for meshery-kind connections uses camelCase
        // (`serverLocation`, see BuildMesheryConnectionPayload); the snake_case
        // sibling above is kept for older records.
        name: 'metadata.serverLocation',
        label: 'Server Location',
        options: {
          display: false,
        },
      },
      {
        name: 'metadata.server',
        label: 'Server',
        options: {
          display: false,
        },
      },
      {
        name: 'name',
        label: 'Name',
        options: {
          sort: true,
          sortThirdClickReset: true,
          customHeadRender: function CustomHead({ index, ...column }, sortColumn, columnMeta) {
            return (
              <SortableTableCell
                index={index}
                columnData={column}
                columnMeta={columnMeta}
                onSort={() => sortColumn(index)}
                icon={<ColumnInfoIcon />}
                tooltip={`The name of the connection, taken from the discovered infrastructure — for example the Kubernetes context name. Hover the name to see the server it points to. [Learn more](${url})`}
              />
            );
          },
          customBodyRender: (value, tableMeta) => {
            const server =
              getColumnValue(tableMeta.rowData, 'metadata.server', nextColumns) ||
              getColumnValue(tableMeta.rowData, 'metadata.serverLocation', nextColumns) ||
              getColumnValue(tableMeta.rowData, 'metadata.server_location', nextColumns);
            const name = getColumnValue(tableMeta.rowData, 'metadata.name', nextColumns);
            const kind = getColumnValue(tableMeta.rowData, 'kind', nextColumns);
            const connectionId = getColumnValue(tableMeta.rowData, 'id', nextColumns);
            const iconSrc = normalizeStaticImagePath(
              getColumnValue(tableMeta.rowData, 'kindLogo', nextColumns) ||
                getFallbackImageBasedOnKind(kind),
            );

            // Only attach handlePing for kinds that support a chip ping. An
            // always-defined no-op handler still makes the chip swallow row
            // clicks (stopPropagation) even when it cannot ping.
            let handlePing: (() => void) | undefined;
            if (kind === CONNECTION_KINDS.KUBERNETES) {
              handlePing = () =>
                ping(
                  getColumnValue(tableMeta.rowData, 'metadata.name', nextColumns),
                  getColumnValue(tableMeta.rowData, 'metadata.server', nextColumns),
                  connectionId,
                );
            } else if (kind === CONNECTION_KINDS.GRAFANA) {
              handlePing = () =>
                pingGrafana(connectionId, getColumnValue(tableMeta.rowData, 'name', nextColumns));
            } else if (kind === CONNECTION_KINDS.PROMETHEUS) {
              handlePing = () =>
                pingPrometheus(
                  connectionId,
                  getColumnValue(tableMeta.rowData, 'name', nextColumns),
                );
            }

            return (
              <>
                <TooltipWrappedConnectionChip
                  tooltip={server ? `Server: ${server}` : ''}
                  title={kind === CONNECTION_KINDS.KUBERNETES ? name : value || name || kind}
                  status={getColumnValue(tableMeta.rowData, 'status', nextColumns)}
                  onDelete={() => handleDeleteConnection(connectionId)}
                  handlePing={handlePing}
                  iconSrc={iconSrc}
                  width="12rem"
                />
                {kind === 'kubernetes' && (
                  <CustomTextTooltip
                    placement="top"
                    title="Learn more about connection status and how to [troubleshoot Kubernetes connections](https://docs.meshery.io/guides/troubleshooting/meshery-operator-meshsync)"
                  >
                    <div style={{ display: 'inline-block' }}>
                      <IconButton
                        color="default"
                        onClick={(event) => {
                          event.stopPropagation();
                          event.preventDefault();
                        }}
                      >
                        <InfoOutlinedIcon height={20} width={20} />
                      </IconButton>
                    </div>
                  </CustomTextTooltip>
                )}
              </>
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
              <DefaultTableCell
                columnData={column}
                icon={<ColumnInfoIcon />}
                tooltip={`Meshery Environments allow you to logically group related Connections and their associated Credentials. [Learn more](${envUrl})`}
              />
            );
          },
          customBodyRender: (value, tableMeta) => {
            const cleanedEnvs =
              value?.map((environment) => ({
                label: environment.name,
                value: environment.id,
              })) || [];

            return (
              isEnvironmentsSuccess && (
                <div onClick={(event) => event.stopPropagation()}>
                  <Grid2 size={{ xs: 12 }} style={{ height: '5rem', width: '15rem' }}>
                    <Grid2 size={{ xs: 12 }} style={{ marginTop: '2rem', cursor: 'pointer' }}>
                      <MultiSelectWrapper
                        updating={updatingConnection.current}
                        onChange={(selected, unselected) =>
                          handleEnvironmentSelect(
                            getColumnValue(tableMeta.rowData, 'id', nextColumns),
                            getColumnValue(tableMeta.rowData, 'name', nextColumns),
                            cleanedEnvs,
                            selected,
                            unselected,
                          )
                        }
                        options={environmentOptions}
                        value={cleanedEnvs}
                        placeholder={`Select or create an environment`}
                        noOptionsMessage={() =>
                          'No matching environments. Type to create a new one.'
                        }
                        isSelectAll={true}
                        menuPlacement={'bottom'}
                        disabled={
                          !CAN(
                            Keys.WorkspaceManagementAssignConnectionsToEnvironment.id,
                            Keys.WorkspaceManagementAssignConnectionsToEnvironment.function,
                          )
                        }
                      />
                    </Grid2>
                  </Grid2>
                </div>
              )
            );
          },
        },
      },
      {
        name: 'kind',
        label: 'Kind',
        options: {
          sort: true,
          sortThirdClickReset: true,
          customHeadRender: function CustomHead({ index, ...column }, sortColumn, columnMeta) {
            return (
              <SortableTableCell
                index={index}
                columnData={column}
                columnMeta={columnMeta}
                onSort={() => sortColumn(index)}
                icon={<ColumnInfoIcon />}
                tooltip={`The kind of infrastructure this connection points to — for example kubernetes, prometheus, grafana, or github. Kind determines which actions and lifecycle states are available. [Learn more](${url})`}
              />
            );
          },
        },
      },
      {
        name: 'type',
        label: 'Category',
        options: {
          sort: true,
          sortThirdClickReset: true,
          customHeadRender: function CustomHead({ index, ...column }, sortColumn, columnMeta) {
            return (
              <SortableTableCell
                index={index}
                columnData={column}
                columnMeta={columnMeta}
                onSort={() => sortColumn(index)}
                icon={<ColumnInfoIcon />}
                tooltip={`The broad classification of the connection: platform, telemetry, or collaboration. [Learn more](${url})`}
              />
            );
          },
        },
      },
      {
        // Connections arrive in the v1beta3 camelCase wire shape (see
        // server/models/connections type aliases): subType, createdAt,
        // updatedAt. Column names must match those row fields; the server's
        // snake_case sort columns are mapped in toServerSortOrder.
        name: 'subType',
        label: 'Sub Category',
        options: {
          sort: true,
          sortThirdClickReset: true,
          customHeadRender: function CustomHead({ index, ...column }, sortColumn, columnMeta) {
            return (
              <SortableTableCell
                index={index}
                columnData={column}
                columnMeta={columnMeta}
                onSort={() => sortColumn(index)}
                icon={<ColumnInfoIcon />}
                tooltip={`A finer classification within the category — for example cloud, identity, metrics, chat, git, or orchestration. [Learn more](${url})`}
              />
            );
          },
        },
      },
      {
        name: 'updatedAt',
        label: 'Updated At',
        options: {
          sort: true,
          sortThirdClickReset: true,
          display: false,
          customHeadRender: function CustomHead({ index, ...column }, sortColumn, columnMeta) {
            return (
              <SortableTableCell
                index={index}
                columnData={column}
                columnMeta={columnMeta}
                onSort={() => sortColumn(index)}
                icon={<ColumnInfoIcon />}
                tooltip="When this connection was last modified in Meshery, such as a status or metadata change. Values show relative time (for example, 2 hours ago). Hover the value for the full local date and time."
              />
            );
          },
          // Same timestamp treatment as Discovered At when enabled via View Columns.
          customBodyRender: renderTimestampCell,
        },
      },
      {
        name: 'createdAt',
        label: 'Discovered At',
        options: {
          sort: true,
          sortThirdClickReset: true,
          customHeadRender: function CustomHead({ index, ...column }, sortColumn, columnMeta) {
            return (
              <SortableTableCell
                index={index}
                columnData={column}
                columnMeta={columnMeta}
                onSort={() => sortColumn(index)}
                icon={<ColumnInfoIcon />}
                tooltip="When Meshery first recorded this connection through discovery or registration. This timestamp is set at creation and is not updated on later MeshSync events. Values show relative time (for example, 2 hours ago). Hover the value for the full local date and time."
              />
            );
          },
          // Relative time in-cell, full local datetime on hover (inline tooltip target).
          customBodyRender: renderTimestampCell,
        },
      },
      {
        name: 'ConnectionID',
        label: 'Connection ID',
        options: {
          sort: true,
          sortThirdClickReset: true,
          customHeadRender: function CustomHead({ index, ...column }, sortColumn, columnMeta) {
            return (
              <SortableTableCell
                index={index}
                columnData={column}
                columnMeta={columnMeta}
                onSort={() => sortColumn(index)}
                icon={<ColumnInfoIcon />}
                tooltip={`Meshery's unique identifier (UUID) for this connection. Use it to reference the connection from the API or mesheryctl; click the value to copy it. [Learn more](${url})`}
              />
            );
          },
          customBodyRender: (value, tableMeta) => {
            const connectionId = getColumnValue(tableMeta.rowData, 'id', nextColumns);
            return <FormatId id={connectionId} />;
          },
        },
      },
      {
        name: 'status',
        label: 'Status',
        options: {
          sort: true,
          sortThirdClickReset: true,
          customHeadRender: function CustomHead({ index, ...column }, sortColumn, columnMeta) {
            return (
              <SortableTableCell
                index={index}
                columnData={column}
                columnMeta={columnMeta}
                onSort={() => sortColumn(index)}
                icon={<ColumnInfoIcon />}
                tooltip={`Every connection can be in one of the states at any given point of time. Eg: Connected, Registered, Discovered, etc. It allow users more control over whether the discovered infrastructure is to be managed or not (registered for use or not). [Learn more](${url})`}
              />
            );
          },
          customBodyRender: function CustomBody(value, tableMeta) {
            const kind = getColumnValue(tableMeta.rowData, 'kind', nextColumns);
            const disabled =
              value === 'deleted'
                ? true
                : !CAN(
                    Keys.LifecycleManagementChangeConnectionState.id,
                    Keys.LifecycleManagementChangeConnectionState.function,
                  );

            return (
              <ConnectionStatusSelect
                status={value}
                transitionMap={transitionMapByKind?.[kind]}
                disabled={disabled}
                onChange={(nextStatus) =>
                  handleStatusChange(
                    nextStatus,
                    getColumnValue(tableMeta.rowData, 'id', nextColumns),
                    kind,
                    getColumnValue(tableMeta.rowData, 'status', nextColumns),
                  )
                }
              />
            );
          },
        },
      },
      {
        name: 'Actions',
        label: 'Actions',
        options: {
          filter: false,
          sort: false,
          searchable: false,
          customHeadRender: function CustomHead({ ...column }) {
            return (
              <TableCell>
                <b>{column.label}</b>
              </TableCell>
            );
          },
          customBodyRender: function CustomBody(_, tableMeta) {
            return (
              <Box sx={{ display: 'flex', justifyContent: 'flex-start' }}>
                {getColumnValue(tableMeta.rowData, 'kind', nextColumns) ===
                CONNECTION_KINDS.KUBERNETES ? (
                  <IconButton
                    aria-label="more"
                    id="long-button"
                    aria-haspopup="true"
                    onClick={(event) => handleActionMenuOpen(event, tableMeta)}
                  >
                    <MoreVertIcon style={iconMedium} />
                  </IconButton>
                ) : (
                  '-'
                )}
              </Box>
            );
          },
        },
      },
      {
        name: 'nextStatus',
        label: 'nextStatus',
        options: {
          display: false,
        },
      },
      {
        name: 'kindLogo',
        label: 'kindLogo',
        options: {
          display: false,
        },
      },
      {
        name: 'metadata.name',
        label: 'Name',
        options: {
          display: false,
        },
      },
    ];

    return nextColumns;
  }, [
    envUrl,
    environmentOptions,
    handleActionMenuOpen,
    handleDeleteConnection,
    handleEnvironmentSelect,
    handleStatusChange,
    isEnvironmentsSuccess,
    ping,
    pingGrafana,
    pingPrometheus,
    transitionMapByKind,
    updatingConnection,
    url,
  ]);
};
