import React from 'react';
import { getStatus, timeAgo } from '../../../../utils/k8s-utils';
import { SINGLE_VIEW } from '../config';

import { Title } from '../../view';

import { TooltipWrappedConnectionChip } from '../../../connections/ConnectionChip';
import useKubernetesHook from '../../../hooks/useKubernetesHook';
import { DefaultTableCell, SortableTableCell } from '../sortable-table-cell';
import { CONNECTION_KINDS } from '../../../../utils/Enum';
import { getK8sContextFromClusterId } from '../../../../utils/multi-ctx';
import { FormatId } from '@/components/DataFormatter';

export const NamespaceTableConfig = (
  switchView,
  meshSyncResources,
  k8sConfig,
  connectionMetadataState,
  workloadType,
) => {
  const ping = useKubernetesHook();
  return {
    name: 'Namespace',
    colViews: [
      ['id', 'na'],
      ['metadata.name', 'xs'],
      ['apiVersion', 'na'],
      ['cluster_id', 'xs'],
      ['metadata.creationTimestamp', 'l'],
      ['status.attribute', 'm'],
    ],
    columns: [
      {
        name: 'id',
        label: 'ID',
        options: {
          display: false,
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
                handlePing={() => ping(context.name, context.server, context.connection_id)}
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
      {
        name: 'status.attribute',
        label: 'Status',
        options: {
          sort: false,
          customHeadRender: function CustomHead({ ...column }) {
            return <DefaultTableCell columnData={column} />;
          },
          customBodyRender: function CustomBody(val) {
            const phase = getStatus(val);
            return <>{phase}</>;
          },
        },
      },
    ],
  };
};
