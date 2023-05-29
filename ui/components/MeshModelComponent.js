import {  withStyles } from '@material-ui/core'
import { withSnackbar } from 'notistack';
import React, { useState, useEffect } from 'react'
import MUIDataTable from 'mui-datatables';
import { TableCell, Tooltip, TableSortLabel } from '@material-ui/core';
import { getComponentsDetail, getModelsDetail, getRelationshipsDetail, searchModels, searchComponents } from '../api/meshmodel'
import debounce from '../utils/debounce';
import { MODELS, COMPONENTS, RELATIONSHIPS } from '../constants/navigator';

const meshmodelStyles = (theme) => ({
  wrapperClss : {
    flexGrow : 1,
    maxWidth : '100%',
    height : 'auto',
  },
  tab : {
    minWidth : 40,
    paddingLeft : 0,
    paddingRight : 0,
    "&.Mui-selected" : {
      color : theme.palette.view === 'dark' ? "#00B39F" : theme.palette.primary,
    },
  },
  tabs : {
    "& .MuiTabs-indicator" : {
      backgroundColor : theme.palette.view === 'dark' ? "#00B39F" : theme.palette.primary,
    },
  },
  dashboardSection : {
    padding : theme.spacing(2), borderRadius : 4,height : "100%",overflowY : "scroll"
  },


})

const MeshModelComponent = ({ view, classes }) => {
  const [resourcesDetail, setResourcesDetail] = useState();
  const [isRequestCancelled, setRequestCancelled] = useState(false);
  const [count, setCount] = useState();
  const [page, setPage] = useState(1);
  const [searchText, setSearchText] = useState(null);


  const getModels = async (page) => {
    try {
      const { total_count, models } = await getModelsDetail(page);
      setCount(total_count);
      if (!isRequestCancelled) {
        setResourcesDetail(models);
      }
    } catch (error) {
      console.error('Failed to fetch models:', error);
    }
  };

  const getComponents = async (page) => {
    try {
      const { total_count, components } = await getComponentsDetail(page);
      setCount(total_count);
      if (!isRequestCancelled) {
        setResourcesDetail(components);
      }
    } catch (error) {
      console.error('Failed to fetch components:', error);
    }
  };

  const getRelationships = async (page) => {
    try {
      const { total_count, relationships } = await getRelationshipsDetail(page);
      setCount(total_count);
      if (!isRequestCancelled) {
        setResourcesDetail(relationships);
      }
    } catch (error) {
      console.error('Failed to fetch relationships:', error);
    }
  };

  const getSearchedModels =  async (searchText) => {
    try {
      const { total_count, models } = await searchModels(searchText);
      setCount(total_count);
      if (!isRequestCancelled) {
        setResourcesDetail(models ? models : []);
      }
    } catch (error) {
      console.error('Failed to fetch components:', error);
    }

  };
  const getSearchedComponents = async (searchText) => {
    try {
      const { total_count, components } = await searchComponents(searchText);
      setCount(total_count);
      if (!isRequestCancelled) {
        setResourcesDetail(components ? components : []);
      }
    } catch (error) {
      console.error('Failed to fetch components:', error);
    }
  };

  useEffect(() => {
    setRequestCancelled(false);

    if (view === MODELS && searchText === null) {
      getModels(page);
    } else if (view === COMPONENTS && searchText === null) {
      getComponents(page);
    } else if (view === RELATIONSHIPS) {
      getRelationships(page);
    } else if (view === MODELS && searchText) {
      getSearchedModels(searchText);
    } else if (view === COMPONENTS && searchText) {
      getSearchedComponents(searchText);
    }

    return () => {
      setRequestCancelled(true);
    };
  }, [view, page, searchText]);

  const meshmodel_columns = [
    {
      name : (view === COMPONENTS || view === RELATIONSHIPS) ? 'kind' : 'displayName',
      label : `Name`,
      options : {
        sort : true,
        sortDescFirst : true,
        sortThirdClickReset : true,
        searchable : view === RELATIONSHIPS ? false : true,
        customHeadRender : function CustomHead({ index, ...column }, sortColumn) {
          return (
            <TableCell align={"start"} key={index} onClick={() => sortColumn(index)}>
              <TableSortLabel active={column.sortDirection != null} direction={column.sortDirection || "asc"}>
                <b>{column.label}</b>
              </TableSortLabel>
            </TableCell>
          );
        },
        customBodyRender : (value) => (
          <Tooltip title={value} placement="top">
            <div>{value}</div>
          </Tooltip>
        )
      },
    },
    {
      name : (view === COMPONENTS || view === RELATIONSHIPS) ? 'apiVersion' : 'version',
      label : (view === COMPONENTS || view === RELATIONSHIPS) ? 'Api Version' : 'Version',
      options : {
        sort : true,
        sortDescFirst : true,
        sortThirdClickReset : true,
        searchable : view === RELATIONSHIPS ?  false : true,
        customHeadRender : function CustomHead({ index, ...column }, sortColumn) {
          return (
            <TableCell align={"start"} key={index} onClick={() => sortColumn(index)}>
              <TableSortLabel active={column.sortDirection != null} direction={column.sortDirection || "asc"}>
                <b>{column.label}</b>
              </TableSortLabel>
            </TableCell>
          );
        },
        customBodyRender : (value) => (
          <Tooltip title={value} placement="top">
            <div>{value}</div>
          </Tooltip>
        ),
      },
    },
    {
      name : 'category',
      label : 'Category Name',
      options : {
        sort : false,
        display : view === MODELS ? 'true' : 'false',
        searchable : true,
        customHeadRender : function CustomHead({ index, ...column }) {
          return (
            <TableCell align={"start"} key={index}>
              <TableSortLabel>
                <b>{column.label}</b>
              </TableSortLabel>
            </TableCell>
          );
        },
        customBodyRender : (value) => {
          if (!(view === RELATIONSHIPS || view === COMPONENTS)) {
            const { modelDisplayName, name } = value;
            return (
              <Tooltip title={view === MODELS ? name : modelDisplayName} placement="top">
                <div>{view === MODELS ? name : modelDisplayName}</div>
              </Tooltip>
            )
          }

        },
      },
    },
    {
      name : 'metadata',
      label : 'Model',
      options : {
        sort : false,
        display : view === COMPONENTS ? 'true' : 'false',
        searchable : true,
        customHeadRender : function CustomHead({ index, ...column }) {
          return (
            <TableCell align={"start"} key={index}>
              <TableSortLabel>
                <b>{column.label}</b>
              </TableSortLabel>
            </TableCell>
          );
        },
        customBodyRender : (value) => {
          if (!(view === MODELS || view === RELATIONSHIPS)) {
            const { modelDisplayName } = value
            return (
              <Tooltip title={modelDisplayName} placement="top">
                <div>{modelDisplayName}</div>
              </Tooltip>
            )
          }
        },
      },
    },
    {
      name : 'metadata',
      label : 'Sub Category',
      options : {
        sort : false,
        display : view === COMPONENTS ? 'true' : 'false',
        searchable : true,
        customHeadRender : function CustomHead({ index, ...column }) {
          return (
            <TableCell align={"start"} key={index}>
              <TableSortLabel>
                <b>{column.label}</b>
              </TableSortLabel>
            </TableCell>
          );
        },
        customBodyRender : (value) => {
          const { subCategory } = value
          return (
            <Tooltip title={subCategory} placement="top">
              <div>{subCategory}</div>
            </Tooltip>
          )
        },
      },
    },
    {
      name : 'model',
      label : 'Model',
      options : {
        sort : false,
        display : view === RELATIONSHIPS ? 'true' : 'false',
        searchable : false,
        customHeadRender : function CustomHead({ index, ...column }) {
          return (
            <TableCell align={"start"} key={index}>
              <TableSortLabel>
                <b>{column.label}</b>
              </TableSortLabel>
            </TableCell>
          );
        },
        customBodyRender : (value) => {
          if (view === RELATIONSHIPS) {
            const { displayName } = value
            return (
              <Tooltip title={displayName} placement="top">
                <div>{displayName}</div>
              </Tooltip>
            )
          }
        },
      },
    },

  ]

  const meshmodel_options = {
    rowsPerPage : 25,
    count : count,
    sort : true,
    download : false,
    print : false,
    filter : false,
    selectableRows : false,
    search : view === RELATIONSHIPS ? false : true,
    serverSide : true,
    onChangePage : debounce((p) =>  setPage(p+1), 200),
    onSearchChange : debounce((searchText) => (setSearchText(searchText))),
  }

  return (
    <div >
      <div data-test="workloads">
        <MUIDataTable
          title={<div className={classes.tableHeader}></div>}
          data={resourcesDetail && resourcesDetail}
          columns={meshmodel_columns}
          options={meshmodel_options}
        />
      </div>
    </div>
  )
}

export default withStyles(meshmodelStyles)((withSnackbar(MeshModelComponent)));

