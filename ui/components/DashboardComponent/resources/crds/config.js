import React, { useMemo } from 'react';
import { timeAgo } from '../../../../utils/k8s-utils';
import { getK8sClusterIdsFromCtxId, getK8sContextFromClusterId } from '@/utils/multi-ctx';
import { getAllCustomResourceDefinitionsKinds, SINGLE_VIEW } from '../config';
import { Title } from '../../view';
import { TootltipWrappedConnectionChip } from '../../../connections/ConnectionChip';
import useKubernetesHook from '../../../hooks/useKubernetesHook';
import { DefaultTableCell, SortableTableCell } from '../sortable-table-cell';
import { CONNECTION_KINDS } from '../../../../utils/Enum';
import { FormatId } from '@/components/DataFormatter';
import { useGetMeshSyncResourceKindsQuery } from '@/rtk-query/meshsync';

export const CustomResourceConfig = (
  switchView,
  meshSyncResources,
  k8sConfig,
  connectionMetadataState,
  _workloadType,
  selectedK8sContexts,
) => {
  const ping = useKubernetesHook();
  const clusterIds = getK8sClusterIdsFromCtxId(selectedK8sContexts, k8sConfig);
  const isClusterIdsEmpty = clusterIds.size === 0 || clusterIds.length === 0;

  const { data: clusterSummary } = useGetMeshSyncResourceKindsQuery(
    {
      page: 0,
      pagesize: 'all',
      clusterIds: clusterIds,
    },
    {
      skip: isClusterIdsEmpty,
    },
  );
  const customResources = useMemo(() => {
    return getAllCustomResourceDefinitionsKinds(clusterSummary?.kinds);
  }, [clusterSummary?.kinds]);
  const customResourceConfigs = {};

  customResources?.forEach((resource) => {
    customResourceConfigs[resource?.Kind] = {
      name: resource?.Kind,
      model: resource?.Model, // model is used to identify the resource image
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
                  kind={resource?.Kind}
                  model={resource?.Model}
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
            customBodyRender: function CustomBody(value) {
              return <>{value}</>;
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
    };
  });

  return customResourceConfigs;
};
