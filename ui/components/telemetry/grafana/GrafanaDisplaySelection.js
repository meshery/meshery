import React from 'react';
import PropTypes from 'prop-types';
import { NoSsr } from '@layer5/sistent';
import { Chip, Box, styled } from '@layer5/sistent';
import MUIDataTable from 'mui-datatables';
import { UsesSistent } from '@/components/SistentWrapper';

const Root = styled(Box)(({ theme }) => ({
  padding: theme.spacing(5),
}));

const PanelChips = styled(Box)(() => ({
  display: 'flex',
  flexWrap: 'wrap',
}));

const StyledChip = styled(Chip)(({ theme }) => ({
  margin: theme.spacing(0.25),
}));

const GrafanaDisplaySelection = ({ boardPanelConfigs, deleteSelectedBoardPanelConfig }) => {
  const selectedValsForDisplay = boardPanelConfigs.map((cf) => ({
    board: cf.board?.title || '',
    panels: (
      <PanelChips>
        {cf.panels.map((panel, ind) => (
          <StyledChip key={`${panel.id}_-_${ind}`} label={panel.title} />
        ))}
      </PanelChips>
    ),
    template_variables: (
      <PanelChips>
        {cf.templateVars
          ? cf.templateVars.map((tv, ind) =>
              tv && tv !== '' ? <StyledChip key={`${tv}-_-${ind}`} label={tv} /> : null,
            )
          : []}
      </PanelChips>
    ),
  }));

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
    <UsesSistent>
      <NoSsr>
        <Root>
          <MUIDataTable
            key={`gds_${new Date().getTime()}`}
            title="Meshery Results"
            data={selectedValsForDisplay}
            columns={columns}
            options={options}
          />
        </Root>
      </NoSsr>
    </UsesSistent>
  );
};

GrafanaDisplaySelection.propTypes = {
  classes: PropTypes.object.isRequired,
  boardPanelConfigs: PropTypes.array.isRequired,
  deleteSelectedBoardPanelConfig: PropTypes.func.isRequired,
};

export default GrafanaDisplaySelection;
