import React from 'react';
import { timeAgo } from '../../../../utils/k8s-utils';
import { getK8sContextFromClusterId } from '../../../../utils/multi-ctx';
import { SINGLE_VIEW } from '../config';
import { Title } from '../../view';
import { TooltipWrappedConnectionChip } from '../../../connections/ConnectionChip';
import { ResizableCell } from '../../../../utils/utils';
import { DefaultTableCell, SortableTableCell } from '../sortable-table-cell';
import { CONNECTION_KINDS } from '../../../../utils/Enum';
import { FormatId } from '@/components/data-formatter';

export const buildPodColumns = ({
  switchView,
  meshSyncResources,
  k8sConfig,
  connectionMetadataState,
  workloadType,
  ping,
}) => ({
  name: 'Pod',
  colViews: [
    ['id', 'na'],
    ['metadata.name', 'xs'],
    ['apiVersion', 'na'],
    ['status.attribute', 's'],
    ['status.attribute', 's'],
    ['status.attribute', 'm'],
    ['metadata.namespace', 'm'],
    ['spec.attribute', 'm'],
    ['cluster_id', 'xs'],
    ['metadata.creationTimestamp', 'l'],
    ['status.attribute', 'm'],
    ['spec.attribute', 'm'],
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
        setCellProps: () => ({ style: { paddingRight: '0px', width: '0%' } }),
        setCellHeaderProps: () => ({ style: { paddingRight: '0px' } }),
        customHeadRender: function CustomHead({ ...column }) {
          return <DefaultTableCell columnData={column} />;
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
      },
    },
    {
      name: 'status.attribute',
      label: 'Restarts',
      options: {
        sort: false,
        customHeadRender: function CustomHead({ ...column }) {
          return <DefaultTableCell columnData={column} />;
        },
        customBodyRender: function CustomBody(value) {
          const parsedStatus = JSON.parse(value);
          const totalRestarts = parsedStatus?.containerStatuses?.reduce(
            (sum, container) => sum + (container.restartCount || 0),
            0,
          );
          return <>{totalRestarts || 0}</>;
        },
      },
    },
    {
      name: 'spec.attribute',
      label: 'Containers',
      options: {
        sort: false,
        customHeadRender: function CustomHead({ ...column }) {
          return <DefaultTableCell columnData={column} />;
        },
        customBodyRender: function CustomBody(value) {
          const parsedSpec = JSON.parse(value);
          return <>{parsedSpec.containers.length}</>;
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
              <TooltipWrappedConnectionChip
                title={context.name}
                iconSrc={
                  connectionMetadataState
                    ? connectionMetadataState[CONNECTION_KINDS.KUBERNETES]?.icon
                    : ''
                }
                width="10.5rem"
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
