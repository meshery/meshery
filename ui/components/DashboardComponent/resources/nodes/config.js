import React from 'react';
import { getResourceStr, resourceParsers, timeAgo } from '../../../../utils/k8s-utils';
import { getClusterNameFromClusterId } from '../../../../utils/multi-ctx';
import { SINGLE_VIEW } from '../config';
import { ConnectionChip } from '../../../connections/ConnectionChip';

export const NodeTableConfig = (switchView, meshSyncResources, k8sConfig) => {
  return {
    name: 'Node',
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
        name: 'apiVersion',
        label: 'API version',
        options: {
          sort: false,
        },
      },
      {
        name: 'status.attribute',
        label: 'CPU',
        options: {
          sort: false,
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
        name: 'status.attribute',
        label: 'Internal IP',
        options: {
          sort: false,
          sortThirdClickReset: true,
          customBodyRender: function CustomBody(val) {
            let attribute = JSON.parse(val);
            let addresses = attribute?.addresses;
            let internalIP =
              addresses.find((address) => address.type === 'InternalIP')?.address || '';
            return <>{internalIP}</>;
          },
        },
      },
      {
        name: 'status.attribute',
        label: 'External IP',
        options: {
          sort: false,
          sortThirdClickReset: true,
          customBodyRender: function CustomBody(val) {
            let attribute = JSON.parse(val);
            let addresses = attribute?.addresses;
            let externalIP =
              addresses.find((address) => address.type === 'ExternalIP')?.address || '';
            return <>{externalIP}</>;
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
  };
};
