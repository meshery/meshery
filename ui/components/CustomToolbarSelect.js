import React from "react";
import IconButton from "@material-ui/core/IconButton";
import Tooltip from "@material-ui/core/Tooltip";
// import CompareArrowsIcon from "@material-ui/icons/CompareArrows";
import IndeterminateCheckBoxIcon from "@material-ui/icons/IndeterminateCheckBox";
import WavesIcon from "@material-ui/icons/Waves";
import { withStyles } from "@material-ui/core/styles";
import { NoSsr } from "@material-ui/core";
import MesheryChartDialog from "./MesheryChartDialog";
import MesheryChart from "./MesheryChart";
import {connect} from "react-redux";
import { bindActionCreators } from "redux";
import { updateResultsSelection, clearResultsSelection } from "../lib/store";

const defaultToolbarSelectStyles = {
  iconButton: {
    marginRight: "24px",
    top: "50%",
    display: "inline-block",
    position: "relative",
    transform: "translateY(-50%)",
  },
  icon: {
    color: "#000",
  },
  inverseIcon: {
    transform: "rotate(90deg)",
  },
};

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
    state = {
        dialogOpen: false,
        data: [],
    }

    handleDialogClose = () => {
        this.setState({dialogOpen: false});
    }

    handleDialogOpen = () => {
        this.setState({dialogOpen: true});
    }

  handleClickDeselectAll = () => {
    this.props.setSelectedRows([]);
    this.props.clearResultsSelection();
  };

  handleCompareSelected = () => {
    //   console.log(`selected rows: ${JSON.stringify(this.props.selectedRows.data)}`);
    //   const self = this;
      let data = [];
    //   this.props.selectedRows.data.map(({dataIndex}) => {
    //     // console.log(`data for selected rows: ${JSON.stringify(self.props.results[dataIndex])}`);
    //     data.push(self.props.results[dataIndex]);
    //   });
    //   this.setState({data, dialogOpen: true});
    const rs = this.props.results_selection;
    Object.keys(rs).map((k1) => {
        Object.keys(rs[k1]).map((k2) => {
          if (typeof rs[k1][k2] !== 'undefined'){
            data.push(rs[k1][k2]);
          }
        });
    })
    this.setState({data, dialogOpen: true});
    // console.log(`block users with dataIndexes: ${this.props.selectedRows.data.map(row => row.dataIndex)}`);
  };

  render() {
    const { classes } = this.props;

    return (
        <NoSsr>
      <div className={"custom-toolbar-select"}>
        <Tooltip title={"Deselect ALL"}>
          <IconButton className={classes.iconButton} onClick={this.handleClickDeselectAll}>
            <IndeterminateCheckBoxIcon className={classes.icon} />
          </IconButton>
        </Tooltip>
        {/* <Tooltip title={"Inverse selection"}>
          <IconButton className={classes.iconButton} onClick={this.handleClickInverseSelection}>
            <CompareArrowsIcon className={[classes.icon, classes.inverseIcon].join(" ")} />
          </IconButton>
        </Tooltip> */}
        <Tooltip title={"Compare selected"}>
          <IconButton className={classes.iconButton} onClick={this.handleCompareSelected}>
            <WavesIcon className={classes.icon} />
          </IconButton>
        </Tooltip>
      </div>

      <MesheryChartDialog handleClose={this.handleDialogClose} open={this.state.dialogOpen} content={
          <MesheryChart data={this.state.data} />
      } />
      </NoSsr>
    );
  }
}



const mapDispatchToProps = dispatch => {
    return {
        updateResultsSelection: bindActionCreators(updateResultsSelection, dispatch),
        clearResultsSelection: bindActionCreators(clearResultsSelection, dispatch),
    }
}
  
const mapStateToProps = state => {
    const results_selection = state.get("results_selection").toObject();
    return {results_selection};
}

export default withStyles(defaultToolbarSelectStyles, { name: "CustomToolbarSelect" })(
    connect(
        mapStateToProps,
        mapDispatchToProps
      )(CustomToolbarSelect));