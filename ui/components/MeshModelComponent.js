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
      color : theme.palette.type === 'dark' ? "#00B39F" : theme.palette.primary,
    },
  },
  tabs : {
    "& .MuiTabs-indicator" : {
      backgroundColor : theme.palette.type === 'dark' ? "#00B39F" : theme.palette.primary,
    },
  },
  dashboardSection : {
    padding : theme.spacing(2), borderRadius : 4,height : "100%",overflowY : "scroll"
  },


})

const MeshModelComponent = ({ type, classes }) => {
  const [resourcesDetail, setResourcesDetail] = useState();
  const [isRequestCancelled, setRequestCancelled] = useState(false);
  const [page, setPage] = useState(1);

  const getModels = async (page) => {
    try {
      const data = await getModelsDetail(page);
      if (!isRequestCancelled) {
        setResourcesDetail(data);
      }
    } catch (error) {
      console.error('Failed to fetch models:', error);
    }
  };

  const getComponents = async (page) => {
    try {
      const data = await getComponentsDetail(page);
      if (!isRequestCancelled) {
        setResourcesDetail(data);
      }
    } catch (error) {
      console.error('Failed to fetch components:', error);
    }
  };

  useEffect(() => {

    setRequestCancelled(false);


    if (type === 'models') {
      getModels(page);
    } else if (type === 'components') {
      getComponents(page);
    }

    return () => {
      setRequestCancelled(true);
    };
  }, [type, page]);




  const meshmodel_columns = [
    {
      name : type === 'components' ? 'kind' : 'displayName',
      label : `Name`,
      options : {
        sort : true,
        search : true,
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
      name : type === 'components' ? 'apiVersion' : 'version',
      label : type === 'components' ? 'Api Version' : 'Version',
      options : {
        sort : true,
        search : true,
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
      name : type === 'components' ? 'metadata' : 'category',
      label : type === 'components' ? 'Model' : 'Category Name',
      options : {
        sort : false,
        search : true,
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
            <Tooltip title={type === 'models' ? name : modelDisplayName} placement="top">
              <div>{type === 'models' ? name : modelDisplayName}</div>
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
        search : true,
        display : type === 'components' ? 'true' : 'false',
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
    count : 500,
    sort : true,
    download : false,
    print : false,
    filter : false,
    selectableRows : false,
    search : true,
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

