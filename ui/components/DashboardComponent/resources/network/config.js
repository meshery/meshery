import React from 'react';
import { timeAgo } from '../../../../utils/k8s-utils';
import { getClusterNameFromClusterId } from '../../../../utils/multi-ctx';
import { SINGLE_VIEW } from '../config';

import { Title } from '../../view';

import { ConnectionChip } from '../../../connections/ConnectionChip';

export const NetWorkTableConfig = (switchView, meshSyncResources, k8sConfig) => {
  return {
    Service: {
      name: 'Service',
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
            sortThirdClickReset: true,
            customBodyRender: function CustomBody(value, tableMeta) {
              return (
                <Title
                  onClick={() => switchView(SINGLE_VIEW, meshSyncResources[tableMeta.rowIndex])}
                  data={
                    meshSyncResources[tableMeta.rowIndex]
                      ? meshSyncResources[tableMeta.rowIndex].component_metadata.metadata
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
            sort: false,
          },
        },
        {
          name: 'spec.attribute',
          label: 'Type',
          options: {
            sort: false,
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
            customBodyRender: function CustomBody(val, tableMeta) {
              let attribute = JSON.parse(val);
              let ports = attribute?.ports;

              const showViewAll = ports?.length > 2;
              return (
                <>
                  <div style={{ display: 'flex' }}>
                    {ports?.slice(0, 2).map((p, i) => (
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
            sort: false,
            sortThirdClickReset: true,
            customBodyRender: function CustomBody(val) {
              let clusterName = getClusterNameFromClusterId(val, k8sConfig);
              return (
                <>
                  <ConnectionChip title={clusterName} iconSrc="/static/img/kubernetes.svg" />
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
            sortThirdClickReset: true,
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
            sortThirdClickReset: true,
            customBodyRender: function CustomBody(value, tableMeta) {
              return (
                <Title
                  onClick={() => switchView(SINGLE_VIEW, meshSyncResources[tableMeta.rowIndex])}
                  data={
                    meshSyncResources[tableMeta.rowIndex]
                      ? meshSyncResources[tableMeta.rowIndex].component_metadata.metadata
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
            sort: false,
          },
        },
        {
          name: 'metadata.namespace',
          label: 'Namespace',
          options: {
            sort: false,
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
            sort: false,
            sortThirdClickReset: true,
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
            sortThirdClickReset: true,
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
            sortThirdClickReset: true,
            customBodyRender: function CustomBody(value, tableMeta) {
              return (
                <Title
                  onClick={() => switchView(SINGLE_VIEW, meshSyncResources[tableMeta.rowIndex])}
                  data={
                    meshSyncResources[tableMeta.rowIndex]
                      ? meshSyncResources[tableMeta.rowIndex].component_metadata.metadata
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
            sort: false,
          },
        },
        {
          name: 'spec.attribute',
          label: 'Rules',
          options: {
            sort: false,
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
            sort: false,
            sortThirdClickReset: true,
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
            sortThirdClickReset: true,
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
            sortThirdClickReset: true,
            customBodyRender: function CustomBody(value, tableMeta) {
              return (
                <Title
                  onClick={() => switchView(SINGLE_VIEW, meshSyncResources[tableMeta.rowIndex])}
                  data={
                    meshSyncResources[tableMeta.rowIndex]
                      ? meshSyncResources[tableMeta.rowIndex].component_metadata.metadata
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
            sort: false,
          },
        },
        {
          name: 'spec.attribute',
          label: 'Controller',
          options: {
            sort: false,
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
            sort: false,
            sortThirdClickReset: true,
            customBodyRender: function CustomBody(val) {
              let clusterName = getClusterNameFromClusterId(val, k8sConfig);
              return (
                <>
                  <ConnectionChip title={clusterName} iconSrc="/static/img/kubernetes.svg" />
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
            sortThirdClickReset: true,
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
            sortThirdClickReset: true,
            customBodyRender: function CustomBody(value, tableMeta) {
              return (
                <Title
                  onClick={() => switchView(SINGLE_VIEW, meshSyncResources[tableMeta.rowIndex])}
                  data={
                    meshSyncResources[tableMeta.rowIndex]
                      ? meshSyncResources[tableMeta.rowIndex].component_metadata.metadata
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
            sort: false,
          },
        },
        {
          name: 'spec.attribute',
          label: 'Ports',
          options: {
            sort: false,
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
            sort: false,
            sortThirdClickReset: true,
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
            sortThirdClickReset: true,
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
