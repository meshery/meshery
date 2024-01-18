import { withStyles } from '@material-ui/core';
import { withSnackbar } from 'notistack';
import React, { useState, useEffect } from 'react';
import { Paper } from '@material-ui/core';
import UploadIcon from '@mui/icons-material/Upload';
import DoNotDisturbOnIcon from '@mui/icons-material/DoNotDisturbOn';
import {
  OVERVIEW,
  MODELS,
  COMPONENTS,
  RELATIONSHIPS,
  REGISTRANTS,
  GRAFANA,
  PROMETHEUS,
} from '../../constants/navigator';
import { SORT } from '../../constants/endpoints';
import useStyles from '../../assets/styles/general/tool.styles';
import MesheryTreeView from './MesheryTreeView';
import MeshModelDetails from './MeshModelDetails';
import { toLower } from 'lodash';
import { DisableButton } from './MeshModel.style';
import CircularProgress from '@mui/material/CircularProgress';
import { Colors } from '../../themes/app';
import { useRouter } from 'next/router';
import { Provider } from 'react-redux';
import { store } from '../../store';
import { ErrorBoundary } from '../General/ErrorBoundary';
import {
  useLazyGetMeshModelsQuery,
  useLazyGetComponentsQuery,
  useLazyGetRelationshipsQuery,
  useLazyGetRegistrantsQuery,
} from '@/rtk-query/meshModel';
import NoSsr from '@material-ui/core/NoSsr';

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

const useMeshModelComponentRouter = () => {
  const router = useRouter();
  const { query } = router;

  const searchQuery = query.searchText || null;
  const selectedTab = query.tab === GRAFANA || query.tab === PROMETHEUS ? OVERVIEW : query.tab;
  const selectedPageSize = query.pagesize || 14;

  return { searchQuery, selectedTab, selectedPageSize };
};

const MeshModelComponent_ = ({
  modelsCount,
  componentsCount,
  relationshipsCount,
  registrantCount,
  settingsRouter,
}) => {
  const router = useRouter();
  const { handleChangeSelectedTab } = settingsRouter(router);
  const [resourcesDetail, setResourcesDetail] = useState([]);
  // const [isRequestCancelled, setRequestCancelled] = useState(false);
  const [, setCount] = useState();
  const { selectedTab, searchQuery, selectedPageSize } = useMeshModelComponentRouter();
  const [page, setPage] = useState({
    Models: 0,
    Components: 0,
    Relationships: 0,
    Registrants: 0,
  });
  const [searchText, setSearchText] = useState(searchQuery);
  const [rowsPerPage, setRowsPerPage] = useState(selectedPageSize);
  const [sortOrder] = useState({
    sort: SORT.ASCENDING,
    order: '',
  });
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
  // const [loading, setLoading] = useState(false);

  /**
   * RTK Lazy Queries
   */
  const [getMeshModelsData] = useLazyGetMeshModelsQuery();
  const [getComponentsData] = useLazyGetComponentsQuery();
  const [getRelationshipsData] = useLazyGetRelationshipsQuery();
  const [getRegistrantsData] = useLazyGetRegistrantsQuery();

  useEffect(() => {
    if (selectedTab && selectedTab !== OVERVIEW) {
      setAnimate(true);
      setConvert(true);
      setView(selectedTab);
    }
  }, [selectedTab]);

  const getMeshModels = () => {
    getMeshModelsData(
      {
        page: searchText ? 1 : page.Models + 1,
        pagesize: searchText ? 'all' : rowsPerPage,
        components: true,
        relationships: true,
        paginated: true,
        search: searchText ? searchText : '',
      },
      true,
    )
      .unwrap()
      .then((res) => {
        console.log('get response', res);
        setResourcesDetail((prev) => {
          if (searchText) {
            return [...res.models];
          } else {
            return [...prev, ...res.models];
          }
        });
        if (rowsPerPage !== 14) {
          setRowsPerPage(14);
        }
      })
      .catch((error) => {
        console.error('Failed to fetch models:', error.data);
      });
  };

  const getComponents = () => {
    getComponentsData(
      {
        page: searchText ? 1 : page.Components + 1,
        pagesize: searchText ? 'all' : rowsPerPage,
        search: searchText ? searchText : '',
        trim: false,
      },
      true,
    )
      .unwrap()
      .then((res) => {
        setCount(res.total_count);
        if (res?.components) {
          setResourcesDetail((prev) => {
            if (searchText) {
              return [...res.components];
            } else {
              return [...prev, ...res.components];
            }
          });
          // setSortOrder(sortOrder);
        }
        if (rowsPerPage !== 14) {
          setRowsPerPage(14);
        }
      })
      .catch((error) => {
        console.error('Failed to fetch components:', error.data);
      });
  };

  const getRelationships = async () => {
    getRelationshipsData(
      {
        page: searchText ? 1 : page.Relationships + 1,
        pagesize: searchText ? 'all' : rowsPerPage,
        paginated: true,
        search: searchText ? searchText : '',
      },
      true,
    )
      .unwrap()
      .then((res) => {
        setCount(res?.total_count);
        if (res?.relationships) {
          setResourcesDetail((prev) => {
            if (searchText) {
              return [...res.relationships];
            } else {
              return [...prev, ...res.relationships];
            }
          });
        }
        if (rowsPerPage !== 14) {
          setRowsPerPage(14);
        }
      })
      .catch((error) => {
        console.error('Failed to fetch relationships:', error.data);
      });
  };

  //TODO: Optimise it to support dynamic loading
  const getRegistrants = async () => {
    try {
      const { data: registrantsData } = await getRegistrantsData(
        {
          page: searchText ? 1 : page.Registrants + 1,
          pagesize: searchText ? 'all' : rowsPerPage,
          search: searchText ? searchText : '',
        },
        true,
      );

      let registrants = registrantsData.registrants;
      if (registrants) {
        let tempRegistrants = [];

        for (let registrant of registrants) {
          let hostname = toLower(registrant?.hostname);
          const { data: modelRes } = await getMeshModelsData(
            {
              page: page?.Models + 1,
              pagesize: rowsPerPage,
              registrant: hostname,
              components: true,
              relationships: true,
            },
            true,
          );

          if (modelRes.models) {
            const updatedRegistrant = {
              ...registrant,
              models: modelRes.models,
            };

            tempRegistrants.push(updatedRegistrant);
          } else {
            tempRegistrants.push(registrant);
          }
        }

        setCount(registrants.total_count);

        let tempResourcesDetail = [];
        tempRegistrants.forEach((registrant) => {
          let oldRegistrant = resourcesDetail.find(
            (resource) => resource?.hostname === registrant?.hostname,
          );

          if (oldRegistrant !== undefined) {
            let newModels = [...oldRegistrant.models, ...registrant.models];
            registrant.models = newModels;
          }

          tempResourcesDetail.push(registrant);
        });

        setResourcesDetail(tempRegistrants);
      }
      setRowsPerPage(14);
    } catch (error) {
      console.error('Failed to fetch registrants:', error);
    }
  };

  // TODO: This is wrong logic, backend should enforce this
  let filteredData =
    !view === MODELS
      ? resourcesDetail
      : checked
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

  useEffect(() => {
    filteredData =
      !view === MODELS
        ? resourcesDetail
        : checked
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
  }, [checked]);

  useEffect(() => {
    if (view === MODELS) {
      getMeshModels();
    } else if (view === COMPONENTS) {
      getComponents();
    } else if (view === RELATIONSHIPS) {
      getRelationships(page, sortOrder);
    } else if (view === REGISTRANTS) {
      getRegistrants();
    }
  }, [view, page, searchText, rowsPerPage]);

  return (
    <div data-test="workloads">
      <div
        className={`${StyleClass.meshModelToolbar} ${animate ? StyleClass.toolWrapperAnimate : ''}`}
      >
        <DisableButton
          disabled
          variant="contained"
          style={{
            visibility: `${animate ? 'visible' : 'hidden'}`,
          }}
          size="large"
          startIcon={<UploadIcon />}
        >
          Import
        </DisableButton>
        <DisableButton
          disabled
          variant="contained"
          size="large"
          style={{
            visibility: `${animate ? 'visible' : 'hidden'}`,
          }}
          startIcon={<DoNotDisturbOnIcon />}
        >
          Ignore
        </DisableButton>
      </div>
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
            className={`${StyleClass.cardStyle} ${animate ? StyleClass.cardStyleAnimate : ''} ${
              view === MODELS && animate ? StyleClass.activeTab : ''
            }`}
            onClick={() => {
              setView(() => {
                handleChangeSelectedTab(MODELS);
                return MODELS;
              });
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
                setConvert(true);
              }
            }}
          >
            <span
              style={{
                fontWeight: `${animate ? 'normal' : 'bold'}`,
                fontSize: `${animate ? '1rem' : '3rem'}`,
                marginLeft: `${animate && '4px'}`,
              }}
            >
              {animate ? `(${modelsCount})` : `${modelsCount}`}
            </span>
            Models
          </Paper>
          <Paper
            elevation={3}
            className={`${StyleClass.cardStyle} ${animate ? StyleClass.cardStyleAnimate : ''} ${
              view === COMPONENTS && animate ? StyleClass.activeTab : ''
            }`}
            onClick={() => {
              setView(() => {
                handleChangeSelectedTab(COMPONENTS);
                return COMPONENTS;
              });
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
                setConvert(true);
              }
            }}
          >
            <span
              style={{
                fontWeight: `${animate ? 'normal' : 'bold'}`,
                fontSize: `${animate ? '1rem' : '3rem'}`,
                marginLeft: `${animate && '4px'}`,
              }}
            >
              {animate ? `(${componentsCount})` : `${componentsCount}`}
            </span>
            Components
          </Paper>
          <Paper
            elevation={3}
            className={`${StyleClass.cardStyle} ${animate ? StyleClass.cardStyleAnimate : ''} ${
              view === RELATIONSHIPS && animate ? StyleClass.activeTab : ''
            }`}
            onClick={() => {
              setView(() => {
                handleChangeSelectedTab(RELATIONSHIPS);
                return RELATIONSHIPS;
              });
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
                setConvert(true);
              }
            }}
          >
            <span
              style={{
                fontWeight: `${animate ? 'normal' : 'bold'}`,
                fontSize: `${animate ? '1rem' : '3rem'}`,
                marginLeft: `${animate && '4px'}`,
              }}
            >
              {animate ? `(${relationshipsCount})` : `${relationshipsCount}`}
            </span>
            Relationships
          </Paper>
          <Paper
            elevation={3}
            className={`${StyleClass.cardStyle} ${animate ? StyleClass.cardStyleAnimate : ''} ${
              view === REGISTRANTS && animate ? StyleClass.activeTab : ''
            }`}
            onClick={() => {
              setView(() => {
                handleChangeSelectedTab(REGISTRANTS);
                return REGISTRANTS;
              });
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
                setConvert(true);
              }
            }}
          >
            <span
              style={{
                fontWeight: `${animate ? 'normal' : 'bold'}`,
                fontSize: `${animate ? '1rem' : '3rem'}`,
                marginLeft: `${animate && '4px'}`,
              }}
            >
              {animate ? `(${registrantCount})` : `${registrantCount}`}
            </span>
            Registrants
          </Paper>
        </div>
        {convert && (
          <div
            className={`${StyleClass.treeWrapper} ${convert ? StyleClass.treeWrapperAnimate : ''}`}
          >
            <div
              className={StyleClass.treeContainer}
              style={{
                alignItems: filteredData.length === 0 ? 'center' : '',
              }}
            >
              {filteredData.length === 0 ? (
                <CircularProgress sx={{ color: Colors.keppelGreen }} />
              ) : (
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
                  searchText={searchText}
                />
              )}
            </div>
            <MeshModelDetails view={view} show={show} rela={rela} regi={regi} comp={comp} />
          </div>
        )}
      </div>
    </div>
  );
};

const MeshModelComponent = (props) => {
  return (
    <NoSsr>
      <ErrorBoundary
        FallbackComponent={() => null}
        onError={(e) => console.error('Error in NotificationCenter', e)}
      >
        <Provider store={store}>
          <MeshModelComponent_ {...props} />
        </Provider>
      </ErrorBoundary>
    </NoSsr>
  );
};

export default withStyles(meshmodelStyles)(withSnackbar(MeshModelComponent));
