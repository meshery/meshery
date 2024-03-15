import React from 'react';
import { getResourceStr, resourceParsers, timeAgo } from '../../../../utils/k8s-utils';
import { getK8sContextFromClusterId } from '../../../../utils/multi-ctx';
import { SINGLE_VIEW } from '../config';

import { Title } from '../../view';

import { TootltipWrappedConnectionChip } from '../../../connections/ConnectionChip';
import useKubernetesHook from '../../../hooks/useKubernetesHook';
import { DefaultTableCell, SortableTableCell } from '../sortable-table-cell';
import { CONNECTION_KINDS } from '../../../../utils/Enum';
import { FormatId } from '@/components/DataFormatter';

export const NodeTableConfig = (
  switchView,
  meshSyncResources,
  k8sConfig,
  connectionMetadataState,
) => {
  const ping = useKubernetesHook();
  return {
    name: 'Node',
    colViews: [
      ['id', 'na'],
      ['metadata.name', 'xs'],
      ['apiVersion', 's'],
      ['status.attribute', 'm'],
      ['status.attribute', 'm'],
      ['cluster_id', 'xs'],
      ['status.attribute', 'm'],
      ['status.attribute', 'm'],
      ['metadata.creationTimestamp', 'l'],
    ],
    columns: [
      {
        name: 'id',
        label: 'ID',
        options: {
          display: false,
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
                data={
                  meshSyncResources[tableMeta.rowIndex]
                    ? meshSyncResources[tableMeta.rowIndex]?.component_metadata
                    : {}
                }
                value={value}
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
        name: 'status.attribute',
        label: 'CPU',
        options: {
          sort: false,
          customHeadRender: function CustomHead({ ...column }) {
            return <DefaultTableCell columnData={column} />;
          },
          customBodyRender: function CustomBody(val) {
            let attribute = JSON.parse(val);
            let capacity = attribute?.capacity;
            let cpu = getResourceStr(resourceParsers['cpu'](capacity?.cpu), 'cpu');
            return <>{cpu}</>;
          },
        },
      },
      {
        name: 'status.attribute',
        label: 'Memory',
        options: {
          sort: false,
          customHeadRender: function CustomHead({ ...column }) {
            return <DefaultTableCell columnData={column} />;
          },
          customBodyRender: function CustomBody(val) {
            let attribute = JSON.parse(val);
            let capacity = attribute?.capacity;
            let memory = getResourceStr(resourceParsers['memory'](capacity?.memory), 'memory');
            return <>{memory}</>;
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
              <>
                <TootltipWrappedConnectionChip
                  title={context.name}
                  iconSrc={
                    connectionMetadataState
                      ? connectionMetadataState[CONNECTION_KINDS.KUBERNETES]?.icon
                      : ''
                  }
                  handlePing={() => ping(context.name, context.server, context.connection_id)}
                />
              </>
            );
          },
        },
      },
      {
        name: 'status.attribute',
        label: 'Internal IP',
        options: {
          sort: false,
          customHeadRender: function CustomHead({ ...column }) {
            return <DefaultTableCell columnData={column} />;
          },
          customBodyRender: function CustomBody(val) {
            let attribute = JSON.parse(val);
            let addresses = attribute?.addresses || [];
            let internalIP =
              addresses?.find((address) => address.type === 'InternalIP')?.address || '';
            return <>{internalIP}</>;
          },
        },
      },
      {
        name: 'status.attribute',
        label: 'External IP',
        options: {
          sort: false,
          customHeadRender: function CustomHead({ ...column }) {
            return <DefaultTableCell columnData={column} />;
          },
          customBodyRender: function CustomBody(val) {
            let attribute = JSON.parse(val);
            let addresses = attribute?.addresses || [];
            let externalIP = addresses?.find((address) => address.type === 'ExternalIP')
              ?.address || <span style={{ display: 'flex', justifyContent: 'center' }}>-</span>;
            return <>{externalIP}</>;
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
  };
};
