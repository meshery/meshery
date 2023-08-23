import React, { useEffect, useState } from 'react';
import IconButton from '@material-ui/core/IconButton';
import Tooltip from '@material-ui/core/Tooltip';
import IndeterminateCheckBoxIcon from '@material-ui/icons/IndeterminateCheckBox';
import CompareArrowsIcon from '@material-ui/icons/CompareArrows';
import { withStyles } from '@material-ui/core/styles';
import { NoSsr } from '@material-ui/core';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import GetAppIcon from '@material-ui/icons/GetApp';
import MesheryChartDialog from './MesheryChartDialog';
import MesheryChart from './MesheryChart';
import { clearResultsSelection } from '../lib/store';

const defaultToolbarSelectStyles = {
  iconButton : {
    marginRight : '24px',
    top : '50%',
    display : 'inline-block',
    position : 'relative',
  },
  icon : {
    color : '#000',
  },
  inverseIcon : {
    transform : 'rotate(90deg)',
  },
};

function CustomToolbarSelect({ classes, results_selection, setSelectedRows, clearResultsSelection }) {
  const [dialogOpen, setDialogOpen] = useState(false);
  const [data, setData] = useState([]);
  const [, setChartCompare] = useState([]);
  const fullData = [];

  const handleDialogClose = () => {
    setDialogOpen(false);
  };

  // Never been used
  // const handleDialogOpen = () => {
  //   setDialogOpen(true);
  // };

  const handleClickDeselectAll = () => {
    setSelectedRows([]);
    clearResultsSelection();
  };

  const handleCompareSelected = () => {
    const rs = results_selection;

    // Reset the data and chartCompare states.
    setData([]);
    setChartCompare([]);

    Object.keys(rs).forEach((k1) => {
      Object.keys(rs[k1]).forEach((k2) => {
        if (typeof rs[k1][k2] !== 'undefined') {
          // Directly update the data state.
          setData(prevData => [...prevData, rs[k1][k2].runner_results]);

          const row = rs[k1][k2].runner_results;
          const startTime = new Date(row.StartTime);
          const endTime = new Date(startTime.getTime() + row.ActualDuration / 1000000);
          const boardConfig = rs[k1][k2].server_board_config;
          const serverMetrics = rs[k1][k2].server_metrics;

          // Directly update the chartCompare state.
          setChartCompare(prevChartCompare => [...prevChartCompare, {
            label : row.Labels,
            startTime,
            endTime,
            boardConfig,
            serverMetrics,
          }]);
        }
      });
    });
    setDialogOpen(true);
  };

  useEffect(() => {
    const rs = results_selection;
    Object.keys(rs).forEach((k1) => {
      Object.keys(rs[k1]).forEach((k2) => {
        if (typeof rs[k1][k2] !== 'undefined') {
          fullData.push(rs[k1][k2]);
        }
      });
    });
  }, [results_selection]);

  return (
    <NoSsr>
      <div className="custom-toolbar-select">
        <Tooltip title="Deselect ALL">
          <IconButton className={classes.iconButton} onClick={handleClickDeselectAll}>
            <IndeterminateCheckBoxIcon className={classes.icon} />
          </IconButton>
        </Tooltip>
        {fullData.length === 1 && (
          <Tooltip title="Download">
            <IconButton
              className={classes.iconButton}
              key="download"
              aria-label="download"
              color="inherit"
              href={`/api/perf/profile/result/${encodeURIComponent(fullData[0].meshery_id)}`}
              download={`${fullData[0].name}_test_result.json`}
            >
              <GetAppIcon className={classes.icon} />
            </IconButton>
          </Tooltip>
        )}
        <Tooltip title="Compare selected">
          <IconButton className={classes.iconButton} onClick={handleCompareSelected}>
            <CompareArrowsIcon className={classes.icon} />
          </IconButton>
        </Tooltip>
      </div>
      <MesheryChartDialog
        handleClose={handleDialogClose}
        open={dialogOpen}
        content={<div><MesheryChart data={data} /></div>}
      />
    </NoSsr>
  );
}

const mapDispatchToProps = (dispatch) => ({
  clearResultsSelection : bindActionCreators(clearResultsSelection, dispatch),
});

const mapStateToProps = (state) => {
  const results_selection = state.get('results_selection').toObject();
  return { results_selection };
};

export default withStyles(defaultToolbarSelectStyles, { name : 'CustomToolbarSelect' })(
  connect(
    mapStateToProps,
    mapDispatchToProps,
  )(CustomToolbarSelect),
);
