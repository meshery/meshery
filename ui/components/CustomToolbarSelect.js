import React, { useState } from 'react';
import { IconButton, Tooltip, styled } from '@layer5/sistent';
import IndeterminateCheckBoxIcon from '@mui/icons-material/IndeterminateCheckBox';
import CompareArrowsIcon from '@mui/icons-material/CompareArrows';
import { NoSsr } from '@mui/material';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import GetAppIcon from '@mui/icons-material/GetApp';
import MesheryChartDialog from './MesheryChartDialog';
import MesheryChart from './MesheryChart';
import { clearResultsSelection } from '../lib/store';

const StyledIconButton = styled(IconButton)(({ theme }) => ({
  marginRight: theme.spacing(3),
  top: '50%',
  display: 'inline-block',
  position: 'relative',
}));

const StyledIcon = styled('span')(({ theme }) => ({
  color:
    theme.palette.mode === 'dark'
      ? theme.palette.background.brand.pressed
      : theme.palette.text.default,
}));

function CustomToolbarSelect({ results_selection, setSelectedRows, clearResultsSelection }) {
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
          setData((prevData) => [...prevData, rs[k1][k2].runner_results]);

          const row = rs[k1][k2].runner_results;
          const startTime = new Date(row.StartTime);
          const endTime = new Date(startTime.getTime() + row.ActualDuration / 1000000);
          const boardConfig = rs[k1][k2].server_board_config;
          const serverMetrics = rs[k1][k2].server_metrics;

          // Directly update the chartCompare state.
          setChartCompare((prevChartCompare) => [
            ...prevChartCompare,
            {
              label: row.Labels,
              startTime,
              endTime,
              boardConfig,
              serverMetrics,
            },
          ]);
        }
      });
    });
    setDialogOpen(true);
  };

  const rs = results_selection;
  Object.keys(rs).forEach((k1) => {
    Object.keys(rs[k1]).forEach((k2) => {
      if (typeof rs[k1][k2] !== 'undefined') {
        fullData.push(rs[k1][k2]);
      }
    });
  });

  return (
    <NoSsr>
      <div className="custom-toolbar-select">
        <Tooltip title="Deselect ALL">
          <StyledIconButton onClick={handleClickDeselectAll}>
            <StyledIcon>
              <IndeterminateCheckBoxIcon />
            </StyledIcon>
          </StyledIconButton>
        </Tooltip>
        {fullData.length === 1 && (
          <Tooltip title="Download">
            <StyledIconButton
              key="download"
              aria-label="download"
              color="inherit"
              href={`/api/perf/profile/result/${encodeURIComponent(fullData[0].meshery_id)}`}
              download={`${fullData[0].name}_test_result.json`}
            >
              <StyledIcon>
                <GetAppIcon />
              </StyledIcon>
            </StyledIconButton>
          </Tooltip>
        )}
        <Tooltip title="Compare selected">
          <StyledIconButton onClick={handleCompareSelected}>
            <StyledIcon>
              <CompareArrowsIcon />
            </StyledIcon>
          </StyledIconButton>
        </Tooltip>
      </div>
      <MesheryChartDialog
        handleClose={handleDialogClose}
        open={dialogOpen}
        content={
          <div>
            <MesheryChart data={data} />
          </div>
        }
      />
    </NoSsr>
  );
}

const mapDispatchToProps = (dispatch) => ({
  clearResultsSelection: bindActionCreators(clearResultsSelection, dispatch),
});

const mapStateToProps = (state) => {
  const results_selection = state.get('results_selection').toObject();
  return { results_selection };
};

export default connect(mapStateToProps, mapDispatchToProps)(CustomToolbarSelect);
