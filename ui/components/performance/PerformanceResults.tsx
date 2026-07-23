import React, { useEffect, useState, useRef } from 'react';
import MUIDataTable from '@sistent/mui-datatables';
import Moment from 'react-moment';
import CustomToolbarSelect from '../CustomToolbarSelect';
import MesheryChart from '../MesheryChart';
import GenericModal from '../shared/Modal/GenericModal';
import { Info as InfoIcon, Reply as ReplyIcon } from '@/assets/icons';
import NodeDetails from './NodeDetails';
import FacebookIcon from './assets/facebookIcon';
import LinkedinIcon from './assets/linkedinIcon';
import TwitterIcon from './assets/twitterIcon';
import { iconMedium, iconLarge } from '../../css/icons.styles';
import { TwitterShareButton, LinkedinShareButton, FacebookShareButton } from 'react-share';
import { useNotification } from '../../utils/hooks/useNotification';
import { EVENT_TYPES } from '../../lib/event-types';
import {
  Tab,
  Tabs,
  IconButton,
  Paper,
  Popper,
  ClickAwayListener,
  useTheme,
  NoSsr,
  Fade,
  BarChartIcon,
} from '@sistent/sistent';

import { DefaultTableCell, SortableTableCell } from '../connections/common';
import { useDispatch, useSelector } from 'react-redux';
import { isLocalProvider } from '@/utils/provider';
import { updateProgressAction } from '@/store/slices/mesheryUi';
import { updateResultsSelection } from '@/store/slices/prefTest';
import { useGetPerformanceProfileResultsQuery } from '@meshery/schemas/mesheryApi';

const COL_MAPPING = {
  QPS: 3,
  P99: 6,
};

function generateResultsForDisplay(results) {
  if (Array.isArray(results)) {
    return results.map((record) => {
      const runnerResults = record.runnerResults || {};
      const data = {
        name: record.name,
        mesh: record.mesh,
        testStartTime: runnerResults.StartTime || record.testStartTime,
        qps: runnerResults.ActualQPS?.toFixed(1) || 'unavailable',
        duration: runnerResults.ActualDuration
          ? (runnerResults.ActualDuration / 1000000000).toFixed(1)
          : 'unavailable',
        threads: runnerResults.NumThreads,
      };

      if (runnerResults?.DurationHistogram?.Percentiles) {
        runnerResults.DurationHistogram.Percentiles.forEach(({ Percentile, Value }) => {
          data[`p${Percentile}`.replace('.', '_')] = Value.toFixed(3);
        });
      } else {
        data.p50 = 0;
        data.p75 = 0;
        data.p90 = 0;
        data.p99 = 0;
        data.p99_9 = 0;
      }

      return data;
    });
  }

  return [];
}

function generateColumnsForDisplay(
  sortOrder,
  setSelectedProfileIdxForChart,
  setSelectedProfileIdxForNodeDetails,
  theme,
  handleSocialExpandClick,
  handleClickAway,
  socialExpand,
  anchorEl,
  socialMessage,
) {
  const columns = [
    {
      name: 'name',
      label: 'Name',
      options: {
        filter: false,
        sort: true,
        searchable: true,
        customHeadRender: function CustomHead({ index, ...column }, sortColumn, columnMeta) {
          return (
            <SortableTableCell
              index={index}
              columnData={column}
              columnMeta={columnMeta}
              onSort={() => sortColumn(index)}
            />
          );
        },
      },
    },
    {
      name: 'mesh',
      label: 'Mesh',
      options: {
        filter: false,
        sort: true,
        searchable: true,
        customHeadRender: function CustomHead({ index, ...column }, sortColumn, columnMeta) {
          return (
            <SortableTableCell
              index={index}
              columnData={column}
              columnMeta={columnMeta}
              onSort={() => sortColumn(index)}
            />
          );
        },
      },
    },
    {
      name: 'testStartTime',
      label: 'Start Time',
      options: {
        filter: false,
        sort: true,
        searchable: true,
        customHeadRender: function CustomHead({ index, ...column }, sortColumn, columnMeta) {
          return (
            <SortableTableCell
              index={index}
              columnData={column}
              columnMeta={columnMeta}
              onSort={() => sortColumn(index)}
            />
          );
        },
        customBodyRender: function CustomBody(value) {
          return <Moment format="LLLL">{value}</Moment>;
        },
      },
    },
    {
      name: 'qps',
      label: 'QPS',
      options: {
        filter: false,
        sort: false,
        searchable: false,
        customHeadRender: function CustomHead({ ...column }) {
          return <DefaultTableCell columnData={column} />;
        },
      },
    },
    {
      name: 'duration',
      label: 'Duration',
      options: {
        filter: false,
        sort: false,
        searchable: false,
        customHeadRender: function CustomHead({ ...column }) {
          return <DefaultTableCell columnData={column} />;
        },
      },
    },

    {
      name: 'p50',
      label: 'P50',
      options: {
        filter: false,
        sort: false,
        searchable: false,
        customHeadRender: function CustomHead({ ...column }) {
          return <DefaultTableCell columnData={column} />;
        },
      },
    },

    {
      name: 'p99_9',
      label: 'P99.9',
      options: {
        filter: false,
        sort: false,
        searchable: false,
        customHeadRender: function CustomHead({ ...column }) {
          return <DefaultTableCell columnData={column} />;
        },
      },
    },
    {
      name: 'Chart',
      options: {
        filter: false,
        sort: false,
        searchable: false,
        customHeadRender: function CustomHead({ ...column }) {
          return <DefaultTableCell columnData={column} />;
        },
        customBodyRender: function CustomBody(value, tableMeta) {
          return (
            <IconButton
              aria-label="more"
              data-testid="open-performance-result-bar-chart"
              color="inherit"
              onClick={() => setSelectedProfileIdxForChart(tableMeta.rowIndex)}
            >
              <BarChartIcon style={iconMedium} /> {/* can change it to large */}
            </IconButton>
          );
        },
      },
    },
    {
      name: 'Node Details',
      options: {
        filter: false,
        sort: false,
        searchable: false,
        customHeadRender: function CustomHead({ ...column }) {
          return <DefaultTableCell columnData={column} />;
        },
        customBodyRender: function CustomBody(value, tableMeta) {
          return (
            <IconButton
              aria-label="more"
              color="inherit"
              onClick={() => setSelectedProfileIdxForNodeDetails(tableMeta.rowIndex)}
            >
              <InfoIcon style={iconMedium} /> {/* can change it to large */}
            </IconButton>
          );
        },
      },
    },
    {
      name: 'Share Results',
      options: {
        filter: false,
        sort: false,
        searchable: false,
        customHeadRender: function CustomHead({ ...column }) {
          return <DefaultTableCell columnData={column} />;
        },
        customBodyRender: function CustomBody(_, tableMeta) {
          return (
            <>
              <IconButton aria-label="Share" onClick={(e) => handleSocialExpandClick(e, tableMeta)}>
                <ReplyIcon
                  style={{
                    transform: 'scaleX(-1)',
                    color: theme.palette.icon.default,
                    ...iconLarge,
                  }}
                />
              </IconButton>
              <Popper
                open={socialExpand[tableMeta.rowIndex]}
                anchorEl={anchorEl[tableMeta.rowIndex]}
                transition
              >
                {({ TransitionProps }) => (
                  <ClickAwayListener onClickAway={() => handleClickAway(tableMeta.rowIndex)}>
                    <Fade {...TransitionProps} timeout={350}>
                      <Paper sx={{ padding: theme.spacing(1) }}>
                        <TwitterShareButton
                          style={{ margin: theme.spacing(0.4) }}
                          url={'https://meshery.io'}
                          title={socialMessage}
                          hashtags={['opensource']}
                        >
                          {/* <img src={`/static/img/twitter.svg`} /> */}
                          <TwitterIcon />
                        </TwitterShareButton>
                        <LinkedinShareButton
                          style={{ margin: theme.spacing(0.4) }}
                          url={'https://meshery.io'}
                          summary={socialMessage}
                        >
                          {/* <img src={`/static/img/linkedin.svg`} /> */}
                          <LinkedinIcon />
                        </LinkedinShareButton>
                        <FacebookShareButton
                          style={{ margin: theme.spacing(0.4) }}
                          url={'https://meshery.io'}
                          quote={socialMessage}
                          hashtag={'#opensource'}
                        >
                          {/* <img src={`/static/img/facebook.svg`} /> */}
                          <FacebookIcon />
                        </FacebookShareButton>
                      </Paper>
                    </Fade>
                  </ClickAwayListener>
                )}
              </Popper>
            </>
          );
        },
      },
    },
  ];

  return columns.map((column) => {
    if (column.name === sortOrder.split(' ')[0]) {
      column.options.sortDirection = sortOrder.split(' ')[1];
    }

    return column;
  });
}

function getSocialMessageForPerformanceTest(rps, percentile) {
  return `I achieved ${rps.trim()} RPS running my service at a P99.9 of ${percentile} ms using @mesheryio with @smp_spec! Find out how fast your service is with`;
}

function generateSelectedRows(results_selection, page, pageSize) {
  const rowsSelected = [];
  Object.keys(results_selection).forEach((pg) => {
    if (parseInt(pg) !== page) {
      Object.keys(results_selection[parseInt(pg)]).forEach((ind) => {
        const val = (parseInt(pg) + 1) * pageSize + parseInt(ind) + 1;
        rowsSelected.push(val);
      });
    } else {
      Object.keys(results_selection[page]).forEach((ind) => {
        const val = parseInt(ind);
        rowsSelected.push(val);
      });
    }
  });

  return rowsSelected;
}

function ResultChart({ result, handleTabChange, tabValue }) {
  if (!result) return <div />;

  // const getMuiTheme = () => createTheme({
  //   overrides : {
  //     MuiTab : {
  //       textColorInherit : {
  //         textTransform : "none",
  //         backgroundColor : "#eaeff1"
  //       }
  //     }
  //   }
  // })
  const row = result.runnerResults;
  if (!row) return <div />;

  return (
    <Paper
      style={{
        width: '100%',
        maxWidth: '90vw',
        padding: '0.5rem',
      }}
    >
      <Tabs value={tabValue} onChange={handleTabChange}>
        <Tab label="Performance Chart" />
        <Tab label="Node Details" />
      </Tabs>

      {tabValue == 0 ? (
        <div>
          <div>
            <MesheryChart
              rawdata={[result && result.runnerResults ? result : {}]}
              data={[result && result.runnerResults ? result.runnerResults : {}]}
            />
          </div>
        </div>
      ) : tabValue == 1 ? (
        <div>
          <NodeDetails result={row} />
        </div>
      ) : (
        <div />
      )}
    </Paper>
  );
}

function ResultNodeDetails({ result, handleTabChange, tabValue }) {
  if (!result) return <div />;
  const chartData = result.runnerResults;

  const row = result.runnerResults;
  if (!row) return <div />;

  return (
    <Paper
      style={{
        width: '100%',
        maxWidth: '90vw',
        padding: '0.5rem',
      }}
    >
      <Tabs
        value={tabValue}
        onChange={handleTabChange}
        TabIndicatorProps={{
          style: {
            backgroundColor: '#00B39F',
          },
        }}
      >
        <Tab label="Performance Chart" />
        <Tab label="Node Details" />
      </Tabs>
      {tabValue == 1 ? (
        <div>
          <NodeDetails result={chartData} />
        </div>
      ) : tabValue == 0 ? (
        <div>
          <div>
            <MesheryChart
              rawdata={[result && result.runnerResults ? result : {}]}
              data={[result && result.runnerResults ? result.runnerResults : {}]}
            />
          </div>
        </div>
      ) : (
        <div />
      )}
    </Paper>
  );
}

/**
 *
 * @param {{
 *  updateProgress?: any,
 *  results_selection?: any,
 *  user?: any
 *  updateResultsSelection?: any,
 *  endpoint: string,
 *  CustomHeader?: JSX.Element
 *  elevation?: Number
 * }} props
 */
function MesheryResults({ endpoint, CustomHeader = <div />, elevation = 4 }) {
  const [page, setPage] = useState(0);
  const [search, setSearch] = useState('');
  const [sortOrder, setSortOrder] = useState('');
  const [count, setCount] = useState(0);
  const [pageSize, setPageSize] = useState(10);
  const [results, setResults] = useState([]);
  const [selectedRowChart, setSelectedRowChart] = useState();
  const [selectedRowNodeDetails, setSelectedRowNodeDetails] = useState();
  const [tabValue, setTabValue] = useState(0);
  const [socialExpand, setSocialExpand] = useState([false]);
  const [anchorEl, setAnchorEl] = useState([]);
  const [socialMessage, setSocialMessage] = useState();
  const theme = useTheme();
  const { providerCapabilities } = useSelector((state) => state.ui);
  const dispatch = useDispatch();
  const { results_selection } = useSelector((state) => state.prefTest);
  const searchTimeout = useRef();
  const { notify } = useNotification();
  const endpointParts = endpoint.split('/');
  const profileID = endpointParts[endpointParts.length - 2];
  const {
    data: performanceResultsData,
    isFetching,
    isError,
    error,
  } = useGetPerformanceProfileResultsQuery(
    {
      performanceProfileId: profileID,
      page: `${page}`,
      pagesize: `${pageSize}`,
      search: search || '',
      order: sortOrder || '',
    },
    {
      skip: !profileID,
    },
  );

  useEffect(() => {
    dispatch(updateProgressAction({ showProgress: isFetching }));
  }, [dispatch, isFetching]);

  useEffect(() => {
    if (!performanceResultsData) return;

    setCount(performanceResultsData.totalCount);
    setPageSize(performanceResultsData.pageSize);
    setResults(performanceResultsData.results || []);
  }, [performanceResultsData]);

  useEffect(() => {
    if (isError) handleError(error);
  }, [isError, error]);

  const handleSocialExpandClick = (e, tableMeta) => {
    let socialExpandUpdate = [...socialExpand];
    socialExpandUpdate[tableMeta.rowIndex] = !socialExpand[tableMeta.rowIndex];
    setSocialExpand(socialExpandUpdate);

    let anchorElUpdate = [...anchorEl];
    anchorElUpdate[tableMeta.rowIndex] = e.currentTarget;
    setAnchorEl(anchorElUpdate);
    setSocialMessage(
      getSocialMessageForPerformanceTest(
        tableMeta.rowData[COL_MAPPING.QPS],
        tableMeta.rowData[COL_MAPPING.P99],
      ),
    );
  };

  const handleClickAway = (index) => {
    let socialExpandUpdate = [...socialExpand];
    socialExpandUpdate[index] = !socialExpand[index];
    setSocialExpand(socialExpandUpdate);
  };

  function handleError(error) {
    dispatch(updateProgressAction({ showProgress: false }));
    notify({
      message: `There was an error fetching results: ${error}`,
      event_type: EVENT_TYPES.ERROR,
      details: error.toString(),
    });
  }

  const columns = generateColumnsForDisplay(
    sortOrder,
    (idx) => {
      setSelectedRowChart(results[idx]);
      setTabValue(0);
    },
    (idx) => {
      setSelectedRowNodeDetails(results[idx]);
      setTabValue(1);
    },
    theme,
    handleSocialExpandClick,
    handleClickAway,
    socialExpand,
    anchorEl,
    socialMessage,
  );

  const options = {
    elevation: elevation,
    filter: false,
    sort: !isLocalProvider(providerCapabilities),
    search: !isLocalProvider(providerCapabilities),
    filterType: 'textField',
    responsive: 'standard',
    resizableColumns: true,
    selectableRows: 'multiple',
    serverSide: true,
    count,
    rowsPerPage: pageSize,
    rowsPerPageOptions: [10, 20, 25],
    fixedHeader: true,
    page,
    rowsSelected: generateSelectedRows(results_selection, page, pageSize),
    print: false,
    download: false,
    onRowsSelect: (_, allRowsSelected) => {
      // const rs = self.props.results_selection;
      const res = {};
      allRowsSelected.forEach(({ dataIndex }) => {
        if (dataIndex < pageSize) {
          if (res[dataIndex]) delete res[dataIndex];
          else res[dataIndex] = results[dataIndex];
        }
      });

      dispatch(updateResultsSelection({ page, results: res }));
    },

    onTableChange: (action, tableState) => {
      const sortInfo = tableState.announceText ? tableState.announceText.split(' : ') : [];
      let order = '';
      if (tableState.activeColumn) {
        order = `${columns[tableState.activeColumn].name} desc`;
      }

      switch (action) {
        case 'changePage':
          setPage(tableState.page);
          break;
        case 'changeRowsPerPage':
          setPageSize(tableState.rowsPerPage);
          break;
        case 'search':
          if (searchTimeout.current) {
            clearTimeout(searchTimeout.current);
          }
          // @ts-ignore
          searchTimeout.current = setTimeout(() => {
            if (search !== tableState.searchText) {
              setSearch(tableState.searchText || '');
            }
          }, 500);
          break;
        case 'sort':
          if (sortInfo.length === 2) {
            if (sortInfo[1] === 'ascending') order = `${columns[tableState.activeColumn].name} asc`;
            else order = `${columns[tableState.activeColumn].name} desc`;
          }

          if (order !== sortOrder) setSortOrder(order);
          break;
      }
    },
    customToolbarSelect: function CustomToolbarSelectComponent(
      selectedRows,
      displayData,
      setSelectedRows,
    ) {
      return (
        <CustomToolbarSelect
          selectedRows={selectedRows}
          displayData={displayData}
          setSelectedRows={setSelectedRows}
          results={results}
        />
      );
    },
  };

  function handleTabChange(event, newValue) {
    setTabValue(newValue);
  }

  return (
    <NoSsr>
      <MUIDataTable
        title={CustomHeader}
        data={generateResultsForDisplay(results)}
        columns={columns}
        // @ts-ignore
        options={options}
      />

      <GenericModal
        open={!!selectedRowChart}
        // @ts-ignore
        Content={
          <ResultChart
            result={selectedRowChart}
            handleTabChange={handleTabChange}
            tabValue={tabValue}
          />
        }
        handleClose={() => setSelectedRowChart(undefined)}
      />

      <GenericModal
        open={!!selectedRowNodeDetails}
        // @ts-ignore
        Content={
          <ResultNodeDetails
            result={selectedRowNodeDetails}
            handleTabChange={handleTabChange}
            tabValue={tabValue}
          />
        }
        handleClose={() => setSelectedRowNodeDetails(undefined)}
      />
    </NoSsr>
  );
}
export default MesheryResults;
