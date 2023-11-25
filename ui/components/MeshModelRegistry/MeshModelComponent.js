import { withStyles } from '@material-ui/core';
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
import MesheryTreeView from './MesheryTreeView';
import MeshModelDetails from './MeshModelDetails';
import { toLower } from 'lodash';
import { DisableButton } from './MeshModel.style';
import CircularProgress from '@mui/material/CircularProgress';
import { Colors } from '../../themes/app';

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

const MeshModelComponent = ({
  modelsCount,
  componentsCount,
  relationshipsCount,
  registrantCount,
}) => {
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

  const getModels = async (page) => {
    try {
      const { models } = await getMeshModels(page?.Models + 1, rowsPerPage);

      if (!isRequestCancelled && models) {
        const updatedModels = [];

        for (const model of models) {
          const { components } = await getComponentFromModelApi(model.name);
          const { relationships } = await getRelationshipFromModelApi(model.name);
          model.components = components;
          model.relationships = relationships;
          updatedModels.push(model);
        }

        setResourcesDetail((prev) =>
          [...prev, ...updatedModels].sort((a, b) => a.displayName.localeCompare(b.displayName)),
        );
      }
    } catch (error) {
      console.error('Failed to fetch models:', error);
    }
  };

  const getComponents = async (page, sortOrder) => {
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

      if (!isRequestCancelled) {
        const updatedModels = [];

        for (const model of models || []) {
          const { components } = await getComponentFromModelApi(model.name);
          const { relationships } = await getRelationshipFromModelApi(model.name);
          model.components = components;
          model.relationships = relationships;
          updatedModels.push(model);
        }

        setCount(total_count);
        setResourcesDetail(updatedModels);
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
      const { total_count, registrants } = await getMeshModelRegistrants(
        page?.Registrants + 1,
        rowsPerPage,
      );

      if (!isRequestCancelled && registrants) {
        let tempRegistrants = [];

        for (const registrant of registrants) {
          let hostname = toLower(registrant?.hostname);
          const { models } = await getMeshModelsByRegistrants(
            page?.Models + 1,
            rowsPerPage,
            hostname,
          ); // page+1 due to server side indexing starting from 1

          if (models) {
            const updatedModels = [];

            for (const model of models) {
              const { components } = await getComponentFromModelApi(model.name);
              const { relationships } = await getRelationshipFromModelApi(model.name);
              model.components = components;
              model.relationships = relationships;
              updatedModels.push(model);
            }

            registrant.models = updatedModels;
            tempRegistrants.push(registrant);
          } else {
            tempRegistrants.push(registrant);
          }
        }

        setCount(total_count);

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
    } catch (error) {
      console.error('Failed to fetch registrants:', error);
    }
  };

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

  useEffect(
    () => {
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
    setRequestCancelled(false);

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
    } else if (view === MODELS && searchText === '') {
      getModels(page);
    } else if (view === COMPONENTS && searchText === '') {
      getComponents(page, sortOrder);
    }

    return () => {
      setRequestCancelled(true);
    };
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

export default withStyles(meshmodelStyles)(withSnackbar(MeshModelComponent));
