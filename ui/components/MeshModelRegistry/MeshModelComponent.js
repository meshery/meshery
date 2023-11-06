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
  getMeshModelsByRegistrants,
} from '../../api/meshmodel';
import {
  OVERVIEW,
  MODELS,
  COMPONENTS,
  RELATIONSHIPS,
  REGISTRANTS,
} from '../../constants/navigator';
import { SORT } from '../../constants/endpoints';
import useStyles from '../../assets/styles/general/tool.styles';
import { Colors } from '../../themes/app';
import MesheryTreeView from './MesheryTreeView';
import MeshModelDetails from './MeshModelDetails';
import { toLower } from 'lodash';
// import { useGetMeshModelQuery, useGetLoggedInUserQuery } from '../rtk-query/meshModel';

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
  const [page, setPage] = useState({
    Models: 0,
    Components: 0,
    Relationships: 0,
    Registrants: 0,
  });

  const [searchText, setSearchText] = useState(null);
  const [rowsPerPage] = useState(14);
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
  const [checked, setChecked] = useState(true);
  // const [filteredData, setFilteredData] = useState([]);
  // const { data: modeldata } = useGetMeshModelQuery();

  const getModels = async (page) => {
    console.log('models:', page);
    try {
      const { models } = await getMeshModels(page?.Models + 1, rowsPerPage); // page+1 due to server side indexing starting from 1
      const componentPromises = models.map(async (model) => {
        const { components } = await getComponentFromModelApi(model.name);
        const { relationships } = await getRelationshipFromModelApi(model.name);
        model.components = components;
        model.relationships = relationships;
      });

      await Promise.all(componentPromises);
      if (!isRequestCancelled && models) {
        setResourcesDetail((prev) => [...prev, ...models]);
      }
    } catch (error) {
      console.error('Failed to fetch models:', error);
    }
  };

  const getComponents = async (page, sortOrder) => {
    // if (typeof sortOrder === "undefined" || sortOrder === null) {
    //   setSortOrder("");
    // }
    console.log('components:', page);
    try {
      const { total_count, components } = await getComponentsDetailWithPageSize(
        page?.Components + 1,
        rowsPerPage,
        sortOrder.sort,
        sortOrder.order,
      ); // page+1 due to server side indexing starting from 1
      setCount(total_count);
      if (!isRequestCancelled && components) {
        setResourcesDetail((prev) => [...prev, ...components]);
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
    console.log('relationships:', page);
    try {
      const { total_count, relationships } = await getRelationshipsDetailWithPageSize(
        page?.Relationships + 1,
        rowsPerPage,
        sortOrder.sort,
        sortOrder.order,
      );
      setCount(total_count);
      if (!isRequestCancelled && relationships) {
        setResourcesDetail((prev) => [...prev, ...relationships]);
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
    console.log('registrants:', page);
    try {
      // const { models } = await getMeshModels(page?.Models + 1, rowsPerPage); // page+1 due to server side indexing starting from 1
      // const componentPromises = models.map(async (model) => {
      //   const { components } = await getComponentFromModelApi(model.name);
      //   const { relationships } = await getRelationshipFromModelApi(model.name);
      //   model.components = components;
      //   model.relationships = relationships;
      // });
      const { total_count, registrants } = await getMeshModelRegistrants(
        page?.Registrants + 1,
        rowsPerPage,
      );
      let tempRegistrants = [];
      console.log('registrants val:', registrants);
      let registrantPromise = registrants.map(async (registrant) => {
        let hostname = toLower(registrant?.hostname);
        console.log('inside');
        console.log('registrant1', registrant);
        const { models } = await getMeshModelsByRegistrants(
          page?.Models + 1,
          rowsPerPage,
          hostname,
        ); // page+1 due to server side indexing starting from 1
        if (models) {
          const componentPromises = models.map(async (model) => {
            const { components } = await getComponentFromModelApi(model.name);
            const { relationships } = await getRelationshipFromModelApi(model.name);
            model.components = components;
            model.relationships = relationships;
          });

          await Promise.all(componentPromises);
          registrant.models = models;
          tempRegistrants.push(registrant);
        } else {
          tempRegistrants.push(registrant);
        }
      });
      await Promise.all(registrantPromise);
      setCount(total_count);
      if (!isRequestCancelled && registrants) {
        console.log('registrants', registrants);
        let tempResourcesDetail = [];
        tempRegistrants.map((registrant) => {
          console.log('registrant', registrant);
          let oldRegistrant = resourcesDetail.find(
            (resource) => resource?.hostname === registrant?.hostname,
          );
          if (oldRegistrant !== undefined) {
            let newModels = [...oldRegistrant.models, ...registrant.models];
            registrant.models = newModels;
          }

          console.log('newRegistrant', registrant);
          tempResourcesDetail.push(registrant);
        });

        console.log('tempResourcesDetail', tempResourcesDetail);
        setResourcesDetail(tempRegistrants);
      }
    } catch (error) {
      console.error('Failed to fetch registrants:', error);
    }
  };

  // const handleToggleDuplicates = () => {
  //   setChecked(!checked);
  // };

  let filteredData = checked
    ? resourcesDetail // Show all data, including duplicates
    : resourcesDetail.filter((item, index, self) => {
        // Filter out duplicates based on your criteria (e.g., name and version)
        return (
          index ===
          self.findIndex(
            (otherItem) => item.name === otherItem.name && item.version === otherItem.version,
          )
        );
      });
  // const filteredData = resourcesDetail;

  useEffect(
    () => {
      console.log('setting filteredData');
      filteredData = checked
        ? resourcesDetail // Show all data, including duplicates
        : resourcesDetail.filter((item, index, self) => {
            // Filter out duplicates based on your criteria (e.g., name and version)
            return (
              index ===
              self.findIndex(
                (otherItem) => item.name === otherItem.name && item.version === otherItem.version,
              )
            );
          });
    },
    resourcesDetail,
    checked,
  );

  useEffect(() => {
    console.log('view', view);
    console.log('searchText', searchText);
    setRequestCancelled(false);
    // console.log(page);

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
  }, [view, page, searchText, rowsPerPage]);

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

  // const customInlineStyle = {
  //   marginBottom: '0.5rem',
  //   marginTop: '1rem',
  //   transition: 'all 0.5s',
  //   opacity: !convert ? '1' : '0',
  // };

  return (
    <div data-test="workloads">
      {/* {!convert && ( */}
      <div className={`${StyleClass.toolWrapper} ${animate ? StyleClass.toolWrapperAnimate : ''}`}>
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
          style={{
            background: '#dddddd',
            color: 'white',
            visibility: `${animate ? 'visible' : 'hidden'}`,
          }}
          size="large"
          startIcon={<UploadIcon />}
        >
          Import
        </Button>
        <Button
          disabled
          variant="contained"
          size="large"
          style={{
            background: '#dddddd',
            color: 'white',
            visibility: `${animate ? 'visible' : 'hidden'}`,
          }}
          startIcon={<DoNotDisturbOnIcon />}
        >
          Ignore
        </Button>
      </div>
      {/* )} */}
      {/* <ResponsiveDataTable
        data={filteredData}
        columns={meshmodel_columns}
        options={meshmodel_options}
        tableCols={tableCols}
        updateCols={updateCols}
        columnVisibility={columnVisibility}
      /> */}
      <div
        className={`${StyleClass.mainContainer} ${animate ? StyleClass.mainContainerAnimate : ''}`}
      >
        <div
          className={`${StyleClass.innerContainer} ${
            animate ? StyleClass.innerContainerAnimate : ''
          }`}
        >
          <Paper
            elevation={3}
            className={`${StyleClass.overviewTab} ${animate ? StyleClass.overviewTabAnimate : ''}`}
            onClick={() => {
              setConvert(false);
              setAnimate(false);
              // setTimeout(() => {
              // }, 1000);
            }}
          >
            Overview
          </Paper>
          <Paper
            elevation={3}
            className={`${StyleClass.cardStyle} ${animate ? StyleClass.cardStyleAnimate : ''}`}
            style={{
              backgroundColor: `${view === MODELS && animate ? 'white' : ''}`,
              color: `${view === MODELS && animate ? 'black' : ''}`,
            }}
            onClick={() => {
              setView(MODELS);
              setPage({
                Models: 0,
                Components: 0,
                Relationships: 0,
                Registrants: 0,
              });
              if (view !== MODELS) {
                setSearchText(null);
                setResourcesDetail([]);
              }
              if (!animate) {
                setAnimate(true);
                setTimeout(() => {
                  setConvert(true);
                }, 1100);
              }
            }}
          >
            <span
              style={{
                fontWeight: `${animate ? 'normal' : 'bold'}`,
                fontSize: `${animate ? '1rem' : '3rem'}`,
                transition: 'all 0.3s',
                marginLeft: `${animate && '4px'}`,
              }}
            >
              {animate ? `(${modelsCount})` : `${modelsCount}`}
            </span>
            Models
          </Paper>
          <Paper
            elevation={3}
            className={`${StyleClass.cardStyle} ${animate ? StyleClass.cardStyleAnimate : ''}`}
            style={{
              backgroundColor: `${view === COMPONENTS && animate ? 'white' : ''}`,
              color: `${view === COMPONENTS && animate ? 'black' : ''}`,
            }}
            onClick={() => {
              setView(COMPONENTS);
              setPage({
                Models: 0,
                Components: 0,
                Relationships: 0,
                Registrants: 0,
              });
              if (view !== COMPONENTS) {
                setSearchText(null);
                setResourcesDetail([]);
              }
              if (!animate) {
                setAnimate(true);
                setTimeout(() => {
                  setConvert(true);
                }, 1100);
              }
            }}
          >
            <span
              style={{
                fontWeight: `${animate ? 'normal' : 'bold'}`,
                fontSize: `${animate ? '1rem' : '3rem'}`,
                transition: 'all 0.3s',
                marginLeft: `${animate && '4px'}`,
              }}
            >
              {animate ? `(${componentsCount})` : `${componentsCount}`}
            </span>
            Components
          </Paper>
          <Paper
            elevation={3}
            className={`${StyleClass.cardStyle} ${animate ? StyleClass.cardStyleAnimate : ''}`}
            style={{
              backgroundColor: `${view === RELATIONSHIPS && animate ? 'white' : ''}`,
              color: `${view === RELATIONSHIPS && animate ? 'black' : ''}`,
            }}
            onClick={() => {
              setView(RELATIONSHIPS);
              setPage({
                Models: 0,
                Components: 0,
                Relationships: 0,
                Registrants: 0,
              });
              if (view !== RELATIONSHIPS) {
                setSearchText(null);
                setResourcesDetail([]);
              }
              if (!animate) {
                setAnimate(true);
                setTimeout(() => {
                  setConvert(true);
                }, 1100);
              }
            }}
          >
            <span
              style={{
                fontWeight: `${animate ? 'normal' : 'bold'}`,
                fontSize: `${animate ? '1rem' : '3rem'}`,
                transition: 'all 0.3s',
                marginLeft: `${animate && '4px'}`,
              }}
            >
              {animate ? `(${relationshipsCount})` : `${relationshipsCount}`}
            </span>
            Relationships
          </Paper>
          <Paper
            elevation={3}
            className={`${StyleClass.cardStyle} ${animate ? StyleClass.cardStyleAnimate : ''}`}
            style={{
              backgroundColor: `${view === REGISTRANTS && animate ? 'white' : ''}`,
              color: `${view === REGISTRANTS && animate ? 'black' : ''}`,
            }}
            onClick={() => {
              setView(REGISTRANTS);
              setPage({
                Models: 0,
                Components: 0,
                Relationships: 0,
                Registrants: 0,
              });
              if (view !== REGISTRANTS) {
                setSearchText(null);
                setResourcesDetail([]);
              }
              if (!animate) {
                setAnimate(true);
                setTimeout(() => {
                  setConvert(true);
                }, 1100);
              }
            }}
          >
            <span
              style={{
                fontWeight: `${animate ? 'normal' : 'bold'}`,
                fontSize: `${animate ? '1rem' : '3rem'}`,
                transition: 'all 0.3s',
                marginLeft: `${animate && '4px'}`,
              }}
            >
              {animate ? `(1)` : `1`}
            </span>
            Registrants
          </Paper>
        </div>
        {convert && (
          <div
            className={`${StyleClass.treeWrapper} ${convert ? StyleClass.treeWrapperAnimate : ''}`}
          >
            {console.log('filteredData', filteredData)}
            {console.log('view', view)}

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
                setPage={setPage}
                checked={checked}
                setChecked={setChecked}
              />
            </div>
            <MeshModelDetails view={view} show={show} rela={rela} regi={regi} comp={comp} />
          </div>
        )}
      </div>
    </div>
  );
};

export default withStyles(meshmodelStyles)(withSnackbar(MeshModelComponent));
