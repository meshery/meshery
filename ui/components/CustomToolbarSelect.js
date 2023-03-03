import React from 'react';
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

const defaultToolbarSelectStyles = { iconButton : {
  marginRight : '24px',
  top : '50%',
  display : 'inline-block',
  position : 'relative',
},
icon : { color : '#000', },
inverseIcon : { transform : 'rotate(90deg)', }, };

class CustomToolbarSelect extends React.Component {
//   handleClickInverseSelection = () => {
//     const nextSelectedRows = this.props.displayData.reduce((nextSelectedRows, _, index) => {
//       if (!this.props.selectedRows.data.find(selectedRow => selectedRow.index === index)) {
//         nextSelectedRows.push(index);
//       }

  //       return nextSelectedRows;
  //     }, []);

  //     this.props.setSelectedRows(nextSelectedRows);
  //   };
  state = { dialogOpen : false,
    data : [],
    chartCompare : [], // will persist start, end times, chart config and metrics for each result to be compared
  }

  handleDialogClose = () => {
    this.setState({ dialogOpen : false });
  }

  handleDialogOpen = () => {
    this.setState({ dialogOpen : true });
  }

  handleClickDeselectAll = () => {
    this.props.setSelectedRows([]);
    this.props.clearResultsSelection();
  };

  handleCompareSelected = () => {
    //   console.log(`selected rows: ${JSON.stringify(this.props.selectedRows.data)}`);
    //   const self = this;
    const data = [];
    const chartCompare = [];
    //   this.props.selectedRows.data.map(({dataIndex}) => {
    //     // console.log(`data for selected rows: ${JSON.stringify(self.props.results[dataIndex])}`);
    //     data.push(self.props.results[dataIndex]);
    //   });
    //   this.setState({data, dialogOpen: true});
    const rs = this.props.results_selection;
    Object.keys(rs).map((k1) => {
      Object.keys(rs[k1]).map((k2) => {
        if (typeof rs[k1][k2] !== 'undefined') {
          data.push(rs[k1][k2].runner_results);

          const row = rs[k1][k2].runner_results;
          const startTime = new Date(row.StartTime);
          const endTime = new Date(startTime.getTime() + row.ActualDuration / 1000000);
          const boardConfig = rs[k1][k2].server_board_config;
          const serverMetrics = rs[k1][k2].server_metrics;

          chartCompare.push({
            label : row.Labels,
            startTime,
            endTime,
            boardConfig,
            serverMetrics,
          });
        }
      });
    });
    this.setState({ data, chartCompare, dialogOpen : true });
    // console.log(`block users with dataIndexes: ${this.props.selectedRows.data.map(row => row.dataIndex)}`);
  };

  render() {
    const { classes, results_selection } = this.props;
    const fullData = [];
    const rs = results_selection;
    Object.keys(rs).map((k1) => {
      Object.keys(rs[k1]).map((k2) => {
        if (typeof rs[k1][k2] !== 'undefined') {
          fullData.push(rs[k1][k2]);
        }
      });
    });
    return (
      <NoSsr>
        <div className="custom-toolbar-select">
          <Tooltip title="Deselect ALL">
            <IconButton className={classes.iconButton} onClick={this.handleClickDeselectAll}>
              <IndeterminateCheckBoxIcon className={classes.icon} />
            </IconButton>
          </Tooltip>
          {/* <Tooltip title={"Inverse selection"}>
          <IconButton className={classes.iconButton} onClick={this.handleClickInverseSelection}>
            <CompareArrowsIcon className={[classes.icon, classes.inverseIcon].join(" ")} />
          </IconButton>
        </Tooltip> */}
          { fullData.length == 1 && (
            <Tooltip title="Download">
              <IconButton
                className={classes.iconButton}
                key="download"
                aria-label="download"
                color="inherit"
                // onClick={() => self.props.closeSnackbar(key) }
                href={`/api/perf/profile/result/${encodeURIComponent(fullData[0].meshery_id)}`}
                download={`${fullData[0].name}_test_result.json`}
              >
                <GetAppIcon className={classes.icon} />
              </IconButton>
            </Tooltip>
          )}
          <Tooltip title="Compare selected">
            <IconButton className={classes.iconButton} onClick={this.handleCompareSelected}>
              <CompareArrowsIcon className={classes.icon} />
            </IconButton>
          </Tooltip>
        </div>

        <MesheryChartDialog
          handleClose={this.handleDialogClose}
          open={this.state.dialogOpen}
          content={(
            <div>
              <MesheryChart data={this.state.data} />
            </div>
          )}
        />
      </NoSsr>
    );
  }
}


const mapDispatchToProps = (dispatch) => ({ clearResultsSelection : bindActionCreators(clearResultsSelection, dispatch), });

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
