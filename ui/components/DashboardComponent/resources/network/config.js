import React from 'react';
import { getStatus, timeAgo } from '../../../../utils/k8s-utils';
import { getK8sContextFromClusterId } from '../../../../utils/multi-ctx';
import { SINGLE_VIEW } from '../config';

import { Title } from '../../view';

import { TootltipWrappedConnectionChip } from '../../../connections/ConnectionChip';
import useKubernetesHook from '../../../hooks/useKubernetesHook';
import { DefaultTableCell, SortableTableCell } from '../sortable-table-cell';
import { CONNECTION_KINDS } from '../../../../utils/Enum';
import { FormatId } from '@/components/DataFormatter';

export const NetWorkTableConfig = (
  switchView,
  meshSyncResources,
  k8sConfig,
  connectionMetadataState,
  workloadType,
) => {
  const ping = useKubernetesHook();
  return {
    Service: {
      name: 'Service',
      colViews: [
        ['id', 'na'],
        ['metadata.name', 'xs'],
        ['apiVersion', 'na'],
        ['spec.attribute', 's'],
        ['spec.attribute', 's'],
        ['status.attribute', 'na'],
        ['spec.attribute', 'm'],
        ['metadata.namespace', 'm'],
        ['cluster_id', 'xs'],
        ['metadata.creationTimestamp', 'l'],
        ['status.attribute', 's'],
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
          name: 'spec.attribute',
          label: 'Type',
          options: {
            sort: false,
            customHeadRender: function CustomHead({ ...column }) {
              return <DefaultTableCell columnData={column} />;
            },
            customBodyRender: function CustomBody(val) {
              let attribute = JSON.parse(val);
              let type = attribute?.type;
              return <>{type}</>;
            },
          },
        },
        {
          name: 'spec.attribute',
          label: 'Cluster IP',
          options: {
            sort: false,
            customHeadRender: function CustomHead({ ...column }) {
              return <DefaultTableCell columnData={column} />;
            },
            customBodyRender: function CustomBody(val) {
              let attribute = JSON.parse(val);
              let clusterIP = attribute?.clusterIP;
              return <>{clusterIP}</>;
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
              let loadbalancer = attribute?.loadbalancer;
              let ingresses = loadbalancer?.ingress;
              return (
                <>
                  {ingresses?.map((ingress, i) => {
                    return (
                      <>
                        {ingress.hostname}
                        {i < ingresses.length - 1 && ','}
                      </>
                    );
                  })}
                </>
              );
            },
          },
        },
        {
          name: 'spec.attribute',
          label: 'Ports',
          options: {
            sort: false,
            customHeadRender: function CustomHead({ ...column }) {
              return <DefaultTableCell columnData={column} />;
            },
            customBodyRender: function CustomBody(val, tableMeta) {
              let attribute = JSON.parse(val);
              let ports = attribute?.ports;

              const showViewAll = ports?.length > 1;
              return (
                <>
                  <div style={{ display: 'flex' }}>
                    {ports?.slice(0, 1).map((p, i) => (
                      <div key={i}>
                        {`${p.port}/${p.targetPort}:${p.protocol}`}
                        {ports.length > 1 && i === 0 && ','}
                      </div>
                    ))}
                    {showViewAll && (
                      <span
                        style={{
                          color: 'inherit',
                          textDecorationLine: 'underline',
                          cursor: 'pointer',
                          marginLeft: '0.5rem',
                          width: 'max-content',
                        }}
                        onClick={() =>
                          switchView(SINGLE_VIEW, meshSyncResources[tableMeta.rowIndex])
                        }
                      >
                        View all
                      </span>
                    )}
                  </div>
                </>
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
            customBodyRender: function CustomBody(value) {
              const status = getStatus(value);
              return <>{status}</>;
            },
          },
        },
      ],
    },
    Endpoints: {
      name: 'Endpoints',
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
    EndpointSlice: {
      name: 'EndpointSlice',
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
            display: false,
            customBodyRender: (value) => <FormatId id={value} />,
          },
        },
        {
          name: 'metadata.name',
          label: 'Name',
          options: {
            sort: false,
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
          },
        },
        {
          name: 'cluster_id',
          label: 'Cluster',
          options: {
            sort: true,
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
            customBodyRender: function CustomBody(value) {
              let time = timeAgo(value);
              return <>{time}</>;
            },
          },
        },
      ],
    },
    Ingress: {
      name: 'Ingress',
      colViews: [
        ['id', 'na'],
        ['metadata.name', 'xs'],
        ['apiVersion', 'na'],
        ['status.attribute', 'm'],
        ['spec.attribute', 'm'],
        ['spec.attribute', 'm'],
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
          name: 'status.attribute',
          label: 'LoadBalancer IP',
          options: {
            sort: false,
            customHeadRender: function CustomHead({ ...column }) {
              return <DefaultTableCell columnData={column} />;
            },
            customBodyRender: function CustomBody(val) {
              let attribute = JSON.parse(val);
              let loadBalancer = attribute?.loadBalancer;
              let ingress = loadBalancer?.ingress;
              const IPs = ingress?.map((ingress) => ingress.ip);
              return <>{IPs?.join(',')}</>;
            },
          },
        },
        {
          name: 'spec.attribute',
          label: 'Rules',
          options: {
            sort: false,
            customHeadRender: function CustomHead({ ...column }) {
              return <DefaultTableCell columnData={column} />;
            },
            customBodyRender: function CustomBody(val) {
              let attribute = JSON.parse(val);

              let rules = attribute?.rules;
              return (
                <>
                  {rules?.map((rule, i) => {
                    return (
                      <>
                        {`${rule.host}`}
                        {i < rules.length - 1 && ','}
                      </>
                    );
                  })}
                </>
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
    IngressClass: {
      name: 'IngressClass',
      colViews: [
        ['id', 'na'],
        ['metadata.name', 'xs'],
        ['apiVersion', 'na'],
        ['spec.attribute', 'm'],
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
          name: 'spec.attribute',
          label: 'Controller',
          options: {
            sort: false,
            customHeadRender: function CustomHead({ ...column }) {
              return <DefaultTableCell columnData={column} />;
            },
            customBodyRender: function CustomBody(val) {
              let attribute = JSON.parse(val);
              let controller = attribute?.controller;
              return <>{controller}</>;
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
    NetworkPolicy: {
      name: 'NetworkPolicy',
      colViews: [
        ['id', 'na'],
        ['metadata.name', 'xs'],
        ['apiVersion', 'na'],
        ['spec.attribute', 'm'],
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
          name: 'spec.attribute',
          label: 'Ports',
          options: {
            sort: false,
            customHeadRender: function CustomHead({ ...column }) {
              return <DefaultTableCell columnData={column} />;
            },
            customBodyRender: function CustomBody(val) {
              let attribute = JSON.parse(val);
              let policyTypes = attribute?.policyTypes;
              return (
                <>
                  {policyTypes?.map((policy, i) => {
                    return (
                      <>
                        {`${policy}`}
                        {i < policyTypes.length - 1 && ','}
                      </>
                    );
                  })}
                </>
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
