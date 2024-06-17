import { useState, useEffect } from 'react';
import MUIDataTable from 'mui-datatables';
import { TableCell, Tooltip, TableSortLabel } from '@material-ui/core';

import debounce from '../utils/debounce';
import { getDuplicateModels, getDuplicateComponents } from '../api/meshmodel';
import { MODELS, COMPONENTS } from '../constants/navigator';

const DuplicatesDataTable = ({ view, rowData, classes }) => {
  const [resourcesDetail, setResourcesDetail] = useState();
  const [count, setCount] = useState();
  const [page, setPage] = useState(0);
  const [, setSearchText] = useState(null);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const { kind, model, version } = rowData;

  const getDuplicatedModels = async (model, version) => {
    const { total_count, models } = await getDuplicateModels(model, version);
    setCount(total_count);
    setResourcesDetail(models);
  };

  const getDuplicatedComponents = async (componentKind, apiVersion, modelName) => {
    const { total_count, components } = await getDuplicateComponents(
      componentKind,
      modelName,
      apiVersion,
    );
    setCount(total_count);
    setResourcesDetail(components);
  };

  const meshmodel_columns = [
    {
      name: view === COMPONENTS ? 'kind' : 'displayName',
      label: `Name (${count})`,
      options: {
        sort: true,
        sortDescFirst: true,
        sortThirdClickReset: true,
        searchable: true,
        customHeadRender: function CustomHead({ index, ...column }, sortColumn) {
          return (
            <TableCell align={'start'} key={index} onClick={() => sortColumn(index)}>
              <TableSortLabel
                active={column.sortDirection != null}
                direction={column.sortDirection || 'asc'}
              >
                <b>{column.label}</b>
              </TableSortLabel>
            </TableCell>
          );
        },
        customBodyRender: (value) => (
          <Tooltip title={value} placement="top">
            <div>{value}</div>
          </Tooltip>
        ),
      },
    },
    {
      name: view === COMPONENTS ? 'apiVersion' : 'version',
      label: view === COMPONENTS ? 'Api Version' : 'Version',
      options: {
        sort: true,
        sortDescFirst: true,
        sortThirdClickReset: true,
        searchable: true,
        customHeadRender: function CustomHead({ index, ...column }, sortColumn) {
          return (
            <TableCell align={'start'} key={index} onClick={() => sortColumn(index)}>
              <TableSortLabel
                active={column.sortDirection != null}
                direction={column.sortDirection || 'asc'}
              >
                <b>{column.label}</b>
              </TableSortLabel>
            </TableCell>
          );
        },
        customBodyRender: (value) => (
          <Tooltip title={value} placement="top">
            <div>{value}</div>
          </Tooltip>
        ),
      },
    },
    {
      name: 'category',
      label: 'Category Name',
      options: {
        sort: false,
        display: view === MODELS ? 'true' : 'false',
        searchable: true,
        customHeadRender: function CustomHead({ index, ...column }) {
          return (
            <TableCell align={'start'} key={index}>
              <TableSortLabel>
                <b>{column.label}</b>
              </TableSortLabel>
            </TableCell>
          );
        },
        customBodyRender: (value) => {
          if (!(view === COMPONENTS)) {
            const { modelDisplayName, name } = value;
            return (
              <Tooltip title={view === MODELS ? name : modelDisplayName} placement="top">
                <div>{view === MODELS ? name : modelDisplayName}</div>
              </Tooltip>
            );
          }
        },
      },
    },
    {
      name: 'metadata',
      label: 'Model',
      options: {
        sort: false,
        display: view === COMPONENTS ? 'true' : 'false',
        searchable: true,
        customHeadRender: function CustomHead({ index, ...column }) {
          return (
            <TableCell align={'start'} key={index}>
              <TableSortLabel>
                <b>{column.label}</b>
              </TableSortLabel>
            </TableCell>
          );
        },
        customBodyRender: (value) => {
          if (!(view === MODELS)) {
            const { modelDisplayName } = value;
            return (
              <Tooltip title={modelDisplayName} placement="top">
                <div>{modelDisplayName}</div>
              </Tooltip>
            );
          }
        },
      },
    },
    {
      name: 'metadata',
      label: 'Sub Category',
      options: {
        sort: false,
        display: view === COMPONENTS ? 'true' : 'false',
        searchable: true,
        customHeadRender: function CustomHead({ index, ...column }) {
          return (
            <TableCell align={'start'} key={index}>
              <TableSortLabel>
                <b>{column.label}</b>
              </TableSortLabel>
            </TableCell>
          );
        },
        customBodyRender: (value) => {
          const { subCategory } = value;
          return (
            <Tooltip title={subCategory} placement="top">
              <div>{subCategory}</div>
            </Tooltip>
          );
        },
      },
    },
  ];

  const meshmodel_options = {
    rowsPerPage: rowsPerPage,
    rowsPerPageOptions: [10, 25, 50, 100],
    page: page,
    count: count,
    sort: true,
    download: false,
    print: false,
    filter: false,
    selectableRows: 'none',
    serverSide: true,
    onChangePage: debounce((p) => setPage(p), 200),
    onSearchChange: debounce((searchText) => setSearchText(searchText)),
    onChangeRowsPerPage: debounce((rowsPerPage) => {
      setRowsPerPage(rowsPerPage);
      setPage(0);
    }),
  };

  useEffect(() => {
    if (view === MODELS) getDuplicatedModels(model, version);
    else if (view === COMPONENTS) getDuplicatedComponents(kind, model, version);
  }, []);

  return (
    <div style={{ width: '100%' }}>
      <MUIDataTable
        title={<div className={classes.tableHeader}></div>}
        data={resourcesDetail && resourcesDetail}
        columns={meshmodel_columns}
        options={meshmodel_options}
        classes={classes}
      />
    </div>
  );
};

export default DuplicatesDataTable;
