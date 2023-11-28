import React from 'react';
import { timeAgo } from '../../../../utils/k8s-utils';
import { getClusterNameFromClusterId } from '../../../../utils/multi-ctx';
import { SINGLE_VIEW } from '../config';
import { Title } from '../../view';

import { ConnectionChip } from '../../../connections/ConnectionChip';

export const StorageTableConfig = (switchView, meshSyncResources, k8sConfig) => {
  return {
    PersistentVolume: {
      name: 'PersistentVolume',
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
            sort: false,
          },
        },
        {
          name: 'spec.attribute',
          label: 'Storage Class',
          options: {
            sort: false,
            customBodyRender: function CustomBody(val) {
              let attribute = JSON.parse(val);
              let storageClassName = attribute?.StorageClassName;
              return <>{storageClassName}</>;
            },
          },
        },
        {
          name: 'spec.attribute',
          label: 'Capacity',
          options: {
            sort: false,
            customBodyRender: function CustomBody(val) {
              let attribute = JSON.parse(val);
              let capacity = attribute?.capacity;
              let storage = capacity?.storage;
              return <>{storage}</>;
            },
          },
        },
        {
          name: 'status.attribute',
          label: 'Status',
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
    PersistentVolumeClaim: {
      name: 'PersistentVolumeClaim',
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
            sort: false,
          },
        },
        {
          name: 'spec.attribute',
          label: 'Storage Class',
          options: {
            sort: false,
            customBodyRender: function CustomBody(val) {
              let attribute = JSON.parse(val);
              let storageClassName = attribute?.StorageClassName;
              return <>{storageClassName}</>;
            },
          },
        },
        {
          name: 'spec.attribute',
          label: 'Size',
          options: {
            sort: false,
            customBodyRender: function CustomBody(val) {
              let attribute = JSON.parse(val);
              let resources = attribute?.resources;
              let requests = resources?.requests;
              let storage = requests?.storage;
              return <>{storage}</>;
            },
          },
        },
        {
          name: 'status.attribute',
          label: 'Status',
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
    StorageClass: {
      name: 'StorageClass',
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
            sort: false,
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
