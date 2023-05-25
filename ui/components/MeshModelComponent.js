import {  withStyles } from '@material-ui/core'
import { withSnackbar } from 'notistack';
import React, { useState, useEffect } from 'react'
import MUIDataTable from 'mui-datatables';
import { TableCell, Tooltip, TableSortLabel } from '@material-ui/core';
import { getComponentsDetail, getModelsDetail } from '../api/meshmodel'

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

  useEffect(() => {

    setRequestCancelled(false);


    if (view === 'models') {
      getModels(page);
    } else if (view === 'components') {
      getComponents(page);
    }

    return () => {
      setRequestCancelled(true);
    };
  }, [view, page]);




  const meshmodel_columns = [
    {
      name : view === 'components' ? 'kind' : 'displayName',
      label : `Name`,
      options : {
        sort : true,
        sortDescFirst : true,
        sortThirdClickReset : true,
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
      name : view === 'components' ? 'apiVersion' : 'version',
      label : view === 'components' ? 'Api Version' : 'Version',
      options : {
        sort : true,
        sortDescFirst : true,
        sortThirdClickReset : true,
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
      name : view === 'components' ? 'metadata' : 'category',
      label : view === 'components' ? 'Model' : 'Category Name',
      options : {
        sort : false,
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
          const { modelDisplayName, name } = value;
          return (
            <Tooltip title={view === 'models' ? name : modelDisplayName} placement="top">
              <div>{view === 'models' ? name : modelDisplayName}</div>
            </Tooltip>
          )

        },
      },
    },
    {
      name : 'metadata',
      label : 'Sub Category',
      options : {
        sort : false,
        display : view === 'components' ? 'true' : 'false',
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


  ]

  const meshmodel_options = {
    rowsPerPage : 25,
    count : count,
    sort : true,
    download : false,
    print : false,
    filter : false,
    search : false,
    selectableRows : false,
    serverSide : true,
    onChangePage : (p) =>  setPage(p+1),
  }

  return (
    <div >
      <div data-test="workloads">
        {resourcesDetail && (
          <MUIDataTable
            title={<div className={classes.tableHeader}></div>}
            data={resourcesDetail}
            columns={meshmodel_columns}
            options={meshmodel_options}
          />
        )}
      </div>
    </div>
  )
}

export default withStyles(meshmodelStyles)((withSnackbar(MeshModelComponent)));

