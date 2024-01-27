import React from 'react';
import { timeAgo } from '../../../../utils/k8s-utils';
import { getK8sContextFromClusterId } from '../../../../utils/multi-ctx';
import { SINGLE_VIEW } from '../config';
import { Title } from '../../view';

import { TootltipWrappedConnectionChip } from '../../../connections/ConnectionChip';
import { ResizableCell } from '../../../../utils/utils';
import useKubernetesHook from '../../../hooks/useKubernetesHook';
import { DefaultTableCell, SortableTableCell } from '../sortable-table-cell';
import { CONNECTION_KINDS } from '../../../../utils/Enum';
import { FormatId } from '@/components/DataFormatter';

export const WorkloadTableConfig = (
  switchView,
  meshSyncResources,
  k8sConfig,
  connectionMetadataState,
) => {
  const ping = useKubernetesHook();
  return {
    PODS: {
      name: 'Pod',
      colViews: [
        ['id', 'na'],
        ['metadata.name', 'xs'],
        ['apiVersion', 's'],
        ['status.attribute', 's'],
        ['status.attribute', 's'],
        ['status.attribute', 'm'],
        ['metadata.namespace', 'm'],
        ['spec.attribute', 'm'],
        ['cluster_id', 'xs'],
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
            setCellProps: () => ({ style: { paddingRight: '0px', width: '0%' } }),
            setCellHeaderProps: () => ({ style: { paddingRight: '0px' } }),
            customHeadRender: function CustomHead({ ...column }) {
              return <DefaultTableCell columnData={column} />;
            },
          },
        },
        {
          name: 'status.attribute',
          label: 'Phase',
          options: {
            sort: false,
            setCellProps: () => ({ style: { paddingLeft: '0px' } }),
            setCellHeaderProps: () => ({ style: { paddingLeft: '0px' } }),
            customHeadRender: function CustomHead({ ...column }) {
              return <DefaultTableCell columnData={column} />;
            },
            customBodyRender: function CustomBody(val) {
              let attribute = JSON.parse(val);
              let phase = attribute?.phase;
              return <>{phase}</>;
            },
          },
        },
        {
          name: 'status.attribute',
          label: 'Host IP',
          options: {
            sort: false,
            customHeadRender: function CustomHead({ ...column }) {
              return <DefaultTableCell columnData={column} />;
            },
            customBodyRender: function CustomBody(val) {
              let attribute = JSON.parse(val);
              let hostIP = attribute?.hostIP;
              return <>{hostIP}</>;
            },
          },
        },
        {
          name: 'status.attribute',
          label: 'Pod IP',
          options: {
            sort: false,
            customHeadRender: function CustomHead({ ...column }) {
              return <DefaultTableCell columnData={column} />;
            },
            customBodyRender: function CustomBody(val) {
              let attribute = JSON.parse(val);
              let podIP = attribute?.podIP;
              return <>{podIP}</>;
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
          name: 'spec.attribute',
          label: 'Node',
          options: {
            sort: false,
            customHeadRender: function CustomHead({ ...column }) {
              return <DefaultTableCell columnData={column} />;
            },
            customBodyRender: function CustomBody(val) {
              let attribute = JSON.parse(val);
              let nodeName = attribute?.nodeName;
              return (
                <>
                  <ResizableCell value={nodeName} />
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
                    width="10.5rem"
                    handlePing={() => ping(context.name, context.server, context.connection_id)}
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
    DEPLOYMENT: {
      name: 'Deployment',
      colViews: [
        ['id', 'na'],
        ['metadata.name', 'xs'],
        ['apiVersion', 's'],
        ['status.attribute', 's'],
        ['spec.attribute', 's'],
        ['metadata.namespace', 'm'],
        ['cluster_id', 'xs'],
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
          name: 'status.attribute',
          label: 'Replicas',
          options: {
            sort: false,
            customHeadRender: function CustomHead({ ...column }) {
              return <DefaultTableCell columnData={column} />;
            },
            customBodyRender: function CustomBody(val) {
              let attribute = JSON.parse(val);
              let replicas = attribute?.replicas;
              return <>{replicas}</>;
            },
          },
        },
        {
          name: 'spec.attribute',
          label: 'Restart Policy',
          options: {
            sort: false,
            customHeadRender: function CustomHead({ ...column }) {
              return <DefaultTableCell columnData={column} />;
            },
            customBodyRender: function CustomBody(val) {
              let attribute = JSON.parse(val);
              let template = attribute?.template;
              let spec = template.spec;
              let restartPolicy = spec.restartPolicy;
              return <>{restartPolicy}</>;
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
                    handlePing={(event) => {
                      event.preventDefault();
                      ping(context.name, context.server, context.connection_id);
                    }}
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
    DAEMONSETS: {
      name: 'DaemonSet',
      colViews: [
        ['id', 'na'],
        ['metadata.name', 'xs'],
        ['apiVersion', 's'],
        ['spec.attribute', 's'],
        ['metadata.namespace', 'm'],
        ['cluster_id', 'xs'],
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
          name: 'spec.attribute',
          label: 'Node Selector',
          options: {
            sort: false,
            customHeadRender: function CustomHead({ ...column }) {
              return <DefaultTableCell columnData={column} />;
            },
            customBodyRender: function CustomBody(val) {
              let attribute = JSON.parse(val);
              let template = attribute?.template;
              let spec = template.spec;
              let nodeSelector = spec.nodeSelector;
              return <>{JSON.stringify(nodeSelector)}</>;
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
              let context = getK8sContextFromClusterId(val, k8sConfig);
              return (
                <TootltipWrappedConnectionChip
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
      ],
    },
    STATEFULSETS: {
      name: 'StatefulSet',
      colViews: [
        ['id', 'na'],
        ['metadata.name', 'xs'],
        ['apiVersion', 's'],
        ['status.attribute', 's'],
        ['metadata.namespace', 'm'],
        ['cluster_id', 'xs'],
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
          name: 'status.attribute',
          label: 'Replicas',
          options: {
            sort: false,
            customHeadRender: function CustomHead({ ...column }) {
              return <DefaultTableCell columnData={column} />;
            },
            customBodyRender: function CustomBody(val) {
              let attribute = JSON.parse(val);
              let replicas = attribute?.replicas;
              return <>{replicas}</>;
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
              let context = getK8sContextFromClusterId(val, k8sConfig);
              return (
                <TootltipWrappedConnectionChip
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
      ],
    },
    REPLICASETS: {
      name: 'ReplicaSet',
      colViews: [
        ['id', 'na'],
        ['metadata.name', 'xs'],
        ['apiVersion', 's'],
        ['spec.attribute', 's'],
        ['status.attribute', 's'],
        ['status.attribute', 'm'],
        ['metadata.namespace', 'm'],
        ['cluster_id', 'xs'],
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
          name: 'spec.attribute',
          label: 'Desired Replicas',
          options: {
            sort: false,
            customHeadRender: function CustomHead({ ...column }) {
              return <DefaultTableCell columnData={column} />;
            },
            customBodyRender: function CustomBody(val) {
              let attribute = JSON.parse(val);
              let replicas = attribute?.replicas;
              return <>{replicas}</>;
            },
          },
        },
        {
          name: 'status.attribute',
          label: 'Current Replicas',
          options: {
            sort: false,
            customHeadRender: function CustomHead({ ...column }) {
              return <DefaultTableCell columnData={column} />;
            },
            customBodyRender: function CustomBody(val) {
              let attribute = JSON.parse(val);
              let replicas = attribute?.replicas;
              return <>{replicas}</>;
            },
          },
        },
        {
          name: 'status.attribute',
          label: 'Ready Replicas',
          options: {
            sort: false,
            customHeadRender: function CustomHead({ ...column }) {
              return <DefaultTableCell columnData={column} />;
            },
            customBodyRender: function CustomBody(val) {
              let attribute = JSON.parse(val);
              let readyReplicas = attribute?.readyReplicas;
              return <>{readyReplicas}</>;
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
            sortThirdClickReset: true,
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
              let context = getK8sContextFromClusterId(val, k8sConfig);
              return (
                <TootltipWrappedConnectionChip
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
      ],
    },
    REPLICATIONCONTROLLERS: {
      name: 'ReplicationController',
      colViews: [
        ['id', 'na'],
        ['metadata.name', 'xs'],
        ['apiVersion', 's'],
        ['spec.attribute', 's'],
        ['status.attribute', 's'],
        ['metadata.namespace', 'm'],
        ['cluster_id', 'xs'],
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
          name: 'spec.attribute',
          label: 'Desired Replicas',
          options: {
            sort: false,
            customHeadRender: function CustomHead({ ...column }) {
              return <DefaultTableCell columnData={column} />;
            },
            customBodyRender: function CustomBody(val) {
              let attribute = JSON.parse(val);
              let replicas = attribute?.replicas;
              return <>{replicas}</>;
            },
          },
        },
        {
          name: 'status.attribute',
          label: 'Current Replicas',
          options: {
            sort: false,
            customHeadRender: function CustomHead({ ...column }) {
              return <DefaultTableCell columnData={column} />;
            },
            customBodyRender: function CustomBody(val) {
              let attribute = JSON.parse(val);
              let replicas = attribute?.replicas;
              return <>{replicas}</>;
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
              let context = getK8sContextFromClusterId(val, k8sConfig);
              return (
                <TootltipWrappedConnectionChip
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
      ],
    },
    JOBS: {
      name: 'Job',
      colViews: [
        ['id', 'na'],
        ['metadata.name', 'xs'],
        ['apiVersion', 's'],
        ['metadata.namespace', 'm'],
        ['cluster_id', 'xs'],
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
              let context = getK8sContextFromClusterId(val, k8sConfig);
              return (
                <TootltipWrappedConnectionChip
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
      ],
    },
    CRONJOBS: {
      name: 'CronJob',
      colViews: [
        ['id', 'na'],
        ['metadata.name', 'xs'],
        ['apiVersion', 's'],
        ['spec.attribute', 's'],
        ['spec.attribute', 's'],
        ['metadata.namespace', 'm'],
        ['cluster_id', 'xs'],
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
          name: 'spec.attribute',
          label: 'Schedule',
          options: {
            sort: false,
            customHeadRender: function CustomHead({ ...column }) {
              return <DefaultTableCell columnData={column} />;
            },
            customBodyRender: function CustomBody(val) {
              let attribute = JSON.parse(val);
              let schedule = attribute?.schedule;
              return <>{schedule}</>;
            },
          },
        },
        {
          name: 'spec.attribute',
          label: 'Suspend',
          options: {
            sort: false,
            customHeadRender: function CustomHead({ ...column }) {
              return <DefaultTableCell columnData={column} />;
            },
            customBodyRender: function CustomBody(val) {
              let attribute = JSON.parse(val);
              let suspend = attribute?.suspend;
              return <>{suspend}</>;
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
              let context = getK8sContextFromClusterId(val, k8sConfig);
              return (
                <TootltipWrappedConnectionChip
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
      ],
    },
  };
};
