import React from 'react';
import { timeAgo } from '../../../../utils/k8s-utils';
import { getClusterNameFromClusterId } from '../../../../utils/multi-ctx';
import { SINGLE_VIEW } from '../config';
import { Title } from '../../view';

import { ConnectionChip } from '../../../connections/ConnectionChip';
import { ConditionalTooltip } from '../../../../utils/utils';

export const WorkloadTableConfig = (switchView, meshSyncResources, k8sConfig) => {
  return {
    PODS: {
      name: 'Pod',
      columns: [
        {
          name: 'id',
          label: 'ID',
          options: {
            display: false,
            customBodyRender: (value) => <ConditionalTooltip value={value} maxLength={10} />,
          },
        },
        {
          name: 'metadata.name',
          label: 'Name',
          options: {
            sort: false,
            sortThirdClickReset: true,
            customBodyRender: function CustomBody(value, tableMeta) {
              console.log(meshSyncResources, value, ';;;;;pp');
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
          name: 'status.attribute',
          label: 'Phase',
          options: {
            sort: false,
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
          name: 'spec.attribute',
          label: 'Node',
          options: {
            sort: false,
            customBodyRender: function CustomBody(val) {
              let attribute = JSON.parse(val);
              let nodeName = attribute?.nodeName;
              return <>{nodeName}</>;
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
    DEPLOYMENT: {
      name: 'Deployment',
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
          name: 'status.attribute',
          label: 'Replicas',
          options: {
            sort: false,
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
    DAEMONSETS: {
      name: 'DaemonSet',
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
          label: 'Node Selector',
          options: {
            sort: false,
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
    STATEFULSETS: {
      name: 'StatefulSet',
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
          name: 'status.attribute',
          label: 'Replicas',
          options: {
            sort: false,
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
    REPLICASETS: {
      name: 'ReplicaSet',
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
          label: 'Desired Replicas',
          options: {
            sort: false,
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
    REPLICATIONCONTROLLERS: {
      name: 'ReplicationController',
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
          label: 'Desired Replicas',
          options: {
            sort: false,
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
    JOBS: {
      name: 'Job',
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
    CRONJOBS: {
      name: 'CronJob',
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
          label: 'Schedule',
          options: {
            sort: false,
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
