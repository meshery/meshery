import { withStyles } from '@material-ui/core';
import { withSnackbar } from 'notistack';
import React, { useState, useEffect, useCallback } from 'react';
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
// import { SORT } from '../../constants/endpoints';
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
import { removeDuplicateVersions } from './helper';

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
      color: theme.palette.secondary.focused,
    },
  },
  tabs: {
    '& .MuiTabs-indicator': {
      backgroundColor: theme.palette.secondary.focused,
    },
  },
  dashboardSection: {
    padding: theme.spacing(2),
    borderRadius: 4,
    height: '100%',
    overflowY: 'scroll',
  },
  duplicatesModelStyle: {
    backgroundColor: theme.palette.secondary.focused,
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
  // const [sortOrder] = useState({
  //   sort: SORT.ASCENDING,
  //   order: '',
  // });
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
  const [checked, setChecked] = useState(false);
  // const [loading, setLoading] = useState(false);

  /**
   * RTK Lazy Queries
   */
  const [getMeshModelsData] = useLazyGetMeshModelsQuery();
  const [getComponentsData] = useLazyGetComponentsQuery();
  const [getRelationshipsData] = useLazyGetRelationshipsQuery();
  const [getRegistrantsData] = useLazyGetRegistrantsQuery();

  const fetchData = useCallback(async () => {
    try {
      let response;
      switch (view) {
        case MODELS:
          response = await getMeshModelsData(
            {
              params: {
                page: searchText ? 1 : page.Models + 1,
                pagesize: searchText || checked ? 'all' : rowsPerPage,
                components: true,
                relationships: true,
                paginated: true,
                search: searchText || '',
              },
            },
            true,
          );
          break;
        case COMPONENTS:
          response = await getComponentsData(
            {
              params: {
                page: searchText ? 1 : page.Components + 1,
                pagesize: searchText ? 'all' : rowsPerPage,
                search: searchText || '',
                trim: false,
              },
            },
            true,
          );
          setCount(response.total_count);
          break;
        case RELATIONSHIPS:
          response = await getRelationshipsData(
            {
              params: {
                page: searchText ? 1 : page.Relationships + 1,
                pagesize: searchText ? 'all' : rowsPerPage,
                paginated: true,
                search: searchText || '',
              },
            },
            true,
          );
          setCount(response?.total_count);
          break;
        case REGISTRANTS:
          response = await getRegistrants();
          break;
        default:
          break;
      }

      if (response.data) {
        const newData =
          searchText || checked
            ? [...response.data[view.toLowerCase()]]
            : [...resourcesDetail, ...response.data[view.toLowerCase()]];
        setResourcesDetail(newData);

        if (rowsPerPage !== 14) {
          setRowsPerPage(14);
        }
      }
    } catch (error) {
      console.error(`Failed to fetch ${view.toLowerCase()}:`, error);
    }
  }, [
    getMeshModelsData,
    getComponentsData,
    getRelationshipsData,
    getRegistrantsData,
    view,
    searchText,
    page,
    rowsPerPage,
    resourcesDetail,
    checked,
  ]);

  const getRegistrants = async () => {
    let registrantResponse;
    let response;
    registrantResponse = await getRegistrantsData(
      {
        params: {
          page: searchText ? 1 : page.Registrants + 1,
          pagesize: searchText ? 'all' : rowsPerPage,
          search: searchText || '',
        },
      },
      true,
    );
    if (registrantResponse.data && registrantResponse.data.registrants) {
      const registrants = registrantResponse.data.registrants;
      const tempResourcesDetail = [];

      for (let registrant of registrants) {
        let hostname = toLower(registrant?.hostname);
        const { data: modelRes } = await getMeshModelsData(
          {
            params: {
              page: page?.Models + 1,
              pagesize: 'all',
              registrant: hostname,
              components: true,
              relationships: true,
            },
          },
          true,
        );
        const updatedRegistrant = {
          ...registrant,
          models: removeDuplicateVersions(modelRes.models) || [],
        };
        tempResourcesDetail.push(updatedRegistrant);
      }
      response = {
        data: {
          registrants: tempResourcesDetail,
        },
      };
    }
    setCount(response.total_count);
    setRowsPerPage(14);
    return response;
  };

  const handleTabClick = (selectedView) => {
    if (view !== selectedView) {
      setSearchText(null);
      setResourcesDetail([]);
    }
    setView(() => {
      handleChangeSelectedTab(selectedView);
      return selectedView;
    });
    setPage({
      Models: 0,
      Components: 0,
      Relationships: 0,
      Registrants: 0,
    });
    if (!animate) {
      setAnimate(true);
      setConvert(true);
    }
  };

  const modifyData = () => {
    if (view === MODELS) {
      return removeDuplicateVersions(
        checked ? resourcesDetail.filter((model) => model.duplicates > 0) : resourcesDetail,
      );
    } else {
      return resourcesDetail;
    }
  };
  useEffect(() => {
    if (selectedTab && selectedTab !== OVERVIEW) {
      setAnimate(true);
      setConvert(true);
      setView(selectedTab);
    }
  }, [selectedTab]);

  useEffect(() => {
    fetchData();
  }, [view, page, searchText, rowsPerPage, checked]);

  return (
    <div data-test="workloads">
      <TabBar animate={animate} />
      <div
        className={`${StyleClass.mainContainer} ${animate ? StyleClass.mainContainerAnimate : ''}`}
      >
        <div
          className={`${StyleClass.innerContainer} ${
            animate ? StyleClass.innerContainerAnimate : ''
          }`}
        >
          <TabCard
            label="Models"
            count={modelsCount}
            active={view === MODELS && animate}
            animate={animate}
            onClick={() => handleTabClick(MODELS)}
          />
          <TabCard
            label="Components"
            count={componentsCount}
            active={view === COMPONENTS && animate}
            animate={animate}
            onClick={() => handleTabClick(COMPONENTS)}
          />
          <TabCard
            label="Relationships"
            count={relationshipsCount}
            active={view === RELATIONSHIPS && animate}
            animate={animate}
            onClick={() => handleTabClick(RELATIONSHIPS)}
          />
          <TabCard
            label="Registrants"
            count={registrantCount}
            active={view === REGISTRANTS && animate}
            animate={animate}
            onClick={() => handleTabClick(REGISTRANTS)}
          />
        </div>
        {convert && (
          <div
            className={`${StyleClass.treeWrapper} ${convert ? StyleClass.treeWrapperAnimate : ''}`}
          >
            <div
              className={StyleClass.treeContainer}
              style={{
                alignItems: resourcesDetail.length === 0 ? 'center' : '',
              }}
            >
              {resourcesDetail.length === 0 ? (
                <CircularProgress sx={{ color: Colors.keppelGreen }} />
              ) : (
                <MesheryTreeView
                  data={modifyData()}
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

const TabBar = ({ animate }) => {
  const StyleClass = useStyles();
  return (
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
  );
};

const TabCard = ({ label, count, active, onClick, animate }) => {
  const StyleClass = useStyles();
  return (
    <Paper
      elevation={3}
      className={`${StyleClass.cardStyle} ${animate ? StyleClass.cardStyleAnimate : ''} ${
        active ? StyleClass.activeTab : ''
      }`}
      onClick={onClick}
    >
      <span
        style={{
          fontWeight: `${animate ? 'normal' : 'bold'}`,
          fontSize: `${animate ? '1rem' : '3rem'}`,
          marginLeft: `${animate && '4px'}`,
        }}
      >
        {animate ? `(${count})` : `${count}`}
      </span>
      {label}
    </Paper>
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
