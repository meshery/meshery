import React from 'react';
import { timeAgo } from '../../../../utils/k8s-utils';
import {
  getClusterNameFromClusterId,
  getConnectionIdFromClusterId,
} from '../../../../utils/multi-ctx';
import { SINGLE_VIEW } from '../config';
import { Title } from '../../view';
import { TootltipWrappedConnectionChip } from '../../../connections/ConnectionChip';
import useKubernetesHook from '../../../hooks/useKubernetesHook';
import { DefaultTableCell, SortableTableCell } from '../sortable-table-cell';
import { CONNECTION_KINDS } from '../../../../utils/Enum';

export const SecurityTypesConfig = (
  switchView,
  meshSyncResources,
  k8sConfig,
  connectionMetadataState,
) => {
  const ping = useKubernetesHook();
  return {
    ServiceAccount: {
      name: 'ServiceAccount',
      colViews: [
        ['id', 'na'],
        ['metadata.name', 'xs'],
        ['apiVersion', 's'],
        ['cluster_id', 'xs'],
        ['metadata.creationTimestamp', 'l'],
      ],
      columns: [
        {
          name: 'id',
          label: 'ID',
          options: {
            display: false,
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
                      ? meshSyncResources[tableMeta.rowIndex]?.component_metadata?.metadata
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
              let clusterName = getClusterNameFromClusterId(val, k8sConfig);
              return (
                <>
                  <a
                    href={'#'}
                    target="_blank"
                    rel="noopener noreferrer"
                    style={{
                      color: 'inherit',
                      textDecorationLine: 'underline',
                      cursor: 'pointer',
                      marginBottom: '0.5rem',
                    }}
                  >
                    {clusterName}
                  </a>
                </>
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
    },
    ClusterRole: {
      name: 'ClusterRole',
      colViews: [
        ['id', 'na'],
        ['metadata.name', 'xs'],
        ['apiVersion', 's'],
        ['cluster_id', 'xs'],
        ['metadata.creationTimestamp', 'l'],
      ],
      columns: [
        {
          name: 'id',
          label: 'ID',
          options: {
            display: false,
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
                      ? meshSyncResources[tableMeta.rowIndex]?.component_metadata?.metadata
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
          name: 'metadata.namespace',
          label: 'Namespace',
          options: {
            sort: false,
            customHeadRender: function CustomHead({ ...column }) {
              return <DefaultTableCell columnData={column} />;
            },
            customBodyRender: function CustomBody(value, tableMeta) {
              return (
                <>
                  <div
                    style={{
                      color: 'inherit',
                      textDecorationLine: 'underline',
                      cursor: 'pointer',
                      marginBottom: '0.5rem',
                    }}
                    onClick={() => switchView(SINGLE_VIEW, meshSyncResources[tableMeta.rowIndex])}
                  >
                    {value}
                  </div>
                </>
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
              let clusterName = getClusterNameFromClusterId(val, k8sConfig);
              let connectionId = getConnectionIdFromClusterId(val, k8sConfig);
              return (
                <>
                  <TootltipWrappedConnectionChip
                    title={clusterName}
                    iconSrc={
                      connectionMetadataState
                        ? connectionMetadataState[CONNECTION_KINDS.KUBERNETES]?.icon
                        : ''
                    }
                    handlePing={() => ping(clusterName, val, connectionId)}
                  />
                </>
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
    },
    Role: {
      name: 'Role',
      colViews: [
        ['id', 'na'],
        ['metadata.name', 'xs'],
        ['apiVersion', 's'],
        ['cluster_id', 'xs'],
        ['metadata.creationTimestamp', 'l'],
      ],
      columns: [
        {
          name: 'id',
          label: 'ID',
          options: {
            display: false,
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
                      ? meshSyncResources[tableMeta.rowIndex]?.component_metadata?.metadata
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
              let clusterName = getClusterNameFromClusterId(val, k8sConfig);
              return (
                <>
                  <a
                    href={'#'}
                    target="_blank"
                    rel="noopener noreferrer"
                    style={{
                      color: 'inherit',
                      textDecorationLine: 'underline',
                      cursor: 'pointer',
                      marginBottom: '0.5rem',
                    }}
                  >
                    {clusterName}
                  </a>
                </>
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
    },
    ClusterRoleBinding: {
      name: 'ClusterRoleBinding',
      colViews: [
        ['id', 'na'],
        ['metadata.name', 'xs'],
        ['apiVersion', 's'],
        ['cluster_id', 'xs'],
        ['metadata.creationTimestamp', 'l'],
      ],
      columns: [
        {
          name: 'id',
          label: 'ID',
          options: {
            display: false,
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
                      ? meshSyncResources[tableMeta.rowIndex]?.component_metadata?.metadata
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
              let clusterName = getClusterNameFromClusterId(val, k8sConfig);
              return (
                <>
                  <a
                    href={'#'}
                    target="_blank"
                    rel="noopener noreferrer"
                    style={{
                      color: 'inherit',
                      textDecorationLine: 'underline',
                      cursor: 'pointer',
                      marginBottom: '0.5rem',
                    }}
                  >
                    {clusterName}
                  </a>
                </>
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
    },
    RoleBinding: {
      name: 'RoleBinding',
      colViews: [
        ['id', 'na'],
        ['metadata.name', 'xs'],
        ['apiVersion', 's'],
        ['cluster_id', 'xs'],
        ['metadata.creationTimestamp', 'l'],
      ],
      columns: [
        {
          name: 'id',
          label: 'ID',
          options: {
            display: false,
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
                      ? meshSyncResources[tableMeta.rowIndex]?.component_metadata?.metadata
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
              let clusterName = getClusterNameFromClusterId(val, k8sConfig);
              return (
                <>
                  <a
                    href={'#'}
                    target="_blank"
                    rel="noopener noreferrer"
                    style={{
                      color: 'inherit',
                      textDecorationLine: 'underline',
                      cursor: 'pointer',
                      marginBottom: '0.5rem',
                    }}
                  >
                    {clusterName}
                  </a>
                </>
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
    },
  };
};
