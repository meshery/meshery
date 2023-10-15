import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { withStyles } from '@material-ui/core/styles';
import { Chip, NoSsr } from '@material-ui/core';
import MUIDataTable from 'mui-datatables';

const grafanaStyles = (theme) => ({
  root: { padding: theme.spacing(5) },
  panelChips: { display: 'flex', flexWrap: 'wrap' },
  panelChip: { margin: theme.spacing(0.25) },
});

class GrafanaDisplaySelection extends Component {
  render() {
    const { classes, boardPanelConfigs, deleteSelectedBoardPanelConfig } = this.props;
    const selectedValsForDisplay = [];
    boardPanelConfigs.forEach((cf) => {
      selectedValsForDisplay.push({
        board: cf.board && cf.board.title ? cf.board.title : '',
        panels: cf.panels.map((panel, ind) => (
          <Chip key={`${panel.id}_-_${ind}`} label={panel.title} className={classes.panelChip} />
        )),
        template_variables: cf.templateVars
          ? cf.templateVars.map((tv, ind) => {
              if (tv && tv !== '') {
                return <Chip key={`${tv}-_-${ind}`} label={tv} className={classes.panelChip} />;
              }
              return null;
            })
          : [],
      });
    });

    const columns = [
      { name: 'board', label: 'Board' },
      { name: 'panels', label: 'Panels' },
      { name: 'template_variables', label: 'Template Variables' },
    ];
    const options = {
      filter: false,
      sort: false,
      search: false,
      filterType: 'textField',
      responsive: 'stacked',
      count: selectedValsForDisplay.length,
      print: false,
      download: false,
      pagination: false,
      viewColumns: false,
      onRowsDelete: (rowsDeleted) => {
        const delRows = rowsDeleted.data.map(({ dataIndex }) => dataIndex);
        deleteSelectedBoardPanelConfig(delRows);
        return false;
      },
    };
    return (
      <NoSsr>
        <MUIDataTable
          key={`gds_${new Date().getTime()}`}
          title="Meshery Results"
          data={selectedValsForDisplay}
          columns={columns}
          options={options}
        />
      </NoSsr>
    );
  }
}

GrafanaDisplaySelection.propTypes = {
  classes: PropTypes.object.isRequired,
  boardPanelConfigs: PropTypes.array.isRequired,
  deleteSelectedBoardPanelConfig: PropTypes.func.isRequired,
};

export default withStyles(grafanaStyles)(GrafanaDisplaySelection);
