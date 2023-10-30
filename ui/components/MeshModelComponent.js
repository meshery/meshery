import { Button, withStyles } from '@material-ui/core';
import { withSnackbar } from 'notistack';
import React, { useState, useEffect } from 'react';
import { Paper } from '@material-ui/core';
import UploadIcon from '@mui/icons-material/Upload';
import DoNotDisturbOnIcon from '@mui/icons-material/DoNotDisturbOn';
import {
  getComponentsDetailWithPageSize,
  getMeshModels,
  getRelationshipsDetailWithPageSize,
  getComponentFromModelApi,
  getRelationshipFromModelApi,
  searchModels,
  searchComponents,
  getMeshModelRegistrants,
} from '../api/meshmodel';
import { OVERVIEW, MODELS, COMPONENTS, RELATIONSHIPS, REGISTRANTS } from '../constants/navigator';
import { SORT } from '../constants/endpoints';
import useStyles from '../assets/styles/general/tool.styles';
import { Colors } from '../themes/app';
import MesheryTreeView from './MesheryTreeView';
// import { useGetMeshModelQuery } from '../rtk-query/meshModel';

//TODO : This Should derive the indices of rendered rows
// const ROWS_INDICES = {
//   KIND : 0 ,
//   VERSION : 1 ,
//   MODEL : 3 ,
// }
const meshmodelStyles = (theme) => ({
  wrapperClss: {
    flexGrow: 1,
    maxWidth: '100%',
    height: 'auto',
  },
  tab: {
    minWidth: 40,
    paddingLeft: 0,
    paddingRight: 0,
    '&.Mui-selected': {
      color: theme.palette.view === 'dark' ? '#00B39F' : theme.palette.primary,
    },
  },
  tabs: {
    '& .MuiTabs-indicator': {
      backgroundColor: theme.palette.view === 'dark' ? '#00B39F' : theme.palette.primary,
    },
  },
  dashboardSection: {
    padding: theme.spacing(2),
    borderRadius: 4,
    height: '100%',
    overflowY: 'scroll',
  },
  duplicatesModelStyle: {
    backgroundColor: theme.palette.type === 'dark' ? '#00B39F' : theme.palette.primary,
  },
});

const MeshModelComponent = ({ modelsCount, componentsCount, relationshipsCount }) => {
  const [resourcesDetail, setResourcesDetail] = useState([]);
  const [isRequestCancelled, setRequestCancelled] = useState(false);
  const [, setCount] = useState();
  const [page] = useState(0);
  const [searchText, setSearchText] = useState(null);
  const [rowsPerPage] = useState(10);
  const [sortOrder, setSortOrder] = useState({
    sort: SORT.ASCENDING,
    order: '',
  });
  // const [checked, setChecked] = useState(false);
  const StyleClass = useStyles();
  const [view, setView] = useState(OVERVIEW);
  const [convert, setConvert] = useState(false);
  const [show, setShow] = useState({
    model: {},
    components: [],
    relationships: [],
  });
  const [comp, setComp] = useState({});
  const [rela, setRela] = useState({});
  const [animate, setAnimate] = useState(false);
  const [regi, setRegi] = useState({});
  // const modeldata = useGetMeshModelQuery({
  //   page: 1,
  //   pageSize: 10,
  // });

  const getModels = async (page) => {
    try {
      const { models } = await getMeshModels(page + 1, rowsPerPage); // page+1 due to server side indexing starting from 1
      const componentPromises = models.map(async (model) => {
        const { components } = await getComponentFromModelApi(model.name);
        const { relationships } = await getRelationshipFromModelApi(model.name);
        model.components = components;
        model.relationships = relationships;
      });

      await Promise.all(componentPromises);
      if (!isRequestCancelled) {
        setResourcesDetail(models);
      }
    } catch (error) {
      console.error('Failed to fetch models:', error);
    }
  };

  const getComponents = async (page, sortOrder) => {
    // if (typeof sortOrder === "undefined" || sortOrder === null) {
    //   setSortOrder("");
    // }

    try {
      const { total_count, components } = await getComponentsDetailWithPageSize(
        page + 1,
        rowsPerPage,
        sortOrder.sort,
        sortOrder.order,
      ); // page+1 due to server side indexing starting from 1
      setCount(total_count);
      if (!isRequestCancelled) {
        setResourcesDetail(components);
        setSortOrder(sortOrder);
      }
    } catch (error) {
      console.error('Failed to fetch components:', error);
    }
  };

  const getRelationships = async (page, sortOrder) => {
    // if (typeof sortOrder === "undefined" || sortOrder === null) {
    //   setSortOrder("");
    // }

    try {
      const { total_count, relationships } = await getRelationshipsDetailWithPageSize(
        page + 1,
        rowsPerPage,
        sortOrder.sort,
        sortOrder.order,
      );
      setCount(total_count);
      if (!isRequestCancelled) {
        setResourcesDetail(relationships);
        setSortOrder(sortOrder);
      }
    } catch (error) {
      console.error('Failed to fetch relationships:', error);
    }
  };

  const getSearchedModels = async (searchText) => {
    try {
      const { total_count, models } = await searchModels(searchText);
      const componentPromises = models.map(async (model) => {
        const { components } = await getComponentFromModelApi(model.name);
        const { relationships } = await getRelationshipFromModelApi(model.name);
        model.components = components;
        model.relationships = relationships;
      });

      await Promise.all(componentPromises);
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

  const getRegistrants = async (page) => {
    try {
      const { total_count, registrants } = await getMeshModelRegistrants(page, rowsPerPage);
      setCount(total_count);
      if (!isRequestCancelled) {
        setResourcesDetail(registrants);
      }
    } catch (error) {
      console.error('Failed to fetch registrants:', error);
    }
  };

  // const handleToggleDuplicates = () => {
  //   setChecked(!checked);
  // };

  // const filteredData = checked
  //   ? resourcesDetail // Show all data, including duplicates
  //   : resourcesDetail.filter((item, index, self) => {
  //       // Filter out duplicates based on your criteria (e.g., name and version)
  //       return (
  //         index ===
  //         self.findIndex(
  //           (otherItem) => item.name === otherItem.name && item.version === otherItem.version,
  //         )
  //       );
  //     });
  const filteredData = resourcesDetail;

  useEffect(() => {
    setRequestCancelled(false);

    // if(view === MODELS){
    //   getModels(page);
    // } else if (view === COMPONENTS){
    //   getComponents(page,sortOrder);
    // } else if (view === RELATIONSHIPS){
    //   getRelationships(page,sortOrder);
    // } else if (view === REGISTRANTS){
    //   getRegistrants(page);
    // }

    if (view === MODELS && searchText === null) {
      getModels(page);
    } else if (view === COMPONENTS && searchText === null) {
      getComponents(page, sortOrder);
    } else if (view === RELATIONSHIPS) {
      getRelationships(page, sortOrder);
    } else if (view === MODELS && searchText) {
      getSearchedModels(searchText);
    } else if (view === COMPONENTS && searchText) {
      getSearchedComponents(searchText);
    } else if (view === REGISTRANTS && searchText === null) {
      getRegistrants(page);
    }

    return () => {
      setRequestCancelled(true);
    };
  }, [view, page, searchText, rowsPerPage, show, comp, rela, regi]);

  // const meshmodel_columns = [
  //   {
  //     name: view === COMPONENTS || view === RELATIONSHIPS ? 'kind' : 'displayName',
  //     label: `Name`,
  //     options: {
  //       sort: view === COMPONENTS || view === RELATIONSHIPS ? true : false,
  //       searchable: view === RELATIONSHIPS ? false : true,
  //       customHeadRender: function CustomHead({ index, ...column }, sortColumn) {
  //         return (
  //           <TableCell align={'start'} key={index} onClick={() => sortColumn(index)}>
  //             <TableSortLabel
  //               active={column.sortDirection != null}
  //               direction={column.sortDirection || 'asc'}
  //             >
  //               <b>{column.label}</b>
  //             </TableSortLabel>
  //           </TableCell>
  //         );
  //       },
  //       customBodyRender: (value) => (
  //         <Tooltip title={value} placement="top">
  //           <div>{value}</div>
  //         </Tooltip>
  //       ),
  //     },
  //   },
  //   {
  //     name: view === COMPONENTS || view === RELATIONSHIPS ? 'apiVersion' : 'version',
  //     label: view === COMPONENTS || view === RELATIONSHIPS ? 'Api Version' : 'Version',
  //     options: {
  //       sort: false,
  //       searchable: view === RELATIONSHIPS ? false : true,
  //       customHeadRender: function CustomHead({ index, ...column }) {
  //         return (
  //           <TableCell align={'start'} key={index}>
  //             <b>{column.label}</b>
  //           </TableCell>
  //         );
  //       },
  //       customBodyRender: (value) => (
  //         <Tooltip title={value} placement="top">
  //           <div>{value}</div>
  //         </Tooltip>
  //       ),
  //     },
  //   },
  //   {
  //     name: 'category',
  //     label: 'Category Name',
  //     options: {
  //       sort: false,
  //       display: view === MODELS ? 'true' : 'false',
  //       searchable: true,
  //       customHeadRender: function CustomHead({ index, ...column }) {
  //         return (
  //           <TableCell align={'start'} key={index}>
  //             <b>{column.label}</b>
  //           </TableCell>
  //         );
  //       },
  //       customBodyRender: (value) => {
  //         if (!(view === RELATIONSHIPS || view === COMPONENTS)) {
  //           const { modelDisplayName, name } = value;
  //           return (
  //             <Tooltip title={view === MODELS ? name : modelDisplayName} placement="top">
  //               <div>{view === MODELS ? name : modelDisplayName}</div>
  //             </Tooltip>
  //           );
  //         }
  //       },
  //     },
  //   },
  //   {
  //     name: 'metadata',
  //     label: 'Model',
  //     options: {
  //       sort: false,
  //       display: view === COMPONENTS ? 'true' : 'false',
  //       searchable: true,
  //       customHeadRender: function CustomHead({ index, ...column }) {
  //         return (
  //           <TableCell align={'start'} key={index}>
  //             <b>{column.label}</b>
  //           </TableCell>
  //         );
  //       },
  //       customBodyRender: (value) => {
  //         if (!(view === MODELS || view === RELATIONSHIPS)) {
  //           const { modelDisplayName } = value;
  //           return (
  //             <Tooltip title={modelDisplayName} placement="top">
  //               <div>{modelDisplayName}</div>
  //             </Tooltip>
  //           );
  //         }
  //       },
  //     },
  //   },
  //   {
  //     name: 'metadata',
  //     label: 'Sub Category',
  //     options: {
  //       sort: false,
  //       display: view === COMPONENTS ? 'true' : 'false',
  //       searchable: true,
  //       customHeadRender: function CustomHead({ index, ...column }) {
  //         return (
  //           <TableCell align={'start'} key={index}>
  //             <b>{column.label}</b>
  //           </TableCell>
  //         );
  //       },
  //       customBodyRender: (value) => {
  //         const { subCategory } = value;
  //         return (
  //           <Tooltip title={subCategory} placement="top">
  //             <div>{subCategory}</div>
  //           </Tooltip>
  //         );
  //       },
  //     },
  //   },
  //   {
  //     name: 'model',
  //     label: 'Model',
  //     options: {
  //       sort: false,
  //       display: view === RELATIONSHIPS ? 'true' : 'false',
  //       searchable: false,
  //       customHeadRender: function CustomHead({ index, ...column }) {
  //         return (
  //           <TableCell align={'start'} key={index}>
  //             <b>{column.label}</b>
  //           </TableCell>
  //         );
  //       },
  //       customBodyRender: (value) => {
  //         if (view === RELATIONSHIPS) {
  //           const { displayName } = value;
  //           return (
  //             <Tooltip title={displayName} placement="top">
  //               <div>{displayName}</div>
  //             </Tooltip>
  //           );
  //         }
  //       },
  //     },
  //   },
  //   {
  //     name: 'duplicates',
  //     label: 'Duplicates',
  //     options: {
  //       sort: false,
  //       searchable: true,
  //       customHeadRender: function CustomHead({ index, ...column }) {
  //         if (view !== RELATIONSHIPS)
  //           return (
  //             <TableCell align={'start'} key={index}>
  //               <b>{column.label}</b>
  //             </TableCell>
  //           );
  //       },
  //       customBodyRender: (value) => {
  //         if (view !== RELATIONSHIPS)
  //           return (
  //             <Tooltip title={value} placement="top">
  //               <div>{value}</div>
  //             </Tooltip>
  //           );
  //       },
  //     },
  //   },
  // ];

  // const meshmodel_options = {
  //   rowsPerPage: rowsPerPage,
  //   rowsPerPageOptions: [10, 25],
  //   page: page,
  //   count: count,
  //   sort: true,
  //   search: false,
  //   viewColumns: false,
  //   download: false,
  //   print: false,
  //   filter: false,
  //   selectableRows: false,
  //   // search: view === RELATIONSHIPS ? false : true,
  //   serverSide: true,
  //   // expandableRows : (view !== RELATIONSHIPS && checked === true) && true,
  //   onChangePage: debounce((p) => setPage(p), 200),
  //   onSearchChange: debounce((searchText) => setSearchText(searchText)),
  //   onChangeRowsPerPage: debounce((rowsPerPage) => {
  //     setRowsPerPage(rowsPerPage);
  //     setPage(0);
  //   }),
  //   onTableChange: (action, tableState) => {
  //     const sortInfo = tableState.announceText ? tableState.announceText.split(' : ') : [];
  //     let order = {
  //       sort: '',
  //       order: '',
  //     };

  //     if (tableState.activeColumn || tableState.activeColumn === 0) {
  //       //order = `${meshmodel_columns[tableState.activeColumn].name} desc`;
  //       order = {
  //         order: meshmodel_columns[tableState.activeColumn].name,
  //         sort: SORT.ASCENDING,
  //       };
  //       console.log('name', meshmodel_columns[tableState.activeColumn].name);
  //       switch (action) {
  //         case 'sort':
  //           if (sortInfo.length == 2) {
  //             if (sortInfo[1] === 'ascending') {
  //               order.sort = SORT.ASCENDING;
  //             } else {
  //               order.sort = SORT.DESCENDING;
  //             }
  //           }

  //           if (
  //             order !== sortOrder &&
  //             view === COMPONENTS &&
  //             meshmodel_columns[tableState.activeColumn].name === 'kind'
  //           ) {
  //             getComponents(page, order);
  //           }

  //           if (
  //             order !== sortOrder &&
  //             view === RELATIONSHIPS &&
  //             meshmodel_columns[tableState.activeColumn].name === 'kind'
  //           ) {
  //             getRelationships(page, order);
  //           }

  //           break;

  //         case 'default':
  //           break;
  //       }
  //     }
  //   },
  //   // renderExpandableRow : (rowData) => {
  //   //   //TODO: Index The data by id and then extract directly from api resp rather than component props
  //   //   const data = {
  //   //     kind : rowData[ROWS_INDICES.KIND]?.props?.children?.props?.children,
  //   //     model : rowData[ROWS_INDICES.MODEL]?.props?.children?.props?.children,
  //   //     version : rowData[ROWS_INDICES.VERSION]?.props?.children?.props?.children,
  //   //   }
  //   //   return (
  //   //     rowData[6].props.children.props.children > 0 ? (
  //   //       <TableCell
  //   //         colSpan={6}
  //   //         sx={{
  //   //           padding : "0.5rem",
  //   //           backgroundColor : "rgba(0, 0, 0, 0.05)"
  //   //         }}
  //   //       >
  //   //         <Grid
  //   //           container
  //   //           xs={12}
  //   //           spacing={1}
  //   //           sx={{
  //   //             margin : "auto",
  //   //             backgroundColor : "#f3f1f1",
  //   //             paddingLeft : "0.5rem",
  //   //             borderRadius : "0.25rem",
  //   //             width : "inherit"
  //   //           }}
  //   //         >
  //   //           <DuplicatesDataTable
  //   //             view={view}
  //   //             rowData={data}
  //   //             classes={classes}
  //   //           >
  //   //           </DuplicatesDataTable>
  //   //         </Grid>
  //   //       </TableCell>
  //   //     ) : (
  //   //       <TableCell
  //   //         colSpan={6}
  //   //         sx={{
  //   //           padding : "0.5rem",
  //   //         }}
  //   //       >
  //   //         <Grid
  //   //           container
  //   //           spacing={1}
  //   //           sx={{
  //   //             justifyContent : "center",
  //   //             margin : "auto",
  //   //             paddingLeft : "0.5rem",
  //   //             borderRadius : "0.25rem",
  //   //           }}
  //   //         >
  //   //           <b>No duplicates found</b>
  //   //         </Grid>
  //   //       </TableCell>
  //   //     )
  //   //   );
  //   // },
  // };

  // const [tableCols, updateCols] = useState(meshmodel_columns);

  // const [columnVisibility, setColumnVisibility] = useState(() => {
  //   // let showCols = updateVisibleColumns(colViews, width);
  //   // Initialize column visibility based on the original columns' visibility
  //   const initialVisibility = {};
  //   meshmodel_columns.forEach((col) => {
  //     initialVisibility[col.name] = col.options?.display !== false;
  //   });
  //   return initialVisibility;
  // });

  const customInlineStyle = {
    marginBottom: '0.5rem',
    marginTop: '1rem',
  };

  // const cardStyle = {
  //   background: '#51636B',
  //   color: 'white',
  //   height: '10rem',
  //   width: '13rem',
  //   display: 'flex',
  //   justifyContent: 'center',
  //   alignItems: 'center',
  //   marginRight: '1rem',
  //   flexDirection: 'column',
  //   cursor: 'pointer',
  // };

  return (
    <div data-test="workloads">
      {!convert && (
        <div className={StyleClass.toolWrapper} style={customInlineStyle}>
          {/* {view !== RELATIONSHIPS && (
            <FormControlLabel
              control={
                <Switch
                  color="primary"
                  checked={checked}
                  onChange={handleToggleDuplicates}
                  inputProps={{ 'aria-label': 'controlled' }}
                />
              }
              label="Duplicates"
            />
          )} */}
          <Button
            disabled
            variant="contained"
            style={{ background: '#dddddd', color: 'white', marginRight: '1rem' }}
            size="large"
            startIcon={<UploadIcon />}
          >
            Import
          </Button>
          <Button
            disabled
            variant="contained"
            size="large"
            style={{ background: '#dddddd', color: 'white' }}
            startIcon={<DoNotDisturbOnIcon />}
          >
            Ignore
          </Button>
        </div>
      )}
      {/* <ResponsiveDataTable
        data={filteredData}
        columns={meshmodel_columns}
        options={meshmodel_options}
        tableCols={tableCols}
        updateCols={updateCols}
        columnVisibility={columnVisibility}
      /> */}
      {convert ? (
        <div className={StyleClass.mainContainer}>
          <div
            className={StyleClass.cardContainer}
            style={{
              display: 'flex',
              flexDirection: 'row',
              marginBottom: '1rem',
            }}
          >
            <Paper
              elevation={3}
              className={StyleClass.cardStyle}
              onClick={() => {
                setView(MODELS);
                setAnimate(true);
                setConvert(false);
              }}
            >
              <span style={{ fontWeight: 'bold', fontSize: '3rem' }}>{modelsCount}</span>
              Models
            </Paper>
            <Paper
              elevation={3}
              className={StyleClass.cardStyle}
              onClick={() => {
                setView(COMPONENTS);
                setAnimate(true);
                setConvert(false);
              }}
            >
              <span style={{ fontWeight: 'bold', fontSize: '3rem' }}>{componentsCount}</span>
              Components
            </Paper>
            <Paper
              elevation={3}
              className={StyleClass.cardStyle}
              onClick={() => {
                setView(RELATIONSHIPS);
                setAnimate(true);
                setConvert(false);
              }}
            >
              <span style={{ fontWeight: 'bold', fontSize: '3rem' }}>{relationshipsCount}</span>
              Relationships
            </Paper>
            <Paper
              elevation={3}
              className={StyleClass.cardStyle}
              onClick={() => {
                setView(REGISTRANTS);
                setAnimate(true);
                setConvert(false);
              }}
            >
              <span style={{ fontWeight: 'bold', fontSize: '3rem' }}>1</span>
              Registrants
            </Paper>
          </div>
        </div>
      ) : (
        <div style={{ backgroundColor: 'white', borderRadius: '6px' }}>
          <div
            style={{
              backgroundColor: '#51636B',
              borderRadius: '6px 6px 0px 0px',
              color: 'white',
              paddingTop: '1rem',
              paddingLeft: '1rem',
              paddingRight: '1rem',
              display: 'flex',
              flexDirection: 'row',
            }}
          >
            <Button
              style={{
                backgroundColor: view === OVERVIEW ? 'white' : '#677a84',
                color: view === OVERVIEW ? 'black' : 'white',
                borderRadius: '8px 8px 0px 0px',
                marginRight: '1rem',
                padding: view === OVERVIEW ? '0.5rem 2rem' : '0.6rem 2rem',
              }}
              onClick={() => {
                setView(OVERVIEW);
                setConvert(true);
                setAnimate(false);
              }}
            >
              Overview
            </Button>
            <Button
              style={{
                backgroundColor: view === MODELS ? 'white' : '#677a84',
                color: view === MODELS ? 'black' : 'white',
                borderRadius: '8px 8px 0px 0px',
                marginRight: '1rem',
                padding: view === MODELS ? '0.5rem 2rem' : '0.6rem 2rem',
              }}
              onClick={() => {
                setView(MODELS);
                setSearchText(null);
              }}
            >
              Models ({modelsCount})
            </Button>
            <Button
              style={{
                backgroundColor: view === COMPONENTS ? 'white' : '#677a84',
                color: view === COMPONENTS ? 'black' : 'white',
                borderRadius: '8px 8px 0px 0px',
                marginRight: '1rem',
                padding: view === COMPONENTS ? '0.5rem 2rem' : '0.6rem 2rem',
              }}
              onClick={() => {
                setView(COMPONENTS);
                setSearchText(null);
              }}
            >
              Components ({componentsCount})
            </Button>
            <Button
              style={{
                backgroundColor: view === RELATIONSHIPS ? 'white' : '#677a84',
                color: view === RELATIONSHIPS ? 'black' : 'white',
                borderRadius: '8px 8px 0px 0px',
                marginRight: '1rem',
                padding: view === RELATIONSHIPS ? '0.5rem 2rem' : '0.6rem 2rem',
              }}
              onClick={() => {
                setView(RELATIONSHIPS);
                setSearchText(null);
              }}
            >
              Relationships ({relationshipsCount})
            </Button>
            <Button
              style={{
                backgroundColor: view === REGISTRANTS ? 'white' : '#677a84',
                color: view === REGISTRANTS ? 'black' : 'white',
                borderRadius: '8px 8px 0px 0px',
                marginRight: '1rem',
                padding: view === REGISTRANTS ? '0.5rem 2rem' : '0.6rem 2rem',
              }}
              onClick={() => {
                setView(REGISTRANTS);
                setSearchText(null);
              }}
            >
              Registrants (1)
            </Button>
          </div>
          <div className={StyleClass.treeWrapper}>
            <div style={{ height: '30rem', width: '50%', margin: '1rem' }}>
              <MesheryTreeView
                data={filteredData}
                view={view}
                show={show}
                setShow={setShow}
                comp={comp}
                setComp={setComp}
                rela={rela}
                setRela={setRela}
                regi={regi}
                setRegi={setRegi}
                setSearchText={setSearchText}
              />
            </div>
            <div
              className={
                (view === MODELS && !show.model.displayName) ||
                (view === COMPONENTS && !comp.displayName) ||
                (view === RELATIONSHIPS && !rela.kind) ||
                (view === REGISTRANTS && !regi.hostname)
                  ? StyleClass.emptyDetailsContainer
                  : StyleClass.detailsContainer
              }
            >
              {((view === MODELS && !show.model.displayName) ||
                (view === COMPONENTS && !comp.displayName) ||
                (view === RELATIONSHIPS && !rela.kind) ||
                (view === REGISTRANTS && !regi.hostname)) && (
                <p style={{ color: '#969696' }}>No {view} selected</p>
              )}
              {view === MODELS && (
                <>
                  {show.model.displayName && (
                    <div>
                      <p
                        style={{
                          marginTop: '0px',
                          fontSize: '20px',
                          fontWeight: 'bold',
                        }}
                      >
                        {show.model.displayName}
                      </p>
                      <div
                        style={{
                          display: 'flex',
                          flexDirection: 'row',
                        }}
                      >
                        <div
                          style={{
                            display: 'flex',
                            flexDirection: 'column',
                            width: '50%',
                            paddingRight: '1rem',
                          }}
                        >
                          <p
                            style={{
                              padding: '0',
                              margin: '0',
                              fontSize: '16px',
                              fontWeight: '600',
                            }}
                          >
                            Version
                          </p>
                          <p style={{ padding: '0', margin: '0', fontSize: '14px' }}>
                            {show.model.version}
                          </p>
                          <p
                            style={{
                              padding: '0',
                              margin: '0',
                              marginTop: '12px',
                              fontSize: '16px',
                              fontWeight: '600',
                            }}
                          >
                            Registrant
                          </p>
                          <p style={{ padding: '0', margin: '0', fontSize: '14px' }}>
                            {show.model.hostname}
                          </p>
                          <p
                            style={{
                              padding: '0',
                              margin: '0',
                              marginTop: '12px',
                              fontSize: '16px',
                              fontWeight: '600',
                            }}
                          >
                            Components
                          </p>
                          <p style={{ padding: '0', margin: '0', fontSize: '14px' }}>
                            {show.model.components === null ? '0' : show.model.components.length}
                          </p>
                        </div>
                        <div style={{ display: 'flex', flexDirection: 'column', width: '50%' }}>
                          <p
                            style={{
                              padding: '0',
                              margin: '0',
                              fontSize: '16px',
                              fontWeight: '600',
                            }}
                          >
                            Category
                          </p>
                          <p style={{ padding: '0', margin: '0', fontSize: '14px' }}>
                            {show.model.category?.name}
                          </p>
                          <p
                            style={{
                              padding: '0',
                              margin: '0',
                              marginTop: '12px',
                              fontSize: '16px',
                              fontWeight: '600',
                            }}
                          >
                            Duplicates
                          </p>
                          <p style={{ padding: '0', margin: '0', fontSize: '14px' }}>
                            {show.model.duplicates}
                          </p>
                          <p
                            style={{
                              padding: '0',
                              margin: '0',
                              marginTop: '12px',
                              fontWeight: '600',
                              fontSize: '16px',
                            }}
                          >
                            Relationships
                          </p>
                          <p style={{ padding: '0', margin: '0', fontSize: '14px' }}>
                            {show.model.relationships === null
                              ? '0'
                              : show.model.relationships.length}
                          </p>
                        </div>
                      </div>
                    </div>
                  )}
                  {show.components.length !== 0 && (
                    <div>
                      <hr style={{ marginTop: '1rem' }} />
                      {show.components.map((component) => (
                        <div>
                          <p
                            style={{
                              fontSize: '20px',
                              margin: '0',
                              fontWeight: 'bold',
                            }}
                          >
                            {component.displayName}
                          </p>
                          <div
                            style={{
                              display: 'flex',
                              flexDirection: 'row',
                            }}
                          >
                            <div
                              style={{
                                display: 'flex',
                                flexDirection: 'column',
                                width: '50%',
                                paddingRight: '1rem',
                              }}
                            >
                              <p
                                style={{
                                  padding: '0',
                                  margin: '0',
                                  fontSize: '16px',
                                  fontWeight: '600',
                                }}
                              >
                                API Version
                              </p>
                              <p style={{ padding: '0', margin: '0', fontSize: '14px' }}>
                                {component.apiVersion}
                              </p>
                            </div>
                            <div style={{ display: 'flex', flexDirection: 'column', width: '50%' }}>
                              <p
                                style={{
                                  padding: '0',
                                  margin: '0',
                                  fontSize: '16px',
                                  fontWeight: '600',
                                }}
                              >
                                Sub Category
                              </p>
                              <p style={{ padding: '0', margin: '0', fontSize: '14px' }}>
                                {component.kind}
                              </p>
                            </div>
                          </div>
                        </div>
                      ))}
                    </div>
                  )}
                  {show.relationships.length !== 0 && (
                    <div>
                      <hr style={{ marginTop: '1rem' }} />
                      {show.relationships.map((rela) => (
                        <div>
                          <p
                            style={{
                              fontSize: '20px',
                              margin: '0',
                              fontWeight: 'bold',
                            }}
                          >
                            {rela.kind}
                          </p>
                          <div
                            style={{
                              display: 'flex',
                              flexDirection: 'row',
                            }}
                          >
                            <div
                              style={{
                                display: 'flex',
                                flexDirection: 'column',
                                width: '50%',
                                paddingRight: '1rem',
                              }}
                            >
                              <p
                                style={{
                                  padding: '0',
                                  margin: '0',
                                  fontSize: '16px',
                                  fontWeight: '600',
                                }}
                              >
                                API Version
                              </p>
                              <p style={{ padding: '0', margin: '0', fontSize: '14px' }}>
                                {rela.apiVersion}
                              </p>
                            </div>
                            <div style={{ display: 'flex', flexDirection: 'column', width: '50%' }}>
                              <p
                                style={{
                                  padding: '0',
                                  margin: '0',
                                  fontSize: '16px',
                                  fontWeight: '600',
                                }}
                              >
                                Sub Type
                              </p>
                              <p style={{ padding: '0', margin: '0', fontSize: '14px' }}>
                                {rela.subType}
                              </p>
                            </div>
                          </div>
                        </div>
                      ))}
                    </div>
                  )}
                </>
              )}
              {view === COMPONENTS && comp.displayName && (
                <div>
                  <p
                    style={{
                      fontSize: '20px',
                      fontWeight: 'bold',
                      marginTop: '0',
                    }}
                  >
                    {comp.displayName}
                  </p>
                  <div
                    style={{
                      display: 'flex',
                      flexDirection: 'row',
                    }}
                  >
                    <div
                      style={{
                        display: 'flex',
                        flexDirection: 'column',
                        width: '50%',
                        paddingRight: '1rem',
                      }}
                    >
                      <p style={{ padding: '0', margin: '0', fontSize: '16px', fontWeight: '600' }}>
                        API Version
                      </p>
                      <p style={{ padding: '0', margin: '0', fontSize: '14px' }}>
                        {comp.apiVersion}
                      </p>
                      <p
                        style={{
                          padding: '0',
                          margin: '0',
                          marginTop: '12px',
                          fontSize: '16px',
                          fontWeight: '600',
                        }}
                      >
                        Model Name
                      </p>
                      <p style={{ padding: '0', margin: '0', fontSize: '14px' }}>
                        {comp.model?.displayName}
                      </p>
                      <p
                        style={{
                          padding: '0',
                          margin: '0',
                          marginTop: '12px',
                          fontSize: '16px',
                          fontWeight: '600',
                        }}
                      >
                        Kind
                      </p>
                      <p style={{ padding: '0', margin: '0', fontSize: '14px' }}>{comp.kind}</p>
                    </div>
                    <div style={{ display: 'flex', flexDirection: 'column', width: '50%' }}>
                      <p style={{ padding: '0', margin: '0', fontSize: '16px', fontWeight: '600' }}>
                        Sub Category
                      </p>
                      <p style={{ padding: '0', margin: '0', fontSize: '14px' }}>{comp.kind}</p>
                      <p
                        style={{
                          padding: '0',
                          margin: '0',
                          marginTop: '12px',
                          fontSize: '16px',
                          fontWeight: '600',
                        }}
                      >
                        Registrant
                      </p>
                      <p style={{ padding: '0', margin: '0', fontSize: '14px' }}>
                        {comp.displayhostname}
                      </p>
                      <p
                        style={{
                          padding: '0',
                          margin: '0',
                          marginTop: '12px',
                          fontSize: '16px',
                          fontWeight: '600',
                        }}
                      >
                        Duplicates
                      </p>
                      <p style={{ padding: '0', margin: '0', fontSize: '14px' }}>
                        {comp.duplicates}
                      </p>
                    </div>
                  </div>
                </div>
              )}
              {view === RELATIONSHIPS && rela.kind && (
                <div>
                  <p
                    style={{
                      fontSize: '20px',
                      marginTop: '0',
                      fontWeight: 'bold',
                    }}
                  >
                    {rela.kind}
                  </p>
                  <p style={{ fontWeight: '600', margin: '0' }}>Description</p>
                  <p style={{ margin: '0', fontSize: '14px' }}>{rela.metadata?.description}</p>
                  <div
                    style={{
                      display: 'flex',
                      marginTop: '12px',
                      flexDirection: 'row',
                    }}
                  >
                    <div
                      style={{
                        display: 'flex',
                        flexDirection: 'column',
                        width: '50%',
                        paddingRight: '1rem',
                      }}
                    >
                      <p style={{ padding: '0', margin: '0', fontSize: '16px', fontWeight: '600' }}>
                        API Version
                      </p>
                      <p style={{ padding: '0', margin: '0', fontSize: '14px' }}>
                        {rela.apiVersion}
                      </p>
                      <p
                        style={{
                          padding: '0',
                          margin: '0',
                          marginTop: '12px',
                          fontWeight: '600',
                          fontSize: '16px',
                        }}
                      >
                        Model Name
                      </p>
                      <p style={{ padding: '0', margin: '0', fontSize: '14px' }}>
                        {rela.model?.displayName}
                      </p>
                    </div>
                    <div style={{ display: 'flex', flexDirection: 'column', width: '50%' }}>
                      <p style={{ padding: '0', margin: '0', fontSize: '16px', fontWeight: '600' }}>
                        Sub Type
                      </p>
                      <p style={{ padding: '0', margin: '0', fontSize: '14px' }}>{rela.subType}</p>
                      <p
                        style={{
                          padding: '0',
                          margin: '0',
                          marginTop: '12px',
                          fontSize: '16px',
                          fontWeight: '600',
                        }}
                      >
                        Registrant
                      </p>
                      <p style={{ padding: '0', margin: '0', fontSize: '14px' }}>
                        {rela.displayhostname}
                      </p>
                    </div>
                  </div>
                </div>
              )}
              {view === REGISTRANTS && regi.hostname && (
                <div>
                  <p
                    style={{
                      fontSize: '20px',
                      marginTop: '0',
                      fontWeight: 'bold',
                    }}
                  >
                    {regi.hostname}
                  </p>
                  <div
                    style={{
                      display: 'flex',
                      flexDirection: 'row',
                    }}
                  >
                    <div
                      style={{
                        display: 'flex',
                        flexDirection: 'column',
                        width: '50%',
                        paddingRight: '1rem',
                      }}
                    >
                      <p style={{ padding: '0', margin: '0', fontSize: '16px', fontWeight: '600' }}>
                        Models
                      </p>
                      <p style={{ padding: '0', margin: '0', fontSize: '14px' }}>
                        {regi.summary?.models}
                      </p>
                      <p
                        style={{
                          padding: '0',
                          margin: '0',
                          marginTop: '12px',
                          fontSize: '16px',
                          fontWeight: '600',
                        }}
                      >
                        Components
                      </p>
                      <p style={{ padding: '0', margin: '0', fontSize: '14px' }}>
                        {regi.summary?.components}
                      </p>
                    </div>
                    <div style={{ display: 'flex', flexDirection: 'column', width: '50%' }}>
                      <p style={{ padding: '0', margin: '0', fontSize: '16px', fontWeight: '600' }}>
                        Relationships
                      </p>
                      <p style={{ padding: '0', margin: '0', fontSize: '14px' }}>
                        {regi.summary?.relationships}
                      </p>
                      <p
                        style={{
                          padding: '0',
                          margin: '0',
                          marginTop: '12px',
                          fontSize: '16px',
                          fontWeight: '600',
                        }}
                      >
                        Policies
                      </p>
                      <p style={{ padding: '0', margin: '0', fontSize: '14px' }}>
                        {regi.summary?.policies}
                      </p>
                    </div>
                  </div>
                </div>
              )}
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default withStyles(meshmodelStyles)(withSnackbar(MeshModelComponent));
