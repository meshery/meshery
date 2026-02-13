import React from 'react';
import { NoSsr, TableCell, TableSortLabel } from '@sistent/sistent';
import MUIDataTable from '@sistent/mui-datatables';

function NodeDetails(props: any) {
  const chartData = props.result;

  const columns = [
    {
      name: 'hostname',
      label: 'Hostname',
      options: {
        filter: false,
        sort: true,
        searchable: true,
        customHeadRender: function CustomHead({ index, ...column }: any, sortColumn: any) {
          return (
            <TableCell key={index} onClick={() => sortColumn(index)}>
              <TableSortLabel
                active={column.sortDirection != null}
                direction={column.sortDirection || 'asc'}
              >
                <b>{column.label}</b>
              </TableSortLabel>
            </TableCell>
          );
        },
      },
    },
    {
      name: 'cpu',
      label: 'CPU',
      options: {
        filter: false,
        sort: true,
        searchable: true,
        customHeadRender: function CustomHead({ index, ...column }: any, sortColumn: any) {
          return (
            <TableCell key={index} onClick={() => sortColumn(index)}>
              <TableSortLabel
                active={column.sortDirection != null}
                direction={column.sortDirection || 'asc'}
              >
                <b>{column.label}</b>
              </TableSortLabel>
            </TableCell>
          );
        },
      },
    },
    {
      name: 'memory',
      label: 'Memory',
      options: {
        filter: false,
        sort: true,
        searchable: true,
        customHeadRender: function CustomHead({ index, ...column }: any, sortColumn: any) {
          return (
            <TableCell key={index} onClick={() => sortColumn(index)}>
              <TableSortLabel
                active={column.sortDirection != null}
                direction={column.sortDirection || 'asc'}
              >
                <b>{column.label}</b>
              </TableSortLabel>
            </TableCell>
          );
        },
      },
    },
    {
      name: 'arch',
      label: 'Arch',
      options: {
        filter: false,
        sort: true,
        searchable: true,
        customHeadRender: function CustomHead({ index, ...column }: any, sortColumn: any) {
          return (
            <TableCell key={index} onClick={() => sortColumn(index)}>
              <TableSortLabel
                active={column.sortDirection != null}
                direction={column.sortDirection || 'asc'}
              >
                <b>{column.label}</b>
              </TableSortLabel>
            </TableCell>
          );
        },
      },
    },
    {
      name: 'os',
      label: 'OS',
      options: {
        filter: false,
        sort: true,
        searchable: true,
        customHeadRender: function CustomHead({ index, ...column }: any, sortColumn: any) {
          return (
            <TableCell key={index} onClick={() => sortColumn(index)}>
              <TableSortLabel
                active={column.sortDirection != null}
                direction={column.sortDirection || 'asc'}
              >
                <b>{column.label}</b>
              </TableSortLabel>
            </TableCell>
          );
        },
      },
    },
    {
      name: 'kubeletVersion',
      label: 'Kubelet Version',
      options: {
        filter: false,
        sort: true,
        searchable: true,
        customHeadRender: function CustomHead({ index, ...column }: any, sortColumn: any) {
          return (
            <TableCell key={index} onClick={() => sortColumn(index)}>
              <TableSortLabel
                active={column.sortDirection != null}
                direction={column.sortDirection || 'asc'}
              >
                <b>{column.label}</b>
              </TableSortLabel>
            </TableCell>
          );
        },
      },
    },
    {
      name: 'containerRuntime',
      label: 'Container Runtime',
      options: {
        filter: false,
        sort: true,
        searchable: true,
        customHeadRender: function CustomHead({ index, ...column }: any, sortColumn: any) {
          return (
            <TableCell key={index} onClick={() => sortColumn(index)}>
              <TableSortLabel
                active={column.sortDirection != null}
                direction={column.sortDirection || 'asc'}
              >
                <b>{column.label}</b>
              </TableSortLabel>
            </TableCell>
          );
        },
      },
    },
  ];

  let data: any[][] = [];

  const options = {
    filter: false,
    selectableRows: 'none' as const,
  };

  let server = chartData?.kubernetes?.server_version;

  chartData?.kubernetes?.nodes?.map((node: any) => {
    let arr: any[] = [];
    let m = node?.allocatable_memory;
    const memString = String(m).slice(0, String(m).length - 2);
    const mem = (parseFloat(memString) * 0.000001024).toPrecision(5);
    arr.push(node?.hostname);
    arr.push(node?.allocatable_cpu);
    arr.push(mem + 'Gi');
    arr.push(node?.architecture);
    arr.push(node?.operating_system);
    arr.push(node?.kubelet_version);
    arr.push(node?.container_runtime_version);

    data.push(arr);
  });

  return (
    <NoSsr>
      <MUIDataTable
        title={<div style={{ fontSize: 18 }}>Kubernetes Server Version: {server}</div>}
        data={data}
        options={options}
        columns={columns}
      />
    </NoSsr>
  );
}

export default NodeDetails;
