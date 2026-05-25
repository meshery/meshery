import React from 'react';
import { timeAgo } from '../../../../utils/k8s-utils';
import { SINGLE_VIEW } from '../config';
import { Title } from '../../view';
import { TooltipWrappedConnectionChip } from '../../../connections/ConnectionChip';
import { DefaultTableCell, SortableTableCell } from '../sortable-table-cell';
import { CONNECTION_KINDS } from '../../../../utils/Enum';
import { getK8sContextFromClusterId } from '../../../../utils/multi-ctx';
import { FormatId } from '@/components/data-formatter';

export const buildRuntimeClassColumns = ({
  switchView,
  meshSyncResources,
  k8sConfig,
  connectionMetadataState,
  workloadType,
  ping,
}) => ({
  name: 'RuntimeClass',
  colViews: [
    ['id', 'na'],
    ['metadata.name', 'xs'],
    ['apiVersion', 'na'],
    ['metadata.namespace', 'm'],
    ['cluster_id', 'xs'],
    ['metadata.creationTimestamp', 'l'],
  ],
  columns: [
    {
      name: 'id',
      label: 'ID',
      options: {
        customHeadRender: function CustomHead({ ...column }) {
          return <DefaultTableCell columnData={column} />;
        },
        customBodyRender: (value) => <FormatId id={value} />,
      },
    },
    {
      name: 'metadata.name',
      label: 'Name',
      options: {
        sort: false,
        customHeadRender: function CustomHead({ ...column }) {
          return <DefaultTableCell columnData={column} />;
        },
        customBodyRender: function CustomBody(value, tableMeta) {
          return (
            <Title
              onClick={() => switchView(SINGLE_VIEW, meshSyncResources[tableMeta.rowIndex])}
              value={value}
              kind={workloadType}
            />
          );
        },
      },
    },
    {
      name: 'apiVersion',
      label: 'API version',
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
            />
          );
        },
      },
    },
    {
      name: 'metadata.namespace',
      label: 'Namespace',
      options: {
        sort: false,
        customHeadRender: function CustomHead({ ...column }) {
          return <DefaultTableCell columnData={column} />;
        },
      },
    },
    {
      name: 'cluster_id',
      label: 'Cluster',
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
            />
          );
        },
        customBodyRender: function CustomBody(val) {
          let context = getK8sContextFromClusterId(val, k8sConfig);
          return (
            <TooltipWrappedConnectionChip
              title={context.name}
              iconSrc={
                connectionMetadataState
                  ? connectionMetadataState[CONNECTION_KINDS.KUBERNETES]?.icon
                  : ''
              }
              handlePing={() => ping(context.name, context.server, context.connectionId)}
            />
          );
        },
      },
    },
    {
      name: 'metadata.creationTimestamp',
      label: 'Age',
      options: {
        sort: false,
        customHeadRender: function CustomHead({ ...column }) {
          return <DefaultTableCell columnData={column} />;
        },
        customBodyRender: function CustomBody(value) {
          let time = timeAgo(value);
          return <>{time}</>;
        },
      },
    },
  ],
});
