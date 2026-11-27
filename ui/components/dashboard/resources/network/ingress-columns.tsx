import React from 'react';
import { timeAgo } from '../../../../utils/k8s-utils';
import { getK8sContextFromClusterId } from '../../../../utils/multi-ctx';
import { SINGLE_VIEW } from '../config';
import { Title } from '../../view';
import { TooltipWrappedConnectionChip } from '../../../connections/ConnectionChip';
import { DefaultTableCell, SortableTableCell } from '../sortable-table-cell';
import { CoreConnectionKinds } from '../../../../utils/Enum';
import { FormatId } from '@/components/data-formatter';

export const buildIngressColumns = ({
  switchView,
  meshSyncResources,
  k8sConfig,
  connectionMetadataState,
  workloadType,
  ping,
}) => ({
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
          // A load balancer reports an IP or, on some clouds such as AWS, a hostname.
          const addresses = ingress?.map(({ ip, hostname }) => ip || hostname).filter(Boolean);
          return <>{addresses?.join(', ')}</>;
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
          // A rule without a host matches every host, which kubectl shows as "*".
          return <>{rules?.map((rule) => rule.host || '*').join(', ')}</>;
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
              <TooltipWrappedConnectionChip
                title={context.name}
                iconSrc={
                  connectionMetadataState
                    ? connectionMetadataState[CoreConnectionKinds.kubernetes]?.icon
                    : ''
                }
                handlePing={() => ping(context.name, context.server, context.connectionId)}
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
});
