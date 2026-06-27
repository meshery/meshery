import React, { useMemo } from 'react';
import {
  CustomTooltip,
  MenuItem,
  Box,
  IconButton,
  Grid2,
  FormControl,
  TableCell,
  InfoOutlinedIcon,
} from '@sistent/sistent';
import { ConnectionStyledSelect } from './styles';
import { FormatId } from '../data-formatter';
import { MoreVertIcon } from '@sistent/sistent';
import { iconMedium } from '../../css/icons.styles';
import { CONNECTION_KINDS } from '../../utils/Enum';
import { ConnectionStateChip, TooltipWrappedConnectionChip } from './ConnectionChip';
import { DefaultTableCell, SortableTableCell } from './common';
import { getColumnValue } from '../../utils/utils';
import MultiSelectWrapper from '../multi-select-wrapper';
import CAN from '@/utils/can';
import { keys } from '@/utils/permission_constants';
import { CustomTextTooltip } from '../meshery-mesh-interface/PatternService/CustomTextTooltip';
import { formatDate } from '../data-formatter';
import { getFallbackImageBasedOnKind, normalizeStaticImagePath } from '@/utils/fallback';
import { getNextStates } from './ConnectionTable.constants';
import type { ConnectionTransitionMap } from './ConnectionTable.constants';
import type { EnvironmentOption, RowData } from './ConnectionTable.types';

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
    event: any,
    connectionId: string,
    connectionKind: string,
    connectionStatus: string,
  ) => void | Promise<void>;
  handleActionMenuOpen: (event: any, tableMeta: RowData) => void;
  ping: (name: string, server: string, id: string) => void;
  pingGrafana: (connectionID: string, name?: string) => void;
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
                icon={null}
                tooltip=""
              />
            );
          },
          customBodyRender: (value, tableMeta) => {
            const server =
              getColumnValue(tableMeta.rowData, 'metadata.server', nextColumns) ||
              getColumnValue(tableMeta.rowData, 'metadata.server_location', nextColumns);
            const name = getColumnValue(tableMeta.rowData, 'metadata.name', nextColumns);
            const kind = getColumnValue(tableMeta.rowData, 'kind', nextColumns);
            const iconSrc = normalizeStaticImagePath(
              getColumnValue(tableMeta.rowData, 'kindLogo', nextColumns) ||
                getFallbackImageBasedOnKind(kind),
            );

            return (
              <>
                <TooltipWrappedConnectionChip
                  tooltip={server ? `Server: ${server}` : ''}
                  title={kind === CONNECTION_KINDS.KUBERNETES ? name : value || name || kind}
                  status={getColumnValue(tableMeta.rowData, 'status', nextColumns)}
                  onDelete={() =>
                    handleDeleteConnection(getColumnValue(tableMeta.rowData, 'id', nextColumns))
                  }
                  handlePing={() => {
                    const rowKind = getColumnValue(tableMeta.rowData, 'kind', nextColumns);
                    if (rowKind === CONNECTION_KINDS.KUBERNETES) {
                      ping(
                        getColumnValue(tableMeta.rowData, 'metadata.name', nextColumns),
                        getColumnValue(tableMeta.rowData, 'metadata.server', nextColumns),
                        getColumnValue(tableMeta.rowData, 'id', nextColumns),
                      );
                    } else if (rowKind === CONNECTION_KINDS.GRAFANA) {
                      pingGrafana(
                        getColumnValue(tableMeta.rowData, 'id', nextColumns),
                        getColumnValue(tableMeta.rowData, 'name', nextColumns),
                      );
                    }
                  }}
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
                icon={
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
                }
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
                        placeholder={`Assigned Environments`}
                        isSelectAll={true}
                        menuPlacement={'bottom'}
                        disabled={
                          !CAN(
                            keys.ASSIGN_CONNECTIONS_TO_ENVIRONMENT.action,
                            keys.ASSIGN_CONNECTIONS_TO_ENVIRONMENT.subject,
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
                icon={null}
                tooltip=""
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
                icon={null}
                tooltip=""
              />
            );
          },
        },
      },
      {
        name: 'sub_type',
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
                icon={null}
                tooltip=""
              />
            );
          },
        },
      },
      {
        name: 'updated_at',
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
                icon={null}
                tooltip=""
              />
            );
          },
        },
      },
      {
        name: 'created_at',
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
                icon={null}
                tooltip=""
              />
            );
          },
          customBodyRender: function CustomBody(value) {
            const renderValue = formatDate(value);
            return (
              <CustomTooltip title={renderValue} placement="top" arrow interactive>
                <span>{renderValue}</span>
              </CustomTooltip>
            );
          },
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
                icon={null}
                tooltip=""
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
                icon={
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
                }
                tooltip={`Every connection can be in one of the states at any given point of time. Eg: Connected, Registered, Discovered, etc. It allow users more control over whether the discovered infrastructure is to be managed or not (registered for use or not). [Learn more](${url})`}
              />
            );
          },
          customBodyRender: function CustomBody(value, tableMeta) {
            const currentStatus = value;
            const kind = getColumnValue(tableMeta.rowData, 'kind', nextColumns);

            const nextStatus = getNextStates(transitionMapByKind?.[kind], currentStatus);
            nextStatus.push(currentStatus);

            const disabled =
              value === 'deleted'
                ? true
                : !CAN(keys.CHANGE_CONNECTION_STATE.action, keys.CHANGE_CONNECTION_STATE.subject);

            return (
              <FormControl>
                <ConnectionStyledSelect
                  labelId="connection-status-select-label"
                  id="connection-status-select"
                  disabled={disabled}
                  value={value}
                  defaultValue={value}
                  onClick={(event) => event.stopPropagation()}
                  onChange={(event) =>
                    handleStatusChange(
                      event,
                      getColumnValue(tableMeta.rowData, 'id', nextColumns),
                      getColumnValue(tableMeta.rowData, 'kind', nextColumns),
                      getColumnValue(tableMeta.rowData, 'status', nextColumns),
                    )
                  }
                  disableUnderline
                  MenuProps={{
                    anchorOrigin: {
                      vertical: 'bottom',
                      horizontal: 'left',
                    },
                    transformOrigin: {
                      vertical: 'top',
                      horizontal: 'left',
                    },
                    getContentAnchorEl: null,
                    MenuListProps: { disablePadding: true },
                    PaperProps: { square: true },
                  }}
                >
                  {nextStatus.length === 1 && (
                    <MenuItem disabled>No transitions Available</MenuItem>
                  )}
                  {nextStatus.map((status) => (
                    <MenuItem
                      disabled={status === value}
                      style={{
                        padding: 0,
                        display: status === value ? 'none' : 'flex',
                        justifyContent: 'center',
                      }}
                      value={status}
                      key={status}
                    >
                      <ConnectionStateChip status={status} actionable={status !== value} />
                    </MenuItem>
                  ))}
                </ConnectionStyledSelect>
              </FormControl>
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
              <Box sx={{ display: 'flex', justifyContent: 'center' }}>
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
    transitionMapByKind,
    updatingConnection,
    url,
  ]);
};
