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

export const NetWorkTableConfig = (
  switchView,
  meshSyncResources,
  k8sConfig,
  connectionMetadataState,
) => {
  const ping = useKubernetesHook();
  return {
    Service: {
      name: 'Service',
      colViews: [
        ['id', 'na'],
        ['metadata.name', 'xs'],
        ['apiVersion', 's'],
        ['spec.attribute', 's'],
        ['spec.attribute', 's'],
        ['status.attribute', 'na'],
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
                        {i < 1 && ','}
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
    Endpoints: {
      name: 'Endpoints',
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
    Ingress: {
      name: 'Ingress',
      colViews: [
        ['id', 'na'],
        ['metadata.name', 'xs'],
        ['apiVersion', 's'],
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
          label: 'Rules',
          options: {
            sort: false,
            customHeadRender: function CustomHead({ ...column }) {
              return <DefaultTableCell columnData={column} />;
            },
            customBodyRender: function CustomBody(val) {
              let attribute = JSON.parse(val);
              let ingressRules = attribute?.ingressRule;
              return (
                <>
                  {ingressRules?.map((rule, i) => {
                    return (
                      <>
                        {`${rule.host}`}
                        {i < ingressRules.length - 1 && ','}
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
    IngressClass: {
      name: 'IngressClass',
      colViews: [
        ['id', 'na'],
        ['metadata.name', 'xs'],
        ['apiVersion', 's'],
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
    NetworkPolicy: {
      name: 'NetworkPolicy',
      colViews: [
        ['id', 'na'],
        ['metadata.name', 'xs'],
        ['apiVersion', 's'],
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
